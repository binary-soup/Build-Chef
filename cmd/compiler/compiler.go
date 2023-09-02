package compiler

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

var compileFlags = []string{"-Wall", "-std=c++17"}

type Compiler struct {
	Indent string
	Log    *log.Logger
	Recipe *recipe.Recipe
}

func (c Compiler) CompileObjects() bool {
	tracker := newTracker(c.Recipe.SourceDir)

	tracker.LoadCache(c.Recipe.Path)
	defer tracker.SaveCache(c.Recipe.Path)

	for i, src := range c.Recipe.SourceFiles {
		obj := c.Recipe.ObjectFiles[i]

		if !tracker.NeedsCompiling(src, obj) {
			continue
		}

		if ok := c.compileObject(src, obj, float32(i)/float32(len(c.Recipe.SourceFiles))); !ok {
			return false
		}
	}
	return true
}

func (c Compiler) compileObject(src string, obj string, percent float32) bool {
	return c.compile(style.Create, percent, []string{"-c"}, obj, src)
}

func (c Compiler) CompileExecutable() bool {
	style.FileV2.Printf("%s+ [%d] objects\n", c.Indent, len(c.Recipe.ObjectFiles))

	sources := append([]string{filepath.Join(c.Recipe.Path, "main.cxx")}, c.Recipe.ObjectFiles...)
	return c.compile(style.BoldCreate, 1.0, []string{}, c.Recipe.Executable, sources...)
}

func (c Compiler) compile(createStyle style.Style, percent float32, flags []string, out string, sources ...string) bool {
	style.BoldInfo.Printf("%s[%3d%%] ", c.Indent, int(percent*100))
	style.File.Print(c.Recipe.TrimSourceDir(sources[0]))
	fmt.Print(" -> ")

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

	c.logCommand(cmd.String(), sources[0], len(sources)-1, !err)
	return !err
}

func (c Compiler) logCommand(cmdStr string, src string, numObjs int, res bool) {
	log := "[Compile " + src

	if numObjs > 0 {
		log += fmt.Sprintf(" + (%d) objects", numObjs)
	}
	log += "] "

	if res {
		log += "PASS"
	} else {
		log += "FAIL"
	}

	c.Log.Println(log, cmdStr)
}
