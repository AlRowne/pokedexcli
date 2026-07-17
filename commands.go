package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"

	"github.com/AlRowne/pokedexcli/internal/pokeapi"
)

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func commandExit(cfg *config, s string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, s string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, com := range getCommands() {
		fmt.Printf("%s: %s\n", com.name, com.description)
	}
	return nil
}
func commandMap(cfg *config, s string) error {
	locationAreas, err := pokeapi.GetLocationAreas(cfg.Next)
	if err != nil {
		return err
	}

	for _, result := range locationAreas.Results {
		cfg.KnownLocations[result.Name] = struct{}{}
		fmt.Println(result.Name)
	}
	cfg.Next = stringOrEmpty(locationAreas.Next)
	cfg.Previous = stringOrEmpty(locationAreas.Previous)
	return nil
}

func commandMapb(cfg *config, s string) error {
	if cfg.Previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreas, err := pokeapi.GetLocationAreas(cfg.Previous)
	if err != nil {
		return err
	}

	for _, result := range locationAreas.Results {
		cfg.KnownLocations[result.Name] = struct{}{}
		fmt.Println(result.Name)
	}
	cfg.Next = stringOrEmpty(locationAreas.Next)
	cfg.Previous = stringOrEmpty(locationAreas.Previous)
	return nil
}

func commandExplore(cfg *config, s string) error {
	if s == "" {
		return errors.New("please provide an area to explore. use the map/mapb command")
	}
	fmt.Printf("Exploring %s...\n", s)
	response, err := pokeapi.GetPokemonEncounters(s)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")

	for _, encounter := range response.PokemonEncounters {
		cfg.KnownPokemon[encounter.Pokemon.Name] = struct{}{}
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, s string) error {
	if s == "" {
		return errors.New("please provide a Pokemon to catch. use the explore command to find pokemon")
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", s)
	pokemon, err := pokeapi.GetPokemon(s)
	if err != nil {
		return err
	}
	threshold := 40
	if pokemon.BaseExperience <= 0 || rand.Intn(pokemon.BaseExperience) < threshold {
		fmt.Printf("%s was caught!\n", pokemon.Name)
		cfg.Pokedex[pokemon.Name] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemon.Name)
	}
	return nil
}

func commandInspect(cfg *config, s string) error {
	if s == "" {
		return errors.New("no pokemon name provided")
	}
	pokemon, ok := cfg.Pokedex[s]
	if !ok {
		return errors.New("you have not caught that pokemon")
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(cfg *config, s string) error {
	if len(cfg.Pokedex) < 1 {
		return errors.New("your pokedex is empty")
	}
	var names []string
	for _, pokemon := range cfg.Pokedex {
		names = append(names, pokemon.Name)
	}
	sort.Strings(names)

	fmt.Println("Your Pokedex:")
	for _, name := range names {
		fmt.Printf(" - %s\n", name)
	}
	return nil
}

type cliCommand struct {
	name         string
	description  string
	callback     func(*config, string) error
	argCompleter func(*config) []string
}

func sortedKeys[V any](m map[string]V) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getCommands() map[string]cliCommand {
	cliCommands := map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 Location-Areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 Location-Areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Shows all Pokemon in the provided area",
			callback:    commandExplore,
			argCompleter: func(cfg *config) []string {
				return sortedKeys(cfg.KnownLocations)
			},
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon. If it's caught, it gets added to the PokeDex",
			callback:    commandCatch,
			argCompleter: func(cfg *config) []string {
				return sortedKeys(cfg.KnownPokemon)
			},
		},
		"inspect": {
			name:        "inspect",
			description: "Print various stats for a Pokemon that's already in the PokeDex",
			callback:    commandInspect,
			argCompleter: func(cfg *config) []string {
				return sortedKeys(cfg.Pokedex)
			},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Show your PokeDex",
			callback:    commandPokedex,
		},
	}
	return cliCommands
}
