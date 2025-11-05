package jobicy

import (
	"encoding/xml"
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"time"
)

type jobicyItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"` // Este é apenas um snippet, vamos ignorá-lo
	FullContent string `xml:"encoded"`     // Tag <content:encoded>
	Location    string `xml:"location"`
}

type jobicyRss struct {
	Channel struct {
		Items []jobicyItem `xml:"item"`
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
	res, err := r.client.Get(r.rssURL)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar feed Jobicy: %w", err)
	}
	defer res.Body.Close()

	var feed jobicyRss
	if err := xml.NewDecoder(res.Body).Decode(&feed); err != nil {
		return nil, fmt.Errorf("erro ao decodificar XML do Jobicy: %w", err)
	}

	var jobs []domain.Job
	for _, item := range feed.Channel.Items {
		jobs = append(jobs, domain.Job{
			Title:           item.Title,
			Link:            item.Link,
			GUID:            item.GUID,
			SourceFeed:      "Jobicy",
			Location:        item.Location,
			FullDescription: item.FullContent,
		})
	}

	return jobs, nil
}

