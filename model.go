package vapi

// ModelConfig contains LLM settings
type ModelConfig struct {
	Provider                  string         `json:"provider,omitempty"`
	Model                     string         `json:"model,omitempty"`
	EmotionRecognitionEnabled bool           `json:"emotionRecognitionEnabled,omitempty"`
	KnowledgeBase             *KnowledgeBase `json:"knowledgeBase,omitempty"`
	KnowledgeBaseID           *string        `json:"knowledgeBaseId,omitempty"`
	MaxTokens                 float64        `json:"maxTokens,omitempty"`
	Messages                  []ModelMessage `json:"messages,omitempty"`
	NumFastTurns              *float64       `json:"numFastTurns,omitempty"`
	Temperature               float64        `json:"temperature,omitempty"`
	ToolIDs                   []string       `json:"toolIds,omitempty"`
	Tools                     []Tool         `json:"tools,omitempty"`
}

// ModelMessage represents a single message in the model conversation
type ModelMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
