package main

import "fmt"

type Interpreter interface {
	EvalVisitor
	StmtVisitor
	GetGlobalEnv() *Env
	ExecuteBlock(stmts []Stmt, env *Env) error
}

type Callable interface {
	fmt.Stringer
	Arity() int
	Call(intp Interpreter, args []interface{}) (interface{}, error)
}
