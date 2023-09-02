package compiler

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"

	"github.com/binary-soup/bchef/recipe"
)

var includeRegex = regexp.MustCompile(`^#include "([^"]+.(h|hxx))"$`)

func newTracker(dir string) sourceTracker {
	return sourceTracker{
		cache:    sourceCache{},
		dir:      dir,
		includes: map[string][]string{},
		mods:     map[string]int64{},
	}
}

type sourceTracker struct {
	cache    sourceCache
	dir      string
	includes map[string][]string
	mods     map[string]int64
}

func (t *sourceTracker) LoadCache(r *recipe.Recipe) {
	t.cache = sourceCache{}
	t.cache.Load(r)
}

func (t *sourceTracker) SaveCache(r *recipe.Recipe) {
	for file, includes := range t.includes {
		t.cache.UpdateEntry(file, t.getMod(file), includes)
	}
	t.cache.Save(r)
}

func (t sourceTracker) CalcChangedIndices(sources []string, objects []string) []int {
	indices := []int{}

	for i, src := range sources {
		if t.isFileNewer(src, t.getMod(objects[i])) {
			indices = append(indices, i)
		}
	}

	return indices
}

func (t sourceTracker) isFileNewer(file string, compare int64) bool {
	mod := t.getMod(file)
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

func (t sourceTracker) getMod(file string) int64 {
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

	includes, ok = t.cache.GetIncludes(file, t.getMod(file))
	if !ok {
		includes = t.parseIncludes(file)
	}

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
