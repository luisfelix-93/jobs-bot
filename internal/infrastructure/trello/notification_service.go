package trello

import (
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"net/url"
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

func (t *TrelloNotifier) Notify(job domain.Job) error {
	apiURL := "https://api.trello.com/1/cards"

	cardDesc := fmt.Sprintf("Link da Vaga:\n&%s", job.Link)
	data := url.Values{}
	data.Set("key", t.apiKey)
	data.Set("token", t.token)
	data.Set("idList", t.listID)
	data.Set("name", job.Title)
	data.Set("desc", cardDesc)

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