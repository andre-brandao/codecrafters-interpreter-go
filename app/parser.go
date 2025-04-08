package main

import (
	"fmt"

	exp "github.com/codecrafters-io/interpreter-starter-go/app/expr"
	st "github.com/codecrafters-io/interpreter-starter-go/app/stmt"
	tok "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Parser struct {
	tokens  []tok.Token
	current int
}

func NewParser(tokens []tok.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}
func (p *Parser) ParseExpression() exp.Expr {
	return p.expression()
}

func (p *Parser) Parse() []st.Stmt {
	// create an arraylist of statements
	statements := make([]st.Stmt, 0)

	for !p.isAtEnd() {
		// statements = append(statements, p.statement())
		statements = append(statements, p.declaration())

	}
	// return p.expression()
	return statements
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tok.EOF
}

func (p *Parser) expression() exp.Expr {
	return p.equality()
}

func (p *Parser) declaration() st.Stmt {
	defer func() {
		if r := recover(); r != nil {
			p.synchronize()
		}
	}()

	if p.match(tok.VAR) {
		return p.varDeclaration()
	}

	return p.statement()
}

func (p *Parser) varDeclaration() st.Stmt {
	name := p.consume(tok.IDENTIFIER, "Expect variable name.")

	var initializer exp.Expr = nil
	if p.match(tok.EQUAL) {
		initializer = p.expression()
	}
	p.consume(tok.SEMICOLON, "Expect ';' after variable declaration.")

	return &st.Var{
		Name:        name,
		Initializer: initializer,
	}
}
func (p *Parser) statement() st.Stmt {
	if p.match(tok.PRINT) {
		return p.printStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) printStatement() st.Stmt {
	value := p.expression()
	p.consume(tok.SEMICOLON, "Expect ';' after value.")
	return &st.Print{
		Expression: value,
	}
}

func (p *Parser) expressionStatement() st.Stmt {
	expr := p.expression()
	p.consume(tok.SEMICOLON, "Expect ';' after expression.")
	return &st.Expression{
		Expression: expr,
	}
}

func (p *Parser) equality() exp.Expr {
	expr := p.comparison()

	for p.match(tok.BANG_EQUAL, tok.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()

		expr = exp.NewBinary(expr, operator, right)
	}

	return expr
}

func (p *Parser) comparison() exp.Expr {
	expr := p.term()

	for p.match(tok.GREATER, tok.GREATER_EQUAL, tok.LESS, tok.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()

		expr = exp.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) term() exp.Expr {
	expr := p.factor()

	for p.match(tok.MINUS, tok.PLUS) {
		operator := p.previous()
		right := p.factor()

		expr = exp.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) factor() exp.Expr {
	expr := p.unary()

	for p.match(tok.SLASH, tok.STAR) {
		operator := p.previous()
		right := p.unary()

		expr = exp.NewBinary(expr, operator, right)
	}
	return expr
}

func (p *Parser) unary() exp.Expr {
	if p.match(tok.BANG, tok.MINUS) {
		operator := p.previous()
		right := p.unary()

		return exp.NewUnary(operator, right)
	}
	return p.primary()
}

func (p *Parser) primary() exp.Expr {
	if p.match(tok.FALSE) {
		return exp.NewLiteral(false)
	}
	if p.match(tok.TRUE) {
		return exp.NewLiteral(true)
	}
	if p.match(tok.NIL) {
		return exp.NewLiteral(nil)
	}

	if p.match(tok.NUMBER, tok.STRING) {
		return exp.NewLiteral(p.previous().Literal)
	}

	// if p.match(tok.IDENTIFIER) {
	// 	return exp.NewVariable(p.previous())
	// }

	if p.match(tok.LEFT_PAREN) {
		expr := p.expression()

		p.consume(tok.RIGHT_PAREN, "Expect ')' after expression.")
		return exp.NewGrouping(expr)
	}
	// return nil, fmt.Errorf("Expect expression.", p.peek().Line)
	// return nil, p.Error(p.peek(), "Expect expression.")
	panic(p.Error(p.peek(), "Expect expression."))
}

func (p *Parser) match(types ...tok.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) consume(t tok.TokenType, message string) tok.Token {
	if p.check(t) {
		return p.advance()
	}
	panic(message)
}

func (p *Parser) check(t tok.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}
func (p *Parser) advance() tok.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() tok.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) peek() tok.Token {
	return p.tokens[p.current]
}

func (p *Parser) Error(token tok.Token, message string) (err error) {
	Error(token, message)

	return fmt.Errorf("[line %d] Error at '%s': %s", token.Line, string(token.Lexeme), message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == tok.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case tok.CLASS, tok.FUN, tok.VAR, tok.FOR, tok.IF, tok.WHILE, tok.PRINT, tok.RETURN:
			return
		}
		p.advance()
	}
}
