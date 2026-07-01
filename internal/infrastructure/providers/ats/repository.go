package ats

import (
	"fmt"
	"log"
	"sync"

	"jobs-bot/config"
	"jobs-bot/internal/domain"
	"jobs-bot/internal/infrastructure/providers/ats/greenhouse"
)

type Repository struct {
	catalogDir        string
	greenhouseClient  AtsClient
	requestedAts      config.AtsConfig
}

func NewRepository(catalogDir string, cfg *config.Config, requestedAts config.AtsConfig) *Repository {
	return &Repository{
		catalogDir:       catalogDir,
		greenhouseClient: greenhouse.NewClient(cfg.GreenhouseAPIKey),
		requestedAts:     requestedAts,
	}
}

func (r *Repository) FetchJobs() ([]domain.Job, error) {
	// Load catalog
	cat, err := LoadCatalog(r.catalogDir)
	if err != nil {
		return nil, fmt.Errorf("ats repository: failed to load catalog: %w", err)
	}

	// Resolve requested companies
	companies, err := cat.ResolveCompanies(r.requestedAts.Collections, r.requestedAts.Companies)
	if err != nil {
		return nil, fmt.Errorf("ats repository: failed to resolve companies: %w", err)
	}

	if len(companies) == 0 {
		return nil, nil
	}

	var wg sync.WaitGroup
	jobsChan := make(chan []domain.Job, len(companies))

	for _, comp := range companies {
		wg.Add(1)
		go func(company CompanyCatalogEntry) {
			defer wg.Done()
			
			var fetched []domain.Job
			var fetchErr error

			switch company.Provider {
			case "greenhouse":
				fetched, fetchErr = r.greenhouseClient.FetchJobs(company.BoardToken)
			default:
				fetchErr = fmt.Errorf("unsupported provider %q", company.Provider)
			}

			if fetchErr != nil {
				// Socratic gate requirement 3: Log individual errors and proceed
				log.Printf("[ATS] ERRO ao buscar vagas para %s (%s): %v", company.Name, company.Provider, fetchErr)
				return
			}

			// Populate company name on each fetched job from catalog details
			for i := range fetched {
				fetched[i].Company = company.Name
			}

			jobsChan <- fetched
		}(comp)
	}

	// Wait for all goroutines to finish and close the channel
	wg.Wait()
	close(jobsChan)

	var allJobs []domain.Job
	for jobs := range jobsChan {
		allJobs = append(allJobs, jobs...)
	}

	return allJobs, nil
}
