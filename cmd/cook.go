package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/cmd/compiler"
	"github.com/binary-soup/bchef/cmd/compiler/gxx"
	"github.com/binary-soup/bchef/common"
	"github.com/binary-soup/bchef/config"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

// TODO: fix possible command injection?

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
	release := cmd.boolFlag("release", false, "build in release mode")
	outputPath := cmd.stringFlag("o", ".", "path to output directory")
	cmd.parseFlags(args)

	t, err := cmd.loadRecipeTree(*path)
	if err != nil {
		return err
	}

	debug := !*release
	if debug {
		style.BoldError.Println("[DEBUG MODE]")
	}

	res := t.Traverse(cookVisitor{
		outputPath:   *outputPath,
		compilerName: cfg.Compiler,
		debug:        debug,
	})

	if res {
		style.BoldSuccess.Println("Bon Appétit!")
	} else {
		style.BoldError.Println("Burnt!")
	}

	return nil
}

type cookVisitor struct {
	outputPath   string
	compilerName string
	debug        bool
}

func (v cookVisitor) Visit(r *recipe.Recipe, index int) bool {
	os.MkdirAll(r.GetObjectPath(v.debug), 0755)
	file, _ := os.Create(r.JoinPath(COMPILE_LOG_FILE))

	defer file.Close()
	log := log.New(file, "", log.Ltime)

	impl, err := v.chooseCompiler(r, v.compilerName)
	if err != nil {
		common.PrintError(err)
		return false
	}

	style.BoldFileV2.Printf("[%d] %s:\n", index+1, r.FullPath())

	com := compiler.NewCompiler(log, impl, v.createCompilerOptions(r, v.debug))
	return v.compile(r, log, com)
}

func (cookVisitor) chooseCompiler(r *recipe.Recipe, compilerName string) (compiler.CompilerImpl, error) {
	// TODO: support other compilers
	switch compilerName {
	case "", "g++":
		return gxx.NewGXXCompiler(r.Includes, r.LinkedStaticLibs, r.LibraryPaths, r.LinkedSharedLibs), nil
	default:
		return nil, fmt.Errorf("unsupported compiler \"%s\"", compilerName)
	}
}

func (cookVisitor) createCompilerOptions(r *recipe.Recipe, debug bool) compiler.Options {
	opts := compiler.Options{
		Debug:  false,
		PIC:    r.TargetType == recipe.TARGET_SHARED_LIBRARY,
		Macros: []string{},
	}

	if debug {
		opts.Debug = true
		opts.Macros = append(opts.Macros, "NDEBUG")
	}

	return opts
}

func (v cookVisitor) compile(r *recipe.Recipe, log *log.Logger, c compiler.Compiler) bool {
	log.Println("[Compilation Start]")
	indices, targetChanged := c.ComputeChangedSources(r, filepath.Join(v.outputPath, r.GetTarget(c.Opts.Debug)))

	if len(indices) == 0 && !targetChanged {
		style.Success.Println(INDENT + "[UP TO DATE]")
		return v.pass(log)
	}

	if ok := v.compileObjects(r, indices, log, c); !ok {
		return v.fail(log)
	}
	if ok := v.compileTarget(r, log, c); !ok {
		return v.fail(log)
	}
	return v.pass(log)
}

func (v cookVisitor) compileObjects(r *recipe.Recipe, indices []int, log *log.Logger, c compiler.Compiler) bool {
	fmt.Print(INDENT)
	style.Header.Println("Cooking...")

	style.InfoV2.Printf("%s+ [%d] changed %s\n", INDENT+INDENT, len(indices), common.SelectPlural("source", "sources", len(indices)))

	for i, index := range indices {
		src := r.SourceFiles[index]
		v.printCompileFile(r, float32(i)/float32(len(indices)), src)

		obj := r.ObjectFiles[index]
		res := c.CompileObject(r.JoinPath(src), r.JoinObjectPath(obj, c.Opts.Debug))

		if !res {
			return false
		}
		style.Create.Println(obj)
	}
	return true
}

func (v cookVisitor) compileTarget(r *recipe.Recipe, log *log.Logger, c compiler.Compiler) bool {
	fmt.Print(INDENT)
	style.Header.Println("Combining...")

	count := len(r.ObjectFiles)
	if count > 0 {
		style.InfoV2.Printf("%s+ [%d] %s\n", INDENT+INDENT, count, common.SelectPlural("object", "objects", count))
	}

	count = len(r.LinkedSharedLibs)
	if count > 0 {
		style.InfoV2.Printf("%s+ [%d] shared %s\n", INDENT+INDENT, count, common.SelectPlural("library", "libraries", count))
	}

	count = len(r.LinkedStaticLibs)
	if count > 0 {
		style.InfoV2.Printf("%s+ [%d] static %s\n", INDENT+INDENT, count, common.SelectPlural("library", "libraries", count))
	}

	objects := make([]string, len(r.ObjectFiles))
	for i, obj := range r.ObjectFiles {
		objects[i] = r.JoinObjectPath(obj, c.Opts.Debug)
	}

	v.printCompileFile(r, 1.0, "(ALL)")

	target := r.GetTarget(c.Opts.Debug)
	res := v.createTarget(r.TargetType, c, filepath.Join(v.outputPath, target), objects)

	if res {
		style.BoldCreate.Println(target)
	}
	return res
}

func (cookVisitor) createTarget(targetType int, c compiler.Compiler, target string, objs []string) bool {
	switch targetType {
	case recipe.TARGET_EXECUTABLE:
		return c.CreateExecutable(target, objs...)
	case recipe.TARGET_STATIC_LIBRARY:
		return c.CreateStaticLibrary(target, objs...)
	case recipe.TARGET_SHARED_LIBRARY:
		return c.CreateSharedLibrary(target, objs...)
	default:
		return false
	}
}

func (v cookVisitor) printCompileFile(r *recipe.Recipe, percent float32, src string) {
	style.BoldInfo.Printf("%s[%3d%%] ", INDENT+INDENT, int(percent*100))
	style.File.Print(src)
	fmt.Print(" -> ")
}

func (cookVisitor) fail(log *log.Logger) bool {
	log.Println("[Compilation Failed]")
	return false
}

func (cookVisitor) pass(log *log.Logger) bool {
	log.Println("[Compilation Success]")
	return true
}
