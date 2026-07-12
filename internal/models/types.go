package models

type UserRequest struct {
	URL string `json:"url"`
}

type APIResponse struct {
	Success       bool   `json:"success"`
	GeminiSummary string `json:"gemini_summary,omitempty"`
	Error         string `json:"error,omitempty"`
}
