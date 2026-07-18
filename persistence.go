package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/AlRowne/pokedexcli/internal/pokeapi"
)

type saveData struct {
	Pokedex        map[string]pokeapi.Pokemon `json:"pokedex"`
	KnownLocations map[string]struct{}        `json:"known_locations"`
	KnownPokemon   map[string]struct{}        `json:"known_pokemon"`
}

func getSavePath() (string, error) {
	path, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	path = filepath.Join(path, "pokedexcli", "pokedex.json")
	return path, nil
}

func saveState(cfg *config) error {
	data, err := json.Marshal(saveData{
		Pokedex:        cfg.Pokedex,
		KnownLocations: cfg.KnownLocations,
		KnownPokemon:   cfg.KnownPokemon,
	})
	if err != nil {
		return err
	}
	savePath, err := getSavePath()
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(savePath), 0o755); err != nil {
		return err
	}
	if err = os.WriteFile(savePath, data, 0o644); err != nil {
		return err
	}
	return nil
}

func loadState(cfg *config) error {
	path, err := getSavePath()
	if err != nil {
		return err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}
	var sd saveData
	if err = json.Unmarshal(data, &sd); err != nil {
		return err
	}
	if sd.Pokedex != nil {
		cfg.Pokedex = sd.Pokedex
	}
	if sd.KnownLocations != nil {
		cfg.KnownLocations = sd.KnownLocations
	}
	if sd.KnownPokemon != nil {
		cfg.KnownPokemon = sd.KnownPokemon
	}
	return nil
}
