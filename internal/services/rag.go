package services

import (
	"log"
	"website-summarizer/internal/db"
)

func IngestURL(url string) error {
	content, err := FetchWebpage(url)
	if err != nil {
		return err
	}

	chunks := ChunkText(content, 500, 50)
	log.Printf("Ingesting %d chunks from %s", len(chunks), url)

	for i, chunk := range chunks {
		vector, err := GetEmbedding(chunk)
		if err != nil {
			// This catches the skipped chunks and keeps the loop running
			log.Printf("Skipping chunk %d: %v", i, err)
			continue
		}

		err = db.UpsertChunk("website_knowledge", uint64(i), vector, chunk)
		if err != nil {
			log.Printf("Upsert error on chunk %d: %v", i, err)
		}
	}
	return nil
}
