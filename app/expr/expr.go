package expr

import (
	token "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

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
	Operator token.Token
	Right    Expr
}

func (b *Binary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitBinaryExpr(b)
}

var _ Expr = (*Binary)(nil)

func NewBinary(left Expr, operator token.Token, right Expr) *Binary {
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
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

var _ Expr = (*Unary)(nil)

func NewUnary(operator token.Token, right Expr) *Unary {
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
