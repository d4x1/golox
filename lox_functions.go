package main

import (
	"errors"
)

type LoxFunction struct {
	name        string
	declaration FunctionStmt
	closure     *Env
}

func newLoxFunction(stmt FunctionStmt, env *Env) *LoxFunction {
	return &LoxFunction{
		declaration: stmt,
		name:        stmt.name.Lexeme,
		closure:     env,
	}
}

func (f *LoxFunction) String() string {
	return "<function: " + f.name + ">"
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) Call(intp Interpreter, args []interface{}) (interface{}, error) {
	env := newEnvWithEnclosing(f.closure)
	for i, v := range f.declaration.params {
		env.Define(v.Lexeme, args[i])
	}
	if err := intp.ExecuteBlock(f.declaration.stmts, env); err != nil {
		var returnValue Return
		if errors.As(err, &returnValue) {
			return returnValue.Value, nil
		}
		return nil, err
	}
	return nil, nil
}
