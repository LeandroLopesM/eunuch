package builtin

import (
	"unicode"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func charOp(op string) BuiltinExec {
	exec := func (a rune, b rune) bool {
		switch op {
		case ">": return a > b
		case "<": return a < b
		case "=": return a == b
		case ">=": return a >= b
		case "<=": return a <= b
		}

		panic("Unreachable")
	}

	return func (e *Engine) error {
		// Stack pops inverted
		charB := util.Assert(e.Pop()).Value.(rune)
		charA := util.Assert(e.Pop()).Value.(rune)

		e.Push(MkBool(exec(charA, charB)))
		return nil
	}
}

func charCiOp(op string) BuiltinExec {
	exec := func (a rune, b rune) bool {
		switch op {
		case ">": return a > b
		case "<": return a < b
		case "=": return a == b
		case ">=": return a >= b
		case "<=": return a <= b
		}

		panic("Unreachable")
	}

	return func (e *Engine) error {
		// Stack pops inverted
		charB := unicode.ToUpper(util.Assert(e.Pop()).Value.(rune))
		charA :=  unicode.ToUpper(util.Assert(e.Pop()).Value.(rune))

		e.Push(MkBool(exec(charA, charB)))
		return nil
	}
}