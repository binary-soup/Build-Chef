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
	verify := cmd.boolFlag("verify", false, "verify filepaths are correct")
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
	style.Header.Println("Target:")
	cmd.reviewTarget(r)

	style.Header.Println("Source Files:")
	for _, src := range r.SourceFiles {
		cmd.reviewSourceFile(r, src)
	}

	style.Header.Println("Include Directories:")
	for _, include := range r.Includes {
		cmd.reviewSystemPath(include)
	}

	style.Header.Println("Shared Library Paths:")
	for _, path := range r.LibraryPaths {
		cmd.reviewSystemPath(path)
	}

	style.Header.Println("Shared Libraries:")
	for _, lib := range r.LinkedSharedLibs {
		cmd.reviewSharedLibrary(lib, append(systemPaths, r.LibraryPaths...))
	}

	style.Header.Println("Static Libraries:")
	for _, lib := range r.LinkedStaticLibs {
		cmd.reviewStaticLibrary(lib)
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

func (cmd ReviewCommand) reviewTarget(r *recipe.Recipe) {
	cmd.reviewFilepath(
		cmd.verifyPath(r.JoinPath(filepath.Dir(r.Target))),
		cmd.verified(), cmd.notFound("path"),
		style.BoldCreate.Format("%s (%s)", r.Target, r.GetTargetType()),
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

func (cmd ReviewCommand) reviewSharedLibrary(lib string, paths []string) {
	cmd.reviewFilepath(
		cmd.verifySharedLibrary(paths, lib),
		cmd.verified(), cmd.unverified(),
		style.FileV2.String(lib),
	)
}

func (cmd ReviewCommand) verifySharedLibrary(paths []string, lib string) func() bool {
	return func() bool {
		return cmd.sharedLibraryExists(paths, lib)
	}
}

func (cmd ReviewCommand) sharedLibraryExists(paths []string, lib string) bool {
	for _, path := range paths {
		if cmd.pathExists(filepath.Join(path, "lib"+lib+".so")) {
			return true
		}
	}
	return false
}

func (cmd ReviewCommand) reviewStaticLibrary(lib string) {
	cmd.reviewFilepath(
		cmd.verifyPath(lib),
		cmd.verified(), cmd.notFound("library"),
		style.FileV2.String(lib),
	)
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
