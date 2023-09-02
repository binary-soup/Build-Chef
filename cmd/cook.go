package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/binary-soup/bchef/cmd/compiler"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

// NOTE: add -g for debug
// TODO: handle command injection

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
	com := compiler.Compiler{Indent: INDENT, Log: log, Recipe: r}

	log.Println("[Compilation Start]")

	style.Header.Println("Prepping...")
	if ok := com.CompileObjects(); !ok {
		return cmd.fail(log, "Bad Ingredients!")
	}

	style.Header.Println("Cooking...")
	if ok := com.CompileExecutable(); !ok {
		return cmd.fail(log, "Burnt!")
	}

	return cmd.pass(log, "Bon Appétit!")
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
