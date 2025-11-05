package trello

import (
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

var htmlTagRegex = regexp.MustCompile("<[^>]*>")

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

// --- ATUALIZADO ---
// A assinatura agora aceita o segundo argumento 'analysis'
// para corresponder à sua interface domain.NotificationService
func (t *TrelloNotifier) Notify(job domain.Job, analysis domain.ResumeAnalysis) error {
	apiURL := "https://api.trello.com/1/cards"

	// Título do Card (com origem)
	cardName := fmt.Sprintf("[%s] %s", job.SourceFeed, job.Title)

	// Limpa a descrição HTML
	cleanDescription := htmlTagRegex.ReplaceAllString(job.FullDescription, "")

	// Formata a análise (de forma genérica)
	// Isso irá imprimir os nomes dos campos e seus valores.
	// Ex: {Score:10 KeywordsFound:[Go, K8s]}
	analysisDetails := fmt.Sprintf("%+v", analysis)

	// Descrição do Card (com todos os novos dados)
	cardDesc := fmt.Sprintf(
		"**ORIGEM:** %s\n\n**LINK DA VAGA:**\n%s\n\n---\n\n**ANÁLISE DO CURRÍCULO:**\n%s\n\n---\n\n**DESCRIÇÃO DA VAGA:**\n%s",
		job.SourceFeed,
		job.Link,
		analysisDetails,  // Adiciona a análise aqui
		cleanDescription,
	)
	// --- FIM DA ATUALIZAÇÃO ---

	data := url.Values{}
	data.Set("key", t.apiKey)
	data.Set("token", t.token)
	data.Set("idList", t.listID)
	data.Set("name", cardName)
	data.Set("desc", cardDesc)

	resp, err := t.client.PostForm(apiURL, data)
	if err != nil {
		return fmt.Errorf("erro ao enviar request para o Trello: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro da API Trello, status: %s, body: %s", resp.Status, resp.Body)
	}

	return nil
}