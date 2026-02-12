package main

import (
	"log"
	"os"

	"jobs-bot/config"
	"jobs-bot/internal/application"
	"jobs-bot/internal/domain"
	deepseekai "jobs-bot/internal/infrastructure/deepseek"
	"jobs-bot/internal/infrastructure/email"
	"jobs-bot/internal/infrastructure/findwork"
	"jobs-bot/internal/infrastructure/jobicy"
	"jobs-bot/internal/infrastructure/jsearch"
	"jobs-bot/internal/infrastructure/linkedin"
	"jobs-bot/internal/infrastructure/mongodb"
	"jobs-bot/internal/infrastructure/trello"
	"jobs-bot/internal/infrastructure/weworkremotely"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	profiles, err := config.LoadProfiles("profiles.yaml")
	if err != nil {
		log.Fatalf("Erro ao carregar perfis: %v", err)
	}

	store, err := mongodb.NewMongoJobStore(cfg.MongoURI, "jobs-bot")
	if err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	defer store.Close()

	var aiAnalyzer domain.AIAnalyzer
	if cfg.DeepSeekAPIKey != "" {
		aiAnalyzer = deepseekai.NewAnalyzer(cfg.DeepSeekAPIKey)
		log.Println("DeepSeek AI ativado.")
	} else {
		log.Println("AVISO: DeepSeek API Key não encontrada. Funcionalidade de IA desativada (fallback apenas).")
	}

	emailService := email.NewEmailService(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword, cfg.EmailTo)

	log.Printf("MongoDB conectado. Processando %d perfil(is)...", len(profiles))

	var allStats []domain.ProfileStats

	for _, profile := range profiles {
		log.Printf("=== Processando perfil: %s ===", profile.Name)

		resumeBytes, err := os.ReadFile(profile.ResumePath)
		if err != nil {
			log.Printf("[%s] ERRO ao ler currículo '%s': %v", profile.Name, profile.ResumePath, err)
			continue
		}
		resumeContent := string(resumeBytes)

		repos := buildRepos(profile.Sources, cfg)
		if len(repos) == 0 {
			log.Printf("[%s] Nenhuma fonte de vagas configurada. Pulando.", profile.Name)
			continue
		}

		notifier := trello.NewTrelloNotifier(cfg.TrelloAPIKey, cfg.TrelloAPIToken, profile.TrelloListID)
		jobFilter := domain.NewJobFilter(profile.PositiveKeywords, profile.NegativeKeywords)
		resumeAnalyzer := domain.NewResumeAnalyzer()

		appService := application.NewJobService(
			repos,
			notifier,
			jobFilter,
			resumeAnalyzer,
			aiAnalyzer,
			store,
			resumeContent,
			profile.PositiveKeywords,
			profile.Name,
			cfg.JobLimit,
		)

		stats, err := appService.ProcessNewJobs()
		if err != nil {
			log.Printf("[%s] Erro ao processar vagas: %v", profile.Name, err)
		} else {
			allStats = append(allStats, stats)
		}
	}

	log.Println("Todos os perfis processados.")
	log.Println("Enviando email de resumo...")

	if err := emailService.SendSummary(allStats); err != nil {
		log.Printf("ERRO ao enviar email de resumo: %v", err)
	} else {
		log.Println("Email de resumo enviado com sucesso.")
	}
}

func buildRepos(sources config.Sources, cfg *config.Config) []domain.JobRepository {
	repos := make([]domain.JobRepository, 0, 5)

	if sources.JobicyURL != "" {
		log.Println("  + Fonte: Jobicy")
		repos = append(repos, jobicy.NewRssRepository(sources.JobicyURL))
	}
	if sources.WwrURL != "" {
		log.Println("  + Fonte: WeWorkRemotely")
		repos = append(repos, weworkremotely.NewRssRepository(sources.WwrURL))
	}
	if sources.LinkedInURL != "" {
		log.Println("  + Fonte: LinkedIn")
		repos = append(repos, linkedin.NewRssRepository(sources.LinkedInURL))
	}
	if sources.JSearchQuery != "" && cfg.JSearchAPIKey != "" {
		log.Println("  + Fonte: JSearch (RapidAPI)")
		repos = append(repos, jsearch.NewRepository(cfg.JSearchAPIKey, sources.JSearchQuery))
	} else if sources.JSearchQuery != "" {
		log.Println("  - Fonte: JSearch IGNORADA (API Key não configurada)")
	}

	if (sources.FindworkSearch != "" || sources.FindworkLocation != "") && cfg.FindworkAPIKey != "" {
		log.Println("  + Fonte: Findwork.dev")
		repos = append(repos, findwork.NewRepository(cfg.FindworkAPIKey, sources.FindworkSearch, sources.FindworkLocation))
	} else if sources.FindworkSearch != "" || sources.FindworkLocation != "" {
		log.Println("  - Fonte: Findwork IGNORADA (API Key não configurada)")
	}

	return repos
}
