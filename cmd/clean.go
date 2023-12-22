package cmd

import (
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/cmd/compiler"
	"github.com/binary-soup/bchef/config"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

func NewCleanCommand() CleanCommand {
	return CleanCommand{
		command: newCommand("clean", "remove all created files"),
	}
}

type CleanCommand struct {
	command
}

func (cmd CleanCommand) Run(_ config.Config, args []string) error {
	path := cmd.pathFlag()
	cache := cmd.boolFlag("cache", false, "also delete cache files")
	cmd.parseFlags(args)

	t, err := cmd.loadRecipeTree(*path)
	if err != nil {
		return err
	}

	style.Header.Println("Doing the Dishes...")
	defer style.BoldSuccess.Println("Squeaky Clean!")

	t.Traverse(cleanVisitor{
		cache: *cache,
	})
	return nil
}

type cleanVisitor struct {
	cache bool
}

func (v cleanVisitor) Visit(r *recipe.Recipe, index int) bool {
	if v.cache {
		v.removeCacheFile(r, COMPILE_LOG_FILE)
		v.removeCacheFile(r, compiler.SOURCE_CACHE_FILE)
	}

	for _, obj := range r.ObjectFiles {
		v.removeObject(r, obj, true)
		v.removeObject(r, obj, false)
	}
	return true
}

func (v cleanVisitor) removeObject(r *recipe.Recipe, obj string, debug bool) {
	v.removeFile(r.ObjectPath, filepath.Join(r.GetMode(debug), obj), style.Delete)
}

func (v cleanVisitor) removeCacheFile(r *recipe.Recipe, file string) {
	v.removeFile(r.Path, file, style.Delete)
}

func (cleanVisitor) removeFile(path string, file string, deleteStyle style.Style) {
	if err := os.Remove(filepath.Join(path, file)); err == nil {
		deleteStyle.Println(INDENT + "x " + file)
	}
}
