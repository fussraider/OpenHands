package llm

import (
	"context"
	"fmt"
	"openhands-go/server/config"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type LLMService struct {
	config config.LLMConfig
	model  llms.Model
}

func NewLLMService(cfg config.LLMConfig) (*LLMService, error) {
	// If API Key is empty, we might want to return a mock or error.
	// For now, let's assume OpenAI compatible if Key/BaseURL are present.

	// If no API Key, we can't really initialize a real LLM.
	// We'll handle this in Complete or return a mock implementation of llms.Model if needed.
	// But `langchaingo` generally requires valid config.

	var model llms.Model
	var err error

	if cfg.APIKey != "" {
		opts := []openai.Option{
			openai.WithToken(cfg.APIKey),
			openai.WithModel(cfg.Model),
		}
		if cfg.BaseURL != "" {
			opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
		}

		model, err = openai.New(opts...)
		if err != nil {
			return nil, err
		}
	}

	return &LLMService{
		config: cfg,
		model:  model,
	}, nil
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (s *LLMService) Complete(ctx context.Context, messages []Message) (string, error) {
	// Mock implementation if no API Key (for testing/dev)
	if s.model == nil {
		return "This is a mock response from the Go backend LLM service (langchaingo integration pending config).", nil
	}

	// Convert messages to langchaingo format
	content := []llms.MessageContent{}
	for _, msg := range messages {
		role := llms.ChatMessageTypeHuman
		if msg.Role == "system" {
			role = llms.ChatMessageTypeSystem
		} else if msg.Role == "assistant" {
			role = llms.ChatMessageTypeAI
		}

		content = append(content, llms.TextParts(role, msg.Content))
	}

	completion, err := s.model.GenerateContent(ctx, content)
	if err != nil {
		return "", err
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no choices in LLM response")
	}

	return completion.Choices[0].Content, nil
}
