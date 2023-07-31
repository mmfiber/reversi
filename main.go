package main

import (
	"flag"
	"fmt"
	"reversi/src/log"
	"reversi/src/ui/gui"
)

// flags
var (
	mode string
)
var logger = log.NewLogger()

type Runner interface {
	Run()
}

func getRunner() (Runner, error) {
	switch mode {
	case "cli":
		return gui.New(), nil
	// case "worker":
	// 	fmt.Println("s:", s)
	default:
		return nil, fmt.Errorf("err")
	}
}

// init func is called after all the variable declarations in the package have evaluated their initializers,
// and those are evaluated only after all the imported packages have been initialized.
// Besides initializations that cannot be expressed as declarations,
// a common use of init functions is to verify or repair correctness of the program state before real execution begins.
// see more: https://go.dev/doc/effective_go#init
func init() {
	flag.StringVar(&mode, "mode", "", "Select mode one of 'cli' and 'worker'")
}

func main() {
	flag.Parse()
	if runner, err := getRunner(); err == nil {
		runner.Run()
	} else {
		logger.Error(err)
	}
}
