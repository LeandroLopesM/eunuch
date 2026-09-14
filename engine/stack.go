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

func (self *Stack[T]) Empty() bool {
	return self.ptr == 0
}

func (self *Stack[T]) Pop() (T, error) {
	if self.ptr-1 < self.guard {
		var def T
		return def, errors.New("Stack underflow")
	}

	self.ptr--
	return self.raw[self.ptr+1], nil
}

func (self *Stack[T]) Peek() (T, error) {
	val, err := self.Pop()
	self.ptr++

	return val, err
}

// func (self *Stack[T]) grow() {
// 	last := self.raw
// 	self.raw = make([]T, self.size*2)
// 	copy(last, self.raw)
// }

func (self *Stack[T]) Push(v T) {
	if self.ptr+1 >= len(self.raw) {
		panic("Stack overflow")
	}

	self.ptr++
	self.raw[self.ptr] = v
}

func (self *Engine) stackGuard() int {
	if v, e := self.stackHistory.Peek(); e == nil {
		return v
	} else {
		return 0
	}
}

func (self *Engine) saveStack() {
	self.stackHistory.Push(self.stack.ptr)
	self.stack.guard = self.stack.ptr
}

func (self *Engine) loadStack() {
	_ = util.Assert(self.stackHistory.Pop())
	self.stack.guard = self.stackGuard()
}
