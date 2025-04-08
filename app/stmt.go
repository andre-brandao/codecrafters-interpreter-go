package main

type StmtVisitor interface {
	VisitPrintStmt(stmt *Print) interface{}
	VisitExpressionStmt(stmt *Expression) interface{}
	// VisitBlockStmt(stmt *Block) interface{}
	// VisitClassStmt(stmt *Class) interface{}
	// VisitFunctionStmt(stmt *Function) interface{}
	// VisitIfStmt(stmt *If) interface{}
	// VisitReturnStmt(stmt *Return) interface{}
	// VisitVarStmt(stmt *Var) interface{}
	// VisitWhileStmt(stmt *While) interface{}
}

type Stmt interface {
	Accept(visitor StmtVisitor) interface{}
}

type Expression struct {
	Expression Expr
}

func (e *Expression) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitExpressionStmt(e)
}

var _ Stmt = &Expression{}

type Print struct {
	Expression Expr
}

func (p *Print) Accept(visitor StmtVisitor) interface{} {
	return visitor.VisitPrintStmt(p)
}

var _ Stmt = &Print{}

// type Block struct {
// 	Statements []Stmt
// }

// func (b *Block) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitBlockStmt(b)
// }

// var _ Stmt = &Block{}

// type Class struct {
// 	Name       Token
// 	SuperClass *Var
// 	Methods    []*Function
// }

// func (c *Class) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitClassStmt(c)
// }

// var _ Stmt = &Class{}
// type Function struct {
// 	Name   Token
// 	Params []Token
// 	Body   []Stmt
// }

// func (f *Function) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitFunctionStmt(f)
// }

// var _ Stmt = &Function{}

// type If struct {
// 	Condition  Expr
// 	ThenBranch Stmt
// 	ElseBranch Stmt
// }

// func (i *If) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitIfStmt(i)
// }

// var _ Stmt = &If{}

// type Return struct {
// 	Keyword Token
// 	Value   Expr
// }

// func (r *Return) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitReturnStmt(r)
// }

// var _ Stmt = &Return{}

// type Var struct {
// 	Name        Token
// 	Initializer Expr
// }

// func (v *Var) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitVarStmt(v)
// }

// var _ Stmt = &Var{}

// type While struct {
// 	Condition Expr
// 	Body      Stmt
// }

// func (w *While) Accept(visitor StmtVisitor) interface{} {
// 	return visitor.VisitWhileStmt(w)
// }

// var _ Stmt = &While{}
