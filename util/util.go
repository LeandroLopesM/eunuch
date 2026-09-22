package util

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func If[T any](cond bool, vtrue, vfalse T) T {
    if cond {
        return vtrue
    }
    return vfalse
}

func Logger(opts Options) {
	log.SetReportTimestamp(false)
	log.SetLevel(If(*opts.Verbose, log.DebugLevel, log.InfoLevel))
}

type Options struct {
	Verbose *bool
	Repl bool
	Files []string
}

func Args() Options {
	flag.CommandLine.Name()
	opt := Options{
		Verbose: flag.Bool("verbose", false, "Enable verbose logging"),
		Repl: false,
		Files: []string{},
	}

	flag.BoolVar(&opt.Repl, "repl", false, "Run as in REPL mode")

	help := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *help {
		flag.Usage();
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		if !opt.Repl {
			print("No files provided, running in REPL")
		}
		opt.Repl = true
	} else {
		opt.Files = flag.Args()
	}

	return opt
}

func Assert[T any](a T, err error) T {
	if err != nil {
		panic(fmt.Sprintf("Assert non-err failed: %s", err))
	}

	return a
}

type Option[T any] struct {
	val T
	inUse bool // can't be 'empty' because the default constructor is { nil, false } and would give false-positivies
}

func (o Option[T]) Unwrap() T {
	if !o.inUse {
		panic("Unwrap on nil option")
	}

	return o.val
}

func (o Option[T]) Try() (T, error) {
	var tmp T
	
	if !o.inUse {
		return tmp, errors.New("Empty option")
	}

	return o.val, nil
}

func None[T any]() Option[T] {
	return Option[T]{
		inUse: false,
	}
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		val: val,
		inUse: true,
	}
}