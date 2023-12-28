package gxx

import (
	"os/exec"
	"path/filepath"

	"github.com/binary-soup/bchef/cmd/compiler"
)

const (
	COMPILER = "g++"
	WARNINGS = "-Wall"
	STANDARD = "-std=c++17"
)

func NewGXXCompiler(d compiler.Data) GXX {
	gxx := GXX{
		includes:     make([]string, len(d.Includes)),
		staticLibs:   d.StaticLibraries,
		libraryPaths: make([]string, len(d.LibraryPaths)),
		sharedLibs:   make([]string, len(d.SharedLibraries)),
		runtimePath:  d.RuntimePath,
	}

	for i, include := range d.Includes {
		gxx.includes[i] = "-I" + include
	}

	for i, path := range d.LibraryPaths {
		gxx.libraryPaths[i] = "-L" + path
	}

	for i, lib := range d.SharedLibraries {
		gxx.sharedLibs[i] = "-l" + lib
	}

	return gxx
}

type GXX struct {
	includes     []string
	staticLibs   []string
	libraryPaths []string
	sharedLibs   []string
	runtimePath  string
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

	if len(gxx.runtimePath) > 0 {
		args = append(args, filepath.Join("-Wl,-rpath,$ORIGIN", gxx.runtimePath))
	}

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
	args = append(args, gxx.libraryPaths...)
	args = append(args, gxx.sharedLibs...)

	if len(gxx.runtimePath) > 0 {
		args = append(args, filepath.Join("-Wl,-rpath,$ORIGIN", gxx.runtimePath))
	}

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
