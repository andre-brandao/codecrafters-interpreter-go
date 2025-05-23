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
	// return p.equality()
	return p.assignment()
}

func (p *Parser) declaration() st.Stmt {
	defer func() {
		if r := recover(); r != nil {
			p.synchronize()
		}
	}()

	if p.match(tok.FUN) {
		return p.function("function")
	}

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

	if p.match(tok.FOR) {
		return p.forStatement()
	}

	if p.match(tok.IF) {
		return p.ifStatement()
	}

	if p.match(tok.PRINT) {
		return p.printStatement()
	}

	if p.match(tok.RETURN) {
		return p.returnStatement()
	}

	if p.match(tok.WHILE) {
		return p.whileStatement()
	}

	if p.match(tok.LEFT_BRACE) {
		return &st.Block{
			Statements: p.block(),
		}
	}

	return p.expressionStatement()
}

func (p *Parser) forStatement() st.Stmt {
	p.consume(tok.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer st.Stmt = nil
	if p.match(tok.SEMICOLON) {
		initializer = nil
	} else if p.match(tok.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition exp.Expr = nil
	if !p.check(tok.SEMICOLON) {
		condition = p.expression()
	}
	p.consume(tok.SEMICOLON, "Expect ';' after loop condition.")

	var increment exp.Expr = nil
	if !p.check(tok.RIGHT_PAREN) {
		increment = p.expression()
	}

	p.consume(tok.RIGHT_PAREN, "Expect ')' after for clauses.")

	body := p.statement()

	if increment != nil {
		body = &st.Block{
			Statements: []st.Stmt{
				body,
				&st.Expression{
					Expression: increment,
				},
			},
		}
	}

	if condition == nil {
		condition = &exp.Literal{
			Value: true,
		}
	}

	body = &st.While{
		Condition: condition,
		Body:      body,
	}

	if initializer != nil {
		body = &st.Block{
			Statements: []st.Stmt{
				initializer, body,
			},
		}
	}

	return body
}

func (p *Parser) ifStatement() st.Stmt {
	p.consume(tok.LEFT_PAREN, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(tok.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch st.Stmt = nil
	if p.match(tok.ELSE) {
		elseBranch = p.statement()
	}

	return &st.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() st.Stmt {
	value := p.expression()
	p.consume(tok.SEMICOLON, "Expect ';' after value.")
	return &st.Print{
		Expression: value,
	}
}

func (p *Parser) returnStatement() st.Stmt {
	keyword := p.previous()
	var value exp.Expr = nil

	if !p.check(tok.SEMICOLON) {
		value = p.expression()
	}
	p.consume(tok.SEMICOLON, "Expect ';' after return value.")

	return &st.Return{
		Keyword: keyword,
		Value:   value,
	}
}

func (p *Parser) whileStatement() st.Stmt {
	p.consume(tok.LEFT_PAREN, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(tok.RIGHT_PAREN, "Expect ')' after condition.")

	body := p.statement()

	return &st.While{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) expressionStatement() st.Stmt {
	expr := p.expression()
	p.consume(tok.SEMICOLON, "Expect ';' after expression.")
	return &st.Expression{
		Expression: expr,
	}
}

func (p *Parser) function(kind string) *st.Function {
	name := p.consume(tok.IDENTIFIER, fmt.Sprintf("Expected %s name.", kind))

	p.consume(tok.LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	parameters := make([]tok.Token, 0)
	if !p.check(tok.RIGHT_PAREN) {
		for {
			if len(parameters) >= 255 {
				p.Error(p.peek(), "Can't have more than 255 parameters.")
			}
			parameters = append(parameters, p.consume(tok.IDENTIFIER, "Expect parameter name."))
			// if
			if !p.match(tok.COMMA) {
				break
			}
		}
	}

	p.consume(tok.RIGHT_PAREN, "Expect ')' after parameters.")

	p.consume(tok.LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))

	body := p.block()

	return &st.Function{
		Name:   name,
		Params: parameters,
		Body:   body,
	}
}

func (p *Parser) block() []st.Stmt {
	statements := make([]st.Stmt, 0)

	for !p.check(tok.RIGHT_BRACE) && !p.isAtEnd() {
		// fmt.Print(statements)
		statements = append(statements, p.declaration())
	}
	p.consume(tok.RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p *Parser) assignment() exp.Expr {
	// expr := p.equality()
	expr := p.or()

	if p.match(tok.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if variable, ok := expr.(*exp.Variable); ok {
			name := variable.Name
			return &exp.Assign{
				Name:  name,
				Value: value,
			}

		}

		p.Error(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) or() exp.Expr {
	expr := p.and()
	for p.match(tok.OR) {
		operator := p.previous()
		right := p.and()

		expr = &exp.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) and() exp.Expr {
	expr := p.equality()
	for p.match(tok.AND) {
		operator := p.previous()
		right := p.equality()

		expr = &exp.Logical{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) equality() exp.Expr {
	expr := p.comparison()

	for p.match(tok.BANG_EQUAL, tok.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()

		expr = &exp.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}

	return expr
}

func (p *Parser) comparison() exp.Expr {
	expr := p.term()

	for p.match(tok.GREATER, tok.GREATER_EQUAL, tok.LESS, tok.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()

		// expr = exp.NewBinary(expr, operator, right)
		expr = &exp.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) term() exp.Expr {
	expr := p.factor()

	for p.match(tok.MINUS, tok.PLUS) {
		operator := p.previous()
		right := p.factor()

		// expr = exp.NewBinary(expr, operator, right)
		expr = &exp.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) factor() exp.Expr {
	expr := p.unary()

	for p.match(tok.SLASH, tok.STAR) {
		operator := p.previous()
		right := p.unary()

		// expr = exp.NewBinary(expr, operator, right)
		expr = &exp.Binary{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) unary() exp.Expr {
	if p.match(tok.BANG, tok.MINUS) {
		operator := p.previous()
		right := p.unary()

		// return exp.NewUnary(operator, right)
		return &exp.Unary{
			Operator: operator,
			Right:    right,
		}
	}
	// return p.primary()
	return p.call()
}

func (p *Parser) finishCall(callee exp.Expr) exp.Expr {
	arguments := make([]exp.Expr, 0)

	if !p.check(tok.RIGHT_PAREN) {
		for {
			if len(arguments) >= 255 {
				p.Error(p.peek(), "Can't have more than 255 arguments.")
			}

			arguments = append(arguments, p.expression())
			if !p.match(tok.COMMA) {
				break
			}
		}

	}

	paren := p.consume(tok.RIGHT_PAREN, "Expect ')' after arguments.")

	return &exp.Call{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (p *Parser) call() exp.Expr {
	expr := p.primary()
	for true {
		if p.match(tok.LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr
}

func (p *Parser) primary() exp.Expr {
	if p.match(tok.FALSE) {
		// return exp.NewLiteral(false)
		return &exp.Literal{
			Value: false,
		}
	}
	if p.match(tok.TRUE) {
		// return exp.NewLiteral(true)
		return &exp.Literal{
			Value: true,
		}
	}
	if p.match(tok.NIL) {
		// return exp.NewLiteral(nil)
		return &exp.Literal{
			Value: nil}
	}

	if p.match(tok.NUMBER, tok.STRING) {
		// return exp.NewLiteral(p.previous().Literal)
		return &exp.Literal{
			Value: p.previous().Literal,
		}
	}

	if p.match(tok.IDENTIFIER) {
		return &exp.Variable{
			Name: p.previous(),
		}
	}

	if p.match(tok.LEFT_PAREN) {
		expr := p.expression()

		p.consume(tok.RIGHT_PAREN, "Expect ')' after expression.")
		// return exp.NewGrouping(expr)
		return &exp.Grouping{
			Expression: expr,
		}
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
	// panic(message)
	// panic(err.NewRuntimeError(, message))
	panic(p.Error(tok.NewToken(t, p.peek().Lexeme, nil, p.peek().Line), message))
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
