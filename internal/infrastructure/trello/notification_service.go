package trello

import (
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"net/url"
	"regexp"
	"strings"
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

func (t *TrelloNotifier) Notify(job domain.Job, analysis domain.ResumeAnalysis, aiAnalysis *domain.AIAnalysis) error {
	apiURL := "https://api.trello.com/1/cards"

	cardName := fmt.Sprintf("[%s] %s", job.SourceFeed, job.Title)

	cleanDescription := htmlTagRegex.ReplaceAllString(job.FullDescription, "")

	analysisDetails := fmt.Sprintf("**Match de Palavras-Chave:** %.2f%%\n**Encontradas:** %v\n**Faltantes:** %v",
		analysis.MatchPercentage, analysis.FoundKeywords, analysis.MissingKeywords)

	aiDetails := ""
	if aiAnalysis != nil {
		cardName = fmt.Sprintf("[AI: %d] %s", aiAnalysis.Score, cardName)
		aiDetails = fmt.Sprintf("\n\n---\n\n**ANÁLISE IA (%s - Score: %d)**\n\n**Recomendação:** %s\n**Resumo:** %s\n**Pontos Fortes:**\n- %s\n**Gaps:**\n- %s",
			aiAnalysis.Source,
			aiAnalysis.Score,
			strings.ToUpper(aiAnalysis.Recommendation),
			aiAnalysis.Summary,
			strings.Join(aiAnalysis.Strengths, "\n- "),
			strings.Join(aiAnalysis.Gaps, "\n- "),
		)
	}

	cardDesc := fmt.Sprintf(
		"**ORIGEM:** %s\n\n**LINK DA VAGA:**\n%s\n\n---\n\n**ANÁLISE DE KEYWORDS:**\n%s%s\n\n---\n\n**DESCRIÇÃO DA VAGA:**\n%s",
		job.SourceFeed,
		job.Link,
		analysisDetails,
		aiDetails,
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
