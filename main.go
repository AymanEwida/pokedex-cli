package main

import (
	"fmt"
	"log"
	"time"

	"github.com/AymanEwida/pokedex-cli/internal"
	"github.com/AymanEwida/pokedex-cli/internal/pokeapi"
	"github.com/AymanEwida/pokedex-cli/lib"
	"github.com/eiannone/keyboard"
)

func main() {
	err := keyboard.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer keyboard.Close()

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

		"fight": {
			name:        "fight <pokemon-name>",
			description: "fight a pokemon with user pokemon collection and earn EV points to grow you pokemons stats",
			callback:    CommandFight,
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

	prevCommands := lib.NewStack[string]()
	redoCommands := lib.NewStack[string]()

	buffer := ""
	hittedKey := ""

	fmt.Print("Pokedex > ")

	for {
		fmt.Printf("%v", hittedKey)

		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println(err)

			continue
		}

		if key == 127 {
			hittedKey = ""

			if len(buffer) > 0 {
				buffer = buffer[:len(buffer)-1]
			}

			fmt.Printf("\rPokedex > %v", buffer)
			fmt.Print(" ")
			fmt.Printf("\rPokedex > %v", buffer)

			continue
		}

		switch key {
		case keyboard.KeyCtrlC:
			return

		case keyboard.KeyArrowUp:
			hittedKey = ""

			prevBufferSize := len(buffer)

			prevCommand, _ := prevCommands.Pop()
			if len(prevCommand) > 0 {
				redoCommands.Push(prevCommand)
			}

			buffer = prevCommand

			fmt.Printf("\rPokedex > %v", buffer)
			for range prevBufferSize {
				fmt.Print(" ")
			}
			fmt.Printf("\rPokedex > %v", buffer)

			continue

		case keyboard.KeyArrowDown:
			hittedKey = ""

			prevBufferSize := len(buffer)

			redoCommand, _ := redoCommands.Pop()
			if len(redoCommand) > 0 {
				prevCommands.Push(redoCommand)
			}

			buffer, _ = redoCommands.Peek()

			fmt.Printf("\rPokedex > %v", buffer)
			for range prevBufferSize {
				fmt.Print(" ")
			}
			fmt.Printf("\rPokedex > %v", buffer)

			continue

		case keyboard.KeySpace:
			buffer += " "
			hittedKey = " "

			continue

		default:
			if key != keyboard.KeyEnter {
				buffer += string(char)
				hittedKey = string(char)

				continue
			} else {
				fmt.Println()
			}
		}

		words := lib.CleanInput(buffer)

		if len(words) == 0 {
			buffer = ""
			hittedKey = ""

			fmt.Print("Pokedex > ")

			continue
		}

		command, ok := commands[words[0]]

		if !ok {
			fmt.Printf("Unknown command: %v\n", words[0])
		} else {
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

		lib.MoveStackToStack(&redoCommands, &prevCommands)
		prevCommands.Push(buffer)

		buffer = ""
		hittedKey = ""

		fmt.Print("Pokedex > ")
	}
}
