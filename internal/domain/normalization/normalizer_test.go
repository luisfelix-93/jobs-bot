package normalization

import (
	"testing"

	"jobs-bot/internal/domain"
)

func TestSeniorityNormalizer(t *testing.T) {
	norm := NewSeniorityNormalizer()

	tests := []struct {
		title    string
		expected string
	}{
		{"Senior Go Engineer", "Senior"},
		{"Tech Lead - Python", "Lead"},
		{"Junior Developer", "Junior"},
		{"Software Engineer II", "Mid"},
		{"Staff Engineer", "Staff"},
		{"Principal Architect", "Principal"},
		{"Pleno Backend Developer", "Mid"},
	}

	for _, tc := range tests {
		job := domain.Job{Title: tc.title}
		norm.Normalize(&job)
		if job.Seniority != tc.expected {
			t.Errorf("For title %q, expected seniority %q, got %q", tc.title, tc.expected, job.Seniority)
		}
	}
}

func TestWorkModeNormalizer(t *testing.T) {
	norm := NewWorkModeNormalizer()

	tests := []struct {
		title       string
		location    string
		description string
		expected    string
	}{
		{"Go Developer", "Remote", "", "Remote"},
		{"Hybrid Software Engineer", "New York", "", "Hybrid"},
		{"On-site DevOps", "Dallas, TX", "", "On-site"},
		{"Engineer", "Worldwide", "", "Remote"},
		{"Backend Developer", "São Paulo", "This is a 100% remote job for everyone", "Remote"},
		{"Developer", "London", "We offer a model híbrido of work", "Hybrid"},
	}

	for _, tc := range tests {
		job := domain.Job{
			Title:           tc.title,
			Location:        tc.location,
			FullDescription: tc.description,
		}
		norm.Normalize(&job)
		if job.WorkMode != tc.expected {
			t.Errorf("For title %q, location %q, expected work mode %q, got %q", tc.title, tc.location, tc.expected, job.WorkMode)
		}
	}
}

func TestEmploymentTypeNormalizer(t *testing.T) {
	norm := NewEmploymentTypeNormalizer()

	tests := []struct {
		title       string
		description string
		expected    string
	}{
		{"CLT Backend Go Developer", "", "FullTime"},
		{"Contractor Python Engineer", "", "Contract"},
		{"Part-time React Developer", "", "PartTime"},
		{"Go Dev", "We are hiring under PJ contract", "Contract"},
	}

	for _, tc := range tests {
		job := domain.Job{Title: tc.title, FullDescription: tc.description}
		norm.Normalize(&job)
		if job.EmploymentType != tc.expected {
			t.Errorf("Expected employment type %q, got %q", tc.expected, job.EmploymentType)
		}
	}
}

func TestTitleNormalizer(t *testing.T) {
	norm := NewTitleNormalizer()

	tests := []struct {
		title    string
		company  string
		expected string
	}{
		{"Google - Software Engineer (Remote)", "Google", "Software Engineer"},
		{"Backend Developer [Hybrid]", "", "Backend Developer"},
		{"DevOps Engineer - 100% Remote", "", "DevOps Engineer"},
	}

	for _, tc := range tests {
		job := domain.Job{Title: tc.title, Company: tc.company}
		norm.Normalize(&job)
		if job.NormalizedTitle != tc.expected {
			t.Errorf("For title %q, expected normalized %q, got %q", tc.title, tc.expected, job.NormalizedTitle)
		}
	}
}

func TestSkillsExtractor(t *testing.T) {
	norm := NewSkillsExtractor()

	job := domain.Job{
		Title:           "Go Developer",
		FullDescription: "Looking for a Go developer with experience in Kubernetes, AWS and React. Know Python too.",
	}
	norm.Normalize(&job)

	expectedSkills := map[string]bool{
		"Go":         true,
		"Kubernetes": true,
		"AWS":        true,
		"React":      true,
		"Python":     true,
	}

	if len(job.Skills) != len(expectedSkills) {
		t.Errorf("Expected %d skills, got %d: %v", len(expectedSkills), len(job.Skills), job.Skills)
	}

	for _, s := range job.Skills {
		if !expectedSkills[s] {
			t.Errorf("Unexpected skill extracted: %s", s)
		}
	}
}

func TestSalaryNormalizer(t *testing.T) {
	norm := NewSalaryNormalizer()

	tests := []struct {
		title       string
		description string
		expectedMin float64
		expectedMax float64
		expectedCur string
	}{
		{"Go Engineer ($120k-$150k)", "", 120000, 150000, "USD"},
		{"Developer", "Salary: USD 80,000 to 100,000", 80000, 100000, "USD"},
		{"Lead Engineer", "We pay €100k", 100000, 100000, "EUR"},
		{"Engineer", "BRL 10.000 - BRL 12.000", 10000, 12000, "BRL"},
	}

	for _, tc := range tests {
		job := domain.Job{Title: tc.title, FullDescription: tc.description}
		norm.Normalize(&job)
		if job.SalaryMin != tc.expectedMin || job.SalaryMax != tc.expectedMax || job.SalaryCurrency != tc.expectedCur {
			t.Errorf("For %q, expected min %.0f, max %.0f, cur %q; got min %.0f, max %.0f, cur %q",
				tc.title, tc.expectedMin, tc.expectedMax, tc.expectedCur, job.SalaryMin, job.SalaryMax, job.SalaryCurrency)
		}
	}
}

func TestLocationNormalizer(t *testing.T) {
	norm := NewLocationNormalizer()

	tests := []struct {
		loc      string
		expected string
	}{
		{"USA", "United States"},
		{"UK ", "United Kingdom"},
		{"BR", "Brazil"},
		{"New York, US", "New York, United States"},
		{"Anywhere", "Remote"},
	}

	for _, tc := range tests {
		job := domain.Job{Location: tc.loc}
		norm.Normalize(&job)
		if job.Location != tc.expected {
			t.Errorf("For location %q, expected %q, got %q", tc.loc, tc.expected, job.Location)
		}
	}
}

func TestPipeline(t *testing.T) {
	pipeline := NewPipeline(
		NewSeniorityNormalizer(),
		NewTitleNormalizer(),
	)

	jobs := []domain.Job{
		{Title: "Senior Developer - Remote", Company: "Company"},
	}

	jobs = pipeline.NormalizeAll(jobs)

	if jobs[0].Seniority != "Senior" {
		t.Errorf("Expected Seniority to be 'Senior', got %q", jobs[0].Seniority)
	}
	if jobs[0].NormalizedTitle != "Developer" {
		t.Errorf("Expected NormalizedTitle to be 'Developer', got %q", jobs[0].NormalizedTitle)
	}
}
