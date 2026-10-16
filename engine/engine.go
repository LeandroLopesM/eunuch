package engine

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/leandrolopesm/eunuch/core"
	. "github.com/leandrolopesm/eunuch/core"
	"github.com/leandrolopesm/eunuch/parser"
	"github.com/leandrolopesm/eunuch/util"
	"github.com/logrusorgru/aurora/v4"
)

type BuiltinExec func(e *Engine) error
type StagingFunc func(Scheme, *Engine) error

type Builtin struct {
	Argc    int
	VarArgs bool // Is the function var_arg

	Prepare util.Option[StagingFunc] // Custom function for checking parameters

	Call BuiltinExec
}

type Engine struct {
	stack        Stack[Value]
	stackHistory Stack[int] // Defines the lower bounds for the current stackPtr

	currentFn util.Option[Scheme]

	vars  map[string](Value)
	funcs map[string](Builtin)
}

func New() Engine {
	return Engine{
		stack:        NewStack[Value](),
		stackHistory: NewStack[int](),

		currentFn: util.None[Scheme](),

		vars:  make(map[string]Value),
		funcs: make(map[string]Builtin),
	}
}

func formatStackTrace(strace error) error {
	calls := strings.Split(strace.Error(), "|")
	var out  = strings.Builder {}
	var idx = 0

	if core.ErrorStyle.Has(core.StackTrace) {
		for range calls[:len(calls)-1] {
			out.WriteString(fmt.Sprintf("%s%s\n", strings.Repeat(". ", idx+1), calls[idx]))
			idx++
		}

		idx++
	}

	var errMsg string
	if core.ErrorStyle.Has(core.Colored) {
		errMsg = aurora.Red(calls[len(calls)-1]).String()
	} else {
		errMsg = calls[len(calls)-1]
	}

	out.WriteString(fmt.Sprintf("%s%s", strings.Repeat(". ", idx), errMsg))

	return errors.New(out.String())
}

func (eng *Engine) GenArgv() []Value {
	var vals []Value
	for {
		if val,err := eng.Pop(); err != nil {
			break
		} else {
			vals = append(vals, val)
		}
	}

	slices.Reverse(vals)
	return vals
}

func (eng *Engine) ExecuteFile(file string) error {
	if v, e := os.ReadFile(file); e != nil {
		return e
	} else {
		return eng.ExecuteStr(file, string(v))
	}
}

func printAst(AST []Value, depth int) {
	for _, val := range AST {
		log.Debugf("%s> %s", strings.Repeat(". ", depth), Format(val))

		if args,ok := val.(Scheme); ok {
			printAst(args.Params(), depth + 1)
		}
	}
}

func (eng *Engine) ExecuteStr(file, code string) error {
	fileAST, err := parser.ParseFile(file, code)

	if err != nil {
		return err
	}

	if log.GetLevel() == log.DebugLevel {
		printAst(fileAST, 0)
	}

	for _, expr := range fileAST {
		if err := eng.Evaluate(expr); err != nil {
			return formatStackTrace(err)
		}
	}

	return nil
}

func (eng *Engine) Evaluate(rawValue Value) error {
	switch typedVal := rawValue.(type) {
	case Scheme:
		return eng.runScheme(typedVal)
	case Symbol:
		if typedVal, err := eng.GetVar(string(typedVal)); err != nil {
			return err
		} else if err := eng.Evaluate(typedVal); err != nil {
			return err
		}
	default:
		eng.Push(rawValue)
	}

	return nil
}

func (eng *Engine) GetVar(name string) (Value, error) {
	if v, ok := eng.vars[name]; !ok {
		return Null, fmt.Errorf("undefined variable '%s'", name)
	} else {
		return v, nil
	}
}

func (eng *Engine) SetVar(name string, val Value) {
	eng.vars[name] = val
}

func (eng *Engine) checkScheme(scheme Scheme) error {
	var actualFn Builtin
	if fn, ok := eng.funcs[scheme.Name()]; !ok {
		return fmt.Errorf("undefined function '%s'", scheme.Name())
	} else {
		actualFn = fn
	}

	if len(scheme.Params()) != actualFn.Argc && !actualFn.VarArgs {
		return fmt.Errorf("scheme '%s': Expected %d args, got %d", scheme.Name(), actualFn.Argc, len(scheme.Params()))
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
	return fmt.Errorf("%s (%s):|%w", currScheme.Position.ToString(), currScheme.Name(), err)
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

func (eng *Engine) Pop() (Value, error) {
	return eng.stack.Pop()
}

func (eng *Engine) Push(v Value) {
	eng.stack.Push(v)
}
