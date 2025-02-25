package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/jimmerzeel/pokedexcli/internal/pokecache"
)

func GetLocationNames(url string, cache *pokecache.Cache) ([]string, string, string, error) {
	// check if result is in the cache
	if item, present := cache.Get(url); present {
		var mainResponse MainLocationResponse

		// unmarshal cached data
		if err := json.Unmarshal(item, &mainResponse); err != nil {
			return []string{}, "", "", err
		}

		// prepare results from cache data
		var locationNames []string
		for _, loc := range mainResponse.Results {
			locationNames = append(locationNames, loc.Name)
		}

		return locationNames, mainResponse.Next, mainResponse.Previous, nil
	}

	// make HTTP GET request
	res, err := http.Get(url)
	if err != nil {
		return []string{}, "", "", err
	}
	defer res.Body.Close()

	// read the data from the HTTP request
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []string{}, "", "", err
	}

	// unmarshal the data into a slice of bytes
	var mainResponse MainLocationResponse
	if err = json.Unmarshal(data, &mainResponse); err != nil {
		return []string{}, "", "", err
	}

	var locationNames []string
	for _, loc := range mainResponse.Results {
		locationNames = append(locationNames, loc.Name)
	}

	cache.Add(url, data)

	return locationNames, mainResponse.Next, mainResponse.Previous, nil
}

func GetPokemonAtLocation(fullURL string, cache *pokecache.Cache) ([]string, error) {
	// check if result is in the cache
	if item, present := cache.Get(fullURL); present {
		var pokemonAtLocationResponse PokemonAtLocationResponse

		// unmarshal cached data
		if err := json.Unmarshal(item, &pokemonAtLocationResponse); err != nil {
			return []string{}, err
		}

		// prepare results from cache data
		var pokemonAtLocation []string
		for _, encounter := range pokemonAtLocationResponse.PokemonEncounters {
			pokemonAtLocation = append(pokemonAtLocation, encounter.Pokemon.Name)
		}

		return pokemonAtLocation, nil
	}

	// make HTTP request
	res, err := http.Get(fullURL)
	if err != nil {
		return []string{}, err
	}

	// read the data from the HTTP request
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []string{}, nil
	}
	defer res.Body.Close()

	// unmarshal the JSON from the response
	var pokemonAtLocationResponse PokemonAtLocationResponse
	if err = json.Unmarshal(data, &pokemonAtLocationResponse); err != nil {
		return []string{}, err
	}

	// get the pokemon that are at the location
	var pokemonAtLocation []string
	for _, encounter := range pokemonAtLocationResponse.PokemonEncounters {
		pokemonAtLocation = append(pokemonAtLocation, encounter.Pokemon.Name)
	}

	cache.Add(fullURL, data)

	return pokemonAtLocation, nil
}
