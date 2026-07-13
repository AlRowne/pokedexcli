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

func commandExit(cfg *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, com := range getCommands() {
		fmt.Printf("%s: %s\n", com.name, com.description)
	}
	return nil
}
func commandMap(cfg *config) error {
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

func commandMapb(cfg *config) error {
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

func stringOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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
			description: "displays a help message",
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
		commands := getCommands()
		val, ok := commands[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		if err := val.callback(&cfg); err != nil {
			fmt.Println(err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Scanner Error: %v\n", err)
	}
}
