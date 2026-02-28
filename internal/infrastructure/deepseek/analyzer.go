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

	// Sanitize possible Markdown code fences in a tolerant way
	// 1. Trim surrounding whitespace (including CRLF/space variations).
	content = strings.TrimSpace(content)

	// 2. Remove leading ``` fence, with or without language tag, and with any newline style.
	if strings.HasPrefix(content, "```") {
		if idx := strings.Index(content, "\n"); idx != -1 {
			// Drop the first line (the fence), keep the rest.
			content = content[idx+1:]
		} else {
			// No newline found; just strip the fence markers and any language tag.
			content = strings.TrimPrefix(content, "```json")
			content = strings.TrimPrefix(content, "```")
		}
	}

	// 3. Trim again to normalize after removing the leading fence.
	content = strings.TrimSpace(content)

	// 4. Remove trailing ``` fence, regardless of preceding newline style.
	if strings.HasSuffix(content, "```") {
		content = content[:len(content)-3]
	}

	// 5. Final trim before JSON parsing.
	content = strings.TrimSpace(content)
	var analysis domain.AIAnalysis
	if err := json.Unmarshal([]byte(content), &analysis); err != nil {
		contentPreview := content
		const maxPreviewLen = 200
		if len(contentPreview) > maxPreviewLen {
			contentPreview = contentPreview[:maxPreviewLen]
		}
		return nil, fmt.Errorf("erro ao parsear JSON do DeepSeek: %w (conteúdo truncado, %d bytes, primeiros %d bytes: %s)", err, len(content), len(contentPreview), contentPreview)
	}

	analysis.Source = "deepseek"

	return &analysis, nil
}