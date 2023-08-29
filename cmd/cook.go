package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

type CookCmd struct{}

// NOTE: add -g for debug
// TODO: handle command injection

func (CookCmd) exec(src string, out string) bool {
	fmt.Print("  ", style.File.Format(src), " -> ")

	cmd := exec.Command("g++", "-Wall", "-std=c++17", "-o", out, src)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_, exit := cmd.Run().(*exec.ExitError)
	if !exit {
		style.BoldCreate.Println(out)
	}

	return exit
}

func (cmd CookCmd) Run(r *recipe.Recipe) error {
	style.Header.Println("Cooking...")

	fail := cmd.exec("main.cxx", r.Name)

	if fail {
		style.BoldError.Println("Burnt!")
	} else {
		style.BoldSuccess.Println("Bon Appétit!")
	}

	return nil
}
