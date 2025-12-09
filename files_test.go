package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadFileContent_ValidFile_ReturnsContent(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.md")
	expectedContent := "# Test\n\nContent here"
	err := os.WriteFile(testFile, []byte(expectedContent), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Act
	content, err := LoadFileContent(testFile)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if content != expectedContent {
		t.Errorf("got %q, want %q", content, expectedContent)
	}
}

func TestLoadFileContent_NonExistentFile_ReturnsError(t *testing.T) {
	// Act
	_, err := LoadFileContent("/nonexistent/file.md")

	// Assert
	if err == nil {
		t.Error("expected error for nonexistent file, got nil")
	}
}

func TestLoadFileContent_EmptyFile_ReturnsEmptyString(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "empty.md")
	err := os.WriteFile(testFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Act
	content, err := LoadFileContent(testFile)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if content != "" {
		t.Errorf("got %q, want empty string", content)
	}
}

func TestSaveFileContent_ValidPath_CreatesFile(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "save-test.md")
	content := "# Saved Content\n\nThis was saved"

	// Act
	err := SaveFileContent(testFile, content)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify file was created with correct content
	savedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("failed to read saved file: %v", err)
	}
	if string(savedContent) != content {
		t.Errorf("got %q, want %q", string(savedContent), content)
	}
}

func TestSaveFileContent_InvalidPath_ReturnsError(t *testing.T) {
	// Act
	err := SaveFileContent("/invalid/nonexistent/path/file.md", "content")

	// Assert
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}

func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"lowercase md", "test.md", true},
		{"uppercase MD", "test.MD", true},
		{"mixed case Md", "test.Md", true},
		{"markdown extension", "test.markdown", true},
		{"uppercase MARKDOWN", "test.MARKDOWN", true},
		{"txt file", "test.txt", false},
		{"no extension", "test", false},
		{"empty string", "", false},
		{"md in middle", "test.md.txt", false},
		{"just md", ".md", true},
		{"path with md", "/path/to/file.md", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsMarkdownFile(tt.filename)
			if got != tt.want {
				t.Errorf("IsMarkdownFile(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestFindMarkdownFiles_EmptyDirectory_ReturnsEmptySlice(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()

	// Act
	files, err := FindMarkdownFiles(tmpDir)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(files) != 0 {
		t.Errorf("expected 0 files, got %d", len(files))
	}
}

func TestFindMarkdownFiles_WithMarkdownFiles_ReturnsAll(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()

	// Create test files
	testFiles := []string{"test1.md", "test2.MD", "test3.markdown", "test4.txt"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tmpDir, filename)
		err := os.WriteFile(filePath, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", filename, err)
		}
	}

	// Act
	files, err := FindMarkdownFiles(tmpDir)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	// Should find 3 markdown files (test1.md, test2.MD, test3.markdown)
	// Should NOT find test4.txt
	if len(files) != 3 {
		t.Errorf("expected 3 markdown files, got %d: %v", len(files), files)
	}
}

func TestFindMarkdownFiles_NestedDirectories_ReturnsAllRecursively(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	// Create files in root and subdirectory
	rootFile := filepath.Join(tmpDir, "root.md")
	subFile := filepath.Join(subDir, "sub.md")

	err = os.WriteFile(rootFile, []byte("root content"), 0644)
	if err != nil {
		t.Fatalf("failed to create root file: %v", err)
	}
	err = os.WriteFile(subFile, []byte("sub content"), 0644)
	if err != nil {
		t.Fatalf("failed to create sub file: %v", err)
	}

	// Act
	files, err := FindMarkdownFiles(tmpDir)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if len(files) != 2 {
		t.Errorf("expected 2 markdown files, got %d", len(files))
	}
}

func TestFindMarkdownFiles_NonExistentDirectory_ReturnsError(t *testing.T) {
	// Act
	_, err := FindMarkdownFiles("/nonexistent/directory/path")

	// Assert
	if err == nil {
		t.Error("expected error for nonexistent directory, got nil")
	}
}

func TestCreateMarkdownFile_ValidPath_CreatesFileWithDefaultContent(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "new.md")

	// Act
	err := CreateMarkdownFile(testFile)

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify file exists and has default content
	content, err := os.ReadFile(testFile)
	if err != nil {
		t.Errorf("failed to read created file: %v", err)
	}

	expectedContent := "# New Document\n\nStart writing..."
	if string(content) != expectedContent {
		t.Errorf("got %q, want %q", string(content), expectedContent)
	}
}

func TestCreateMarkdownFile_InvalidPath_ReturnsError(t *testing.T) {
	// Act
	err := CreateMarkdownFile("/invalid/nonexistent/path/file.md")

	// Assert
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
