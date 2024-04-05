package main

import (
	"crypto/sha1"
	"io"
	"os"
)

func hashBytes(bytes []byte) []byte {
	h := sha1.New()
	h.Write(bytes)
	return h.Sum(nil)
}

func hashFile(filePath string) ([]byte, error) {
	h := sha1.New()
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}
