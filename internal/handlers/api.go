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

	prompt := fmt.Sprintf("Summarize this in 3 sentences: %s", content)
	var answer string

	var status string
	// --- NEW LOGIC START ---
	// If the user specifically picks 'ollama', use it.
	// Otherwise, try Gemini, and fallback to Ollama if it fails.
	if req.Model == "ollama" {
		log.Println("DEBUG: System is routing to Ollama (Local).")
		answer, err = services.AskOllama(prompt)
	} else {
		log.Println("DEBUG: System is routing to Gemini (Cloud).")
		answer, err = services.AskGemini(prompt)

		// The Fallback: If Gemini hits a rate limit or fails, switch to local model
		if err != nil {
			log.Printf("DEBUG: Gemini failed (%v), triggering Failover to Ollama.", err)
			status = "Gemini failed, switching to local Ollama..."
			answer, err = services.AskOllama(prompt)
		}
	}
	// --- NEW LOGIC END ---

	// Final error check if both models fail
	if err != nil || answer == "" {
		json.NewEncoder(w).Encode(map[string]string{"error": "AI service failure"})
		return
	}

	// 3. Respond
	json.NewEncoder(w).Encode(map[string]string{"gemini_summary": answer, "status": status})
}
