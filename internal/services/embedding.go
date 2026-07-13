package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Modern Ollama embed request API
type OllamaEmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// Modern Ollama embed response API
type OllamaEmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

func GetEmbedding(text string) ([]float32, error) {
	// 1. Prevent sending empty text (which causes the 400 Bad Request)
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, fmt.Errorf("text chunk is empty, skipping")
	}

	reqData := OllamaEmbedRequest{
		Model: "nomic-embed-text",
		Input: text,
	}
	body, err := json.Marshal(reqData)
	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	// 2. Use the modern /api/embed endpoint
	resp, err := http.Post("http://localhost:11434/api/embed", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("could not connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	// 3. Catch the 400 Bad Request and print the actual error from Ollama
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result OllamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("could not decode Ollama response: %w", err)
	}

	// 4. Validate the new nested array structure
	if len(result.Embeddings) == 0 || len(result.Embeddings[0]) == 0 {
		return nil, fmt.Errorf("received empty embedding array from Ollama")
	}

	log.Printf("Successfully generated vector of size: %d", len(result.Embeddings[0]))
	return result.Embeddings[0], nil
}
