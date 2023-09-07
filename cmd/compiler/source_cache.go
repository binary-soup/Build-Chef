package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/binary-soup/bchef/parser"
)

const SOURCE_CACHE_FILE = ".bchef/source_cache.txt"

type sourceInfo struct {
	Mod      int64
	Includes []string
}
type sourceCache map[string]sourceInfo

func (c sourceCache) Load(path string) {
	file, err := os.Open(c.cacheFile(path))
	if err != nil {
		return
	}

	defer file.Close()
	p := parser.New(file, 0)

	for line, hasNext := p.Next(); hasNext; line, hasNext = p.Next() {
		tokens := strings.Split(line, ",")

		includes := []string{}
		for _, token := range tokens[2:] {
			if len(token) == 0 {
				continue
			}
			includes = append(includes, token)
		}

		c[tokens[0]] = sourceInfo{
			Mod:      c.parseMod(tokens[1]),
			Includes: includes,
		}
	}
}

func (sourceCache) cacheFile(path string) string {
	return filepath.Join(path, SOURCE_CACHE_FILE)
}

func (c sourceCache) parseMod(token string) int64 {
	num, err := strconv.ParseInt(token, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func (c sourceCache) Save(path string) {
	file, _ := os.Create(c.cacheFile(path))
	defer file.Close()

	for key, val := range c {
		fmt.Fprintf(file, "%s,%d,", key, val.Mod)

		for _, include := range val.Includes {
			fmt.Fprintf(file, "%s,", include)
		}
		fmt.Fprintln(file)
	}
}

func (c sourceCache) UpdateEntry(file string, mod int64, includes []string) {
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
