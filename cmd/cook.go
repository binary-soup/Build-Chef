package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/binary-soup/bchef/cmd/compiler"
	"github.com/binary-soup/bchef/cmd/compiler/gxx"
	"github.com/binary-soup/bchef/config"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

// TODO: fix command injection

const COMPILE_LOG_FILE = ".bchef/compile_log.txt"

func NewCookCommand() CookCommand {
	return CookCommand{
		command: newCommand("cook", "build the project"),
	}
}

type CookCommand struct {
	command
}

func (cmd CookCommand) Run(cfg config.Config, args []string) error {
	path := cmd.pathFlag()
	release := cmd.flagSet.Bool("release", false, "build in release mode")
	cmd.parseFlags(args)

	r, err := cmd.loadRecipe(*path)
	if err != nil {
		return err
	}

	return cmd.cook(r, cfg.Compiler, !*release)
}

func (cmd CookCommand) cook(r *recipe.Recipe, compilerName string, debug bool) error {
	os.MkdirAll(r.GetObjectPath(debug), 0755)
	file, _ := os.Create(r.JoinPath(COMPILE_LOG_FILE))

	defer file.Close()
	log := log.New(file, "", log.Ltime)

	impl, err := cmd.chooseCompiler(r, compilerName)
	if err != nil {
		return err
	}

	com := compiler.NewCompiler(log, impl, cmd.createCompilerOptions(debug))
	cmd.compile(r, log, com)

	return nil
}

func (CookCommand) chooseCompiler(r *recipe.Recipe, compilerName string) (compiler.CompilerImpl, error) {
	switch compilerName {
	case "", "g++":
		return gxx.NewGXXCompiler(r.Includes, r.StaticLibs, r.LibraryPaths, r.SharedLibs), nil
	default:
		return nil, fmt.Errorf("unsupported compiler \"%s\"", compilerName)
	}
}

func (CookCommand) createCompilerOptions(debug bool) compiler.Options {
	if debug {
		return compiler.Options{
			Debug:  true,
			Macros: []string{"NDEBUG"},
		}
	} else {
		return compiler.Options{
			Debug:  false,
			Macros: []string{},
		}
	}
}

func (cmd CookCommand) compile(r *recipe.Recipe, log *log.Logger, c compiler.Compiler) bool {
	log.Println("[Compilation Start]")

	if c.Opts.Debug {
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
		res := c.CompileObject(r.JoinPath(src), r.JoinObjectPath(obj, c.Opts.Debug))

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

	count = len(r.SharedLibs)
	if count > 0 {
		style.InfoV2.Printf("%s+ [%d] shared %s\n", INDENT, count, style.SelectPlural("library", "libraries", count))
	}

	count = len(r.StaticLibs)
	if count > 0 {
		style.InfoV2.Printf("%s+ [%d] static %s\n", INDENT, count, style.SelectPlural("library", "libraries", count))
	}

	objects := make([]string, len(r.ObjectFiles))
	for i, obj := range r.ObjectFiles {
		objects[i] = r.JoinObjectPath(obj, c.Opts.Debug)
	}

	cmd.printCompileFile(r, 1.0, "(ALL)")

	exec := r.GetExecutable(c.Opts.Debug)
	res := c.CompileExecutable(r.JoinPath(exec), objects...)

	if res {
		style.BoldCreate.Println(exec)
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
