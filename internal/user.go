package internal

import "github.com/AymanEwida/pokedex-cli/internal/pokeapi"

type User struct {
	Pokedex map[string]pokeapi.Pokemon
}

func NewUser() User {
	return User{
		Pokedex: map[string]pokeapi.Pokemon{},
	}
}

func (u *User) AddPokemon(pokemon pokeapi.Pokemon) {
	u.Pokedex[pokemon.Name] = pokemon
}

func (u *User) GetPokemon(pokemonName string) (pokeapi.Pokemon, bool) {
	pokemon, ok := u.Pokedex[pokemonName]

	return pokemon, ok
}
