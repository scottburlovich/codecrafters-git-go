package lib

import (
	"fmt"
	"os"
)

func HandleError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
