package main

import "fmt"

type LoxInstance struct {
	class *LoxClass
}

func newLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class: class,
	}
}

func (i *LoxInstance) String() string {
	return fmt.Sprintf("<class: %s's instance>", i.class.name)
}
