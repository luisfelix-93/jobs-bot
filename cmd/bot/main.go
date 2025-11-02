package main

import (
	"io/ioutil"
	"jobs-bot/config"
	"jobs-bot/internal/application"
	"jobs-bot/internal/domain"
	"jobs-bot/internal/infrastructure/linkedin"
	"jobs-bot/internal/infrastructure/trello"
	"log"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar as configurações: %v", err)
	}
	
	resumeBytes, err := ioutil.ReadFile(cfg.ResumePath)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo de currículo: %v", err)
	}

	resumeContent := string(resumeBytes)
	linkedinInRepo := linkedin.NewRssRepository(cfg.LinkedInRssURL)
	trelloNotifier := trello.NewTrelloNotifier(cfg.TrelloAPIKey, cfg.TrelloAPIToken, cfg.TrelloListID)
	jobFilter := domain.NewJobFilter(cfg.PositiveKeywords, cfg.NegativeKeywords)
	resumeAnalyzer := domain.NewResumeAnalyzer() // NOVO
	
	appService := application.NewJobService(
		linkedinInRepo,
		trelloNotifier,
		jobFilter,
		resumeAnalyzer, // NOVO
		resumeContent,  // NOVO
		cfg.JobLimit,
	)

	if err := appService.ProcessNewJobs(); err != nil {
		log.Fatalf("O bot encontrou um erro fatal: %v", err)
	}

}