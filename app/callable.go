package main

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/environment"
	"github.com/codecrafters-io/interpreter-starter-go/app/stmt"
)

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments []any) any
	String() string
}

// implements LoxCallabel
type LoxFunction struct {
	declaration stmt.Function
}

func NewLoxFunction(declaration *stmt.Function) *LoxFunction {
	return &LoxFunction{
		declaration: *declaration,
	}
}

func (lf *LoxFunction) call(interp *Interpreter, arguments []any) any {
	env := environment.NewEnvironment(&interp.Globals)

	for i := 0; i < len(lf.declaration.Params); i++ {
		env.Define(string(lf.declaration.Params[i].Lexeme), arguments[i])
	}

	var returnValue any = nil

	func() {
		defer func() {
			if r := recover(); r != nil {
				if value, ok := r.(*Return); ok {
					returnValue = value.Value
				} else {
					panic("Unknow error")
				}
			}
		}()
		interp.executeBlock(lf.declaration.Body, env)
	}()

	return returnValue
}

func (lf *LoxFunction) arity() int {
	return len(lf.declaration.Params)
}

func (lf *LoxFunction) String() string {
	return fmt.Sprintf("<fn %s>", string(lf.declaration.Name.Lexeme))
}

var _ LoxCallable = (*LoxFunction)(nil)
