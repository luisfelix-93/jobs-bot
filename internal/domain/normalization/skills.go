package normalization

import (
	"regexp"
	// "strings"

	"jobs-bot/internal/domain"
)

// SkillsExtractor extracts technical skills from job titles and descriptions.
type SkillsExtractor struct {
	skillsMap map[string]*regexp.Regexp
}

func NewSkillsExtractor() *SkillsExtractor {
	skills := []string{
		"Go", "Golang", "Python", "Ruby", "Java", "Rust",
		"TypeScript", "JavaScript", "React", "Vue", "Angular", "Node.js", "Next.js",
		"Docker", "Kubernetes", "AWS", "GCP", "Azure", "Terraform", "Ansible", "Jenkins",
		"CI/CD", "MongoDB", "PostgreSQL", "MySQL", "Redis", "Elasticsearch",
		"Prometheus", "Grafana", "Linux", "Git", "REST", "GraphQL", "gRPC", "SRE", "DevOps", "Bash",
		"SQL Server", "ASP.NET",
	}

	skillsMap := make(map[string]*regexp.Regexp)
	for _, skill := range skills {
		pattern := `(?i)\b` + regexp.QuoteMeta(skill) + `\b`
		skillsMap[skill] = regexp.MustCompile(pattern)
	}

	// Custom patterns for special chars
	skillsMap[".NET"] = regexp.MustCompile(`(?i)(?:\b|_)` + regexp.QuoteMeta(".NET") + `\b`)
	skillsMap["C#"] = regexp.MustCompile(`(?i)(?:\b|_)C#(?:\b|_)`)

	return &SkillsExtractor{
		skillsMap: skillsMap,
	}
}

func (n *SkillsExtractor) Normalize(job *domain.Job) {
	text := job.Title + " " + job.FullDescription
	var extracted []string
	seen := make(map[string]bool)

	for skill, rx := range n.skillsMap {
		if rx.MatchString(text) {
			canonical := skill
			if canonical == "Golang" {
				canonical = "Go"
			}
			if !seen[canonical] {
				seen[canonical] = true
				extracted = append(extracted, canonical)
			}
		}
	}

	job.Skills = extracted
}
