package main

import (
	"bytes"
	"compress/zlib"
	"io"
)

func compressBytes(byteSlice []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(byteSlice); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decompressBytes(byteSlice []byte) (io.ReadCloser, error) {
	return zlib.NewReader(bytes.NewReader(byteSlice))
}

func readAndDecompressBlob(hash string) io.ReadCloser {
	blob, err := readBlob(hash)
	if err != nil {
		handleError("Error reading file: %s\n", err)
	}
	zBlob, err := decompressBytes(blob)
	if err != nil {
		handleError("Error decompressing file: %s\n", err)
	}
	return zBlob
}

func readAndDecompressFile(file string) (io.ReadCloser, error) {
	compressedFile, err := readFile(file)
	if err != nil {
		return nil, err
	}

	decompressedFile, err := decompressBytes(compressedFile)
	if err != nil {
		return nil, err
	}

	return decompressedFile, nil
}

func compressAndWriteFile(file string, data []byte) error {
	compressedData, err := compressBytes(data)
	if err != nil {
		return err
	}

	return writeFile(file, compressedData)
}
