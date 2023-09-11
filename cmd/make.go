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

const HEADER_TEMPLATE = `
#ifndef {{.}}
#define {{.}}



#endif
`

const SOURCE_TEMPLATE = `
#include "{{.}}"


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

	style.Header.Println("Making...")

	header := name + ".hxx"
	if err := cmd.createHeader(r.JoinPath(header), header); err != nil {
		return err
	}

	if err := cmd.createSource(r.JoinPath(name+".cxx"), header); err != nil {
		return err
	}

	style.BoldSuccess.Println("New Ingredients Ready!")
	return nil
}

func (cmd MakeCommand) createHeader(path string, name string) error {
	return cmd.createFile(path, func(file *os.File) {
		headerGuards := common.ToUpper(common.ReplaceChar(name, "/.", '_'))
		cmd.parseTemplate("header", HEADER_TEMPLATE).Execute(file, headerGuards)
	})
}

func (cmd MakeCommand) createSource(path string, header string) error {
	return cmd.createFile(path, func(file *os.File) {
		cmd.parseTemplate("source", SOURCE_TEMPLATE).Execute(file, header)
	})
}

func (cmd MakeCommand) createFile(path string, fill func(*os.File)) error {
	if !cmd.canCreateFile(path) {
		style.File.Printf("%s- %s (already exists)\n", INDENT, path)
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Join(fmt.Errorf("error creating \"%s\"", path), err)
	}
	defer file.Close()

	fill(file)
	style.Create.Println(INDENT + "+ " + path)
	return nil
}

func (cmd MakeCommand) canCreateFile(path string) bool {
	if cmd.overwrite {
		return true
	}

	_, err := os.Stat(path)
	return err != nil
}

func (MakeCommand) parseTemplate(name string, tmpl string) *template.Template {
	return template.Must(template.New(name).Parse(strings.TrimLeft(tmpl, "\n")))
}
