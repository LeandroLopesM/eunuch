package builtin

import (
	"fmt"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func vectorRef(eng *Engine) error {
	idx := util.Assert(eng.Pop()).Raw.(int64)
	vec := util.Assert(eng.Pop()).Raw.(Vector)

	if int(idx) > len(vec) || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d", idx, len(vec))
	}

	eng.Push(vec[idx])
	return nil
}

func vectorSet(eng *Engine) error {
	val := util.Assert(eng.Pop())
	idx := util.Assert(eng.Pop()).Raw.(int64)
	vecVarName := util.Assert(eng.Pop()).Raw.(string)
	vecVar,err := eng.GetVar(vecVarName)

	if err != nil {
		return err
	} else if vecVar.TypeString() != "Vector" {
		return fmt.Errorf("expected 'Vector', got '%s'", vecVar.TypeString())
	}

	asArr := vecVar.Raw.(Vector)
	if int(idx) > len(asArr) - 1 || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d ", idx, len(asArr))
	}

	asArr[idx] = val

	eng.SetVar(vecVarName, VectorVal(asArr))

	return nil
}

func vectorCreate(eng *Engine) error {
	var vecLen = util.Assert(eng.Pop()).Raw.(int64)

	if vecLen < 0 {
		return fmt.Errorf("invalid index %d", vecLen)
	}

	var tmp = make(Vector, vecLen)
	eng.Push(VectorVal(tmp))

	return nil
}

func vector(eng *Engine) error {
	var strs Vector
	for {
		if val, err := eng.Pop(); err != nil {
			break
		} else {
			strs = append(strs, val)
		}
	}

	var out Vector
	idx := len(strs) - 1

	for idx >= 0 {
		out = append(out, strs[idx])
		idx--
	}

	eng.Push(VectorVal(out))
	return nil
}

