package jobicy

import (
	"encoding/json"
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"strconv"
	"time"
)

type jobicyAPIResponse struct {
	Jobs []jobicyJob `json:"jobs"`
}

type jobicyJob struct {
	ID             int    `json:"id"`
	URL            string `json:"url"`
	JobTitle       string `json:"jobTitle"`
	JobGeo         string `json:"jobGeo"`
	JobDescription string `json:"jobDescription"`
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
	req, err := http.NewRequest("GET", r.rssURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request para a API do Jobicy: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	res, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados da Jobicy API: %w", err)
	}
	defer res.Body.Close()
	
	var apiResponse jobicyAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do Jobicy: %w", err)
	}

	var jobs []domain.Job
	for _, item := range apiResponse.Jobs {
		jobs = append(jobs, domain.Job{
			Title:           item.JobTitle,
			Link:            item.URL,
			GUID:            strconv.Itoa(item.ID),
			SourceFeed:      "Jobicy",
			Location:        item.JobGeo,
			FullDescription: item.JobDescription,
			
		})
	}
	return jobs, nil
}

