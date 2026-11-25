package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cfigueroab91/pokedexcli/internal/pokeapi"
)

type config struct {
	commands         map[string]cliCommand
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon    map[string]pokeapi.Pokemon
}

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, args []string) error
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore <area_name>",
			description: "Lists the Pokemon found in a location area",
			callback:    commandExplore,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempts to catch a Pokemon and add it to your Pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "Shows details of a Pokemon you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Lists all the Pokemon you have caught",
			callback:    commandPokedex,
		},
	}
}

func startRepl(cfg *config) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")

		if !scanner.Scan() {
			break
		}

		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}

		command, ok := cfg.commands[words[0]]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}

		if err := command.callback(cfg, words[1:]); err != nil {
			fmt.Println(err)
		}
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
