package normalization

import (
	"strings"

	"jobs-bot/internal/domain"
)

// LocationNormalizer standardizes and cleans location descriptions.
type LocationNormalizer struct{}

func NewLocationNormalizer() *LocationNormalizer {
	return &LocationNormalizer{}
}

func (n *LocationNormalizer) Normalize(job *domain.Job) {
	loc := strings.TrimSpace(job.Location)
	loc = strings.Join(strings.Fields(loc), " ")

	locLower := strings.ToLower(loc)

	switch locLower {
	case "us", "usa", "united states of america":
		loc = "United States"
	case "uk", "united kingdom of great britain and northern ireland":
		loc = "United Kingdom"
	case "br", "brazil", "brasil":
		loc = "Brazil"
	case "anywhere", "remote", "remoto", "worldwide", "global":
		loc = "Remote"
	}

	if strings.HasSuffix(locLower, ", us") || strings.HasSuffix(locLower, ", usa") {
		parts := strings.Split(loc, ",")
		if len(parts) > 1 {
			parts[len(parts)-1] = " United States"
			loc = strings.Join(parts, ",")
		}
	}

	job.Location = loc
}
