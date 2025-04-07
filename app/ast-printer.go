package main

import (
    "fmt"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
    return &AstPrinter{}
}

func (p *AstPrinter) Print(expr Expr) string {
    return expr.Accept(p).(string)
}

func (p *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
    return p.parenthesize(string(expr.Operator.Lexeme), expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
    return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
    if expr.Value == nil {
        return "nil"
    }
    return fmt.Sprintf("%v", expr.Value)
}

func (p *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
    return p.parenthesize(string(expr.Operator.Lexeme), expr.Right)
}

func (p *AstPrinter) parenthesize(name string, exprs ...Expr) string {
    var result string
    result += "(" + name
    
    for _, expr := range exprs {
        result += " "
        result += expr.Accept(p).(string)
    }
    
    result += ")"
    return result
}

var _ ExprVisitor = (*AstPrinter)(nil)