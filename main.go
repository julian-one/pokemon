package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"time"

	"pokemon/internal/api"
	"pokemon/internal/database"

	"github.com/jmoiron/sqlx"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	db, err := database.Initialize()
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", FetchSavedPokemon(db))
	mux.HandleFunc("GET /pokemon", FetchPokemon(db))

	handler := loggingMiddleware(mux)

	slog.Info("Server starting", "port", 8080)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

func FetchSavedPokemon(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p, err := database.GetAll(db)
		if err != nil {
			slog.Error("Failed to fetch all pokemon", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Error fetching all pokemon"}); err != nil {
				slog.Error("Failed to encode error response", "error", err)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(p); err != nil {
			slog.Error("Failed to encode response", "error", err)
		}
	}
}

func FetchPokemon(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Query().Get("name")
		if name == "" {
			slog.Warn("Request missing pokemon name")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Name is required"}); err != nil {
				slog.Error("Failed to encode error response", "error", err)
			}
			return
		}

		p, err := api.Search(name)
		if err != nil {
			slog.Error("Failed to fetch pokemon from API", "name", name, "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Error fetching pokemon from the api"}); err != nil {
				slog.Error("Failed to encode error response", "error", err)
			}
			return
		}

		exists, err := database.Exists(db, p.Id)
		if err != nil {
			slog.Error("Failed to check if pokemon exists", "name", name, "id", p.Id, "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if err := json.NewEncoder(w).Encode(map[string]string{"error": "Error checking if pokemon exists"}); err != nil {
				slog.Error("Failed to encode error response", "error", err)
			}
			return
		}

		if !exists {
			err = database.Save(db, p)
			if err != nil {
				slog.Error("Failed to save pokemon", "name", name, "id", p.Id, "error", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Error saving pokemon"}); err != nil {
					slog.Error("Failed to encode error response", "error", err)
				}
				return
			}
			slog.Info("Saved new pokemon", "name", name, "id", p.Id)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(p); err != nil {
			slog.Error("Failed to encode response", "error", err)
		}
	}
}
