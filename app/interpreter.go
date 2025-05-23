package main

import (
	"fmt"
	"os"

	env "github.com/codecrafters-io/interpreter-starter-go/app/environment"
	err "github.com/codecrafters-io/interpreter-starter-go/app/err"
	exp "github.com/codecrafters-io/interpreter-starter-go/app/expr"
	st "github.com/codecrafters-io/interpreter-starter-go/app/stmt"
	"github.com/codecrafters-io/interpreter-starter-go/app/token"
	tok "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Interpreter struct {
	Globals    env.Environment
	enviroment *env.Environment
	locals     map[exp.Expr]int
}

func NewInterpreter() *Interpreter {
	globals := *env.NewEnvironment(nil)

	globals.Define("clock", &clock{})
	return &Interpreter{
		Globals:    globals,
		enviroment: &globals,
		locals:     make(map[exp.Expr]int),
	}
}

func (i *Interpreter) Resolve(expr exp.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) InterpretExpression(expr exp.Expr) {

	defer func() {
		if r := recover(); r != nil {
			// fmt.Println("recovered")
			runTimeError, ok := r.(*err.RuntimeError)
			if ok {
				fmt.Fprint(os.Stderr, runTimeError.Error())
				hadRuntimeError = true
			} else {
				fmt.Fprint(os.Stderr, r)
				fmt.Fprint(os.Stderr, "Unknown error")
			}
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
			runTimeError, ok := r.(*err.RuntimeError)
			if ok {
				fmt.Fprint(os.Stderr, runTimeError.Error())
			} else {
				fmt.Fprintln(os.Stderr, "Unknown error")
				fmt.Fprintln(os.Stderr, r)
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
func (i *Interpreter) VisitLogicalExpr(expr *exp.Logical) interface{} {
	left := i.evaluate(expr.Left)

	if expr.Operator.Type == tok.OR {
		if isTruthy(left) {
			return left
		}
	} else {
		if !isTruthy(left) {
			return left
		}
	}

	return i.evaluate(expr.Right)
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

		panic(err.NewRuntimeError(expr.Operator, "Operands must be two numbers or two strings."))
	}

	return nil
}

func (i *Interpreter) VisitCallExpr(expr *exp.Call) any {
	callee := i.evaluate(expr.Callee)

	arguments := make([]any, 0)

	for _, arg := range expr.Arguments {
		arguments = append(arguments, i.evaluate(arg))
	}

	function, ok := callee.(LoxCallable)

	if len(arguments) != function.arity() {
		panic(err.NewRuntimeError(expr.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.arity(), len(arguments))))
	}

	if !ok {
		panic(err.NewRuntimeError(expr.Paren, "Can only call functions and classes."))
	}

	return function.call(i, arguments)
}

func (i *Interpreter) VisitUnaryExpr(expr *exp.Unary) any {
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

//	func (i *Interpreter) VisitVarExpr(expr *exp.Var) interface{} {
//		return i.enviroment.Get(expr.Name)
//	}
func (i *Interpreter) VisitVariableExpr(expr *exp.Variable) any {
	return i.lookUpVariable(expr.Name, expr)
	// return i.enviroment.Get(expr.Name)
}

func (i *Interpreter) lookUpVariable(name token.Token, expr *exp.Variable) any {
	distance := i.locals[expr]
	if distance != -1 {
		return i.enviroment.GetAt(distance, name)
	}
	return i.Globals.Get(name)
}

func (i *Interpreter) evaluate(expr exp.Expr) any {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt st.Stmt) any {
	// fmt.Println("execute stmt ----------- ", stmt)
	// i.enviroment.Print()

	return stmt.Accept(i)
}
func (i *Interpreter) executeBlock(statements []st.Stmt, environment *env.Environment) {
	previous := i.enviroment
	i.enviroment = environment
	defer func() {
		// if r := recover(); r != nil {
		// 	// recover from panic
		// 	if value, ok := r.(*Return); ok {
		// 		panic(value)
		// 	} else {
		// 		panic("Unknow error")
		// 	}
		// 	// fmt.Println("recover from panic executeBlock")
		// 	// fmt.Println(r)
		// }
		i.enviroment = previous
	}()
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i *Interpreter) VisitBlockStmt(stmt *st.Block) any {
	i.executeBlock(stmt.Statements, env.NewEnvironment(i.enviroment))
	return nil
}

func checkNumberOperand(operator tok.Token, operand any) {
	if !isNumber(operand) {
		panic(err.NewRuntimeError(operator, "Operand must be a number."))
	}
}

func checkNumberOperands(operator tok.Token, left, right any) {
	if !isNumber(left) || !isNumber(right) {
		panic(err.NewRuntimeError(operator, "Operands must be numbers."))
	}
}

// Stmt Visitor
func (i *Interpreter) VisitExpressionStmt(stmt *st.Expression) any {
	i.evaluate(stmt.Expression)
	return nil
}

func (i *Interpreter) VisitFunctionStmt(stmt *st.Function) any {
	function := NewLoxFunction(stmt, *i.enviroment)
	i.enviroment.Define(string(stmt.Name.Lexeme), function)
	return nil
}

func (i *Interpreter) VisitIfStmt(stmt *st.If) any {
	if isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}
	return nil
}

func (i *Interpreter) VisitPrintStmt(stmt *st.Print) any {
	value := i.evaluate(stmt.Expression)
	fmt.Println(stringfy(value))
	return nil
}

func (i *Interpreter) VisitReturnStmt(stmt *st.Return) any {
	var value any = nil
	if stmt.Value != nil {
		value = i.evaluate(stmt.Value)
	}
	panic(&Return{Value: value})
}

func (i *Interpreter) VisitVarStmt(stmt *st.Var) any {
	var value any = nil

	if stmt.Initializer != nil {
		value = i.evaluate(stmt.Initializer)
	}
	i.enviroment.Define(string(stmt.Name.Lexeme), value)
	return nil
}

func (i *Interpreter) VisitWhileStmt(stmt *st.While) any {
	for isTruthy(i.evaluate(stmt.Condition)) {
		i.execute(stmt.Body)
	}

	return nil
}
func (i *Interpreter) VisitAssignExpr(expr *exp.Assign) any {
	value := i.evaluate(expr.Value)

	// i.enviroment.Assign(expr.Name, value)
	distance, exists := i.locals[expr]
	if exists {
		i.enviroment.AssignAt(distance, expr.Name, value)
	} else {
		i.enviroment.Assign(expr.Name, value)
	}
	return value
}

// func (i *Interpreter) VisitUnaryStmt(stmt *st.UnaryStmt) interface{} {
// 	return nil
// }

// func (i *Interpreter) VisitVariableStmt(stmt *st.Variable) interface{} {
// 	return nil
// }

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
