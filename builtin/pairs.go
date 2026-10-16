package builtin

import (
	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func pairGet(idx int) BuiltinExec {
	return func (eng *Engine) error {
		pair := util.Assert(eng.Pop()).(Pair)

		eng.Push(pair[idx])
		return nil
	}
}

func pairSet(idx int) BuiltinExec {
	return func (eng *Engine) error {
		pairVarName := string(util.Assert(eng.Pop()).(Symbol))
		pairVarVal := util.Assert(eng.GetVar(pairVarName)).(Pair)

		newVal := util.Assert(eng.Pop())
		pairVarVal[idx] = newVal

		eng.SetVar(pairVarName, Pair{pairVarVal})
		return nil
	}
}

func newPair(eng *Engine) error {
	cdr := util.Assert(eng.Pop())
	car := util.Assert(eng.Pop())

	eng.Push(Pair{ car, cdr })
	return nil
}