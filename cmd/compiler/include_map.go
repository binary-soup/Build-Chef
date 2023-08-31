package compiler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/binary-soup/bchef/parser"
)

const (
	INCLUDE_CACHE_FILE = ".bchef/include_cache.txt"
)

var includeRegex = regexp.MustCompile(`^#include "([^"]+.(h|hxx))"$`)

type sourceInfo struct {
	Mod      int64
	Includes []string
}
type sourceMap map[string]sourceInfo

func (m sourceMap) LoadCache() {
	file, err := os.Open(INCLUDE_CACHE_FILE)
	if err != nil {
		return
	}
	defer file.Close()

	p := parser.New(file)
	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
		tokens := strings.Split(line, ",")
		m[tokens[0]] = sourceInfo{
			Mod:      m.parseModTime(tokens[1]),
			Includes: tokens[2:],
		}
	}
}

func (sourceMap) parseModTime(mod string) int64 {
	num, err := strconv.ParseInt(mod, 10, 64)
	if err != nil {
		return 0
	}
	return num
}

func (m sourceMap) SaveCache() {
	file, _ := os.Create(INCLUDE_CACHE_FILE)
	defer file.Close()

	for key, val := range m {
		fmt.Fprintf(file, "%s,%d,%s\n", key, val.Mod, strings.Join(val.Includes, ","))
	}
}

func (m sourceMap) ParseSourceFile(src string, path string) {
	file, err := os.Open(src)
	if err != nil {
		return
	}
	defer file.Close()

	info := sourceInfo{
		Mod:      m.fileModTime(src),
		Includes: []string{},
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		match := includeRegex.FindStringSubmatch(scanner.Text())
		if match == nil {
			continue
		}
		info.Includes = append(info.Includes, filepath.Join(path, match[1]))
	}
	m[src] = info
}

func (sourceMap) fileModTime(file string) int64 {
	info, err := os.Stat(file)
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
}

func (m sourceMap) IsFileChanged(file string) bool {
	info, ok := m[file]
	if !ok {
		return true
	}

	return m.fileModTime(file) > info.Mod
}
