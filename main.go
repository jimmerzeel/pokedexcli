package main

import (
	"time"

	"github.com/jimmerzeel/pokedexcli/internal/pokecache"
)

func main() {
	cfg := &config{}
	cache := pokecache.NewCache(7 * time.Second)
	startRepl(cfg, cache)
}
