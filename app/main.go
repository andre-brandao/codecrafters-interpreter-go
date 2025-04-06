package main

import (
	"fmt"
	"os"
)

func run(line string) {
    fmt.Println(line)
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
	}
}

func runFile(filename string) {
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	if len(fileContents) > 0 {
		lines := string(fileContents)
		for _, line := range lines {
		    
			run(string(line))
		}
	} else {
		fmt.Println("EOF  null -- file is empty")
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
	// fileContents, err := os.ReadFile(filename)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
	// 	os.Exit(1)
	// }

	// if len(fileContents) > 0 {
	// 	panic("Scanner not implemented")
	// } else {
	// 	fmt.Println("EOF  null") // Placeholder, remove this line when implementing the scanner
	// }
}
