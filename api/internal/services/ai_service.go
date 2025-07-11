package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"promptforge/internal/config"
	"promptforge/internal/models"
)

// AIService interface for all AI providers
type AIService interface {
	CallAI(messages []models.Message, temperature float64, maxTokens int, model string, provider config.AIProvider) (string, error)
}

// UnifiedAIService implements AIService for multiple providers
type UnifiedAIService struct {
	client *http.Client
}

func NewUnifiedAIService() *UnifiedAIService {
	return &UnifiedAIService{
		client: &http.Client{},
	}
}

// CallAI routes to the appropriate provider
func (s *UnifiedAIService) CallAI(messages []models.Message, temperature float64, maxTokens int, model string, provider config.AIProvider) (string, error) {
	switch provider {
	case config.ProviderOpenAI:
		return s.callOpenAI(messages, temperature, maxTokens, model)
	case config.ProviderAzureOpenAI:
		return s.callAzureOpenAI(messages, temperature, maxTokens, model)
	case config.ProviderAnthropic:
		return s.callAnthropic(messages, temperature, maxTokens, model)
	case config.ProviderOllama:
		return s.callOllama(messages, temperature, maxTokens, model)
	default:
		return "", fmt.Errorf("unsupported AI provider: %s", provider)
	}
}

// CallWithDefaultProvider uses the configured default provider
func (s *UnifiedAIService) CallWithDefaultProvider(messages []models.Message, temperature float64, maxTokens int, model string) (string, error) {
	return s.CallAI(messages, temperature, maxTokens, model, config.AppConfig.DefaultProvider)
}

func (s *UnifiedAIService) callOpenAI(messages []models.Message, temperature float64, maxTokens int, model string) (string, error) {
	if config.AppConfig.OpenAI.APIKey == "" {
		return "", fmt.Errorf("OpenAI API key not configured")
	}

	// Default to gpt-4 if no model specified
	if model == "" {
		model = "gpt-4"
	}

	requestBody := models.OpenAIRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
	}

	if maxTokens > 0 {
		requestBody.MaxTokens = maxTokens
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/chat/completions", strings.TrimSuffix(config.AppConfig.OpenAI.BaseURL, "/"))
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.OpenAI.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("OpenAI API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp models.OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", err
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func (s *UnifiedAIService) callAzureOpenAI(messages []models.Message, temperature float64, maxTokens int, model string) (string, error) {
	if config.AppConfig.AzureOpenAI.APIKey == "" {
		return "", fmt.Errorf("Azure OpenAI API key not configured")
	}

	endpoint := config.GetEndpointURL(model)

	// o3 model only supports temperature = 1
	if model == "o3" {
		temperature = 1
	}

	requestBody := models.OpenAIRequest{
		Messages:    messages,
		Temperature: temperature,
	}

	if maxTokens > 0 {
		// o3 model uses max_completion_tokens instead of max_tokens
		if model == "o3" {
			requestBody.MaxCompletionTokens = maxTokens
		} else {
			requestBody.MaxTokens = maxTokens
		}
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("api-key", config.AppConfig.AzureOpenAI.APIKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Azure OpenAI API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp models.OpenAIResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		return "", err
	}

	if len(openAIResp.Choices) == 0 {
		return "", fmt.Errorf("no response from Azure OpenAI")
	}

	return openAIResp.Choices[0].Message.Content, nil
}

func (s *UnifiedAIService) callAnthropic(messages []models.Message, temperature float64, maxTokens int, model string) (string, error) {
	if config.AppConfig.Anthropic.APIKey == "" {
		return "", fmt.Errorf("Anthropic API key not configured")
	}

	// Default to claude-3-5-sonnet if no model specified
	if model == "" {
		model = "claude-3-5-sonnet-20241022"
	}

	// Anthropic temperature range is 0-1, while OpenAI is 0-2
	// Convert OpenAI temperature range to Anthropic range
	if temperature > 1.0 {
		temperature /= 2.0 // Scale down from 0-2 to 0-1
	}
	// Ensure temperature is within Anthropic's acceptable range
	if temperature < 0 {
		temperature = 0
	} else if temperature > 1 {
		temperature = 1
	}

	// Convert OpenAI format messages to Anthropic format
	var anthropicMessages []models.AnthropicMessage
	var systemMessage string

	for _, msg := range messages {
		if msg.Role == "system" {
			systemMessage = msg.Content
		} else {
			anthropicMessages = append(anthropicMessages, models.AnthropicMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	requestBody := models.AnthropicRequest{
		Model:       model,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		Messages:    anthropicMessages,
	}

	if systemMessage != "" {
		requestBody.System = systemMessage
	}

	if maxTokens == 0 {
		requestBody.MaxTokens = 1000 // Anthropic requires max_tokens
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/v1/messages", strings.TrimSuffix(config.AppConfig.Anthropic.BaseURL, "/"))
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", config.AppConfig.Anthropic.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Anthropic API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var anthropicResp models.AnthropicResponse
	if err := json.Unmarshal(body, &anthropicResp); err != nil {
		return "", err
	}

	if len(anthropicResp.Content) == 0 {
		return "", fmt.Errorf("no response from Anthropic")
	}

	return anthropicResp.Content[0].Text, nil
}

func (s *UnifiedAIService) callOllama(messages []models.Message, temperature float64, maxTokens int, model string) (string, error) {
	if config.AppConfig.Ollama.BaseURL == "" {
		return "", fmt.Errorf("Ollama base URL not configured")
	}

	// Default to gemma3:12b if no model specified (user's preferred model)
	if model == "" {
		model = "gemma3:12b"
	}

	requestBody := models.OllamaRequest{
		Model:       model,
		Messages:    messages,
		Temperature: temperature,
		Stream:      false, // We want a complete response, not streaming
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	endpoint := fmt.Sprintf("%s/api/chat", strings.TrimSuffix(config.AppConfig.Ollama.BaseURL, "/"))
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Ollama API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var ollamaResp models.OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", err
	}

	if !ollamaResp.Done {
		return "", fmt.Errorf("incomplete response from Ollama")
	}

	return ollamaResp.Message.Content, nil
}
