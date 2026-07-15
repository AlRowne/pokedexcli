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
	Pokedex  map[string]pokeapi.Pokemon
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
