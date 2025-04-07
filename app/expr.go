package main

// type Operator TokenType

// const (
// 	OP_PLUS  Operator = Operator(PLUS)
// 	OP_MINUS Operator = Operator(MINUS)
// 	OP_STAR  Operator = Operator(STAR)
// 	OP_SLASH Operator = Operator(SLASH)

// 	OP_EQUAL_EQUAL   Operator = Operator(EQUAL_EQUAL)
// 	OP_BANG_EQUAL    Operator = Operator(BANG_EQUAL)
// 	OP_GREATER       Operator = Operator(GREATER)
// 	OP_GREATER_EQUAL Operator = Operator(GREATER_EQUAL)
// 	OP_LESS          Operator = Operator(LESS)
// 	OP_LESS_EQUAL    Operator = Operator(LESS_EQUAL)

// 	OP_AND Operator = Operator(AND)
// 	OP_OR  Operator = Operator(OR)
// )

// func (op Operator) String() string {
// 	return TokenType(op).String()
// }

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

// EXPR
type Expr interface {
	Accept(visitor ExprVisitor) interface{}
}

// BINARY EXPR
type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

var _ Expr = (*Binary)(nil)

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{
		Left:     left,
		Operator: operator,
		Right:    right,
	}
}

// LITERAL EXPR
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

var _ Expr = (*Literal)(nil)

func NewLiteral(value interface{}) *Literal {
	return &Literal{Value: value}
}

// UNARY EXPR

type Unary struct {
	Operator Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

var _ Expr = (*Unary)(nil)

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{
		Operator: operator,
		Right:    right,
	}
}

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

var _ Expr = (*Grouping)(nil)

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{Expression: expression}
}
