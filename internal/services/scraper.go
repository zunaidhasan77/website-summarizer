package services

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func FetchWebpage(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("website returned status code: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read content: %v", err)
	}

	return cleanHTML(string(bodyBytes)), nil
}

// Separate function for the cleanup logic
func cleanHTML(html string) string {
	// 1. Manually replace common junk tags with spaces
	// This is safer and avoids the panic error
	text := strings.ReplaceAll(html, "<script", " <script")
	text = strings.ReplaceAll(text, "</script>", "> ")
	text = strings.ReplaceAll(text, "<style", " <style")
	text = strings.ReplaceAll(text, "</style>", "> ")

	// 2. Now remove all tags (everything between < and >)
	reTags := regexp.MustCompile("<[^>]*>")
	text = reTags.ReplaceAllString(text, " ")

	// 3. Clean up whitespace
	text = strings.Join(strings.Fields(text), " ")

	return text
}
