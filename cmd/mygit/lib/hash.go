package lib

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

func HashBytes(b []byte) ([]byte, error) {
	h := sha1.New()
	_, err := h.Write(b)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func HashFile(filePath string) ([]byte, error) {
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

func ValidateHash(hash string) error {
	if len(hash) != 40 {
		return fmt.Errorf("invalid hash: %s", hash)
	}
	return nil
}
