package application

import (
	"fmt"
	"jobs-bot/internal/domain"
	"jobs-bot/internal/domain/normalization"
	"log"
	"sort"
	"sync"
	"time"
)

const aiScoreThreshold = 50

type JobService struct {
	repos          []domain.JobRepository
	notifier       domain.NotificationService
	filter         *domain.JobFilter
	analyzer       *domain.ResumeAnalyzer
	aiAnalyzer     domain.AIAnalyzer
	store          domain.JobStore
	pipeline       *normalization.Pipeline
	resumeContent  string
	filterKeywords []string
	profileName    string
	limit          int
}

func NewJobService(
	repos []domain.JobRepository,
	notifier domain.NotificationService,
	filter *domain.JobFilter,
	analyzer *domain.ResumeAnalyzer,
	aiAnalyzer domain.AIAnalyzer,
	store domain.JobStore,
	pipeline *normalization.Pipeline,
	resumeContent string,
	filterKeywords []string,
	profileName string,
	limit int,
) *JobService {
	return &JobService{
		repos:          repos,
		notifier:       notifier,
		filter:         filter,
		analyzer:       analyzer,
		aiAnalyzer:     aiAnalyzer,
		store:          store,
		pipeline:       pipeline,
		resumeContent:  resumeContent,
		filterKeywords: filterKeywords,
		profileName:    profileName,
		limit:          limit,
	}
}

func (s *JobService) ProcessNewJobs() (domain.ProfileStats, error) {
	stats := domain.ProfileStats{ProfileName: s.profileName}
	log.Printf("[%s] Iniciando busca em todas as fontes...", s.profileName)

	var allJobs []domain.Job
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, repo := range s.repos {
		wg.Add(1)
		go func(r domain.JobRepository) {
			defer wg.Done()
			jobs, err := r.FetchJobs()
			if err != nil {
				log.Printf("[%s] ERRO ao buscar em um repositório: %v", s.profileName, err)
				return
			}
			mu.Lock()
			allJobs = append(allJobs, jobs...)
			mu.Unlock()
		}(repo)
	}
	wg.Wait()

	stats.TotalFound = len(allJobs)
	log.Printf("[%s] Encontradas %d vagas no total. Normalizando...", s.profileName, len(allJobs))

	if s.pipeline != nil {
		allJobs = s.pipeline.NormalizeAll(allJobs)

		var withSeniority, withWorkMode, withSalary int
		for _, j := range allJobs {
			if j.Seniority != "" {
				withSeniority++
			}
			if j.WorkMode != "" {
				withWorkMode++
			}
			if j.SalaryMin > 0 {
				withSalary++
			}
		}
		log.Printf("[%s] Normalização concluída. Stats - Seniority: %d/%d, WorkMode: %d/%d, Salary: %d/%d",
			s.profileName, withSeniority, len(allJobs), withWorkMode, len(allJobs), withSalary, len(allJobs))
	}

	rankedJobs := s.filter.FilterAndRankJobs(allJobs)
	stats.TotalFiltered = len(rankedJobs)
	log.Printf("[%s] Após filtragem, %d vagas elegíveis.", s.profileName, len(rankedJobs))

	if len(rankedJobs) == 0 {
		log.Printf("[%s] Nenhuma vaga corresponde aos critérios.", s.profileName)
		return stats, nil
	}

	var notifiedJobs []domain.ProcessedJob
	newJobsProcessed := 0

	for _, job := range rankedJobs {
		if s.limit > 0 && newJobsProcessed >= s.limit {
			break
		}

		guid := fmt.Sprintf("%s-%s", job.SourceFeed, job.GUID)

		exists, err := s.store.Exists(guid, s.profileName)
		if err != nil {
			log.Printf("[%s] ERRO ao verificar dedup para '%s': %v", s.profileName, job.Title, err)
			continue
		}
		if exists {
			stats.TotalSkipped++
			continue
		}

		keywordAnalysis := s.analyzer.Analyze(s.resumeContent, job.FullDescription, s.filterKeywords)

		aiAnalysis := s.analyzeWithAI(job)

		shouldNotify := true
		if aiAnalysis != nil && aiAnalysis.Score < aiScoreThreshold {
			shouldNotify = false
			stats.BelowThreshold++
			log.Printf("[%s] Vaga '%s' abaixo do threshold (score: %d). Salvando sem notificar.", s.profileName, job.Title, aiAnalysis.Score)
		}

		if shouldNotify {
			log.Printf("[%s] Enviando vaga: %s (AI Score: %s, Keyword Match: %.2f%%)",
				s.profileName, job.Title, formatAIScore(aiAnalysis), keywordAnalysis.MatchPercentage)

			if err := s.notifier.Notify(job, keywordAnalysis, aiAnalysis); err != nil {
				log.Printf("[%s] ERRO ao notificar '%s': %v", s.profileName, job.Title, err)
			}
			stats.TotalNotified++
		}

		processedJob := domain.ProcessedJob{
			GUID:            guid,
			Source:          job.SourceFeed,
			Profile:         s.profileName,
			Title:           job.Title,
			Link:            job.Link,
			Location:        job.Location,
			Description:     job.FullDescription,
			KeywordAnalysis: keywordAnalysis,
			AIAnalysis:      aiAnalysis,
			Notified:        shouldNotify,
			NotifiedAt:      time.Now(),
			CreatedAt:       time.Now(),
			TTLExpireAt:     time.Now().Add(90 * 24 * time.Hour),
			Company:         job.Company,
			Seniority:       job.Seniority,
			WorkMode:        job.WorkMode,
			EmploymentType:  job.EmploymentType,
			Skills:          job.Skills,
			SalaryMin:       job.SalaryMin,
			SalaryMax:       job.SalaryMax,
			SalaryCurrency:  job.SalaryCurrency,
			NormalizedTitle: job.NormalizedTitle,
		}

		if err := s.store.Save(processedJob); err != nil {
			log.Printf("[%s] ERRO ao salvar job '%s' no MongoDB: %v", s.profileName, job.Title, err)
		}

		if shouldNotify {
			notifiedJobs = append(notifiedJobs, processedJob)
		}

		newJobsProcessed++
	}

	// Sort notified jobs by score (AI or MatchPercentage) to pick top ones
	sort.Slice(notifiedJobs, func(i, j int) bool {
		scoreI := getScore(notifiedJobs[i])
		scoreJ := getScore(notifiedJobs[j])
		return scoreI > scoreJ
	})

	if len(notifiedJobs) > 5 {
		stats.TopJobs = notifiedJobs[:5]
	} else {
		stats.TopJobs = notifiedJobs
	}

	log.Printf("[%s] Concluído: %d novas processadas, %d notificadas, %d duplicadas, %d abaixo do threshold.",
		s.profileName, newJobsProcessed, stats.TotalNotified, stats.TotalSkipped, stats.BelowThreshold)
	return stats, nil
}

