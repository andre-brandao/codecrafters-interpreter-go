package printer

import (
	"fmt"

	exp "github.com/codecrafters-io/interpreter-starter-go/app/expr"
)

type AstPrinter struct {
}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{}
}

func (p *AstPrinter) Print(expr exp.Expr) string {
	return expr.Accept(p).(string)
}

func (p *AstPrinter) VisitBinaryExpr(expr *exp.Binary) interface{} {
	return p.parenthesize(string(expr.Operator.Lexeme), expr.Left, expr.Right)
}

func (p *AstPrinter) VisitGroupingExpr(expr *exp.Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *AstPrinter) VisitLiteralExpr(expr *exp.Literal) interface{} {
	if expr.Value == nil {
		return "nil"
	}
	// return fmt.Sprintf("%v", expr.Value)
	switch v := expr.Value.(type) {
	case float64:
		if v == float64(int(v)) {
			return fmt.Sprintf("%.1f", v) // Ensures 1234.0 for whole numbers
		} else {
			return fmt.Sprintf("%g", v) // Keeps the precision for non-whole numbers

		}
	case string:
		return fmt.Sprintf("%s", v)
	case []rune:
		return fmt.Sprintf("%s", string(v))
	default:
		return fmt.Sprintf("%v", expr.Value)
	}
}

func (p *AstPrinter) VisitUnaryExpr(expr *exp.Unary) interface{} {
	return p.parenthesize(string(expr.Operator.Lexeme), expr.Right)
}

// func (p *AstPrinter) VisitVarExpr(expr *exp.Var) interface{} {
// 	// TODO: implement this
// 	return nil
// }

func (p *AstPrinter) VisitVariableExpr(expr *exp.Variable) interface{} {
	// return nil
	return fmt.Sprintf("%s", string(expr.Name.Lexeme))
}

func (p *AstPrinter) VisitAssignExpr(expr *exp.Assign) interface{} {
	// return nil
	return fmt.Sprintf("%s", string(expr.Name.Lexeme))
}

func (p *AstPrinter) VisitLogicalExpr(expr *exp.Logical) interface{} {
	// return p.parenthesize(string(expr.Operator.Lexeme), expr.Left, expr.Right)
	// return nil
	return p.parenthesize(string(expr.Operator.Lexeme), expr.Left, expr.Right)
}

func (p *AstPrinter) VisitCallExpr(expr *exp.Call) interface{} {
	// return nil
	// return p.parenthesizeSlice("call", expr.Callee, expr.Arguments)
	return p.parenthesizeSlice("call", append(expr.Arguments, expr.Callee))
}

func (p *AstPrinter) parenthesize(name string, exprs ...exp.Expr) string {
	var result string
	result += "(" + name

	for _, expr := range exprs {
		result += " "
		result += expr.Accept(p).(string)
	}

	result += ")"
	return result
}
func (p *AstPrinter) parenthesizeSlice(name string, exprs []exp.Expr) string {
	var result string
	result += "(" + name

	for _, expr := range exprs {
		result += " "
		result += expr.Accept(p).(string)
	}

	result += ")"
	return result
}

var _ exp.ExprVisitor = (*AstPrinter)(nil)
