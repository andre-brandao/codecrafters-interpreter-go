package main

type LoxCallable interface {
	arity() int
	call(interpreter *Interpreter, arguments []any) any
	String() string
}
