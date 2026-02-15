package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := scanner.Text()
		cleanedInput := cleanInput(input)
		command := cleanedInput[0]
		_, exist := commandMap[command]
		if !exist {
			fmt.Println("Unknown command")
		} else {
			err := commandMap[command].callback()
			fmt.Println(err)
		}
	}
}
