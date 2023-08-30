package cmd

import (
	"log"
	"os"

	"github.com/binary-soup/bchef/cmd/compiler"
	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

// NOTE: add -g for debug
// TODO: handle command injection

const (
	COMPILE_LOG_FILE = ".bchef/compile_log.txt"
)

type CookCmd struct{}

func (CookCmd) fail(log *log.Logger, message string) bool {
	log.Println("[Compilation Failed]")
	style.BoldError.Println(message)
	return false
}

func (CookCmd) pass(log *log.Logger, message string) bool {
	log.Println("[Compilation Success]")
	style.BoldSuccess.Println(message)
	return true
}

func (c CookCmd) Run(r *recipe.Recipe) bool {
	os.MkdirAll(recipe.OBJECT_DIR, 0755)
	file, _ := os.Create(COMPILE_LOG_FILE)
	defer file.Close()

	log := log.New(file, "", log.Ltime)
	com := compiler.Compiler{Indent: INDENT, Log: log, Recipe: r}

	log.Println("[Compilation Start]")

	style.Header.Println("Prepping...")
	if ok := com.CompileObjects(); !ok {
		return c.fail(log, "Bad Ingredients!")
	}

	style.Header.Println("Cooking...")
	if ok := com.CompileExecutable(); !ok {
		return c.fail(log, "Burnt!")
	}

	return c.pass(log, "Bon Appétit!")
}
