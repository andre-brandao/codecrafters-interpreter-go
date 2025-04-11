package main

import (
	"github.com/codecrafters-io/interpreter-starter-go/app/token"

	exp "github.com/codecrafters-io/interpreter-starter-go/app/expr"
	st "github.com/codecrafters-io/interpreter-starter-go/app/stmt"
)

type FunctionType int

const (
	FunctionTypeNone FunctionType = iota
	FunctionTypeFunction
)

type ScopeStack []map[string]bool

func (s ScopeStack) isEmpty() bool {
	return len(s) == 0
}

func (s ScopeStack) Push(scope map[string]bool) {
	s = append(s, scope)
}

func (s ScopeStack) Pop() ScopeStack {
	return s[:len(s)-1]
}

func (s ScopeStack) Peek() map[string]bool {
	return s[len(s)-1]
}

type Resolver struct {
	interpreter     Interpreter
	scopes          ScopeStack
	currentFunction FunctionType
}

func NewResolver(interpreter Interpreter) Resolver {
	return Resolver{
		interpreter:     interpreter,
		scopes:          make(ScopeStack, 0),
		currentFunction: FunctionTypeNone,
	}
}

func (r *Resolver) VisitBlockStmt(stmt *st.Block) any {
	r.beginScope()
	r.resolveStmts(stmt.Statements)
	r.endScope()
	return nil
}

func (r *Resolver) VisitExpressionStmt(stmt *st.Expression) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitIfStmt(stmt *st.If) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.ThenBranch)
	if stmt.ElseBranch != nil {
		r.resolveStmt(stmt.ElseBranch)
	}
	return nil
}

func (r *Resolver) VisitPrintStmt(stmt *st.Print) any {
	r.resolveExpr(stmt.Expression)
	return nil
}

func (r *Resolver) VisitReturnStmt(stmt *st.Return) any {

	if r.currentFunction == FunctionTypeNone {
		Error(stmt.Keyword, "Can't return from top-level code.")
	}
	if stmt.Value != nil {
		r.resolveExpr(stmt.Value)
	}
	return nil
}

func (r *Resolver) VisitVarStmt(stmt *st.Var) any {
	r.declare(stmt.Name)

	if stmt.Initializer != nil {
		r.resolveExpr(stmt.Initializer)
	}
	r.define(stmt.Name)
	return nil
}

func (r *Resolver) VisitWhileStmt(stmt *st.While) any {
	r.resolveExpr(stmt.Condition)
	r.resolveStmt(stmt.Body)
	return nil
}

func (r *Resolver) VisitFunctionStmt(stmt *st.Function) any {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	r.resolveFunction(stmt, FunctionTypeFunction)
	return nil
}

func (r *Resolver) VisitVariableExpr(expr *exp.Variable) any {
	// if !r.scopes.isEmpty() && (r.scopes.Peek()[string(expr.Name.Lexeme)] == false) {
	// 	Error(expr.Name, "Can't read local variable in its own initializer.")
	// }
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitAssignExpr(expr *exp.Assign) any {
	r.resolveExpr(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil
}

func (r *Resolver) VisitBinaryExpr(expr *exp.Binary) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitCallExpr(expr *exp.Call) any {
	r.resolveExpr(expr.Callee)
	for _, arg := range expr.Arguments {
		r.resolveExpr(arg)
	}
	return nil
}

func (r *Resolver) VisitGroupingExpr(expr *exp.Grouping) any {
	r.resolveExpr(expr.Expression)
	return nil
}

func (r *Resolver) VisitLiteralExpr(expr *exp.Literal) any {
	return nil
}

func (r *Resolver) VisitLogicalExpr(expr *exp.Logical) any {
	r.resolveExpr(expr.Left)
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) VisitUnaryExpr(expr *exp.Unary) any {
	r.resolveExpr(expr.Right)
	return nil
}

func (r *Resolver) declare(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes.Peek()

	_, exists := scope[string(name.Lexeme)]
	if exists {
		Error(name, "Already a variable with this name in this scope.")
	}
	scope[string(name.Lexeme)] = false

	// r.scopes.Peek()[string(name.Lexeme)] = false

}

func (r *Resolver) define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}
	r.scopes.Peek()[string(name.Lexeme)] = true
}

func (r *Resolver) resolveLocal(expr exp.Expr, name token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if r.scopes[i][string(name.Lexeme)] {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *Resolver) resolveStmts(statements []st.Stmt) {
	for _, statement := range statements {
		statement.Accept(r)
	}
}

func (r *Resolver) resolveFunction(fn *st.Function, typ FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = typ

	r.beginScope()
	for _, param := range fn.Params {
		r.declare(param)
		r.define(param)
	}
	r.resolveStmts(fn.Body)
	r.endScope()

	r.currentFunction = enclosingFunction
}

func (r *Resolver) resolveStmt(statement st.Stmt) {
	statement.Accept(r)
}

func (r *Resolver) resolveExpr(expr exp.Expr) {
	expr.Accept(r)
}
func (r *Resolver) beginScope() {
	r.scopes = append(r.scopes, make(map[string]bool))
}
func (r *Resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

var _ exp.ExprVisitor = (*Resolver)(nil)
var _ st.StmtVisitor = (*Resolver)(nil)
