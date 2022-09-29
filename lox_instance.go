package main

import "fmt"

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func newLoxInstance(class *LoxClass) *LoxInstance {
	instance := &LoxInstance{
		class: class,
	}
	instance.fields = make(map[string]interface{})
	return instance
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("<class: %s's instance>", i.class.name)
}

func (i *LoxInstance) Get(name token) (interface{}, error) {
	v, ok := i.fields[name.Lexeme]
	if ok {
		return v, nil
	}

	if v, err := i.class.FindMethod(name.Lexeme); err != nil {
		return nil, err
	} else if v != nil {
		return v.Bind(i)
	}

	return nil, fmt.Errorf("%s not found in this instance", name.Lexeme)
}

func (i *LoxInstance) Set(name token, value interface{}) error {
	i.fields[name.Lexeme] = value
	return nil
}
