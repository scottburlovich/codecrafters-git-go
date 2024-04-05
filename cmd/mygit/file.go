package main

import (
	"os"
)

func readFile(file string) ([]byte, error) {
	return os.ReadFile(file)
}

func writeFile(file string, data []byte) error {
	return os.WriteFile(file, data, 0644)
}
