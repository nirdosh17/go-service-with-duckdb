package main

import (
	"database/sql"
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type Storage struct {
	DB *sql.DB
}

func NewStorage(ddbFilePath string) *Storage {
	ddb := initDuckDB(ddbFilePath)
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

func initDuckDB(ddbFilePath string) *sql.DB {
	db, err := sql.Open("duckdb", ddbFilePath+"?access_mode=read_only")
	if err != nil {
		panic(err)
	}
	if db.Ping() != nil {
		panic("failed to ping db!")
	}

	return db
}
