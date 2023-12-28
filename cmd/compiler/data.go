package compiler

import (
	"path/filepath"

	"github.com/binary-soup/bchef/recipe"
)

func NewData(r *recipe.Recipe) Data {
	return Data{
		Includes:        r.Includes,
		StaticLibraries: r.LinkedStaticLibs,
		LibraryPaths:    r.LibraryPaths,
		SharedLibraries: r.LinkedSharedLibs,
		RuntimePath:     "",
	}
}

type Data struct {
	Includes []string

	StaticLibraries []string

	LibraryPaths    []string
	SharedLibraries []string
	RuntimePath     string
}

func (d *Data) LinkStaticLayer(outputDir string, target string) {
	d.StaticLibraries = append(d.StaticLibraries, filepath.Join(outputDir, target))
}

func (d *Data) LinkSharedLayer(outputDir string, target string) {
	d.LibraryPaths = append(d.LibraryPaths, outputDir)
	d.SharedLibraries = append(d.SharedLibraries, target)

	d.RuntimePath = "." //TODO: ensure relative to linking target
}
