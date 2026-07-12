package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"website-summarizer/internal/models"
	"website-summarizer/internal/services"
)

func HandleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
		return
	}

	// 1. Scrape
	content, err := services.FetchWebpage(req.URL)
	if err != nil || content == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "Could not scrape website"})
		return
	}

	// 2. Summarize (Add error checking here!)
	prompt := fmt.Sprintf("Summarize this in 3 sentences: %s", content)
	answer, err := services.AskGemini(prompt) // Ensure this returns an error
	if err != nil || answer == "" {
		log.Printf("Gemini Error: %v", err) // Check your terminal for this!
		json.NewEncoder(w).Encode(map[string]string{"error": "Gemini failed to generate summary"})
		return
	}

	// 3. Respond
	json.NewEncoder(w).Encode(map[string]string{"gemini_summary": answer})
}
