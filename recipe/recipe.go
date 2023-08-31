package recipe

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/binary-soup/bchef/parser"
)

const (
	OBJECT_DIR = ".bchef/obj"
)

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

type Recipe struct {
	Path        string
	Name        string
	SourceDir   string
	SourceFiles []string
	ObjectFiles []string
}

func (r Recipe) TrimSourceDir(src string) string {
	return strings.TrimPrefix(src, r.SourceDir+"/")
}

func (Recipe) TrimObjectDir(obj string) string {
	return strings.TrimPrefix(obj, OBJECT_DIR+"/")
}

func (r *Recipe) IsSourceChanged(index int) bool {
	return r.fileModTime(r.SourceFiles[index]) > r.fileModTime(r.ObjectFiles[index])
}

func (r *Recipe) parse(reader io.Reader) {
	p := parser.New(reader)

	r.Name = p.NextLine()
	r.SourceDir = strings.TrimRight(p.NextLine(), "/")

	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
		r.SourceFiles = append(r.SourceFiles, filepath.Join(r.SourceDir, line))
	}

	r.ObjectFiles = make([]string, len(r.SourceFiles))
	for i, src := range r.SourceFiles {
		r.ObjectFiles[i] = r.pathToObject(r.TrimSourceDir(src))
	}
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

func (Recipe) fileModTime(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
}
