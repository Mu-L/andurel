package files

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbvlabs/andurel/pkg/cache"
)

func TestDefaultPermissions(t *testing.T) {
	perms := DefaultPermissions()

	if perms.FilePrivate != 0o600 {
		t.Errorf("expected FilePrivate to be 0600, got %v", perms.FilePrivate)
	}
	if perms.FilePublic != 0o644 {
		t.Errorf("expected FilePublic to be 0644, got %v", perms.FilePublic)
	}
	if perms.DirDefault != 0o755 {
		t.Errorf("expected DirDefault to be 0755, got %v", perms.DirDefault)
	}
	if perms.DirExecutable != 0o755 {
		t.Errorf("expected DirExecutable to be 0755, got %v", perms.DirExecutable)
	}
}

func TestNewUnifiedFileManager(t *testing.T) {
	fm := NewUnifiedFileManager()

	if fm == nil {
		t.Fatal("expected non-nil file manager")
	}
	if fm.cache == nil {
		t.Error("expected non-nil cache")
	}

	perms := fm.GetPermissions()
	if perms.FilePrivate != 0o600 {
		t.Errorf("expected FilePrivate to be 0600, got %v", perms.FilePrivate)
	}
}

func TestWriteFile(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewUnifiedFileManager()

	testPath := filepath.Join(tmpDir, "subdir", "test.txt")
	content := "test content"

	err := fm.WriteFile(testPath, content)
	if err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Verify file exists
	if !fm.FileExists(testPath) {
		t.Error("file should exist after WriteFile")
	}

	// Verify content
	readContent, err := fm.ReadFile(testPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if readContent != content {
		t.Errorf("expected content %q, got %q", content, readContent)
	}

	// Verify directory was created
	dirPath := filepath.Join(tmpDir, "subdir")
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Error("directory should have been created")
	}
}

func TestWriteFileWithPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewUnifiedFileManager()

	testPath := filepath.Join(tmpDir, "test.txt")
	content := "test content"
	perm := os.FileMode(0o755)

	err := fm.WriteFileWithPermissions(testPath, content, perm)
	if err != nil {
		t.Fatalf("WriteFileWithPermissions failed: %v", err)
	}

	// Verify file permissions
	stat, err := os.Stat(testPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if stat.Mode().Perm() != perm {
		t.Errorf("expected permissions %v, got %v", perm, stat.Mode().Perm())
	}
}

func TestReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewUnifiedFileManager()

	testPath := filepath.Join(tmpDir, "test.txt")
	expectedContent := "test content\nline 2"

	// Write file directly
	err := os.WriteFile(testPath, []byte(expectedContent), 0o644)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Read using file manager
	content, err := fm.ReadFile(testPath)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}
	if content != expectedContent {
		t.Errorf("expected content %q, got %q", expectedContent, content)
	}
}

func TestReadFileNonExistent(t *testing.T) {
	fm := NewUnifiedFileManager()

	_, err := fm.ReadFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestFileExists(t *testing.T) {
	cache.ClearFileSystemCache()
	defer cache.ClearFileSystemCache()
	
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "test.txt")

	// File should not exist initially
	fm := NewUnifiedFileManager()
	if fm.FileExists(testPath) {
		t.Error("file should not exist initially")
	}

	// Create file
	err := os.WriteFile(testPath, []byte("test"), 0o644)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Clear cache to re-check
	cache.ClearFileSystemCache()
	
	// File should exist now
	if !fm.FileExists(testPath) {
		t.Error("file should exist after creation")
	}
}

func TestEnsureDir(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewUnifiedFileManager()

	testPath := filepath.Join(tmpDir, "subdir1", "subdir2", "subdir3")

	err := fm.EnsureDir(testPath)
	if err != nil {
		t.Fatalf("EnsureDir failed: %v", err)
	}

	// Verify directory exists
	stat, err := os.Stat(testPath)
	if err != nil {
		t.Fatalf("directory should exist: %v", err)
	}
	if !stat.IsDir() {
		t.Error("path should be a directory")
	}
}

func TestEnsureDirWithPermissions(t *testing.T) {
	tmpDir := t.TempDir()
	fm := NewUnifiedFileManager()

	testPath := filepath.Join(tmpDir, "testdir")
	perm := os.FileMode(0o700)

	err := fm.EnsureDirWithPermissions(testPath, perm)
	if err != nil {
		t.Fatalf("EnsureDirWithPermissions failed: %v", err)
	}

	// Verify directory permissions
	stat, err := os.Stat(testPath)
	if err != nil {
		t.Fatalf("Stat failed: %v", err)
	}
	if stat.Mode().Perm() != perm {
		t.Errorf("expected permissions %v, got %v", perm, stat.Mode().Perm())
	}
}

