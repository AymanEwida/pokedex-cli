package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/AymanEwida/pokedex-cli/internal"
	"github.com/AymanEwida/pokedex-cli/internal/pokeapi"
)

type CliCommand struct {
	name        string
	description string
	callback    func(*Config, []string) error
}

type Config struct {
	Client   pokeapi.Client
	Next     string
	Previous string
	User     internal.User
}

func CommandExit(config *Config, params []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")

	os.Exit(0)

	return nil
}

func CommandHelp(commands map[string]CliCommand) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, v := range commands {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}

	return nil
}

func CommandMapf(config *Config, params []string) error {
	if config.Next == "" {
		return errors.New("you have explored all location areas in the  all location areas in the Pokemon world")
	}

	areaLocations, err := config.Client.GetAreaLocations(config.Next)
	if err != nil {
		return err
	}

	config.Next = areaLocations.Next
	config.Previous = areaLocations.Previous

	for _, location := range areaLocations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func CommandMapb(config *Config, params []string) error {
	if config.Previous == "" {
		return errors.New("you're on the first page")
	}

	areaLocations, err := config.Client.GetAreaLocations(config.Previous)
	if err != nil {
		return err
	}

	config.Next = areaLocations.Next
	config.Previous = areaLocations.Previous

	for _, location := range areaLocations.Results {
		fmt.Println(location.Name)
	}

	return nil
}

func CommandExplore(config *Config, params []string) error {
	if len(params) == 0 {
		return errors.New("you need to provide <location-name>\nsee 'help' command for more info")
	}

	locationName := params[0]

	pokemonAreaLocation, err := config.Client.GetPokemonAreaLoaction(locationName)
	if err != nil {
		return err
	}

	fmt.Printf("Exploring %s...\n", locationName)
	fmt.Println("Found Pokemon:")

	for _, pokemon := range pokemonAreaLocation.PokemonEncounters {
		fmt.Printf("- %s\n", pokemon.Pokemon.Name)
	}

	return nil
}

func CommandCatch(config *Config, params []string) error {
	if len(params) == 0 {
		return errors.New("you need to provide <pokemon-name>\nsee 'help' command for more info")
	}

	pokemonName := params[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := config.Client.CatchPokemon(pokemonName)
	if err != nil {
		return err
	}

	config.User.AddPokemon(pokemon)

	fmt.Printf("%s was caught!\n", pokemonName)

	return nil
}

func CommandInspect(config *Config, params []string) error {
	if len(params) == 0 {
		return errors.New("you need to provide <pokemon-name>\nsee 'help' command for more info")
	}

	pokemonName := params[0]

	pokemon, ok := config.User.GetPokemon(pokemonName)
	if !ok {
		return fmt.Errorf("you have not caught: %s", pokemonName)
	}

	fmt.Printf("Name: %v\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Printf("Base Experience: %v\n", pokemon.BaseExperience)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("\t-%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("\t-%v\n", pokemonType.Type.Name)
	}

	return nil
}

func CommandPokedex(config *Config, params []string) error {
	if (len(config.User.Pokedex)) == 0 {
		return errors.New("you have not caught any pokemon yet!\nyou can catch pokemons with the 'catch' command")
	}

	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.User.Pokedex {
		fmt.Printf("\t-%v\n", pokemon.Name)
	}

	return nil
}
