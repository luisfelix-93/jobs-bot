package normalization

import (
	"regexp"
	"strconv"
	"strings"

	"jobs-bot/internal/domain"
)

// SalaryNormalizer parses salary ranges and values from titles and descriptions.
type SalaryNormalizer struct {
	rangeRegex  *regexp.Regexp
	singleRegex *regexp.Regexp
}

func NewSalaryNormalizer() *SalaryNormalizer {
	return &SalaryNormalizer{
		// Match patterns like: $120k-$150k, USD 80k to 100k, 120.000 - 150.000, $120,000 - $150,000
		// NOTE: Comma/dot-separated formats MUST come before \d{2,3} to prevent partial matches (e.g. "100" instead of "100,000")
		rangeRegex:  regexp.MustCompile(`(?i)(USD|EUR|GBP|BRL|\$|€|£)?\s*(\d{1,3},\d{3}|\d{1,3}\.\d{3}|\d{2,3}(?:\.\d+)?)\s*(k)?\s*(?:\-|\bto\b)\s*(?:USD|EUR|GBP|BRL|\$|€|£)?\s*(\d{1,3},\d{3}|\d{1,3}\.\d{3}|\d{2,3}(?:\.\d+)?)\s*(k)?`),
		singleRegex: regexp.MustCompile(`(?i)(USD|EUR|GBP|BRL|\$|€|£)\s*(\d{1,3},\d{3}|\d{1,3}\.\d{3}|\d{2,3}(?:\.\d+)?)\s*(k)?`),
	}
}

func (n *SalaryNormalizer) Normalize(job *domain.Job) {
	if job.SalaryMin > 0 {
		return
	}

	text := job.Title + "\n" + job.FullDescription

	if matches := n.rangeRegex.FindStringSubmatch(text); len(matches) > 0 {
		minVal := parseValue(matches[2], matches[3])
		maxVal := parseValue(matches[4], matches[5])
		currency := detectCurrency(matches[1], text)

		job.SalaryMin = minVal
		job.SalaryMax = maxVal
		job.SalaryCurrency = currency
		return
	}

	if matches := n.singleRegex.FindStringSubmatch(text); len(matches) > 0 {
		val := parseValue(matches[2], matches[3])
		currency := detectCurrency(matches[1], text)

		job.SalaryMin = val
		job.SalaryMax = val
		job.SalaryCurrency = currency
		return
	}
}

func parseValue(valStr, kStr string) float64 {
	// Remove comma or dots (for thousands separators)
	valStr = strings.ReplaceAll(valStr, ",", "")
	// If it has dot but it is followed by 3 digits (e.g. 120.000), replace dot too
	if strings.Contains(valStr, ".") {
		parts := strings.Split(valStr, ".")
		if len(parts) == 2 && len(parts[1]) == 3 {
			valStr = strings.ReplaceAll(valStr, ".", "")
		}
	}
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return 0
	}
	if strings.ToLower(kStr) == "k" {
		val = val * 1000
	}
	return val
}

func detectCurrency(symbol, fullText string) string {
	symbol = strings.ToUpper(strings.TrimSpace(symbol))
	switch symbol {
	case "$", "USD":
		return "USD"
	case "€", "EUR":
		return "EUR"
	case "£", "GBP":
		return "GBP"
	case "BRL":
		return "BRL"
	}
	if strings.Contains(strings.ToLower(fullText), "usd") || strings.Contains(fullText, "$") {
		return "USD"
	}
	return "USD"
}
