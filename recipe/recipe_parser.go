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

	if len(rec.Target) == 0 {
		return r.Error("missing target keyword")
	}

	if len(rec.SourceFiles) == 0 {
		return r.Error("at least one source file required")
	}

	return nil
}

type recipeParser struct{}

func (recipeParser) ParseKeyword(rec *Recipe, r *reader.Reader, keyword string, tokens []string) error {
	switch keyword {
	case "TARGET":
		return rec.parseTargetKeyword(r, tokens)
	case "SOURCES":
		return rec.parseSourcesKeyword(r, tokens)
	case "INCLUDES":
		return rec.parseIncludesKeyword(r, tokens)
	case "LINK_SHARED_LIBS":
		return rec.parseLinkSharedLibsKeyword(r, tokens)
	case "LINK_STATIC_LIBS":
		return rec.parseLinkStaticLibsKeyword(r, tokens)
	case "PACKAGE":
		return rec.parsePackageKeyword(r, tokens)
	case "LAYER":
		return rec.parseLayerKeyword(r, tokens)
	default:
		return r.Errorf("unknown keyword \"%s\"", keyword)
	}
}
