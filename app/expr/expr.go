package expr

import (
	token "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
	VisitVarExpr(expr *Var) interface{}
	VisitVariableExpr(expr *Variable) interface{}
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

// LITERAL EXPR
type Literal struct {
	Value interface{}
}

func (l *Literal) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLiteralExpr(l)
}

var _ Expr = (*Literal)(nil)

// UNARY EXPR

type Unary struct {
	Operator token.Token
	Right    Expr
}

func (u *Unary) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitUnaryExpr(u)
}

var _ Expr = (*Unary)(nil)

type Grouping struct {
	Expression Expr
}

func (g *Grouping) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitGroupingExpr(g)
}

var _ Expr = (*Grouping)(nil)

type Var struct {
	Name        token.Token
	Initializer Expr
}

func (v *Var) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVarExpr(v)
}

var _ Expr = (*Var)(nil)

type Variable struct {
	Name token.Token
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}
