package main

import (
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
	
	linkedinInRepo := linkedin.NewRssRepository(cfg.LinkedInRssURL)
	trelloNotifier := trello.NewTrelloNotifier(cfg.TrelloAPIKey, cfg.TrelloAPIToken, cfg.TrelloListID)
	
	jobFilter := domain.NewJobFilter(cfg.PositiveKeywords, cfg.NegativeKeywords)
	
	appService := application.NewJobService(linkedinInRepo, trelloNotifier, jobFilter, cfg.JobLimit)

	if err := appService.ProcessNewJobs(); err != nil {
		log.Fatalf("O bot encontrou um erro fatal: %v", err)
	}

}