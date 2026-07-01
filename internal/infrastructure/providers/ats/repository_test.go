package ats

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"jobs-bot/config"
	"jobs-bot/internal/domain"
)

type MockAtsClient struct {
	FetchFunc func(boardToken string) ([]domain.Job, error)
}

func (m *MockAtsClient) FetchJobs(boardToken string) ([]domain.Job, error) {
	return m.FetchFunc(boardToken)
}

func TestRepositoryFetchJobs(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "repo-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	collectionsData := `
collections:
  test-col:
    - ok-comp
    - fail-comp
`
	greenhouseData := `
companies:
  ok-comp:
    name: "OK Company"
    board_token: "ok-token"
  fail-comp:
    name: "Fail Company"
    board_token: "fail-token"
`
	if err := os.WriteFile(filepath.Join(tempDir, "collections.yaml"), []byte(collectionsData), 0644); err != nil {
		t.Fatalf("failed to write collections.yaml: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tempDir, "greenhouse.yaml"), []byte(greenhouseData), 0644); err != nil {
		t.Fatalf("failed to write greenhouse.yaml: %v", err)
	}

	mockClient := &MockAtsClient{
		FetchFunc: func(boardToken string) ([]domain.Job, error) {
			if boardToken == "ok-token" {
				return []domain.Job{
					{GUID: "1", Title: "Job 1", SourceFeed: "mock"},
				}, nil
			}
			if boardToken == "fail-token" {
				return nil, errors.New("http error")
			}
			return nil, errors.New("unknown token")
		},
	}

	repo := &Repository{
		catalogDir:       tempDir,
		greenhouseClient: mockClient,
		requestedAts: config.AtsConfig{
			Collections: []string{"test-col"},
		},
	}

	jobs, err := repo.FetchJobs()
	if err != nil {
		t.Fatalf("FetchJobs failed: %v", err)
	}

	// Only OK Company jobs should be returned
	// Fail Company's error should be caught, logged, and skipped without breaking execution
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}

	if jobs[0].Company != "OK Company" {
		t.Errorf("expected Company 'OK Company', got %q", jobs[0].Company)
	}
	if jobs[0].Title != "Job 1" {
		t.Errorf("expected Title 'Job 1', got %q", jobs[0].Title)
	}
}
