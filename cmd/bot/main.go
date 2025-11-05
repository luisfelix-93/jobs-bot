package main

import (
	"log"
	"os"

	"jobs-bot/config"
	"jobs-bot/internal/application"
	"jobs-bot/internal/domain"
	"jobs-bot/internal/infrastructure/jobicy"
	"jobs-bot/internal/infrastructure/linkedin"
	"jobs-bot/internal/infrastructure/trello"
	"jobs-bot/internal/infrastructure/weworkremotely"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	log.Println("Carregando currículo de:", cfg.ResumeFilePath)
	resumeBytes, err := os.ReadFile(cfg.ResumeFilePath)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo do currículo: %v", err)
	}
	resumeContent := string(resumeBytes)

	trelloNotifier := trello.NewTrelloNotifier(cfg.TrelloAPIKey, cfg.TrelloAPIToken, cfg.TrelloListID)

	allRepos := make([]domain.JobRepository, 0, 3)

	if cfg.JobicyRssURL != "" {
		log.Println("URL do Jobicy encontrada, adicionando à busca.")
		allRepos = append(allRepos, jobicy.NewRssRepository(cfg.JobicyRssURL))
	} else {
		log.Println("URL do Jobicy não configurada, pulando.")
	}

	if cfg.WwrRssURL != "" {
		log.Println("URL do WeWorkRemotely encontrada, adicionando à busca.")
		allRepos = append(allRepos, weworkremotely.NewRssRepository(cfg.WwrRssURL))
	} else {
		log.Println("URL do WeWorkRemotely não configurada, pulando.")
	}

	if cfg.LinkedInRssURL != "" {
		log.Println("URL do LinkedIn encontrada, adicionando à busca.")
		allRepos = append(allRepos, linkedin.NewRssRepository(cfg.LinkedInRssURL))
	} else {
		log.Println("URL do LinkedIn não configurada, pulando.")
	}


	jobFilter := domain.NewJobFilter(cfg.PositiveKeywords, cfg.NegativeKeywords)
	resumeAnalyzer := domain.NewResumeAnalyzer()

	appService := application.NewJobService(
		allRepos,
		trelloNotifier,
		jobFilter,
		resumeAnalyzer,
		resumeContent,
		cfg.PositiveKeywords,
		cfg.JobLimit,
	)

	if err := appService.ProcessNewJobs(); err != nil {
		log.Fatalf("O bot encontrou um erro fatal: %v", err)
	}
}