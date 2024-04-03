package main

import (
	"fmt"
	"os"
)

const (
	GitDir       = ".git"
	ObjectsDir   = ".git/objects"
	RefsDir      = ".git/refs"
	HeadFilePath = ".git/HEAD"
)

func initRepo() {
	for _, dir := range []string{GitDir, ObjectsDir, RefsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			handleError("Error creating directory: %s\n", err)
		}
	}
	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(HeadFilePath, headFileContents, 0644); err != nil {
		handleError("Error writing file: %s\n", err)
	}
	fmt.Println("Initialized git directory")
}
