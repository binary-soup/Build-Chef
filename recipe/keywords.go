package recipe

import (
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/common"
	"github.com/binary-soup/bchef/reader"
)

func (*Recipe) whileNotKeyword(r *reader.Reader, do func(string)) error {
	for line, hasNext := r.Next(); hasNext && line[0] != '|'; line, hasNext = r.Next() {
		do(line)
	}

	r.Rewind(1)
	return nil
}

func (*Recipe) peekExtraLine(r *reader.Reader) bool {
	line, hasNext := r.Next()
	if hasNext {
		r.Rewind(1)
	}

	return hasNext && line[0] != '|'
}

func (*Recipe) firstOrEmpty(tokens []string) string {
	if len(tokens) > 0 {
		return tokens[0]
	} else {
		return ""
	}
}

func (rec *Recipe) parseExecutableKeyword(r *reader.Reader, tokens []string) error {
	if len(rec.Executable) > 0 {
		return r.Error("duplicate executable keyword")
	}

	rec.Executable = rec.firstOrEmpty(tokens)
	if len(rec.Executable) == 0 {
		return r.Error("missing or empty executable name")
	}

	return nil
}

func (rec *Recipe) parseSourcesKeyword(r *reader.Reader, tokens []string) error {
	token := rec.firstOrEmpty(tokens)

	if !rec.peekExtraLine(r) {
		return rec.parseSourcesSingle(r, token)
	} else {
		return rec.parseSourcesMulti(r, token)
	}
}

func (rec *Recipe) parseSourcesSingle(r *reader.Reader, src string) error {
	if len(src) == 0 {
		return r.Error("missing or empty source")
	}

	rec.addSource(src)
	return nil
}

func (rec *Recipe) parseSourcesMulti(r *reader.Reader, path string) error {
	return rec.whileNotKeyword(r, func(line string) {
		rec.addSource(filepath.Join(path, line))
	})
}

func (rec *Recipe) addSource(src string) {
	rec.SourceFiles = append(rec.SourceFiles, src)
	rec.ObjectFiles = append(rec.ObjectFiles, common.ReplaceChar(src, "/", '.')+".o")
}

func (rec *Recipe) parseIncludesKeyword(r *reader.Reader, tokens []string) error {
	if !rec.peekExtraLine(r) {
		return rec.parseIncludesSingle(r, rec.firstOrEmpty(tokens))
	} else {
		return rec.parseIncludesMulti(r)
	}
}

func (rec *Recipe) parseIncludesSingle(r *reader.Reader, include string) error {
	if len(include) == 0 {
		return r.Error("empty or missing include")
	}

	rec.addInclude(include)
	return nil
}

func (rec *Recipe) parseIncludesMulti(r *reader.Reader) error {
	return rec.whileNotKeyword(r, func(include string) {
		rec.addInclude(include)
	})
}

func (rec *Recipe) addInclude(include string) {
	if !filepath.IsAbs(include) {
		include = rec.JoinPath(include)
	}
	rec.Includes = append(rec.Includes, include)
}

func (rec *Recipe) parseSharedLibsKeyword(r *reader.Reader, tokens []string) error {
	token := rec.firstOrEmpty(tokens)

	if !rec.peekExtraLine(r) {
		return rec.parseSharedLibsSingle(r, token)
	} else {
		return rec.parseSharedLibsMulti(r, token)
	}
}

func (rec *Recipe) parseSharedLibsSingle(r *reader.Reader, lib string) error {
	if len(lib) == 0 {
		return r.Error("empty or missing library")
	}

	rec.addSharedLibrary(lib)
	return nil
}

func (rec *Recipe) parseSharedLibsMulti(r *reader.Reader, path string) error {
	if len(path) > 0 {
		rec.LibraryPaths = append(rec.LibraryPaths, path)
	}

	return rec.whileNotKeyword(r, func(line string) {
		rec.addSharedLibrary(line)
	})
}

func (rec *Recipe) addSharedLibrary(lib string) {
	rec.SharedLibs = append(rec.SharedLibs, lib)
}

func (rec *Recipe) parseStaticLibsKeyword(r *reader.Reader, tokens []string) error {
	token := rec.firstOrEmpty(tokens)

	if !rec.peekExtraLine(r) {
		return rec.parseStaticLibsSingle(r, token)
	} else {
		return rec.parseStaticLibsMulti(r, token)
	}
}

func (rec *Recipe) parseStaticLibsSingle(r *reader.Reader, lib string) error {
	if len(lib) == 0 {
		return r.Error("empty or missing library")
	}

	rec.addStaticLibrary(lib)
	return nil
}

func (rec *Recipe) parseStaticLibsMulti(r *reader.Reader, path string) error {
	return rec.whileNotKeyword(r, func(line string) {
		rec.addStaticLibrary(filepath.Join(path, line))
	})
}

func (rec *Recipe) addStaticLibrary(lib string) {
	if !filepath.IsAbs(lib) {
		lib = rec.JoinPath(lib)
	}
	rec.StaticLibs = append(rec.StaticLibs, lib)
}

func (rec *Recipe) parsePackageKeyword(r *reader.Reader, tokens []string) error {
	token := rec.firstOrEmpty(tokens)
	if len(token) == 0 {
		return r.Error("missing or empty package name")
	}
	pkg := rec.JoinPath(token)

	file, err := os.Open(pkg)
	if err != nil {
		return r.Errorf("error opening package file \"%s\"", pkg)
	}
	defer file.Close()

	return rec.parsePackage(pkg, file)
}
