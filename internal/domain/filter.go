package domain

import (
	"sort"
	"strings"
)

// NOTA: O struct 'ResumeAnalysis' está definido em 'resume_analyzer.go'.
// O 'FilteredJob' que eu criei antes não é mais necessário.

// JobFilter contém a lógica de negócio pura para filtrar e rankear vagas.
// (Versão simplificada sem AllowedLocations).
type JobFilter struct {
	PositiveKeywords []string
	NegativeKeywords []string
}

// NewJobFilter é o construtor para o nosso filtro.
func NewJobFilter(positive, negative []string) *JobFilter {
	return &JobFilter{
		PositiveKeywords: positive,
		NegativeKeywords: negative,
	}
}

// FilterAndRankJobs agora retorna apenas []domain.Job.
// Sua única responsabilidade é encontrar os jobs mais relevantes
// baseado nas palavras-chave encontradas na *descrição da vaga*.
func (f *JobFilter) FilterAndRankJobs(jobs []Job, limit int) []Job {
	// struct temporário para ajudar na ordenação
	type rankedJob struct {
		job   Job
		score int
	}

	var rankedJobs []rankedJob

	for _, job := range jobs {
		// Combinamos título e descrição para a busca
		fullText := strings.ToLower(job.Title + " " + job.FullDescription)

		// PASSO 1: Filtro de Palavras Negativas
		if f.containsNegativeKeyword(fullText) {
			continue // Pula esta vaga
		}

		// PASSO 2: Cálculo de Score
		// Apenas calcula o score baseado na descrição da vaga
		score := f.calculateScore(fullText)

		// Apenas adiciona vagas que pontuaram (score > 0)
		if score > 0 {
			rankedJobs = append(rankedJobs, rankedJob{
				job:   job,
				score: score,
			})
		}
	}

	// Ordena as vagas pelo score, da maior para a menor
	sort.Slice(rankedJobs, func(i, j int) bool {
		return rankedJobs[i].score > rankedJobs[j].score
	})

	// Prepara o retorno simples
	var bestJobs []Job
	for i := 0; i < len(rankedJobs) && i < limit; i++ {
		bestJobs = append(bestJobs, rankedJobs[i].job)
	}

	return bestJobs
}

// calculateScore apenas retorna um int
func (f *JobFilter) calculateScore(fullText string) int {
	score := 0
	for _, keyword := range f.PositiveKeywords {
		if strings.Contains(fullText, strings.ToLower(keyword)) {
			score++
		}
	}
	return score
}

// containsNegativeKeyword verifica a presença de palavras que desqualificam a vaga.
func (f *JobFilter) containsNegativeKeyword(fullText string) bool {
	for _, keyword := range f.NegativeKeywords {
		if strings.Contains(fullText, strings.ToLower(keyword)) {
			return true
		}
	}
	return false
}