package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/AlRowne/pokedexcli/internal/pokecache"
)

const baseURL = "https://pokeapi.co/api/v2/location-area/"
const cacheInterval = 30 * time.Minute

var cache = pokecache.NewCache(cacheInterval)

type LocationAreaResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetLocationAreas(url string) (LocationAreaResponse, error) {
	if url == "" {
		url = baseURL
	}

	cacheRes, ok := cache.Get(url)
	if ok {
		var loc LocationAreaResponse
		if err := json.Unmarshal(cacheRes, &loc); err != nil {
			return LocationAreaResponse{}, err
		}
		return loc, nil
	}
	res, err := http.Get(url)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	defer res.Body.Close()

	var loc LocationAreaResponse
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreaResponse{}, err
	}
	if err := json.Unmarshal(data, &loc); err != nil {
		return LocationAreaResponse{}, err
	}
	cache.Add(url, data)

	return loc, nil
}
