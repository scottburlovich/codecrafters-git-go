package lib

import (
	"bytes"
	"fmt"
	"io"
)

func CreateBlob(fileContents []byte) []byte {
	headerStr := fmt.Sprintf("blob %d", len(fileContents))
	header := append([]byte(headerStr), 0)
	return append(header, fileContents...)
}

func ReadBlob(hash string) ([]byte, error) {
	return ReadFile(fmt.Sprintf(ObjectsDir+"/%s/%s", hash[:2], hash[2:]))
}

func ReadFileContentsFromDecompressedBlob(zBlob io.ReadCloser) []byte {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, zBlob); err != nil {
		HandleError("Error reading decompressed blob: %s\n", err)
	}

	return bytes.SplitN(buf.Bytes(), []byte("\x00"), 2)[1]
}
