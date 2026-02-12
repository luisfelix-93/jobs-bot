package application

import (
	"fmt"
	"jobs-bot/internal/domain"
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
	log.Printf("[%s] Encontradas %d vagas no total. Filtrando...", s.profileName, len(allJobs))

	bestJobs := s.filter.FilterAndRankJobs(allJobs, s.limit)
	stats.TotalFiltered = len(bestJobs)
	log.Printf("[%s] Após filtragem, %d vagas selecionadas.", s.profileName, len(bestJobs))

	if len(bestJobs) == 0 {
		log.Printf("[%s] Nenhuma vaga nova corresponde aos critérios.", s.profileName)
		return stats, nil
	}

	var notifiedJobs []domain.ProcessedJob

	for _, job := range bestJobs {
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
		}

		if err := s.store.Save(processedJob); err != nil {
			log.Printf("[%s] ERRO ao salvar job '%s' no MongoDB: %v", s.profileName, job.Title, err)
		}

		if shouldNotify {
			notifiedJobs = append(notifiedJobs, processedJob)
		}
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

	log.Printf("[%s] Concluído: %d notificadas, %d duplicadas, %d abaixo do threshold.",
		s.profileName, stats.TotalNotified, stats.TotalSkipped, stats.BelowThreshold)
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
