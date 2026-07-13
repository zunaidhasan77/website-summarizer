package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"website-summarizer/internal/db"
	"website-summarizer/internal/models"
	"website-summarizer/internal/services"
)

func HandleChat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("API Error: Invalid request body: %v", err)
		json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Invalid JSON request"})
		return
	}

	// Default to summarize if mode is missing
	if req.Mode == "" {
		req.Mode = "summarize"
	}

	log.Printf("Processing request - Mode: %s, URL: %s, Model: %s", req.Mode, req.URL, req.Model)

	var finalPrompt string

	// 1. Build the Prompt
	if req.Mode == "summarize" {
		content, err := services.FetchWebpage(req.URL)
		if err != nil {
			json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Failed to read website"})
			return
		}
		if len(content) > 5000 {
			content = content[:5000] + "... [truncated]"
		}
		finalPrompt = fmt.Sprintf("You are an expert summarizer. Please provide a clear, 3-paragraph summary of the following website content:\n\n%s", content)
	} else {

		// --- JUST-IN-TIME INGESTION ---
		// This makes the system "flow" automatically from the prompt
		log.Printf("Just-in-Time Ingestion: %s", req.URL)
		err := services.IngestURL(req.URL)
		if err != nil {
			log.Printf("Ingestion error: %v", err)
			json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Failed to ingest URL automatically"})
			return
		}

		// --- RAG QUERY ---
		if strings.TrimSpace(req.Message) == "" {
			json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "No question provided"})
			return
		}
		queryVector, err := services.GetEmbedding(req.Message)
		if err != nil {
			json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Failed to embed question"})
			return
		}
		chunks, err := db.SearchChunks("website_knowledge", queryVector, 3)
		if err != nil {
			json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "Database search failed"})
			return
		}
		contextData := strings.Join(chunks, "\n\n---\n\n")
		finalPrompt = fmt.Sprintf("You are a helpful assistant. Use the following context to answer the user's question. If the answer is not in the context, say 'I cannot answer this based on the provided website'.\n\nCONTEXT:\n%s\n\nQUESTION:\n%s", contextData, req.Message)
	}

	// 2. Call your specific Service functions
	var aiResponse string
	var err error

	if req.Model == "gemini" {
		aiResponse, err = services.AskGemini(finalPrompt)

		if err != nil {
			log.Printf("CRITICAL: Gemini call failed (%v). Triggering fallback to Ollama...", err)

			// Fallback: Try Ollama
			aiResponse, err = services.AskOllama(finalPrompt)

			if err != nil {
				log.Printf("FALLBACK FAILED: Ollama also returned an error: %v", err)
			} else {
				log.Println("SUCCESS: Fallback to Ollama successful.")
			}
		}
	} else {
		aiResponse, err = services.AskOllama(finalPrompt)
	}

	if err != nil {
		log.Printf("LLM Generation Error: %v", err)
		json.NewEncoder(w).Encode(models.APIResponse{Success: false, Error: "AI generation failed: " + err.Error()})
		return
	}

	// 3. Return the actual AI response
	if req.Mode == "summarize" {
		json.NewEncoder(w).Encode(models.APIResponse{
			Success:       true,
			GeminiSummary: aiResponse,
		})
	} else {
		json.NewEncoder(w).Encode(models.APIResponse{
			Success:   true,
			ChatReply: aiResponse,
		})
	}
}

func HandleIngest(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	if err := services.IngestURL(url); err != nil {
		http.Error(w, "Ingestion failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Ingestion successful! Check your database."))
}
