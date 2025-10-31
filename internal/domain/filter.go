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
		score := f.calculateJobScore(job.Title)
		if score >  0 && !f.containsNegativeKeyword(job.Title) {
			rankedJobs = append(rankedJobs, rankedJob{job: job, score: score})
		}
	}

	sort.Slice(rankedJobs, func(i, j int) bool {
		return rankedJobs[i].score > rankedJobs[j].score
	})

	var bestJobs []Job
	for i := 0; i < len(rankedJobs) && i < limit; i++ {
		bestJobs = append(bestJobs, rankedJobs[i].job)
	}

	return bestJobs
}

func (f *JobFilter) calculateJobScore(title string) int {
score := 0
	lowerTitle := strings.ToLower(title)
	for _, keyword := range f.PositiveKeywords {
		if strings.Contains(lowerTitle, strings.ToLower(keyword)) {
			score++
		}
	}
	return score
}

func (f *JobFilter) containsNegativeKeyword(title string) bool {
	lowerTitle := strings.ToLower(title)
	for _, keyword := range f.NegativeKeywords {
		if strings.Contains(lowerTitle, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}