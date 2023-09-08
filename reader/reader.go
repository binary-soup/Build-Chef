package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/binary-soup/bchef/style"
)

func New(r io.Reader, name string, historyLength int) Reader {
	return Reader{
		scanner:    bufio.NewScanner(r),
		name:       name,
		line:       0,
		historyPtr: 0,
		history:    NewLoopedList[string](historyLength),
	}
}

type Reader struct {
	scanner    *bufio.Scanner
	name       string
	line       int
	historyPtr int
	history    LoopedList[string]
}

func (r *Reader) Rewind(count int) {
	if count < 0 {
		count = 0
	} else if count > r.history.Size() {
		count = r.history.Size()
	}
	r.historyPtr = -count
}

func (r *Reader) Next() (string, bool) {
	if r.historyPtr < 0 {
		r.historyPtr += 1
		return r.history.Get(r.historyPtr - 1), true
	}

	for r.scanner.Scan() {
		r.line++
		line := strings.TrimLeft(r.scanner.Text(), " \t")

		if len(line) > 0 && line[0] != '#' {
			r.history.Push(line)
			return line, true
		}
	}

	return "", false
}

func (r Reader) Error(message string) error {
	return errors.New(style.BoldError.Format("[%s:%d] ", r.name, r.line) + message)
}

func (r Reader) Errorf(format string, a ...any) error {
	return r.Error(fmt.Sprintf(format, a...))
}
