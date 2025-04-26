package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// copyFile копирует файл из src в dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return nil
}

// backup копирует все файлы из sourceDir в backupDir
func backup(sourceDir, backupDir string) error {
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(backupDir, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, os.ModePerm)
		}

		fmt.Printf("Backing up: %s -> %s\n", path, destPath)
		return copyFile(path, destPath)
	})

	return err
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go [source directory] [backup directory]")
		return
	}

	sourceDir := os.Args[1]
	backupDir := os.Args[2]

	// Добавим папку с текущей датой для удобства версионности
	timestamp := time.Now().Format("20060102_150405")
	backupDirWithTimestamp := filepath.Join(backupDir, timestamp)

	err := backup(sourceDir, backupDirWithTimestamp)
	if err != nil {
		fmt.Println("Backup failed:", err)
	} else {
		fmt.Println("Backup completed successfully!")
	}
}

