package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/marcboeker/go-duckdb"
)

var DefaultSeedDataSize = 100_000

func randate(startYear, endYear int) time.Time {
	min := time.Date(startYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Date(endYear, 1, 0, 0, 0, 0, 0, time.UTC).Unix()
	delta := max - min
	sec := rand.Int63n(delta) + min
	return time.Unix(sec, 0).UTC()
}

func main() {
	db, err := duckdb()
	if err != nil {
		panic(err)
	}

	seedDataSize, err := strconv.Atoi(os.Getenv("SEED_COUNT"))
	if err != nil {
		log.Println("failed to parse env SEED_COUNT, defaulting to", DefaultSeedDataSize)
		seedDataSize = DefaultSeedDataSize
	}

	log.Printf("inserting %v rows...", seedDataSize)
	start := time.Now()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Failed to start transaction:", err)
	}
	stmt, err := tx.Prepare("INSERT INTO users (id, joined_date, name, email) VALUES (?, ?, ?, ?)")
	if err != nil {
		log.Fatal("Failed to prepare statement:", err)
	}
	defer stmt.Close()

	var date string
	for i := 1; i <= seedDataSize; i++ {
		date = randate(2021, 2022).Format(time.RFC3339)
		_, err := stmt.Exec(i, date, gofakeit.Name(), gofakeit.Email())
		if err != nil {
			log.Fatal("Failed to insert row:", err)
		}

		// track %
		fmt.Printf("%0.0f%% complete\r", float32((i*100)/seedDataSize))
	}

	if err := tx.Commit(); err != nil {
		log.Fatal("Failed to commit transaction:", err)
	}

	log.Printf("%v rows populated in %v\n", seedDataSize, time.Since(start))
}

func duckdb() (*sql.DB, error) {
	db, err := sql.Open("duckdb", "test.duckdb")
	if err != nil {
		log.Println("failed to open test.duckdb", err)
		return nil, err
	}

	log.Println("opened test.duckdb duckdb database file")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER,
			joined_date TIMESTAMP,
			name VARCHAR,
			email VARCHAR,
			PRIMARY KEY(id)
		);
	`)

	if err != nil {
		log.Println("failed to create table", err)
		return nil, err
	}

	log.Println("table 'users' created successfully")

	return db, nil
}
