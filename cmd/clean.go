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

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.clean(r, *cache)
	return nil
}

func (cmd CleanCommand) clean(r *recipe.Recipe, cache bool) {
	style.Header.Println("Doing the Dishes...")

	if cache {
		cmd.removeCacheFile(r, COMPILE_LOG_FILE)
		cmd.removeCacheFile(r, compiler.SOURCE_CACHE_FILE)
	}

	for _, obj := range r.ObjectFiles {
		cmd.removeObject(r, obj, true)
		cmd.removeObject(r, obj, false)
	}

	style.BoldSuccess.Println("Squeaky Clean!")
}

func (cmd CleanCommand) removeObject(r *recipe.Recipe, obj string, debug bool) {
	cmd.removeFile(r.ObjectPath, filepath.Join(r.GetMode(debug), obj), style.Delete)
}

func (cmd CleanCommand) removeCacheFile(r *recipe.Recipe, file string) {
	cmd.removeFile(r.Path, file, style.Delete)
}

func (CleanCommand) removeFile(path string, file string, deleteStyle style.Style) {
	if err := os.Remove(filepath.Join(path, file)); err == nil {
		deleteStyle.Println(INDENT + "x " + file)
	}
}
