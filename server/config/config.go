package config

import (
	"os"
	"strconv"

	"github.com/pelletier/go-toml/v2"
)

type AppMode string

const (
	AppModeOpenHands AppMode = "oss"
	AppModeSaaS      AppMode = "saas"
)

type Config struct {
	AppMode       AppMode       `toml:"app_mode"`
	FileStorePath string        `toml:"file_store_path"`
	LLM           LLMConfig     `toml:"llm"`
	Agent         AgentConfig   `toml:"agent"`
	Sandbox       SandboxConfig `toml:"sandbox"`
	Server        ServerConfig  `toml:"server"`
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
	if err == nil {
		if err := toml.Unmarshal(data, AppConfig); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	// Override with environment variables
	if host := os.Getenv("OPENHANDS_HOST"); host != "" {
		AppConfig.Server.Host = host
	}
	if port := os.Getenv("OPENHANDS_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			AppConfig.Server.Port = p
		}
	}
	if llmModel := os.Getenv("LLM_MODEL"); llmModel != "" {
		AppConfig.LLM.Model = llmModel
	}
	if llmAPIKey := os.Getenv("LLM_API_KEY"); llmAPIKey != "" {
		AppConfig.LLM.APIKey = llmAPIKey
	}
	if llmBaseURL := os.Getenv("LLM_BASE_URL"); llmBaseURL != "" {
		AppConfig.LLM.BaseURL = llmBaseURL
	}
	if sandboxRuntime := os.Getenv("SANDBOX_RUNTIME"); sandboxRuntime != "" {
		AppConfig.Sandbox.Runtime = sandboxRuntime
	}
	if fileStorePath := os.Getenv("FILE_STORE_PATH"); fileStorePath != "" {
		AppConfig.FileStorePath = fileStorePath
	}

	// Default to 0.0.0.0 if running in Docker (checking widely used convention)
	if os.Getenv("RUN_AS_OPENHANDS") == "true" && AppConfig.Server.Host == "127.0.0.1" {
		AppConfig.Server.Host = "0.0.0.0"
	}

	return nil
}
