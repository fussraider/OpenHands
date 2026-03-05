package models

// WebClientConfig is the model representing the web client configuration settings
type WebClientConfig struct {
	AppMode             string                `json:"app_mode"`
	PosthogClientKey    *string               `json:"posthog_client_key"`
	GithubAppSlug       *string               `json:"github_app_slug"`
	FeatureFlags        WebClientFeatureFlags `json:"feature_flags"`
	ProvidersConfigured []string              `json:"providers_configured"`
	FaultyModels        []string              `json:"faulty_models"`
	UpdatedAt           string                `json:"updated_at"`
}

// WebClientFeatureFlags configures various toggle features
type WebClientFeatureFlags struct {
	EnableBilling   bool `json:"enable_billing"`
	HideLLMSettings bool `json:"hide_llm_settings"`
	EnableJira      bool `json:"enable_jira"`
	EnableJiraDC    bool `json:"enable_jira_dc"`
	EnableLinear    bool `json:"enable_linear"`
}
