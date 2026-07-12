package models

type UserRequest struct {
	URL   string `json:"url"`
	Model string `json:"model"`
}

type APIResponse struct {
	Success       bool   `json:"success"`
	GeminiSummary string `json:"gemini_summary,omitempty"`
	Error         string `json:"error,omitempty"`
}
