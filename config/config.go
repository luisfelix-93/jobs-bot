package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Config struct {
	TrelloAPIKey   string
	TrelloAPIToken string
	MongoURI       string
	JobLimit       int
	DeepSeekAPIKey string
	JSearchAPIKey  string
	FindworkAPIKey string
	// Email Config
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	EmailTo      string
}

type ProfileConfig struct {
	Name             string   `yaml:"name"`
	ResumePath       string   `yaml:"resume_path"`
	PositiveKeywords []string `yaml:"positive_keywords"`
	NegativeKeywords []string `yaml:"negative_keywords"`
	TrelloListID     string   `yaml:"trello_list_id"`
	Sources          Sources  `yaml:"sources"`
}

type Sources struct {
	JobicyURL        string `yaml:"jobicy_url"`
	WwrURL           string `yaml:"wwr_url"`
	LinkedInURL      string `yaml:"linkedin_url"`
	JSearchQuery     string `yaml:"jsearch_query"`
	FindworkSearch   string `yaml:"findwork_search"`
	FindworkLocation string `yaml:"findwork_location"`
}

type profilesFile struct {
	Profiles []ProfileConfig `yaml:"profiles"`
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

	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort := 587
	if smtpPortStr != "" {
		p, err := strconv.Atoi(smtpPortStr)
		if err == nil {
			smtpPort = p
		}
	}

	cfg := &Config{
		TrelloAPIKey:   getEnv("TRELLO_API_KEY", ""),
		TrelloAPIToken: getEnv("TRELLO_API_TOKEN", ""),
		MongoURI:       getEnv("MONGO_URI", ""),
		JobLimit:       jobLimit,
		DeepSeekAPIKey: os.Getenv("DEEPSEEK_API_KEY"),
		JSearchAPIKey:  os.Getenv("JSEARCH_API_KEY"),
		FindworkAPIKey: os.Getenv("FINDWORK_API_KEY"),

		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     smtpPort,
		SMTPUser:     os.Getenv("SMTP_USER"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		EmailTo:      os.Getenv("EMAIL_TO"),
	}

	return cfg, nil
}

func LoadProfiles(path string) ([]ProfileConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de perfis '%s': %w", path, err)
	}

	var pf profilesFile
	if err := yaml.Unmarshal(data, &pf); err != nil {
		return nil, fmt.Errorf("erro ao parsear YAML de perfis: %w", err)
	}

	if len(pf.Profiles) == 0 {
		return nil, fmt.Errorf("nenhum perfil encontrado em '%s'", path)
	}

	return pf.Profiles, nil
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
