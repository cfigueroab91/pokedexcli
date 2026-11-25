package main

import (
	"fmt"
	"sort"
)

func commandPokedex(cfg *config, args []string) error {
	fmt.Println("Your Pokedex:")

	names := make([]string, 0, len(cfg.caughtPokemon))
	for name := range cfg.caughtPokemon {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}
