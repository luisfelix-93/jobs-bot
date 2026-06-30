package normalization

import "jobs-bot/internal/domain"

// Normalizer defines the interface for modifying jobs to normalize their data.
type Normalizer interface {
	Normalize(job *domain.Job)
}

// Pipeline holds a collection of Normalizers to run in sequence.
type Pipeline struct {
	normalizers []Normalizer
}

// NewPipeline creates a new normalization pipeline with the given normalizers.
func NewPipeline(normalizers ...Normalizer) *Pipeline {
	return &Pipeline{
		normalizers: normalizers,
	}
}

// NormalizeAll runs all normalizers on a slice of jobs.
func (p *Pipeline) NormalizeAll(jobs []domain.Job) []domain.Job {
	for i := range jobs {
		for _, norm := range p.normalizers {
			norm.Normalize(&jobs[i])
		}
	}
	return jobs
}
