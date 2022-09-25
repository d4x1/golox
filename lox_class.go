package main

import "fmt"

type LoxClass struct {
	name string
}

func newLoxClass(name string) *LoxClass {
	return &LoxClass{
		name: name,
	}
}

func (c *LoxClass) String() string {
	return fmt.Sprintf("<class: %s >", c.name)
}

func (c *LoxClass) Arity() int {
	return 0
}

func (c *LoxClass) Call(intp Interpreter, args []interface{}) (interface{}, error) {
	instance := newLoxInstance(c)
	return instance, nil
}
