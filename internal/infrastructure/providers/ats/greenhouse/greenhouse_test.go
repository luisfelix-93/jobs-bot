package greenhouse

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreenhouseClientFetchJobs(t *testing.T) {
	mockResponse := `{
		"jobs": [
			{
				"id": 4048440,
				"title": "Software Engineer - Go",
				"content": "<p>We are looking for a Go engineer...</p>",
				"absolute_url": "https://boards.greenhouse.io/openai/jobs/4048440",
				"location": {
					"name": "San Francisco, CA"
				}
			}
		],
		"meta": {
			"total": 1
		}
	}`

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify URL format
		expectedPath := "/v1/boards/openai/jobs"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %q, got %q", expectedPath, r.URL.Path)
		}
		if r.URL.Query().Get("content") != "true" {
			t.Errorf("expected query param content=true, got %q", r.URL.RawQuery)
		}
		
		// Verify optional token header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && authHeader != "Bearer test-api-key" {
			t.Errorf("unexpected authorization header %q", authHeader)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer mockServer.Close()

	// Test without API key
	client := NewClient("")
	client.baseURL = mockServer.URL

	jobs, err := client.FetchJobs("openai")
	if err != nil {
		t.Fatalf("FetchJobs failed: %v", err)
	}

	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}

	j := jobs[0]
	if j.GUID != "4048440" {
		t.Errorf("expected GUID '4048440', got %q", j.GUID)
	}
	if j.Title != "Software Engineer - Go" {
		t.Errorf("expected Title 'Software Engineer - Go', got %q", j.Title)
	}
	if j.Link != "https://boards.greenhouse.io/openai/jobs/4048440" {
		t.Errorf("expected Link 'https://boards.greenhouse.io/openai/jobs/4048440', got %q", j.Link)
	}
	if j.Location != "San Francisco, CA" {
		t.Errorf("expected Location 'San Francisco, CA', got %q", j.Location)
	}
	if j.FullDescription != "<p>We are looking for a Go engineer...</p>" {
		t.Errorf("expected FullDescription, got %q", j.FullDescription)
	}
	if j.SourceFeed != "ats-greenhouse-openai" {
		t.Errorf("expected SourceFeed 'ats-greenhouse-openai', got %q", j.SourceFeed)
	}

	// Test with API key
	clientWithKey := NewClient("test-api-key")
	clientWithKey.baseURL = mockServer.URL

	_, err = clientWithKey.FetchJobs("openai")
	if err != nil {
		t.Fatalf("FetchJobs with key failed: %v", err)
	}
}
