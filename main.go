package main

import (
	"bufio"
	"fmt"
	"os"
	"pokedexcli/internal/pokecache"
	"time"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	config := Config{
		Next:     nil,
		Previous: nil,
		Cache:    pokecache.NewCache(3 * time.Minute),
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		command := cleanedInput[0]
		_, exist := cliCommandMap[command]
		if !exist {
			fmt.Println("Unknown command")
		} else {
			err := cliCommandMap[command].callback(&config)
			fmt.Print(err)
		}
	}
}
