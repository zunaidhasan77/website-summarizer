package handlers

import (
	"encoding/json"
	"net/http"

	"website-summarizer/internal/models"
	"website-summarizer/internal/services"
)

// HandleChat processes incoming POST requests for the AI agent
func HandleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Enforce POST method rules
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Method not allowed"})
		return
	}

	// Decode the incoming JSON payload
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Invalid JSON payload"})
		return
	}

	// Call the Gemini service layer
	answer, err := services.AskGemini(req.Prompt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: err.Error()})
		return
	}

	// Return successful response
	json.NewEncoder(w).Encode(models.APIResponse{
		Success: true,
		Answer:  answer,
	})
}
