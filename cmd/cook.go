package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/binary-soup/bchef/recipe"
	"github.com/binary-soup/bchef/style"
)

// NOTE: add -g for debug
// TODO: handle command injection

const (
	COMPILE_LOG_FILE = ".bchef/compile_log.txt"
)

var compileFlags = []string{"-Wall", "-std=c++17"}

type CookCmd struct {
}

func (CookCmd) compile(createStyle style.Style, r *recipe.Recipe, log *log.Logger, flags []string, out string, src ...string) bool {
	fmt.Print("  ", style.FileV1.Format(r.TrimSourceDir(src[0])), " -> ")

	args := append(compileFlags, r.IncludeDirs...)
	args = append(args, flags...)
	args = append(args, src...)
	args = append(args, "-o", out)

	cmd := exec.Command("g++", args...)
	cmd.Stderr = os.Stdout

	_, err := cmd.Run().(*exec.ExitError)
	if !err {
		createStyle.Println(r.TrimObjectDir(out))
	}

	log.Printf("[Compile %s] %s\n", src[0], cmd.String())
	return !err
}

func (cmd CookCmd) compileSourceFiles(r *recipe.Recipe, log *log.Logger) bool {
	for i, src := range r.SourceFiles {
		if ok := cmd.compile(style.Create, r, log, []string{"-c"}, r.ObjectFiles[i], src); !ok {
			return false
		}
	}
	return true
}

func (cmd CookCmd) compileExecutable(r *recipe.Recipe, log *log.Logger) bool {
	style.FileV2.Printf("  + [%d] object files\n", len(r.ObjectFiles))

	src := append([]string{"main.cxx"}, r.ObjectFiles...)
	return cmd.compile(style.BoldCreate, r, log, []string{}, r.Name, src...)
}

func (cmd CookCmd) Run(r *recipe.Recipe) {
	os.MkdirAll(recipe.OBJECT_DIR, 0755)
	file, _ := os.Create(COMPILE_LOG_FILE)
	defer file.Close()

	log := log.New(file, "", log.Ltime)
	log.Println("[Compilation Start]")

	style.Header.Println("Prepping...")
	if ok := cmd.compileSourceFiles(r, log); !ok {
		log.Println("[Compilation Failed]")
		style.BoldError.Println("Bad Ingredients!")
		return
	}

	style.Header.Println("Cooking...")
	if ok := cmd.compileExecutable(r, log); !ok {
		log.Println("[Compilation Failed]")
		style.BoldError.Println("Burnt!")
		return
	}

	log.Println("[Compilation Success]")
	style.BoldSuccess.Println("Bon Appétit!")
}
