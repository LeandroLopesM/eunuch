package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	. "github.com/leandrolopesm/eunuch/core"
	"github.com/leandrolopesm/eunuch/parser"
	"github.com/leandrolopesm/eunuch/util"
	"github.com/logrusorgru/aurora/v4"
)

type BuiltinExec func(e *Engine) error
type StagingFunc func(Scheme, *Engine) error

type Builtin struct {
	Args    []Type
	VarArgs bool // Is the function var_arg

	Return Type

	Prepare util.Option[StagingFunc] // Custom function for checking parameters

	Call BuiltinExec
}

type Engine struct {
	file string

	stack        Stack[Unit]
	stackHistory Stack[int] // Defines the lower bounds for the current stackPtr

	currentFn util.Option[Scheme]

	vars  map[string](Unit)
	funcs map[string](Builtin)
}

func New() Engine {
	ret := Engine{
		file: "#ENGINE",

		stack:        NewStack[Unit](),
		stackHistory: NewStack[int](),

		currentFn: util.None[Scheme](),

		vars:  make(map[string]Unit),
		funcs: make(map[string]Builtin),
	}

	return ret
}

func formatStackTrace(strace error) error {
	calls := strings.Split(strace.Error(), "|")
	var out string

	var idx = 0
	for range calls[:len(calls)-1] {
		out += fmt.Sprintf("%s%s\n", strings.Repeat(". ", idx+1), calls[idx])
		idx++
	}

	out += fmt.Sprintf("%s%s", strings.Repeat(". ", idx+1), aurora.Red(calls[len(calls)-1]))

	return errors.New(out)
}

func (eng *Engine) ExecuteFile(file string) error {
	if v, e := os.ReadFile(file); e != nil {
		return e
	} else {
		eng.file = file
		return eng.ExecuteStr(string(v))
	}
}

func (eng *Engine) ExecuteStr(code string) error {
	if eng.file == "#ENGINE" { // If this wasn't called by ExecuteFile
		eng.file = "<anonymous>"
	}

	schemes, err := parser.ParseFile(eng.file, code)

	if err != nil {
		return err
	}

	if log.GetLevel() == log.DebugLevel {
		for _, scheme := range schemes {
			DebugUnit(scheme)
		}
	}

	for _, expr := range schemes {
		if err := eng.Evaluate(expr); err != nil {
			return formatStackTrace(err)
		}
	}

	return nil
}

func (eng *Engine) Evaluate(unit Unit) error {
	switch unit.Type {
	case SchemeType:
		return eng.runScheme(unit.Value.(Scheme))
	case Symbol:
		if val, err := eng.GetVar(unit.Value.(string)); err != nil {
			return err
		} else if err := eng.Evaluate(val); err != nil {
			return err
		}
	default:
		eng.Push(unit)
	}

	return nil
}

func (eng *Engine) GetVar(name string) (Unit, error) {
	if v, ok := eng.vars[name]; !ok {
		return Unit{}, fmt.Errorf("undefined variable '%s'", name)
	} else {
		return v, nil
	}
}

func (eng *Engine) SetVar(name string, val Unit) {
	eng.vars[name] = val
}

func (eng *Engine) checkScheme(scheme Scheme) error {
	var actualFn Builtin
	if fn, ok := eng.funcs[scheme.Name()]; !ok {
		return fmt.Errorf("undefined function '%s'", scheme.Name())
	} else {
		actualFn = fn
	}

	schemeParams := scheme.Params()

	if len(schemeParams) != len(actualFn.Args) && !actualFn.VarArgs {
		return fmt.Errorf("scheme '%s': Expected %d args, got %d", scheme.Name(), len(actualFn.Args), len(schemeParams))
	}

	idx := 0
	for range actualFn.Args {
		inType := schemeParams[idx].Type

		if schemeParams[idx].Type == SchemeType {
			asScheme := schemeParams[idx].Value.(Scheme)
			if e := eng.checkScheme(asScheme); e != nil {
				return e
			}

			// If the function exists, use it's return value as the type
			inType = eng.funcs[asScheme.Name()].Return

			// If the passed argument is a symbol and we dont want a symbol, get its actual type
		} else if schemeParams[idx].Type == Symbol && !actualFn.Args[idx].Matches(Symbol) {
			if v, e := eng.GetVar(schemeParams[idx].Value.(string)); e != nil {
				return e
			} else {
				inType = v.Type
			}
		}

		if !actualFn.Args[idx].Matches(inType) {
			return fmt.Errorf(
				"incorrect argument type for '%s'. Expected '%s' got '%s'",
				scheme.Name(),
				actualFn.Args[idx].ToString(),
				inType.ToString(),
			)
		}

		if !actualFn.VarArgs { // We only match against the first arg, repeating
			idx++
		}
	}

	return nil
}

func (eng *Engine) error(msg string, args ...any) error {
	if eng.currentFn.IsSome() {
		scheme := eng.currentFn.Unwrap()
		return fmt.Errorf("%s (%s): %s", scheme.Position.ToString(), scheme.Name(), fmt.Sprintf(msg, args...))
	}
	return fmt.Errorf("%s", aurora.Red(fmt.Sprintf(msg, args...)))
}

func (eng *Engine) schemeError(currScheme Scheme, err error) error {
	return fmt.Errorf("%s (%s):|%s", currScheme.Position.ToString(), currScheme.Name(), err)
}

func (eng *Engine) runScheme(scheme Scheme) error {
	if err := eng.checkScheme(scheme); err != nil {
		return eng.error("%s", err)
	}

	eng.currentFn = util.Some(scheme)

	// Function must exist (Already checked with checkScheme)
	actualFn := eng.funcs[scheme.Name()]

	eng.saveStack()
	defer eng.loadStack()

	if prepFun, err := actualFn.Prepare.Try(); err == nil {
		if err := prepFun(scheme, eng); err != nil {
			return eng.schemeError(scheme, err)
		}
	} else {
		for _, arg := range scheme.Params() {
			if err := eng.Evaluate(arg); err != nil {
				return eng.schemeError(scheme, err)
			}
		}
	}

	if ret := actualFn.Call(eng); ret != nil {
		return eng.schemeError(scheme, ret)
	}

	eng.currentFn = util.None[Scheme]()

	return nil
}

func (eng *Engine) AddFunc(name string, function Builtin) {
	for k := range eng.funcs {
		if k == name {
			log.Warnf("Attempt to redeclare function %s", name)
			break
		}
	}

	eng.funcs[name] = function
}

func (eng *Engine) Pop() (Unit, error) {
	return eng.stack.Pop()
}

func (eng *Engine) Push(v Unit) {
	eng.stack.Push(v)
}
