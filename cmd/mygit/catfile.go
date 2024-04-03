package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func validateHash(hash string) error {
	if len(hash) != 40 {
		return fmt.Errorf("invalid hash: %s", hash)
	}
	return nil
}

func readFromFile(hash string) ([]byte, error) {
	return os.ReadFile(fmt.Sprintf(".git/objects/%s/%s", hash[:2], hash[2:]))
}

func decompressBlob(blob []byte) (io.ReadCloser, error) {
	return zlib.NewReader(bytes.NewReader(blob))
}

func catFile(hash string) {
	if err := validateHash(hash); err != nil {
		handleError("Error: %s\n", err)
	}

	blob, err := readFromFile(hash)
	if err != nil {
		handleError("Error reading file: %s\n", err)
	}

	zBlob, err := decompressBlob(blob)
	if err != nil {
		handleError("Error decompressing file: %s\n", err)
	}
	defer zBlob.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zBlob); err != nil {
		handleError("Error reading decompressed blob: %s\n", err)
	}

	contents := bytes.SplitN(buf.Bytes(), []byte("\x00"), 2)[1]
	fmt.Printf("%s", contents)
}
