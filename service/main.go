package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// "github.com/gin-contrib/pprof"

	_ "github.com/marcboeker/go-duckdb"
)

func main() {
	router := http.NewServeMux()
	// enable this to profile the service
	// pprof.Register(router)

	storage := NewStorage()

	router.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		uid, err := strconv.Atoi(id)
		if err != nil {
			badRequest(w)
			return
		}

		user, err := storage.GetUserByID(uid)
		if err == nil {
			b, merr := json.Marshal(user)
			if merr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			return
		}

		if err == ErrUserNotFound {
			notFound(w)
			return
		}

		internalError(w)
	})

	log.Println("starting server at localhost:8000")
	if err := http.ListenAndServe(":8000", router); err != nil {
		log.Fatalln("failed to start server", err)
	}
}

func badRequest(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"message": "invalid user id"}`))
}

func notFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{"message": "user not found"}`))
}

func internalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"message": "internal server error"}`))
}
