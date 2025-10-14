package config

// Provider names
const (
	ProviderOpenAI           = "OpenAI"
	ProviderAnthropic        = "Anthropic"
	ProviderGoogle           = "Google"
	ProviderXAI              = "xAI"
	ProviderOpenAICompatible = "OpenAI-Compatible"
)

type Config struct {
	CurrentProvider string `json:"current_provider"`
	CurrentModel    string `json:"current_model"`
	BaseURL         string `json:"base_url,omitempty"` // For OpenAI-Compatible providers
}

// List of available models for each provider
var ProviderModels = map[string][]string{
	ProviderOpenAI: {
		// GPT-5 Models
		"gpt-5",
		"gpt-5-mini",
		"gpt-5-nano",
		// O Models
		"o4-mini",
		"o3",
		"o3-mini",
		"o1",
		"o1-mini",
		// GPT-4 Models
		"gpt-4.1",
		"gpt-4.1-mini",
		"gpt-4.1-nano",
		"gpt-4o",
		"gpt-4o-mini",
		"gpt-4-turbo",
		"gpt-4",
		"gpt-3.5-turbo",
	},
	ProviderAnthropic: {
		"claude-opus-4-1",
		"claude-opus-4-0",
		"claude-sonnet-4-5",
		"claude-sonnet-4-0",
		"claude-3-7-sonnet-latest",
		"claude-3-5-haiku-latest",
	},
	ProviderGoogle: {
		"gemini-2.5-pro",
		"gemini-2.5-flash",
		"gemini-2.5-flash-lite",
		"gemini-2.0-flash",
		"gemini-2.0-flash-lite",
	},
	ProviderXAI: {
		"grok-code-fast-1",
		"grok-4-fast-reasoning",
		"grok-4-fast-non-reasoning",
		"grok-3-mini",
		"grok-3",
	},
	ProviderOpenAICompatible: {
		// No predefined models - user will enter custom model name
	},
}

// Return a list of all providers
func GetProviders() []string {
	return []string{ProviderOpenAI, ProviderAnthropic, ProviderGoogle, ProviderXAI, ProviderOpenAICompatible}
}

func NewConfig() *Config {
	return &Config{
		CurrentProvider: "",
		CurrentModel:    "",
	}
}

// Check if the config has a provider and API key set
func (c *Config) IsConfigured() (bool, []string) {
	missing := []string{}
	ready := true

	// Check for missing fields
	if c.CurrentProvider == "" {
		missing = append(missing, "provider")
		ready = false
	}
	if c.CurrentModel == "" {
		missing = append(missing, "model")
		ready = false
	}
	if c.CurrentProvider != ProviderOpenAICompatible {
		// Check keyring for API key
		hasKey, err := HasAPIKeyInKeyring(c.CurrentProvider)
		if err != nil || !hasKey {
			missing = append(missing, "API key")
			ready = false
		}
	} else {
		if c.BaseURL == "" {
			missing = append(missing, "base URL")
			ready = false
		}
	}
	return ready, missing
}
	
// Set configured provider and model
func (c *Config) SetProvider(provider, model string) {
	c.CurrentProvider = provider
	c.CurrentModel = model
}
