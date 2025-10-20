package create

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CopyAllApi(sourceDir, targetDir string) error {
	if _, err := os.Stat(sourceDir); err != nil {
		return err
	}

	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".api") {
			sourcePath := filepath.Join(sourceDir, file.Name())
			targetPath := filepath.Join(targetDir, file.Name())

			if err := copyFile(sourcePath, targetPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(source, destination string) error {
	srcFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
