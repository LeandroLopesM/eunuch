package builtin

import (
	"fmt"
	"math"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func OrdOp(kind rune) BuiltinExec {
	return func(e *Engine) error {
		var nums []Unit
		var overallType Type

		val, err := e.Pop()
		for err == nil {
			if len(nums) == 0 {
				overallType = val.Type
			} else if val.Type != overallType {
				return fmt.Errorf(
					"%s expects all arguments to be the same type (Was %s, now %s)",
					util.If(kind == '>',
						"(max)",
						"(min)",
					),
					TypeNames[overallType],
					TypeNames[val.Type],
				)
			}

			nums = append(nums, val)
			val, err = e.Pop()
		}

		switch overallType {
		case Float:
			var curr = numAsF(nums[0])
			for _, v := range nums {
				if util.If(kind == '>',
					numAsF(v) > curr,
					numAsF(v) < curr,
				) {
					curr = numAsF(v)
				}
			}

			e.Push(MkFloat(curr))
		default:
			var curr = nums[0].Value.(int64)
			for _, v := range nums {

				if util.If(
					kind == '>',
					v.Value.(int64) > curr,
					v.Value.(int64) < curr,
				) {
					curr = v.Value.(int64)
				}
			}

			e.Push(MkInt(curr))
		}

		return nil
	}
}

func numAsF(num Unit) float64 {
	switch num.Type {
	case Float:
		return num.Value.(float64)
	default:
		return float64(num.Value.(int64))
	}
}

func expt(e *Engine) error {
	lhs, lErr := e.Pop()
	rhs, rErr := e.Pop()

	if lErr != nil {
		return lErr
	} else if rErr != nil {
		return rErr
	}

	var lFloat = numAsF(lhs)
	var rFloat = numAsF(rhs)

	e.Push(MkFloat(math.Pow(lFloat, rFloat)))

	return nil
}

// TODO: (- 4) => -4
// TODO: (/ 4) => 1/4
func MathOp(kind rune) BuiltinExec {
	floatOp := func(a float64, b float64) float64 {
		switch kind {
		case '+':
			return a + b
		case '-':
			return a - b
		case '*':
			return a * b
		case '/':
			return a / b
		}

		panic(fmt.Sprintf("Undefined operation %c", kind))
	}

	intOp := func(a int64, b int64) int64 {
		switch kind {
		case '+':
			return a + b
		case '-':
			return a - b
		case '*':
			return a * b
		case '/':
			return a / b
		}

		panic(fmt.Sprintf("undefined operation %c", kind))
	}

	return func(e *Engine) error {
		filter := Number
		var numbers []Unit
		var outType = Integer // We can be optimistic, right?

		v, err := e.Pop()

		for err == nil {
			if !filter.Matches(v.Type) {
				return fmt.Errorf("expected integer or float, got %s", TypeNames[v.Type])
			}

			if v.Type == Float {
				outType = Float
			}

			numbers = append(numbers, v)

			v, err = e.Pop()
		}

		switch outType {
		case Float:
			var out = numAsF(numbers[0])
			for _, num := range numbers[1:] {
				out = floatOp(out, numAsF(num))
			}

			e.Push(MkFloat(out))
		default:
			var out = numbers[0].Value.(int64)
			for _, num := range numbers[1:] {
				out = intOp(out, num.Value.(int64))
			}

			e.Push(MkInt(out))
		}

		return nil
	}
}