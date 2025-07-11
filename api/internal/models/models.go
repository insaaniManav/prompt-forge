package models

import "time"

// Request/Response structures for OpenAI API
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model               string    `json:"model"`
	Messages            []Message `json:"messages"`
	Temperature         float64   `json:"temperature"`
	MaxTokens           int       `json:"max_tokens,omitempty"`
	MaxCompletionTokens int       `json:"max_completion_tokens,omitempty"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// API Request structures
type CritiqueRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model,omitempty"`
}

type ExecuteRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model,omitempty"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
}

type PromptEngineerRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model,omitempty"`
	Temperature float64   `json:"temperature"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Data    string `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// Dual Analysis Response structures
type DualAnalysisResponse struct {
	Success bool              `json:"success"`
	Data    *DualAnalysisData `json:"data,omitempty"`
	Error   string            `json:"error,omitempty"`
}

type DualAnalysisData struct {
	QuickReport    string `json:"quick_report"`
	DetailedReport string `json:"detailed_report"`
}

// History structures
type HistoryItem struct {
	ID          int64     `json:"id" db:"id"`
	Timestamp   time.Time `json:"timestamp" db:"timestamp"`
	Prompt      string    `json:"prompt" db:"prompt"`
	Model       string    `json:"model" db:"model"`
	Temperature float64   `json:"temperature" db:"temperature"`
	MaxTokens   int       `json:"max_tokens" db:"max_tokens"`
	Success     bool      `json:"success" db:"success"`
	Response    string    `json:"response" db:"response"`
	ErrorMsg    string    `json:"error_msg,omitempty" db:"error_msg"`
}

type SaveHistoryRequest struct {
	Prompt      string  `json:"prompt"`
	Model       string  `json:"model"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	Success     bool    `json:"success"`
	Response    string  `json:"response"`
	ErrorMsg    string  `json:"error_msg,omitempty"`
}

type HistoryResponse struct {
	Success bool          `json:"success"`
	Data    []HistoryItem `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// Conversation structures
type ConversationMessage struct {
	ID             int64     `json:"id" db:"id"`
	ConversationID string    `json:"conversation_id" db:"conversation_id"`
	Role           string    `json:"role" db:"role"`
	Content        string    `json:"content" db:"content"`
	Timestamp      time.Time `json:"timestamp" db:"timestamp"`
}

type Conversation struct {
	ID        string                `json:"id" db:"id"`
	Title     string                `json:"title" db:"title"`
	CreatedAt time.Time             `json:"created_at" db:"created_at"`
	UpdatedAt time.Time             `json:"updated_at" db:"updated_at"`
	Messages  []ConversationMessage `json:"messages,omitempty"`
}

type SaveConversationRequest struct {
	ConversationID string                `json:"conversation_id"`
	Title          string                `json:"title,omitempty"`
	Messages       []ConversationMessage `json:"messages"`
}

type ConversationResponse struct {
	Success bool           `json:"success"`
	Data    []Conversation `json:"data,omitempty"`
	Error   string         `json:"error,omitempty"`
}

type ConversationDetailResponse struct {
	Success bool          `json:"success"`
	Data    *Conversation `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

// Prompt Library structures
type SavedPrompt struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Content     string    `json:"content" db:"content"`
	Description string    `json:"description" db:"description"`
	Category    string    `json:"category" db:"category"`
	Tags        string    `json:"tags" db:"tags"` // JSON string of tag array
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	UsageCount  int       `json:"usage_count" db:"usage_count"`
}

type SavePromptRequest struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type UpdatePromptRequest struct {
	ID          int64    `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type PromptLibraryResponse struct {
	Success bool          `json:"success"`
	Data    []SavedPrompt `json:"data,omitempty"`
	Error   string        `json:"error,omitempty"`
}

type PromptResponse struct {
	Success bool         `json:"success"`
	Data    *SavedPrompt `json:"data,omitempty"`
	Error   string       `json:"error,omitempty"`
}

// Eval Generator structures
type EvalGenerateRequest struct {
	Prompt     string   `json:"prompt"`
	EvalTypes  []string `json:"eval_types"`
	SampleSize int      `json:"sample_size"`
	Model      string   `json:"model,omitempty"`
	Difficulty string   `json:"difficulty,omitempty"`
}

type TestCase struct {
	Input      string `json:"input"`
	Category   string `json:"category"`
	Difficulty string `json:"difficulty"`
	Expected   string `json:"expected,omitempty"`
}

type EvalCriterion struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Weight      int    `json:"weight"`
}

type EvalData struct {
	TestCases  []TestCase      `json:"test_cases"`
	Criteria   []EvalCriterion `json:"criteria"`
	BasePrompt string          `json:"base_prompt"`
	Metadata   EvalMetadata    `json:"metadata"`
}

type EvalMetadata struct {
	GeneratedAt time.Time `json:"generated_at"`
	Model       string    `json:"model"`
	SampleSize  int       `json:"sample_size"`
	EvalTypes   []string  `json:"eval_types"`
	Difficulty  string    `json:"difficulty"`
}

type EvalResponse struct {
	Success bool      `json:"success"`
	Data    *EvalData `json:"data,omitempty"`
	Error   string    `json:"error,omitempty"`
}

// Multi-model execution structures
type MultiModelExecuteRequest struct {
	Prompt      string   `json:"prompt"`
	Models      []string `json:"models"`
	Temperature float64  `json:"temperature"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
}

type ModelExecutionResult struct {
	Model         string      `json:"model"`
	Response      string      `json:"response"`
	Success       bool        `json:"success"`
	Error         string      `json:"error,omitempty"`
	ExecutionTime int64       `json:"execution_time_ms,omitempty"`
	TokenUsage    *TokenUsage `json:"token_usage,omitempty"`
}

type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}

type MultiModelExecuteResponse struct {
	Success bool                   `json:"success"`
	Data    []ModelExecutionResult `json:"data,omitempty"`
	Error   string                 `json:"error,omitempty"`
}

// Anthropic API structures
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AnthropicRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	Temperature float64            `json:"temperature,omitempty"`
	Messages    []AnthropicMessage `json:"messages"`
	System      string             `json:"system,omitempty"`
}

type AnthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AnthropicResponse struct {
	Content []AnthropicContent `json:"content"`
}

// Ollama API structures
type OllamaRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	Stream      bool      `json:"stream"`
}

type OllamaResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}
