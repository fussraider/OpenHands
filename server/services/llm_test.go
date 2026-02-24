package services

import (
	"openhands-go/server/config"
	"openhands-go/server/llm"
	"testing"
)

func TestLLMService(t *testing.T) {
	// Mock config
	cfg := config.LLMConfig{
		Model: "gpt-mock",
		// No API Key -> Mock Response
	}

	service := llm.NewLLMService(cfg)

	messages := []llm.Message{
		{Role: "user", Content: "Hello"},
	}

	resp, err := service.Complete(messages)
	if err != nil {
		t.Fatalf("Complete failed: %v", err)
	}

	if resp == "" {
		t.Error("Response is empty")
	}

	// Expect mock response
	expected := "This is a mock response from the Go backend LLM service."
	if resp != expected {
		t.Errorf("Unexpected response: got %q, want %q", resp, expected)
	}
}
