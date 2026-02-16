package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config) error
}

type Config struct {
	Next     *string
	Previous *string
	Cache    pokecache.Cache
}

var cliCommandMap map[string]cliCommand

func init() {
	cliCommandMap = map[string]cliCommand{
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
			description: "Displays the next 20 Pokemon location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 Pokemon location areas",
			callback:    commandMapb,
		},
	}
}

func commandExit(config *Config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Error from exiting Pokedex CLI")
}

func commandHelp(config *Config) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n")
	for _, v := range cliCommandMap {
		fmt.Printf("\n%s: %s", v.name, v.description)
	}
	return fmt.Errorf("\n")
}

func commandMap(config *Config) error {
	var url string
	if config.Next != nil {
		url = *config.Next
	} else if config.Previous == nil {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		return fmt.Errorf("you're on the last page")
	}
	if entry, ok := config.Cache.Get(url); ok {
		return readLocationAreas(entry, config)
	}
	return getLocationAreas(url, config)
}

func commandMapb(config *Config) error {
	var url string
	if config.Previous != nil {
		url = *config.Previous
	} else if config.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		return fmt.Errorf("you're on the first page\n")
	}
	if entry, ok := config.Cache.Get(url); ok {
		return readLocationAreas(entry, config)
	}
	return getLocationAreas(url, config)
}

func getLocationAreas(url string, config *Config) error {
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error making Get request: %w", err)
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if res.StatusCode > 299 {
		return fmt.Errorf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}
	var response pokeapi.PokeAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
	config.Next = response.Next
	config.Previous = response.Previous
	config.Cache.Add(url, body)
	return fmt.Errorf("")
}

func readLocationAreas(body []byte, config *Config) error {
	var response pokeapi.PokeAPIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("error unmarshalling response: %w", err)
	}
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
	config.Next = response.Next
	config.Previous = response.Previous
	return fmt.Errorf("")
}
