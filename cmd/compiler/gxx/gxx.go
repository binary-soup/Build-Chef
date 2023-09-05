package gxx

import "os/exec"

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

func (gxx GXX) CompileObject(src string, obj string) *exec.Cmd {
	args := []string{WARNINGS, STANDARD}
	args = append(args, gxx.includes...)
	args = append(args, "-c", "-o", obj, src)

	return exec.Command(BINARY, args...)
}

func (gxx GXX) CompileExecutable(src string, out string, objs ...string) *exec.Cmd {
	args := []string{WARNINGS, STANDARD}
	args = append(args, gxx.includes...)
	args = append(args, "-o", out, src)
	args = append(args, objs...)
	args = append(args, gxx.libraries...)

	return exec.Command(BINARY, args...)
}
