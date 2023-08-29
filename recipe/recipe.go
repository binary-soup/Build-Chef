package recipe

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
)

type Recipe struct {
	Path string
	Name string
}

func Load(path string) (*Recipe, error) {
	path = filepath.Join(path, "recipe.txt")

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, errors.New("recipe file not found")
	}
	if err != nil {
		return nil, errors.Join(errors.New("error opening file"), err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	return &Recipe{
		Path: path,
		Name: scanner.Text(),
	}, nil
}
