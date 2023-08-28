package cook

import "errors"

func Run(args []string) error {
	r, err := newRecipe(".")
	if err != nil {
		return errors.Join(errors.New("error loading recipe"), err)
	}

	r.Cook()
	return nil
}
