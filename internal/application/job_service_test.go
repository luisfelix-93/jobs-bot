package application

import (
	"fmt"
	"jobs-bot/internal/domain"
	"testing"
)

// --- Mocks ---

type mockRepo struct {
	jobs []domain.Job
	err  error
}

func (m *mockRepo) FetchJobs() ([]domain.Job, error) {
	return m.jobs, m.err
}

type mockNotifier struct {
	notified []domain.Job
}

func (m *mockNotifier) Notify(job domain.Job, analysis domain.ResumeAnalysis, ai *domain.AIAnalysis) error {
	m.notified = append(m.notified, job)
	return nil
}

type mockStore struct {
	existing map[string]bool
	saved    []domain.ProcessedJob
}

func (m *mockStore) Exists(guid, profile string) (bool, error) {
	key := fmt.Sprintf("%s_%s", guid, profile)
	return m.existing[key], nil
}

func (m *mockStore) Save(job domain.ProcessedJob) error {
	m.saved = append(m.saved, job)
	return nil
}

func (m *mockStore) Close() error { return nil }

type mockAIAnalyzer struct {
	callCount int
	score     int
}

func (m *mockAIAnalyzer) Analyze(resume, jobDesc string) (*domain.AIAnalysis, error) {
	m.callCount++
	return &domain.AIAnalysis{
		Score:          m.score,
		Recommendation: "apply",
		Source:         "mock",
	}, nil
}

// --- Helper ---

func buildService(repos []domain.JobRepository, notifier domain.NotificationService, store domain.JobStore, ai domain.AIAnalyzer, limit int) *JobService {
	filter := domain.NewJobFilter([]string{"go", "docker"}, nil)
	analyzer := domain.NewResumeAnalyzer()
	return NewJobService(repos, notifier, filter, analyzer, ai, store, "resume content with go", []string{"go"}, "test-profile", limit)
}

func makeJobs(count int) []domain.Job {
	jobs := make([]domain.Job, count)
	for i := range jobs {
		jobs[i] = domain.Job{
			Title:           fmt.Sprintf("Go Developer %d", i),
			GUID:            fmt.Sprintf("guid-%d", i),
			SourceFeed:      "TestSource",
			FullDescription: "go docker kubernetes",
		}
	}
	return jobs
}

// --- Tests ---

func TestProcessNewJobs_SkipsDuplicatesAndFindsNew(t *testing.T) {
	jobs := makeJobs(10)
	repo := &mockRepo{jobs: jobs}
	notifier := &mockNotifier{}

	existing := map[string]bool{}
	for i := 0; i < 5; i++ {
		key := fmt.Sprintf("TestSource-guid-%d_test-profile", i)
		existing[key] = true
	}
	store := &mockStore{existing: existing}
	ai := &mockAIAnalyzer{score: 80}

	svc := buildService([]domain.JobRepository{repo}, notifier, store, ai, 50)
	stats, err := svc.ProcessNewJobs()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TotalSkipped != 5 {
		t.Errorf("expected 5 skipped (duplicates), got %d", stats.TotalSkipped)
	}
	if stats.TotalNotified != 5 {
		t.Errorf("expected 5 notified (new jobs), got %d", stats.TotalNotified)
	}
	if len(store.saved) != 5 {
		t.Errorf("expected 5 saved, got %d", len(store.saved))
	}
}

func TestProcessNewJobs_LimitAppliesToNewJobsOnly(t *testing.T) {
	jobs := makeJobs(20)
	repo := &mockRepo{jobs: jobs}
	notifier := &mockNotifier{}

	existing := map[string]bool{}
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("TestSource-guid-%d_test-profile", i)
		existing[key] = true
	}
	store := &mockStore{existing: existing}
	ai := &mockAIAnalyzer{score: 80}

	svc := buildService([]domain.JobRepository{repo}, notifier, store, ai, 5)
	stats, err := svc.ProcessNewJobs()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TotalSkipped != 10 {
		t.Errorf("expected 10 skipped (duplicates), got %d", stats.TotalSkipped)
	}
	if stats.TotalNotified != 5 {
		t.Errorf("expected 5 notified (limited to 5 new), got %d", stats.TotalNotified)
	}
	if len(store.saved) != 5 {
		t.Errorf("expected 5 saved (new only), got %d", len(store.saved))
	}
}

func TestProcessNewJobs_AINotCalledForDuplicates(t *testing.T) {
	jobs := makeJobs(10)
	repo := &mockRepo{jobs: jobs}
	notifier := &mockNotifier{}

	existing := map[string]bool{}
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("TestSource-guid-%d_test-profile", i)
		existing[key] = true
	}
	store := &mockStore{existing: existing}
	ai := &mockAIAnalyzer{score: 80}

	svc := buildService([]domain.JobRepository{repo}, notifier, store, ai, 50)
	_, err := svc.ProcessNewJobs()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ai.callCount != 0 {
		t.Errorf("expected 0 AI calls for all-duplicate jobs, got %d", ai.callCount)
	}
}

func TestProcessNewJobs_AICalledOnlyForNewJobs(t *testing.T) {
	jobs := makeJobs(10)
	repo := &mockRepo{jobs: jobs}
	notifier := &mockNotifier{}

	existing := map[string]bool{}
	for i := 0; i < 7; i++ {
		key := fmt.Sprintf("TestSource-guid-%d_test-profile", i)
		existing[key] = true
	}
	store := &mockStore{existing: existing}
	ai := &mockAIAnalyzer{score: 80}

	svc := buildService([]domain.JobRepository{repo}, notifier, store, ai, 50)
	_, err := svc.ProcessNewJobs()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ai.callCount != 3 {
		t.Errorf("expected 3 AI calls (only for 3 new jobs), got %d", ai.callCount)
	}
}

func TestProcessNewJobs_ZeroLimitProcessesAll(t *testing.T) {
	jobs := makeJobs(10)
	repo := &mockRepo{jobs: jobs}
	notifier := &mockNotifier{}
	store := &mockStore{existing: map[string]bool{}}
	ai := &mockAIAnalyzer{score: 80}

	svc := buildService([]domain.JobRepository{repo}, notifier, store, ai, 0)
	stats, err := svc.ProcessNewJobs()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.TotalNotified != 10 {
		t.Errorf("expected all 10 notified when limit is 0 (unlimited), got %d", stats.TotalNotified)
	}
}
