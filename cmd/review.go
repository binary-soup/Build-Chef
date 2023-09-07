package cmd

import (
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

func NewReviewCommand() ReviewCommand {
	return ReviewCommand{
		command: newCommand("review", "print details about a recipe"),
	}
}

type ReviewCommand struct {
	command
}

func (cmd ReviewCommand) Run(args []string) error {
	path := cmd.pathFlag()
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.info(r)
	return nil
}

func (cmd ReviewCommand) info(r *recipe.Recipe) {
	style.Header.Println("Executable:")

	style.BoldCreate.Println(INDENT + r.Executable)
	style.File.Println(INDENT + r.MainSource)

	style.Header.Println("Source Files:")
	for _, src := range r.SourceFiles {
		style.File.Println(INDENT + src)
	}

	style.Header.Println("Include Directories:")
	for _, include := range r.Includes {
		style.InfoV2.Println(INDENT + include)
	}

	style.Header.Println("Library Paths:")
	for _, path := range r.LibraryPaths {
		style.InfoV2.Println(INDENT + path)
	}

	style.Header.Println("Libraries:")
	for _, lib := range r.Libraries {
		style.FileV2.Println(INDENT + lib)
	}
}
