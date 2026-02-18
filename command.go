package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *Config, argument string) error
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
			description: "Displays the next 20 Pokemon location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 Pokemon location areas",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Displays the Pokemon found in the given location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempt to catch a given Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Displays information about the given Pokemon in the Pokedex",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays all Pokemon caught and registered in the Pokedex",
			callback:    commandPokedex,
		},
	}
}

func commandExit(config *Config, argument string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config, argument string) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
	for _, v := range getCommands() {
		fmt.Printf("%s: %s\n", v.name, v.description)
	}
	return nil
}

func commandMap(config *Config, argument string) error {
	var url string
	if config.Next != nil {
		url = *config.Next
	} else if config.Previous == nil {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		return fmt.Errorf("you're on the last page")
	}
	if err := getLocationAreas(url, config); err != nil {
		return err
	}
	return nil
}

func commandMapb(config *Config, argument string) error {
	var url string
	if config.Previous != nil {
		url = *config.Previous
	} else if config.Next == nil {
		url = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	} else {
		return fmt.Errorf("you're on the first page")
	}
	if err := getLocationAreas(url, config); err != nil {
		return err
	}
	return nil
}

func getLocationAreas(url string, config *Config) error {
	var response pokeapi.LocationPageResponse
	if entry, ok := config.Cache.Get(url); ok {
		if err := json.Unmarshal(entry, &response); err != nil {
			return fmt.Errorf("Error unmarshalling response: %w", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error making Get request: %w", err)
		}
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		}
		if err != nil {
			return fmt.Errorf("Error reading response body: %w", err)
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("Error unmarshalling response: %w", err)
		}
		config.Cache.Add(url, body)
	}
	for _, result := range response.Results {
		fmt.Println(result.Name)
	}
	config.Next = response.Next
	config.Previous = response.Previous
	return nil
}

func commandExplore(config *Config, name string) error {
	if name == "" {
		return fmt.Errorf("Missing location area argument\nUsage: explore <area_name>")
	}
	var response pokeapi.LocationResponse
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", name)
	if entry, ok := config.Cache.Get(url); ok {
		if err := json.Unmarshal(entry, &response); err != nil {
			return fmt.Errorf("Error unmarshalling response: %w", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error making Get request: %w", err)
		}
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		}
		if err != nil {
			return fmt.Errorf("Error reading response body: %w", err)
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("Error unmarshalling response: %w", err)
		}
		config.Cache.Add(url, body)
	}
	fmt.Printf("Exploring %s...\nFound Pokemon:\n", name)
	for _, result := range response.PokemonEncounters {
		fmt.Printf(" - %s\n", result.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *Config, name string) error {
	if name == "" {
		return fmt.Errorf("Missing Pokemon name argument\nUsage: catch <pokemon_name>")
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", name)
	var response pokeapi.Pokemon
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", name)
	if entry, ok := config.Cache.Get(url); ok {
		if err := json.Unmarshal(entry, &response); err != nil {
			return fmt.Errorf("Error unmarshalling response: %w", err)
		}
	} else {
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("Error making Get request: %w", err)
		}
		body, err := io.ReadAll(res.Body)
		defer res.Body.Close()
		if res.StatusCode > 299 {
			return fmt.Errorf("Response failed with status code: %d and\nbody: %s", res.StatusCode, body)
		}
		if err != nil {
			return fmt.Errorf("Error reading response body: %w", err)
		}
		if err := json.Unmarshal(body, &response); err != nil {
			return fmt.Errorf("Error unmarshalling response: %w", err)
		}
		config.Cache.Add(url, body)
	}
	target_rate := 30
	catch_attempt := rand.Intn(response.BaseExperience)
	if catch_attempt < target_rate {
		fmt.Printf("%s was caught!\n", name)
		fmt.Println("You may now inspect it with the inspect command.")
		config.Pokedex[name] = response
	} else {
		fmt.Printf("%s escaped!\n", name)
	}
	return nil
}

func commandInspect(config *Config, name string) error {
	if name == "" {
		return fmt.Errorf("Missing Pokemon name argument\nUsage: inspect <pokemon_name>")
	}
	if pokemon, exists := config.Pokedex[name]; !exists {
		return fmt.Errorf("you have not caught that pokemon")
	} else {
		fmt.Printf("Name: %s\n", pokemon.Name)
		fmt.Printf("Height: %v\n", pokemon.Height)
		fmt.Printf("Weight: %v\n", pokemon.Weight)
		fmt.Println("Stats:")
		for _, pokemon_stat := range pokemon.Stats {
			fmt.Printf("  -%s: %v\n", pokemon_stat.Stat.Name, pokemon_stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, pokemon_type := range pokemon.Types {
			fmt.Printf("  - %s\n", pokemon_type.Type.Name)
		}
	}
	return nil
}

func commandPokedex(config *Config, argument string) error {
	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.Pokedex {
		fmt.Printf(" - %s\n", pokemon.Name)
	}
	return nil
}
