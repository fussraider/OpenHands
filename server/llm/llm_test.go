package llm

import (
	"openhands-go/server/config"
	"testing"
)

func TestNewLLMServiceRequiresAPIKey(t *testing.T) {
	_, err := NewLLMService(config.LLMConfig{Model: "gpt-4"})
	if err == nil {
		t.Fatal("expected error when API key is missing")
	}
}
