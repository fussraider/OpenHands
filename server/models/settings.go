package models

type Settings struct {
	Language                     string  `json:"language,omitempty"`
	Agent                        string  `json:"agent,omitempty"`
	MaxIterations                int     `json:"max_iterations,omitempty"`
	SecurityAnalyzer             string  `json:"security_analyzer,omitempty"`
	ConfirmationMode             bool    `json:"confirmation_mode,omitempty"`
	LLMModel                     string  `json:"llm_model,omitempty"`
	LLMAPIKey                    string  `json:"llm_api_key,omitempty"`
	LLMBaseURL                   string  `json:"llm_base_url,omitempty"`
	UserVersion                  int     `json:"user_version,omitempty"`
	RemoteRuntimeResourceFactor  int     `json:"remote_runtime_resource_factor,omitempty"`
	EnableDefaultCondenser       bool    `json:"enable_default_condenser"`
	EnableSoundNotifications     bool    `json:"enable_sound_notifications"`
	EnableProactiveStarter       bool    `json:"enable_proactive_conversation_starters"`
	EnableSolvabilityAnalysis    bool    `json:"enable_solvability_analysis"`
	UserConsentsToAnalytics      bool    `json:"user_consents_to_analytics,omitempty"`
	SandboxBaseContainerImage    string  `json:"sandbox_base_container_image,omitempty"`
	SandboxRuntimeContainerImage string  `json:"sandbox_runtime_container_image,omitempty"`
	SearchAPIKey                 string  `json:"search_api_key,omitempty"`
	SandboxAPIKey                string  `json:"sandbox_api_key,omitempty"`
	MaxBudgetPerTask             float64 `json:"max_budget_per_task,omitempty"`
	CondenserMaxSize             int     `json:"condenser_max_size,omitempty"`
	Email                        string  `json:"email,omitempty"`
	EmailVerified                bool    `json:"email_verified,omitempty"`
	GitUserName                  string  `json:"git_user_name,omitempty"`
	GitUserEmail                 string  `json:"git_user_email,omitempty"`
}

func DefaultSettings() Settings {
	return Settings{
		Language:                  "en",
		Agent:                     "CodeActAgent",
		LLMModel:                  "gpt-4",
		EnableDefaultCondenser:    true,
		EnableProactiveStarter:    true,
		EnableSolvabilityAnalysis: true,
	}
}
