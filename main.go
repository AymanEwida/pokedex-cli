package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AymanEwida/pokedex-cli/internal"
	"github.com/AymanEwida/pokedex-cli/internal/pokeapi"
)

func cleanInput(text string) []string {
	if len(text) == 0 {
		return []string{}
	}

	text = strings.ToLower(text)
	words := strings.Fields(text)

	return words
}

func main() {
	config := Config{
		Client:   pokeapi.NewClient(time.Duration(5*time.Second), time.Duration(30*60*time.Second)),
		Next:     "https://pokeapi.co/api/v2/location-area?offset=0&limit=20",
		Previous: "",
		User:     internal.NewUser(),
	}

	commands := map[string]CliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
		},

		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    CommandExit,
		},

		"map": {
			name:        "map",
			description: "Explore the Pokemon world by getting the next 20 location areas in the Pokemon world",
			callback:    CommandMapf,
		},

		"mapb": {
			name:        "mapb",
			description: "Explore previous 20 location areas in the Pokemon world",
			callback:    CommandMapb,
		},

		"explore": {
			name:        "explore <location-name>",
			description: "explore command takes a <loaction-name> parameter and returns list of all the Pokémon located there",
			callback:    CommandExplore,
		},

		"catch": {
			name:        "catch <pokemon-name>",
			description: "try to catch a pokemon, you need to provide a <pokemon-name> param",
			callback:    CommandCatch,
		},

		"inspect": {
			name:        "inspect <pokemon-name>",
			description: "inspect a pokemon you have caught, this command takes a <pokemon-name> param",
			callback:    CommandInspect,
		},

		"pokedex": {
			name:        "pokedex",
			description: "prints all the pokemons names you have caught",
			callback:    CommandPokedex,
		},
	}

	reader := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())

		if len(words) == 0 {
			continue
		}

		command, ok := commands[words[0]]

		if !ok {
			fmt.Printf("Unknown command: %v\n", words[0])

			continue
		}

		if command.name == "help" {
			CommandHelp(commands)
		} else {
			params := words[1:]

			err := command.callback(&config, params)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
