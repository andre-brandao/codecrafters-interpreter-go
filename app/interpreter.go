package main

import (
	"fmt"
	"os"

	exp "github.com/codecrafters-io/interpreter-starter-go/app/expr"
	st "github.com/codecrafters-io/interpreter-starter-go/app/stmt"
	tok "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) InterpretExpression(expr exp.Expr) {

	defer func() {
		if r := recover(); r != nil {
			// fmt.Println("recovered")
			runTimeError, ok := r.(*RuntimeError)
			if ok {
				fmt.Fprint(os.Stderr, runTimeError.Error())
			} else {
				fmt.Fprint(os.Stderr, "Unknown error")
			}
			hadRuntimeError = true
			// fmt.Print("interpret error")
			// hadError = true
		}
	}()

	value := i.evaluate(expr)

	fmt.Print(stringfy(value))
}

func (i *Interpreter) Interpret(statements []st.Stmt) {
	defer func() {
		if r := recover(); r != nil {
			runTimeError, ok := r.(*RuntimeError)
			if ok {
				fmt.Fprint(os.Stderr, runTimeError.Error())
			} else {
				fmt.Fprint(os.Stderr, "Unknown error")
			}
			hadRuntimeError = true
		}
	}()

	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitLiteralExpr(expr *exp.Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *exp.Grouping) interface{} {
	return i.evaluate(expr.Expression)
}
func (i *Interpreter) VisitBinaryExpr(expr *exp.Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch op := expr.Operator.Type; op {

	case tok.GREATER:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case tok.GREATER_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case tok.LESS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case tok.LESS_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case tok.MINUS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)

	case tok.SLASH:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case tok.STAR:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)

	case tok.EQUAL_EQUAL:
		return isEqual(left, right)
	case tok.BANG_EQUAL:
		return !isEqual(left, right)

	case tok.PLUS:
		if isNumber(left) && isNumber(right) {

			return left.(float64) + right.(float64)
		}
		if isString(left) && isString(right) {
			return left.(string) + right.(string)
		}
		if isRune(left) && isRune(right) {
			return append(left.([]rune), right.([]rune)...)
		}

		panic(NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings."))
	}

	return nil
}

func (i *Interpreter) VisitUnaryExpr(expr *exp.Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch t := expr.Operator.Type; t {
	case tok.MINUS:
		checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case tok.BANG:
		return !isTruthy(right)

	}

	return nil
}

func (i *Interpreter) evaluate(expr exp.Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt st.Stmt) interface{} {
	return stmt.Accept(i)
}

func checkNumberOperand(operator tok.Token, operand interface{}) {
	if !isNumber(operand) {
		panic(NewRuntimeError(operator, "Operand must be a number."))
	}
}

func checkNumberOperands(operator tok.Token, left, right interface{}) {
	if !isNumber(left) || !isNumber(right) {
		panic(NewRuntimeError(operator, "Operands must be numbers."))
	}
}

// Stmt Visitor
func (i *Interpreter) VisitExpressionStmt(stmt *st.Expression) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *st.Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringfy(value))
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *st.Var) interface{} {
	return nil
}

func (i *Interpreter) VisitUnaryStmt(stmt *st.UnaryStmt) interface{} {
	return nil
}

func (i *Interpreter) VisitVariableStmt(stmt *st.Variable) interface{} {
	return nil
}

// func (i *Interpreter) VisitBlockStmt(stmt *Block) interface{} {
// 	return nil
// }
// func (i *Interpreter) VisitIfStmt(stmt *If) interface{} {
// 	return nil
// }

// func (i *Interpreter) VisitClassStmt(stmt *Class) interface{} {
// 	return nil
// }
// func (i *Interpreter) VisitFunctionStmt(stmt *Function) interface{} {
// 	return nil
// }

// func (i *Interpreter) VisitWhileStmt(stmt *While) interface{} {
// 	return nil
// }
// func (i *Interpreter) VisitReturnStmt(stmt *Return) interface{} {
// 	return nil
// }

var _ exp.ExprVisitor = (*Interpreter)(nil)
var _ st.StmtVisitor = (*Interpreter)(nil)
