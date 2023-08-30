package cmd

import (
	"os"

	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

type CleanCmd struct{}

func (CleanCmd) removeFile(r *recipe.Recipe, file string, deleteStyle style.Style) {
	if err := os.Remove(file); err == nil {
		deleteStyle.Println(INDENT + "x " + r.TrimObjectDir(file))
	}
}

func (cmd CleanCmd) Run(r *recipe.Recipe) bool {
	style.Header.Println("Doing the Dishes...")

	for _, obj := range r.ObjectFiles {
		cmd.removeFile(r, obj, style.Delete)
	}
	cmd.removeFile(r, r.Name, style.BoldDelete)

	style.BoldSuccess.Println("Squeaky Clean!")
	return true
}
