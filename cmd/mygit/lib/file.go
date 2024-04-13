package lib

import (
	"bytes"
	"fmt"
	"io"
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
	zObj, err := compressBytes(obj)
	if err != nil {
		HandleError("Error compressing object: %s\n", err)
	}

	objHashSum := HashBytes(obj)

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

func WriteObjectWithType(obj []byte, objType string) ([]byte, error) {
	buf := bytes.Buffer{}
	fmt.Fprintf(&buf, "%s %d", objType, len(obj))
	buf.WriteByte(0)
	buf.Write(obj)

	return WriteObject(buf.Bytes())
}

func CreateObjectDirectory(hashSum []byte) (string, error) {
	hashString := fmt.Sprintf("%x", hashSum)
	objectDir := fmt.Sprintf(ObjectsDir+"/%s", hashString[:2])
	err := os.MkdirAll(objectDir, 0755)
	if err != nil {
		return "", err
	}

	return objectDir, nil
}

func ObjectFileExists(hashString string) bool {
	objectPath := filepath.Join(ObjectsDir, hashString[:2], hashString[2:])
	_, err := os.Stat(objectPath)
	return !os.IsNotExist(err)
}

func ReadObjectFile(hashString string) ([]byte, string, int, error) {
	objectPath := filepath.Join(ObjectsDir, hashString[:2], hashString[2:])
	zObj, err := ReadFile(objectPath)
	if err != nil {
		return nil, "", 0, err
	}
	obj, err := decompressBytes(zObj)
	if err != nil {
		return nil, "", 0, err
	}

	objData, err := io.ReadAll(obj)
	if err != nil {
		return nil, "", 0, err
	}

	byteIndex := bytes.IndexByte(objData, 0)
	var objType string
	var objSize int

	fmt.Sscanf(string(objData[:byteIndex]), "%s %d", &objType, &objSize)
	if byteIndex+objSize+1 != len(objData) {
		return nil, "", 0, fmt.Errorf("invalid object size")
	}

	return objData[byteIndex+1:], objType, objSize, nil
}

func InitRepository(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return fmt.Errorf("error changing directory: %s\n", err)
	}

	for _, dir := range []string{GitDir, ObjectsDir, RefsDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("error creating directory: %s\n", err)
		}
	}

	headFileContents := []byte("ref: refs/heads/main\n")
	if err := os.WriteFile(HeadFilePath, headFileContents, 0644); err != nil {
		return fmt.Errorf("error writing file: %s\n", err)
	}

	return nil
}

func SplitDirFile(hex string) (string, string) {
	return hex[:2], hex[2:]
}
