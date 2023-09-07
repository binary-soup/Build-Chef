package compiler

import (
	"log"
	"os"
	"os/exec"

	"github.com/binary-soup/bchef/recipe"
)

func NewCompiler(log *log.Logger, impl CompilerImpl, debug bool) Compiler {
	return Compiler{
		log:   log,
		impl:  impl,
		Debug: debug,
	}
}

type CompilerImpl interface {
	CompileObject(debug bool, src string, obj string) *exec.Cmd
	CompileExecutable(debug bool, src string, exec string, objs ...string) *exec.Cmd
}

type Compiler struct {
	log   *log.Logger
	impl  CompilerImpl
	Debug bool
}

func (c Compiler) ComputeChangedSources(r *recipe.Recipe) []int {
	tracker := newTracker(r.Path)

	tracker.LoadCache()
	defer tracker.SaveCache()

	indices := tracker.CalcChangedIndices(r.SourceFiles, r.ObjectFiles, r.GetObjectPath(c.Debug))

	return indices
}

func (c Compiler) CompileObject(src string, obj string) bool {
	cmd := c.impl.CompileObject(c.Debug, src, obj)
	res := c.runCommand(cmd)

	c.logCommand(cmd.String(), src, 0, res)
	return res
}

func (c Compiler) CompileExecutable(src string, exec string, objs ...string) bool {
	cmd := c.impl.CompileExecutable(c.Debug, src, exec, objs...)
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

	c.log.Println(log, cmdStr)
}
