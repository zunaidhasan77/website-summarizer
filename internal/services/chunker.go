package services

import "strings"

// ChunkText splits text into segments to maintain semantic context
func ChunkText(text string, chunkSize int, overlap int) []string {
	var chunks []string
	runes := []rune(text)
	n := len(runes)

	if n <= chunkSize {
		return []string{text}
	}

	for i := 0; i < n; i += (chunkSize - overlap) {
		end := i + chunkSize
		if end > n {
			end = n
		}
		chunk := string(runes[i:end])
		chunks = append(chunks, strings.TrimSpace(chunk))
		if end == n {
			break
		}
	}
	return chunks
}
