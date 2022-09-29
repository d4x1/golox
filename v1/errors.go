package main

import (
	"errors"
	"fmt"
)

var (
	errCastToMapString2Bool  = errors.New("interface{} is not a map[string]bool")
	errCastStmt2FunctionStmt = errors.New("stmt is not a function stmt")
)

type Return struct {
	Value interface{}
}

func NewReturn(v interface{}) Return {
	return Return{
		Value: v,
	}
}

func (r Return) Error() string {
	return fmt.Sprintf("return value is: %v", r.Value)
}
