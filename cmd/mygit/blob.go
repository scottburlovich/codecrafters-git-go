package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func readBlob(hash string) ([]byte, error) {
	return os.ReadFile(fmt.Sprintf(ObjectsDir+"/%s/%s", hash[:2], hash[2:]))
}

func writeBlob(dir, zBlobHashSum string, compressedBlob []byte) error {
	return os.WriteFile(fmt.Sprintf(dir+"/%s", zBlobHashSum[2:]), compressedBlob, 0644)
}

func compressBlob(blob []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(blob); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decompressBlob(blob []byte) (io.ReadCloser, error) {
	return zlib.NewReader(bytes.NewReader(blob))
}

func computeBlobHash(fileContents []byte) ([]byte, string) {
	zBlob := fmt.Sprintf("blob %d\x00%s", len(fileContents), string(fileContents))
	zBlobHash := sha1.New()
	zBlobHash.Write([]byte(zBlob))
	zBlobHashSum := zBlobHash.Sum(nil)
	return []byte(zBlob), fmt.Sprintf("%x", zBlobHashSum)
}

func readAndDecompressBlob(hash string) io.ReadCloser {
	blob, err := readBlob(hash)
	if err != nil {
		handleError("Error reading file: %s\n", err)
	}
	zBlob, err := decompressBlob(blob)
	if err != nil {
		handleError("Error decompressing file: %s\n", err)
	}
	return zBlob
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
		compressedBlob, err := compressBlob(zBlob)
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
