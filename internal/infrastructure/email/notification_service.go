package email

import (
	"fmt"
	"jobs-bot/internal/domain"
	"log"
	"net/smtp"
	"time"
)

type EmailService struct {
	host     string
	port     int
	user     string
	password string
	to       string
}

func NewEmailService(host string, port int, user, password, to string) *EmailService {
	return &EmailService{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		to:       to,
	}
}

func (s *EmailService) SendSummary(stats []domain.ProfileStats) error {
	if s.host == "" || s.to == "" {
		log.Println("AVISO: Envio de e-mail de resumo pulado (SMTP_HOST ou EMAIL_TO não configurado).")
		return nil
	}

	subject := fmt.Sprintf("Subject: Resumo Diario de Vagas - %s\n", time.Now().Format("02/01/2006"))
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body := "<html><body><h2>📊 Resumo Diário de Vagas</h2>"

	for _, stat := range stats {
		body += fmt.Sprintf("<h3>👤 Perfil: %s</h3>", stat.ProfileName)
		body += "<ul>"
		body += fmt.Sprintf("<li>🔍 <b>Encontradas:</b> %d</li>", stat.TotalFound)
		body += fmt.Sprintf("<li>✅ <b>Selecionadas (Filtro):</b> %d</li>", stat.TotalFiltered)
		body += fmt.Sprintf("<li>📢 <b>Notificadas:</b> %d</li>", stat.TotalNotified)
		body += fmt.Sprintf("<li>🚫 <b>Abaixo do Threshold (<50):</b> %d</li>", stat.BelowThreshold)
		body += fmt.Sprintf("<li>🔄 <b>Duplicadas (Ignoradas):</b> %d</li>", stat.TotalSkipped)
		body += "</ul>"

		if len(stat.TopJobs) > 0 {
			body += "<h4>🚀 Top Vagas Notificadas:</h4><ul>"
			for _, job := range stat.TopJobs {
				score := "N/A"
				if job.AIAnalysis != nil {
					score = fmt.Sprintf("%d", job.AIAnalysis.Score)
				} else {
					score = fmt.Sprintf("%.0f%%", job.KeywordAnalysis.MatchPercentage)
				}
				body += fmt.Sprintf("<li><b>[%s]</b> <a href=\"%s\">%s</a> (Score: %s)</li>",
					job.Source, job.Link, job.Title, score)
			}
			body += "</ul>"
		}
		body += "<hr>"
	}

	body += "<small>Gerado automaticamente por Jobs Bot 🤖</small></body></html>"

	msg := []byte(subject + mime + body)
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	if err := smtp.SendMail(addr, auth, s.user, []string{s.to}, msg); err != nil {
		return fmt.Errorf("erro ao enviar email: %w", err)
	}

	return nil
}
