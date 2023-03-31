package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: ./copyfiles [source path prefix] [destination path prefix]")
		os.Exit(1)
	}

	sourcePathPrefix := os.Args[1]
	destPathPrefix := os.Args[2]

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fileName := scanner.Text()
		sourceFilePath := filepath.Join(sourcePathPrefix, fileName)
		destFilePath := filepath.Join(destPathPrefix, fileName)

		if err := copyFile(sourceFilePath, destFilePath); err != nil {
			fmt.Printf("Error copying file %s: %s\n", fileName, err)
		} else {
			fmt.Printf("Copied file %s to %s\n", sourceFilePath, destFilePath)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input: %s\n", err)
		os.Exit(1)
	}
}

func copyFile(sourceFilePath, destFilePath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destFilePath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	if err := destFile.Sync(); err != nil {
		return err
	}

	sourceFileInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return err
	}

	if err := os.Chmod(destFilePath, sourceFileInfo.Mode()); err != nil {
		return err
	}

	if err := os.Chtimes(destFilePath, sourceFileInfo.ModTime(), sourceFileInfo.ModTime()); err != nil {
		return err
	}

	return nil
}
