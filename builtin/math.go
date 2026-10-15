package builtin

import (
	"fmt"
	"log"
	"math"
	"reflect"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
)

func ordOp(kind rune) BuiltinExec {
	return func(e *Engine) error {
		argv := genArgv(e)
		applies := func (prev, new float64) bool {
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
			if applies(tmp.Raw.(float64), arg.Raw.(float64)) {
				tmp = arg
			}
		}

		e.Push(tmp)
		return nil
	}
}

func mathOp(kind rune) BuiltinExec {
	return func(e *Engine) error {
		argv := genArgv(e)

		if len(argv) == 1 {
			return unaryOp(kind, argv[0], e)
		}
		return binaryOp(kind, argv, e)
	}
}

func binaryOp(kind rune, operands []Value, e *Engine) error {
	type number interface {int64|float64}
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

	var (
		intType = reflect.TypeFor[Integer]()
	)

	finalType := intType
	output := 0.0

	for _,op := range operands {
		// If it was an integer and now is a float we need to assume the output will be float
		if op.Type() != finalType && finalType == intType {
			finalType = op.Type()
		}

		switch val := op.Raw.(type) {
		case Integer, Float:
			output = applyOp(output, val.(float64))
		default:
			return fmt.Errorf("invalid type '%s' for mathematical operation", op.TypeString())
		}
	}

	if finalType == intType {
		e.Push(IntegerVal(int64(output)))
	} else {
		e.Push(FloatVal(output))
	}

	return nil
}

func unaryOp(kind rune, operand Value, e *Engine) error {
	switch kind {
	case '-':
		if asInt,ok := operand.Raw.(Integer); ok {
			e.Push(IntegerVal(int64(asInt) * -1))
		} else if asFloat,ok := operand.Raw.(Integer); ok {
			e.Push(FloatVal(float64(asFloat) * -1.0))
		} else {
			return fmt.Errorf("invalid type for mathematical operation '%T'", operand)
		}
	case '/':
		if asInt,ok := operand.Raw.(Integer); ok {
			e.Push(IntegerVal(1/int64(asInt)))
		} else if asFloat,ok := operand.Raw.(Integer); ok {
			e.Push(FloatVal(1/float64(asFloat)))
		} else {
			return fmt.Errorf("invalid type for mathematical operation '%T'", operand)
		}
	default:
		return fmt.Errorf("math operation '%c' does not support unary usage", kind)
	}

	return nil
}

func expt(eng *Engine) error {
	lhs,_ := eng.Pop()
	rhs,_ := eng.Pop()

	var lFloat = lhs.Raw.(float64)
	var rFloat = rhs.Raw.(float64)

	eng.Push(FloatVal(math.Pow(lFloat, rFloat)))

	return nil
}