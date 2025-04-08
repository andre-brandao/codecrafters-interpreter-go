package expr

import (
	token "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type ExprVisitor interface {
	VisitBinaryExpr(expr *Binary) any
	VisitGroupingExpr(expr *Grouping) any
	VisitLiteralExpr(expr *Literal) any
	VisitUnaryExpr(expr *Unary) any
	// VisitVarExpr(expr *Var) any
	VisitCallExpr(expr *Call) any
	VisitVariableExpr(expr *Variable) any
	VisitAssignExpr(expr *Assign) any
	VisitLogicalExpr(expr *Logical) any
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

// type Var struct {
// 	Name        token.Token
// 	Initializer Expr
// }

// func (v *Var) Accept(visitor ExprVisitor) interface{} {
// 	return visitor.VisitVarExpr(v)
// }

// var _ Expr = (*Var)(nil)

type Variable struct {
	Name token.Token
}

func (v *Variable) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitVariableExpr(v)
}

type Assign struct {
	Name  token.Token
	Value Expr
}

func (a *Assign) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitAssignExpr(a)
}

var _ Expr = (*Assign)(nil)

type Logical struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

func (l *Logical) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitLogicalExpr(l)
}

var _ Expr = (*Logical)(nil)

type Call struct {
	Callee    Expr
	Paren     token.Token
	Arguments []Expr
}

func (c *Call) Accept(visitor ExprVisitor) interface{} {
	return visitor.VisitCallExpr(c)
}

var _ Expr = (*Call)(nil)
