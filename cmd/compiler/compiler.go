package compiler

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

var compileFlags = []string{"-Wall", "-std=c++17"}

type Compiler struct {
	Log    *log.Logger
	Recipe *recipe.Recipe
}

func (c Compiler) CompileObjects() bool {
	for i, src := range c.Recipe.SourceFiles {
		if ok := c.compile(style.Create, []string{"-c"}, c.Recipe.ObjectFiles[i], src); !ok {
			return false
		}
	}
	return true
}

func (c Compiler) CompileExecutable() bool {
	style.FileV2.Printf("  + [%d] object files\n", len(c.Recipe.ObjectFiles))

	sources := append([]string{"main.cxx"}, c.Recipe.ObjectFiles...)
	return c.compile(style.BoldCreate, []string{}, c.Recipe.Name, sources...)
}

func (c Compiler) compile(createStyle style.Style, flags []string, out string, sources ...string) bool {
	fmt.Print("  ", style.FileV1.Format(c.Recipe.TrimSourceDir(sources[0])), " -> ")

	args := append(compileFlags, "-I", c.Recipe.SourceDir)
	args = append(args, flags...)
	args = append(args, sources...)
	args = append(args, "-o", out)

	cmd := exec.Command("g++", args...)
	cmd.Stderr = os.Stdout

	_, err := cmd.Run().(*exec.ExitError)
	if !err {
		createStyle.Println(c.Recipe.TrimObjectDir(out))
	}

	res := "PASS"
	if err {
		res = "FAIL"
	}

	c.Log.Printf("[Compile %s] %s %s\n", sources[0], res, cmd.String())
	return !err
}
