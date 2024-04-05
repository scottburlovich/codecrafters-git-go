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
		handleInitCommand()
	case "cat-file":
		handleCatFileCommand()
	case "hash-object":
		handleHashObjectCommand()
	case "ls-tree":
		handleLsTreeCommand()
	default:
		handleError("Unknown command %s\n", command)
	}
}

func handleInitCommand() {
	initRepo()
}

func handleCatFileCommand() {
	if len(os.Args) > 3 {
		catFile(os.Args[3])
	}
}

func handleHashObjectCommand() {
	file := ""
	write := false
	if len(os.Args) > 2 {
		write = os.Args[2] == "-w"
		if write && len(os.Args) > 3 {
			file = os.Args[3]
		} else {
			file = os.Args[2]
		}
	}
	hashObject(file, write)
}

func handleLsTreeCommand() {
	hash := ""
	nameOnly := false
	if len(os.Args) > 2 {
		nameOnly = os.Args[2] == "--name-only"
		if nameOnly && len(os.Args) > 3 {
			hash = os.Args[3]
		} else {
			hash = os.Args[2]
		}
	}
	lsTree(hash, nameOnly)
}

func handleError(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}
