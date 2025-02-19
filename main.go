package main

import (
	"fmt"
	"strings"
)

func cleanInput(text string) []string {
	lowered := strings.ToLower(text)
	trimmed := strings.Fields(lowered)
	return trimmed
}

func main() {
	fmt.Println("Hello, World!")
}
