package builtin

import (
	"fmt"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func vectorRef(e *Engine) error {
	idx := util.Assert(e.Pop()).Value.(int64)
	vec := util.Assert(e.Pop()).Value.(VectorVal)

	if int(idx) > len(vec) || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d", idx, len(vec))
	}

	e.Push(vec[idx])
	return nil
}

func vectorSet(e *Engine) error {
	val := util.Assert(e.Pop())
	idx := util.Assert(e.Pop()).Value.(int64)
	vecVarName := util.Assert(e.Pop()).Value.(string)
	vecVar,err := e.GetVar(vecVarName)

	if err != nil {
		return err
	} else if vecVar.Type != Vector {
		return fmt.Errorf("expected 'Vector', got '%s'", TypeNames[vecVar.Type])
	}

	asArr := vecVar.Value.(VectorVal)
	if int(idx) > len(asArr) - 1 || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d ", idx, len(asArr))
	}

	asArr[idx] = val

	e.SetVar(vecVarName, MkVector(asArr))

	return nil
}

func vectorCreate(e *Engine) error {
	var len = util.Assert(e.Pop()).Value.(int64)

	if len < 0 {
		return fmt.Errorf("invalid index %d", len)
	}

	var tmp = make(VectorVal, len)
	e.Push(MkVector(tmp))

	return nil
}

func vector(e *Engine) error {
	var strs VectorVal
	for {
		if val, err := e.Pop(); err != nil {
			break
		} else {
			strs = append(strs, val)
		}
	}

	var out VectorVal
	idx := len(strs) - 1

	for idx >= 0 {
		out = append(out, strs[idx])
		idx--
	}

	e.Push(MkVector(out))
	return nil
}

