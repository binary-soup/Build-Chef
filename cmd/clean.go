package cmd

import (
	"os"

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
		cmd.removeFile(r, obj, style.Delete)
	}
	cmd.removeFile(r, r.Executable, style.BoldDelete)

	os.Remove(CompileLogFile(r.Path))
	os.Remove(compiler.SourceCacheFile(r.Path))

	style.BoldSuccess.Println("Squeaky Clean!")
}

func (CleanCommand) removeFile(r *recipe.Recipe, file string, deleteStyle style.Style) {
	if err := os.Remove(file); err == nil {
		deleteStyle.Println(INDENT + "x " + r.TrimObjectDir(file))
	}
}
