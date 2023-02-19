package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	r := gin.Default()
	db := setupDuckDB()
	storage := NewStorage(db)

	r.GET("/users_joined_daily", func(c *gin.Context) {
		users, err := storage.AggregatedUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		bytes, err := json.Marshal(users)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.Data(http.StatusOK, "application/json", bytes)
	})

	r.Run(":8000")
}

func setupDuckDB() *sql.DB {
	absPath, _ := filepath.Abs("../prepare-test-data/test.db")
	db, err := sql.Open("duckdb", absPath+"?access_mode=read_only")
	if err != nil {
		panic(err)
	}
	return db
}
