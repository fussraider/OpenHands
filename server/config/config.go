package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type AppMode string

const (
	AppModeOpenHands AppMode = "oss"
	AppModeSaaS      AppMode = "saas"
)

type Config struct {
	AppMode AppMode       `toml:"app_mode"`
	LLM     LLMConfig     `toml:"llm"`
	Agent   AgentConfig   `toml:"agent"`
	Sandbox SandboxConfig `toml:"sandbox"`
	Server  ServerConfig  `toml:"server"`
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

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

var AppConfig *Config

func LoadConfig() error {
	// Default config
	AppConfig = &Config{
		AppMode: AppModeOpenHands,
		LLM: LLMConfig{
			Model: "gpt-4",
		},
		Agent: AgentConfig{
			Name: "CodeActAgent",
		},
		Server: ServerConfig{
			Host: "127.0.0.1",
			Port: 3000,
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
