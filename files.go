package main

import (
	"os"
	"path/filepath"
	"strings"
)

// LoadFileContent reads and returns the content of a file
func LoadFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// SaveFileContent writes content to a file
func SaveFileContent(filePath string, content string) error {
	return os.WriteFile(filePath, []byte(content), 0644)
}

// IsMarkdownFile checks if a filename has a markdown extension
func IsMarkdownFile(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".md") || strings.HasSuffix(lower, ".markdown")
}

// FindMarkdownFiles recursively finds all markdown files in a directory
func FindMarkdownFiles(dirPath string) ([]string, error) {
	var files []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && IsMarkdownFile(path) {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// CreateMarkdownFile creates a new markdown file with default content
func CreateMarkdownFile(filePath string) error {
	defaultContent := "# New Document\n\nStart writing..."
	return SaveFileContent(filePath, defaultContent)
}
