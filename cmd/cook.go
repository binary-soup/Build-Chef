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

func (cmd CookCmd) Run(r *recipe.Recipe) {
	os.MkdirAll(recipe.OBJECT_DIR, 0755)
	file, _ := os.Create(COMPILE_LOG_FILE)
	defer file.Close()

	log := log.New(file, "", log.Ltime)
	com := compiler.Compiler{
		Log:         log,
		IncludeDirs: r.IncludeDirs,
	}

	log.Println("[Compilation Start]")

	style.Header.Println("Prepping...")
	if ok := com.CompileObjects(r.SourceFiles, r.ObjectFiles); !ok {
		log.Println("[Compilation Failed]")
		style.BoldError.Println("Bad Ingredients!")
		return
	}

	style.Header.Println("Cooking...")
	if ok := com.CompileExecutable("main.cxx", r.Name, r.ObjectFiles); !ok {
		log.Println("[Compilation Failed]")
		style.BoldError.Println("Burnt!")
		return
	}

	log.Println("[Compilation Success]")
	style.BoldSuccess.Println("Bon Appétit!")
}
