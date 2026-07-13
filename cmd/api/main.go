package main

import (
	"fmt"
	"log"
	"net/http"

	"website-summarizer/internal/db"
	"website-summarizer/internal/handlers"
)

func main() {
	db.InitDB()

	db.InitQdrant()

	db.EnsureCollection("website_knowledge")

	http.HandleFunc("/api/chat", handlers.HandleChat)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	fmt.Println("🚀 AI Server successfully running on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
