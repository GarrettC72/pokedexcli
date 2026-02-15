package main

import (
	"strings"
)

func cleanInput(text string) []string {
	splitInput := strings.Fields(text)
	lowerCaseInput := []string{}
	for _, word := range splitInput {
		lowerCaseInput = append(lowerCaseInput, strings.ToLower(word))
	}
	return lowerCaseInput
}
