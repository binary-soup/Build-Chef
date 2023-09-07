package cmd

import (
	"errors"
	"flag"
	"fmt"
	"path/filepath"

	"github.com/binary-soup/bchef/config"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

const (
	INDENT = "  "
)

func newCommand(name string, description string) command {
	return command{
		Name:        name,
		Description: description,
		flagSet:     flag.NewFlagSet(name, flag.ExitOnError),
	}
}

type Command interface {
	GetName() string
	PrintUsage()
	Run(cfg config.Config, args []string) error
}

type command struct {
	Name        string
	Description string
	flagSet     *flag.FlagSet
}

func (cmd command) GetName() string {
	return cmd.Name
}

func (cmd command) PrintUsage() {
	fmt.Printf("%s ~ %s\n", style.BoldInfoV2.String(cmd.Name), cmd.Description)
}

func (cmd command) parseFlags(args []string) {
	cmd.flagSet.Usage = func() {
		cmd.PrintUsage()
		style.BoldFileV2.Println("Options:")
		cmd.flagSet.PrintDefaults()
	}
	cmd.flagSet.Parse(args)
}

func (cmd command) pathFlag() *string {
	return cmd.flagSet.String("path", "recipe.txt", "path to the recipe file")
}

func (cmd command) loadRecipe(path string) (*recipe.Recipe, error) {
	r, err := recipe.Load(path)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error loading recipe at %s", style.FileV2.String(path)), err)
	}

	fmt.Println("Recipe loaded from", style.BoldFileV2.String(filepath.Join(r.Path, r.Name)))
	return r, nil
}
