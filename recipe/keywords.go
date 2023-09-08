package recipe

import (
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/reader"
)

func (*Recipe) whileNotKeyword(r *reader.Reader, do func(string)) error {
	for line, hasNext := r.Next(); hasNext && line[0] != '|'; line, hasNext = r.Next() {
		do(line)
	}

	r.Rewind(1)
	return nil
}

func (rec *Recipe) parseExecutableKeyword(r *reader.Reader, tokens []string) error {
	if len(rec.Executable) > 0 {
		return r.Error("duplicate executable keyword")
	}

	if len(tokens) < 1 || len(tokens[0]) == 0 {
		return r.Error("missing or empty executable name")
	}
	rec.Executable = tokens[0]

	if len(tokens) < 2 || len(tokens[1]) == 0 {
		return r.Error("missing or empty main source")
	}
	rec.MainSource = tokens[1]

	return nil
}

func (rec *Recipe) parseSourcesKeyword(r *reader.Reader, tokens []string) error {
	srcDir := "."
	if len(tokens) > 0 {
		srcDir = tokens[0]
	}

	return rec.whileNotKeyword(r, func(line string) {
		src := filepath.Join(srcDir, line)

		rec.SourceFiles = append(rec.SourceFiles, src)
		rec.ObjectFiles = append(rec.ObjectFiles, rec.srcToObject(src))
	})
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

func (rec *Recipe) parseIncludesKeyword(r *reader.Reader) error {
	return rec.whileNotKeyword(r, func(line string) {
		rec.Includes = append(rec.Includes, line)
	})
}

func (rec *Recipe) parseLibrariesKeyword(r *reader.Reader, tokens []string) error {
	if len(tokens) > 0 && len(tokens[0]) > 0 {
		rec.LibraryPaths = append(rec.LibraryPaths, tokens[0])
	}

	return rec.whileNotKeyword(r, func(line string) {
		rec.Libraries = append(rec.Libraries, line)
	})
}

func (rec *Recipe) parsePackageKeyword(r *reader.Reader, tokens []string) error {
	if len(tokens) < 1 || len(tokens[0]) == 0 {
		return r.Error("missing or empty package name")
	}
	pkg := rec.JoinPath(tokens[0])

	file, err := os.Open(pkg)
	if err != nil {
		return r.Errorf("error opening package file \"%s\"", pkg)
	}
	defer file.Close()

	return rec.parsePackage(pkg, file)
}
