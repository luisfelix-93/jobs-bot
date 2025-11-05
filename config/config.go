package config

import (
	"log"
	"os" // Importar 'os'
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	JobicyRssURL     string
	WwrRssURL        string
	LinkedInRssURL   string
	TrelloAPIKey     string
	TrelloAPIToken   string
	TrelloListID     string
	PositiveKeywords []string
	NegativeKeywords []string
	JobLimit         int
	ResumeFilePath   string
}

func LoadConfig() (*Config, error) {
	if os.Getenv("VERCEL_ENV") != "production" {
		_ = godotenv.Load()
	}

	jobLimitStr := getEnv("JOB_LIMIT", "10")
	jobLimit, err := strconv.Atoi(jobLimitStr)
	if err != nil {
		log.Fatalf("JOB_LIMIT inválido: %v", err)
	}

	cfg := &Config{
		// --- VARIÁVEIS OBRIGATÓRIAS ---
		TrelloAPIKey:     getEnv("TRELLO_API_KEY", ""),
		TrelloAPIToken:   getEnv("TRELLO_API_TOKEN", ""),
		TrelloListID:     getEnv("TRELLO_LIST_ID", ""),
		ResumeFilePath:   getEnv("RESUME_FILE_PATH", ""),
		PositiveKeywords: strings.Split(getEnv("POSITIVE_KEYWORDS", ""), ","),
		NegativeKeywords: strings.Split(getEnv("NEGATIVE_KEYWORDS", ""), ","),
		JobLimit:         jobLimit,

		// --- VARIÁVEIS OPCIONAIS (URLs) ---
		// Carrega todas as URLs opcionalmente.
		// Se a string estiver vazia, o 'main.go' vai pular.
		JobicyRssURL:   os.Getenv("JOBICY_RSS_URL"),
		WwrRssURL:      os.Getenv("WWR_RSS_URL"),
		LinkedInRssURL: os.Getenv("LINKEDIN_RSS_URL"),
	}

	return cfg, nil
}

// getEnv é usado apenas para chaves obrigatórias
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "" {
		// Se for uma chave obrigatória e não tiver fallback, vai falhar.
		log.Fatalf("ERRO: Variável de ambiente obrigatória não definida: %s", key)
	}
	return fallback
}