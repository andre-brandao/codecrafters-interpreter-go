package main

import (
	"fmt"
	"strconv"
)

type TokenType int

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var keywords = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

func (t TokenType) String() string {
	switch t {
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case COMMA:
		return "COMMA"
	case DOT:
		return "DOT"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case SLASH:
		return "SLASH"
	case STAR:
		return "STAR"
	case BANG:
		return "BANG"
	case BANG_EQUAL:
		return "BANG_EQUAL"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case GREATER:
		return "GREATER"
	case GREATER_EQUAL:
		return "GREATER_EQUAL"
	case LESS:
		return "LESS"
	case LESS_EQUAL:
		return "LESS_EQUAL"
	case IDENTIFIER:
		return "IDENTIFIER"
	case STRING:
		return "STRING"
	case NUMBER:
		return "NUMBER"
	case AND:
		return "AND"
	case CLASS:
		return "CLASS"
	case ELSE:
		return "ELSE"
	case FALSE:
		return "FALSE"
	case FUN:
		return "FUN"
	case FOR:
		return "FOR"
	case IF:
		return "IF"
	case NIL:
		return "NIL"
	case OR:
		return "OR"
	case PRINT:
		return "PRINT"
	case RETURN:
		return "RETURN"
	case SUPER:
		return "SUPER"
	case THIS:
		return "THIS"
	case TRUE:
		return "TRUE"
	case VAR:
		return "VAR"
	case WHILE:
		return "WHILE"
	case EOF:
		return "EOF"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", t)
	}
}

type Token struct {
	Type    TokenType
	Lexeme  []rune
	Literal any
	Line    int
}

func NewToken(t TokenType, lexeme []rune, literal any, line int) Token {
	return Token{
		Type:    t,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}
func (t Token) String() string {
	var literalStr string

	if t.Literal == nil {
		literalStr = "null"
	} else {
		switch v := t.Literal.(type) {
		case []rune:
			literalStr = string(v)
		case string:
			literalStr = v
		case int:
			literalStr = strconv.Itoa(v)
		// case float64:
		// 	literalStr = strconv.FormatFloat(v, 'f', -1, 64)
		case float64:
			if v == float64(int(v)) {
				literalStr = fmt.Sprintf("%.1f", v) // Ensures 1234.0 for whole numbers
			} else {
				literalStr = fmt.Sprintf("%g", v) // Keeps the precision for non-whole numbers

			}
		case bool:
			literalStr = strconv.FormatBool(v)
		// case []string:
		//     literalStr = strings.Join(v, "")
		default:
			literalStr = fmt.Sprintf("%v", v)
		}
	}

	return fmt.Sprintf("%s %s %s\n", t.Type, string(t.Lexeme), literalStr)
}
