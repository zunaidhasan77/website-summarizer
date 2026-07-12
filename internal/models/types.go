package models

// UserRequest handles the incoming payload from the user or UI
type UserRequest struct {
	Prompt string `json:"prompt"`
}

// APIResponse structures the JSON sent back to the client
type APIResponse struct {
	Success bool   `json:"success"`
	Answer  string `json:"answer"`
	Error   string `json:"error,omitempty"`
}