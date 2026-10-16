package core

import (
	"fmt"
	"reflect"
)

type Value interface {
	ToString() string
	TypeString() string
	Type() reflect.Type
}

func Is[T Value](v Value) bool {
	return v.Type() == reflect.TypeFor[T]()
}

func Format(v Value) string {
	return fmt.Sprintf("%s: %s", v.TypeString(), v.ToString())
}

type Position struct {
	Line, Char int
	File       string
}

func (pos Position) ToString() string {
	return fmt.Sprintf("%s:%d:%d", pos.File, pos.Line, pos.Char)
}
