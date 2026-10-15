package core

import (
	"fmt"
	"reflect"
	"strings"
)

type Value interface {
	ToString() string
	TypeString() string
	Type() reflect.Type
}

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
	Symbol  string
	Integer int64
	Float   float64
	Bool    bool
	String  string
	Char    rune
	Vector  []Value
	Pair    [2]Value
	Scheme  struct {
		Args     []Value
		Position Position
	}
)

func (s Symbol) ToString() string   { return string(s) }
func (s Symbol) Type() reflect.Type { return reflect.TypeFor[Symbol]() }
func (s Symbol) TypeString() string { return "Symbol" }

func (s Integer) ToString() string   { return fmt.Sprintf("%b", s) }
func (s Integer) TypeString() string { return "Integer" }
func (s Integer) Type() reflect.Type { return reflect.TypeFor[Integer]() }

func (s Float) ToString() string   { return fmt.Sprintf("%b", s) }
func (s Float) TypeString() string { return "Float" }
func (s Float) Type() reflect.Type { return reflect.TypeFor[Float]() }

func (s Bool) ToString() string   { return fmt.Sprintf("%v", s) }
func (s Bool) TypeString() string { return "Bool" }
func (s Bool) Type() reflect.Type { return reflect.TypeFor[Bool]() }

func (s String) ToString() string   { return string(s) }
func (s String) TypeString() string { return "String" }
func (s String) Type() reflect.Type { return reflect.TypeFor[String]() }

func (s Char) ToString() string   { return fmt.Sprintf("%v", s) }
func (s Char) TypeString() string { return "Char" }
func (s Char) Type() reflect.Type { return reflect.TypeFor[Char]() }

func (s Vector) TypeString() string { return "Vector" }
func (s Vector) Type() reflect.Type { return reflect.TypeFor[Vector]() }
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

func (s Pair) TypeString() string { return "Pair" }
func (s Pair) Type() reflect.Type { return reflect.TypeFor[Pair]() }
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
		return s.Args[0].ToString()
	}
}

type Position struct {
	Line, Char int
	File       string
}

func (pos Position) ToString() string {
	return fmt.Sprintf("%s:%d:%d", pos.File, pos.Line, pos.Char)
}
