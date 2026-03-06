package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
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
	Github        GithubConfig  `toml:"github"`
	Security      SecurityConfig `toml:"security"`
}

type SecurityConfig struct {
	SecurityAnalyzer string `toml:"security_analyzer"` // "llm" or "invariant"
	ConfirmationMode bool   `toml:"confirmation_mode"`
}

type LLMConfig struct {
	Model             string  `toml:"model"`
	APIKey            string  `toml:"api_key"`
	BaseURL           string  `toml:"base_url"`
	NumRetries        int     `toml:"num_retries"`
	Timeout           int     `toml:"timeout"`
	Temperature       float64 `toml:"temperature"`
	TopP              float64 `toml:"top_p"`
	MaxInputTokens    int     `toml:"max_input_tokens"`
	MaxOutputTokens   int     `toml:"max_output_tokens"`
	CustomLLMProvider string  `toml:"custom_llm_provider"`
}

type AgentConfig struct {
	Name             string `toml:"name"`
	EnableHistoryTruncation bool `toml:"enable_history_truncation"`
	MaxEvents        int    `toml:"max_events"` // Simple condenser threshold
}

type SandboxConfig struct {
	Runtime string `toml:"runtime"`
}

type ServerConfig struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type GithubConfig struct {
	ClientID     string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
	RedirectURL  string `toml:"redirect_url"`
}

var AppConfig *Config

func LoadConfig() error {
	// Attempt to load .env file if it exists, ignoring errors if it doesn't
	_ = godotenv.Load()

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

	// Load from global ~/.openhands/config.toml
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		globalConfigPath := filepath.Join(homeDir, ".openhands", "config.toml")
		data, err := os.ReadFile(globalConfigPath)
		if err == nil {
			if err := toml.Unmarshal(data, AppConfig); err != nil {
				// Don't fail if we can't parse global config, just ignore it or log
			}
		}
	}

	// Load from local config.toml
	data, err := os.ReadFile("config.toml")
	if err == nil {
		// We unmarshal again into AppConfig, effectively deep-merging local into global
		if err := toml.Unmarshal(data, AppConfig); err != nil {
			return err
		}
		slog.Debug("Loaded configuration", "config_file", "config.toml")
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
	if numRetries := os.Getenv("LLM_NUM_RETRIES"); numRetries != "" {
		if n, err := strconv.Atoi(numRetries); err == nil {
			AppConfig.LLM.NumRetries = n
		}
	}
	if timeout := os.Getenv("LLM_TIMEOUT"); timeout != "" {
		if n, err := strconv.Atoi(timeout); err == nil {
			AppConfig.LLM.Timeout = n
		}
	}
	if temperature := os.Getenv("LLM_TEMPERATURE"); temperature != "" {
		if f, err := strconv.ParseFloat(temperature, 64); err == nil {
			AppConfig.LLM.Temperature = f
		}
	}
	if sandboxRuntime := os.Getenv("SANDBOX_RUNTIME"); sandboxRuntime != "" {
		AppConfig.Sandbox.Runtime = sandboxRuntime
	}
	if fileStorePath := os.Getenv("FILE_STORE_PATH"); fileStorePath != "" {
		AppConfig.FileStorePath = fileStorePath
	}
	if ghClientID := os.Getenv("GITHUB_CLIENT_ID"); ghClientID != "" {
		AppConfig.Github.ClientID = ghClientID
	}
	if ghClientSecret := os.Getenv("GITHUB_CLIENT_SECRET"); ghClientSecret != "" {
		AppConfig.Github.ClientSecret = ghClientSecret
	}

	// Default to 0.0.0.0 if running in Docker (checking widely used convention)
	if os.Getenv("RUN_AS_OPENHANDS") == "true" && AppConfig.Server.Host == "127.0.0.1" {
		AppConfig.Server.Host = "0.0.0.0"
	}

	return nil
}
