package cmd

import (
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/cmd/compiler"
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

func (cmd CleanCommand) Run(args []string) error {
	path := cmd.pathFlag()
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.clean(r)
	return nil
}

func (cmd CleanCommand) clean(r *recipe.Recipe) {
	style.Header.Println("Doing the Dishes...")

	for _, obj := range r.ObjectFiles {
		cmd.removeObject(r, obj, true)
		cmd.removeObject(r, obj, false)
	}

	cmd.removeExecutable(r, true)
	cmd.removeExecutable(r, false)

	os.Remove(CompileLogFile(r.Path))
	os.Remove(compiler.SourceCacheFile(r.Path))

	style.BoldSuccess.Println("Squeaky Clean!")
}

func (cmd CleanCommand) removeObject(r *recipe.Recipe, obj string, debug bool) {
	cmd.removeFile(r, r.ObjectPath, filepath.Join(r.GetMode(debug), obj), style.Delete)
}

func (cmd CleanCommand) removeExecutable(r *recipe.Recipe, debug bool) {
	cmd.removeFile(r, r.Path, r.GetExecutable(debug), style.BoldDelete)
}

func (CleanCommand) removeFile(r *recipe.Recipe, path string, file string, deleteStyle style.Style) {
	if err := os.Remove(filepath.Join(path, file)); err == nil {
		deleteStyle.Println(INDENT + "x " + file)
	}
}
