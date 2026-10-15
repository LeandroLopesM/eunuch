package core

import (
	"fmt"
	"reflect"
	"strings"
)

type Value struct {Raw any}

func (v Value) ToString() string {
	switch val := v.Raw.(type) {
	case Symbol, String:
		return val.(string)
	case Integer, Bool, Float, Char:
		return fmt.Sprintf("%v", val)
	case Vector:
		return val.ToString()
	case Pair:
		return val.ToString()

	default:
		return fmt.Sprintf("Custom %v", val)
	}
}

func (v *Value) Type() reflect.Type {
	return reflect.TypeOf(v.Raw)
}

func (v Value) TypeString() string {
	return fmt.Sprintf("%T", v.Raw)
}

func (v Value) Format() string {
	return fmt.Sprintf("%s: %s", v.TypeString(), v.ToString())
}

type (
	Symbol string
	Integer int64
	Float float64
	Bool bool
	String string
	Char rune
	Vector []Value
	Pair [2]Value
	Scheme struct {
		Args []Value
		Position Position
	}
)

func SymbolVal(v string) Value {
	return Value { Raw: v }
}
func IntegerVal(v int64) Value {
	return Value { Raw: v }
}
func FloatVal(v float64) Value {
	return Value { Raw: v }
}
func BoolVal(v bool) Value {
	return Value { Raw: v }
}
func StringVal(v string) Value {
	return Value { Raw: v }
}
func CharVal(v rune) Value {
	return Value { Raw: v }
}
func VectorVal(v []Value) Value {
	return Value { Raw: v }
}
func PairVal(v [2]Value) Value {
	return Value { Raw: v }
}
func SchemeVal(v Scheme) Value {
	return Value { Raw: v }
}

func (v Vector) ToString() string {
	var ret strings.Builder
	for i := range v {
		if i != 0 {
			ret.WriteString(v[i].ToString())
		} else {
			ret.WriteString("#(")
			ret.WriteString(v[i].ToString())
		}
	}
	return ret.String() + ")"
}

func (p Pair) ToString() string {
	return fmt.Sprintf("(%v . %v)", p[0].ToString(), p[1].ToString())
}


func (s *Scheme) Params() []Value {
	if len(s.Args) == 0 {
		return []Value{}
	}

	return s.Args[1:]
}

func (s *Scheme) Name() string {
	if len(s.Args) == 0 {
		return ""
	} else {
		return s.Args[0].Raw.(string)
	}
}

type Position struct {
	Line, Char int
	File       string
}

func (pos Position) ToString() string {
	return fmt.Sprintf("%s:%d:%d", pos.File, pos.Line, pos.Char)
}