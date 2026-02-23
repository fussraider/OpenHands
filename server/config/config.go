package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	LLM     LLMConfig     `toml:"llm"`
	Agent   AgentConfig   `toml:"agent"`
	Sandbox SandboxConfig `toml:"sandbox"`
}

type LLMConfig struct {
	Model   string `toml:"model"`
	APIKey  string `toml:"api_key"`
	BaseURL string `toml:"base_url"`
}

type AgentConfig struct {
	Name string `toml:"name"`
}

type SandboxConfig struct {
	Runtime string `toml:"runtime"`
}

var AppConfig *Config

func LoadConfig() error {
	// Default config
	AppConfig = &Config{
		LLM: LLMConfig{
			Model: "gpt-4",
		},
		Agent: AgentConfig{
			Name: "CodeActAgent",
		},
	}

	data, err := os.ReadFile("config.toml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return toml.Unmarshal(data, AppConfig)
}
