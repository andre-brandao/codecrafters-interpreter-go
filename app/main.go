package main

import (
	"fmt"
	"os"
)

var hadError = false

func report(line int, where string, message string) {
	// fmt.Println("[line", line, "] Error: ", where, message)
	fmt.Printf("[line %d] Error: %s\n", line, message)
	hadError = true
}

func error(line int, message string) {
	report(line, "", message)
}

func run(source string) {
	fmt.Println(source)

	for _, char := range source {
		token := string(char)
		fmt.Println("Symbol: ", token)
	}
}

func runPrompt() {
	input := ""
	for {
		fmt.Print(">> ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		if input == "exit" {
			break
		}
		run(input)
		hadError = false
	}
}

func runFile(filename string) {
	fileContents, err := os.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	if len(fileContents) > 0 {

		runes := []rune(string(fileContents))
		s := NewScanner(runes)

		tokens := s.ScanTokens()

		for _, token := range tokens {
			fmt.Printf(token.String())
		}
		
		if hadError {
		    os.Exit(65)
		}
	} else {
		fmt.Println("EOF  null")
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command == "repl" {
		runPrompt()
		return
	}

	if command != "tokenize" {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	filename := os.Args[2]

	runFile(filename)

}
