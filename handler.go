package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func getPokemonHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %v\n", r)

		query := r.URL.Query()
		name := query.Get("name")

		w.Header().Set("Content-Type", "application/json")

		if name != "" {
			log.Printf("Request to fetch Pokemon with name: %s\n", name)
			pokemon, err := getOnePokemonService(name, db)
			if err != nil {
				log.Printf("Error fetching Pokemon with name %s: %v\n", name, err)
				status := http.StatusInternalServerError
				if err == sql.ErrNoRows {
					status = http.StatusNotFound
				}
				http.Error(w, http.StatusText(status), status)
				return
			}
			log.Printf("Successfully fetched Pokemon with name: %s\n", name)

			if err := json.NewEncoder(w).Encode(pokemon); err != nil {
				log.Printf("Error encoding Pokemon to JSON: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			log.Println("[Default] Request to get all existing Pokemon")
			pokemon, err := getAllPokemonService(db)
			if err != nil {
				log.Printf("Error fetching all Pokemon: %v\n", err)
				http.Error(w, "Failed to fetch pokemon", http.StatusInternalServerError)
				return
			}
			log.Println("Successfully fetched all Pokemon")

			if err := json.NewEncoder(w).Encode(pokemon); err != nil {
				log.Printf("Error encoding all Pokemon to JSON: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	}
}
