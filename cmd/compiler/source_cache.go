package compiler

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/binary-soup/bchef/parser"
)

const (
	SOURCE_CACHE_FILE = ".bchef/source_cache.txt"
)

type sourceInfo struct {
	Mod      int64
	Includes []string
}
type sourceCache map[string]sourceInfo

func (c sourceCache) Load() {
	file, err := os.Open(SOURCE_CACHE_FILE)
	if err != nil {
		return
	}
	defer file.Close()

	p := parser.New(file)
	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
		tokens := strings.Split(line, ",")

		c[tokens[0]] = sourceInfo{
			Mod:      c.parseMod(tokens[1]),
			Includes: tokens[2:],
		}
	}
}

func (c sourceCache) parseMod(token string) int64 {
	num, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func (c sourceCache) Save() {
	file, _ := os.Create(SOURCE_CACHE_FILE)
	defer file.Close()

	for key, val := range c {
		fmt.Fprintf(file, "%s,%d,%s\n", key, val.Mod, strings.Join(val.Includes, ","))
	}
}

func (c sourceCache) UpdateEntry(file string, mod int64, includes []string) {
	if len(file) == 0 {
		return
	}

	c[file] = sourceInfo{
		Mod:      mod,
		Includes: includes,
	}
}

func (c sourceCache) GetIncludes(file string, mod int64) ([]string, bool) {
	info, ok := c[file]
	if !ok {
		return []string{}, false
	}

	return info.Includes, mod <= info.Mod
}
