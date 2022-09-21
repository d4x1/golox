package main

import "time"

type nativeFunctionClock struct{}

func newNativeFunctionClock() *nativeFunctionClock {
	return &nativeFunctionClock{}
}

func (nativeFunctionClock) String() string {
	return "clock"
}

func (nativeFunctionClock) Arity() int {
	return 0
}

func (nativeFunctionClock) Call(intp Interpreter, args []interface{}) (interface{}, error) {
	return time.Now().UnixMilli(), nil
}
