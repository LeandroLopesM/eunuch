package util

import (
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