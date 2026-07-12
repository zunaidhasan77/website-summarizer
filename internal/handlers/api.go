package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"website-summarizer/internal/db"
	"website-summarizer/internal/models"
	"website-summarizer/internal/services"
)

// HandleChat is now just the "controller" that manages the flow.
func HandleChat(w http.ResponseWriter, r *http.Request) {
	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 1. Cache Check
	if summary, err := db.GetSummary(req.URL); err == nil {
		log.Println("DEBUG: Cache hit.")
		respondJSON(w, map[string]string{"gemini_summary": summary, "status": "Cached"})
		return
	}

	// 2. Scrape Content
	content, err := services.FetchWebpage(req.URL)
	if err != nil || content == "" {
		respondError(w, "Could not scrape website", http.StatusInternalServerError)
		return
	}

	// 3. AI Processing
	answer, status := routeAI(req.Model, content)

	// 4. Persistence
	if answer != "" {
		if err := db.SaveSummary(req.URL, answer, req.Model); err != nil {
			log.Printf("ERROR: DB save failed: %v", err)
		}
	} else {
		respondError(w, "AI service failed to generate summary", http.StatusInternalServerError)
		return
	}

	// 5. Final Response
	respondJSON(w, map[string]string{"gemini_summary": answer, "status": status})
}

// routeAI handles the logic of which model to pick and the failover.
func routeAI(model, content string) (string, string) {
	prompt := fmt.Sprintf("Summarize this in 3 sentences: %s", content)

	// Option A: Explicitly request Ollama
	if model == "ollama" {
		log.Println("DEBUG: Routing to Ollama.")
		ans, err := services.AskOllama(prompt)
		if err != nil {
			return "", ""
		}
		return ans, "Success"
	}

	// Option B: Try Gemini with Failover
	log.Println("DEBUG: Routing to Gemini.")
	ans, err := services.AskGemini(prompt)
	if err == nil {
		return ans, "Success"
	}

	// Failover happens here
	log.Printf("DEBUG: Gemini failed (%v). Switching to Ollama.", err)
	ans, err = services.AskOllama(prompt)
	if err != nil {
		return "", ""
	}
	return ans, "Gemini failed, switched to local Ollama"
}

// --- Helper Functions to keep code clean ---

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
