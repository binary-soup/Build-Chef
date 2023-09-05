package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/binary-soup/bchef/style"
)

func New(r io.Reader, historyLength int) Parser {
	return Parser{
		scanner:    bufio.NewScanner(r),
		line:       0,
		historyPtr: 0,
		history:    NewLoopedList[string](historyLength),
	}
}

type Parser struct {
	scanner    *bufio.Scanner
	line       int
	historyPtr int
	history    LoopedList[string]
}

func (p *Parser) Rewind(count int) {
	if count < 0 {
		count = 0
	} else if count > p.history.Size() {
		count = p.history.Size()
	}
	p.historyPtr = -count
}

func (p *Parser) Next() (string, bool) {
	if p.historyPtr < 0 {
		p.historyPtr += 1
		return p.history.Get(p.historyPtr - 1), true
	}

	for p.scanner.Scan() {
		p.line++
		line := strings.TrimLeft(p.scanner.Text(), " \t")

		if len(line) > 0 && line[0] != '#' {
			p.history.Push(line)
			return line, true
		}
	}

	return "", false
}

func (p Parser) Error(message string) error {
	return errors.New(style.BoldError.Format("[%d]: ", p.line) + message)
}

func (p Parser) Errorf(format string, a ...any) error {
	return p.Error(fmt.Sprintf(format, a...))
}
