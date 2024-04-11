package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type commandArgs struct {
	Args         map[string]bool
	TotalArgs    int
	ExpectedArgs []string
	OptionalArgs []string
	HandlerFunc  func(map[string]string)
}

var commandsMap = map[string]commandArgs{
	"init": {
		Args:         map[string]bool{},
		TotalArgs:    0,
		ExpectedArgs: []string{},
		OptionalArgs: []string{},
		HandlerFunc:  initRepo,
	},
	"cat-file": {
		Args: map[string]bool{
			"-p": false,
		},
		TotalArgs:    2,
		ExpectedArgs: []string{"arg1", "-p"},
		OptionalArgs: []string{},
		HandlerFunc:  catFile,
	},
	"hash-object": {
		Args: map[string]bool{
			"-w": false,
		},
		TotalArgs:    1,
		ExpectedArgs: []string{"arg1"},
		OptionalArgs: []string{"-w"},
		HandlerFunc:  hashObject,
	},
	"ls-tree": {
		Args: map[string]bool{
			"--name-only": false,
		},
		TotalArgs:    1,
		ExpectedArgs: []string{"arg1"},
		OptionalArgs: []string{"--name-only"},
		HandlerFunc:  lsTree,
	},
	"write-tree": {
		Args:         map[string]bool{},
		TotalArgs:    0,
		ExpectedArgs: []string{},
		OptionalArgs: []string{},
		HandlerFunc:  writeTree,
	},
	"commit-tree": {
		Args: map[string]bool{
			"-m": true,
			"-p": true,
		},
		TotalArgs:    1,
		ExpectedArgs: []string{"arg1"},
		OptionalArgs: []string{"-p", "-m"},
		HandlerFunc:  commitTree,
	},
}

func getArgs(cmd string, args []string) map[string]string {
	argMap := make(map[string]string)
	cmdArgsMap := commandsMap[cmd].Args
	expectedArgs := commandsMap[cmd].ExpectedArgs
	optionalArgs := commandsMap[cmd].OptionalArgs
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			arg := args[i]
			if _, ok := cmdArgsMap[arg]; ok || contains(optionalArgs, arg) {
				if cmdArgsMap[arg] {
					if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
						argMap[arg] = args[i+1]
						i++
					} else {
						fmt.Printf("Value expected for argument %s\n", arg)
					}
				} else {
					argMap[arg] = ""
					expectedArgs = remove(expectedArgs, arg)
				}
			}
		} else if len(expectedArgs) > 0 {
			argMap["arg"+strconv.Itoa(len(expectedArgs))] = args[i]
			expectedArgs = expectedArgs[1:]
		}
	}

	if len(expectedArgs) > 0 {
		fmt.Printf("Required argument(s) missing: %v\n", expectedArgs)
		os.Exit(1)
	}

	return argMap
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
