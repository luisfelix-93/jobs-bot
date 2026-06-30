package normalization

import (
	"strings"

	"jobs-bot/internal/domain"
)

// SeniorityNormalizer identifies the job seniority from its title.
type SeniorityNormalizer struct{}

func NewSeniorityNormalizer() *SeniorityNormalizer {
	return &SeniorityNormalizer{}
}

func (n *SeniorityNormalizer) Normalize(job *domain.Job) {
	if job.Seniority != "" {
		// Clean up existing seniority if it's set in provider but non-standard
		switch strings.ToLower(job.Seniority) {
		case "junior", "jr", "jr.":
			job.Seniority = "Junior"
			return
		case "mid", "middle", "pleno", "pl":
			job.Seniority = "Mid"
			return
		case "senior", "sr", "sr.":
			job.Seniority = "Senior"
			return
		case "staff":
			job.Seniority = "Staff"
			return
		case "principal":
			job.Seniority = "Principal"
			return
		case "lead", "lead engineer", "tech lead":
			job.Seniority = "Lead"
			return
		}
	}

	title := " " + strings.ToLower(job.Title) + " "

	// Tech Lead / Lead
	if strings.Contains(title, "tech lead") || strings.Contains(title, " lead ") || strings.Contains(title, " lilder ") || strings.Contains(title, " líder ") {
		job.Seniority = "Lead"
		return
	}

	// Principal
	if strings.Contains(title, " principal ") {
		job.Seniority = "Principal"
		return
	}

	// Staff
	if strings.Contains(title, " staff ") {
		job.Seniority = "Staff"
		return
	}

	// Senior
	if strings.Contains(title, " senior ") || strings.Contains(title, " sênior ") || strings.Contains(title, " sr. ") || strings.Contains(title, " sr ") || strings.Contains(title, " iii ") || strings.Contains(title, " iv ") || strings.HasSuffix(title, " iii ") || strings.HasSuffix(title, " iv ") {
		job.Seniority = "Senior"
		return
	}

	// Junior
	if strings.Contains(title, " junior ") || strings.Contains(title, " júnior ") || strings.Contains(title, " jr. ") || strings.Contains(title, " jr ") || strings.Contains(title, " i ") || strings.HasSuffix(title, " i ") {
		job.Seniority = "Junior"
		return
	}

	// Mid / Pleno
	if strings.Contains(title, " pleno ") || strings.Contains(title, " mid ") || strings.Contains(title, " ii ") || strings.HasSuffix(title, " ii ") {
		job.Seniority = "Mid"
		return
	}
}
