package services

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/genai"
)

// AskGemini sends the text prompt to Google's gemini-2.5-flash model
func AskGemini(userPrompt string) (string, error) {
	ctx := context.Background()
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable is not set")
	}

	// Initialize the client using the official GenAI SDK layout
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini client: %v", err)
	}

	// Generate content using the fast, lightweight flash model

	result, err := client.Models.GenerateContent(ctx, "gemini-3.5-flash", genai.Text(userPrompt), nil)
	if err != nil {
		return "", fmt.Errorf("gemini generation failed: %v", err)
	}

	return result.Text(), nil
}
