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
	case "EXECUTABLE":
		return r.Error("\"EXECUTABLE\" unsupported in packages")
	case "SOURCES":
		return r.Error("\"SOURCES\" unsupported in packages")
	case "INCLUDES":
		return rec.parseIncludesKeyword(r, tokens)
	case "SHARED_LIBS":
		return rec.parseSharedLibsKeyword(r, tokens)
	case "STATIC_LIBS":
		return rec.parseStaticLibsKeyword(r, tokens)
	case "PACKAGE":
		return r.Error("sub-packages unsupported")
	default:
		return r.Errorf("unknown keyword \"%s\"", keyword)
	}
}
