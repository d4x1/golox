package main

import "fmt"

type Interpreter interface {
	EvalVisitor
	StmtVisitor

	GetGlobalEnv() *Env
	ExecuteBlock(stmts []Stmt, env *Env) error
	Resolve(expr Expr, distance int) error
}

type Callable interface {
	Arity() int
	Call(intp Interpreter, args []interface{}) (interface{}, error)

	fmt.Stringer
}
