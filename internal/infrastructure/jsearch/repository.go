package jsearch

import (
	"encoding/json"
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"net/url"
	"time"
)

type jsearchResponse struct {
	Status string       `json:"status"`
	Data   []jsearchJob `json:"data"`
}

type jsearchJob struct {
	JobID          string `json:"job_id"`
	JobTitle       string `json:"job_title"`
	EmployerName   string `json:"employer_name"`
	JobDescription string `json:"job_description"`
	JobApplyLink   string `json:"job_apply_link"`
	JobCity        string `json:"job_city"`
	JobCountry     string `json:"job_country"`
	JobIsRemote    bool   `json:"job_is_remote"`
}

type Repository struct {
	apiKey string
	query  string
	client *http.Client
}

func NewRepository(apiKey, query string) *Repository {
	return &Repository{
		apiKey: apiKey,
		query:  query,
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

func (r *Repository) FetchJobs() ([]domain.Job, error) {
	endpoint := "https://jsearch.p.rapidapi.com/search"

	params := url.Values{}
	params.Set("query", r.query)
	params.Set("page", "1")
	params.Set("num_pages", "1")
	params.Set("date_posted", "week")
	params.Set("remote_jobs_only", "true")

	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request JSearch: %w", err)
	}

	req.Header.Set("X-RapidAPI-Key", r.apiKey)
	req.Header.Set("X-RapidAPI-Host", "jsearch.p.rapidapi.com")

	res, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados do JSearch: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("JSearch API retornou status %s", res.Status)
	}

	var apiResponse jsearchResponse
	if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do JSearch: %w", err)
	}

	var jobs []domain.Job
	for _, item := range apiResponse.Data {
		location := item.JobCity
		if item.JobCountry != "" {
			if location != "" {
				location += ", "
			}
			location += item.JobCountry
		}
		if item.JobIsRemote {
			location = "Remote - " + location
		}

		jobs = append(jobs, domain.Job{
			Title:           item.JobTitle,
			Link:            item.JobApplyLink,
			GUID:            item.JobID,
			SourceFeed:      "JSearch",
			Location:        location,
			FullDescription: item.JobDescription,
		})
	}

	return jobs, nil
}
