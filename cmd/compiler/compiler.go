package compiler

import (
	"log"
	"os"
	"os/exec"

	"github.com/binary-soup/bchef/recipe"
)

func NewCompiler(log *log.Logger, impl CompilerImpl, opts Options) Compiler {
	return Compiler{
		log:  log,
		impl: impl,
		Opts: opts,
	}
}

type Options struct {
	Debug  bool
	PIC    bool
	Macros []string
}

type CompilerImpl interface {
	CompileObject(opts Options, src string, obj string) *exec.Cmd
	CreateExecutable(opts Options, exec string, objs ...string) *exec.Cmd
	CreateStaticLibrary(opts Options, lib string, objs ...string) *exec.Cmd
	CreateSharedLibrary(opts Options, lib string, objs ...string) *exec.Cmd
}

type Compiler struct {
	log  *log.Logger
	impl CompilerImpl
	Opts Options
}

func (c Compiler) ComputeChangedSources(r *recipe.Recipe) []int {
	tracker := newTracker(r)

	tracker.LoadCache()
	defer tracker.SaveCache()

	indices := tracker.CalcChangedIndices(r.SourceFiles, r.ObjectFiles, r.GetObjectPath(c.Opts.Debug))

	return indices
}

func (c Compiler) CompileObject(src string, obj string) bool {
	cmd := c.impl.CompileObject(c.Opts, src, obj)
	res := c.runCommand(cmd)

	c.logCommand(cmd.String(), src, 0, res)
	return res
}

func (c Compiler) CreateExecutable(exec string, objs ...string) bool {
	cmd := c.impl.CreateExecutable(c.Opts, exec, objs...)
	res := c.runCommand(cmd)

	c.logCommand(cmd.String(), "", len(objs), res)
	return res
}

func (c Compiler) CreateStaticLibrary(exec string, objs ...string) bool {
	cmd := c.impl.CreateStaticLibrary(c.Opts, exec, objs...)
	res := c.runCommand(cmd)

	c.logCommand(cmd.String(), "", len(objs), res)
	return res
}

func (c Compiler) CreateSharedLibrary(exec string, objs ...string) bool {
	cmd := c.impl.CreateSharedLibrary(c.Opts, exec, objs...)
	res := c.runCommand(cmd)

	c.logCommand(cmd.String(), "", len(objs), res)
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
