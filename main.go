package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/binary-soup/bchef/cook"
)

func run(args []string) error {
	if len(args) == 0 {
		return errors.New("no command given")
	}

	cmd := args[0]
	args = args[1:]

	if cmd == "cook" {
		return cook.Run(args)
	}

	return fmt.Errorf("unknown command \"%s\"", cmd)
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Println(err)
	}
}
