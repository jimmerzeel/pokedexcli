package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jimmerzeel/pokedexcli/internal/pokecache"
)

func GetPokemon(pokemonName string, cache *pokecache.Cache) (PokemonResponse, error) {
	fullURL := "https://pokeapi.co/api/v2/pokemon/" + pokemonName
	// check if result is in the cache
	if item, present := cache.Get(fullURL); present {
		var pokemonResponse PokemonResponse

		// unmarshal cached data
		if err := json.Unmarshal(item, &pokemonResponse); err != nil {
			return PokemonResponse{}, err
		}

		return pokemonResponse, nil
	}

	// make HTTP GET request
	res, err := http.Get(fullURL)
	if err != nil {
		return PokemonResponse{}, err
	}
	defer res.Body.Close()

	// read the data from the HTTP request
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return PokemonResponse{}, err
	}

	// unmarshal the data into a slice of bytes
	var pokemonResponse PokemonResponse
	if err = json.Unmarshal(data, &pokemonResponse); err != nil {
		return PokemonResponse{}, err
	}

	cache.Add(fullURL, data)

	return pokemonResponse, nil
}
