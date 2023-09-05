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
		Path:        filepath.Dir(path),
		Name:        filepath.Base(path),
		SourceFiles: []string{},
		ObjectFiles: []string{},
	}
	r.ObjectDir = filepath.Join(r.Path, ".bchef/obj")

	return &r, r.parse(file)
}

type Recipe struct {
	Name      string
	Path      string
	ObjectDir string

	Executable string
	MainSource string

	SourceFiles []string
	ObjectFiles []string
}

func (Recipe) TrimDir(dir string, file string) string {
	return strings.TrimPrefix(file, dir+"/")
}

func (r Recipe) JoinPath(src string) string {
	return filepath.Join(r.Path, src)
}

func (r Recipe) TrimPath(src string) string {
	return r.TrimDir(r.Path, src)
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
	if len(r.Executable) > 0 {
		return p.Error("duplicate executable keyword")
	}

	if len(tokens) < 1 || len(tokens[0]) == 0 {
		return p.Error("missing or empty executable name")
	}
	r.Executable = filepath.Join(r.Path, tokens[0])

	if len(tokens) < 2 || len(tokens[1]) == 0 {
		return p.Error("missing or empty main source")
	}
	r.MainSource = filepath.Join(r.Path, tokens[1])

	return nil
}

func (r *Recipe) parseSources(p *parser.Parser, tokens []string) error {
	srcDir := "."
	if len(tokens) > 0 {
		srcDir = tokens[0]
	}

	for line, hasNext := p.Next(); hasNext && line[0] != '|'; line, hasNext = p.Next() {
		src := filepath.Join(srcDir, line)

		r.SourceFiles = append(r.SourceFiles, r.JoinPath(src))
		r.ObjectFiles = append(r.ObjectFiles, r.JoinObjectDir(r.srcToObject(src)))
	}

	p.Rewind(1)
	return nil
}

func (r Recipe) srcToObject(path string) string {
	result := make([]rune, len(path))

	for i, char := range path {
		if char == '/' {
			result[i] = '.'
		} else {
			result[i] = char
		}
	}

	return string(result) + ".o"
}
