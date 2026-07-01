package greenhouse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"jobs-bot/internal/domain"
)

type GreenhouseClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewClient(apiKey string) *GreenhouseClient {
	return &GreenhouseClient{
		apiKey:  apiKey,
		baseURL: "https://boards-api.greenhouse.io",
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type greenhouseJobResponse struct {
	Jobs []struct {
		ID          int    `json:"id"`
		Title       string `json:"title"`
		Content     string `json:"content"`
		AbsoluteURL string `json:"absolute_url"`
		Location    struct {
			Name string `json:"name"`
		} `json:"location"`
	} `json:"jobs"`
}

func (c *GreenhouseClient) FetchJobs(boardToken string) ([]domain.Job, error) {
	url := fmt.Sprintf("%s/v1/boards/%s/jobs?content=true", c.baseURL, boardToken)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Inject optional Authorization header if API Key is configured
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	var response greenhouseJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	jobs := make([]domain.Job, 0, len(response.Jobs))
	for _, j := range response.Jobs {
		jobs = append(jobs, domain.Job{
			GUID:            strconv.Itoa(j.ID),
			Title:           j.Title,
			Link:            j.AbsoluteURL,
			Location:        j.Location.Name,
			FullDescription: j.Content,
			SourceFeed:      fmt.Sprintf("ats-greenhouse-%s", boardToken),
		})
	}

	return jobs, nil
}
