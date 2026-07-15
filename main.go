package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/AlRowne/pokedexcli/internal/pokeapi"
)

type config struct {
	Next     string
	Previous string
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
	}
	return cliCommands
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	cfg := config{}

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
		explArea := ""
		if command == "explore" && len(words) < 2 {
			fmt.Println("Please provide an area to explore. (Use the map/mapb command)")
			continue
		} else if command == "explore" && len(words) >= 2 {
			explArea = words[1]
		}

		commands := getCommands()
		val, ok := commands[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := val.callback(&cfg, explArea); err != nil {
			fmt.Println(err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner Error: %v\n", err)
	}
}
