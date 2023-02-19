package main

import (
	"database/sql"
)

type Storage struct {
	DB *sql.DB
}

type DailyAggregatedUser struct {
	Date        string `json:"date"`
	UsersJoined int    `json:"users_joined"`
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) AggregatedUsers() ([]DailyAggregatedUser, error) {
	users := []DailyAggregatedUser{}

	rows, err := s.DB.Query(`
		SELECT
			CAST(joined_date AS VARCHAR) as date,
			MAX(row_number) as users_joined
		FROM (
			SELECT
				joined_date,
				ROW_NUMBER() OVER (PARTITION BY joined_date) as row_number
			FROM users
			WHERE joined_date BETWEEN (now()::TIMESTAMP - INTERVAL '2 years') AND now()
		) GROUP BY joined_date
		ORDER BY joined_date ASC;`,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user DailyAggregatedUser
		if err := rows.Scan(&user.Date, &user.UsersJoined); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	rows.Close()

	return users, nil
}
