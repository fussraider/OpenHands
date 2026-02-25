package services

import (
	"context"
	"openhands-go/server/config"
	"openhands-go/server/llm"
	"testing"
)

func TestLLMService(t *testing.T) {
	// Mock config
	cfg := config.LLMConfig{
		Model: "gpt-mock",
		// No API Key -> Mock Response (handled by our wrapper, not langchaingo)
	}

	service, err := llm.NewLLMService(cfg)
	if err != nil {
		t.Fatalf("NewLLMService failed: %v", err)
	}

	messages := []llm.Message{
		{Role: "user", Content: "Hello"},
	}

	resp, err := service.Complete(context.Background(), messages)
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	if resp == "" {
		t.Error("Response is empty")
	}

	// Expect mock response
	expected := "This is a mock response from the Go backend LLM service (langchaingo integration pending config)."
	if resp != expected {
		t.Errorf("Unexpected response: got %q, want %q", resp, expected)
	}
}
