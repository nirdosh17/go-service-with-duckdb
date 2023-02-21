package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"path/filepath"

	// "github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	router := gin.Default()
	// enable this to profile the service
	// pprof.Register(router)

	db := setupDuckDB()
	storage := NewStorage(db)

	router.GET("/stats", func(c *gin.Context) {
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

	router.Run(":8000")
}

func setupDuckDB() *sql.DB {
	absPath, _ := filepath.Abs("../prepare-test-data/test.duckdb")
	db, err := sql.Open("duckdb", absPath+"?access_mode=read_only")
	if err != nil {
		panic(err)
	}
	return db
}
