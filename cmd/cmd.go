package cmd

import (
	"github.com/binary-soup/bchef/recipe"
)

type Command interface {
	Run(*recipe.Recipe)
}
