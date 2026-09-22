package lib

import (
	"github.com/leandrolopesm/eunuch/core"
	"github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/parser"
)

/* The primary library API
 * for eunuch.
 */

// Core API
type PairVal = core.PairVal
type Position = core.Position
type Scheme = core.Scheme
type Type = core.Type
type Unit = core.Unit
type VectorVal = core.VectorVal

// Engine API
type Engine = engine.Engine
type Builtin = engine.Builtin
type BuiltinExec = engine.BuiltinExec
type Stack[T any] = engine.Stack[T]
type StagingFunc = engine.StagingFunc

// Lexer API
type Lexer = parser.Lexer
type Iterator[T any] = parser.Iterator[T]

// 