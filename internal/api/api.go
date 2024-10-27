package api

import (
	"pokemon/internal/database"

	"github.com/mtslzr/pokeapi-go"
)

func Search(name string) (*database.Pokemon, error) {
	data, err := pokeapi.Pokemon(name)
	if err != nil {
		return nil, err
	}

	p := &database.Pokemon{
		Id:     data.ID,
		Name:   data.Name,
		Height: float64(data.Height),
		Weight: float64(data.Weight),
	}

	return p, nil
}
