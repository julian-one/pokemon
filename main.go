package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"pokemon/internal/api"
	"pokemon/internal/database"

	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := database.Initialize()
	if err != nil {
		log.Fatalf("Error initializing database %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", FetchSavedPokemon(db))
	mux.HandleFunc("/pokemon", FetchPokemon(db))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Server started on port 8080")
}

func FetchSavedPokemon(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := database.GetAll(db)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Error fetching all pokemon: %s", err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(p)
	}
}

func FetchPokemon(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = fmt.Fprintf(w, "Name is required")
		}

		p, err := api.Search(name)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Error fetching pokemon from the api: %s", err)
		}

		exists, err := database.Exists(db, p.Id)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprintf(w, "Error checking if pokemon exists: %s", err)
		}

		if !exists {
			err = database.Save(db, p)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = fmt.Fprintf(w, "Error saving pokemon: %s", err)
			}
		}

		_ = json.NewEncoder(w).Encode(p)
	}
}
