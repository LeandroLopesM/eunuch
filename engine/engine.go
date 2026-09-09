package engine

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	. "github.com/leandrolopesm/eunuch/core"
	"github.com/leandrolopesm/eunuch/parser"
	"github.com/logrusorgru/aurora/v4"
)

type BuiltinExec func(e *Engine) error
type StagingFunc func(Scheme, *Engine) error

type Builtin struct {
	Args    []Type
	VarArgs bool // Is the function var_arg

	Return Type

	Prepare StagingFunc // Custom function for checking parameters
	MustPrepare bool

	Call BuiltinExec
}

type Engine struct {
	file string

	stack        Stack[Unit]
	stackHistory Stack[int] // Defines the lower bounds for the current stackPtr

	vars  map[string](Unit)
	funcs map[string](Builtin)
}

func New() Engine {
	ret := Engine{
		file: "#ENGINE",

		stack:        NewStack[Unit](),
		stackHistory: NewStack[int](),

		vars:  make(map[string]Unit),
		funcs: make(map[string]Builtin),
	}
	
	return ret
}

func formatStackTrace(strace error) error {
	calls := strings.Split(strace.Error(), "|")
	var out string

	for i := range calls {
		out = fmt.Sprintf("%s%s%s\n", out, strings.Repeat(". ", i+1), calls[i])
	}

	return errors.New(out)
}

func (self *Engine) ExecuteFile(file string) error {
	if v, e := os.ReadFile(file); e != nil {
		return e
	} else {
		self.file = file
		return self.ExecuteStr(string(v))
	}
}

func (self *Engine) ExecuteStr(code string) error {
	if self.file == "#ENGINE" { // If this wasn't called by ExecuteFile
		self.file = "<anonymous>"
	}

	schemes, err := parser.Lex(self.file, code)

	if err != nil {
		return err
	}

	if log.GetLevel() == log.DebugLevel {
		for _, scheme := range schemes {
			DebugUnit(scheme)
		}
	}

	for _, expr := range schemes {
		if err := self.evalExpr(expr); err != nil {
			return formatStackTrace(err)
		}
	}

	return nil
}

func (self *Engine) evalExpr(unit Unit) error {
	switch unit.Type {
	case SchemeType: return self.runScheme(unit.Value.(Scheme))
	case Symbol:
		if val, err := self.GetVar(unit.Value.(string)); err != nil {
			return err
		} else {
			self.Push(val)
		}
	case Integer, Float , Bool , String , Char, Vector, Pair:
		self.Push(unit)
	default:
		panic(fmt.Sprintf("Unknown unit type %d", unit.Type))
	}

	return nil
}

func (self *Engine) GetVar(name string) (Unit, error) {
	if v,ok := self.vars[name]; !ok {
		return Unit{}, fmt.Errorf("Undefined variable '%s'", name)
	} else {
		return v, nil
	}
}

func (self *Engine) SetVar(name string, val Unit) {
	self.vars[name] = val;
}

func (self *Engine) checkScheme(scheme Scheme) error {
	var actualFn Builtin
	if fn, ok := self.funcs[scheme.Name]; !ok {
		return fmt.Errorf("Undefined function '%s'", scheme.Name)
	} else {
		actualFn = fn
	}

	if len(scheme.Args) != len(actualFn.Args) && !actualFn.VarArgs {
		return fmt.Errorf("Scheme '%s': Expected %d args, got %d", scheme.Name, len(actualFn.Args), len(scheme.Args))
	}

	idx := 0
	for range actualFn.Args {
		inType := scheme.Args[idx].Type

		if scheme.Args[idx].Type == SchemeType {
			asScheme := scheme.Args[idx].Value.(Scheme)
			if e := self.checkScheme(asScheme); e != nil {
				return e
			}

			// If the function exists, use it's return value as the type
			inType = self.funcs[asScheme.Name].Return
		
		// If the passed argument is a symbol and we dont want a symbol, get its actual type
		} else if scheme.Args[idx].Type == Symbol && !actualFn.Args[idx].Matches(Symbol) {
			if v,e := self.GetVar(scheme.Args[idx].Value.(string)); e != nil {
				return e
			} else {
				inType = v.Type
			}
		}

		if !actualFn.Args[idx].Matches(inType) {
			return fmt.Errorf(
				"Incorrect argument type for '%s'. Expected '%s' got '%s'",
				scheme.Name,
				TypeNames[actualFn.Args[idx]],
				TypeNames[inType],
			)
		}

		if !actualFn.VarArgs { // We only match against the first arg, repeating
			idx++
		}
	}

	return nil
}

func (self *Engine) runScheme(scheme Scheme) error {
	if err := self.checkScheme(scheme); err != nil {
		return err
	}

	fn := self.funcs[scheme.Name] // Function must exist (Already checked with checkScheme)
	self.saveStack()

	if fn.MustPrepare {
		if err := fn.Prepare(scheme, self); err != nil {
			return fmt.Errorf("%s %s: %s", scheme.Position.ToString(), scheme.Name, aurora.Red(err))
		}
	} else {
		var i int
		for _, arg := range scheme.Args {
			self.evalExpr(arg)
	
			if !fn.VarArgs {
				i++
			}
		}
	}


	ret := fn.Call(self)

	self.loadStack()
	if ret != nil {
		return fmt.Errorf("%s %s: %s", scheme.Position.ToString(), scheme.Name, aurora.Red(ret))
	}

	return nil
}

func (self *Engine) AddFunc(name string, fn Builtin) error {
	for k := range self.funcs {
		if k == name {
			return fmt.Errorf("Attemt to redeclare function %s", name)
		}
	}

	self.funcs[name] = fn;

	return nil
}

func (self *Engine) Pop() (Unit, error) {
	return self.stack.Pop()
}

func (self *Engine) Push(v Unit) {
	self.stack.Push(v)
}
