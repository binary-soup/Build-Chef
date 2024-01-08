package recipe

import (
	"io"

	"github.com/binary-soup/bchef/reader"
)

func (rec *Recipe) parsePackage(pkg string, file io.Reader) error {
	r := reader.New(file, pkg, 1)

	err := rec.parse(&r, packageParser{})
	if err != nil {
		return err
	}
	return nil
}

type packageParser struct{}

func (packageParser) ParseKeyword(rec *Recipe, r *reader.Reader, keyword string, tokens []string) error {
	switch keyword {
	case "TARGET":
		return r.Error("\"TARGET\" unsupported in packages")
	case "SOURCES":
		return r.Error("\"SOURCES\" unsupported in packages")
	case "INCLUDES":
		return rec.parseIncludesKeyword(r, tokens)
	case "LINK_SHARED_LIBS":
		return rec.parseLinkSharedLibsKeyword(r, tokens)
	case "LINK_STATIC_LIBS":
		return rec.parseLinkStaticLibsKeyword(r, tokens)
	case "PACKAGES":
		return r.Error("sub-packages unsupported")
	case "LAYERS":
		return r.Error("\"LAYERS\" unsupported in packages")
	default:
		return r.Errorf("unknown keyword \"%s\"", keyword)
	}
}
