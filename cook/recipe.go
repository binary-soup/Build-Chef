package cook

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/binary-soup/bchef/style"
)

type recipe struct {
	Name string
}

func newRecipe(path string) (*recipe, error) {
	file, err := os.Open(filepath.Join(path, "recipe.txt"))
	if os.IsNotExist(err) {
		return nil, errors.New("recipe not found")
	}
	if err != nil {
		return nil, errors.Join(errors.New("error opening file"), err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	return &recipe{
		Name: scanner.Text(),
	}, nil
}

func (r *recipe) exec() bool {
	//NOTE: add -g for debug
	//TODO: handle command injection

	cmd := exec.Command("g++", "-Wall", "-std=c++17", "-o", r.Name, "main.cxx")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	_, ok := cmd.Run().(*exec.ExitError)
	return ok
}

func (r *recipe) Cook() bool {
	fmt.Print(style.New(style.Magenta).Format("main.cxx"), " -> ")
	isExitErr := r.exec()

	if isExitErr {
		return false
	}

	style.New(style.Green, style.Bold).Println(r.Name)
	return true
}
