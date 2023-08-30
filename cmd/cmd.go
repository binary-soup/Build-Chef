package cmd

import (
	"github.com/binary-soup/bchef/recipe"
)

const (
	INDENT = "  "
)

type Command interface {
	Run(*recipe.Recipe) bool
}
