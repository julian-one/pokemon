# Pokemon API

A lightweight Go API that fetches Pokémon data from [PokeAPI](https://pokeapi.co/) and caches it in SQLite.

## Quick Start

```bash
go run main.go
```

Server starts on `http://localhost:8080`

## Endpoints

### `GET /pokemon?name={name}`
Fetch and save a Pokémon by name. Returns height and weight as JSON.

```bash
curl "http://localhost:8080/pokemon?name=pikachu"
```

### `GET /`
List all cached Pokémon.

```bash
curl http://localhost:8080/
```

## Features

- Fetches Pokémon data from PokeAPI
- Stores height and weight in SQLite
- Avoids duplicate saves
- JSON logging with structured output
