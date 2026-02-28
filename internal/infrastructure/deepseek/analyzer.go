package deepseekai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"jobs-bot/internal/domain"

	deepseek "github.com/cohesion-org/deepseek-go"
)

type Analyzer struct {
	client *deepseek.Client
}

func NewAnalyzer(apiKey string) *Analyzer {
	client := deepseek.NewClient(apiKey)
	return &Analyzer{client: client}
}

func (a *Analyzer) Analyze(resumeContent, jobDescription string) (*domain.AIAnalysis, error) {
	prompt := fmt.Sprintf(`Você é um analista de vagas de emprego.
Compare o currículo abaixo com a descrição da vaga e avalie a compatibilidade.

CURRÍCULO:
%s

VAGA:
%s

Retorne APENAS um JSON válido com:
{
  "score": 0-100,
  "strengths": ["competência que o candidato tem e a vaga requer"],
  "gaps": ["competência que a vaga requer e o candidato não tem"],
  "recommendation": "apply" | "review" | "skip",
  "summary": "análise em 2-3 frases"
}`, resumeContent, jobDescription)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	request := &deepseek.ChatCompletionRequest{
		Model: "deepseek-chat",
		Messages: []deepseek.ChatCompletionMessage{
			{
				Role:    deepseek.ChatMessageRoleSystem,
				Content: "Você é um analista de vagas. Responda APENAS com JSON válido, sem markdown ou texto extra.",
			},
			{
				Role:    deepseek.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		JSONMode:    true,
		Temperature: 0.3,
		MaxTokens:   500,
	}

	resp, err := a.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar DeepSeek API: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("DeepSeek retornou resposta sem choices")
	}

	content := resp.Choices[0].Message.Content

	content = strings.TrimPrefix(content, "```json\n")
	content = strings.TrimPrefix(content, "```\n")
	content = strings.TrimSuffix(content, "\n```")
	content = strings.TrimSpace(content)

	var analysis domain.AIAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		return nil, fmt.Errorf("erro ao parsear JSON do DeepSeek: %w (conteúdo: %s)", err, content)
	}

	analysis.Source = "deepseek"

	return &analysis, nil
}