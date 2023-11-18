package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createTableSQL = `
		CREATE TABLE IF NOT EXISTS pokemon (
			"id" INTEGER PRIMARY KEY,
			"name" TEXT,
			"height" INTEGER,
			"weight" INTEGER
		);
	`
	selectAllPokemonSQL = `SELECT id, name, height, weight FROM pokemon`
	selectOnePokemonSQL = `SELECT id, name, height, weight FROM pokemon WHERE name = ? LIMIT 1`
	insertPokemonSQL    = `INSERT INTO pokemon (id, name, height, weight) VALUES (?, ?, ?, ?)`
)

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "pokemon.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getAllPokemonFromDB(db *sql.DB) ([]Pokemon, error) {
	var pokemon []Pokemon
	rows, err := db.Query(selectAllPokemonSQL)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p Pokemon
		if err := rows.Scan(&p.ID, &p.Name, &p.Height, &p.Weight); err != nil {
			return nil, err
		}
		pokemon = append(pokemon, p)
	}

	return pokemon, nil
}

func getOnePokemonFromDB(db *sql.DB, name string) (Pokemon, error) {
	var pokemon Pokemon
	row := db.QueryRow(selectOnePokemonSQL, name)
	err := row.Scan(&pokemon.ID, &pokemon.Name, &pokemon.Height, &pokemon.Weight)

	if err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

func savePokemonToDB(db *sql.DB, pokemon Pokemon) error {
	_, err := db.Exec(insertPokemonSQL, pokemon.ID, pokemon.Name, pokemon.Height, pokemon.Weight)

	if err != nil {
		return err
	}
	return nil
}
