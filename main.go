package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/AlRowne/pokedexcli/internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
	Pokedex  map[string]pokeapi.Pokemon
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

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, string) error
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
		},
		"catch": {
			name:        "catch",
			description: "Try to catch a pokemon. If it's caught, it gets added to the PokeDex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Print various stats for a Pokemon that's already in the Pokedex",
			callback:    commandInspect,
		},
	}
	return cliCommands
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := config{}
	cfg.Pokedex = make(map[string]pokeapi.Pokemon)

	for {
		fmt.Print("Pokedex > ")
		if ok := scanner.Scan(); !ok {
			break
		}
		line := scanner.Text()
		words := cleanInput(line)
		if len(words) == 0 {
			fmt.Println("Please enter a command")
			continue
		}
		command := words[0]
		argument := ""
		if len(words) > 1 {
			argument = words[1]
		}

		commands := getCommands()
		val, ok := commands[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := val.callback(&cfg, argument); err != nil {
			fmt.Println(err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner Error: %v\n", err)
	}
}
