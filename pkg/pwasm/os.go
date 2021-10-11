package pwasm

import (
	"fmt"
	"os"
	"path/filepath"
)

func removeFile(filename string) error {
	fmt.Println(" - ", filename)
	return os.RemoveAll(filename)
}

func writeFile(filename string, data []byte) error {
	fmt.Println(" + ", filename)
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(filename, data, os.ModePerm)
}
