package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

var compileFlags = []string{"-Wall", "-std=c++17"}

// NOTE: add -g for debug
// TODO: handle command injection
type CookCmd struct{}

func (CookCmd) compile(createStyle style.Style, r *recipe.Recipe, flags []string, out string, src ...string) bool {
	fmt.Print("  ", style.FileV1.Format(src[0]), " -> ")

	args := append(compileFlags, flags...)
	args = append(args, r.IncludeDirs...)
	args = append(args, src...)
	args = append(args, "-o", out)

	cmd := exec.Command("g++", args...)
	cmd.Stderr = os.Stdout

	_, err := cmd.Run().(*exec.ExitError)
	if !err {
		createStyle.Println(out)
	}

	return !err
}

func (cmd CookCmd) compileSourceFiles(r *recipe.Recipe) bool {
	for i, src := range r.SourceFiles {
		if ok := cmd.compile(style.Create, r, []string{"-c"}, r.ObjectFiles[i], src); !ok {
			return false
		}
	}
	return true
}

func (cmd CookCmd) compileExecutable(r *recipe.Recipe) bool {
	style.FileV2.Printf("  + [%d] object files\n", len(r.ObjectFiles))

	src := append([]string{"main.cxx"}, r.ObjectFiles...)
	return cmd.compile(style.BoldCreate, r, []string{}, r.Name, src...)
}

func (cmd CookCmd) Run(r *recipe.Recipe) {
	os.MkdirAll(".bchef/objects", 0755)

	style.Header.Println("Prepping...")
	if ok := cmd.compileSourceFiles(r); !ok {
		style.BoldError.Println("Bad Ingredients!")
		return
	}

	style.Header.Println("Cooking...")
	if ok := cmd.compileExecutable(r); !ok {
		style.BoldError.Println("Burnt!")
		return
	}

	style.BoldSuccess.Println("Bon Appétit!")
}
