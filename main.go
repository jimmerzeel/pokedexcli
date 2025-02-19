package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	trimmed := strings.Fields(lowered)
	return trimmed
}

func main() {
	// wait for user input using bufio.NewScanner
	scanner := bufio.NewScanner(os.Stdin)

	// start an infinite for loop
	for {
		// use fmt.Print to print the prompt Pokedex > WITHOUT a newline character
		fmt.Print("Pokedex > ")

		// use the scanner's .Scan and .Text methods to get the user's input as a string
		scanner.Scan()
		input := scanner.Text()

		// clean the user's input by triming any leading or trailing whitespace, and converting it to lowercase. (strings.ToLower and strings.Fields)
		lowerInput := strings.ToLower(input)
		text := strings.Fields(lowerInput)

		// capture the first "word" of the input and use it to print: "Your command was: <first word>"
		firstWord := text[0]

		fmt.Printf("Your command was: %s\n", firstWord)
	}

	// first input: CHARMANDER is better than bulbasaur
	// second input: Pikachu is kinda mean to ash
	// terminate the program by pressing ctrl+c
}