func TestValidateFileNotExists(t *testing.T) {
	cache.ClearFileSystemCache()
	defer cache.ClearFileSystemCache()
	
	tmpDir := t.TempDir()
	fm := NewUnifiedFileManager()

	testPath := filepath.Join(tmpDir, "test.txt")

	// Should pass when file doesn't exist
	err := fm.ValidateFileNotExists(testPath)
	if err != nil {
		t.Errorf("validation should pass for nonexistent file: %v", err)
	}

	// Create file
	os.WriteFile(testPath, []byte("test"), 0o644)

	// Clear cache to re-check
	cache.ClearFileSystemCache()
	
	// Should fail when file exists
	err = fm.ValidateFileNotExists(testPath)
	if err == nil {
		t.Error("validation should fail for existing file")
	} else {
		// Check error type
		fileOpErr, ok := err.(*FileOperationError)
		if !ok {
			t.Fatalf("expected *FileOperationError, got %T", err)
		}
		if !os.IsExist(fileOpErr.Err) {
			t.Error("error should wrap os.ErrExist")
		}
		if fileOpErr.Operation != "validate_not_exists" {
			t.Errorf("expected operation 'validate_not_exists', got %q", fileOpErr.Operation)
		}
	}
}

func TestValidateFileExists(t *testing.T) {
	cache.ClearFileSystemCache()
	defer cache.ClearFileSystemCache()
	
	tmpDir := t.TempDir()
	fm := NewUnifiedFileManager()

	testPath := filepath.Join(tmpDir, "test.txt")

	// Should fail when file doesn't exist
	err := fm.ValidateFileExists(testPath)
	if err == nil {
		t.Error("validation should fail for nonexistent file")
	}

	// Check error type
	if err != nil {
		fileOpErr, ok := err.(*FileOperationError)
		if !ok {
			t.Fatalf("expected *FileOperationError, got %T", err)
		}
		if !os.IsNotExist(fileOpErr.Err) {
			t.Error("error should wrap os.ErrNotExist")
		}
	}

	// Create file
	os.WriteFile(testPath, []byte("test"), 0o644)
	
	// Clear cache to re-check
	cache.ClearFileSystemCache()

	// Should pass when file exists
	err = fm.ValidateFileExists(testPath)
	if err != nil {
		t.Errorf("validation should pass for existing file: %v", err)
	}
}

func TestSetPermissions(t *testing.T) {
	fm := NewUnifiedFileManager()

	customPerms := Permissions{
		FilePrivate:   0o600,
		FilePublic:    0o644,
		DirDefault:    0o700,
		DirExecutable: 0o755,
	}

	fm.SetPermissions(customPerms)

	gotPerms := fm.GetPermissions()
	if gotPerms != customPerms {
		t.Errorf("expected permissions %+v, got %+v", customPerms, gotPerms)
	}
}

func TestFindGoModRoot(t *testing.T) {
	// This test needs to run in the actual project directory
	// since it searches for go.mod
	fm := NewUnifiedFileManager()

	root, err := fm.FindGoModRoot()
	if err != nil {
		t.Fatalf("FindGoModRoot failed: %v", err)
	}

	// Verify go.mod exists at the returned path
	goModPath := filepath.Join(root, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Errorf("go.mod should exist at %s", goModPath)
	}
}

func TestFindGoModRootNotFound(t *testing.T) {
	cache.ClearFileSystemCache()
	defer cache.ClearFileSystemCache()
	
	tmpDir := t.TempDir()
	fm := NewUnifiedFileManager()

	// Change to temp directory that doesn't have go.mod
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)

	os.Chdir(tmpDir)

	_, err := fm.FindGoModRoot()
	if err == nil {
		t.Error("expected error when go.mod not found")
	} else {
		fileOpErr, ok := err.(*FileOperationError)
		if !ok {
			t.Fatalf("expected *FileOperationError, got %T", err)
		}
		if fileOpErr.Operation != "find_gomod_root" {
			t.Errorf("expected operation 'find_gomod_root', got %q", fileOpErr.Operation)
		}
	}
}

func TestInterfaceImplementation(t *testing.T) {
	fm := NewUnifiedFileManager()

	// Verify interface implementations
	var _ Reader = fm
	var _ Writer = fm
	var _ Validator = fm
	var _ ProjectLocator = fm
	var _ SQLCRunner = fm
	var _ Manager = fm
	var _ EnhancedFileManager = fm
}
