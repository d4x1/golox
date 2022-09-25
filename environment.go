package main

import (
	"fmt"
)

type Env struct {
	enclosing *Env
	data      map[string]interface{}
}

func newEnv() *Env {
	env := &Env{}
	env.data = make(map[string]interface{}, 0)
	return env
}

func newEnvWithEnclosing(enclosingEnv *Env) *Env {
	env := &Env{
		enclosing: enclosingEnv,
	}
	env.data = make(map[string]interface{}, 0)
	return env
}

func (env *Env) Define(varName string, value interface{}) {
	env.data[varName] = value
}

func (env *Env) Get(name token) (interface{}, error) {
	if v, ok := env.data[name.Lexeme]; ok {
		return v, nil
	}
	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}
	return nil, fmt.Errorf("undefined variable %s when getting", name.Lexeme)
}

func (env *Env) Ancestor(distance int) (*Env, error) {
	destination := env
	for i := 0; i < distance; i++ {
		if destination.enclosing == nil {
			return destination, fmt.Errorf("nil enclosing")
		}
		destination = destination.enclosing
	}
	return destination, nil
}

func (env *Env) GetAtByVarName(distance int, varName string) (interface{}, error) {
	destinationEnv, err := env.Ancestor(distance)
	if err != nil {
		return nil, err
	}
	v, ok := destinationEnv.data[varName]
	if ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected get, var name %s not found in destination env, maybe something wrong?", varName)
}

func (env *Env) Assign(name token, value interface{}) error {
	if _, ok := env.data[name.Lexeme]; ok {
		env.data[name.Lexeme] = value
		return nil
	}
	if env.enclosing != nil {
		return env.enclosing.Assign(name, value)
	}
	return fmt.Errorf("undefined variable %s when assiging", name.Lexeme)
}

func (env *Env) AssignAt(distance int, name token, value interface{}) error {
	destinationEnv, err := env.Ancestor(distance)
	if err != nil {
		return err
	}
	_, ok := destinationEnv.data[name.Lexeme]
	if !ok {
		return fmt.Errorf("unexpected assign, name %s not found in destination env, maybe something wrong?", name)
	}
	destinationEnv.data[name.Lexeme] = value
	return nil
}
