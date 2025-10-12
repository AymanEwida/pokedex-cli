package lib

import "errors"

type Stack[T any] struct {
	data []T
}

func MoveStackToStack[T any](stackA *Stack[T], stackB *Stack[T]) (n int, err error) {
	for range stackA.Size() {
		curr, err := stackA.Pop()
		if err != nil {
			return n, err
		}

		stackB.Push(curr)

		n++
	}

	return n, nil
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		data: []T{},
	}
}

func (s *Stack[any]) Size() int {
	return len(s.data)
}

func (s *Stack[T]) Push(item T) {
	s.data = append(s.data, item)
}

func (s *Stack[T]) Pop() (last T, err error) {
	if len(s.data) == 0 {
		return last, errors.New("Stack is empty")
	}

	last = s.data[len(s.data)-1]

	s.data = s.data[:len(s.data)-1]

	return last, nil
}

func (s *Stack[T]) Peek() (last T, err error) {
	if len(s.data) == 0 {
		return last, errors.New("Stack is empty")
	}

	return s.data[len(s.data)-1], nil
}
