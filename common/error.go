package common

import (
	"fmt"

	"github.com/binary-soup/bchef/style"
)

func PrintError(err error) {
	fmt.Println(style.BoldError.String("Error:"), err)
}
