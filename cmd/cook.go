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
		command: newCommand("cook", "build the project"),
	}
}

type CookCommand struct {
	command
}

func (cmd CookCommand) Run(args []string) error {
	path := cmd.pathFlag()
	debug := cmd.flagSet.Bool("d", false, "build in debug mode")
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	cmd.cook(r, *debug)
	return nil
}

func (cmd CookCommand) cook(r *recipe.Recipe, debug bool) bool {
	os.MkdirAll(r.GetObjectPath(debug), 0755)
	file, _ := os.Create(CompileLogFile(r.Path))

	defer file.Close()
	log := log.New(file, "", log.Ltime)

	gxx := gxx.NewGXXCompiler(r.Includes, r.LibraryPaths, r.Libraries)
	com := compiler.NewCompiler(log, gxx, debug)

	return cmd.compile(r, log, com)
}

func (cmd CookCommand) compile(r *recipe.Recipe, log *log.Logger, c compiler.Compiler) bool {
	log.Println("[Compilation Start]")

	if c.Debug {
		style.BoldError.Println("[DEBUG MODE]")
	}

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
		cmd.printCompileFile(r, float32(i)/float32(len(indices)), src)

		obj := r.ObjectFiles[index]
		res := c.CompileObject(r.JoinPath(src), r.JoinObjectPath(obj, c.Debug))

		if !res {
			return false
		}
		style.Create.Println(obj)
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

	objects := make([]string, len(r.ObjectFiles))
	for i, obj := range r.ObjectFiles {
		objects[i] = r.JoinObjectPath(obj, c.Debug)
	}

	cmd.printCompileFile(r, 1.0, r.MainSource)
	res := c.CompileExecutable(r.JoinPath(r.MainSource), r.JoinPath(r.Executable), objects...)

	if res {
		style.BoldCreate.Println(r.Executable)
	}
	return res
}

func (cmd CookCommand) printCompileFile(r *recipe.Recipe, percent float32, src string) {
	style.BoldInfo.Printf("%s[%3d%%] ", INDENT, int(percent*100))
	style.File.Print(src)
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
