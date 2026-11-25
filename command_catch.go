package main

import (
	"errors"
	"fmt"
	"math/rand"
)

// Cuanto más alto, más fácil es atrapar. Con 50, un Pokémon de experiencia
// base 100 se atrapa ~50% de las veces y uno de 300 ~17%.
const catchDifficulty = 50

func commandCatch(cfg *config, args []string) error {
	if len(args) < 1 {
		return errors.New("you must provide a pokemon name")
	}
	name := args[0]

	pokemon, err := cfg.pokeapiClient.GetPokemon(name)
	if err != nil {
		return err
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)

	// Una tirada entre 0 y la experiencia base: cuanto mayor es la experiencia,
	// menos probable es sacar un valor por debajo del umbral.
	roll := rand.Intn(max(pokemon.BaseExperience, 1))
	if roll >= catchDifficulty {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)
	fmt.Println("You may now inspect it with the inspect command.")
	cfg.caughtPokemon[pokemon.Name] = pokemon
	return nil
}
