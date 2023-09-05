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
	err = r.parse(file)

	return &r, err
}

type Recipe struct {
	Path string
	Name string

	Executable       string
	ExecutableSource string

	SourceDir   string
	SourceFiles []string

	ObjectDir   string
	ObjectFiles []string
}

func (Recipe) TrimDir(dir string, file string) string {
	return strings.TrimPrefix(file, dir+"/")
}

func (r Recipe) JoinSourceDir(src string) string {
	return filepath.Join(r.SourceDir, src)
}

func (r Recipe) TrimSourceDir(src string) string {
	return r.TrimDir(r.SourceDir, src)
}

func (r Recipe) JoinObjectDir(obj string) string {
	return filepath.Join(r.ObjectDir, obj)
}

func (r Recipe) TrimObjectDir(obj string) string {
	return r.TrimDir(r.ObjectDir, obj)
}

func (r *Recipe) parse(reader io.Reader) error {
	p := parser.New(reader, 1)

	for line, hasNext := p.Next(); hasNext; line, hasNext = p.Next() {
		if line[0] != '|' {
			continue
		}
		tokens := strings.Split(line[1:], " ")

		err := r.parseKeyword(&p, tokens[0], tokens[1:])
		if err != nil {
			return err
		}
	}

	if len(r.Executable) == 0 {
		return p.Error("missing executable keyword")
	}

	return nil
}

func (r *Recipe) parseKeyword(p *parser.Parser, keyword string, tokens []string) error {
	switch keyword {
	case "EXECUTABLE":
		return r.parseExecutable(p, tokens)
	case "SOURCES":
		return r.parseSources(p, tokens)
	default:
		return p.Errorf("unknown keyword \"%s\"", keyword)
	}
}

func (r *Recipe) parseExecutable(p *parser.Parser, tokens []string) error {
	if len(tokens) < 1 || len(tokens[0]) == 0 {
		return p.Error("missing or empty executable name")
	}
	r.Executable = filepath.Join(r.Path, tokens[0])

	if len(tokens) < 2 || len(tokens[1]) == 0 {
		return p.Error("missing or empty executable source")
	}
	r.ExecutableSource = filepath.Join(r.Path, tokens[1])

	return nil
}

func (r *Recipe) parseSources(p *parser.Parser, tokens []string) error {
	dir := "."
	if len(tokens) > 0 {
		dir = tokens[0]
	}
	r.SourceDir = filepath.Join(r.Path, strings.TrimRight(dir, "/"))

	for line, hasNext := p.Next(); hasNext && line[0] != '|'; line, hasNext = p.Next() {
		r.SourceFiles = append(r.SourceFiles, r.JoinSourceDir(line))
	}

	r.ObjectDir = filepath.Join(r.Path, ".bchef/obj")

	r.ObjectFiles = make([]string, len(r.SourceFiles))
	for i, src := range r.SourceFiles {
		r.ObjectFiles[i] = r.pathToObject(r.TrimSourceDir(src))
	}

	p.Rewind(1)
	return nil
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

	return r.JoinObjectDir(string(result) + ".o")
}
