package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
)

type ExploreResponse struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon PokemonInfo `json:"pokemon"`
}

type PokemonInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetPokemonEncounters(locationName string) (ExploreResponse, error) {
	url := locationAreaURL + locationName + "/"

	cacheRes, ok := cache.Get(url)
	if ok {
		var exploreResponse ExploreResponse
		if err := json.Unmarshal(cacheRes, &exploreResponse); err != nil {
			return ExploreResponse{}, err
		}
		return exploreResponse, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return ExploreResponse{}, err
	}
	defer res.Body.Close()

	var exploreResponse ExploreResponse
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return ExploreResponse{}, err
	}
	if err := json.Unmarshal(data, &exploreResponse); err != nil {
		return ExploreResponse{}, err
	}
	cache.Add(url, data)

	return exploreResponse, nil
}
