package engine

import (
	"errors"

	"github.com/leandrolopesm/eunuch/util"
)

const STACK_SIZE = 64

type Stack[T any] struct {
	raw []T
	ptr int

	guard int
	size  int
}

func NewStack[T any]() Stack[T] {
	return Stack[T]{
		raw: make([]T, STACK_SIZE),
		ptr: 0,

		guard: 0,
		size:  STACK_SIZE,
	}
}

func (stack *Stack[T]) Empty() bool {
	return stack.ptr == 0
}

func (stack *Stack[T]) Pop() (T, error) {
	if stack.ptr-1 < stack.guard {
		var def T
		return def, errors.New("Stack underflow")
	}

	stack.ptr--
	return stack.raw[stack.ptr+1], nil
}

func (stack *Stack[T]) Peek() (T, error) {
	val, err := stack.Pop()
	stack.ptr++

	return val, err
}

// func (stack *Stack[T]) grow() {
// 	last := stack.raw
// 	stack.raw = make([]T, stack.size*2)
// 	copy(last, stack.raw)
// }

func (stack *Stack[T]) Push(v T) {
	if stack.ptr+1 >= len(stack.raw) {
		panic("Stack overflow")
	}

	stack.ptr++
	stack.raw[stack.ptr] = v
}

func (stack *Engine) stackGuard() int {
	if v, e := stack.stackHistory.Peek(); e == nil {
		return v
	} else {
		return 0
	}
}

func (stack *Engine) saveStack() {
	stack.stackHistory.Push(stack.stack.ptr)
	stack.stack.guard = stack.stack.ptr
}

func (stack *Engine) loadStack() {
	_ = util.Assert(stack.stackHistory.Pop())
	stack.stack.guard = stack.stackGuard()
}
