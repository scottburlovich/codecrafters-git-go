package lib

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

type TreeObj struct {
	mode    string
	objType string
	name    string
	hash    []byte
}

func TraverseTree(path string) ([]byte, error) {
	pathContent, err := os.ReadDir(path)
	if err != nil {
		HandleError("Error reading directory: %s\n", err)
	}
	treeContent, err := collectTreeContent(path, pathContent)
	if err != nil {
		HandleError("Error collecting tree content: %s\n", err)
	}
	sort.SliceStable(treeContent, func(i, j int) bool {
		return treeContent[i].name < treeContent[j].name
	})
	var treeEntries [][]byte
	for _, t := range treeContent {
		var modeStr string
		if t.mode == ModeBlob {
			modeStr = fmt.Sprintf("%06s", t.mode)
		} else {
			modeStr = fmt.Sprintf("%05s", t.mode)
		}
		modePath := []byte(fmt.Sprintf("%s %s", modeStr, t.name))
		treeEntry := append(modePath, append([]byte{0}, t.hash...)...)
		treeEntries = append(treeEntries, treeEntry)
	}

	combinedTreeEntries := bytes.Join(treeEntries, []byte{})
	headerStr := fmt.Sprintf("tree %d", len(combinedTreeEntries))
	header := append([]byte(headerStr), 0)
	tree := append(header, combinedTreeEntries...)

	return WriteObject(tree)
}

func collectTreeContent(path string, pathContent []os.DirEntry) ([]*TreeObj, error) {
	treeContent := make([]*TreeObj, 0, len(pathContent))
	for _, entry := range pathContent {
		if entry.IsDir() && filepath.Join(path, entry.Name()) == filepath.Join(path, ".git") {
			continue
		}
		treeObj, err := processEntry(path, entry)
		if err != nil {
			return nil, err
		}
		treeContent = append(treeContent, treeObj)
	}
	return treeContent, nil
}

func processEntry(path string, entry os.DirEntry) (*TreeObj, error) {
	mode, objType, hash, err := categorizeAndHandleEntry(path, entry)
	if err != nil {
		return nil, err
	}
	return &TreeObj{
		mode:    mode,
		objType: objType,
		name:    entry.Name(),
		hash:    hash,
	}, nil
}

func categorizeAndHandleEntry(path string, entry os.DirEntry) (string, string, []byte, error) {
	fi, _ := entry.Info()
	modePerm := fi.Mode().Perm()

	if entry.Type().IsDir() {
		newPath := filepath.Join(path, entry.Name())
		hash, err := TraverseTree(newPath)
		return ModeTree, Tree, hash, err
	}

	var mode, objType string
	if (modePerm & 0111) != 0 {
		mode = ModeBlobExec
		objType = Blob
	} else {
		mode = ModeBlob
		objType = Blob
	}
	hash, err := processBlob(path, entry.Name())

	return mode, objType, hash, err
}

func processBlob(path, name string) ([]byte, error) {
	fileContents, err := ReadFile(filepath.Join(path, name))
	if err != nil {
		return nil, err
	}

	blob := CreateBlob(fileContents)
	blobHash, err := HashBytes(blob)
	if err != nil {
		return nil, err
	}

	return blobHash, nil
}

func ReadTree(tree []byte, nameOnly bool) {
	for len(tree) > 0 {
		t, remainingTree := extractTreeData(tree)
		tree = remainingTree
		displayTreeData(t, nameOnly)
	}
}

func extractTreeData(tree []byte) (TreeObj, []byte) {
	modes := getModes()
	var t TreeObj
	modeNameSeparator := bytes.IndexByte(tree, ' ')
	nullByteSeparator := bytes.IndexByte(tree, '\x00')

	t.mode = fmt.Sprintf("%s", tree[:modeNameSeparator])
	t.objType = modes[t.mode]
	t.name = string(tree[modeNameSeparator+1 : nullByteSeparator])
	t.hash = tree[nullByteSeparator+1 : nullByteSeparator+21]
	return t, tree[nullByteSeparator+21:]
}

func getModes() map[string]string {
	return map[string]string{
		ModeBlob:     Blob,
		ModeTree:     Tree,
		ModeBlobExec: Blob,
		ModeSymLink:  Blob,
	}
}

func displayTreeData(t TreeObj, nameOnly bool) {
	if nameOnly {
		fmt.Println(t.name)
	} else {
		fmt.Println(t.mode, t.objType, hex.EncodeToString(t.hash), "  ", t.name)
	}
}

func ReadContentsFromDecompressedTree(zTree io.ReadCloser) []byte {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zTree); err != nil {
		HandleError("Error reading decompressed tree: %s\n", err)
	}

	return bytes.SplitN(buf.Bytes(), []byte("\x00"), 2)[1]
}
