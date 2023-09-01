package compiler

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
)

const (
	INCLUDE_CACHE_FILE = ".bchef/include_cache.txt"
)

var includeRegex = regexp.MustCompile(`^#include "([^"]+.(h|hxx))"$`)

func newTracker(dir string) sourceTracker {
	return sourceTracker{
		dir:      dir,
		includes: map[string][]string{},
		mods:     map[string]int64{},
	}
}

type sourceTracker struct {
	dir      string
	includes map[string][]string
	mods     map[string]int64
}

// func (c *sourceTracker) Load() {
// 	file, err := os.Open(INCLUDE_CACHE_FILE)
// 	if err != nil {
// 		return
// 	}
// 	defer file.Close()

// 	p := parser.New(file)
// 	for line := p.NextLine(); len(line) > 0; line = p.NextLine() {
// 		tokens := strings.Split(line, ",")
// 		t.includes[tokens[0]] = tokens[1:]
// 	}
// }

// func (t sourceTracker) Save() {
// 	file, _ := os.Create(INCLUDE_CACHE_FILE)
// 	defer file.Close()

// 	for key, val := range t.includes {
// 		fmt.Fprintf(file, "%s,%s\n", key, strings.Join(val, ","))
// 	}
// }

func (t sourceTracker) NeedsCompiling(src string, obj string) bool {
	return t.isFileNewer(src, t.getFileMod(obj))
}

func (t sourceTracker) isFileNewer(file string, compare int64) bool {
	mod := t.getFileMod(file)
	if mod > compare {
		return true
	}

	for _, include := range t.getIncludes(file) {
		if newer := t.isFileNewer(include, compare); newer {
			return true
		}
	}
	return false
}

func (t sourceTracker) getFileMod(file string) int64 {
	mod, ok := t.mods[file]
	if ok {
		return mod
	}

	info, err := os.Stat(file)
	if err != nil {
		mod = 0
	} else {
		mod = info.ModTime().Unix()
	}

	t.mods[file] = mod
	return mod
}

func (t sourceTracker) getIncludes(file string) []string {
	includes, ok := t.includes[file]
	if ok {
		return includes
	}

	includes = t.parseIncludes(file)
	t.includes[file] = includes

	return includes
}

func (t sourceTracker) parseIncludes(src string) []string {
	includes := []string{}

	file, err := os.Open(src)
	if err != nil {
		return includes
	}

	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		match := includeRegex.FindStringSubmatch(scanner.Text())
		if match == nil {
			continue
		}
		includes = append(includes, filepath.Join(t.dir, match[1]))
	}

	return includes
}
