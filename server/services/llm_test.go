package services

import (
	"openhands-go/server/config"
	"openhands-go/server/llm"
	"testing"
)

func TestLLMConfigIntegration(t *testing.T) {
	cfg := config.LLMConfig{
		Model:       "gpt-4",
		APIKey:      "mock-key",
		Temperature: 0.7,
		TopP:        0.9,
		MaxOutputTokens: 100,
	}

	service, err := llm.NewLLMService(cfg)
	if err != nil {
		t.Fatalf("Failed to create LLM service: %v", err)
	}

	if service == nil {
		t.Fatal("LLM service is nil")
	}

	// We can't verify internal options of langchaingo easily without digging into private fields
	// But ensuring compilation and initialization passes is a good first step.
}
