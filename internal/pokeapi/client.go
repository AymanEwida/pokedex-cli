package pokeapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/AymanEwida/pokedex-cli/internal/pokecache"
)

const (
	baseAreaLoacationsUrl string = "https://pokeapi.co/api/v2/location-area/"
	basePokemonUrl        string = "https://pokeapi.co/api/v2/pokemon/"
)

type PokemonRarity string

const (
	LOW       PokemonRarity = "low"
	MID       PokemonRarity = "mid"
	HIGH      PokemonRarity = "high"
	LEGENDARY PokemonRarity = "legendary"
)

type Client struct {
	httpClient http.Client
	cache      pokecache.Cache
}

func NewClient(timeout, cacheInterval time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
		cache: pokecache.NewCache(cacheInterval),
	}
}

func (c *Client) GetAreaLocations(url string) (RespAreaLocations, error) {
	if val, ok := c.cache.Get(url); ok {
		var areaLocations RespAreaLocations
		if err := json.Unmarshal(val, &areaLocations); err != nil {
			return RespAreaLocations{}, err
		}

		return areaLocations, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespAreaLocations{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return RespAreaLocations{}, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return RespAreaLocations{}, err
	}

	var areaLocations RespAreaLocations
	if err := json.Unmarshal(data, &areaLocations); err != nil {
		return RespAreaLocations{}, err
	}

	c.cache.Add(url, data)

	return areaLocations, nil
}

func (c *Client) GetPokemonAreaLoaction(areaLocationName string) (RespPokemonAreaLocation, error) {
	url := baseAreaLoacationsUrl + areaLocationName

	if val, ok := c.cache.Get(url); ok {
		var pokemonAreaLocation RespPokemonAreaLocation
		if err := json.Unmarshal(val, &pokemonAreaLocation); err != nil {
			return RespPokemonAreaLocation{}, err
		}

		return pokemonAreaLocation, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return RespPokemonAreaLocation{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return RespPokemonAreaLocation{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		if res.StatusCode == 404 {
			return RespPokemonAreaLocation{}, fmt.Errorf("unknowon location: %s", areaLocationName)
		} else if res.StatusCode > 499 {
			return RespPokemonAreaLocation{}, errors.New("server error")
		} else {
			return RespPokemonAreaLocation{}, errors.New("got an error from the pokemon api")
		}
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return RespPokemonAreaLocation{}, err
	}

	var pokemonAreaLocation RespPokemonAreaLocation
	if err := json.Unmarshal(data, &pokemonAreaLocation); err != nil {
		return RespPokemonAreaLocation{}, err
	}

	c.cache.Add(url, data)

	return pokemonAreaLocation, nil
}

func (c *Client) GetPokemonByName(pokemonName string) (Pokemon, error) {
	url := basePokemonUrl + pokemonName

	if val, ok := c.cache.Get(url); ok {
		var pokemon Pokemon
		if err := json.Unmarshal(val, &pokemon); err != nil {
			return Pokemon{}, err
		}

		return pokemon, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}

	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	var pokemon Pokemon
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, data)

	return pokemon, nil
}

func DecidedPokemonRarity(baseExperience int) PokemonRarity {
	if baseExperience > 300 {
		return LEGENDARY
	} else if baseExperience >= 250 {
		return HIGH
	} else if baseExperience >= 100 {
		return MID
	}

	return LOW
}

func (c *Client) CatchPokemon(pokemonName string) (Pokemon, error) {
	pokemon, err := c.GetPokemonByName(pokemonName)
	if err != nil {
		return Pokemon{}, err
	}

	var probability float64
	switch DecidedPokemonRarity(pokemon.BaseExperience) {
	case LOW:
		probability = 80.0 / 100.0
	case MID:
		probability = 60.0 / 100.0
	case HIGH:
		probability = 40.0 / 100.0
	case LEGENDARY:
		probability = 20.0 / 100.0
	default:
		probability = 0.0
	}

	if rand.Float64() <= probability {
		return pokemon, nil
	}

	return Pokemon{}, fmt.Errorf("%s escaped", pokemonName)
}
