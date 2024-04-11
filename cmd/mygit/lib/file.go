package lib

import (
	"fmt"
	"os"
	"path/filepath"
)

func ReadFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func WriteFile(file string, data []byte) error {
	return os.WriteFile(file, data, 0644)
}

func WriteObject(obj []byte) ([]byte, error) {
	objHashSum, err := HashBytes(obj)
	if err != nil {
		HandleError("Error hashing object: %s\n", err)
	}

	zObj, err := compressBytes(obj)
	if err != nil {
		HandleError("Error compressing object: %s\n", err)
	}

	objectDir, err := CreateObjectDirectory(objHashSum)
	if err != nil {
		HandleError("Error creating object directory: %s\n", err)
	}

	hashString := fmt.Sprintf("%x", objHashSum)
	writePath := filepath.Join(objectDir, fmt.Sprintf("/%s", hashString[2:]))

	err = WriteFile(writePath, zObj)
	if err != nil {
		HandleError("Error writing object: %s\n", err)
	}

	return objHashSum, nil
}
