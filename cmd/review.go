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
}

func (cmd ReviewCommand) Run(cfg config.Config, args []string) error {
	path := cmd.pathFlag()
	verify := cmd.boolFlag("verify", false, "verify filepaths are correct")
	cmd.parseFlags(args)

	t, err := cmd.loadRecipeTree(*path)
	if err != nil {
		return err
	}

	t.Traverse(reviewVisitor{
		verify:      *verify,
		systemPaths: cfg.SystemPaths,
	})
	return nil
}

type reviewVisitor struct {
	verify      bool
	systemPaths []string
}

func (v reviewVisitor) Visit(r *recipe.Recipe, index int) bool {
	style.BoldFileV2.Printf("[%d] %s:\n", index+1, r.FullPath())

	fmt.Print(INDENT)
	style.Header.Println("Target:")
	v.reviewTarget(r)

	fmt.Print(INDENT)
	style.Header.Println("Layers:")
	for _, layer := range r.Layers {
		v.reviewLayer(r, layer)
	}

	fmt.Print(INDENT)
	style.Header.Println("Source Files:")
	for _, src := range r.SourceFiles {
		v.reviewSourceFile(r, src)
	}

	fmt.Print(INDENT)
	style.Header.Println("Include Directories:")
	for _, include := range r.Includes {
		v.reviewSystemPath(include)
	}

	fmt.Print(INDENT)
	style.Header.Println("Shared Library Paths:")
	for _, path := range r.LibraryPaths {
		v.reviewSystemPath(path)
	}

	fmt.Print(INDENT)
	style.Header.Println("Shared Libraries:")
	for _, lib := range r.LinkedSharedLibs {
		v.reviewSharedLibrary(lib, append(v.systemPaths, r.LibraryPaths...))
	}

	fmt.Print(INDENT)
	style.Header.Println("Static Libraries:")
	for _, lib := range r.LinkedStaticLibs {
		v.reviewStaticLibrary(lib)
	}

	return true
}

func (v reviewVisitor) reviewFilepath(verify func() bool, pass string, fail string, entry string) {
	fmt.Print(INDENT, INDENT)

	if v.verify {
		if verify() {
			fmt.Print(pass)
		} else {
			fmt.Print(fail)
		}
	}

	fmt.Println(entry)
}

func (v reviewVisitor) reviewTarget(r *recipe.Recipe) {
	v.reviewFilepath(
		v.verifyPath(r.JoinPath(filepath.Dir(r.Target))),
		v.verified(), v.notFound("path"),
		style.BoldCreate.Format("%s (%s)", r.Target, r.GetTargetType()),
	)
}

func (v reviewVisitor) reviewLayer(r *recipe.Recipe, layer string) {
	v.reviewFilepath(
		v.verifyPath(layer),
		v.verified(), v.notFound("file"),
		style.FileV2.String(layer),
	)
}

func (v reviewVisitor) reviewSourceFile(r *recipe.Recipe, src string) {
	v.reviewFilepath(
		v.verifyPath(r.JoinPath(src)),
		v.verified(), v.notFound("file"),
		style.File.String(src),
	)
}

func (v reviewVisitor) reviewSystemPath(path string) {
	v.reviewFilepath(
		v.verifyPath(path),
		v.verified(), v.notFound("path"),
		style.InfoV2.String(path),
	)
}

func (v reviewVisitor) verifyPath(path string) func() bool {
	return func() bool {
		return v.pathExists(path)
	}
}

func (reviewVisitor) pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (v reviewVisitor) reviewSharedLibrary(lib string, paths []string) {
	v.reviewFilepath(
		v.verifySharedLibrary(paths, lib),
		v.verified(), v.unverified(),
		style.FileV2.String(lib),
	)
}

func (v reviewVisitor) verifySharedLibrary(paths []string, lib string) func() bool {
	return func() bool {
		return v.sharedLibraryExists(paths, lib)
	}
}

func (v reviewVisitor) sharedLibraryExists(paths []string, lib string) bool {
	for _, path := range paths {
		if v.pathExists(filepath.Join(path, "lib"+lib+".so")) {
			return true
		}
	}
	return false
}

func (v reviewVisitor) reviewStaticLibrary(lib string) {
	v.reviewFilepath(
		v.verifyPath(lib),
		v.verified(), v.notFound("library"),
		style.FileV2.String(lib),
	)
}

func (reviewVisitor) verified() string {
	return style.Success.String("[verified] ")
}

func (reviewVisitor) unverified() string {
	return style.Info.String("[unverified] ")
}

func (reviewVisitor) notFound(name string) string {
	return style.Error.Format("[%s not found] ", name)
}
