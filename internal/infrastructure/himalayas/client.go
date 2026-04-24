package himalayas

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	searchEndpoint = "https://himalayas.app/jobs/api/search"
	maxPerPage     = 20
)

// APIClient handles HTTP requests to the Himalayas Remote Jobs API.
// No authentication is required; the API is public and free.
type APIClient struct {
	httpClient *http.Client
}

// NewAPIClient creates a new Himalayas API client.
func NewAPIClient() *APIClient {
	return &APIClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SearchParams holds the query parameters for the search endpoint.
type SearchParams struct {
	Query          string
	EmploymentType string // e.g. "Full Time"
	Country        string // ISO alpha-2 country code, e.g. "BR"
	Worldwide      bool   // if true, returns only jobs open worldwide
	Page           int    // 1-based
}

// Search calls the Himalayas search endpoint and returns one page of results.
func (c *APIClient) Search(ctx context.Context, params SearchParams) (*SearchResponse, error) {
	reqURL, err := url.Parse(searchEndpoint)
	if err != nil {
		return nil, fmt.Errorf("himalayas: invalid endpoint URL: %w", err)
	}

	q := reqURL.Query()
	if params.Query != "" {
		q.Set("q", params.Query)
	}
	if params.EmploymentType != "" {
		q.Set("employment_type", params.EmploymentType)
	}
	if params.Country != "" {
		q.Set("country", params.Country)
	}
	if params.Worldwide {
		q.Set("worldwide", "true")
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	q.Set("limit", strconv.Itoa(maxPerPage))
	reqURL.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("himalayas: error creating request: %w", err)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("himalayas: error executing request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusTooManyRequests:
		return nil, fmt.Errorf("himalayas: rate limit exceeded (429). Try again later")
	case http.StatusBadRequest:
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("himalayas: bad request (400): %s", string(body))
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("himalayas: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("himalayas: error decoding response: %w", err)
	}

	return &result, nil
}
