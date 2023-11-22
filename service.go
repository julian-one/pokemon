package main

import (
	"database/sql"
	"strings"
)

type Pokemon struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Height int    `json:"height"`
	Weight int    `json:"weight"`
}

func getOnePokemonService(db *sql.DB, pokemonName string) (Pokemon, error) {
	lowerCasePokemonName := strings.ToLower(pokemonName)

	pokemon, err := getOnePokemonFromDB(db, lowerCasePokemonName)
	if err != nil {
		if err == sql.ErrNoRows {
			return searchAndSavePokemon(db, lowerCasePokemonName)
		}
		return Pokemon{}, err
	}
	return pokemon, nil
}

func searchAndSavePokemon(db *sql.DB, pokemonName string) (Pokemon, error) {
	pokemon, err := searchAPIByName(pokemonName)
	if err != nil {
		return Pokemon{}, err
	}

	err = savePokemonToDB(db, pokemon)
	if err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

func getAllPokemonService(db *sql.DB) ([]Pokemon, error) {
	pokemon, err := getAllPokemonFromDB(db)
	if err != nil {
		return nil, err
	}
	return pokemon, nil
}
