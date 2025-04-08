package stmt

import (
	expr "github.com/codecrafters-io/interpreter-starter-go/app/expr"
	token "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type StmtVisitor interface {
	VisitPrintStmt(stmt *Print) any
	VisitExpressionStmt(stmt *Expression) any
	VisitVarStmt(stmt *Var) any
	// VisitVariableStmt(stmt *Variable) any
	// VisitUnaryStmt(stmt *UnaryStmt) any
	VisitBlockStmt(stmt *Block) any
	VisitIfStmt(stmt *If) any
	VisitWhileStmt(stmt *While) any
	// VisitForStmt(stmt *For) any
	// VisitClassStmt(stmt *Class) any
	VisitFunctionStmt(stmt *Function) any
	VisitReturnStmt(stmt *Return) interface{}
}

type Stmt interface {
	Accept(visitor StmtVisitor) any
}

type Expression struct {
	Expression expr.Expr
}

func (e *Expression) Accept(visitor StmtVisitor) any {
	return visitor.VisitExpressionStmt(e)
}

var _ Stmt = &Expression{}

type Print struct {
	Expression expr.Expr
}

func (p *Print) Accept(visitor StmtVisitor) any {
	return visitor.VisitPrintStmt(p)
}

var _ Stmt = &Print{}

type Var struct {
	Name        token.Token
	Initializer expr.Expr
}

func (v *Var) Accept(visitor StmtVisitor) any {
	return visitor.VisitVarStmt(v)
}

var _ Stmt = &Var{}

// type UnaryStmt struct {
// 	operator token.Token
// 	right    expr.Expr
// }

// func (u *UnaryStmt) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitUnaryStmt(u)
// }

// var _ Stmt = &UnaryStmt{}

// type Variable struct {
// 	Name token.Token
// }

// func (v *Variable) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitVariableStmt(v)
// }

// var _ Stmt = &Variable{}

type Block struct {
	Statements []Stmt
}

func (b *Block) Accept(visitor StmtVisitor) any {
	return visitor.VisitBlockStmt(b)
}

var _ Stmt = &Block{}

type If struct {
	Condition  expr.Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (i *If) Accept(visitor StmtVisitor) any {
	return visitor.VisitIfStmt(i)
}

var _ Stmt = &If{}

// type Class struct {
// 	Name       Token
// 	SuperClass *Var
// 	Methods    []*Function
// }

// func (c *Class) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitClassStmt(c)
// }

// var _ Stmt = &Class{}

type While struct {
	Condition expr.Expr
	Body      Stmt
}

func (w *While) Accept(visitor StmtVisitor) any {
	return visitor.VisitWhileStmt(w)
}

var _ Stmt = &While{}

type Function struct {
	Name   token.Token
	Params []token.Token
	Body   []Stmt
}

func (f *Function) Accept(visitor StmtVisitor) any {
	return visitor.VisitFunctionStmt(f)
}

var _ Stmt = &Function{}

type Return struct {
	Keyword token.Token
	Value   expr.Expr
}

func (r *Return) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitReturnStmt(r)
}

var _ Stmt = &Return{}
