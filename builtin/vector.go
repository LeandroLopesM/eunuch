package builtin

import (
	"fmt"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func vectorRef(eng *Engine) error {
	idx := util.Assert(eng.Pop()).(Integer)
	vec := util.Assert(eng.Pop()).(Vector)

	if int(idx) > len(vec) || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d", idx, len(vec))
	}

	eng.Push(vec[idx])
	return nil
}

func vectorSet(eng *Engine) error {
	val := util.Assert(eng.Pop())
	idx := util.Assert(eng.Pop()).(Integer)
	vecVarName := string(util.Assert(eng.Pop()).(Symbol))
	vecVar,err := eng.GetVar(vecVarName)

	if err != nil {
		return err
	} else if vecVar.TypeString() != "Vector" {
		return fmt.Errorf("expected 'Vector', got '%s'", vecVar.TypeString())
	}

	asArr := vecVar.(Vector)
	if int(idx) > len(asArr) - 1 || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d ", idx, len(asArr))
	}

	asArr[idx] = val

	eng.SetVar(vecVarName, Vector{asArr})

	return nil
}

func vectorCreate(eng *Engine) error {
	var vecLen = util.Assert(eng.Pop()).(Integer)
	if vecLen < 0 {
		return fmt.Errorf("invalid index %d", vecLen)
	}

	var tmp = make(Vector, vecLen)
	eng.Push(Vector{tmp})

	return nil
}

func vector(eng *Engine) error {
	vals := eng.GenArgv()

	var out Vector

	for _,val := range vals {
		out = append(out, val)
	}

	eng.Push(Vector{out})
	return nil
}

