package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jimmerzeel/pokedexcli/internal/pokeapi"
	"github.com/jimmerzeel/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, ...string) error
}

type config struct {
	next     string
	previous string
}

func startRepl(cfg *config, cache *pokecache.Cache) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		text := cleanInput(scanner.Text())
		if len(text) == 0 {
			continue
		}

		commandName := text[0]
		args := []string{}
		if len(text) > 1 {
			args = text[1:]
		}

		if command, ok := getCommands()[commandName]; ok {
			err := command.callback(cfg, cache, args...)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}

	}
}

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	trimmed := strings.Fields(lowered)
	return trimmed
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
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
			description: "Displays names of the next 20 location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays names of the previous 20 location areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Find Pokemon at a specified location area",
			callback:    commandExplore,
		},
	}
}

func commandExit(cfg *config, cache *pokecache.Cache, args ...string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, cache *pokecache.Cache, args ...string) error {
	fmt.Print("Welcome to the Pokedex!\nUsage:\n\n")

	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	fmt.Println("")
	return nil
}

func commandMap(cfg *config, cache *pokecache.Cache, args ...string) error {
	// if this is the first call, use the base url, otherwise use the one in config
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.next != "" {
		url = cfg.next
	}

	// use pokeapi location-area endpoint to get the location areas
	locations, next, previous, err := pokeapi.GetLocationNames(url, cache)
	if err != nil {
		return err
	}

	// display the names
	for _, name := range locations {
		fmt.Println(name)
	}

	// update the URLs in the config
	cfg.next = next
	cfg.previous = previous

	return nil
}

func commandMapBack(cfg *config, cache *pokecache.Cache, args ...string) error {
	// if this is the first call, use the base url, otherwise use the one in config
	url := "https://pokeapi.co/api/v2/location-area"
	if cfg.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		url = cfg.previous
	}

	// use pokeapi location-area endpoint to get the location areas
	locations, next, previous, err := pokeapi.GetLocationNames(url, cache)
	if err != nil {
		return err
	}

	// display the names
	for _, name := range locations {
		fmt.Println(name)
	}

	// update the URLs in the config
	cfg.next = next
	cfg.previous = previous

	return nil
}

func commandExplore(cfg *config, cache *pokecache.Cache, args ...string) error {
	// make sure a location area is provided
	if len(args) < 1 {
		fmt.Println("explore command needs a location name area")
		return nil
	}

	locationArea := args[0]
	pokemonEncounters, err := pokeapi.GetPokemonAtLocation(locationArea, cache)
	if err != nil {
		fmt.Printf("%s is not a valid location. Use the map or mapb command to find valid location names\n", locationArea)
		return err
	}

	fmt.Printf("Exploring %s\n", locationArea)
	fmt.Println("Found Pokemon:")

	for _, pokemon := range pokemonEncounters {
		fmt.Printf("- %s\n", pokemon)
	}

	return nil
}
