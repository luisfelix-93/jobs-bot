package domain

import (
	"testing"
)

func TestFilterAndRankJobs_ReturnsAllValidJobs(t *testing.T) {
	filter := NewJobFilter([]string{"go", "docker"}, []string{"senior"})

	jobs := []Job{
		{Title: "Go Developer", GUID: "1", FullDescription: "Work with go and docker"},
		{Title: "Docker Engineer", GUID: "2", FullDescription: "docker kubernetes"},
		{Title: "Python Dev", GUID: "3", FullDescription: "python flask"},
		{Title: "Go Senior", GUID: "4", FullDescription: "senior go developer"},
	}

	result := filter.FilterAndRankJobs(jobs)

	// Job 1: matches "go" + "docker" (score 2)
	// Job 2: matches "docker" (score 1)
	// Job 3: no positive keywords (score 0, excluded)
	// Job 4: contains negative keyword "senior" (excluded)
	if len(result) != 2 {
		t.Fatalf("expected 2 jobs, got %d", len(result))
	}

	if result[0].GUID != "1" {
		t.Errorf("expected first job to be GUID '1' (highest score), got '%s'", result[0].GUID)
	}
	if result[1].GUID != "2" {
		t.Errorf("expected second job to be GUID '2', got '%s'", result[1].GUID)
	}
}

func TestFilterAndRankJobs_NoLimitTruncation(t *testing.T) {
	filter := NewJobFilter([]string{"golang"}, nil)

	jobs := make([]Job, 100)
	for i := range jobs {
		jobs[i] = Job{
			Title:           "Golang Dev",
			GUID:            string(rune('A' + i%26)),
			FullDescription: "golang microservices",
		}
	}

	result := filter.FilterAndRankJobs(jobs)

	if len(result) != 100 {
		t.Fatalf("expected all 100 jobs returned without truncation, got %d", len(result))
	}
}

func TestFilterAndRankJobs_EmptyInput(t *testing.T) {
	filter := NewJobFilter([]string{"go"}, nil)
	result := filter.FilterAndRankJobs(nil)

	if len(result) != 0 {
		t.Fatalf("expected 0 jobs for nil input, got %d", len(result))
	}
}

func TestFilterAndRankJobs_AllNegativeExcluded(t *testing.T) {
	filter := NewJobFilter([]string{"go"}, []string{"intern"})

	jobs := []Job{
		{Title: "Go Intern", GUID: "1", FullDescription: "go intern position"},
		{Title: "Go Intern 2", GUID: "2", FullDescription: "go intern role"},
	}

	result := filter.FilterAndRankJobs(jobs)

	if len(result) != 0 {
		t.Fatalf("expected 0 jobs (all excluded by negative keyword), got %d", len(result))
	}
}

func TestFilterAndRankJobs_RankedByScore(t *testing.T) {
	filter := NewJobFilter([]string{"go", "docker", "kubernetes"}, nil)

	jobs := []Job{
		{Title: "A", GUID: "1", FullDescription: "go"},                      // score 1
		{Title: "B", GUID: "2", FullDescription: "go docker kubernetes"},     // score 3
		{Title: "C", GUID: "3", FullDescription: "go docker"},               // score 2
	}

	result := filter.FilterAndRankJobs(jobs)

	if len(result) != 3 {
		t.Fatalf("expected 3 jobs, got %d", len(result))
	}
	if result[0].GUID != "2" {
		t.Errorf("expected GUID '2' first (score 3), got '%s'", result[0].GUID)
	}
	if result[1].GUID != "3" {
		t.Errorf("expected GUID '3' second (score 2), got '%s'", result[1].GUID)
	}
	if result[2].GUID != "1" {
		t.Errorf("expected GUID '1' third (score 1), got '%s'", result[2].GUID)
	}
}
