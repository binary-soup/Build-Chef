package gxx

import (
	"os/exec"

	"github.com/binary-soup/bchef/cmd/compiler"
)

const (
	COMPILER = "g++"
	WARNINGS = "-Wall"
	STANDARD = "-std=c++17"
)

func NewGXXCompiler(includes []string, staticLibs []string, libraryPaths []string, sharedLibs []string) GXX {
	gxx := GXX{
		includes:     make([]string, len(includes)),
		staticLibs:   staticLibs,
		libraryPaths: make([]string, len(libraryPaths)),
		sharedLibs:   make([]string, len(sharedLibs)),
	}

	for i, include := range includes {
		gxx.includes[i] = "-I" + include
	}

	for i, path := range libraryPaths {
		gxx.libraryPaths[i] = "-L" + path
	}

	for i, lib := range sharedLibs {
		gxx.sharedLibs[i] = "-l" + lib
	}

	return gxx
}

type GXX struct {
	includes     []string
	staticLibs   []string
	libraryPaths []string
	sharedLibs   []string
}

func (gxx GXX) CompileObject(opts compiler.Options, src string, obj string) *exec.Cmd {
	args := gxx.createArgs(opts)

	args = append(args, gxx.includes...)
	args = append(args, "-c", "-o", obj, src)

	return exec.Command(COMPILER, args...)
}

func (gxx GXX) CreateExecutable(opts compiler.Options, out string, objs ...string) *exec.Cmd {
	args := []string{WARNINGS, STANDARD}

	args = append(args, "-o", out)
	args = append(args, objs...)
	args = append(args, gxx.staticLibs...)
	args = append(args, gxx.libraryPaths...)
	args = append(args, gxx.sharedLibs...)

	return exec.Command(COMPILER, args...)
}

func (gxx GXX) CreateStaticLibrary(opts compiler.Options, lib string, objs ...string) *exec.Cmd {
	args := []string{"rcs", lib}
	args = append(args, objs...)

	return exec.Command("ar", args...)
}

func (gxx GXX) CreateSharedLibrary(opts compiler.Options, lib string, objs ...string) *exec.Cmd {
	args := []string{WARNINGS, STANDARD, "-shared"}

	args = append(args, "-o", lib)
	args = append(args, objs...)
	args = append(args, gxx.staticLibs...)
	args = append(args, gxx.libraryPaths...)
	args = append(args, gxx.sharedLibs...)

	return exec.Command(COMPILER, args...)
}

func (GXX) createArgs(opts compiler.Options) []string {
	args := []string{WARNINGS, STANDARD}

	if opts.Debug {
		args = append(args, "-g")
	}

	if opts.PIC {
		args = append(args, "-fPIC")
	}

	for _, marco := range opts.Macros {
		args = append(args, "-D"+marco)
	}

	return args
}
