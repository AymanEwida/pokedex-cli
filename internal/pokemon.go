package internal

import (
	"math/rand"

	"github.com/AymanEwida/pokedex-cli/internal/pokeapi"
)

type StatsInfo[T int | float64] struct {
	HP             T
	Attack         T
	Defense        T
	SpecialAttack  T
	SpecialDefense T
	Speed          T
}

type Pokemon struct {
	Level     int
	IVs       StatsInfo[int]
	RealStats StatsInfo[float64]
	Rarity    pokeapi.PokemonRarity
	pokeapi.Pokemon
}

type PokemonBaseStats = []struct {
	BaseStat int
	Effort   int
	Stat     struct{ Name string }
}

func DecideToDoDefense(rarity pokeapi.PokemonRarity) bool {
	var probability float64
	switch rarity {
	case pokeapi.LOW:
		probability = 10.0 / 100.0
	case pokeapi.MID:
		probability = 25.0 / 100.0
	case pokeapi.HIGH:
		probability = 35.0 / 100.0
	case pokeapi.LEGENDARY:
		probability = 45.0 / 100.0
	default:
		probability = 0.0
	}

	if rand.Float64() <= probability {
		return true
	}

	return false
}

func DecideToDoSpecial(rarity pokeapi.PokemonRarity) bool {
	var probability float64
	switch rarity {
	case pokeapi.LOW:
		probability = 15.0 / 100.0
	case pokeapi.MID:
		probability = 30.0 / 100.0
	case pokeapi.HIGH:
		probability = 50.0 / 100.0
	case pokeapi.LEGENDARY:
		probability = 65.0 / 100.0
	default:
		probability = 0.0
	}

	if rand.Float64() <= probability {
		return true
	}

	return false
}

func GetRandomLevelFromRarity(rarity pokeapi.PokemonRarity) int {
	switch rarity {
	case pokeapi.LOW:
		return rand.Intn(10) + 1
	case pokeapi.MID:
		return rand.Intn(11) + 10
	case pokeapi.HIGH:
		return rand.Intn(11) + 20
	case pokeapi.LEGENDARY:
		return rand.Intn(11) + 30
	default:
		return 0
	}
}

func CalcPokemonRealStat(baseStat, iv, ev, level int, isHpStat bool) float64 {
	var realStat float64 = (((2.0 * float64(baseStat)) + float64(iv) + (float64(ev) / 4.0)) / 100.0) * float64(level)

	if isHpStat {
		realStat += float64(level) + 10.0
	} else {
		realStat += 5.0
	}

	return realStat
}

func GetRandomIVsToPokemon() StatsInfo[int] {
	return StatsInfo[int]{
		HP:             rand.Intn(32),
		Attack:         rand.Intn(32),
		Defense:        rand.Intn(32),
		SpecialAttack:  rand.Intn(32),
		SpecialDefense: rand.Intn(32),
		Speed:          rand.Intn(32),
	}
}

func CalcAllPokemonRealStat(
	pokemonBaseStats PokemonBaseStats,
	ivs,
	evs StatsInfo[int],
	level int,
) StatsInfo[float64] {
	return StatsInfo[float64]{
		HP:             CalcPokemonRealStat(pokemonBaseStats[0].BaseStat, ivs.HP, evs.HP, level, true),
		Attack:         CalcPokemonRealStat(pokemonBaseStats[1].BaseStat, ivs.Attack, evs.Attack, level, false),
		Defense:        CalcPokemonRealStat(pokemonBaseStats[2].BaseStat, ivs.Defense, evs.Defense, level, false),
		SpecialAttack:  CalcPokemonRealStat(pokemonBaseStats[3].BaseStat, ivs.SpecialAttack, evs.SpecialAttack, level, false),
		SpecialDefense: CalcPokemonRealStat(pokemonBaseStats[4].BaseStat, ivs.SpecialDefense, evs.SpecialDefense, level, false),
		Speed:          CalcPokemonRealStat(pokemonBaseStats[5].BaseStat, ivs.Speed, evs.Speed, level, false),
	}
}

func NewPokemon(userEVs *StatsInfo[int], pokemon pokeapi.Pokemon) Pokemon {
	ivs := GetRandomIVsToPokemon()

	level := GetRandomLevelFromRarity(pokeapi.DecidedPokemonRarity(pokemon.BaseExperience))

	return Pokemon{
		Level:     level,
		IVs:       ivs,
		RealStats: CalcAllPokemonRealStat(PokemonBaseStats(pokemon.Stats), ivs, *userEVs, level),
		Rarity:    pokeapi.DecidedPokemonRarity(pokemon.BaseExperience),
		Pokemon:   pokemon,
	}
}
