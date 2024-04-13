package handlers

import (
	"encoding/hex"
	"fmt"
	"github.com/codecrafters-io/git-starter-go/cmd/mygit/lib"
	"os"
	"path/filepath"
)

func InitRepo(args map[string]string) {
	err := lib.InitRepository(".")
	if err != nil {
		lib.HandleError("Error initializing git repository: %s\n", err)
	}
	fmt.Println("Initialized git repository")
}

func CatFile(args map[string]string) {
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

func CommitTree(args map[string]string) {
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
	commitHash := lib.HashBytes(commit)

	_, err = lib.WriteObject(commit)
	if err != nil {
		lib.HandleError("Error writing commit: %s\n", err)
	}

	fmt.Printf("%x\n", commitHash)
}

func HashObject(args map[string]string) {
	file := args["arg1"]
	_, write := args["-w"]

	fileContents, err := lib.ReadFile(file)
	if err != nil {
		lib.HandleError("Error reading file: %s\n", err)
	}

	blob := lib.CreateBlob(fileContents)
	blobHashSum := lib.HashBytes(blob)

	if write {
		_, err = lib.WriteObject(blob)
		if err != nil {
			lib.HandleError("Error writing blob: %s\n", err)
		}
	}

	fmt.Printf("%x\n", blobHashSum)
}

func LsTree(args map[string]string) {
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

func WriteTree(args map[string]string) {
	tree, err := lib.TraverseTree(".")
	if err != nil {
		lib.HandleError("Error traversing tree: %s\n", err)
	}
	fmt.Printf("%x\n", tree)
}

func CloneRepository(args map[string]string) {
	remoteURL := args["arg1"]
	localPath := args["arg2"]

	if localPath == "" {
		localPath = filepath.Base(remoteURL)

		if localPath[len(localPath)-4:] == ".git" {
			localPath = localPath[:len(localPath)-4]
		}
	}

	workingDir := os.Getenv("PWD")
	localPath = filepath.Join(workingDir, localPath)

	lib.CloneRepository(remoteURL, localPath)
}
