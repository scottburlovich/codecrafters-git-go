package main

import (
	"fmt"
	"os"
)

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
	case "hash-object":
		fileToHash := ""
		writeFlag := false
		if len(os.Args) > 2 {
			writeFlag = os.Args[2] == "-w"
			if writeFlag && len(os.Args) > 3 {
				fileToHash = os.Args[3]
			} else {
				fileToHash = os.Args[2]
			}
		}
		hashObject(fileToHash, writeFlag)

	default:
		handleError("Unknown command %s\n", command)
	}
}

func handleError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
