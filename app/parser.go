package main

import (
	"fmt"
)

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}
func (p *Parser) ParseExpression() Expr {
	return p.expression()
}

func (p *Parser) Parse() []Stmt {
	// create an arraylist of statements
	statements := make([]Stmt, 0)

	for !p.isAtEnd() {
		// statements = append(statements, p.statement())
		statements = append(statements, p.declaration())

	}
	// return p.expression()
	return statements
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p *Parser) expression() Expr {
	return p.equality()
}

func (p *Parser) declaration() Stmt {
	defer func() {
		if r := recover(); r != nil {
			p.synchronize()

		}
	}()

	if p.match(VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name.")

	var initializer Expr = nil
	if p.match(EQUAL) {
		initializer = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration.")

	return &Var{
		Name:        name,
		Initializer: initializer,
	}
}
func (p *Parser) statement() Stmt {
	if p.match(PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return &Print{
		Expression: value,
	}
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return &Expression{
		Expression: expr,
	}
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()

		expr = NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()
		right := p.term()

		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		operator := p.previous()
		right := p.factor()

		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		operator := p.previous()
		right := p.unary()

		expr = NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right := p.unary()

		return NewUnary(operator, right)
	}
	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return NewLiteral(false)
	}
	if p.match(TRUE) {
		return NewLiteral(true)
	}
	if p.match(NIL) {
		return NewLiteral(nil)
	}

	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().Literal)
	}

	// if p.match(IDENTIFIER) {
	// 	return NewVariable(p.previous())
	// }

	if p.match(LEFT_PAREN) {
		expr := p.expression()

		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	}
	// return nil, fmt.Errorf("Expect expression.", p.peek().Line)
	// return nil, p.Error(p.peek(), "Expect expression.")
	panic(p.Error(p.peek(), "Expect expression."))
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t TokenType, message string) Token {
	if p.check(t) {
		return p.advance()
	}
	panic(message)
}

func (p *Parser) check(t TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}
func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) Error(token Token, message string) (err error) {
	Error(token, message)

	return fmt.Errorf("[line %d] Error at '%v': %s", token.Line, token.Lexeme, message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}
		p.advance()
	}
}
