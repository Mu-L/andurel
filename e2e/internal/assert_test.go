package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAssertFileExists(t *testing.T) {
	tempDir := t.TempDir()
	project := &Project{
		Dir: tempDir,
		T:   t,
	}

	t.Run("file exists", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "exists.txt")
		if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// This should not fail
		AssertFileExists(t, project, "exists.txt")
	})
}

func TestAssertDirExists(t *testing.T) {
	tempDir := t.TempDir()
	project := &Project{
		Dir: tempDir,
		T:   t,
	}

	t.Run("directory exists", func(t *testing.T) {
		testDir := filepath.Join(tempDir, "testdir")
		if err := os.Mkdir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// This should not fail
		AssertDirExists(t, project, "testdir")
	})
}

func TestAssertFilesExist(t *testing.T) {
	tempDir := t.TempDir()
	project := &Project{
		Dir: tempDir,
		T:   t,
	}

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "file3.txt"}
	for _, filename := range testFiles {
		testFile := filepath.Join(tempDir, filename)
		if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	t.Run("all files exist", func(t *testing.T) {
		// This should not fail
		AssertFilesExist(t, project, testFiles)
	})
}

func TestAssertGoVetPasses(t *testing.T) {
	tempDir := t.TempDir()

	// Create a minimal Go project with at least one .go file
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	mainPath := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(mainPath, []byte("package main\n\nfunc main() {}\n"), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	project := &Project{
		Dir: tempDir,
		T:   t,
	}

	// This should not fail for valid Go code
	AssertGoVetPasses(t, project)
}

func TestAssertCommandSucceeds(t *testing.T) {
	t.Run("command succeeds", func(t *testing.T) {
		var err error = nil
		cmdDesc := "test command"

		// This should not fail
		AssertCommandSucceeds(t, err, cmdDesc)
	})
}

func TestAssertOutputContains(t *testing.T) {
	t.Run("output contains expected string", func(t *testing.T) {
		output := "This is a test output string"
		expected := "test output"

		// This should not fail
		AssertOutputContains(t, output, expected)
	})

	t.Run("output contains empty expected", func(t *testing.T) {
		output := "This is a test output string"
		expected := ""

		// Empty string is contained in any string
		AssertOutputContains(t, output, expected)
	})

	t.Run("full match", func(t *testing.T) {
		output := "exact match"
		expected := "exact match"

		// This should not fail
		AssertOutputContains(t, output, expected)
	})
}
