package cmd

import (
	"fmt"
	"os"
	"path/filepath"

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
	verify := cmd.flagSet.Bool("verify", false, "verify filepaths are correct")
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.info(r, *verify)
	return nil
}

func (cmd ReviewCommand) info(r *recipe.Recipe, verify bool) {
	style.Header.Println("Executable:")
	cmd.reviewExecutable(r, r.Executable, verify)
	cmd.reviewSourceFile(r, r.MainSource, verify)

	style.Header.Println("Source Files:")
	for _, src := range r.SourceFiles {
		cmd.reviewSourceFile(r, src, verify)
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

func (cmd ReviewCommand) reviewExecutable(r *recipe.Recipe, exec string, verify bool) {
	cmd.reviewFile(style.BoldCreate, exec, r.JoinPath(filepath.Dir(exec)), "path", verify)
}

func (cmd ReviewCommand) reviewSourceFile(r *recipe.Recipe, src string, verify bool) {
	cmd.reviewFile(style.File, src, r.JoinPath(src), "file", verify)
}

func (cmd ReviewCommand) reviewFile(style style.Style, file string, filepath string, name string, verify bool) {
	fmt.Print(INDENT)

	if verify {
		cmd.verifyPath(filepath, name)
	}
	style.Println(file)
}

func (ReviewCommand) verifyPath(path string, name string) {
	_, err := os.Stat(path)

	if err == nil {
		style.Success.Print("[verified] ")
	} else {
		style.Error.Printf("[%s not found] ", name)
	}
}
