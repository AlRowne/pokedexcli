package main

import (
	"fmt"
	"io"
	"log"

	"github.com/AlRowne/pokedexcli/internal/pokeapi"
	"github.com/chzyer/readline"
)

type config struct {
	Next           string
	Previous       string
	Pokedex        map[string]pokeapi.Pokemon
	KnownLocations map[string]struct{}
	KnownPokemon   map[string]struct{}
}

func main() {

	cfg := config{}
	cfg.Pokedex = make(map[string]pokeapi.Pokemon)
	cfg.KnownLocations = make(map[string]struct{})
	cfg.KnownPokemon = make(map[string]struct{})

	if err := loadState(&cfg); err != nil {
		log.Fatalf("error loading pokedex: %s", err)
	}

	commands := getCommands()
	var pcItems []readline.PrefixCompleterInterface
	for _, c := range commands {
		if c.argCompleter == nil {
			pcItems = append(pcItems, readline.PcItem(c.name))
		} else {
			pcItems = append(pcItems, readline.PcItem(c.name, readline.PcItemDynamic(func(s string) []string {
				return c.argCompleter(&cfg)
			})))
		}
	}
	completer := readline.NewPrefixCompleter(pcItems...)
	readlineCfg := readline.Config{
		Prompt:       "Pokedex > ",
		AutoComplete: completer,
	}
	rl, err := readline.NewEx(&readlineCfg)
	if err != nil {
		log.Fatalf("error: %s", err)
	}
	defer rl.Close()

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
