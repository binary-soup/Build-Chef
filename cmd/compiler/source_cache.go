package compiler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/binary-soup/bchef/parser"
	"github.com/binary-soup/bchef/recipe"
)

func SourceCacheFile(path string) string {
	return filepath.Join(path, ".bchef/source_cache.txt")
}

type sourceInfo struct {
	Mod      int64
	Includes []string
}
type sourceCache map[string]sourceInfo

func (c sourceCache) Load(r *recipe.Recipe) {
	file, err := os.Open(SourceCacheFile(r.Path))
	if err != nil {
		return
	}

	defer file.Close()
	p := parser.New(file)

	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
		tokens := strings.Split(line, ",")

		includes := []string{}
		for _, token := range tokens[2:] {
			if len(token) == 0 {
				continue
			}
			includes = append(includes, r.JoinSourceDir(token))
		}

		c[r.JoinSourceDir(tokens[0])] = sourceInfo{
			Mod:      c.parseMod(tokens[1]),
			Includes: includes,
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

func (c sourceCache) Save(r *recipe.Recipe) {
	file, _ := os.Create(SourceCacheFile(r.Path))
	defer file.Close()

	for key, val := range c {
		fmt.Fprintf(file, "%s,%d,", r.TrimSourceDir(key), val.Mod)

		for _, include := range val.Includes {
			fmt.Fprintf(file, "%s,", r.TrimSourceDir(include))
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
