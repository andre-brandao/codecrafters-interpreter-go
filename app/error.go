package main

import "fmt"

type RuntimeError struct {
	message string
	token   Token
}

func NewRuntimeError(token Token, message string) *RuntimeError {
	return &RuntimeError{
		message: message,
		token:   token,
	}
}

func (e *RuntimeError) Error() string {
	if len(e.token.Lexeme) > 0 {
		return fmt.Sprintf("[line %d] Error at '%s': %s",
			e.token.Line,
			string(e.token.Lexeme),
			e.message)
	}
	return fmt.Sprintf("[line %d] Error: %s", e.token.Line, e.message)
}
