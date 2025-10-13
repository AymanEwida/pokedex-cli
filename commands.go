package main

import (
	"errors"
	"fmt"
	"math/rand"
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

func CommandHelp(commands map[string]CliCommand) {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, v := range commands {
		fmt.Printf("%v: %v\n", v.name, v.description)
	}
}

func CommandMapf(config *Config, params []string) error {
	if config.Next == "" {
		return errors.New("you have explored all location areas in the all location areas in the Pokemon world")
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

	if _, ok := config.User.GetPokemon(pokemonName); ok {
		return fmt.Errorf("you already caught: %s", pokemonName)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := config.Client.CatchPokemon(pokemonName)
	if err != nil {
		return err
	}

	config.User.AddPokemon(internal.NewPokemon(&config.User.EVs, pokemon))

	fmt.Printf("%s was caught!\n", pokemonName)

	return nil
}

func CommandFight(config *Config, params []string) error {
	if len(params) == 0 {
		return errors.New("you need to provide <pokemon-name>\nsee 'help' command for more info")
	}

	if config.User.GetPokemonCount() == 0 {
		return errors.New("you have 0 pokemons in your pokedex, you need at least 1 pokemon to fight")
	}

	pokemonName := params[0]

	pokemon, err := config.Client.GetPokemonByName(pokemonName)
	if err != nil {
		return err
	}

	pokemonRarity := pokeapi.DecidedPokemonRarity(pokemon.BaseExperience)
	level := internal.GetRandomLevelFromRarity(pokemonRarity)

	isPokemonTurn := true
	if rand.Float64() <= 0.5 {
		isPokemonTurn = false
	}

	currentFigher, err := config.User.GetRandomPokemon()
	if err != nil {
		return err
	}

	fmt.Printf("%v(%v) - %v vs. Your Pokedex:\n", pokemonName, level, pokemonRarity)

	pokemonRealStats := internal.CalcAllPokemonRealStat(internal.PokemonBaseStats(pokemon.Stats), internal.GetRandomIVsToPokemon(), internal.StatsInfo[int]{}, level)

	// TODO: make the fight more real with time.Sleep

	fmt.Printf("%v(%v) - %v, HP: %.2f\n", pokemon.Name, level, pokemonRarity, pokemonRealStats.HP)
	fmt.Printf("> Current fight is: %v(%v) - %v, HP: %.2f\n", currentFigher.Name, currentFigher.Level, currentFigher.Rarity, currentFigher.RealStats.HP)

	for config.User.GetAlivePokemonCount() > 0 && pokemonRealStats.HP > 0 {
		if currentFigher.RealStats.HP <= 0 {
			fmt.Printf("> %v is deid\n", currentFigher.Name)
			fmt.Println("> Replacing current fighter....")

			currentFigher, err = config.User.GetRandomPokemon()
			if err != nil {
				return err
			}

			fmt.Printf("> Current fighter is: %v(%v) - %v, HP: %.2f\n", currentFigher.Name, currentFigher.Level, currentFigher.Rarity, currentFigher.RealStats.HP)
		}

		// TODO: convert the fight logic to a different function
		if !isPokemonTurn {
			fmt.Printf("> %v(%v) turn\n", currentFigher.Name, currentFigher.Level)

			isSpecialAttack := internal.DecideToDoSpecial(currentFigher.Rarity)
			attackDamage := currentFigher.RealStats.Attack
			if isSpecialAttack {
				fmt.Printf("> %v is doing special attack\n", currentFigher.Name)

				attackDamage = currentFigher.RealStats.SpecialAttack
			}

			fmt.Printf("> %v is attacking with: %.2f\n", currentFigher.Name, attackDamage)

			isDoingDefense := internal.DecideToDoDefense(pokemonRarity)
			defenseValue := 0.0

			if isDoingDefense {
				fmt.Printf("%v is defending\n", pokemon.Name)

				defenseValue = pokemonRealStats.Defense

				isSpecialDefense := internal.DecideToDoSpecial(pokemonRarity)
				if isSpecialDefense {
					fmt.Printf("%v is doing special defense\n", pokemon.Name)

					defenseValue = pokemonRealStats.SpecialDefense
				}

				fmt.Printf("%v is defending with: %.2f\n", pokemon.Name, defenseValue)
			}

			finalAttack := attackDamage - defenseValue
			if finalAttack < 0.0 {
				finalAttack = 0.0
			}

			fmt.Printf("finalAttack is: %.2f\n", finalAttack)

			pokemonRealStats.HP -= finalAttack
			if pokemonRealStats.HP < 0.0 {
				pokemonRealStats.HP = 0.0
			}
		} else {
			fmt.Printf("%v(%v) turn:\n", pokemon.Name, level)

			isSpecialAttack := internal.DecideToDoSpecial(pokemonRarity)
			attackDamage := pokemonRealStats.Attack
			if isSpecialAttack {
				fmt.Printf("%v is doing special attack\n", pokemon.Name)

				attackDamage = pokemonRealStats.SpecialAttack
			}

			fmt.Printf("%v is attacking with: %.2f\n", pokemon.Name, attackDamage)

			isDoingDefense := internal.DecideToDoDefense(currentFigher.Rarity)
			defenseValue := 0.0

			if isDoingDefense {
				fmt.Printf("> %v is defending\n", currentFigher.Name)

				defenseValue = currentFigher.RealStats.Defense

				isSpecialDefense := internal.DecideToDoSpecial(currentFigher.Rarity)
				if isSpecialDefense {
					fmt.Printf("> %v is doing special defense\n", currentFigher.Name)

					defenseValue = currentFigher.RealStats.SpecialDefense
				}

				fmt.Printf("> %v is defending with: %.2f\n", currentFigher.Name, defenseValue)
			}

			finalAttack := attackDamage - defenseValue
			if finalAttack < 0.0 {
				finalAttack = 0.0
			}

			fmt.Printf("finalAttack is: %.2f\n", finalAttack)

			currentFigher.RealStats.HP -= finalAttack
			if currentFigher.RealStats.HP < 0.0 {
				currentFigher.RealStats.HP = 0.0
			}
		}

		fmt.Printf("%v HP: %.2f\n> %v HP: %.2f\n", pokemon.Name, pokemonRealStats.HP, currentFigher.Name, currentFigher.RealStats.HP)

		isPokemonTurn = !isPokemonTurn
	}

	if config.User.GetAlivePokemonCount() <= 0 {
		fmt.Println("You Lost the battle!")

		return nil
	}

	fmt.Println("You won the battle!")

	config.User.AddEVs(internal.PokemonBaseStats(pokemon.Stats))
	config.User.CalcRealStatsToAllPokemons()

	return nil
}

// TODO: change inspect command to show the new info of the pokemon
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
		fmt.Printf("\t-%v(%v)\n", pokemon.Name, pokemon.Level)
	}

	return nil
}
