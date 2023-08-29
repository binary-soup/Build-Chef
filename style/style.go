package style

import (
	"fmt"
	"strings"
)

const (
	Bold      = "1"
	Underline = "4"
	Red       = "31"
	Green     = "32"
	Yellow    = "33"
	Blue      = "34"
	Magenta   = "35"
	Cyan      = "36"
)

type Style []string

func New(style ...string) Style {
	return style
}

func (s Style) Format(str string) string {
	if len(s) == 0 {
		return str
	}
	return "\033[" + strings.Join(s, ";") + "m" + str + "\033[0m"
}

func (s Style) Print(str string) {
	fmt.Print(s.Format(str))
}

func (s Style) Println(str string) {
	fmt.Println(s.Format(str))
}

func (s Style) Printf(format string, a ...any) {
	s.Print(fmt.Sprintf(format, a...))
}
