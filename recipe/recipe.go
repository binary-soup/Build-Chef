package recipe

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/binary-soup/bchef/parser"
)

func Load(path string) (*Recipe, error) {
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, errors.New("file not found")
	}
	if err != nil {
		return nil, errors.Join(errors.New("error opening file"), err)
	}
	defer file.Close()

	r := Recipe{
		Path: filepath.Dir(path),
		Name: filepath.Base(path),
	}
	r.parse(file)

	return &r, nil
}

type Recipe struct {
	Path string
	Name string

	Executable string

	SourceDir   string
	SourceFiles []string

	ObjectDir   string
	ObjectFiles []string
}

func (r Recipe) TrimSourceDir(src string) string {
	return strings.TrimPrefix(src, r.SourceDir+"/")
}

func (r Recipe) TrimObjectDir(obj string) string {
	return strings.TrimPrefix(obj, r.ObjectDir+"/")
}

func (r *Recipe) parse(reader io.Reader) {
	// TODO: handle invalid recipes
	p := parser.New(reader)

	r.Executable = filepath.Join(r.Path, p.NextLine())

	r.SourceDir = filepath.Join(r.Path, strings.TrimRight(p.NextLine(), "/"))
	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
		r.SourceFiles = append(r.SourceFiles, filepath.Join(r.SourceDir, line))
	}

	r.ObjectDir = filepath.Join(r.Path, ".bchef/obj")
	r.ObjectFiles = make([]string, len(r.SourceFiles))
	for i, src := range r.SourceFiles {
		r.ObjectFiles[i] = r.pathToObject(r.TrimSourceDir(src))
	}
}

func (r Recipe) pathToObject(path string) string {
	result := make([]rune, len(path))

	for i, char := range path {
		if char == '/' {
			result[i] = '.'
		} else {
			result[i] = char
		}
	}

	return filepath.Join(r.ObjectDir, string(result)+".o")
}
