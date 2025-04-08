package main

import "time"

type clock struct{}

func (clock) arity() int                   { return 0 }
func (clock) call(*Interpreter, []any) any { return (time.Now().UnixNano()) }
func (clock) String() string               { return "<native fn>" }

var _ LoxCallable = (*clock)(nil)
