package recipe

import (
	"bufio"
	"io"
)

func newParser(r io.Reader) parser {
	return parser{
		Scanner: bufio.NewScanner(r),
	}
}

type parser struct {
	Scanner *bufio.Scanner
}

func (p parser) NextLine() string {
	for p.Scanner.Scan() {
		line := p.Scanner.Text()

		if len(line) > 0 && line[0] != '#' {
			return line
		}
	}

	return ""
}
