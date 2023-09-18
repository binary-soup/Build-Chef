package compiler

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"

	"github.com/binary-soup/bchef/recipe"
)

var includeRegex = regexp.MustCompile(`^#include "([^"]+.(h|hxx))"$`)

func newTracker(r *recipe.Recipe) sourceTracker {
	return sourceTracker{
		cache:    sourceCache{},
		recipe:   r,
		includes: map[string][]string{},
		mods:     map[string]int64{},
	}
}

type sourceTracker struct {
	cache    sourceCache
	recipe   *recipe.Recipe
	includes map[string][]string
	mods     map[string]int64
}

func (t *sourceTracker) LoadCache() {
	t.cache = sourceCache{}
	t.cache.Load(t.recipe.Path)
}

func (t *sourceTracker) SaveCache() {
	for file, includes := range t.includes {
		t.cache.UpdateEntry(file, t.getMod(file), includes)
	}
	t.cache.Save(t.recipe.Path)
}

func (t sourceTracker) CalcChangedIndices(sources []string, objects []string, objPath string) []int {
	indices := []int{}

	for i, src := range sources {
		if t.isFileNewer(src, t.getMod(filepath.Join(objPath, objects[i]))) {
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

	info, err := os.Stat(t.recipe.JoinPath(file))
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
		includes = t.parseIncludes(t.recipe.JoinPath(file))
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

		include, res := t.resolveInclude(match[1])
		if res {
			includes = append(includes, include)
		}
	}
	return includes
}

func (t sourceTracker) resolveInclude(include string) (string, bool) {
	for _, path := range t.recipe.Includes {
		full := filepath.Join(t.recipe.Path, path, include)

		if t.pathExists(full) {
			return full, true
		}
	}

	return "", false
}

func (sourceTracker) pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
