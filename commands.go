package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

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

func fightTurn(fightTurnOps struct {
	attacker             *internal.Pokemon
	attackerPrefixString string
	defender             *internal.Pokemon
	defenderPrefixString string
},
) {
	fmt.Printf(fightTurnOps.attackerPrefixString+"%v(%v) turn\n", fightTurnOps.attacker.Name, fightTurnOps.attacker.Level)

	time.Sleep(1 * time.Second)
	fmt.Println()

	isSpecialAttack := internal.DecideToDoSpecial(fightTurnOps.attacker.Rarity)
	attackDamage := fightTurnOps.attacker.RealStats.Attack
	if isSpecialAttack {
		fmt.Printf(fightTurnOps.attackerPrefixString+"%v is doing special attack\n", fightTurnOps.attacker.Name)

		time.Sleep(500 * time.Millisecond)

		attackDamage = fightTurnOps.attacker.RealStats.SpecialAttack
	}

	fmt.Printf(fightTurnOps.attackerPrefixString+"%v is attacking with: %.2f\n", fightTurnOps.attacker.Name, attackDamage)

	time.Sleep(500 * time.Millisecond)
	fmt.Println()

	isDoingDefense := internal.DecideToDoDefense(fightTurnOps.defender.Rarity)
	defenseValue := 0.0

	if isDoingDefense {
		fmt.Printf(fightTurnOps.defenderPrefixString+"%v is defending\n", fightTurnOps.defender.Name)

		time.Sleep(500 * time.Millisecond)

		defenseValue = fightTurnOps.defender.RealStats.Defense

		isSpecialDefense := internal.DecideToDoSpecial(fightTurnOps.defender.Rarity)
		if isSpecialDefense {
			fmt.Printf(fightTurnOps.defenderPrefixString+"%v is doing special defense\n", fightTurnOps.defender.Name)

			time.Sleep(500 * time.Millisecond)

			defenseValue = fightTurnOps.defender.RealStats.SpecialDefense
		}

		fmt.Printf(fightTurnOps.defenderPrefixString+"%v is defending with: %.2f\n", fightTurnOps.defender.Name, defenseValue)

		time.Sleep(500 * time.Millisecond)
		fmt.Println()
	} else {
		fmt.Printf(fightTurnOps.defenderPrefixString+"%v is not defending\n", fightTurnOps.defender.Name)

		time.Sleep(500 * time.Millisecond)
		fmt.Println()
	}

	finalAttack := attackDamage - defenseValue
	if finalAttack < 0.0 {
		finalAttack = 0.0
	}

	fmt.Printf("finalAttack is: %.2f\n", finalAttack)

	time.Sleep(500 * time.Millisecond)
	fmt.Println()

	fightTurnOps.defender.RealStats.HP -= finalAttack
	if fightTurnOps.defender.RealStats.HP < 0.0 {
		fightTurnOps.defender.RealStats.HP = 0.0
	}
}

