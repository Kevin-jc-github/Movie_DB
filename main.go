package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

// Createt Database Tables
func createTables(db *sql.DB) {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS movies;
		DROP TABLE IF EXISTS movies_genres;
	`)
	if err != nil {
		log.Fatalf("Failed to drop tables: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE movies (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			year INTEGER,
			rank REAL
		);
		CREATE TABLE movies_genres (
			movie_id INTEGER,
			genre TEXT NOT NULL,
			FOREIGN KEY (movie_id) REFERENCES movies (id)
		);
	`)
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
}

func clearTables(db *sql.DB) {
	_, err := db.Exec("DELETE FROM movies;")
	if err != nil {
		log.Fatalf("Failed to clear movies table: %v", err)
	}
	_, err = db.Exec("DELETE FROM movies_genres;")
	if err != nil {
		log.Fatalf("Failed to clear movies_genres table: %v", err)
	}
	log.Println("Tables cleared successfully.")
}

func populateMoviesTable(db *sql.DB, moviesFile string) {
	moviesCSV, err := os.Open(moviesFile)
	if err != nil {
		log.Fatalf("Failed to open movies file: %v", err)
	}
	defer moviesCSV.Close()

	reader := csv.NewReader(moviesCSV)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	lineNum := 0
	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading line %d: %v. Attempting to fix...", lineNum, err)
			record = fixBrokenCSVLine(record)
			if record == nil {
				log.Printf("Skipping line %d due to irreparable error.", lineNum)
				continue
			}
		}

		id, _ := strconv.Atoi(record[0])
		name := record[1]
		year := 0
		if record[2] != "NULL" {
			year, _ = strconv.Atoi(record[2])
		}
		rank := sql.NullFloat64{}
		if record[3] != "NULL" {
			val, _ := strconv.ParseFloat(record[3], 64)
			rank = sql.NullFloat64{Float64: val, Valid: true}
		}

		_, err = db.Exec(`
			INSERT OR IGNORE INTO movies (id, name, year, rank)
			VALUES (?, ?, ?, ?)`,
			id, name, year, rank)
		if err != nil {
			log.Printf("Failed to insert record with id=%d: %v", id, err)
		}

	}
	log.Println("Movies table populated successfully.")
}

func fixBrokenCSVLine(record []string) []string {
	for i, field := range record {
		if strings.Count(field, "\"")%2 != 0 {
			record[i] = strings.ReplaceAll(field, "\"", "")
		}
	}

	if len(record) < 4 {
		return nil
	}

	return record
}

func populateMoviesGenresTable(db *sql.DB, genresFile string) {
	genresCSV, err := os.Open(genresFile)
	if err != nil {
		log.Fatalf("Failed to open genres file: %v", err)
	}
	defer genresCSV.Close()

	reader := csv.NewReader(genresCSV)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	lineNum := 0
	for {
		lineNum++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading line %d: %v. Attempting to fix...", lineNum, err)
			record = fixBrokenCSVLine(record)
			if record == nil {
				log.Printf("Skipping line %d due to irreparable error.", lineNum)
				continue
			}
		}

		movieID, _ := strconv.Atoi(record[0])
		genre := record[1]

		_, err = db.Exec(`
			INSERT INTO movies_genres (movie_id, genre)
			VALUES (?, ?)`,
			movieID, genre)
		if err != nil {
			log.Printf("Failed to insert line %d into movies_genres table: %v", lineNum, err)
		}
	}
	log.Println("Movies genres table populated successfully.")
}

func queryHighestRatedGenres(db *sql.DB) {
	rows, err := db.Query(`
		SELECT genre, AVG(rank) AS avg_rank
		FROM movies
		JOIN movies_genres ON movies.id = movies_genres.movie_id
		WHERE rank IS NOT NULL
		GROUP BY genre
		ORDER BY avg_rank DESC;
	`)
	if err != nil {
		log.Fatalf("Failed to execute query: %v", err)
	}
	defer rows.Close()

	fmt.Println("Genre\t\tAverage Rank")
	fmt.Println("-----------------------------")
	for rows.Next() {
		var genre string
		var avgRank float64
		err := rows.Scan(&genre, &avgRank)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		fmt.Printf("%s\t\t%.2f\n", genre, avgRank)
	}
}

func main() {
	db, err := sql.Open("sqlite", "./movies.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTables(db)

	clearTables(db)

	populateMoviesTable(db, "IMDB-movies.csv")
	populateMoviesGenresTable(db, "IMDB-movies_genres.csv")

	fmt.Println("Database populated successfully.")

	// execute query result:  What types of movies (by genre) receive the highest ratings/ranks?
	fmt.Println("What types of movies (by genre) receive the highest ratings/ranks?:")
	queryHighestRatedGenres(db)
}
