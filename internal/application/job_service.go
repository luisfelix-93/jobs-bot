package application

import (
	"jobs-bot/internal/domain"
	"log"
)

type JobService struct {
	repo          domain.JobRepository
	notifier      domain.NotificationService
	filter        *domain.JobFilter
	analyzer      *domain.ResumeAnalyzer
	resumeContent string
	limit         int
}

func NewJobService(repo domain.JobRepository, notifier domain.NotificationService, filter *domain.JobFilter, analyzer *domain.ResumeAnalyzer, resumeContent string, limit int) *JobService {
	return &JobService{
		repo:          repo,
		notifier:      notifier,
		filter:        filter,
		analyzer:      analyzer,
		resumeContent: resumeContent,
		limit:         limit,
	}
}


func (s *JobService) ProcessNewJobs() error {
	log.Println("Iniciando busca por novas vagas...")
	allJobs, err := s.repo.FetchJobs()
	if err != nil {
		return err
	}
	log.Printf("Encontradas %d vagas no feed.", len(allJobs))

	bestJobs := s.filter.FilterAndRankJobs(allJobs, s.limit)
	log.Printf("Após filtragem, %d vagas foram selecionadas para notificação.", len(bestJobs))

	if len(bestJobs) == 0 {
		log.Printf("Nenhuma vaga corresponde aos critérios. Encerrando.")
		return nil
	}

	for _, job := range bestJobs {
		analysis := s.analyzer.Analyze(s.resumeContent, job.Description, s.filter.PositiveKeywords)
		log.Printf("Análise para '%s': %.2f%% de compatibilidade.", job.Title, analysis.MatchPercentage)
		log.Printf("Enviando vaga para o trello: %s", job.Title)
		if err := s.notifier.Notify(job, analysis); err != nil {
			log.Printf("Erro ao enviar vaga para o trello: %v", err)
		}
	}
	log.Println("Processo concluído com sucesso.")
	return nil
}