func fight(fightOps struct {
	fighterA             *internal.Pokemon
	fighterB             *internal.Pokemon
	isFighterATurn       bool
	fighterBPrefixString string
},
) {
	if !fightOps.isFighterATurn {
		fightTurn(struct {
			attacker             *internal.Pokemon
			attackerPrefixString string
			defender             *internal.Pokemon
			defenderPrefixString string
		}{
			attacker:             fightOps.fighterB,
			attackerPrefixString: fightOps.fighterBPrefixString,
			defender:             fightOps.fighterA,
			defenderPrefixString: "",
		})
	} else {
		fightTurn(struct {
			attacker             *internal.Pokemon
			attackerPrefixString string
			defender             *internal.Pokemon
			defenderPrefixString string
		}{
			attacker:             fightOps.fighterA,
			attackerPrefixString: "",
			defender:             fightOps.fighterB,
			defenderPrefixString: fightOps.fighterBPrefixString,
		})
	}
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
	ivs := internal.GetRandomIVsToPokemon()
	pokemonRealStats := internal.CalcAllPokemonRealStat(internal.PokemonBaseStats(pokemon.Stats), ivs, internal.StatsInfo[int]{}, level)

	fighting := internal.Pokemon{
		Rarity:    pokemonRarity,
		Level:     level,
		Pokemon:   pokemon,
		RealStats: pokemonRealStats,
		IVs:       ivs,
	}

	isPokemonTurn := true
	if rand.Float64() <= 0.5 {
		isPokemonTurn = false
	}

	currentFigher, err := config.User.GetRandomPokemon()
	if err != nil {
		return err
	}

	fmt.Printf("%v(%v) - %v vs. > Your Pokedex:\n", fighting.Name, fighting.Level, fighting.Rarity)

	time.Sleep(1 * time.Second)
	fmt.Println()

	fmt.Printf("%v(%v) - %v, HP: %.2f\n", fighting.Name, fighting.Level, fighting.Rarity, fighting.RealStats.HP)
	fmt.Printf("> Current fight is: %v(%v) - %v, HP: %.2f\n", currentFigher.Name, currentFigher.Level, currentFigher.Rarity, currentFigher.RealStats.HP)

	time.Sleep(1 * time.Second)
	fmt.Println()

	for config.User.GetAlivePokemonCount() > 0 && fighting.RealStats.HP > 0 {
		if currentFigher.RealStats.HP <= 0 {
			fmt.Printf("> %v is deid\n", currentFigher.Name)

			time.Sleep(1 * time.Second)
			fmt.Println()

			fmt.Println("> Replacing current fighter....")

			time.Sleep(500 * time.Millisecond)
			fmt.Println()

			currentFigher, err = config.User.GetRandomPokemon()
			if err != nil {
				return err
			}

			fmt.Printf("> Current fighter is: %v(%v) - %v, HP: %.2f\n", currentFigher.Name, currentFigher.Level, currentFigher.Rarity, currentFigher.RealStats.HP)

			time.Sleep(1 * time.Second)
			fmt.Println()
		}

		fight(struct {
			fighterA             *internal.Pokemon
			fighterB             *internal.Pokemon
			isFighterATurn       bool
			fighterBPrefixString string
		}{
			fighterA:             &fighting,
			fighterB:             currentFigher,
			isFighterATurn:       isPokemonTurn,
			fighterBPrefixString: "> ",
		})

		fmt.Printf("%v HP: %.2f\n> %v HP: %.2f\n", fighting.Name, fighting.RealStats.HP, currentFigher.Name, currentFigher.RealStats.HP)
		fmt.Println("===============================")

		time.Sleep(1 * time.Second)

		isPokemonTurn = !isPokemonTurn
	}

	if config.User.GetAlivePokemonCount() <= 0 {
		fmt.Println("You Lost the battle, no Evs earned!")

		return nil
	}

	fmt.Println("You won the battle!")
	fmt.Printf("Evs earned from the battle with %v(%v):\n", fighting.Name, fighting.Level)
	fmt.Printf("\t- HP = %v Ev\n", fighting.Stats[0].Effort)
	fmt.Printf("\t- Attack = %v Ev\n", fighting.Stats[1].Effort)
	fmt.Printf("\t- Defense = %v Ev\n", fighting.Stats[2].Effort)
	fmt.Printf("\t- SpecialAttack = %v Ev\n", fighting.Stats[3].Effort)
	fmt.Printf("\t- SpecialDefense = %v Ev\n", fighting.Stats[4].Effort)
	fmt.Printf("\t- Speed = %v Ev\n", fighting.Stats[5].Effort)

	config.User.AddEVs(internal.PokemonBaseStats(pokemon.Stats))
	config.User.CalcRealStatsToAllPokemons()

	fmt.Println()
	fmt.Println("Use the inspect command to see the new stats for your pokemons\ntype 'help' for more info to the inspect command")

	return nil
}

