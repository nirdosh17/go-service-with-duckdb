package main

import (
	"database/sql"
	"errors"
	"path/filepath"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type Storage struct {
	DB *sql.DB
}

func NewStorage() *Storage {
	ddb := initDuckDB()
	return &Storage{DB: ddb}
}

type User struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	JoinedDate time.Time `json:"joined_date"`
}

func (s *Storage) GetUserByID(id int) (User, error) {
	row := s.DB.QueryRow(`
		SELECT
			id,
			name,
			email,
			joined_date
		FROM users
		WHERE id = ?;`, id)
	user := User{}
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.JoinedDate); err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}

func initDuckDB() *sql.DB {
	absPath, _ := filepath.Abs("../prepare-test-data/test.duckdb")
	db, err := sql.Open("duckdb", absPath+"?access_mode=read_only")
	if err != nil {
		panic(err)
	}
	return db
}
