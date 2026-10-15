package builtin

import (
	"fmt"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func stringSet(eng *Engine) error {
	char := util.Assert(eng.Pop()).Raw.(rune)
	idx := util.Assert(eng.Pop()).Raw.(int64)
	strVarName := util.Assert(eng.Pop()).Raw.(string)
	strVar,err := eng.GetVar(strVarName)

	if err != nil {
		return err
	} else if strVar.TypeString() != "String" {
		return fmt.Errorf("expected 'String', got '%s'", strVar.TypeString())
	}

	asArr := []rune(strVar.Raw.(string))
	if int(idx) > len(asArr) - 1 || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d ", idx, len(asArr))
	}

	asArr[idx] = char

	eng.SetVar(strVarName, StringVal(string(asArr)))

	return nil
}

func stringCreate(eng *Engine) error {
	var strLen = util.Assert(eng.Pop()).Raw.(int64)

	if strLen < 0 {
		return fmt.Errorf("invalid index %d", strLen)
	}

	var tmp = make([]rune, strLen)
	eng.Push(StringVal(string(tmp)))

	return nil
}

func stringConcat(eng *Engine) error {
	var strs []string
	for {
		if val, err := eng.Pop(); err != nil {
			break
		} else {
			strs = append(strs, val.Raw.(string))
		}
	}

	var out string
	idx := len(strs) - 1

	for idx >= 0 {
		out = fmt.Sprintf("%s%s", out, strs[idx])
		idx--
	}

	eng.Push(StringVal(out))
	return nil
}

func stringRef(eng *Engine) error {
	idx := util.Assert(eng.Pop()).Raw.(int64)
	str := util.Assert(eng.Pop()).Raw.(string)

	if int(idx) > len([]rune(str)) || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d", idx, len([]rune(str)))
	}

	eng.Push(CharVal([]rune(str)[idx]))
	return nil
}

func stringize(eng *Engine) error {
	var chars []rune
	for {
		if val, err := eng.Pop(); err != nil {
			break
		} else {
			chars = append(chars, val.Raw.(rune))
		}
	}

	var out []rune
	idx := len(chars) - 1

	for idx >= 0 {
		out = append(out, chars[idx])

		idx--
	}

	eng.Push(StringVal(string(out)))
	return nil
}