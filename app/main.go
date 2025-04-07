package main

import (
	"fmt"
	"os"
)

var hadError = false

func report(line int, where string, message string) {
	// fmt.Printf("[line %d] Error: %s\n", line, message)
	fmt.Fprintf(os.Stderr, "[line %d] Error: %s\n", line, message)

	hadError = true
}

func Error(token Token, message string) {
	// report(line, "", message)
	if token.Type == EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, fmt.Sprintf(" at '%v'", token.Lexeme), message)
	}
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

type LoxHandler func([]rune)

func runFile(filename string, handler LoxHandler) {
	fileContents, err := os.ReadFile(filename)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	if len(fileContents) > 0 {

		sourceCode := []rune(string(fileContents))

		handler(sourceCode)

		if hadError {
			os.Exit(65)
		}

		os.Exit(0)
	} else {
		fmt.Println("<|EMPTY FILE|>")
		os.Exit(0)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]
	filename := os.Args[2]

	if command == "repl" {
		runPrompt()
		return
	}

	if command == "tokenize" {
		runFile(filename, func(source []rune) {
			s := NewScanner(source)

			tokens := s.ScanTokens()

			for _, token := range tokens {
				fmt.Printf(token.String())
			}
		})
	}

	if command == "parse" {
		runFile(filename, func(source []rune) {
			s := NewScanner(source)

			tokens := s.ScanTokens()

			p := NewParser(tokens)

			expr, err := p.Parse()
			// fmt.Print(expr)

			if err != nil {
				fmt.Println("Error parsing:", err)
				return
			}

			printer := NewAstPrinter()

			fmt.Println(printer.Print(expr))

		})
		os.Exit(0)
		return
	}

	if command == "ast" {
		astPrinter := NewAstPrinter()
		exp1 := NewUnary(NewToken(MINUS, []rune{'-'}, nil, 1), NewLiteral(123))
		exp2 := NewGrouping(NewLiteral(45.67))

		exp3 := NewBinary(exp1, NewToken(STAR, []rune{'*'}, nil, 1), exp2)
		fmt.Println(astPrinter.Print(exp3))

	}

	//fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
	os.Exit(1)

}
