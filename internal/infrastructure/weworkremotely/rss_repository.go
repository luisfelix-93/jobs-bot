package weworkremotely

import (
	"encoding/xml"
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"strings"
	"time"
)

type wwrItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"` // Contém o HTML da descrição [cite: 3, 15]
	Region      string `xml:"region"`      // Ex: "Anywhere in the World", "California" [cite: 2, 13]
	Country     string `xml:"country"`     // Ex: "🇺🇸 United States of America" [cite: 2]
}

type wwrRss struct {
	Channel struct {
		Items []wwrItem `xml:"item"`
	} `xml:"channel"`
}

type RssRepository struct {
	rssURL string
	client *http.Client
}

func NewRssRepository(rssURL string) *RssRepository {
	return &RssRepository{
		rssURL: rssURL,
		client: &http.Client{Timeout: 15 *time.Second},
	}
}

func (r *RssRepository) FetchJobs() ([]domain.Job, error) {

	// 2. CRIAMOS UMA REQUISIÇÃO MANUAL
	req, err := http.NewRequest("GET", r.rssURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request para WWR: %w", err)
	}

	// 3. DEFINIMOS O USER-AGENT PARA SIMULAR UM NAVEGADOR
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.0.0 Safari/537.36")

	// 4. EXECUTAMOS A REQUISIÇÃO COM NOSSO CLIENT (que tem o timeout)
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar feed WWR: %w", err)
	}
	defer resp.Body.Close()

	// 5. O RESTO DA FUNÇÃO (PARSE DO XML) CONTINUA EXATAMENTE IGUAL
	var feed wwrRss
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("erro ao decodificar XML do WWR: %w", err)
	}

	var jobs []domain.Job
	for _, item := range feed.Channel.Items {
		jobs = append(jobs, domain.Job{
			Title:           item.Title,
			Link:            item.Link,
			GUID:            item.GUID,
			SourceFeed:      "WWR",
			Location:        strings.TrimSpace(item.Country + " " + item.Region),
			FullDescription: item.Description,
		})
	}
	return jobs, nil
}