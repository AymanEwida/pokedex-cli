package internal

import (
	"errors"
	"math/rand"
)

type User struct {
	EVs          StatsInfo[int]
	Pokedex      map[string]Pokemon
	PokemonNames []string
}

func NewUser() User {
	return User{
		PokemonNames: []string{},
		Pokedex:      map[string]Pokemon{},
		EVs:          StatsInfo[int]{},
	}
}

func (u *User) AddPokemon(pokemon Pokemon) {
	u.Pokedex[pokemon.Name] = pokemon
	u.PokemonNames = append(u.PokemonNames, pokemon.Name)
}

func (u *User) GetPokemon(pokemonName string) (Pokemon, bool) {
	pokemon, ok := u.Pokedex[pokemonName]

	return pokemon, ok
}

func (u *User) GetPokemonCount() int {
	return len(u.Pokedex)
}

func (u *User) GetAlivePokemonCount() int {
	count := 0

	for _, pokemon := range u.Pokedex {
		if pokemon.RealStats.HP > 0 {
			count++
		}
	}

	return count
}

func (u *User) GetRandomPokemon() (Pokemon, error) {
	if u.GetPokemonCount() == 0 {
		return Pokemon{}, errors.New("your pokedex is empty")
	}

	if u.GetAlivePokemonCount() == 0 {
		return Pokemon{}, errors.New("all of your pokemons deided")
	}

	pokemon, _ := u.GetPokemon(u.PokemonNames[rand.Intn(len(u.PokemonNames))])
	for pokemon.RealStats.HP == 0 {
		pokemon, _ = u.GetPokemon(u.PokemonNames[rand.Intn(len(u.PokemonNames))])
	}

	return pokemon, nil
}

func (u *User) AddEVs(baseStatsEvs PokemonBaseStats) {
	u.EVs.HP = baseStatsEvs[0].Effort
	u.EVs.Attack = baseStatsEvs[1].Effort
	u.EVs.Defense = baseStatsEvs[2].Effort
	u.EVs.SpecialAttack = baseStatsEvs[3].Effort
	u.EVs.SpecialDefense = baseStatsEvs[4].Effort
	u.EVs.Speed = baseStatsEvs[5].Effort
}

func (u *User) CalcRealStatsToAllPokemons() {
	for _, pokemon := range u.Pokedex {
		pokemon.RealStats = CalcAllPokemonRealStat(PokemonBaseStats(pokemon.Stats), pokemon.IVs, u.EVs, pokemon.Level)
	}
}
