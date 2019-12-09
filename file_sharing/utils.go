package file_sharing

import (
	"fmt"
	"os"
	"path"
)

type File struct {
	Name   string
	Origin string
	Hash   string
}

func ExecutableRelativePath(relativePath string) (string, error) {
	executablePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%v", path.Dir(executablePath), relativePath), nil
}

func createDirectory(directory string) error {
	if !fileExists(directory) {
		return os.MkdirAll(directory, 0755)
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
