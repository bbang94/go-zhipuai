package zhipuai

// common.go defines common types used throughout the zhipuai API.

// Usage Represents the total token usage per request to zhipuai.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
