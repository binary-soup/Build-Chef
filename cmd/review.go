package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/config"
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
	verify bool
}

func (cmd ReviewCommand) Run(cfg config.Config, args []string) error {
	path := cmd.pathFlag()
	verify := cmd.flagSet.Bool("verify", false, "verify filepaths are correct")
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.verify = *verify
	cmd.review(r, cfg.SystemPaths)

	return nil
}

func (cmd ReviewCommand) review(r *recipe.Recipe, systemPaths []string) {
	style.Header.Println("Executable:")
	cmd.reviewExecutable(r, r.Executable)
	cmd.reviewSourceFile(r, r.MainSource)

	style.Header.Println("Source Files:")
	for _, src := range r.SourceFiles {
		cmd.reviewSourceFile(r, src)
	}

	style.Header.Println("Include Directories:")
	for _, include := range r.Includes {
		cmd.reviewSystemPath(include)
	}

	style.Header.Println("Library Paths:")
	for _, path := range r.LibraryPaths {
		cmd.reviewSystemPath(path)
	}

	style.Header.Println("Libraries:")
	for _, lib := range r.Libraries {
		cmd.reviewLibrary(lib, append(systemPaths, r.LibraryPaths...))
	}
}

func (cmd ReviewCommand) reviewFilepath(verify func() bool, pass string, fail string, entry string) {
	fmt.Print(INDENT)

	if cmd.verify {
		if verify() {
			fmt.Print(pass)
		} else {
			fmt.Print(fail)
		}
	}

	fmt.Println(entry)
}

func (cmd ReviewCommand) reviewExecutable(r *recipe.Recipe, exec string) {
	cmd.reviewFilepath(
		cmd.verifyPath(r.JoinPath(filepath.Dir(exec))),
		cmd.verified(), cmd.notFound("path"),
		style.BoldCreate.String(exec),
	)
}

func (cmd ReviewCommand) reviewSourceFile(r *recipe.Recipe, src string) {
	cmd.reviewFilepath(
		cmd.verifyPath(r.JoinPath(src)),
		cmd.verified(), cmd.notFound("file"),
		style.File.String(src),
	)
}

func (cmd ReviewCommand) reviewSystemPath(path string) {
	cmd.reviewFilepath(
		cmd.verifyPath(path),
		cmd.verified(), cmd.notFound("path"),
		style.InfoV2.String(path),
	)
}

func (cmd ReviewCommand) verifyPath(path string) func() bool {
	return func() bool {
		return cmd.pathExists(path)
	}
}

func (ReviewCommand) pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (cmd ReviewCommand) reviewLibrary(lib string, paths []string) {
	cmd.reviewFilepath(
		cmd.verifyLibrary(paths, lib),
		cmd.verified(), cmd.unverified(),
		style.FileV2.String(lib),
	)
}

func (cmd ReviewCommand) verifyLibrary(paths []string, lib string) func() bool {
	return func() bool {
		return cmd.libraryExists(paths, lib)
	}
}

func (cmd ReviewCommand) libraryExists(paths []string, lib string) bool {
	for _, path := range paths {
		for _, ext := range []string{".a", ".so"} {
			if cmd.pathExists(filepath.Join(path, "lib"+lib+ext)) {
				return true
			}
		}
	}
	return false
}

func (ReviewCommand) verified() string {
	return style.Success.String("[verified] ")
}

func (ReviewCommand) unverified() string {
	return style.Info.String("[unverified] ")
}

func (ReviewCommand) notFound(name string) string {
	return style.Error.Format("[%s not found] ", name)
}
