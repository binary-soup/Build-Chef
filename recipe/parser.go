package recipe

import (
	"strings"

	"github.com/binary-soup/bchef/reader"
)

type parser interface {
	ParseKeyword(rec *Recipe, r *reader.Reader, keyword string, tokens []string) error
}

func (rec *Recipe) parse(r *reader.Reader, p parser) error {
	for line, hasNext := r.Next(); hasNext; line, hasNext = r.Next() {
		if line[0] != '|' {
			continue
		}
		tokens := strings.Split(line[1:], " ")

		err := p.ParseKeyword(rec, r, tokens[0], tokens[1:])
		if err != nil {
			return err
		}
	}
	return nil
}
