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
	indices := c.calcChangedSources()

	for i, index := range indices {
		if ok := c.compileObject(c.Recipe.SourceFiles[index], c.Recipe.ObjectFiles[index], float32(i+1)/float32(len(indices))); !ok {
			return false
		}
	}
	return true
}

func (c Compiler) calcChangedSources() []int {
	tracker := newTracker(c.Recipe.SourceDir)

	tracker.LoadCache(c.Recipe)
	defer tracker.SaveCache(c.Recipe)

	indices := tracker.CalcChangedIndices(c.Recipe.SourceFiles, c.Recipe.ObjectFiles)
	style.InfoV2.Printf("%s+ [%d] changed %s\n", c.Indent, len(indices), style.SelectPlural("source", "sources", len(indices)))

	return indices
}

func (c Compiler) compileObject(src string, obj string, percent float32) bool {
	return c.compile(style.Create, percent, []string{"-c"}, c.Recipe.ObjectDir, obj, c.Recipe.SourceDir, src)
}

func (c Compiler) CompileExecutable() bool {
	count := len(c.Recipe.ObjectFiles)
	style.InfoV2.Printf("%s+ [%d] %s\n", c.Indent, count, style.SelectPlural("object", "objects", count))

	sources := append([]string{filepath.Join(c.Recipe.Path, "main.cxx")}, c.Recipe.ObjectFiles...)
	return c.compile(style.BoldCreate, 1.0, []string{}, c.Recipe.Path, c.Recipe.Executable, c.Recipe.Path, sources...)
}

func (c Compiler) compile(createStyle style.Style, percent float32, flags []string, outDir string, out string, srcDir string, sources ...string) bool {
	style.BoldInfo.Printf("%s[%3d%%] ", c.Indent, int(percent*100))
	style.File.Print(c.Recipe.TrimDir(srcDir, sources[0]))
	fmt.Print(" -> ")

	args := append(compileFlags, "-I", c.Recipe.SourceDir)
	args = append(args, flags...)
	args = append(args, sources...)
	args = append(args, "-o", out)

	cmd := exec.Command("g++", args...)
	cmd.Stderr = os.Stdout

	_, err := cmd.Run().(*exec.ExitError)
	if !err {
		createStyle.Println(c.Recipe.TrimDir(outDir, out))
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
