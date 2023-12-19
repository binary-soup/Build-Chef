package recipe

import (
	"errors"
	"os"
	"path/filepath"
)

func Load(path string) (*Recipe, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Join(errors.New("error opening file"), err)
	}
	defer file.Close()

	r := Recipe{
		Path:             filepath.Dir(path),
		Name:             filepath.Base(path),
		SourceFiles:      []string{},
		ObjectFiles:      []string{},
		Includes:         []string{},
		LibraryPaths:     []string{},
		LinkedSharedLibs: []string{},
		LinkedStaticLibs: []string{},
	}

	r.ObjectPath = filepath.Join(r.Path, ".bchef/obj")
	r.Includes = append(r.Includes, r.Path)

	return &r, r.parseRecipe(file)
}

type Recipe struct {
	Name       string
	Path       string
	ObjectPath string

	Executable string

	SourceFiles []string
	ObjectFiles []string

	Includes []string

	LibraryPaths     []string
	LinkedSharedLibs []string
	LinkedStaticLibs []string
}

func (r Recipe) FullPath() string {
	return filepath.Join(r.Path, r.Name)
}

func (r Recipe) JoinPath(src string) string {
	return filepath.Join(r.Path, src)
}

func (Recipe) GetMode(debug bool) string {
	if debug {
		return "debug"
	} else {
		return "release"
	}
}

func (r Recipe) GetExecutable(debug bool) string {
	if debug {
		return r.Executable + "." + r.GetMode(true)
	} else {
		return r.Executable
	}
}

func (r Recipe) GetObjectPath(debug bool) string {
	return filepath.Join(r.ObjectPath, r.GetMode(debug))
}

func (r Recipe) JoinObjectPath(obj string, debug bool) string {
	return filepath.Join(r.GetObjectPath(debug), obj)
}
