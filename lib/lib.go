package lib

import (
	"github.com/leandrolopesm/eunuch/core"
	"github.com/leandrolopesm/eunuch/engine"
)

/* The primary library API
 * for eunuch.
 */

// Core API.

type Symbol = core.Symbol
type Integer = core.Integer
type Float = core.Float
type Bool = core.Bool
type String = core.String
type Char = core.Char
type Vector = core.Vector
type Pair = core.Pair
type Scheme = core.Scheme

var SymbolVal = core.SymbolVal
var IntegerVal = core.IntegerVal
var FloatVal = core.FloatVal
var BoolVal = core.BoolVal
var StringVal = core.StringVal
var CharVal = core.CharVal
var VectorVal = core.VectorVal
var PairVal = core.PairVal
var SchemeVal = core.SchemeVal

// Engine API.
type Engine = engine.Engine
type Builtin = engine.Builtin
type BuiltinExec = engine.BuiltinExec
type Stack[T any] = engine.Stack[T]
type StagingFunc = engine.StagingFunc

var NewEngine = engine.New()
