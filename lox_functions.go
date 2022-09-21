package main

import (
	"errors"
	"fmt"
)

type LoxFunction struct {
	name        string
	declaration FunctionStmt
}

func newLoxFunction(stmt FunctionStmt) *LoxFunction {
	return &LoxFunction{
		declaration: stmt,
		name:        stmt.name.Lexeme,
	}
}

func (f *LoxFunction) String() string {
	return "<function: " + f.name + ">"
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) Call(intp Interpreter, args []interface{}) (interface{}, error) {
	env := newEnvWithEnclosing(intp.GetGlobalEnv())
	for i, v := range f.declaration.params {
		env.Define(v.Lexeme, args[i])
	}
	fmt.Printf("call %+v\n", env)
	if err := intp.ExecuteBlock(f.declaration.stmts, env); err != nil {
		var returnValue Return
		if errors.As(err, &returnValue) {
			return returnValue.Value, nil
		}
		return nil, err
	}
	return nil, nil
}
