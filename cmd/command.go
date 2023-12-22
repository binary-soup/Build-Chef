package cmd

import (
	"errors"
	"flag"
	"fmt"
	"os"
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

func (cmd command) boolFlag(name string, value bool, usage string) *bool {
	return cmd.flagSet.Bool(name, value, usage)
}

func (cmd command) stringFlag(name string, value string, usage string) *string {
	return cmd.flagSet.String(name, value, usage)
}

func (cmd command) pathFlag() *string {
	return cmd.stringFlag("path", "recipe.txt", "path to the recipe file")
}

func (cmd command) loadRecipeTree(path string) (*recipe.RecipeTree, error) {
	stat, err := os.Stat(path)
	if err != nil || stat.IsDir() {
		path = filepath.Join(path, "recipe.txt")
	}

	t, err := recipe.LoadTree(path)
	if err != nil {
		return nil, errors.Join(fmt.Errorf("error loading recipe at %s", style.FileV2.String(path)), err)
	}

	fmt.Println("Recipe tree loaded from", style.BoldFileV2.String(t.Root.FullPath()))
	return t, nil
}
