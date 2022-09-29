package main

import "fmt"

type LoxClass struct {
	superclass *LoxClass
	name       string
	methods    map[string]*LoxFunction
}

func newLoxClass(name string) *LoxClass {
	return &LoxClass{
		name: name,
	}
}

func newLoxClassWithMethods(name string, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:    name,
		methods: methods,
	}
}

func newLoxClassWithSuperClass(name string, superclass *LoxClass, methods map[string]*LoxFunction) *LoxClass {
	return &LoxClass{
		name:       name,
		methods:    methods,
		superclass: superclass,
	}
}

func (c *LoxClass) String() string {
	return fmt.Sprintf("<class: %s >", c.name)
}

func (c *LoxClass) Arity() int {
	initFunction, err := c.FindMethod("init")
	if err != nil {
		return 0
	}
	if initFunction != nil {
		return initFunction.Arity()
	}
	return 0
}

func (c *LoxClass) Call(intp Interpreter, args []interface{}) (interface{}, error) {
	instance := newLoxInstance(c)
	initFunction, err := c.FindMethod("init")
	if err != nil {
		return nil, err
	}
	if initFunction != nil {
		initliazer, err := initFunction.Bind(instance)
		if err != nil {
			return nil, err
		}
		return initliazer.Call(intp, args)
	}
	return instance, nil
}

func (c *LoxClass) FindMethod(name string) (*LoxFunction, error) {
	v, ok := c.methods[name]
	if ok {
		return v, nil
	}
	if c.superclass != nil {
		return c.superclass.FindMethod(name)
	}
	return nil, nil
}
