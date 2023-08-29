package main

import (
	"fmt"
	"os"

	"github.com/binary-soup/bchef/cmd"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

var cmds = map[string]cmd.Command{
	"cook":  cmd.CookCmd{},
	"clean": cmd.CleanCmd{},
}

func RunCommand(name string) error {
	r, err := recipe.Load(".")
	if err != nil {
		return err
	}

	fmt.Println("Recipe loaded from", style.BoldFile.Format(r.Path))

	cmd, exists := cmds[name]
	if !exists {
		return fmt.Errorf("unknown command \"%s\"", name)
	}

	return cmd.Run(r)
}

func main() {
	if len(os.Args) < 2 {
		// TODO: print help
		fmt.Println("no command given")
		return
	}

	if err := RunCommand(os.Args[1]); err != nil {
		fmt.Println(style.BoldError.Format("Error:"), err)
	}
}
