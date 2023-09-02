package cmd

import (
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

func NewInfoCommand() InfoCommand {
	return InfoCommand{
		command: newCommand("info", "print info about a recipe"),
	}
}

type InfoCommand struct {
	command
}

func (cmd InfoCommand) Run(args []string) error {
	path := cmd.pathFlag()
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.info(r)
	return nil
}

func (cmd InfoCommand) info(r *recipe.Recipe) {
	style.Header.Println("Executable:")
	style.BoldCreate.Println(INDENT + r.Executable)

	style.Header.Println("Source Files:")
	for _, src := range r.SourceFiles {
		style.File.Println(INDENT + src)
	}
}
