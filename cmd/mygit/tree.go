package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
)

type TreeObj struct {
	mode    string
	objType string
	name    string
	hash    string
}

func lsTree(hash string, nameOnly bool) {
	modes := map[string]string{
		"100644": "blob",
		"040000": "tree",
		"100755": "blob",
		"120000": "blob",
	}

	path := filepath.Join(ObjectsDir, hash[:2], hash[2:])
	file, err := readAndDecompressFile(path)

	if err != nil {
		handleError("Error reading file: %s\n", err)
	}

	tree := readContentsFromDecompressedTree(file)

	for len(tree) > 0 {
		var t TreeObj

		modeNameSeparator := bytes.IndexByte(tree, ' ')
		nullByteSeparator := bytes.IndexByte(tree, '\x00')

		t.mode = fmt.Sprintf("%06s", tree[:modeNameSeparator])
		t.objType = modes[t.mode]
		t.name = string(tree[modeNameSeparator+1 : nullByteSeparator])
		t.hash = hex.EncodeToString(tree[nullByteSeparator+1 : nullByteSeparator+21])

		tree = tree[nullByteSeparator+21:]

		if nameOnly {
			fmt.Println(t.name)
		} else {
			fmt.Println(t.mode, t.objType, t.hash, "  ", t.name)
		}
	}
}

func readContentsFromDecompressedTree(zTree io.ReadCloser) []byte {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zTree); err != nil {
		handleError("Error reading decompressed tree: %s\n", err)
	}

	return bytes.SplitN(buf.Bytes(), []byte("\x00"), 2)[1]
}
