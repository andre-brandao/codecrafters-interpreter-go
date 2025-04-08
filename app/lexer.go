package main

import (
	"fmt"
	"strconv"

	tok "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Scanner struct {
	source  []rune
	tokens  []tok.Token
	start   int
	current int
	line    int
}

func NewScanner(source []rune) *Scanner {
	return &Scanner{
		source:  source,
		tokens:  make([]tok.Token, 0),
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) ScanTokens() []tok.Token {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, tok.Token{
		Type:    tok.EOF,
		Lexeme:  []rune(""),
		Literal: nil,
		Line:    s.line,
	})

	return s.tokens
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() || (s.source[s.current] != expected) {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		// error(s.line, "Unterminated string.")
		report(s.line, "", "Unterminated string.")
		return
	}

	s.advance()
	value := s.source[s.start+1 : s.current-1]
	s.addToken(tok.STRING, value)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	numStr := string(s.source[s.start:s.current])
	val, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		// error(s.line, "Invalid number format.")
		report(s.line, "", "Invalid number format.")
		return
	}
	s.addToken(tok.NUMBER, val)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := string(s.source[s.start:s.current])
	type_, ok := tok.Keywords[text]

	if !ok {
		type_ = tok.IDENTIFIER
	}
	s.addToken(type_, nil)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(tok.LEFT_PAREN, nil)
	case ')':
		s.addToken(tok.RIGHT_PAREN, nil)
	case '{':
		s.addToken(tok.LEFT_BRACE, nil)
	case '}':
		s.addToken(tok.RIGHT_BRACE, nil)
	case ',':
		s.addToken(tok.COMMA, nil)
	case '.':
		s.addToken(tok.DOT, nil)
	case '-':
		s.addToken(tok.MINUS, nil)
	case '+':
		s.addToken(tok.PLUS, nil)
	case ';':
		s.addToken(tok.SEMICOLON, nil)
	case '*':
		s.addToken(tok.STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(tok.BANG_EQUAL, nil)
		} else {
			s.addToken(tok.BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(tok.EQUAL_EQUAL, nil)
		} else {
			s.addToken(tok.EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(tok.LESS_EQUAL, nil)
		} else {
			s.addToken(tok.LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(tok.GREATER_EQUAL, nil)
		} else {
			s.addToken(tok.GREATER, nil)
		}
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			// while (peek() != '\n' && !isAtEnd()) advance();
			//
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tok.SLASH, nil)
		}
	case ' ', '\r', '\t':
		// Ignore whitespace.
		break
	case '\n':
		s.line++
	case '"':
		s.string()

	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			// error(s.line, fmt.Sprintf("Unexpected character: %c", c))
			report(s.line, "", fmt.Sprintf("Unexpected character: %c", c))
		}
	}
}

func (s *Scanner) addToken(t tok.TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, tok.Token{
		Type:    t,
		Lexeme:  text,
		Literal: literal,
		Line:    s.line,
	})
}
