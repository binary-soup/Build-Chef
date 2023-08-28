package cook

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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

func (r *recipe) Cook() {
	fmt.Println(r.Name)
}
