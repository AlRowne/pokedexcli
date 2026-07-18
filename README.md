# pokedexcli

A small interactive Pokédex REPL for the terminal, written in Go. Explore
location areas, catch Pokémon, and inspect the ones you've caught — all from a
shell-style prompt with command history and tab-completion.

Data comes from the [PokéAPI](https://pokeapi.co/).

This project started as part of the [boot.dev](https://www.boot.dev) "Backend
Developer" path and has since grown a few self-directed extras (readline history,
dynamic tab-completion).

## Install & run

Requires Go 1.26 or newer.

```bash
git clone https://github.com/AlRowne/pokedexcli.git
cd pokedexcli
go run .
```

Or build a binary:

```bash
go build -o pokedexcli .
./pokedexcli
```

## Usage

You'll be dropped into a prompt:

```
Pokedex > help
```

### Commands

| Command             | Description                                                        |
| ------------------- | ------------------------------------------------------------------ |
| `help`              | Show the list of commands                                          |
| `map`               | Show the next 20 location areas                                    |
| `mapb`              | Show the previous 20 location areas                                |
| `explore <area>`    | List the Pokémon that can be found in a location area              |
| `catch <pokemon>`   | Attempt to catch a Pokémon (chance based on its base experience)   |
| `inspect <pokemon>` | Show stats for a Pokémon you've already caught                     |
| `pokedex`           | List every Pokémon in your Pokédex                                 |
| `exit`              | Close the Pokédex                                                  |

### A typical session

```
Pokedex > map
canalave-city-area
eterna-city-area
...
Pokedex > explore canalave-city-area
Exploring canalave-city-area...
Found Pokemon:
 - tentacool
 - staryu
 ...
Pokedex > catch tentacool
Throwing a Pokeball at tentacool...
tentacool was caught!
Pokedex > inspect tentacool
Name: tentacool
Height: 9
Weight: 455
...
```

## Features

- **Tab-completion.** Press `Tab` to complete arguments — `explore` suggests
  location areas you've seen via `map`, `catch` suggests Pokémon you've found via
  `explore`, and `inspect` suggests Pokémon you've actually caught. Suggestions
  are built live from what you've discovered, not the full PokéAPI dataset, so
  you're prompted to explore first.
- **Saved progress.** Your Pokédex and the location/Pokémon names you've
  discovered are saved to disk after every change and reloaded on startup, so
  your progress (and tab-completion suggestions) survive between runs.
- **Command history.** Arrow keys scroll through previously entered commands
  (via [chzyer/readline](https://github.com/chzyer/readline)).
- **Response caching.** HTTP responses are cached in memory for 30 minutes, with
  a concurrency-safe store and a background reaper, so repeated `map`/`explore`
  calls don't re-hit the network.

## Project layout

```
.
├── main.go                     # REPL loop, readline setup, tab-completion wiring
├── commands.go                 # command definitions and their handlers
├── persistence.go              # saving/loading state as JSON
├── repl.go                     # input cleaning
└── internal/
    ├── pokeapi/                # PokéAPI client (location areas, encounters, pokemon)
    └── pokecache/              # concurrency-safe in-memory cache with TTL reaping
```

State (your Pokédex and the location/Pokémon names you've discovered) is saved to
`pokedex.json` in your user config directory (e.g. `~/.config/pokedexcli/` on
Linux) after every change, and reloaded when you start the program again.

## Development

```bash
go test ./...   # run the tests
go vet ./...     # static checks
gofmt -l .       # list files needing formatting
```
