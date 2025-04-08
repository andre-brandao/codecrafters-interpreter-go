package main

import "time"

type clock struct{}

func (clock) arity() int                   { return 0 }
func (clock) call(*Interpreter, []any) any { return float64(time.Now().Unix()) }
func (clock) String() string               { return "<native fn>" }

var _ LoxCallable = (*clock)(nil)
