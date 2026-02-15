package main

import (
	"fmt"
	"os"
)

var commandMap map[string]cliCommand

func init() {
	commandMap = map[string]cliCommand{
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
	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return fmt.Errorf("Error from exiting Pokedex CLI")
}

func commandHelp() error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n")
	for _, v := range commandMap {
		fmt.Printf("\n%s: %s", v.name, v.description)
	}
	return fmt.Errorf("")
}

type cliCommand struct {
	name        string
	description string
	callback    func() error
}
