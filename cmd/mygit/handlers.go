package main

import (
	"fmt"
	"github.com/codecrafters-io/git-starter-go/cmd/mygit/config"
	"github.com/codecrafters-io/git-starter-go/cmd/mygit/lib"
	"os"
	"path/filepath"
)

func initRepo() {
	for _, dir := range []string{config.GitDir, config.ObjectsDir, config.RefsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			lib.HandleError("Error creating directory: %s\n", err)
		}
	}
	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(config.HeadFilePath, headFileContents, 0644); err != nil {
		lib.HandleError("Error writing file: %s\n", err)
	}
	fmt.Println("Initialized git directory")
}

func catFile(hash string) {
	if err := lib.ValidateHash(hash); err != nil {
		lib.HandleError("Error: %s\n", err)
	}

	zBlob := lib.ReadAndDecompressBlob(hash)
	defer zBlob.Close()

	fileContents := lib.ReadFileContentsFromDecompressedBlob(zBlob)
	fmt.Printf("%s", fileContents)
}

func hashObject(file string, write bool) {
	fileContents, err := lib.ReadFile(file)
	if err != nil {
		lib.HandleError("Error reading file: %s\n", err)
	}

	blob := lib.CreateBlob(fileContents)
	blobHashSum, err := lib.HashBytes(blob)
	if err != nil {
		lib.HandleError("Error hashing blob: %s\n", err)
	}

	if write {
		_, err = lib.WriteObject(blob)
		if err != nil {
			lib.HandleError("Error writing blob: %s\n", err)
		}
	}

	fmt.Printf("%x\n", blobHashSum)
}

func lsTree(hash string, nameOnly bool) {
	path := filepath.Join(config.ObjectsDir, hash[:2], hash[2:])
	file, err := lib.ReadAndDecompressFile(path)
	if err != nil {
		lib.HandleError("Error reading file: %s\n", err)
	}
	tree := lib.ReadContentsFromDecompressedTree(file)
	lib.ReadTree(tree, nameOnly)
}

func writeTree(path string) {
	tree, err := lib.TraverseTree(path)
	if err != nil {
		lib.HandleError("Error traversing tree: %s\n", err)
	}
	fmt.Printf("%x\n", tree)
}
