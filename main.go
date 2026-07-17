package main

import (
	"fmt"
	"io"
	"log"

	"github.com/AlRowne/pokedexcli/internal/pokeapi"
	"github.com/chzyer/readline"
)

type config struct {
	Next     string
	Previous string
	Pokedex  map[string]pokeapi.Pokemon
}

func main() {
	rl, err := readline.New("Pokedex > ")
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer rl.Close()

	cfg := config{}
	cfg.Pokedex = make(map[string]pokeapi.Pokemon)

	for {
		line, err := rl.Readline()
		if err == io.EOF || err == readline.ErrInterrupt {
			fmt.Println("Exiting...")
			break
		}
		if err != nil {
			log.Fatalf("error: %s", err)
		}
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
}
