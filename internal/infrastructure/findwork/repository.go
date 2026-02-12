package findwork

import (
	"encoding/json"
	"fmt"
	"jobs-bot/internal/domain"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type findworkResponse struct {
	Count   int           `json:"count"`
	Results []findworkJob `json:"results"`
}

type findworkJob struct {
	ID          int    `json:"id"`
	Role        string `json:"role"`
	CompanyName string `json:"company_name"`
	Text        string `json:"text"`
	URL         string `json:"url"`
	Location    string `json:"location"`
}

type Repository struct {
	apiKey   string
	search   string
	location string
	client   *http.Client
}

func NewRepository(apiKey, search, location string) *Repository {
	return &Repository{
		apiKey:   apiKey,
		search:   search,
		location: location,
		client:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (r *Repository) FetchJobs() ([]domain.Job, error) {
	endpoint := "https://findwork.dev/api/jobs/"

	params := url.Values{}
	if r.search != "" {
		params.Set("search", r.search)
	}
	if r.location != "" {
		params.Set("location", r.location)
	}

	reqURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar request Findwork: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", r.apiKey))

	res, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados do Findwork: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Findwork API retornou status %s", res.Status)
	}

	var apiResponse findworkResponse
	if err := json.NewDecoder(res.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("erro ao decodificar JSON do Findwork: %w", err)
	}

	var jobs []domain.Job
	for _, item := range apiResponse.Results {
		jobs = append(jobs, domain.Job{
			Title:           item.Role,
			Link:            item.URL,
			GUID:            strconv.Itoa(item.ID),
			SourceFeed:      "Findwork",
			Location:        item.Location,
			FullDescription: item.Text,
		})
	}

	return jobs, nil
}
