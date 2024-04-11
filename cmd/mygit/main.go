package main

import (
	"github.com/codecrafters-io/git-starter-go/cmd/mygit/lib"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		lib.HandleError("usage: ./your-git.sh <command> [<args>...]\n", "")
	}
	command := os.Args[1]
	args := os.Args[2:]
	argsMap := getArgs(command, args)

	commandStruct, ok := commandsMap[command]

	if ok {
		commandStruct.HandlerFunc(argsMap)
	} else {
		lib.HandleError("Unknown command %s\n", command)
	}
}
