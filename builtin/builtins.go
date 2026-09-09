package builtin

import (
	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func RegisterSelf(eng *Engine) {
	eng.AddFunc("display", []Type{Any}, true, None, display)
	eng.AddFunc("newline", []Type{}, false, None, newline)

	eng.AddFunc("string?", []Type{Any}, false, Bool, isX(String))
	eng.AddFunc("symbol?", []Type{Any}, false, Bool, isX(Symbol))
	eng.AddFunc("boolean?", []Type{Any}, false, Bool, isX(Bool))
	eng.AddFunc("integer?", []Type{Any}, false, Bool, isX(Integer))
	eng.AddFunc("rational?", []Type{Any}, false, Bool, isX(Float))
	
	eng.AddFunc("+", []Type{Number}, true, Number, MathOp('+'))
	eng.AddFunc("-", []Type{Number}, true, Number, MathOp('-'))
	eng.AddFunc("*", []Type{Number}, true, Number, MathOp('*'))
	eng.AddFunc("/", []Type{Number}, true, Number, MathOp('/'))
	eng.AddFunc("expt", []Type{Number}, true, Number, expt)

	eng.AddFunc("max", []Type{Number}, true, Number, OrdOp('>'))
	eng.AddFunc("min", []Type{Number}, true, Number, OrdOp('<'))


	eng.AddFunc("eqv", []Type{Any}, true, Bool, eqv)
	
	eng.AddFunc("define", []Type{Symbol, Any}, false, None, define)
	
	eng.AddFunc("string", []Type{Char}, true, String, stringize)
	eng.AddFunc("string-ref", []Type{String, Integer}, false, Char, stringRef)
	eng.AddFunc("string-append", []Type{String}, true, String, stringConcat)
	eng.AddFunc("make-string", []Type{Integer}, false, String, stringCreate) // This shouldnt get used much
	eng.AddFunc("string-set!", []Type{Symbol, Integer, Char}, false, String, stringSet) // This shouldnt get used much
	
	eng.AddFunc("vector", []Type{Any}, true, Vector, vector)
	eng.AddFunc("vector-ref", []Type{Vector, Integer}, false, Any, vector)
	eng.AddFunc("make-vector", []Type{Integer}, false, Vector, vectorCreate) // This shouldnt get used much
	eng.AddFunc("vector-set!", []Type{Symbol, Integer, Any}, false, Vector, vectorSet) // This shouldnt get used much
	
	eng.AddFunc("char=?", []Type{Char}, false, Bool, charOp("="))
	eng.AddFunc("char>?", []Type{Char}, false, Bool, charOp(">"))
	eng.AddFunc("char<?", []Type{Char}, false, Bool, charOp("<"))
	eng.AddFunc("char>=?", []Type{Char}, false, Bool, charOp(">="))
	eng.AddFunc("char<=?", []Type{Char}, false, Bool, charOp("<="))
	
	eng.AddFunc("char-ci=?", []Type{Char}, false, Bool, charCiOp("="))
	eng.AddFunc("char-ci>?", []Type{Char}, false, Bool, charCiOp(">"))
	eng.AddFunc("char-ci<?", []Type{Char}, false, Bool, charCiOp("<"))
	eng.AddFunc("char-ci>=?", []Type{Char}, false, Bool, charCiOp(">="))
	eng.AddFunc("char-ci<=?", []Type{Char}, false, Bool, charCiOp("<="))
	
	eng.AddFunc("car", []Type{Pair}, false, Any, pairGet(0))
	eng.AddFunc("cdr", []Type{Pair}, false, Any, pairGet(1))
	eng.AddFunc("set-car!", []Type{Symbol}, false, None, pairSet(0))
	eng.AddFunc("set-cdr!", []Type{Symbol}, false, None, pairSet(1))
	eng.AddFunc("cons", []Type{Any, Any}, false, Pair, newPair)

	eng.AddFunc("quote", []Type{Any}, false, Any, quote)
}

func quote(e *Engine) error {
	e.Push(util.Assert(e.Pop()))

	return nil
}

func define(e *Engine) error {
	value := util.Assert(e.Pop())
	name := util.Assert(e.Pop())
	
	e.SetVar(name.Value.(string), value)
	return nil
}

func isX(which Type) BuiltinExec {
	return func(e *Engine) error {
		if val, err := e.Pop(); e != nil {
			return err
		} else {
			e.Push(MkBool(val.Type == which))
		}

		return nil
	}
}

func newline(e *Engine) error {
	print("\n")

	return nil
}

func eqv(e *Engine) error {
	lhs := util.Assert(e.Pop()) // We can assert because the argCount was checked
	rhs := util.Assert(e.Pop())

	if lhs.Type != rhs.Type {
		e.Push(MkBool(false))
	}

	e.Push(MkBool(lhs == rhs))
	return nil
}

func display(e *Engine) error {
	var args []Unit

	for {
		if val, err := e.Pop(); err != nil {
			break
		} else {
			args = append(args, val)
		}
	}

	idx := len(args) - 1;

	for idx >= 0 {
		print(SprintUnit(args[idx]))

		idx--
	}

	return nil
}
