package recipe

import (
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/common"
	"github.com/binary-soup/bchef/reader"
)

func (*Recipe) whileNotKeyword(r *reader.Reader, do func(string) error) error {
	for line, hasNext := r.Next(); hasNext && line[0] != '|'; line, hasNext = r.Next() {
		if err := do(line); err != nil {
			return err
		}
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

func (r *Recipe) firstOrEmpty(tokens []string) string {
	return r.indexOrEmpty(tokens, 0)
}

func (*Recipe) indexOrEmpty(tokens []string, index int) string {
	if len(tokens) > index {
		return tokens[index]
	} else {
		return ""
	}
}

func (rec *Recipe) parseTargetKeyword(r *reader.Reader, tokens []string) error {
	if len(rec.Target) > 0 {
		return r.Error("duplicate target keyword")
	}

	typeStr := rec.indexOrEmpty(tokens, 0)

	switch typeStr {
	case "EXECUTABLE":
		rec.TargetType = TARGET_EXECUTABLE
	case "STATIC_LIBRARY":
		rec.TargetType = TARGET_STATIC_LIBRARY
	case "SHARED_LIBRARY":
		rec.TargetType = TARGET_SHARED_LIBRARY
	case "":
		return r.Error("missing or empty target type")
	default:
		return r.Errorf("invalid target type \"%s\"", typeStr)
	}

	rec.Target = rec.indexOrEmpty(tokens, 1)
	if len(rec.Target) == 0 {
		return r.Error("missing or empty target name")
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
	return rec.whileNotKeyword(r, func(line string) error {
		rec.addSource(filepath.Join(path, line))
		return nil
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
	return rec.whileNotKeyword(r, func(include string) error {
		rec.addInclude(include)
		return nil
	})
}

func (rec *Recipe) addInclude(include string) {
	if !filepath.IsAbs(include) {
		include = rec.JoinPath(include)
	}
	rec.Includes = append(rec.Includes, include)
}

func (rec *Recipe) parseLinkSharedLibsKeyword(r *reader.Reader, tokens []string) error {
	token := rec.firstOrEmpty(tokens)

	if !rec.peekExtraLine(r) {
		return rec.parseLinkSharedLibsSingle(r, token)
	} else {
		return rec.parseLinkSharedLibsMulti(r, token)
	}
}

func (rec *Recipe) parseLinkSharedLibsSingle(r *reader.Reader, lib string) error {
	if len(lib) == 0 {
		return r.Error("empty or missing library")
	}

	rec.addLinkedSharedLibrary(lib)
	return nil
}

func (rec *Recipe) parseLinkSharedLibsMulti(r *reader.Reader, path string) error {
	if len(path) > 0 {
		rec.LibraryPaths = append(rec.LibraryPaths, path)
	}

	return rec.whileNotKeyword(r, func(line string) error {
		rec.addLinkedSharedLibrary(line)
		return nil
	})
}

func (rec *Recipe) addLinkedSharedLibrary(lib string) {
	rec.LinkedSharedLibs = append(rec.LinkedSharedLibs, lib)
}

func (rec *Recipe) parseLinkStaticLibsKeyword(r *reader.Reader, tokens []string) error {
	token := rec.firstOrEmpty(tokens)

	if !rec.peekExtraLine(r) {
		return rec.parseLinkStaticLibsSingle(r, token)
	} else {
		return rec.parseLinkStaticLibsMulti(r, token)
	}
}

func (rec *Recipe) parseLinkStaticLibsSingle(r *reader.Reader, lib string) error {
	if len(lib) == 0 {
		return r.Error("empty or missing library")
	}

	rec.addLinkedStaticLibrary(lib)
	return nil
}

func (rec *Recipe) parseLinkStaticLibsMulti(r *reader.Reader, path string) error {
	return rec.whileNotKeyword(r, func(line string) error {
		rec.addLinkedStaticLibrary(filepath.Join(path, line))
		return nil
	})
}

func (rec *Recipe) addLinkedStaticLibrary(lib string) {
	if !filepath.IsAbs(lib) {
		lib = rec.JoinPath(lib)
	}
	rec.LinkedStaticLibs = append(rec.LinkedStaticLibs, lib)
}

func (rec *Recipe) parsePackageKeyword(r *reader.Reader, tokens []string) error {
	token := rec.firstOrEmpty(tokens)

	if !rec.peekExtraLine(r) {
		return rec.parsePackageSingle(r, token)
	} else {
		return rec.parsePackageMulti(r, token)
	}
}

func (rec *Recipe) parsePackageSingle(r *reader.Reader, pkg string) error {
	if len(pkg) == 0 {
		return r.Error("missing or empty package name")
	}

	return rec.parsePackageFile(r, pkg)
}

func (rec *Recipe) parsePackageMulti(r *reader.Reader, path string) error {
	return rec.whileNotKeyword(r, func(line string) error {
		return rec.parsePackageFile(r, filepath.Join(path, line))
	})
}

func (rec *Recipe) parsePackageFile(r *reader.Reader, pkg string) error {
	if !filepath.IsAbs(pkg) {
		pkg = rec.JoinPath(pkg)
	}

	file, err := os.Open(pkg)
	if err != nil {
		return r.Errorf("error opening package file \"%s\"", pkg)
	}
	defer file.Close()

	return rec.parsePackage(pkg, file)
}

func (rec *Recipe) parseLayerKeyword(r *reader.Reader, tokens []string) error {
	token := rec.firstOrEmpty(tokens)

	if !rec.peekExtraLine(r) {
		return rec.parseLayerSingle(r, token)
	} else {
		return rec.parseLayerMulti(r, token)
	}
}

func (rec *Recipe) parseLayerSingle(r *reader.Reader, layer string) error {
	if len(layer) == 0 {
		return r.Error("missing or empty layer name")
	}

	return rec.addLayer(r, layer)
}

func (rec *Recipe) parseLayerMulti(r *reader.Reader, path string) error {
	return rec.whileNotKeyword(r, func(line string) error {
		return rec.addLayer(r, filepath.Join(path, line))
	})
}

func (rec *Recipe) addLayer(r *reader.Reader, layer string) error {
	if !filepath.IsAbs(layer) {
		layer = rec.JoinPath(layer)
	}

	stat, err := os.Stat(layer)
	if err != nil {
		return r.Errorf("error opening layer file \"%s\"", layer)
	}

	if stat.IsDir() {
		layer = filepath.Join(layer, "recipe.txt")
	}

	rec.Layers = append(rec.Layers, layer)
	return nil
}
