package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/cmd/compiler"
	"github.com/binary-soup/bchef/cmd/compiler/gxx"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

// TODO: fix command injection

func CompileLogFile(path string) string {
	return filepath.Join(path, ".bchef/compile_log.txt")
}

func NewCookCommand() CookCommand {
	return CookCommand{
		command: newCommand("cook", "compile the project"),
	}
}

type CookCommand struct {
	command
}

func (cmd CookCommand) Run(args []string) error {
	path := cmd.pathFlag()
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.cook(r)
	return nil
}

func (cmd CookCommand) cook(r *recipe.Recipe) bool {
	os.MkdirAll(r.ObjectDir, 0755)
	file, _ := os.Create(CompileLogFile(r.Path))

	defer file.Close()
	log := log.New(file, "", log.Ltime)

	return cmd.compile(
		r, log, compiler.NewCompiler(log, gxx.NewGXXCompiler(r.Includes, r.LibraryPaths, r.Libraries)),
	)
}

func (cmd CookCommand) compile(r *recipe.Recipe, log *log.Logger, c compiler.Compiler) bool {
	log.Println("[Compilation Start]")

	if ok := cmd.compileObjects(r, log, c); !ok {
		return cmd.fail(log, "Bad Ingredients!")
	}
	if ok := cmd.compileExecutable(r, log, c); !ok {
		return cmd.fail(log, "Burnt!")
	}
	return cmd.pass(log, "Bon Appétit!")
}

func (cmd CookCommand) compileObjects(r *recipe.Recipe, log *log.Logger, c compiler.Compiler) bool {
	style.Header.Println("Prepping...")

	indices := c.ComputeChangedSources(r)
	style.InfoV2.Printf("%s+ [%d] changed %s\n", INDENT, len(indices), style.SelectPlural("source", "sources", len(indices)))

	for i, index := range indices {
		src := r.SourceFiles[index]
		obj := r.ObjectFiles[index]

		cmd.printCompileFile(r, float32(i)/float32(len(indices)), src)
		res := c.CompileObject(src, obj)

		if !res {
			return false
		}
		style.Create.Println(r.TrimObjectDir(obj))
	}
	return true
}

func (cmd CookCommand) compileExecutable(r *recipe.Recipe, log *log.Logger, c compiler.Compiler) bool {
	style.Header.Println("Cooking...")

	count := len(r.ObjectFiles)
	if count > 0 {
		style.InfoV2.Printf("%s+ [%d] %s\n", INDENT, count, style.SelectPlural("object", "objects", count))
	}

	count = len(r.Libraries)
	if count > 0 {
		style.InfoV2.Printf("%s+ [%d] %s\n", INDENT, count, style.SelectPlural("library", "libraries", count))
	}

	cmd.printCompileFile(r, 1.0, r.MainSource)
	res := c.CompileExecutable(r.MainSource, r.Executable, r.ObjectFiles...)

	if res {
		style.BoldCreate.Println(r.TrimPath(r.Executable))
	}
	return res
}

func (cmd CookCommand) printCompileFile(r *recipe.Recipe, percent float32, src string) {
	style.BoldInfo.Printf("%s[%3d%%] ", INDENT, int(percent*100))
	style.File.Print(r.TrimPath(src))
	fmt.Print(" -> ")
}

func (CookCommand) fail(log *log.Logger, message string) bool {
	log.Println("[Compilation Failed]")
	style.BoldError.Println(message)
	return false
}

func (CookCommand) pass(log *log.Logger, message string) bool {
	log.Println("[Compilation Success]")
	style.BoldSuccess.Println(message)
	return true
}
