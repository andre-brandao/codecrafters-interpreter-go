package main

import (
	"fmt"
)

func FmtError(token Token, message string) error {
	return fmt.Errorf("Error at line %d: %s at '%v'", token.Line, message, token.Lexeme)
}
