package main

import (
	"encoding/hex"
	"fmt"
	"github.com/codecrafters-io/git-starter-go/cmd/mygit/lib"
	"os"
	"path/filepath"
)

func initRepo(args map[string]string) {
	for _, dir := range []string{lib.GitDir, lib.ObjectsDir, lib.RefsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			lib.HandleError("Error creating directory: %s\n", err)
		}
	}
	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(lib.HeadFilePath, headFileContents, 0644); err != nil {
		lib.HandleError("Error writing file: %s\n", err)
	}
	fmt.Println("Initialized git directory")
}

func catFile(args map[string]string) {
	hash := args["arg1"]
	_, prettyPrint := args["-p"]

	if err := lib.ValidateHash(hash); err != nil {
		lib.HandleError("Error: %s\n", err)
	}

	zBlob := lib.ReadAndDecompressBlob(hash)
	defer zBlob.Close()

	fileContents := lib.ReadFileContentsFromDecompressedBlob(zBlob)

	if prettyPrint {
		fmt.Printf("%s", fileContents)
	}
}

func commitTree(args map[string]string) {
	tree := args["arg1"]
	parent, hasParent := args["-p"]
	message := args["-m"]

	treeHash, err := hex.DecodeString(tree)
	if err != nil {
		lib.HandleError("Error decoding tree hash: %s\n", err)
	}

	var parentHash []byte
	if hasParent {
		parentHash, err = hex.DecodeString(parent)
		if err != nil {
			lib.HandleError("Error decoding parent hash: %s\n", err)
		}
	}

	author := lib.DefaultAuthor
	authorEmail := lib.DefaultAuthorEmail

	commit := lib.CreateCommit(treeHash, parentHash, message, author, authorEmail)
	commitHash, err := lib.HashBytes(commit)
	if err != nil {
		lib.HandleError("Error hashing commit: %s\n", err)
	}

	_, err = lib.WriteObject(commit)
	if err != nil {
		lib.HandleError("Error writing commit: %s\n", err)
	}

	fmt.Printf("%x\n", commitHash)
}

func hashObject(args map[string]string) {
	file := args["arg1"]
	_, write := args["-w"]

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

func lsTree(args map[string]string) {
	hash := args["arg1"]
	_, nameOnly := args["--name-only"]
	path := filepath.Join(lib.ObjectsDir, hash[:2], hash[2:])
	file, err := lib.ReadAndDecompressFile(path)
	if err != nil {
		lib.HandleError("Error reading file: %s\n", err)
	}
	tree := lib.ReadContentsFromDecompressedTree(file)
	lib.ReadTree(tree, nameOnly)
}

func writeTree(args map[string]string) {
	tree, err := lib.TraverseTree(".")
	if err != nil {
		lib.HandleError("Error traversing tree: %s\n", err)
	}
	fmt.Printf("%x\n", tree)
}
