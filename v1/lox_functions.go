package main

import (
	"errors"
)

type LoxFunction struct {
	name           string
	declaration    FunctionStmt
	closure        *Env
	isInitlializer bool
}

func newLoxFunction(stmt FunctionStmt, env *Env, isInitlializer bool) *LoxFunction {
	return &LoxFunction{
		declaration:    stmt,
		name:           stmt.name.Lexeme,
		closure:        env,
		isInitlializer: isInitlializer,
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
			if f.isInitlializer {
				return f.closure.GetAtByVarName(0, "this")
			}
			return returnValue.Value, nil
		}
		return nil, err
	}
	if f.isInitlializer {
		return f.closure.GetAtByVarName(0, "this")
	}
	return nil, nil
}

func (f *LoxFunction) Bind(instance *LoxInstance) (*LoxFunction, error) {
	env := newEnvWithEnclosing(f.closure)
	env.Define("this", instance)
	return newLoxFunction(f.declaration, env, f.isInitlializer), nil
}
