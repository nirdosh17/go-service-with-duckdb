package main

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	_ "github.com/marcboeker/go-duckdb"
)

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

	st := time.Now()
	i := 1
	max := 1000000 // 1 million
	log.Println("inserting 1 million randomly generated rows in users table...")
	var date string
	for i <= max {
		date = randate(2021, 2022).Format("2006-01-02")
		_, err := db.Exec(`INSERT INTO users (id, joined_date, name, email) VALUES (?, ?, ?, ?)`, i, date, gofakeit.Name(), gofakeit.Email())
		if err != nil {
			log.Println("failed to insert", err)
			break
		}
		if i%100000 == 0 {
			log.Printf("%0.0f%% complete\n", float32((i*100)/max))
		}
		i++
	}
	log.Println("done! time taken:", time.Since(st))
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
			joined_date DATE,
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
