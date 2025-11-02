package linkedin

import (
	"encoding/xml"
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"time"
	"github.com/jaytaylor/html2text"
)

type RssFeed struct {
	Channel struct {
		Items []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			GUID        string `xml:"guid"`
			Description string `xml:"description"`
		} `xml:"item"`
	} `xml:"channel"`
}

type RssRepository struct {
	rssURL string
	client *http.Client
}

func NewRssRepository(rssURL string) *RssRepository {
	return &RssRepository{
		rssURL: rssURL,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (r *RssRepository) FetchJobs() ([]domain.Job, error) {
	resp, err := r.client.Get(r.rssURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar feed RSS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status não esperado do feed RSS: %s", resp.Status)
	}

	var feed RssFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("erro ao decodificar XML: %w", err)
	}
	var jobs []domain.Job
	for _, item := range feed.Channel.Items {
		plainDescription, err := html2text.FromString(item.Description, html2text.Options{PrettyTables: true})
		if err != nil {
			plainDescription = item.Description // Fallback to original if conversion fails
		}
		jobs = append(jobs, domain.Job{
			Title: item.Title,
			Link:  item.Link,
			GUID:  item.GUID,
			Description: plainDescription,
		})
	}

	return jobs, nil
}