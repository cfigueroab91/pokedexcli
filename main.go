package main

import (
	"time"

	"github.com/cfigueroab91/pokedexcli/internal/pokeapi"
)

func main() {
	startURL := pokeapi.LocationAreaURL

	cfg := &config{
		commands:         getCommands(),
		pokeapiClient:    pokeapi.NewClient(5*time.Second, 5*time.Minute),
		nextLocationsURL: &startURL,
		caughtPokemon:    make(map[string]pokeapi.Pokemon),
	}
	startRepl(cfg)
}
