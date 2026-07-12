package services

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func AskOllama(prompt string) (string, error) {
	// Ollama's local API expects this format
	data := map[string]interface{}{
		"model":  "llama3",
		"prompt": prompt,
		"stream": false,
	}
	payload, _ := json.Marshal(data)

	// Sending the request to your local Ollama server
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	// Returning the "response" field from Ollama's JSON
	return result["response"].(string), nil
}