func getScore(job domain.ProcessedJob) float64 {
	if job.AIAnalysis != nil {
		return float64(job.AIAnalysis.Score)
	}
	return job.KeywordAnalysis.MatchPercentage
}

func (s *JobService) analyzeWithAI(job domain.Job) *domain.AIAnalysis {
	if s.aiAnalyzer == nil {
		return nil
	}

	analysis, err := s.aiAnalyzer.Analyze(s.resumeContent, job.FullDescription)
	if err != nil {
		log.Printf("[%s] DeepSeek falhou para '%s': %v. Usando fallback keyword.", s.profileName, job.Title, err)

		keywordResult := s.analyzer.Analyze(s.resumeContent, job.FullDescription, s.filterKeywords)
		return &domain.AIAnalysis{
			Score:          int(keywordResult.MatchPercentage),
			Strengths:      keywordResult.FoundKeywords,
			Gaps:           keywordResult.MissingKeywords,
			Recommendation: classifyByScore(int(keywordResult.MatchPercentage)),
			Summary:        "Análise por keyword matching (fallback - DeepSeek indisponível).",
			Source:         "keyword_fallback",
		}
	}

	return analysis
}

func classifyByScore(score int) string {
	if score >= 70 {
		return "apply"
	}
	if score >= 50 {
		return "review"
	}
	return "skip"
}

func formatAIScore(ai *domain.AIAnalysis) string {
	if ai == nil {
		return "N/A"
	}
	return fmt.Sprintf("%d (%s)", ai.Score, ai.Source)
}
