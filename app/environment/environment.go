package environment

import (
	"fmt"

	err "github.com/codecrafters-io/interpreter-starter-go/app/err"
	tok "github.com/codecrafters-io/interpreter-starter-go/app/token"
)

type Environment struct {
	values    map[string]any
	enclosing *Environment
}

func NewEnvironment(env *Environment) *Environment {
	return &Environment{
		values:    make(map[string]any),
		enclosing: env,
	}
}

func (e *Environment) Get(name tok.Token) any {
	if value, ok := e.values[string(name.Lexeme)]; ok {
		return value
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	panic(err.NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", string(name.Lexeme))))
}

func (e *Environment) Assign(name tok.Token, value any) {
	if _, ok := e.values[string(name.Lexeme)]; ok {
		e.values[string(name.Lexeme)] = value
		return
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
		return
	}
	panic(err.NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", string(name.Lexeme))))
}

func (e *Environment) Define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) Print() {
	for k, v := range e.values {
		fmt.Println(k, v)
	}
}
