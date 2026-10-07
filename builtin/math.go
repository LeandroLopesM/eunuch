package builtin

import (
	"fmt"
	"math"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func OrdOp(kind rune) BuiltinExec {
	return func(eng *Engine) error {
		var nums []Unit
		var overallType Type

		val, err := eng.Pop()
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
					overallType.ToString(),
					val.Type.ToString(),
				)
			}

			nums = append(nums, val)
			val, err = eng.Pop()
		}

		switch overallType {
		case Float:
			var curr = numAsF(nums[0])
			for _, num := range nums {
				if util.If(kind == '>',
					numAsF(num) > curr,
					numAsF(num) < curr,
				) {
					curr = numAsF(num)
				}
			}

			eng.Push(MkFloat(curr))
		default:
			var curr = nums[0].Value.(int64)
			for _, num := range nums {
				if util.If(
					kind == '>',
					num.Value.(int64) > curr,
					num.Value.(int64) < curr,
				) {
					curr = num.Value.(int64)
				}
			}

			eng.Push(MkInt(curr))
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

func expt(eng *Engine) error {
	lhs, lErr := eng.Pop()
	rhs, rErr := eng.Pop()

	if lErr != nil {
		return lErr
	} else if rErr != nil {
		return rErr
	}

	var lFloat = numAsF(lhs)
	var rFloat = numAsF(rhs)

	eng.Push(MkFloat(math.Pow(lFloat, rFloat)))

	return nil
}

// TODO: (- 4) => -4.
// TODO: (/ 4) => 1/4.
func MathOp(kind rune) BuiltinExec {
	floatOp := func(lhs float64, rhs float64) float64 {
		switch kind {
		case '+':
			return lhs + rhs
		case '-':
			return lhs - rhs
		case '*':
			return lhs * rhs
		case '/':
			return lhs / rhs
		}

		panic(fmt.Sprintf("Undefined operation %c", kind))
	}

	intOp := func(lhs int64, rhs int64) int64 {
		switch kind {
		case '+':
			return lhs + rhs
		case '-':
			return lhs - rhs
		case '*':
			return lhs * rhs
		case '/':
			return lhs / rhs
		}

		panic(fmt.Sprintf("undefined operation %c", kind))
	}

	return func(eng *Engine) error {
		filter := Number
		var numbers []Unit
		var outType = Integer // We can be optimistic, right?

		nextUnit, err := eng.Pop()

		for err == nil {
			if !filter.Matches(nextUnit.Type) {
				return fmt.Errorf("expected integer or float, got %s", nextUnit.Type.ToString())
			}

			if nextUnit.Type == Float {
				outType = Float
			}

			numbers = append(numbers, nextUnit)

			nextUnit, err = eng.Pop()
		}

		switch outType {
		case Float:
			var out = numAsF(numbers[0])
			for _, num := range numbers[1:] {
				out = floatOp(out, numAsF(num))
			}

			eng.Push(MkFloat(out))
		default:
			var out = numbers[0].Value.(int64)
			for _, num := range numbers[1:] {
				out = intOp(out, num.Value.(int64))
			}

			eng.Push(MkInt(out))
		}

		return nil
	}
}