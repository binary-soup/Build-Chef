package compiler

import (
	"log"
	"os"
	"os/exec"

	"github.com/binary-soup/bchef/recipe"
)

func NewCompiler(log *log.Logger, impl CompilerImpl) Compiler {
	return Compiler{
		Log:  log,
		Impl: impl,
	}
}

type CompilerImpl interface {
	CompileObject(src string, obj string) *exec.Cmd
	CompileExecutable(src string, exec string, objs ...string) *exec.Cmd
}

type Compiler struct {
	Log  *log.Logger
	Impl CompilerImpl
}

func (Compiler) ComputeChangedSources(r *recipe.Recipe) []int {
	tracker := newTracker(r.Path)

	tracker.LoadCache(r)
	defer tracker.SaveCache(r)

	indices := tracker.CalcChangedIndices(r.SourceFiles, r.ObjectFiles)

	return indices
}

func (c Compiler) CompileObject(src string, obj string) bool {
	cmd := c.Impl.CompileObject(src, obj)
	res := c.runCommand(cmd)

	c.logCommand(cmd.String(), src, 0, res)
	return res
}

func (c Compiler) CompileExecutable(src string, exec string, objs ...string) bool {
	cmd := c.Impl.CompileExecutable(src, exec, objs...)
	res := c.runCommand(cmd)

	c.logCommand(cmd.String(), src, len(objs), res)
	return res
}

func (c Compiler) runCommand(cmd *exec.Cmd) bool {
	cmd.Stderr = os.Stdout
	err := cmd.Run()

	_, ok := err.(*exec.ExitError)
	return !ok
}

func (c Compiler) logCommand(cmdStr string, src string, numObjs int, res bool) {
	log := "[Compile " + src + "] "

	if res {
		log += "PASS"
	} else {
		log += "FAIL"
	}

	c.Log.Println(log, cmdStr)
}
