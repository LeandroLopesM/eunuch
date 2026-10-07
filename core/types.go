package core

import (
	"fmt"
)

type Scheme struct {
	Args []Unit
	Position Position
}

func (s *Scheme) Params() []Unit {
	if len(s.Args) == 0 {
		return []Unit{}
	}

	return s.Args[1:]
}

func (s *Scheme) Name() string {
	if len(s.Args) == 0 {
		return ""
	} else {
		return s.Args[0].Value.(string)
	}
}

type PairVal = [2]Unit
type VectorVal = []Unit

func (pos Position) ToString() string {
	return fmt.Sprintf("%s:%d:%d", pos.File, pos.Line, pos.Char)
}

type Type int

const (
	// default type for anything, ignored
	NullType Type = iota 

	SchemeType
	Symbol
	Integer
	Float
	Bool
	String
	Char
	Vector
	Pair

	// Misc filters for arguments, not actual unit values
	None
	Any
	Number
)

func (t Type) ToString() string {
	switch t {
	case SchemeType: return "Scheme"
	case Symbol: return "Symbol"
	case Integer: return "Integer"
	case Float: return "Float"
	case Bool: return "Bool"
	case String: return "String"
	case Char: return "Char"
	case Vector: return "Vector"
	case Pair: return "Pair"

	default:
		return "Unknown"
	}
}

func (tf Type) Matches(other Type) bool {
	switch tf {
	case Any:
		return true
	case Number:
		return other == Integer || other == Float
	default:
		return other == tf
	}
}

type Position struct {
	Line, Char int
	File       string
}

type Unit struct {
	Type Type

	Value any
}

var Null = Unit { Type: NullType, Value: nil }

// For user-defined types.
// Be very careful using since a lot of functions accept 'Any'
// but then subsequently unwrap it, causing a panic
func MkCustom(v any) Unit {
	return Unit{
		Type:  Any,
		Value: v,
	}
}
func MkScheme(v Scheme) Unit {
	return Unit{
		Type:  SchemeType,
		Value: v,
	}
}
func MkInt(v int64) Unit {
	return Unit{
		Type:  Integer,
		Value: v,
	}
}
func MkFloat(v float64) Unit {
	return Unit{
		Type:  Float,
		Value: v,
	}
}
func MkSymbol(v string) Unit {
	return Unit{
		Type:  Symbol,
		Value: v,
	}
}
func MkBool(v bool) Unit {
	return Unit{
		Type:  Bool,
		Value: v,
	}
}
func MkString(v string) Unit {
	return Unit{
		Type:  String,
		Value: v,
	}
}
func MkChar(v rune) Unit {
	return Unit{
		Type:  Char,
		Value: v,
	}
}
func MkVector(v []Unit) Unit {
	return Unit{
		Type:  Vector,
		Value: v,
	}
}
func MkPair(v PairVal) Unit {
	return Unit{
		Type:  Pair,
		Value: v,
	}
}
