package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

const pokemonURL = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name           string        `json:"name"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
	Moves          []PokemonMove `json:"moves"`
}

type PokemonStat struct {
	BaseStat int `json:"base_stat"`
	Stat     struct {
		Name string `json:"name"`
	} `json:"stat"`
}
type PokemonType struct {
	Type struct {
		Name string `json:"name"`
	} `json:"type"`
}

type PokemonMove struct {
	Move struct {
		Name string `json:"name"`
	} `json:"move"`
}

func GetPokemon(name string) (Pokemon, error) {
	url := pokemonURL + name + "/"

	cacheRes, ok := cache.Get(url)
	if ok {
		var pokemon Pokemon
		if err := json.Unmarshal(cacheRes, &pokemon); err != nil {
			return Pokemon{}, err
		}
		return pokemon, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	var pokemon Pokemon
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}
	if err := json.Unmarshal(data, &pokemon); err != nil {
		return Pokemon{}, err
	}
	cache.Add(url, data)

	return pokemon, nil
}
