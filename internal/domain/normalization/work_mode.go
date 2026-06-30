package normalization

import (
	"strings"

	"jobs-bot/internal/domain"
)

// WorkModeNormalizer detects the work mode (Remote, Hybrid, On-site).
type WorkModeNormalizer struct{}

func NewWorkModeNormalizer() *WorkModeNormalizer {
	return &WorkModeNormalizer{}
}

func (n *WorkModeNormalizer) Normalize(job *domain.Job) {
	if job.WorkMode != "" {
		switch strings.ToLower(job.WorkMode) {
		case "remote", "remoto", "home office", "telework":
			job.WorkMode = "Remote"
			return
		case "hybrid", "híbrido", "misto":
			job.WorkMode = "Hybrid"
			return
		case "on-site", "onsite", "presencial", "in office", "office":
			job.WorkMode = "On-site"
			return
		}
	}

	location := " " + strings.ToLower(job.Location) + " "
	title := " " + strings.ToLower(job.Title) + " "
	desc := " " + strings.ToLower(job.FullDescription) + " "

	// Remote signals
	if strings.Contains(location, "remote") || strings.Contains(location, "remoto") || strings.Contains(location, "worldwide") ||
		strings.Contains(title, "remote") || strings.Contains(title, "remoto") || strings.Contains(title, "home office") || strings.Contains(title, "anywhere") {
		job.WorkMode = "Remote"
		return
	}

	// Hybrid signals
	if strings.Contains(location, "hybrid") || strings.Contains(location, "híbrido") ||
		strings.Contains(title, "hybrid") || strings.Contains(title, "híbrido") {
		job.WorkMode = "Hybrid"
		return
	}

	// On-site signals
	if strings.Contains(location, "on-site") || strings.Contains(location, "onsite") || strings.Contains(location, "presencial") ||
		strings.Contains(title, "on-site") || strings.Contains(title, "onsite") || strings.Contains(title, "presencial") {
		job.WorkMode = "On-site"
		return
	}

	// Description heuristics
	if strings.Contains(desc, "100% remote") || strings.Contains(desc, "fully remote") || strings.Contains(desc, "trabalho remoto") || strings.Contains(desc, "home office") || strings.Contains(desc, "remote-first") || strings.Contains(desc, "remote") || strings.Contains(desc, "remoto") {
		job.WorkMode = "Remote"
		return
	}
	if strings.Contains(desc, "hybrid") || strings.Contains(desc, "híbrido") || strings.Contains(desc, "hibrido") {
		job.WorkMode = "Hybrid"
		return
	}
	if strings.Contains(desc, "on-site") || strings.Contains(desc, "onsite") || strings.Contains(desc, "presencial") || strings.Contains(desc, "in-office") || strings.Contains(desc, "in office") || strings.Contains(desc, "work from office") {
		job.WorkMode = "On-site"
		return
	}
}
