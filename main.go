package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
	"time"
)

type Config struct {
	Next     *string
	Previous *string
	Cache    pokecache.Cache
	Pokedex  map[string]pokeapi.Pokemon
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{
		Next:     nil,
		Previous: nil,
		Cache:    pokecache.NewCache(3 * time.Minute),
		Pokedex:  make(map[string]pokeapi.Pokemon),
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		command := cleanedInput[0]
		var argument string
		if len(cleanedInput) > 1 {
			argument = cleanedInput[1]
		}
		if foundCommand, exist := getCommands()[command]; !exist {
			fmt.Println("Unknown command")
		} else {
			err := foundCommand.callback(&config, argument)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
