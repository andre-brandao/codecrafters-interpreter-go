package main

import (
	"fmt"
	"os"

	printer "github.com/codecrafters-io/interpreter-starter-go/app/printer"
	tok "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

var hadError = false

var hadRuntimeError = false

func report(line int, where string, message string) {
	// fmt.Printf("[line %d] Error: %s\n", line, message)
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %s\n", line, where, message)

	hadError = true
}

func Error(token tok.Token, message string) {
	// report(line, "", message)
	if token.Type == tok.EOF {
		report(token.Line, " at end", message)
	} else {
		report(token.Line, fmt.Sprintf(" at '%s'", string(token.Lexeme)), message)
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
			return
		}
		if hadRuntimeError {
			os.Exit(70)
			return
		}

		os.Exit(0)
	} else {
		fmt.Println("EOF  null")
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

	switch command {
	case "repl":
		runPrompt()
	case "tokenize":
		runFile(filename, func(source []rune) {
			s := NewScanner(source)

			tokens := s.ScanTokens()

			for _, token := range tokens {
				fmt.Printf(token.String())
			}
		})

	case "parse":
		runFile(filename, func(source []rune) {
			s := NewScanner(source)

			tokens := s.ScanTokens()

			p := NewParser(tokens)

			defer func() {
				if r := recover(); r != nil {
					// fmt.Println("recovered")
					// fmt.Fprint(os.Stderr, r)
					hadError = true

				}
			}()
			expr := p.ParseExpression()

			if hadError {
				return
			}

			printer := printer.NewAstPrinter()

			fmt.Println(printer.Print(expr))

		})
		os.Exit(0)
		return
	case "evaluate", "eval":
		runFile(filename, func(source []rune) {
			s := NewScanner(source)
			tokens := s.ScanTokens()

			p := NewParser(tokens)
			defer func() {
				if r := recover(); r != nil {
					// fmt.Println("recovered")
					// fmt.Fprint(os.Stderr, r)
					hadError = true

				}
			}()
			expr := p.ParseExpression()

			if hadError {
				return
			}
			interpreter := NewInterpreter()
			interpreter.InterpretExpression(expr)
			if hadRuntimeError {
				return
			}
		})
		return

	case "run":

		runFile(filename, func(source []rune) {
			s := NewScanner(source)
			tokens := s.ScanTokens()

			p := NewParser(tokens)
			statements := p.Parse()

			if hadError {
				return
			}

			interpreter := NewInterpreter()
			resolver := NewResolver(*interpreter)

			if hadError {
				return
			}
			resolver.resolveStmts(statements)
			interpreter.Interpret(statements)
			if hadRuntimeError {
				return
			}
		})
		return

	default:
		fmt.Fprintln(os.Stderr, "Unknown command:", command)

	}

	os.Exit(1)

}
