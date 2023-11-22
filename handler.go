package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func getPokemonHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		query := r.URL.Query()
		name := query.Get("name")

		var err error
		var pokemon interface{}

		if name != "" {
			pokemon, err = getOnePokemonService(db, name)
		} else {
			pokemon, err = getAllPokemonService(db)
		}

		if err != nil {
			handleError(w, "Error fetching Pokemon", err)
			return
		}

		if err := json.NewEncoder(w).Encode(pokemon); err != nil {
			handleError(w, "Error encoding Pokemon to JSON", err)
		}
	}
}

func handleError(w http.ResponseWriter, msg string, err error) {
	status := http.StatusInternalServerError

	if err == sql.ErrNoRows {
		status = http.StatusNotFound
		msg = "Pokemon not found."
	} else if strings.Contains(err.Error(), "invalid character 'N' looking for beginning of value") {
		status = http.StatusNotFound
		msg = "Pokemon not found. Details: invalid input format"
	}

	log.Printf("%s: %v\n", msg, err)
	http.Error(w, msg, status)
}
