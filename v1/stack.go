package main

import (
	"container/list"
	"errors"
)

var (
	errEmptyStack = errors.New("empty stack")
)

type stack struct {
	list *list.List
}

func newStack() *stack {
	return &stack{
		list: list.New(),
	}
}

func (s *stack) Push(value interface{}) error {
	s.list.PushBack(value)
	return nil
}

func (s *stack) Pop() (interface{}, error) {
	e := s.list.Back()
	if e != nil {
		s.list.Remove(e)
		return e.Value, nil
	}
	return nil, errEmptyStack
}

func (s *stack) IsEmpty() bool {
	return s.list.Len() == 0
}

func (s *stack) Size() int {
	return s.list.Len()
}

func (s *stack) Peek() (interface{}, error) {
	e := s.list.Back()
	if e != nil {
		return e.Value, nil
	}
	return nil, errEmptyStack
}

// 这里不暴露底层的 Element 更好。
func (s *stack) Back() *list.Element {
	return s.list.Back()
}
