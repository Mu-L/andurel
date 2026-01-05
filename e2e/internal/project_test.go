package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewProject(t *testing.T) {
	andurelBinary := "/tmp/andurel"
	project := NewProject(t, andurelBinary)

	if project == nil {
		t.Fatal("NewProject returned nil")
	}

	if project.Name != "testapp" {
		t.Errorf("Expected project name 'testapp', got '%s'", project.Name)
	}

	if project.BinaryPath != andurelBinary {
		t.Errorf("Expected binary path '%s', got '%s'", andurelBinary, project.BinaryPath)
	}

	if project.Database != "" {
		t.Errorf("Expected empty database, got '%s'", project.Database)
	}

	if !filepath.IsAbs(project.Dir) {
		t.Errorf("Expected absolute path for project directory, got '%s'", project.Dir)
	}

	if project.T != t {
		t.Error("Expected project.T to be the passed testing.T")
	}
}

func TestNewProjectWithDatabase(t *testing.T) {
	andurelBinary := "/tmp/andurel"
	database := "postgresql://localhost/testdb"
	project := NewProjectWithDatabase(t, andurelBinary, database)

	if project == nil {
		t.Fatal("NewProjectWithDatabase returned nil")
	}

	if project.Name != "testapp" {
		t.Errorf("Expected project name 'testapp', got '%s'", project.Name)
	}

	if project.BinaryPath != andurelBinary {
		t.Errorf("Expected binary path '%s', got '%s'", andurelBinary, project.BinaryPath)
	}

	if project.Database != database {
		t.Errorf("Expected database '%s', got '%s'", database, project.Database)
	}

	if !filepath.IsAbs(project.Dir) {
		t.Errorf("Expected absolute path for project directory, got '%s'", project.Dir)
	}
}

func TestProjectFileExists(t *testing.T) {
	tempDir := t.TempDir()
	project := &Project{
		Dir: tempDir,
		T:   t,
	}

	// Test non-existent file
	if project.FileExists("nonexistent.txt") {
		t.Error("FileExists should return false for non-existent file")
	}

	// Test existing file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if !project.FileExists("test.txt") {
		t.Error("FileExists should return true for existing file")
	}
}

func TestProjectDirExists(t *testing.T) {
	tempDir := t.TempDir()
	project := &Project{
		Dir: tempDir,
		T:   t,
	}

	// Test non-existent directory
	if project.DirExists("nonexistent") {
		t.Error("DirExists should return false for non-existent directory")
	}

	// Test existing directory
	testDir := filepath.Join(tempDir, "testdir")
	if err := os.Mkdir(testDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	if !project.DirExists("testdir") {
		t.Error("DirExists should return true for existing directory")
	}

	// Test file (should return false)
	testFile := filepath.Join(tempDir, "testfile.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if project.DirExists("testfile.txt") {
		t.Error("DirExists should return false for files")
	}
}

func TestProjectGoVet(t *testing.T) {
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

	err := project.GoVet()
	// Should pass for simple valid Go code
	if err != nil {
		t.Errorf("GoVet should succeed for simple valid Go code, got: %v", err)
	}
}

func TestProjectGoBuild(t *testing.T) {
	tempDir := t.TempDir()

	// Create a minimal Go project
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

	// Test that GoBuild can be called (actual build may succeed or fail)
	// Just testing that the function runs without panicking
	target := "."  // Build to current directory
	err := project.GoBuild(target)
	// The build should succeed for simple Go code
	if err != nil {
		t.Logf("GoBuild failed: %v", err)
	}
}
