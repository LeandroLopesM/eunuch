package builtin

import (
	"reflect"
	"slices"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func RegisterSelf(eng *Engine) {
	eng.AddFunc("display", Builtin{VarArgs: true, Call: display})
	eng.AddFunc("newline", Builtin{Argc: 0, Call: newline})

	eng.AddFunc("string?", Builtin{Argc: 1, Call: isX[String]()})
	eng.AddFunc("symbol?", Builtin{Argc: 1, Call: isX[Symbol]()})
	eng.AddFunc("boolean?", Builtin{Argc: 1, Call: isX[Bool]()})
	eng.AddFunc("integer?", Builtin{Argc: 1, Call: isX[Integer]()})
	eng.AddFunc("rational?", Builtin{Argc: 1, Call: isX[Float]()})
	eng.AddFunc("eqv", Builtin{VarArgs: true, Call: eqv})

	eng.AddFunc("+", Builtin{VarArgs: true, Call: MathOp('+')})
	eng.AddFunc("-", Builtin{VarArgs: true, Call: MathOp('-')})
	eng.AddFunc("*", Builtin{VarArgs: true, Call: MathOp('*')})
	eng.AddFunc("/", Builtin{VarArgs: true, Call: MathOp('/')})
	eng.AddFunc("expt", Builtin{VarArgs: true, Call: expt})

	eng.AddFunc("max", Builtin{VarArgs: true, Call: OrdOp('>')})
	eng.AddFunc("min", Builtin{VarArgs: true, Call: OrdOp('<')})

	eng.AddFunc("string", Builtin{VarArgs: true, Call: stringize})
	eng.AddFunc("string-ref", Builtin{Argc: 2, Call: stringRef})
	eng.AddFunc("string-append", Builtin{VarArgs: true, Call: stringConcat})
	eng.AddFunc("make-string", Builtin{Argc: 1, Call: stringCreate})            // This shouldnt get used muc}h
	eng.AddFunc("string-set!", Builtin{Argc: 3, Call: stringSet})

	eng.AddFunc("vector", Builtin{VarArgs: true, Call: vector})
	eng.AddFunc("vector-ref", Builtin{Argc: 2, Call: vectorRef})
	eng.AddFunc("make-vector", Builtin{Argc: 1, Call: vectorCreate})           // This shouldnt get used muc}h
	eng.AddFunc("vector-set!", Builtin{Argc: 3, Call: vectorSet}) // This shouldnt get used muc}h

	eng.AddFunc("char=?", Builtin{Argc: 1, Call: charOp("=")})
	eng.AddFunc("char>?", Builtin{Argc: 1, Call: charOp(">")})
	eng.AddFunc("char<?", Builtin{Argc: 1, Call: charOp("<")})
	eng.AddFunc("char>=?", Builtin{Argc: 1, Call: charOp(">=")})
	eng.AddFunc("char<=?", Builtin{Argc: 1, Call: charOp("<=")})

	eng.AddFunc("char-ci=?", Builtin{Argc: 1, Call: charCiOp("=")})
	eng.AddFunc("char-ci>?", Builtin{Argc: 1, Call: charCiOp(">")})
	eng.AddFunc("char-ci<?", Builtin{Argc: 1, Call: charCiOp("<")})
	eng.AddFunc("char-ci>=?", Builtin{Argc: 1, Call: charCiOp(">=")})
	eng.AddFunc("char-ci<=?", Builtin{Argc: 1, Call: charCiOp("<=")})

	eng.AddFunc("car", Builtin{Argc: 1, Call: pairGet(0)})
	eng.AddFunc("cdr", Builtin{Argc: 1, Call: pairGet(1)})
	eng.AddFunc("set-car!", Builtin{Argc: 1, Call: pairSet(0)})
	eng.AddFunc("set-cdr!", Builtin{Argc: 1, Call: pairSet(1)})
	eng.AddFunc("cons", Builtin{Argc: 2, Call: newPair})

	eng.AddFunc(
		"define",
		Builtin{
			Argc: 2,
			Call:    define,
			Prepare: util.Some[StagingFunc](
				func(s Scheme, eng *Engine) error {
					params := s.Params()
					eng.Push(params[0])
					return eng.Evaluate(params[1])
				}),
		})

	eng.AddFunc(
		"quote",
		Builtin{
			Argc: 1,
			Call:    quote,
			Prepare: util.Some[StagingFunc](
				func(s Scheme, eng *Engine) error {
					// Don't evaluate the argument,
					// see quote() for why.

					eng.Push(s.Args[1])
					return nil
				}),
		},
	)
}

func genArgv(e *Engine) []Value {
	var vals []Value
	for {
		if val,err := e.Pop(); err != nil {
			break
		} else {
			vals = append(vals, val)
		}
	}

	slices.Reverse(vals)
	return vals
}

/*
 *	As i understand, (quote) essentially
 *	defers evaluation of an expression
 *	until the related symbol is used.
 *	For example:
 *		(define x '(+ 1 y))
 *		(define y 3)
 *		(display x) // 1 + 3 = 4
 *	It would fail to evaluate 'y' but,
 *	since it won't be evaluated until it
 *	is used, it doesn't fail.
 *
 *	!! This is how i interpreted the
 *	spec and, therefore, is open to critique.
 */
func quote(eng *Engine) error {
	eng.Push(util.Assert(eng.Pop()))
	return nil
}

func define(eng *Engine) error {
	value := util.Assert(eng.Pop())
	name := util.Assert(eng.Pop())

	eng.SetVar(name.Raw.(string), value)
	return nil
}

func isX[T any]() BuiltinExec {
	return func(eng *Engine) error {
		if val, err := eng.Pop(); err != nil {
			return err
		} else {
			eng.Push(BoolVal(reflect.TypeFor[T]() == reflect.TypeOf(val.Raw)))
		}

		return nil
	}
}

func newline(eng *Engine) error {
	print("\n")

	return nil
}

func eqv(eng *Engine) error {
	lhs := util.Assert(eng.Pop()) // We can assert because the argCount was checked
	rhs := util.Assert(eng.Pop())

	if lhs.Type() != rhs.Type() {
		eng.Push(BoolVal(false))
	}

	eng.Push(BoolVal(lhs == rhs))
	return nil
}

func display(eng *Engine) error {
	args := genArgv(eng)
	for _,arg := range args {
		print(arg.ToString())
	}

	return nil
}
