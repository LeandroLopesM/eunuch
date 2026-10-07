package builtin

import (
	"github.com/charmbracelet/log"
	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func pairGet(idx int) BuiltinExec {
	return func (eng *Engine) error {
		pair := util.Assert(eng.Pop()).Value.(PairVal)

		eng.Push(pair[idx])
		return nil
	}
}

func pairSet(idx int) BuiltinExec {
	return func (eng *Engine) error {
		pairVarName := util.Assert(eng.Pop()).Value.(string)
		pairVarVal := util.Assert(eng.GetVar(pairVarName)).Value.(PairVal)

		newVal := util.Assert(eng.Pop())
		pairVarVal[idx] = newVal

		eng.SetVar(pairVarName, MkPair(pairVarVal))
		return nil
	}
}

func newPair(eng *Engine) error {
	cdr := util.Assert(eng.Pop())
	car := util.Assert(eng.Pop())

	log.Debugf("New pair %s", SprintUnit(MkPair(PairVal{car, cdr})));

	eng.Push(MkPair(PairVal{ car, cdr }))

	return nil
}