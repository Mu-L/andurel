package upgrade

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func setupGitRepo(t *testing.T) string {
	tmpDir := t.TempDir()
	
	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git: %v", err)
	}
	
	// Configure git
	exec.Command("git", "config", "user.email", "test@example.com").Run()
	exec.Command("git", "config", "user.name", "Test User").Run()
	
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	cmd.Run()
	
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	cmd.Run()
	
	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("initial"), 0644)
	
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = tmpDir
	cmd.Run()
	
	cmd = exec.Command("git", "commit", "-m", "initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to create initial commit: %v", err)
	}
	
	return tmpDir
}

func TestNewGitAnalyzer(t *testing.T) {
	analyzer := NewGitAnalyzer("/test/path")
	
	if analyzer == nil {
		t.Fatal("expected non-nil analyzer")
	}
	
	if analyzer.projectRoot != "/test/path" {
		t.Errorf("expected projectRoot /test/path, got %s", analyzer.projectRoot)
	}
}

func TestGitAnalyzer_IsClean(t *testing.T) {
	tmpDir := setupGitRepo(t)
	analyzer := NewGitAnalyzer(tmpDir)
	
	// Should be clean initially
	clean, err := analyzer.IsClean()
	if err != nil {
		t.Fatalf("IsClean failed: %v", err)
	}
	if !clean {
		t.Error("expected clean repository")
	}
	
	// Add a new file
	testFile := filepath.Join(tmpDir, "newfile.txt")
	os.WriteFile(testFile, []byte("new content"), 0644)
	
	// Should not be clean now
	clean, err = analyzer.IsClean()
	if err != nil {
		t.Fatalf("IsClean failed: %v", err)
	}
	if clean {
		t.Error("expected dirty repository")
	}
}

func TestGitAnalyzer_IsClean_NotGitRepo(t *testing.T) {
	tmpDir := t.TempDir()
	analyzer := NewGitAnalyzer(tmpDir)
	
	_, err := analyzer.IsClean()
	if err == nil {
		t.Error("expected error for non-git repository")
	}
}

func TestGitAnalyzer_GetFirstCommit(t *testing.T) {
	tmpDir := setupGitRepo(t)
	analyzer := NewGitAnalyzer(tmpDir)
	
	commit, err := analyzer.getFirstCommit()
	if err != nil {
		t.Fatalf("getFirstCommit failed: %v", err)
	}
	
	if commit == "" {
		t.Error("expected non-empty commit hash")
	}
	
	if len(commit) != 40 {
		t.Errorf("expected 40 character commit hash, got %d characters", len(commit))
	}
}

func TestGitAnalyzer_GetModifiedFiles(t *testing.T) {
	tmpDir := setupGitRepo(t)
	analyzer := NewGitAnalyzer(tmpDir)
	
	// Modify the initial file
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("modified content"), 0644)
	
	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	cmd.Run()
	
	cmd = exec.Command("git", "commit", "-m", "modify test.txt")
	cmd.Dir = tmpDir
	cmd.Run()
	
	// Add a new file
	newFile := filepath.Join(tmpDir, "new.txt")
	os.WriteFile(newFile, []byte("new file"), 0644)
	
	cmd = exec.Command("git", "add", "new.txt")
	cmd.Dir = tmpDir
	cmd.Run()
	
	cmd = exec.Command("git", "commit", "-m", "add new.txt")
	cmd.Dir = tmpDir
	cmd.Run()
	
	// Get modified files
	modifiedFiles, err := analyzer.GetModifiedFiles()
	if err != nil {
		t.Fatalf("GetModifiedFiles failed: %v", err)
	}
	
	// Should include both test.txt and new.txt
	if !modifiedFiles["test.txt"] {
		t.Error("expected test.txt to be in modified files")
	}
	
	if !modifiedFiles["new.txt"] {
		t.Error("expected new.txt to be in modified files")
	}
}

func TestGitAnalyzer_GetFileFromInitialCommit(t *testing.T) {
	tmpDir := setupGitRepo(t)
	analyzer := NewGitAnalyzer(tmpDir)
	
	// Get the initial version of test.txt
	content, err := analyzer.GetFileFromInitialCommit("test.txt")
	if err != nil {
		t.Fatalf("GetFileFromInitialCommit failed: %v", err)
	}
	
	if string(content) != "initial" {
		t.Errorf("expected content 'initial', got %q", string(content))
	}
	
	// Modify the file
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("modified"), 0644)
	
	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	cmd.Run()
	
	cmd = exec.Command("git", "commit", "-m", "modify")
	cmd.Dir = tmpDir
	cmd.Run()
	
	// Should still get initial version
	content, err = analyzer.GetFileFromInitialCommit("test.txt")
	if err != nil {
		t.Fatalf("GetFileFromInitialCommit failed: %v", err)
	}
	
	if string(content) != "initial" {
		t.Errorf("expected initial content 'initial', got %q", string(content))
	}
}

func TestGitAnalyzer_GetFileFromInitialCommit_NonExistent(t *testing.T) {
	tmpDir := setupGitRepo(t)
	analyzer := NewGitAnalyzer(tmpDir)
	
	// Try to get a file that doesn't exist
	content, err := analyzer.GetFileFromInitialCommit("nonexistent.txt")
	if err != nil {
		t.Fatalf("GetFileFromInitialCommit failed: %v", err)
	}
	
	// Should return nil for non-existent files
	if content != nil {
		t.Errorf("expected nil content for non-existent file, got %q", string(content))
	}
}
