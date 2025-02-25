package main

import (
	"time"

	"github.com/jimmerzeel/pokedexcli/internal/pokeapi"
	"github.com/jimmerzeel/pokedexcli/internal/pokecache"
)

func main() {
	cfg := &config{
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}
	cache := pokecache.NewCache(7 * time.Second)
	startRepl(cfg, cache)
}
