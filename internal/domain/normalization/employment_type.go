package normalization

import (
	"strings"

	"jobs-bot/internal/domain"
)

// EmploymentTypeNormalizer parses and normalizes the job's employment type.
type EmploymentTypeNormalizer struct{}

func NewEmploymentTypeNormalizer() *EmploymentTypeNormalizer {
	return &EmploymentTypeNormalizer{}
}

func (n *EmploymentTypeNormalizer) Normalize(job *domain.Job) {
	if job.EmploymentType != "" {
		switch strings.ToLower(job.EmploymentType) {
		case "full time", "full-time", "fulltime", "clt", "tempo integral", "permanent":
			job.EmploymentType = "FullTime"
			return
		case "contract", "contractor", "pj", "prestador de serviços", "freelance", "freelancer":
			job.EmploymentType = "Contract"
			return
		case "part time", "part-time", "parttime", "meio período":
			job.EmploymentType = "PartTime"
			return
		}
	}

	title := " " + strings.ToLower(job.Title) + " "
	desc := " " + strings.ToLower(job.FullDescription) + " "

	// Check Title
	if strings.Contains(title, " clt ") || strings.Contains(title, " full-time ") || strings.Contains(title, " full time ") || strings.Contains(title, " permanente ") {
		job.EmploymentType = "FullTime"
		return
	}
	if strings.Contains(title, " pj ") || strings.Contains(title, " contractor ") || strings.Contains(title, " contract ") || strings.Contains(title, " freelance ") || strings.Contains(title, " freelancer ") {
		job.EmploymentType = "Contract"
		return
	}
	if strings.Contains(title, " part-time ") || strings.Contains(title, " part time ") {
		job.EmploymentType = "PartTime"
		return
	}

	// Check Description
	if strings.Contains(desc, " clt ") || strings.Contains(desc, " fulltime ") || strings.Contains(desc, " full-time ") || strings.Contains(desc, " full time ") || strings.Contains(desc, " permanente ") {
		job.EmploymentType = "FullTime"
		return
	}
	if strings.Contains(desc, " pj ") || strings.Contains(desc, " contractor ") || strings.Contains(desc, " prestador ") || strings.Contains(desc, " freelance ") || strings.Contains(desc, " freelancer ") {
		job.EmploymentType = "Contract"
		return
	}
}
