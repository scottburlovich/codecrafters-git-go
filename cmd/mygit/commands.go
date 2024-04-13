package main

import (
	"fmt"
	"github.com/codecrafters-io/git-starter-go/cmd/mygit/handlers"
	"os"
	"strconv"
	"strings"
)

type commandArgs struct {
	Args         map[string]bool
	ExpectedArgs []string
	OptionalArgs []string
	HandlerFunc  func(map[string]string)
}

var commandsMap = map[string]commandArgs{
	"init": {
		Args:         map[string]bool{},
		ExpectedArgs: []string{},
		OptionalArgs: []string{},
		HandlerFunc:  handlers.InitRepo,
	},
	"cat-file": {
		Args: map[string]bool{
			"-p": false,
		},
		ExpectedArgs: []string{"arg1", "-p"},
		OptionalArgs: []string{},
		HandlerFunc:  handlers.CatFile,
	},
	"hash-object": {
		Args: map[string]bool{
			"-w": false,
		},
		ExpectedArgs: []string{"arg1"},
		OptionalArgs: []string{"-w"},
		HandlerFunc:  handlers.HashObject,
	},
	"ls-tree": {
		Args: map[string]bool{
			"--name-only": false,
		},
		ExpectedArgs: []string{"arg1"},
		OptionalArgs: []string{"--name-only"},
		HandlerFunc:  handlers.LsTree,
	},
	"write-tree": {
		Args:         map[string]bool{},
		ExpectedArgs: []string{},
		OptionalArgs: []string{},
		HandlerFunc:  handlers.WriteTree,
	},
	"commit-tree": {
		Args: map[string]bool{
			"-m": true,
			"-p": true,
		},
		ExpectedArgs: []string{"arg1"},
		OptionalArgs: []string{"-p", "-m"},
		HandlerFunc:  handlers.CommitTree,
	},
	"clone": {
		Args:         map[string]bool{},
		ExpectedArgs: []string{"arg1", "arg2"},
		OptionalArgs: []string{"arg2"},
		HandlerFunc:  handlers.CloneRepository,
	},
}

func getArgs(cmd string, args []string) map[string]string {
	argMap := make(map[string]string)
	posArgCount := 1

	cmdArgsMap := commandsMap[cmd].Args
	expectedArgs := commandsMap[cmd].ExpectedArgs
	optionalArgs := commandsMap[cmd].OptionalArgs

	for i := 0; i < len(args); i++ {
		arg := args[i]

		if strings.HasPrefix(arg, "-") {
			expectsVal, isArg := cmdArgsMap[arg]

			if isArg {
				if expectsVal && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
					argMap[arg] = args[i+1]
					i++
				} else {
					argMap[arg] = ""
				}
			} else {
				fmt.Printf("Unknown argument: %v\n", arg)
				os.Exit(1)
			}
		} else {
			if posArgCount > len(expectedArgs) {
				fmt.Printf("Unexpected argument: %v\n", arg)
				os.Exit(1)
			}

			argMap["arg"+strconv.Itoa(posArgCount)] = arg
			posArgCount++
		}
	}

	for _, expectedArg := range expectedArgs {
		if _, ok := argMap[expectedArg]; !ok {
			// Only error out if the missing argument is not optional
			if !contains(optionalArgs, expectedArg) {
				fmt.Printf("Required argument missing: %s\n", expectedArg)
				os.Exit(1)
			}
		}
	}

	return argMap
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
