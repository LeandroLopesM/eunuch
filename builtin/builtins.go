package builtin

import (
	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func RegisterSelf(eng *Engine) {
	eng.AddFunc("display", Builtin{Args: []Type{Any}, VarArgs: true, Return: None, Call: display})
	eng.AddFunc("newline", Builtin{Args: []Type{}, VarArgs: false, Return: None, Call: newline})

	eng.AddFunc("string?", Builtin{Args: []Type{Any}, VarArgs: false, Return: Bool, Call: isX(String)})
	eng.AddFunc("symbol?", Builtin{Args: []Type{Any}, VarArgs: false, Return: Bool, Call: isX(Symbol)})
	eng.AddFunc("boolean?", Builtin{Args: []Type{Any}, VarArgs: false, Return: Bool, Call: isX(Bool)})
	eng.AddFunc("integer?", Builtin{Args: []Type{Any}, VarArgs: false, Return: Bool, Call: isX(Integer)})
	eng.AddFunc("rational?", Builtin{Args: []Type{Any}, VarArgs: false, Return: Bool, Call: isX(Float)})
	eng.AddFunc("eqv", Builtin{Args: []Type{Any}, VarArgs: true, Return: Bool, Call: eqv})

	eng.AddFunc("+", Builtin{Args: []Type{Number}, VarArgs: true, Return: Number, Call: MathOp('+')})
	eng.AddFunc("-", Builtin{Args: []Type{Number}, VarArgs: true, Return: Number, Call: MathOp('-')})
	eng.AddFunc("*", Builtin{Args: []Type{Number}, VarArgs: true, Return: Number, Call: MathOp('*')})
	eng.AddFunc("/", Builtin{Args: []Type{Number}, VarArgs: true, Return: Number, Call: MathOp('/')})
	eng.AddFunc("expt", Builtin{Args: []Type{Number}, VarArgs: true, Return: Number, Call: expt})
	eng.AddFunc("max", Builtin{Args: []Type{Number}, VarArgs: true, Return: Number, Call: OrdOp('>')})
	eng.AddFunc("min", Builtin{Args: []Type{Number}, VarArgs: true, Return: Number, Call: OrdOp('<')})

	eng.AddFunc("string", Builtin{Args: []Type{Char}, VarArgs: true, Return: String, Call: stringize})
	eng.AddFunc("string-ref", Builtin{Args: []Type{String, Integer}, VarArgs: false, Return: Char, Call: stringRef})
	eng.AddFunc("string-append", Builtin{Args: []Type{String}, VarArgs: true, Return: String, Call: stringConcat})
	eng.AddFunc("make-string", Builtin{Args: []Type{Integer}, VarArgs: false, Return: String, Call: stringCreate})            // This shouldnt get used muc}h
	eng.AddFunc("string-set!", Builtin{Args: []Type{Symbol, Integer, Char}, VarArgs: false, Return: String, Call: stringSet}) // This shouldnt get used muc}h

	eng.AddFunc("vector", Builtin{Args: []Type{Any}, VarArgs: true, Return: Vector, Call: vector})
	eng.AddFunc("vector-ref", Builtin{Args: []Type{Vector, Integer}, VarArgs: false, Return: Any, Call: vectorRef})
	eng.AddFunc("make-vector", Builtin{Args: []Type{Integer}, VarArgs: false, Return: Vector, Call: vectorCreate})           // This shouldnt get used muc}h
	eng.AddFunc("vector-set!", Builtin{Args: []Type{Symbol, Integer, Any}, VarArgs: false, Return: Vector, Call: vectorSet}) // This shouldnt get used muc}h

	eng.AddFunc("char=?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charOp("=")})
	eng.AddFunc("char>?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charOp(">")})
	eng.AddFunc("char<?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charOp("<")})
	eng.AddFunc("char>=?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charOp(">=")})
	eng.AddFunc("char<=?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charOp("<=")})

	eng.AddFunc("char-ci=?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charCiOp("=")})
	eng.AddFunc("char-ci>?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charCiOp(">")})
	eng.AddFunc("char-ci<?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charCiOp("<")})
	eng.AddFunc("char-ci>=?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charCiOp(">=")})
	eng.AddFunc("char-ci<=?", Builtin{Args: []Type{Char}, VarArgs: false, Return: Bool, Call: charCiOp("<=")})

	eng.AddFunc("car", Builtin{Args: []Type{Pair}, VarArgs: false, Return: Any, Call: pairGet(0)})
	eng.AddFunc("cdr", Builtin{Args: []Type{Pair}, VarArgs: false, Return: Any, Call: pairGet(1)})
	eng.AddFunc("set-car!", Builtin{Args: []Type{Symbol}, VarArgs: false, Return: None, Call: pairSet(0)})
	eng.AddFunc("set-cdr!", Builtin{Args: []Type{Symbol}, VarArgs: false, Return: None, Call: pairSet(1)})
	eng.AddFunc("cons", Builtin{Args: []Type{Any, Any}, VarArgs: false, Return: Pair, Call: newPair})

	eng.AddFunc(
		"define",
		Builtin{
			Args:    []Type{Symbol, Any},
			VarArgs: false,
			Return:  None,
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
			Args:    []Type{Any},
			VarArgs: false,
			Return:  Any,
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

	eng.SetVar(name.Value.(string), value)
	return nil
}

func isX(which Type) BuiltinExec {
	return func(eng *Engine) error {
		if val, err := eng.Pop(); err != nil {
			return err
		} else {
			eng.Push(MkBool(val.Type == which))
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

	if lhs.Type != rhs.Type {
		eng.Push(MkBool(false))
	}

	eng.Push(MkBool(lhs == rhs))
	return nil
}

func display(eng *Engine) error {
	var args []Unit

	for {
		if val, err := eng.Pop(); err != nil {
			break
		} else {
			args = append(args, val)
		}
	}

	idx := len(args) - 1

	for idx >= 0 {
		print(SprintUnit(args[idx]))

		idx--
	}

	return nil
}
