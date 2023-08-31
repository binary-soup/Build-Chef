package compiler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/binary-soup/bchef/parser"
)

const (
	INCLUDE_CACHE_FILE = ".bchef/include_cache.txt"
)

var includeRegex = regexp.MustCompile(`^#include "([^"]+.(h|hxx))"$`)

type includeMap map[string][]string

func (m includeMap) LoadCache() {
	file, err := os.Open(INCLUDE_CACHE_FILE)
	if err != nil {
		return
	}
	defer file.Close()

	p := parser.New(file)
	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
		tokens := strings.Split(line, ",")
		m[tokens[0]] = tokens[1:]
	}
}

func (m includeMap) SaveCache() {
	file, _ := os.Create(INCLUDE_CACHE_FILE)
	defer file.Close()

	for key, val := range m {
		fmt.Fprintf(file, "%s,%s\n", key, strings.Join(val, ","))
	}
}

func (m includeMap) ParseSourceFile(src string, path string) {
	file, err := os.Open(src)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		match := includeRegex.FindStringSubmatch(scanner.Text())
		if match == nil {
			continue
		}

		m[src] = append(m[src], filepath.Join(path, match[1]))
	}
}
