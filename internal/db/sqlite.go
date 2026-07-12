package db

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/glebarez/go-sqlite"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./summaries.db")
	if err != nil {
		panic(err)
	}

	// Create table if it doesn't exist
	sqlStmt := `CREATE TABLE IF NOT EXISTS summaries (
		url TEXT PRIMARY KEY,
		summary TEXT,
		model TEXT
	);`
	DB.Exec(sqlStmt)
}

func GetSummary(url string) (string, error) {

	url = strings.TrimSpace(url)
	log.Printf("DEBUG: Searching DB for URL: '%s'", url)

	var summary string
	err := DB.QueryRow("SELECT summary FROM summaries WHERE url = ?", url).Scan(&summary)

	if err != nil {
		log.Printf("DEBUG: GetSummary error: %v", err)
	} else {
		log.Printf("DEBUG: Found summary for '%s'", url)
	}
	return summary, err
}

func SaveSummary(url, summary, model string) error {
	// 1. Prepare the statement
	stmt, err := DB.Prepare("INSERT OR REPLACE INTO summaries(url, summary, model) VALUES(?,?,?)")
	if err != nil {
		return err // Return the error so we know if the SQL is bad
	}
	// 2. CRITICAL: Close the statement when done to prevent resource leaks
	defer stmt.Close()

	// 3. Execute
	_, err = stmt.Exec(url, summary, model)
	return err // Return the result of the execution
}
