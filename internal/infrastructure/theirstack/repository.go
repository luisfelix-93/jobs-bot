package theirstack

import (
	"context"
	"fmt"
	"log"

	"jobs-bot/config"
	"jobs-bot/internal/domain"
)

// JobRepository TheirStack API implementation
type JobRepository struct {
	client  *APIClient
	profile config.ProfileConfig
}

// NewJobRepository initializes the TheirStack API repository
func NewJobRepository(apiKey string, profile config.ProfileConfig) *JobRepository {
	client := NewAPIClient(apiKey)
	if profile.Sources.TheirStackURL != "" {
		client.apiURL = profile.Sources.TheirStackURL
	}
	return &JobRepository{
		client:  client,
		profile: profile,
	}
}

// FetchJobs retrieves remote jobs utilizing positive_keywords from the profile
func (r *JobRepository) FetchJobs() ([]domain.Job, error) {
	if r.client.apiKey == "" {
		return nil, fmt.Errorf("TheirStack API Key not configured")
	}

	isRemote := true
	maxAge := 3

	// Build query
	req := SearchJobsRequest{
		Page:  0,
		Limit: 25,
		OrderBy: []Order{
			{Field: "date_posted", Desc: true},
		},
		JobTitleOr:              r.profile.PositiveKeywords,
		JobDescriptionPatternOr: r.profile.PositiveKeywords, // Use the keywords on description search too
		Remote:                  &isRemote,
		PostedAtMaxAgeDays:      &maxAge,
	}

	resp, err := r.client.SearchJobs(context.Background(), req)
	if err != nil {
		return nil, fmt.Errorf("TheirStack Fetch Error: %w", err)
	}

	var jobs []domain.Job
	for _, tsJob := range resp.Data {
		jobs = append(jobs, domain.Job{
			Title:           tsJob.JobTitle,
			Link:            tsJob.URL,
			GUID:            fmt.Sprintf("%d", tsJob.ID),
			SourceFeed:      "TheirStack",
			Location:        "Remote (" + tsJob.CountryCode + ")",
			FullDescription: tsJob.Description,
		})
	}

	log.Printf("[TheirStack] Fetched %d jobs for profile %s", len(jobs), r.profile.Name)
	return jobs, nil
}
