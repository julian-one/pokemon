package main

import (
	"encoding/json"

	"github.com/mtslzr/pokeapi-go"
	"github.com/mtslzr/pokeapi-go/structs"
)

func fetchPokemonData(pokemonName string) (structs.Pokemon, error) {
	pokemonData, err := pokeapi.Pokemon(pokemonName)
	if err != nil {
		return structs.Pokemon{}, err
	}
	return pokemonData, nil
}

func encodePokemonToJSON(pokemon structs.Pokemon) ([]byte, error) {
	jsonData, err := json.Marshal(pokemon)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}

func decodeJSONToPokemon(jsonData []byte) (Pokemon, error) {
	var pokemon Pokemon
	err := json.Unmarshal(jsonData, &pokemon)
	if err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}

func searchAPIByName(pokemonName string) (Pokemon, error) {
	pokemonData, err := fetchPokemonData(pokemonName)
	if err != nil {
		return Pokemon{}, err
	}

	jsonBytes, err := encodePokemonToJSON(pokemonData)
	if err != nil {
		return Pokemon{}, err
	}

	pokemon, err := decodeJSONToPokemon(jsonBytes)
	if err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
