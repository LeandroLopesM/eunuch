package main

import (
	"context"
	"fmt"

	log "github.com/charmbracelet/log"
	"github.com/leandrolopesm/eunuch/builtin"
	"github.com/leandrolopesm/eunuch/core"
	"github.com/leandrolopesm/eunuch/engine"
	"github.com/leandrolopesm/eunuch/util"

	aurora "github.com/logrusorgru/aurora/v4"

	"github.com/nyaosorg/go-readline-ny"
)

func main() {
	opt := util.Args()
	util.Logger(opt)

	lispEngine := engine.New()
	builtin.RegisterSelf(&lispEngine)

	if opt.Repl {
		runREPL(lispEngine)
	} else {
		for _, file := range opt.Files {
			if err := lispEngine.ExecuteFile(file); err != nil {
				fmt.Print(aurora.Red("Execution failed:"), "\n", err)
			}
		}

		if v, e := lispEngine.Pop(); e == nil { // If the last call pushed a value, print it
			fmt.Printf("\n;; %v", core.SprintUnit(v))
		}
	}
}

func runREPL(lispEngine engine.Engine) {
	var editor readline.Editor

	for {
		if text, err := editor.ReadLine(context.Background()); err != nil {
			if err.Error() == "EOF" {
				break
			}

			log.Errorf("Failed to read input: %s", err)
		} else {
			if err := lispEngine.ExecuteStr(text); err != nil {
				fmt.Print(aurora.Red("Execution failed:\n"), err)
			}

			if v, e := lispEngine.Pop(); e == nil { // If the last call pushed a value, print it
				core.PrintUnit(v)
			}
		}
	}
}