package internal

import (
	"encoding/json"
	"io"
	"net/http"
)

// match the API structure
type MainLocationResponse struct {
	Count    int                    `json:"count"`
	Next     string                 `json:"next"`
	Previous string                 `json:"previous"`
	Results  []LocationAreaResponse `json:"results"`
}

type LocationAreaResponse struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func getLocationNames(url string) ([]string, string, string, error) {
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

	return locationNames, mainResponse.Next, mainResponse.Previous, nil
}
