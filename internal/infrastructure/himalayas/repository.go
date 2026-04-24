package himalayas

import (
	"context"
	"fmt"
	"log"
	"strings"

	"jobs-bot/internal/domain"
)

const maxPages = 5 // safety cap to avoid hammering the API

// JobRepository fetches remote job listings from the Himalayas Jobs API.
// It implements domain.JobRepository.
type JobRepository struct {
	client      *APIClient
	profileName string
	query       string
}

// NewJobRepository creates a JobRepository for the given profile.
// The query string should be built from the profile's positive_keywords.
func NewJobRepository(profileName string, query string) *JobRepository {
	return &JobRepository{
		client:      NewAPIClient(),
		profileName: profileName,
		query:       query,
	}
}

// FetchJobs retrieves jobs eligible for Brazil: first fetches jobs explicitly
// open to Brazil (country=BR), then fetches worldwide-open jobs, and merges
// both lists deduplicating by GUID.
func (r *JobRepository) FetchJobs() ([]domain.Job, error) {
	seen := make(map[string]struct{})
	var allJobs []domain.Job

	// Pass 1: vagas abertas para o Brasil
	brJobs, err := r.fetchPages(SearchParams{
		Query:          r.query,
		EmploymentType: "Full Time",
		Country:        "BR",
	})
	if err != nil {
		return nil, err
	}
	for _, j := range brJobs {
		if _, dup := seen[j.GUID]; !dup {
			seen[j.GUID] = struct{}{}
			allJobs = append(allJobs, j)
		}
	}
	log.Printf("[Himalayas][%s] Vagas abertas para BR: %d", r.profileName, len(brJobs))

	// Pass 2: vagas worldwide (sem restrição geográfica)
	wwJobs, err := r.fetchPages(SearchParams{
		Query:          r.query,
		EmploymentType: "Full Time",
		Worldwide:      true,
	})
	if err != nil {
		return nil, err
	}
	for _, j := range wwJobs {
		if _, dup := seen[j.GUID]; !dup {
			seen[j.GUID] = struct{}{}
			allJobs = append(allJobs, j)
		}
	}
	log.Printf("[Himalayas][%s] Vagas worldwide adicionadas: %d | Total final: %d",
		r.profileName, len(wwJobs), len(allJobs))

	return allJobs, nil
}

// fetchPages pages through one search query up to maxPages.
func (r *JobRepository) fetchPages(base SearchParams) ([]domain.Job, error) {
	var jobs []domain.Job
	page := 1

	for page <= maxPages {
		base.Page = page
		resp, err := r.client.Search(context.Background(), base)
		if err != nil {
			return nil, fmt.Errorf("himalayas FetchJobs page %d: %w", page, err)
		}

		for _, j := range resp.Jobs {
			jobs = append(jobs, mapToDomain(j))
		}

		if len(jobs) >= resp.TotalCount || len(resp.Jobs) == 0 {
			break
		}
		page++
	}

	return jobs, nil
}

func mapToDomain(j Job) domain.Job {
	location := buildLocation(j)
	return domain.Job{
		Title:           j.Title,
		Link:            j.ApplicationLink,
		GUID:            j.GUID,
		SourceFeed:      "Himalayas",
		Location:        location,
		FullDescription: j.Description,
	}
}

func buildLocation(j Job) string {
	if len(j.LocationRestrictions) == 0 {
		return "Worldwide (Remote)"
	}
	names := make([]string, 0, len(j.LocationRestrictions))
	for _, loc := range j.LocationRestrictions {
		names = append(names, loc.Name)
	}
	return strings.Join(names, ", ")
}
