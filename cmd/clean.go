package cmd

import (
	"os"

	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

type CleanCmd struct{}

func (CleanCmd) removeFile(file string, deleteStyle style.Style) {
	if err := os.Remove(file); err == nil {
		deleteStyle.Println("  x " + file)
	}
}

func (cmd CleanCmd) Run(r *recipe.Recipe) {
	style.Header.Println("Doing the Dishes...")

	for _, obj := range r.ObjectFiles {
		cmd.removeFile(obj, style.Delete)
	}
	cmd.removeFile(r.Name, style.BoldDelete)

	style.BoldSuccess.Println("Squeaky Clean!")
}
