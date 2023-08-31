package parser

import (
	"bufio"
	"io"
)

func New(r io.Reader) Parser {
	return Parser{
		Scanner: bufio.NewScanner(r),
	}
}

type Parser struct {
	Scanner *bufio.Scanner
}

func (p Parser) NextLine() string {
	for p.Scanner.Scan() {
		line := p.Scanner.Text()

		if len(line) > 0 && line[0] != '#' {
			return line
		}
	}

	return ""
}
