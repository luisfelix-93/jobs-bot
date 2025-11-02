package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	
	"github.com/joho/godotenv" // Importe a nova biblioteca
)

// Config armazena todas as configurações da aplicação.
type Config struct {
	LinkedInRssURL   string
	TrelloAPIKey     string
	TrelloAPIToken   string
	TrelloListID     string
	PositiveKeywords []string
	NegativeKeywords []string
	ResumePath       string
	JobLimit         int
}

// LoadConfig carrega as configurações das variáveis de ambiente.
func LoadConfig() (*Config, error) {
	// --- ADICIONE ESTA LINHA ---
	// Carrega os valores do arquivo .env para o ambiente.
	// Ignora o erro se o arquivo .env não for encontrado.
	_ = godotenv.Load()
	// --------------------------

	jobLimitStr := getEnv("JOB_LIMIT", "10")
	jobLimit, err := strconv.Atoi(jobLimitStr)
	if err != nil {
		log.Fatalf("JOB_LIMIT inválido: %v", err)
	}

	return &Config{
		LinkedInRssURL:   getEnv("LINKEDIN_RSS_URL", ""),
		TrelloAPIKey:     getEnv("TRELLO_API_KEY", ""),
		TrelloAPIToken:   getEnv("TRELLO_API_TOKEN", ""),
		TrelloListID:     getEnv("TRELLO_LIST_ID", ""),
		PositiveKeywords: strings.Split(getEnv("POSITIVE_KEYWORDS", ""), ","),
		NegativeKeywords: strings.Split(getEnv("NEGATIVE_KEYWORDS", ""), ","),
		ResumePath:       getEnv("RESUME_PATH", ""),
		JobLimit:         jobLimit,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if fallback == "" {
		log.Fatalf("ERRO: Variável de ambiente obrigatória não definida: %s", key)
	}
	return fallback
}