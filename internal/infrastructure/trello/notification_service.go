package trello

import (
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TrelloNotifier struct {
	apiKey string
	token  string
	listID string
	client *http.Client
}

func NewTrelloNotifier(apiKey, token, listID string) *TrelloNotifier {
	return &TrelloNotifier{
		apiKey: apiKey,
		token:  token,
		listID: listID,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (t *TrelloNotifier) Notify(job domain.Job, analysis domain.ResumeAnalysis) error {
	apiURL := "https://api.trello.com/1/cards"

	// Crie uma descrição muito mais rica!
	cardDesc := fmt.Sprintf(
		"**Link da Vaga:**\n%s\n\n---\n\n**Compatibilidade: %.2f%%**\n\n**Palavras-Chave Encontradas:**\n`%s`\n\n**Palavras-Chave Faltando (Adicionar ao CV!):**\n`%s`",
		job.Link,
		analysis.MatchPercentage,
		strings.Join(analysis.FoundKeywords, ", "),
		strings.Join(analysis.MissingKeywords, ", "),
	)

	// Adicione a pontuação ao título do card
	cardName := fmt.Sprintf("[%.f%%] %s", analysis.MatchPercentage, job.Title)

	data := url.Values{}
	data.Set("key", t.apiKey)
	data.Set("token", t.token)
	data.Set("idList", t.listID)
	data.Set("name", cardName) // Título com a pontuação
	data.Set("desc", cardDesc) // Descrição rica com a análise

	resp, err := t.client.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("erro ao enviar notificação para Trello: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status não esperado ao enviar notificação para Trello: %s", resp.Status)
	}

	fmt.Println("Notificação enviada com sucesso para Trello!")
	return nil
}