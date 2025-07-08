package config

import (
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {
	// Save original env vars
	originalProvider := os.Getenv("DEFAULT_AI_PROVIDER")
	originalOpenAIKey := os.Getenv("OPENAI_API_KEY")
	originalAnthropicKey := os.Getenv("ANTHROPIC_API_KEY")

	// Clean up after test
	defer func() {
		os.Setenv("DEFAULT_AI_PROVIDER", originalProvider)
		os.Setenv("OPENAI_API_KEY", originalOpenAIKey)
		os.Setenv("ANTHROPIC_API_KEY", originalAnthropicKey)
	}()

	// Test with environment variables set
	os.Setenv("DEFAULT_AI_PROVIDER", "openai")
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	os.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")

	InitConfig()

	if AppConfig == nil {
		t.Fatal("AppConfig should not be nil after initialization")
	}

	if AppConfig.DefaultProvider != ProviderOpenAI {
		t.Errorf("Expected default provider to be openai, got %s", AppConfig.DefaultProvider)
	}

	if AppConfig.OpenAI.APIKey != "test-openai-key" {
		t.Errorf("Expected OpenAI API key to be 'test-openai-key', got %s", AppConfig.OpenAI.APIKey)
	}

	if AppConfig.Anthropic.APIKey != "test-anthropic-key" {
		t.Errorf("Expected Anthropic API key to be 'test-anthropic-key', got %s", AppConfig.Anthropic.APIKey)
	}
}

func TestGetDefaultProvider(t *testing.T) {
	tests := []struct {
		envValue string
		expected AIProvider
	}{
		{"openai", ProviderOpenAI},
		{"anthropic", ProviderAnthropic},
		{"azure-openai", ProviderAzureOpenAI},
		{"invalid", ProviderAnthropic}, // Should default to Anthropic
		{"", ProviderAnthropic},        // Should default to Anthropic (as per getEnv default)
	}

	originalProvider := os.Getenv("DEFAULT_AI_PROVIDER")
	defer os.Setenv("DEFAULT_AI_PROVIDER", originalProvider)

	for _, test := range tests {
		os.Setenv("DEFAULT_AI_PROVIDER", test.envValue)
		result := getDefaultProvider()
		if result != test.expected {
			t.Errorf("For env value '%s', expected %s, got %s", test.envValue, test.expected, result)
		}
	}
}

func TestGetEnv(t *testing.T) {
	tests := []struct {
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{"TEST_KEY", "default", "env_value", "env_value"},
		{"NONEXISTENT_KEY", "default", "", "default"},
	}

	for _, test := range tests {
		if test.envValue != "" {
			os.Setenv(test.key, test.envValue)
			defer os.Unsetenv(test.key)
		}

		result := getEnv(test.key, test.defaultValue)
		if result != test.expected {
			t.Errorf("For key '%s', expected '%s', got '%s'", test.key, test.expected, result)
		}
	}
}

func TestGetEndpointURL(t *testing.T) {
	// Initialize AppConfig for testing
	AppConfig = &Config{
		AzureOpenAI: AzureOpenAIConfig{
			BaseURL:    "https://test.openai.azure.com",
			APIVersion: "2024-02-15-preview",
		},
	}

	tests := []struct {
		model    string
		expected string
	}{
		{"gpt-4.1", "https://test.openai.azure.com/openai/deployments/gpt-4.1/chat/completions?api-version=2024-02-15-preview"},
		{"o3", "https://test.openai.azure.com/openai/deployments/o3/chat/completions?api-version=2024-02-15-preview"},
		{"unknown-model", "https://test.openai.azure.com/openai/deployments/gpt-4.1/chat/completions?api-version=2024-02-15-preview"}, // Should fallback to gpt-4.1
	}

	for _, test := range tests {
		result := GetEndpointURL(test.model)
		if result != test.expected {
			t.Errorf("For model '%s', expected '%s', got '%s'", test.model, test.expected, result)
		}
	}
}
