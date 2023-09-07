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
		Path:         filepath.Dir(path),
		Name:         filepath.Base(path),
		SourceFiles:  []string{},
		ObjectFiles:  []string{},
		Includes:     []string{},
		LibraryPaths: []string{},
		Libraries:    []string{},
	}

	r.ObjectPath = filepath.Join(r.Path, ".bchef/obj")
	r.Includes = append(r.Includes, r.Path)

	return &r, r.parse(file)
}

type Recipe struct {
	Name       string
	Path       string
	ObjectPath string

	Executable string
	MainSource string

	SourceFiles []string
	ObjectFiles []string

	Includes []string

	LibraryPaths []string
	Libraries    []string
}

func (r Recipe) JoinPath(src string) string {
	return filepath.Join(r.Path, src)
}

func (r Recipe) GetDebugDir(debug bool) string {
	if debug {
		return "debug"
	} else {
		return "release"
	}
}

func (r Recipe) GetObjectPath(debug bool) string {
	return filepath.Join(r.ObjectPath, r.GetDebugDir(debug))
}

func (r Recipe) JoinObjectPath(obj string, debug bool) string {
	return filepath.Join(r.GetObjectPath(debug), obj)
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
	case "INCLUDES":
		return r.parseIncludes(p)
	case "LIBRARIES":
		return r.parseLibraries(p, tokens)
	case "PACKAGE":
		return r.parsePackage(p, tokens)
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
	r.Executable = tokens[0]

	if len(tokens) < 2 || len(tokens[1]) == 0 {
		return p.Error("missing or empty main source")
	}
	r.MainSource = tokens[1]

	return nil
}

func (r *Recipe) whileNotKeyword(p *parser.Parser, do func(string)) error {
	for line, hasNext := p.Next(); hasNext && line[0] != '|'; line, hasNext = p.Next() {
		do(line)
	}

	p.Rewind(1)
	return nil
}

func (r *Recipe) parseSources(p *parser.Parser, tokens []string) error {
	srcDir := "."
	if len(tokens) > 0 {
		srcDir = tokens[0]
	}

	return r.whileNotKeyword(p, func(line string) {
		src := filepath.Join(srcDir, line)

		r.SourceFiles = append(r.SourceFiles, src)
		r.ObjectFiles = append(r.ObjectFiles, r.srcToObject(src))
	})
}

func (r *Recipe) parseIncludes(p *parser.Parser) error {
	return r.whileNotKeyword(p, func(line string) {
		r.Includes = append(r.Includes, line)
	})
}

func (r *Recipe) parseLibraries(p *parser.Parser, tokens []string) error {
	if len(tokens) > 0 && len(tokens[0]) > 0 {
		r.LibraryPaths = append(r.LibraryPaths, tokens[0])
	}

	return r.whileNotKeyword(p, func(line string) {
		r.Libraries = append(r.Libraries, line)
	})
}

func (r *Recipe) parsePackage(p *parser.Parser, tokens []string) error {
	return nil //TODO: implement
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
