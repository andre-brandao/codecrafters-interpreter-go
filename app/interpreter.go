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

func (i *Interpreter) Interpret(expr Expr) {

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

var _ ExprVisitor = (*Interpreter)(nil)

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
