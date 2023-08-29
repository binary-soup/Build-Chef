package recipe

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	INCLUDE_FLAG = "-I"
	OBJECT_DIR   = ".bchef/obj"
)

type Recipe struct {
	Path        string
	Name        string
	IncludeDirs []string
	SourceFiles []string
	ObjectFiles []string
}

func (r Recipe) GetSourceDir() string {
	return r.IncludeDirs[0][len(INCLUDE_FLAG):]
}

func (r Recipe) TrimSourceDir(src string) string {
	return strings.TrimPrefix(src, r.GetSourceDir()+"/")
}

func (Recipe) TrimObjectDir(obj string) string {
	return strings.TrimPrefix(obj, OBJECT_DIR+"/")
}

func (Recipe) pathToObject(path string) string {
	result := make([]rune, len(path))

	for i, char := range path {
		if char == '/' {
			result[i] = '.'
		} else {
			result[i] = char
		}
	}

	return filepath.Join(OBJECT_DIR, string(result)+".o")
}

func (rec *Recipe) parse(r io.Reader) {
	p := newParser(r)

	rec.Name = p.NextLine()
	dir := strings.TrimRight(p.NextLine(), "/")

	rec.IncludeDirs = []string{INCLUDE_FLAG + dir}

	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
		rec.SourceFiles = append(rec.SourceFiles, filepath.Join(dir, line))
	}

	rec.ObjectFiles = make([]string, len(rec.SourceFiles))
	for i, src := range rec.SourceFiles {
		rec.ObjectFiles[i] = rec.pathToObject(rec.TrimSourceDir(src))
	}
}

func Load(path string) (*Recipe, error) {
	path = filepath.Join(path, "recipe.txt")

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, errors.New("recipe file not found")
	}
	if err != nil {
		return nil, errors.Join(errors.New("error opening file"), err)
	}
	defer file.Close()

	r := Recipe{Path: path}
	r.parse(file)

	return &r, nil
}
