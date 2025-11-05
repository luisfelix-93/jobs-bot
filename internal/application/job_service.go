// Arquivo: internal/application/job_service.go
package application

import (
	"jobs-bot/internal/domain"
	"log"
	"sync"
)

// JobService agora contém o analisador e o conteúdo do currículo
type JobService struct {
	repos          []domain.JobRepository
	notifier       domain.NotificationService
	filter         *domain.JobFilter
	analyzer       *domain.ResumeAnalyzer // <-- ADICIONADO
	resumeContent  string                 // <-- ADICIONADO
	filterKeywords []string               // <-- ADICIONADO (para o analyzer)
	limit          int
}

// NewJobService é o construtor que injeta as novas dependências
func NewJobService(
	repos []domain.JobRepository,
	notifier domain.NotificationService,
	filter *domain.JobFilter,
	analyzer *domain.ResumeAnalyzer, // <-- ADICIONADO
	resumeContent string,            // <-- ADICIONADO
	filterKeywords []string,          // <-- ADICIONADO
	limit int,
) *JobService {
	return &JobService{
		repos:          repos,
		notifier:       notifier,
		filter:         filter,
		analyzer:       analyzer,       // <-- ADICIONADO
		resumeContent:  resumeContent,  // <-- ADICIONADO
		filterKeywords: filterKeywords, // <-- ADICIONADO
		limit:          limit,
	}
}

// ProcessNewJobs executa o fluxo completo
func (s *JobService) ProcessNewJobs() error {
	log.Println("Iniciando busca em todas as fontes...")

	var allJobs []domain.Job
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, repo := range s.repos {
		wg.Add(1)
		go func(r domain.JobRepository) {
			defer wg.Done()
			jobs, err := r.FetchJobs()
			if err != nil {
				log.Printf("ERRO ao buscar em um repositório: %v", err)
				return
			}
			mu.Lock()
			allJobs = append(allJobs, jobs...)
			mu.Unlock()
		}(repo)
	}
	wg.Wait()

	log.Printf("Encontradas %d vagas no total. Filtrando...", len(allJobs))

	// PASSO 1: O filtro encontra as vagas mais relevantes
	bestJobs := s.filter.FilterAndRankJobs(allJobs, s.limit)
	log.Printf("Após filtragem, %d vagas foram selecionadas para notificação.", len(bestJobs))

	if len(bestJobs) == 0 {
		log.Println("Nenhuma vaga nova corresponde aos critérios. Encerrando.")
		return nil
	}

	// PASSO 2: Iterar, ANALISAR o currículo, e NOTIFICAR
	for _, job := range bestJobs {
		log.Printf("Analisando vaga: %s", job.Title)

		// Executa a análise do currículo contra a descrição da vaga
		analysis := s.analyzer.Analyze(s.resumeContent, job.FullDescription, s.filterKeywords)

		log.Printf("Enviando vaga para o Trello: %s (Match: %.2f%%)", job.Title, analysis.MatchPercentage)

		// Finalmente, chama o Notify com os DOIS argumentos
		if err := s.notifier.Notify(job, analysis); err != nil {
			log.Printf("ERRO ao notificar sobre a vaga '%s': %v", job.Title, err)
		}
	}

	log.Println("Processo concluído com sucesso.")
	return nil
}