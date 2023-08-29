package cmd

import (
	"os"

	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

type CleanCmd struct{}

func (CleanCmd) Run(r *recipe.Recipe) error {
	style.Header.Println("Doing the Dishes...")

	err := os.Remove(r.Name)
	if err == nil {
		style.BoldDelete.Println("  x " + r.Name)
	}

	return nil
}
