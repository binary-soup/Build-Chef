package gxx

import (
	"os/exec"

	"github.com/binary-soup/bchef/cmd/compiler"
)

const (
	BINARY   = "g++"
	WARNINGS = "-Wall"
	STANDARD = "-std=c++17"
	DEBUG    = "-g"
)

func NewGXXCompiler(includes []string, libraryPaths []string, libraries []string) GXX {
	gxx := GXX{
		includes:     make([]string, len(includes)),
		libraryPaths: make([]string, len(libraryPaths)),
		libraries:    make([]string, len(libraries)),
	}

	for i, include := range includes {
		gxx.includes[i] = "-I" + include
	}

	for i, path := range libraryPaths {
		gxx.libraryPaths[i] = "-L" + path
	}

	for i, lib := range libraries {
		gxx.libraries[i] = "-l" + lib
	}

	return gxx
}

type GXX struct {
	includes     []string
	libraryPaths []string
	libraries    []string
}

func (gxx GXX) CompileObject(opts compiler.Options, src string, obj string) *exec.Cmd {
	args := gxx.createArgs(opts)

	args = append(args, gxx.includes...)
	args = append(args, "-c", "-o", obj, src)

	return exec.Command(BINARY, args...)
}

func (gxx GXX) CompileExecutable(opts compiler.Options, src string, out string, objs ...string) *exec.Cmd {
	args := gxx.createArgs(opts)

	args = append(args, gxx.includes...)
	args = append(args, "-o", out, src)
	args = append(args, objs...)
	args = append(args, gxx.libraries...)

	return exec.Command(BINARY, args...)
}

func (GXX) createArgs(opts compiler.Options) []string {
	args := []string{WARNINGS, STANDARD}

	if opts.Debug {
		args = append(args, DEBUG)
	}

	for _, marco := range opts.Macros {
		args = append(args, "-D"+marco)
	}

	return args
}
