package compiler

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"

	"github.com/binary-soup/bchef/recipe"
)

var includeRegex = regexp.MustCompile(`^\s*#include\s+"([^"]+.(h|hxx))"\s*$`)

func newTracker(r *recipe.Recipe, debug bool) sourceTracker {
	return sourceTracker{
		cache:    sourceCache{},
		recipe:   r,
		debug:    debug,
		includes: map[string][]string{},
		mods:     map[string]int64{},
	}
}

type sourceTracker struct {
	cache    sourceCache
	recipe   *recipe.Recipe
	debug    bool
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

func (t sourceTracker) CalcChangedIndices(target string) ([]int, bool) {
	indices := []int{}
	var maxMod int64 = 0

	for i, src := range t.recipe.SourceFiles {
		mod := t.getMod(t.recipe.JoinObjectDir(t.recipe.ObjectFiles[i], t.debug))
		maxMod = t.maxMod(maxMod, mod)

		if t.isFileNewer(src, mod) {
			indices = append(indices, i)
		}
	}

	return indices, t.getModByPath(target) < maxMod
}

func (t sourceTracker) maxMod(mod1 int64, mod2 int64) int64 {
	if mod1 >= mod2 {
		return mod1
	}
	return mod2
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

	mod = t.getModByPath(t.recipe.JoinPath(file))

	t.mods[file] = mod
	return mod
}

func (t sourceTracker) getModByPath(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.ModTime().Unix()
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
