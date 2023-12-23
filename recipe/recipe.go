package recipe

import (
	"errors"
	"os"
	"path/filepath"
)

const (
	TARGET_EXECUTABLE     = iota
	TARGET_STATIC_LIBRARY = iota
	TARGET_SHARED_LIBRARY = iota
)

func loadRecipe(path string) (*Recipe, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Join(errors.New("error opening file"), err)
	}
	defer file.Close()

	r := &Recipe{
		Path:             filepath.Dir(path),
		Name:             filepath.Base(path),
		SourceFiles:      []string{},
		ObjectFiles:      []string{},
		Includes:         []string{},
		LibraryPaths:     []string{},
		LinkedSharedLibs: []string{},
		LinkedStaticLibs: []string{},
		Layers:           []string{},
	}

	return r, r.parseRecipe(file)
}

type Recipe struct {
	Name string
	Path string

	Target     string
	TargetType int

	SourceFiles []string
	ObjectFiles []string

	Includes []string

	LibraryPaths     []string
	LinkedSharedLibs []string
	LinkedStaticLibs []string

	Layers []string
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

func (r Recipe) GetTarget(debug bool) string {
	target := r.Target
	if debug {
		target = r.Target + ".d"
	}

	switch r.TargetType {
	case TARGET_STATIC_LIBRARY:
		return "lib" + target + ".a"
	case TARGET_SHARED_LIBRARY:
		return "lib" + target + ".so"
	default: // TARGET_EXECUTABLE
		return target
	}
}

func (r Recipe) GetTargetType() string {
	switch r.TargetType {
	case TARGET_STATIC_LIBRARY:
		return "STATIC_LIBRARY"
	case TARGET_SHARED_LIBRARY:
		return "SHARED_LIBRARY"
	default: // TARGET_EXECUTABLE
		return "EXECUTABLE"
	}
}

func (r Recipe) GetObjectDir(debug bool) string {
	return filepath.Join(".bchef/obj", r.GetMode(debug))
}

func (r Recipe) JoinObjectDir(obj string, debug bool) string {
	return filepath.Join(r.GetObjectDir(debug), obj)
}

func (r Recipe) GetObjectPath(debug bool) string {
	return r.JoinPath(r.GetObjectDir(debug))
}

func (r Recipe) JoinObjectPath(obj string, debug bool) string {
	return r.JoinPath(r.JoinObjectDir(obj, debug))
}
