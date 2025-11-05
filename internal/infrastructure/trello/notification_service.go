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

func (t *TrelloNotifier) Notify(job domain.Job, analysis domain.ResumeAnalysis) error {
	apiURL := "https://api.trello.com/1/cards"

	cardName := fmt.Sprintf("[%s] %s", job.SourceFeed, job.Title)

	cleanDescription := htmlTagRegex.ReplaceAllString(job.FullDescription, "")

	analysisDetails := fmt.Sprintf("%+v", analysis)

	cardDesc := fmt.Sprintf(
		"**ORIGEM:** %s\n\n**LINK DA VAGA:**\n%s\n\n---\n\n**ANÁLISE DO CURRÍCULO:**\n%s\n\n---\n\n**DESCRIÇÃO DA VAGA:**\n%s",
		job.SourceFeed,
		job.Link,
		analysisDetails,  
		cleanDescription,
	)
	

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