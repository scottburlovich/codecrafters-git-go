package main

import (
	"fmt"
	"os"
	// Uncomment this block to pass the first stage!
	// "os"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	if len(os.Args) < 2 {
		handleError("usage: mygit <command> [<args>...]\n", "")
	}

	command := os.Args[1]

	switch command {
	case "init":
		initRepo()
	case "cat-file":
		catFile(os.Args[3])
	default:
		handleError("Unknown command %s\n", command)
	}
}

func handleError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
