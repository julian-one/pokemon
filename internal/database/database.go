package database

import (
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func Initialize() (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite3", "pokemon.db")
	if err != nil {
		return nil, err
	}

	f, err := os.ReadFile("model.sql")
	if err != nil {
		return nil, err
	}
	content := string(f)

	_, err = db.Exec(content)
	if err != nil {
		return nil, err
	}

	return db, nil
}

type Pokemon struct {
	Id     int     `json:"pokemon_id" db:"pokemon_id"`
	Name   string  `json:"name"       db:"name"`
	Height float64 `json:"height"     db:"height"`
	Weight float64 `json:"weight"     db:"weight"`
}

func GetAll(db *sqlx.DB) ([]Pokemon, error) {
	pokemon := make([]Pokemon, 0)
	err := db.Select(&pokemon, `SELECT * FROM pokemon`)
	if err != nil {
		return nil, err
	}
	return pokemon, nil
}

func Exists(db *sqlx.DB, id int) (bool, error) {
	var exists bool
	err := db.Get(&exists, `SELECT EXISTS(SELECT 1 FROM pokemon WHERE pokemon_id = ?)`, id)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func Save(db *sqlx.DB, p *Pokemon) error {
	_, err := db.Exec(
		`INSERT INTO pokemon (pokemon_id, name, height, weight) VALUES (?, ?, ?, ?)`,
		p.Id, p.Name, p.Height, p.Weight)
	if err != nil {
		return err
	}
	return nil
}
