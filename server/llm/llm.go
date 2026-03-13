package llm

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"openhands-go/server/config"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type LLMService struct {
	config config.LLMConfig
	model  llms.Model
}

func NewLLMService(cfg config.LLMConfig) (*LLMService, error) {
	var model llms.Model
	var err error

	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, errors.New("error authenticating with the LLM provider: please configure LLM_API_KEY")
	}

	opts := []openai.Option{
		openai.WithToken(cfg.APIKey),
		openai.WithModel(cfg.Model),
	}
	if cfg.BaseURL != "" {
		opts = append(opts, openai.WithBaseURL(cfg.BaseURL))
	}
	// Apply enhanced config
	// langchaingo openai provider doesn't expose generic options for temperature/top_p in New(),
	// they are usually call options. But some providers allow default options.
	// Checking langchaingo/llms/openai source (from memory/knowledge), it doesn't support setting default temp/top_p in New().
	// We have to pass them in GenerateContent.
	// However, we can store them in LLMService and apply them in CompleteWithTools.

	if cfg.BaseURL != "" {
		slog.Debug("Rewrote openhands model URL", "model", cfg.Model, "base_url", cfg.BaseURL)
	}

	model, err = openai.New(opts...)
	if err != nil {
		slog.Debug("Error getting model info", "error", err)
		return nil, err
	}

	slog.Debug("Got model info from litellm proxy", "model", cfg.Model)
	slog.Debug("LLM: model supports function calling", "model", cfg.Model)

	return &LLMService{
		config: cfg,
		model:  model,
	}, nil
}

type Message struct {
	Role       string          `json:"role"`
	Content    string          `json:"content"`
	ToolCalls  []llms.ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string          `json:"tool_call_id,omitempty"` // For tool output
}

type ToolCallResponse struct {
	Content   string
	ToolCalls []llms.ToolCall
}

func (s *LLMService) Complete(ctx context.Context, messages []Message) (string, error) {
	resp, err := s.CompleteWithTools(ctx, messages, nil)
	if err != nil {
		return "", err
	}
	return resp.Content, nil
}

func (s *LLMService) CompleteWithTools(ctx context.Context, messages []Message, tools []llms.Tool) (*ToolCallResponse, error) {
	if s.model == nil {
		return nil, errors.New("LLM service is not initialized")
	}

	content := []llms.MessageContent{}
	for _, msg := range messages {
		var role llms.ChatMessageType
		if msg.Role == "system" {
			role = llms.ChatMessageTypeSystem
		} else if msg.Role == "assistant" {
			role = llms.ChatMessageTypeAI
		} else if msg.Role == "tool" {
			role = llms.ChatMessageTypeTool
		} else {
			role = llms.ChatMessageTypeHuman
		}

		parts := []llms.ContentPart{}

		if msg.Role == "tool" {
			// For tool response, we use ToolCallResponse part
			parts = append(parts, llms.ToolCallResponse{
				ToolCallID: msg.ToolCallID,
				Name:       "", // Name is often optional or inferred by ID? Langchaingo might require it?
				Content:    msg.Content,
			})
		} else {
			if msg.Content != "" {
				parts = append(parts, llms.TextPart(msg.Content))
			}

			// Create parts for tool calls
			for _, tc := range msg.ToolCalls {
				parts = append(parts, tc)
			}
		}

		content = append(content, llms.MessageContent{
			Role:  role,
			Parts: parts,
		})
	}

	opts := []llms.CallOption{}
	if len(tools) > 0 {
		opts = append(opts, llms.WithTools(tools))
	}

	// Apply config options
	if s.config.Temperature != 0 {
		opts = append(opts, llms.WithTemperature(s.config.Temperature))
	}
	if s.config.TopP != 0 {
		opts = append(opts, llms.WithTopP(s.config.TopP))
	}
	if s.config.MaxOutputTokens != 0 {
		opts = append(opts, llms.WithMaxTokens(s.config.MaxOutputTokens))
	}

	completion, err := s.model.GenerateContent(ctx, content, opts...)
	if err != nil {
		return nil, err
	}

	if len(completion.Choices) == 0 {
		slog.Debug("No completion messages!")
		return nil, fmt.Errorf("no choices in LLM response")
	}

	choice := completion.Choices[0]
	slog.Debug("Cost calculation not supported for this model.", "model", s.config.Model)

	return &ToolCallResponse{
		Content:   choice.Content,
		ToolCalls: choice.ToolCalls,
	}, nil
}
