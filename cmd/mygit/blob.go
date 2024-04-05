package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

func readBlob(hash string) ([]byte, error) {
	return readFile(fmt.Sprintf(ObjectsDir+"/%s/%s", hash[:2], hash[2:]))
}

func writeBlob(dir, zBlobHashSum string, compressedBlob []byte) error {
	return writeFile(fmt.Sprintf(dir+"/%s", zBlobHashSum[2:]), compressedBlob)
}

func computeBlobHash(fileContents []byte) ([]byte, string) {
	zBlob := fmt.Sprintf("blob %d\x00%s", len(fileContents), string(fileContents))
	zBlobHashSum := hashBytes([]byte(zBlob))
	return []byte(zBlob), fmt.Sprintf("%x", zBlobHashSum)
}

func readFileContentsFromDecompressedBlob(zBlob io.ReadCloser) []byte {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zBlob); err != nil {
		handleError("Error reading decompressed blob: %s\n", err)
	}

	return bytes.SplitN(buf.Bytes(), []byte("\x00"), 2)[1]
}

func hashObject(file string, write bool) {
	fileContents, err := readFile(file)
	if err != nil {
		handleError("Error reading file: %s\n", err)
	}

	zBlob, zBlobHashSum := computeBlobHash(fileContents)

	if write {
		compressedBlob, err := compressBytes(zBlob)
		if err != nil {
			handleError("Error compressing file: %s\n", err)
		}

		objectDir, err := createObjectDirectory(zBlobHashSum)
		if err != nil {
			handleError("Error creating object directory: %s\n", err)
		}

		err = writeBlob(objectDir, zBlobHashSum, compressedBlob)
		if err != nil {
			handleError("Error writing file: %s\n", err)
		}
	}

	fmt.Printf("%s\n", zBlobHashSum)
}

func createObjectDirectory(zBlobHashSum string) (string, error) {
	objectDir := fmt.Sprintf(ObjectsDir+"/%s", zBlobHashSum[:2])
	err := os.MkdirAll(objectDir, 0755)
	if err != nil {
		return "", err
	}
	return objectDir, nil
}