func CommandSimulatFight(config *Config, params []string) error {
	if len(params) < 2 {
		return errors.New("you need to provide: <fighterA-name> <fighterB-name>\nsee 'help' command for more info")
	}

	fighterAName := params[0]
	fighterBName := params[1]

	fighterA, ok := config.User.GetPokemon(fighterAName)
	if !ok {
		return fmt.Errorf("%v is not in your pokedex\nboth of the fighters need to be in your pokedex\nuse the 'catch' command to catch pokemons", fighterAName)
	}

	fighterB, ok := config.User.GetPokemon(fighterBName)
	if !ok {
		return fmt.Errorf("%v is not in your pokedex\nboth of the fighters need to be in your pokedex\nuse the 'catch' command to catch pokemons", fighterBName)
	}

	fighterA.Level = 50
	fighterA.RealStats = internal.CalcAllPokemonRealStat(internal.PokemonBaseStats(fighterA.Stats), fighterA.IVs, config.User.EVs, 50)

	fighterB.Level = 50
	fighterB.RealStats = internal.CalcAllPokemonRealStat(internal.PokemonBaseStats(fighterB.Stats), fighterB.IVs, config.User.EVs, 50)

	isFighterATurn := true
	if rand.Float64() <= 0.5 {
		isFighterATurn = false
	}

	fmt.Printf("%v(%v) - %v vs. %v(%v) - %v:\n", fighterA.Name, fighterA.Level, fighterA.Rarity, fighterB.Name, fighterB.Level, fighterB.Rarity)

	time.Sleep(1 * time.Second)
	fmt.Println()

	fmt.Printf("%v(%v) - %v, HP: %.2f\n", fighterA.Name, fighterA.Level, fighterA.Rarity, fighterA.RealStats.HP)
	fmt.Printf("%v(%v) - %v, HP: %.2f\n", fighterB.Name, fighterB.Level, fighterB.Rarity, fighterB.RealStats.HP)

	time.Sleep(1 * time.Second)
	fmt.Println()

	for fighterA.RealStats.HP > 0 && fighterB.RealStats.HP > 0 {
		fight(struct {
			fighterA             *internal.Pokemon
			fighterB             *internal.Pokemon
			isFighterATurn       bool
			fighterBPrefixString string
		}{
			fighterA:             &fighterA,
			fighterB:             &fighterB,
			isFighterATurn:       isFighterATurn,
			fighterBPrefixString: "",
		})

		fmt.Printf("%v HP: %.2f\n%v HP: %.2f\n", fighterA.Name, fighterA.RealStats.HP, fighterB.Name, fighterB.RealStats.HP)
		fmt.Println("===============================")

		time.Sleep(1 * time.Second)

		isFighterATurn = !isFighterATurn
	}

	time.Sleep(1 * time.Second)
	fmt.Println()

	if fighterA.RealStats.HP < 0 {
		fmt.Printf("%v(%v) - %v won!\n", fighterB.Name, fighterB.Level, fighterB.Rarity)
	} else {
		fmt.Printf("%v(%v) - %v won!\n", fighterA.Name, fighterA.Level, fighterA.Rarity)
	}

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
	fmt.Printf("Level: %v\n", pokemon.Level)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)
	fmt.Printf("Base Experience: %v\n", pokemon.BaseExperience)
	fmt.Printf("Rarity: %v\n", pokemon.Rarity)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("\t-%v: %v, Ev: %v\n", stat.Stat.Name, stat.BaseStat, stat.Effort)
	}

	fmt.Println("Real Stat:")
	fmt.Printf("\t-HP: %v\n", pokemon.RealStats.HP)
	fmt.Printf("\t-Attack: %v\n", pokemon.RealStats.Attack)
	fmt.Printf("\t-Defense: %v\n", pokemon.RealStats.Defense)
	fmt.Printf("\t-SpecialAttack: %v\n", pokemon.RealStats.SpecialAttack)
	fmt.Printf("\t-SpecialDefense: %v\n", pokemon.RealStats.SpecialDefense)
	fmt.Printf("\t-Speed: %v\n", pokemon.RealStats.Speed)

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

	fmt.Println()

	fmt.Println("Your Evs:")
	fmt.Printf("\t-HP: %v\n", config.User.EVs.HP)
	fmt.Printf("\t-Attack: %v\n", config.User.EVs.Attack)
	fmt.Printf("\t-Defense: %v\n", config.User.EVs.Defense)
	fmt.Printf("\t-SpecialAttack: %v\n", config.User.EVs.SpecialAttack)
	fmt.Printf("\t-SpecialDefense: %v\n", config.User.EVs.SpecialDefense)
	fmt.Printf("\t-Speed: %v\n", config.User.EVs.Speed)

	return nil
}
