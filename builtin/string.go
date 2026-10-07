package builtin

import (
	"fmt"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func stringSet(e *Engine) error {
	char := util.Assert(e.Pop()).Value.(rune)
	idx := util.Assert(e.Pop()).Value.(int64)
	strVarName := util.Assert(e.Pop()).Value.(string)
	strVar,err := e.GetVar(strVarName)

	if err != nil {
		return err
	} else if strVar.Type != String {
		return fmt.Errorf("expected 'String', got '%s'", TypeNames[strVar.Type])
	}

	asArr := []rune(strVar.Value.(string))
	if int(idx) > len(asArr) - 1 || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d ", idx, len(asArr))
	}

	asArr[idx] = char

	e.SetVar(strVarName, MkString(string(asArr)))

	return nil
}

func stringCreate(e *Engine) error {
	var len = util.Assert(e.Pop()).Value.(int64)

	if len < 0 {
		return fmt.Errorf("invalid index %d", len)
	}

	var tmp = make([]rune, len)
	e.Push(MkString(string(tmp)))

	return nil
}

func stringConcat(e *Engine) error {
	var strs []string
	for {
		if val, err := e.Pop(); err != nil {
			break
		} else {
			strs = append(strs, val.Value.(string))
		}
	}

	var out string
	idx := len(strs) - 1

	for idx >= 0 {
		out = fmt.Sprintf("%s%s", out, strs[idx])
		idx--
	}

	e.Push(MkString(out))
	return nil
}

func stringRef(e *Engine) error {
	idx := util.Assert(e.Pop()).Value.(int64)
	str := util.Assert(e.Pop()).Value.(string)

	if int(idx) > len([]rune(str)) || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d", idx, len([]rune(str)))
	}

	e.Push(MkChar([]rune(str)[idx]))
	return nil
}

func stringize(e *Engine) error {
	var chars []rune
	for {
		if val, err := e.Pop(); err != nil {
			break
		} else {
			chars = append(chars, val.Value.(rune))
		}
	}

	var out []rune
	idx := len(chars) - 1

	for idx >= 0 {
		out = append(out, chars[idx])

		idx--
	}

	e.Push(MkString(string(out)))
	return nil
}