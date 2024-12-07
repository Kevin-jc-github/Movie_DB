package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

func main() {
	// Get current directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Current working directory: %s\n", wd)

	schemaPath := wd + "/schema.sql"

	// Read schema.sql file
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to open schema file: %v", err)
	}

	// Open db（Will automatically create movies.db if it is not exist）
	db, err := sql.Open("sqlite", "./movies.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// initializing the schema db
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatalf("Failed to execute schema: %v", err)
	}

	fmt.Println("Database schema initialized successfully!")
}
