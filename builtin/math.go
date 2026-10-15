package builtin

import (
	"fmt"
	"log"
	"math"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func ordOp(kind rune) BuiltinExec {
	return func(e *Engine) error {
		argv := e.GenArgv()
		applies := func (prev, new Float) bool {
			switch kind {
			case '>':
				return new > prev
			case '<':
				return new < prev
			}
			panic("unreachable")
		}

		var tmp Value
		for _,arg := range argv {
			if applies(tmp.(Float), arg.(Float)) {
				tmp = arg
			}
		}

		e.Push(tmp)
		return nil
	}
}

func mathOp(kind rune) BuiltinExec {
	return func(e *Engine) error {
		argv := e.GenArgv()

		if len(argv) == 1 {
			return unaryOp(kind, argv[0], e)
		}
		return binaryOp(kind, argv, e)
	}
}

func binaryOp(kind rune, operands []Value, e *Engine) error {
	type number interface {int64|Float}
	var applyOp = func(lhs, rhs float64) float64 {
		switch kind {
		case '-':
			return lhs - rhs
		case '+':
			return lhs + rhs
		case '/':
			return lhs / rhs
		case '*':
			return lhs * rhs
		}

		log.Fatalf("Unknown operator %c", kind)
		panic("Unreachable")
	}

	finalType := "Integer"
	output := 0.0

	for _,op := range operands {
		// If it was an integer and now is a float we need to assume the output will be float
		if op.TypeString() != finalType && finalType == "Integer" {
			finalType = op.TypeString()
		}

		var currVal float64
		switch val := op.(type) {
		// These two can't be merged, the compiler isn't *that* smart
		case Integer:
			currVal = float64(val)
		case Float:
			currVal = float64(val)
		default:
			return fmt.Errorf("invalid type '%s' for mathematical operation", op.TypeString())
		}

		output = applyOp(output, currVal)
	}

	// Mod is to check the output is not decimal (We dont want to accidentally truncate it)
	if finalType == "Integer" && math.Mod(output, 1.0) == 0 {
		e.Push(Integer(int64(output)))
	} else {
		e.Push(Float(output))
	}

	return nil
}

func unaryOp(kind rune, operand Value, e *Engine) error {
	switch kind {
	case '-':
		if asInt,ok := operand.(Integer); ok {
			e.Push(Integer(int64(asInt) * -1))
		} else if asFloat,ok := operand.(Integer); ok {
			e.Push(Float(Float(asFloat) * -1.0))
		} else {
			return fmt.Errorf("invalid type for mathematical operation '%T'", operand)
		}
	case '/':
		if asInt,ok := operand.(Integer); ok {
			e.Push(Integer(1/int64(asInt)))
		} else if asFloat,ok := operand.(Integer); ok {
			e.Push(Float(1/Float(asFloat)))
		} else {
			return fmt.Errorf("invalid type for mathematical operation '%T'", operand)
		}
	default:
		return fmt.Errorf("math operation '%c' does not support unary usage", kind)
	}

	return nil
}

func expt(eng *Engine) error {
	var lFloat = float64(util.Assert(eng.Pop()).(Float))
	var rFloat = float64(util.Assert(eng.Pop()).(Float))

	eng.Push(Float(math.Pow(lFloat, rFloat)))

	return nil
}