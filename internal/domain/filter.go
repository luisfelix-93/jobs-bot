package domain

import (
	"sort"
	"strings"
)
 
type JobFilter struct {
	PositiveKeywords []string
	NegativeKeywords []string
}

func NewJobFilter(positive, negative []string) *JobFilter {
	return &JobFilter{
		PositiveKeywords: positive,
		NegativeKeywords: negative,
	}
}

func (f *JobFilter) FilterAndRankJobs(jobs []Job, limit int) []Job {

	type rankedJob struct {
		job   Job
		score int
	}

	var rankedJobs []rankedJob

	for _, job := range jobs {
		fullText := strings.ToLower(job.Title + " " + job.FullDescription)

		if f.containsNegativeKeyword(fullText) {
			continue
		}

		score := f.calculateScore(fullText)

		if score > 0 {
			rankedJobs = append(rankedJobs, rankedJob{
				job:   job,
				score: score,
			})
		}
	}

	sort.Slice(rankedJobs, func(i, j int) bool {
		return rankedJobs[i].score > rankedJobs[j].score
	})

	// Prepara o retorno simples
	var bestJobs []Job
	for i := 0; i < len(rankedJobs) && i < limit; i++ {
		bestJobs = append(bestJobs, rankedJobs[i].job)
	}

	return bestJobs
}


func (f *JobFilter) calculateScore(fullText string) int {
	score := 0
	for _, keyword := range f.PositiveKeywords {
		if strings.Contains(fullText, strings.ToLower(keyword)) {
			score++
		}
	}
	return score
}

func (f *JobFilter) containsNegativeKeyword(fullText string) bool {
	for _, keyword := range f.NegativeKeywords {
		if strings.Contains(fullText, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}