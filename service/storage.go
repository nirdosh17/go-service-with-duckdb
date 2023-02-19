package main

import (
	"database/sql"
	"time"
)

type Storage struct {
	DB *sql.DB
}

type DailyAggregatedUser struct {
	Date  time.Time `json:"date"`
	Count int       `json:"count"`
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) AggregatedUsers() ([]DailyAggregatedUser, error) {
	users := []DailyAggregatedUser{}

	rows, err := s.DB.Query(`
		SELECT
			created as date,
			MAX(row_number) as count
		FROM (
			SELECT
				created,
				ROW_NUMBER() OVER (PARTITION BY created) as row_number
			FROM users
			WHERE created BETWEEN (now()::TIMESTAMP - INTERVAL '2 years') AND now()
		) GROUP BY created
		ORDER BY created ASC;`,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user DailyAggregatedUser
		if err := rows.Scan(&user.Date, &user.Count); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	rows.Close()

	return users, nil
}
