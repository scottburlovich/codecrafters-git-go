package main

import (
	"fmt"
)

func validateHash(hash string) error {
	if len(hash) != 40 {
		return fmt.Errorf("invalid hash: %s", hash)
	}
	return nil
}

func catFile(hash string) {
	if err := validateHash(hash); err != nil {
		handleError("Error: %s\n", err)
	}

	zBlob := readAndDecompressBlob(hash)
	defer zBlob.Close()

	fileContents := readFileContentsFromDecompressedBlob(zBlob)
	fmt.Printf("%s", fileContents)
}
