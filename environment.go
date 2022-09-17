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

func (env *Env) Define(name string, value interface{}) {
	env.data[name] = value
}

func (env *Env) Get(name string) (interface{}, error) {
	if v, ok := env.data[name]; ok {
		return v, nil
	}
	if env.enclosing != nil {
		return env.enclosing.Get(name)
	}
	return nil, fmt.Errorf("undefined %s ", name)
}

func (env *Env) Assign(name string, value interface{}) error {
	if _, ok := env.data[name]; ok {
		env.data[name] = value
		return nil
	}
	if env.enclosing != nil {
		return env.enclosing.Assign(name, value)
	}
	return fmt.Errorf("undefined variable %s ", name)
}
