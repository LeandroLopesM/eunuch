package builtin

import (
	"fmt"

	. "github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"
)

func stringSet(eng *Engine) error {
	char := util.Assert(eng.Pop()).(Char)
	idx := util.Assert(eng.Pop()).(Integer)
	strVarName := util.Assert(eng.Pop()).ToString()
	strVar,err := eng.GetVar(strVarName)

	if err != nil {
		return err
	} else if strVar.TypeString() != "String" {
		return fmt.Errorf("expected String, got '%s'", Format(strVar))
	}

	asArr := []rune(strVar.ToString())
	if int(idx) > len(asArr) - 1 || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d ", idx, len(asArr))
	}

	asArr[idx] = rune(char)
	eng.SetVar(strVarName, String(asArr))

	return nil
}

func stringCreate(eng *Engine) error {
	var strLen = util.Assert(eng.Pop()).(Integer)

	if strLen < 0 {
		return fmt.Errorf("invalid index %d", strLen)
	}

	var tmp = make([]rune, strLen)
	eng.Push(String(string(tmp)))

	return nil
}

func stringConcat(eng *Engine) error {
	var strs []string

	for _,arg := range eng.GenArgv() {
		eng.Evaluate(arg)
		if val, err := eng.Pop(); err == nil {
			strs = append(strs, val.ToString())
		}
	}

	var out string
	idx := len(strs) - 1

	for idx >= 0 {
		out = fmt.Sprintf("%s%s", out, strs[idx])
		idx--
	}

	eng.Push(String(out))
	return nil
}

func stringRef(eng *Engine) error {
	idx := util.Assert(eng.Pop()).(Integer)
	str := util.Assert(eng.Pop()).(String)

	if int(idx) > len([]rune(str)) || idx < 0 {
		return fmt.Errorf("index %d out of bounds for %d", idx, len([]rune(str)))
	}

	eng.Push(Char([]rune(str)[idx]))
	return nil
}

func stringize(eng *Engine) error {
	var chars []rune
	for {
		if val, err := eng.Pop(); err != nil {
			break
		} else {
			chars = append(chars, rune(val.(Char)))
		}
	}

	var out []rune
	idx := len(chars) - 1

	for idx >= 0 {
		out = append(out, chars[idx])

		idx--
	}

	eng.Push(String(string(out)))
	return nil
}