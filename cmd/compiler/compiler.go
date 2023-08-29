package compiler

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/binary-soup/bchef/style"
)

var compileFlags = []string{"-Wall", "-std=c++17"}

type Compiler struct {
	Log         *log.Logger
	IncludeDirs []string
}

func (c Compiler) CompileObject(src string, obj string) bool {
	return c.compile(style.Create, []string{"-c"}, obj, src)
}

func (c Compiler) CompileObjects(sources []string, objects []string) bool {
	for i, src := range sources {
		if ok := c.CompileObject(src, objects[i]); !ok {
			return false
		}
	}
	return true
}

func (c Compiler) CompileExecutable(src string, exe string, objs []string) bool {
	style.FileV2.Printf("  + [%d] object files\n", len(objs))

	sources := append([]string{src}, objs...)
	return c.compile(style.BoldCreate, []string{}, exe, sources...)
}

func (c Compiler) compile(createStyle style.Style, flags []string, out string, sources ...string) bool {
	fmt.Print("  ", style.FileV1.Format(sources[0]), " -> ")

	args := append(compileFlags, c.IncludeDirs...)
	args = append(args, flags...)
	args = append(args, sources...)
	args = append(args, "-o", out)

	cmd := exec.Command("g++", args...)
	cmd.Stderr = os.Stdout

	_, err := cmd.Run().(*exec.ExitError)
	if !err {
		createStyle.Println(out)
	}

	res := "PASS"
	if err {
		res = "FAIL"
	}

	c.Log.Printf("[Compile %s] %s %s\n", sources[0], res, cmd.String())
	return !err
}
