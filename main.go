package main

import (
	"log"
	"net/http"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("Error initializing database: ", err)
	}

	http.HandleFunc("/pokemon/", getPokemonHandler(db))
	http.ListenAndServe(":8080", nil)
}
