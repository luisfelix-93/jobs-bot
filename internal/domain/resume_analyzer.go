// NOVO ARQUIVO: luisfelix-93/jobs-bot/jobs-bot-c698ff856bf701e11b39469df507bdfca44d838b/internal/domain/resume_analyzer.go

package domain

import "strings"

// ResumeAnalysis armazena o resultado da comparação
type ResumeAnalysis struct {
	MatchPercentage float64
	MissingKeywords []string
	FoundKeywords   []string
}

// ResumeAnalyzer realiza a análise
type ResumeAnalyzer struct{}

func NewResumeAnalyzer() *ResumeAnalyzer {
	return &ResumeAnalyzer{}
}


func (r *ResumeAnalyzer) Analyze(resumeContent, jobDescriptionContent string, keywords []string) ResumeAnalysis {
	lowerResume := strings.ToLower(resumeContent)
	lowerJobDesc := strings.ToLower(jobDescriptionContent)

	var keywordsInJob []string
	var foundInResume []string
	var missingFromResume []string

	
	for _, kw := range keywords {
		if strings.Contains(lowerJobDesc, strings.ToLower(kw)) {
			keywordsInJob = append(keywordsInJob, kw)
		}
	}

	if len(keywordsInJob) == 0 {
		return ResumeAnalysis{MatchPercentage: 0} 
	}

	
	for _, jobKw := range keywordsInJob {
		if strings.Contains(lowerResume, strings.ToLower(jobKw)) {
			foundInResume = append(foundInResume, jobKw)
		} else {
			missingFromResume = append(missingFromResume, jobKw)
		}
	}

	// 3. Calcula a pontuação
	score := float64(len(foundInResume)) / float64(len(keywordsInJob)) * 100

	return ResumeAnalysis{
		MatchPercentage: score,
		MissingKeywords: missingFromResume,
		FoundKeywords:   foundInResume,
	}
}