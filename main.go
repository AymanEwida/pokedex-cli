package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AymanEwida/pokedex-cli/internal"
	"github.com/AymanEwida/pokedex-cli/internal/pokeapi"
	"github.com/AymanEwida/pokedex-cli/lib"
	"github.com/eiannone/keyboard"
	"github.com/joho/godotenv"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func main() {
	godotenv.Load()

	err := keyboard.Open()
	if err != nil {
		log.Fatalln(err)
	}

	defer keyboard.Close()

	openaiClient := openai.NewClient(option.WithAPIKey(os.Getenv("OPENAI_API_KEY")))

	config := Config{
		Client:   pokeapi.NewClient(openaiClient, time.Duration(5*time.Second), time.Duration(30*60*time.Second)),
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

		"simulate-fight": {
			name:        "simulate-fight <fighterA-name> <fighterB-name>",
			description: "simulate fights between two pokemons from your Pokedex",
			callback:    CommandSimulatFight,
		},

		"mix": {
			name:        "mix <pokemon1-name> <pokemon2-name>",
			description: "mix two pokemons from your Pokedex to create a new pokemon",
			callback:    CommandMix,
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
	idx := len(buffer)

	for {
		fmt.Printf("\rPokedex > %v", buffer[:idx])

		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println(err)

			continue
		}

		if key == 127 {
			if len(buffer) > 0 {
				if idx == len(buffer) {
					buffer = buffer[:len(buffer)-1]
				} else {
					buffer = buffer[:idx-1] + buffer[idx:]
				}

				idx -= 1
			}

			fmt.Printf("\rPokedex > %v", buffer)
			fmt.Print(" ")
			fmt.Printf("\rPokedex > %v", buffer)

			continue
		}

		switch key {
		case keyboard.KeyCtrlC:
			return

		case keyboard.KeyArrowLeft:
			if idx > 0 {
				idx -= 1
			}

			continue

		case keyboard.KeyArrowRight:
			if idx < len(buffer) {
				idx += 1
			}

			continue

		case keyboard.KeyArrowUp:
			prevBufferSize := len(buffer)

			prevCommand, _ := prevCommands.Pop()
			if len(prevCommand) > 0 {
				redoCommands.Push(prevCommand)
			}

			buffer = prevCommand
			idx = len(buffer)

			fmt.Printf("\rPokedex > %v", buffer)
			for range prevBufferSize {
				fmt.Print(" ")
			}
			fmt.Printf("\rPokedex > %v", buffer)

			continue

		case keyboard.KeyArrowDown:
			prevBufferSize := len(buffer)

			redoCommand, _ := redoCommands.Pop()
			if len(redoCommand) > 0 {
				prevCommands.Push(redoCommand)
			}

			buffer, _ = redoCommands.Peek()
			idx = len(buffer)

			fmt.Printf("\rPokedex > %v", buffer)
			for range prevBufferSize {
				fmt.Print(" ")
			}
			fmt.Printf("\rPokedex > %v", buffer)

			continue

		case keyboard.KeySpace:
			if idx == len(buffer) {
				buffer += " "
			} else {
				buffer = buffer[:idx] + " " + buffer[idx:]
			}

			idx += 1

			fmt.Printf("\rPokedex > %v", buffer)

			continue

		default:
			if key != keyboard.KeyEnter {
				if idx == len(buffer) {
					buffer += string(char)
				} else {
					buffer = buffer[:idx] + string(char) + buffer[idx:]
				}

				idx += 1

				fmt.Printf("\rPokedex > %v", buffer)

				continue
			} else {
				fmt.Println()
			}
		}

		words := lib.CleanInput(buffer)

		if len(words) == 0 {
			buffer = ""
			idx = len(buffer)

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
		idx = len(buffer)
	}
}
