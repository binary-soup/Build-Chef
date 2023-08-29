package recipe

import (
	"errors"
	"io"
	"os"
	"path/filepath"
)

type Recipe struct {
	Path        string
	Name        string
	IncludeDirs []string
	SourceFiles []string
	ObjectFiles []string
}

func (Recipe) sourceToObject(src string) string {
	result := make([]rune, len(src))

	for i, char := range src {
		if char == '/' {
			result[i] = '.'
		} else {
			result[i] = char
		}
	}

	return filepath.Join(".bchef", "objects", string(result)+".o")
}

func (rec *Recipe) parse(r io.Reader) {
	p := newParser(r)

	rec.Name = p.NextLine()
	dir := p.NextLine()

	rec.IncludeDirs = []string{"-I" + dir}

	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
		rec.SourceFiles = append(rec.SourceFiles, filepath.Join(dir, line))
	}

	rec.ObjectFiles = make([]string, len(rec.SourceFiles))
	for i, src := range rec.SourceFiles {
		rec.ObjectFiles[i] = rec.sourceToObject(src)
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

	rec := Recipe{Path: path}
	rec.parse(file)

	return &rec, nil
}
