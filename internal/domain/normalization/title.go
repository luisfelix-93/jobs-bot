package normalization

import (
	"regexp"
	"strings"

	"jobs-bot/internal/domain"
)

// TitleNormalizer cleans up the job title to make it standard and searchable.
type TitleNormalizer struct {
	remoteRegex      *regexp.Regexp
	parenthesesRegex *regexp.Regexp
	seniorityRegex   *regexp.Regexp
}

func NewTitleNormalizer() *TitleNormalizer {
	return &TitleNormalizer{
		// Match common remote/hybrid/workplace tags in titles
		remoteRegex:      regexp.MustCompile(`(?i)\s*[\-\/|([{\[\]\}]*\s*(100%|fully)?\s*(remote|remoto|hybrid|híbrido|on\-site|onsite|home office|anywhere|office|wfh)\s*[)\]}]*`),
		parenthesesRegex: regexp.MustCompile(`\s*[\(\[\{]\s*[\)\]\}]`),
		// Match seniority-level words (already captured by SeniorityNormalizer)
		seniorityRegex: regexp.MustCompile(`(?i)\b(senior|sênior|junior|júnior|pleno|mid[\-\s]?level|staff|principal|lead|sr\.?|jr\.?)\b`),
	}
}

func (n *TitleNormalizer) Normalize(job *domain.Job) {
	title := job.Title

	// Remove company name prefixes if present (e.g. "Google - Software Engineer")
	if job.Company != "" {
		companyLower := strings.ToLower(job.Company)
		titleLower := strings.ToLower(title)
		if strings.HasPrefix(titleLower, companyLower) {
			prefixLen := len(job.Company)
			if len(title) > prefixLen {
				title = title[prefixLen:]
				title = strings.TrimLeft(title, " -:|/\\")
			}
		}
	}

	// Clean tags
	title = n.remoteRegex.ReplaceAllString(title, " ")
	title = n.parenthesesRegex.ReplaceAllString(title, " ")
	// Strip seniority words (already captured separately)
	title = n.seniorityRegex.ReplaceAllString(title, " ")

	// Trim and squeeze whitespace
	title = strings.Join(strings.Fields(title), " ")
	title = strings.Trim(title, " -:|/,")

	if title == "" {
		job.NormalizedTitle = job.Title
	} else {
		job.NormalizedTitle = title
	}
}
