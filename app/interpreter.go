package main

import (
	"fmt"
	"os"
)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) InterpretExpression(expr Expr) {

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

func (i *Interpreter) Interpret(statements []Stmt) {
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

func (i *Interpreter) VisitLiteralExpr(expr *Literal) interface{} {
	return expr.Value
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) interface{} {
	return i.evaluate(expr.Expression)
}
func (i *Interpreter) VisitBinaryExpr(expr *Binary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch op := expr.Operator.Type; op {

	case GREATER:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) >= right.(float64)
	case LESS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) <= right.(float64)
	case MINUS:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) - right.(float64)

	case SLASH:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) / right.(float64)
	case STAR:
		checkNumberOperands(expr.Operator, left, right)
		return left.(float64) * right.(float64)

	case EQUAL_EQUAL:
		return isEqual(left, right)
	case BANG_EQUAL:
		return !isEqual(left, right)

	case PLUS:
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

func (i *Interpreter) VisitUnaryExpr(expr *Unary) interface{} {
	right := i.evaluate(expr.Right)

	switch t := expr.Operator.Type; t {
	case MINUS:
		checkNumberOperand(expr.Operator, right)
		return -right.(float64)
	case BANG:
		return !isTruthy(right)

	}

	return nil
}

func (i *Interpreter) evaluate(expr Expr) interface{} {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt Stmt) interface{} {
	return stmt.Accept(i)
}

func checkNumberOperand(operator Token, operand interface{}) {
	if !isNumber(operand) {
		panic(NewRuntimeError(operator, "Operand must be a number."))
	}
}

func checkNumberOperands(operator Token, left, right interface{}) {
	if !isNumber(left) || !isNumber(right) {
		panic(NewRuntimeError(operator, "Operands must be numbers."))
	}
}

// Stmt Visitor
func (i *Interpreter) VisitExpressionStmt(stmt *Expression) interface{} {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *Print) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringfy(value))
	return nil
}

func (i *Interpreter) VisitVarStmt(stmt *Var) interface{} {
	return nil
}

func (i *Interpreter) VisitUnaryStmt(stmt *UnaryStmt) interface{} {
    return nil
}

func (i *Interpreter) VisitVariableStmt(stmt *Variable) interface{} {
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

var _ ExprVisitor = (*Interpreter)(nil)
var _ StmtVisitor = (*Interpreter)(nil)
