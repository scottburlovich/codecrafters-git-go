package main

import (
	"fmt"
	"os"
)

func readFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func createObjectDirectory(zBlobHashSum string) (string, error) {
	objectDir := fmt.Sprintf(ObjectsDir+"/%s", zBlobHashSum[:2])
	err := os.MkdirAll(objectDir, 0755)
	if err != nil {
		return "", err
	}
	return objectDir, nil
}
