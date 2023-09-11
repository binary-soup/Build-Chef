package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/binary-soup/bchef/common"
	"github.com/binary-soup/bchef/config"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

const HEADER_TMPL = `
#ifndef {{.}}
#define {{.}}



#endif
`

func NewMakeCommand() MakeCommand {
	return MakeCommand{
		command: newCommand("make", "make a new header/source"),
	}
}

type MakeCommand struct {
	command
	overwrite bool
}

func (cmd MakeCommand) Run(_ config.Config, args []string) error {
	path := cmd.pathFlag()
	name := cmd.flagSet.String("name", "", "name of the new file")
	overwrite := cmd.flagSet.Bool("overwrite", false, "overwrite files if already exists (no undo)")
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.overwrite = *overwrite
	return cmd.create(r, *name)
}

func (cmd MakeCommand) create(r *recipe.Recipe, name string) error {
	if len(name) == 0 {
		return errors.New("missing or empty name")
	}
	name += ".hxx"

	style.Header.Println("Making...")

	if err := cmd.createHeader(r.JoinPath(name), name); err != nil {
		return err
	}

	style.BoldSuccess.Println("New Ingredients Ready!")
	return nil
}

func (cmd MakeCommand) createHeader(path string, name string) error {
	if err := cmd.checkFileExists(path); err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Join(fmt.Errorf("error creating \"%s\"", path), err)
	}
	defer file.Close()

	cmd.fillHeader(file, name)
	style.Create.Println(INDENT + "+ " + path)

	return nil
}

func (cmd MakeCommand) fillHeader(file *os.File, name string) {
	headerGuards := common.ToUpper(common.ReplaceChar(name, "/.", '_'))
	cmd.parseTemplate("header", HEADER_TMPL).Execute(file, headerGuards)
}

func (cmd MakeCommand) checkFileExists(path string) error {
	if cmd.overwrite {
		return nil
	}

	_, err := os.Stat(path)
	if err == nil {
		return fmt.Errorf("file \"%s\" already exists", path)
	}

	return nil
}

func (MakeCommand) parseTemplate(name string, tmpl string) *template.Template {
	return template.Must(template.New(name).Parse(strings.TrimLeft(tmpl, "\n")))
}
