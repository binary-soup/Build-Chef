package recipe

import (
	"io"

	"github.com/binary-soup/bchef/reader"
)

func (rec *Recipe) parseRecipe(file io.Reader) error {
	r := reader.New(file, rec.FullPath(), 1)

	err := rec.parse(&r, recipeParser{})
	if err != nil {
		return err
	}

	if len(rec.Executable) == 0 {
		return r.Error("missing executable keyword")
	}

	if len(rec.SourceFiles) == 0 {
		return r.Error("at least one source file required")
	}

	return nil
}

type recipeParser struct{}

func (recipeParser) ParseKeyword(rec *Recipe, r *reader.Reader, keyword string, tokens []string) error {
	switch keyword {
	case "EXECUTABLE":
		return rec.parseExecutableKeyword(r, tokens)
	case "SOURCES":
		return rec.parseSourcesKeyword(r, tokens)
	case "INCLUDES":
		return rec.parseIncludesKeyword(r, tokens)
	case "SHARED_LIBS":
		return rec.parseSharedLibsKeyword(r, tokens)
	case "STATIC_LIBS":
		return rec.parseStaticLibsKeyword(r, tokens)
	case "PACKAGE":
		return rec.parsePackageKeyword(r, tokens)
	default:
		return r.Errorf("unknown keyword \"%s\"", keyword)
	}
}
