package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/binary-soup/bchef/parser"
)

func SourceCacheFile(path string) string {
	return filepath.Join(path, ".bchef/source_cache.txt")
}

type sourceInfo struct {
	Mod      int64
	Includes []string
}
type sourceCache map[string]sourceInfo

func (c sourceCache) Load(path string) {
	file, err := os.Open(SourceCacheFile(path))
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

func (c sourceCache) Save(path string) {
	file, _ := os.Create(SourceCacheFile(path))
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
