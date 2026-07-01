package ats

import "jobs-bot/internal/domain"

type AtsClient interface {
	FetchJobs(boardToken string) ([]domain.Job, error)
}
