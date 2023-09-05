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

	style.Info.Print(INDENT + "Name:   ")
	style.BoldCreate.Println(r.TrimPath(r.Executable))

	style.Info.Print(INDENT + "Source: ")
	style.File.Println(r.TrimPath(r.MainSource))

	style.Header.Println("Source Files:")
	for _, src := range r.SourceFiles {
		style.File.Println(INDENT + r.TrimPath(src))
	}
}
