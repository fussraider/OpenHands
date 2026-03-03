package models

// WebClientConfig is the model representing the web client configuration settings
type WebClientConfig struct {
	AppMode             string                `json:"appMode"`
	PosthogClientKey    *string               `json:"posthogClientKey"`
	GithubAppSlug       *string               `json:"githubAppSlug"`
	FeatureFlags        WebClientFeatureFlags `json:"featureFlags"`
	ProvidersConfigured []string              `json:"providersConfigured"`
	FaultyModels        []string              `json:"faultyModels"`
	UpdatedAt           string                `json:"updatedAt"`
}

// WebClientFeatureFlags configures various toggle features
type WebClientFeatureFlags struct {
	EnableBilling   bool `json:"enableBilling"`
	HideLLMSettings bool `json:"hideLLMSettings"`
	EnableJira      bool `json:"enableJira"`
	EnableJiraDC    bool `json:"enableJiraDC"`
	EnableLinear    bool `json:"enableLinear"`
}
