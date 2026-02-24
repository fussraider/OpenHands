package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"openhands-go/server/config"
	"time"
)

type LLMService struct {
	config config.LLMConfig
	client *http.Client
}

func NewLLMService(cfg config.LLMConfig) *LLMService {
	return &LLMService{
		config: cfg,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type CompletionResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func (s *LLMService) Complete(messages []Message) (string, error) {
	// Mock implementation if no API Key (for testing/dev)
	if s.config.APIKey == "" {
		return "This is a mock response from the Go backend LLM service.", nil
	}

	reqBody := CompletionRequest{
		Model:    s.config.Model,
		Messages: messages,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := s.config.BaseURL + "/chat/completions" // Assumption for OpenAI-compatible API
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("LLM API returned status: %d", resp.StatusCode)
	}

	var completionResp CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&completionResp); err != nil {
		return "", err
	}

	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no choices in LLM response")
	}

	return completionResp.Choices[0].Message.Content, nil
}
