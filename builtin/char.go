package builtin

import (
	"unicode"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func charOp(op string) BuiltinExec {
	exec := func (lhs rune, rhs rune) bool {
		switch op {
		case ">": return lhs > rhs
		case "<": return lhs < rhs
		case "=": return lhs == rhs
		case ">=": return lhs >= rhs
		case "<=": return lhs <= rhs
		}

		panic("Unreachable")
	}

	return func (e *Engine) error {
		// Stack pops inverted
		charB := rune(util.Assert(e.Pop()).(Char))
		charA := rune(util.Assert(e.Pop()).(Char))

		e.Push(Bool(exec(charA, charB)))
		return nil
	}
}

func charCiOp(op string) BuiltinExec {
	exec := func (lhs rune, rhs rune) bool {
		switch op {
		case ">": return lhs > rhs
		case "<": return lhs < rhs
		case "=": return lhs == rhs
		case ">=": return lhs >= rhs
		case "<=": return lhs <= rhs
		}

		panic("Unreachable")
	}

	return func (e *Engine) error {
		// Stack pops inverted
		charB := unicode.ToUpper(rune(util.Assert(e.Pop()).(Char)))
		charA :=  unicode.ToUpper(rune(util.Assert(e.Pop()).(Char)))

		e.Push(Bool(exec(charA, charB)))
		return nil
	}
}