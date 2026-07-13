package models

type UserRequest struct {
	Mode    string `json:"mode"`    // "summarize" or "chat"
	Message string `json:"message"` // The user's question (if in chat mode)
	URL     string `json:"url"`
	Model   string `json:"model"`
}

type APIResponse struct {
	Success       bool   `json:"success"`
	GeminiSummary string `json:"gemini_summary,omitempty"`
	ChatReply     string `json:"chat_reply,omitempty"`
	Error         string `json:"error,omitempty"`
}
