package theirstack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIClient_SearchJobs(t *testing.T) {
	// Mock Server that returns a valid JSON response
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"id": 123,
					"job_title": "Software Engineer",
					"company": "Test Company",
					"url": "https://example.com/job/123",
					"date_posted": "2023-11-20",
					"location": "Remote",
					"country_code": "US",
					"description": "Great remote job"
				}
			],
			"has_more": false,
			"total": 1
		}`))
	}))
	defer ts.Close()

	client := NewAPIClient("fake-key")
	client.apiURL = ts.URL // Inject mocked URL

	req := SearchJobsRequest{
		Page:  0,
		Limit: 1,
	}

	resp, err := client.SearchJobs(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no errors, got %v", err)
	}

	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 job, got %d", len(resp.Data))
	}

	job := resp.Data[0]
	if job.ID != 123 {
		t.Errorf("expected job ID 123, got %d", job.ID)
	}
	if job.JobTitle != "Software Engineer" {
		t.Errorf("expected Software Engineer, got %s", job.JobTitle)
	}
}

func TestAPIClient_RateLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Retry-After", "5")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer ts.Close()

	client := NewAPIClient("fake-key")
	client.apiURL = ts.URL

	_, err := client.SearchJobs(context.Background(), SearchJobsRequest{})
	if err == nil {
		t.Fatal("expected an error, got nil")
	}

	expectedErrFragment := "retry after 5 seconds"
	if err.Error() != "rate limit exceeded, retry after 5 seconds" {
		t.Errorf("expected error containing %q, got %q", expectedErrFragment, err.Error())
	}
}
