package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

const pokemonURL = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
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
