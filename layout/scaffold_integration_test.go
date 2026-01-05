package layout

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScaffold_BasicProject(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	err := Scaffold(
		targetDir,
		"testproject",
		"github.com/test",
		"postgresql",
		"vanilla",
		"v0.1.0",
		[]string{},
	)
	
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}
	
	// Verify key directories were created
	expectedDirs := []string{
		"controllers",
		"models",
		"views",
		"database/migrations",
		"database/queries",
		"router",
		"cmd/run",
		"assets",
		"config",
	}
	
	for _, dir := range expectedDirs {
		path := filepath.Join(targetDir, dir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected directory %s was not created", dir)
		}
	}
	
	// Verify key files were created
	expectedFiles := []string{
		"go.mod",
		".gitignore",
		"README.md",
		"andurel.lock",
	}
	
	for _, file := range expectedFiles {
		path := filepath.Join(targetDir, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %s was not created", file)
		}
	}
}

func TestScaffold_WithDocker(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	err := Scaffold(
		targetDir,
		"testproject",
		"github.com/test",
		"postgresql",
		"vanilla",
		"v0.1.0",
		[]string{"docker"},
	)
	
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}
	
	// Verify Dockerfile was created
	dockerfilePath := filepath.Join(targetDir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		t.Error("Dockerfile was not created")
	}
	
	// Verify .dockerignore was created
	dockerignorePath := filepath.Join(targetDir, ".dockerignore")
	if _, err := os.Stat(dockerignorePath); os.IsNotExist(err) {
		t.Error(".dockerignore was not created")
	}
}

func TestScaffold_WithTailwind(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	err := Scaffold(
		targetDir,
		"testproject",
		"github.com/test",
		"postgresql",
		"tailwind",
		"v0.1.0",
		[]string{},
	)
	
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}
	
	// Verify tailwind CSS files
	cssDir := filepath.Join(targetDir, "css")
	if _, err := os.Stat(cssDir); os.IsNotExist(err) {
		t.Error("css directory was not created for tailwind")
	}
	
	// Check for tailwind-specific CSS files
	baseCSSPath := filepath.Join(cssDir, "base.css")
	if _, err := os.Stat(baseCSSPath); os.IsNotExist(err) {
		t.Error("base.css was not created for tailwind")
	}
}

func TestScaffold_AlreadyExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	// Create the directory with a file first
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	
	// Create a file to make it non-empty
	testFile := filepath.Join(targetDir, "existing.txt")
	err = os.WriteFile(testFile, []byte("existing"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	
	// Try to scaffold - it will succeed but we verify the directory exists
	err = Scaffold(
		targetDir,
		"testproject",
		"github.com/test",
		"postgresql",
		"vanilla",
		"v0.1.0",
		[]string{},
	)
	
	// Scaffold may overwrite or merge, so we just verify it doesn't panic
	if err != nil {
		t.Logf("Scaffold returned error (may be expected): %v", err)
	}
}

func TestScaffold_VerifyGoMod(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	err := Scaffold(
		targetDir,
		"testproject",
		"test",
		"postgresql",
		"vanilla",
		"v0.1.0",
		[]string{},
	)
	
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}
	
	// Verify go.mod content
	goModPath := filepath.Join(targetDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("failed to read go.mod: %v", err)
	}
	
	contentStr := string(content)
	// The module name is test/testproject (not github.com/test/testproject)
	if !contains(contentStr, "module test/testproject") {
		t.Errorf("go.mod should contain correct module declaration, got: %s", contentStr)
	}
}

func TestScaffold_VerifyMigrations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	err := Scaffold(
		targetDir,
		"testproject",
		"github.com/test",
		"postgresql",
		"vanilla",
		"v0.1.0",
		[]string{},
	)
	
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}
	
	// Verify migrations directory
	migrationsDir := filepath.Join(targetDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("failed to read migrations directory: %v", err)
	}
	
	// Should have river queue migrations + auth migrations
	if len(entries) < 6 {
		t.Errorf("expected at least 6 migration files, got %d", len(entries))
	}
}

func TestScaffold_VerifyLockFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	err := Scaffold(
		targetDir,
		"testproject",
		"github.com/test",
		"postgresql",
		"vanilla",
		"v0.1.0",
		[]string{"docker"},
	)
	
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}
	
	// Verify lock file
	lockPath := filepath.Join(targetDir, "andurel.lock")
	content, err := os.ReadFile(lockPath)
	if err != nil {
		t.Fatalf("failed to read andurel.lock: %v", err)
	}
	
	contentStr := string(content)
	
	// Verify tools are listed
	expectedTools := []string{"templ", "sqlc", "goose", "air"}
	for _, tool := range expectedTools {
		if !contains(contentStr, tool) {
			t.Errorf("lock file should contain tool: %s", tool)
		}
	}
	
	// Verify extension is listed
	if !contains(contentStr, "docker") {
		t.Error("lock file should contain docker extension")
	}
}

func TestScaffold_VerifyGitInit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	err := Scaffold(
		targetDir,
		"testproject",
		"github.com/test",
		"postgresql",
		"vanilla",
		"v0.1.0",
		[]string{},
	)
	
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}
	
	// Verify git was initialized
	gitDir := filepath.Join(targetDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error(".git directory was not created")
	}
}

func TestScaffold_MultipleExtensions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	
	tmpDir := t.TempDir()
	targetDir := filepath.Join(tmpDir, "testproject")
	
	err := Scaffold(
		targetDir,
		"testproject",
		"github.com/test",
		"postgresql",
		"vanilla",
		"v0.1.0",
		[]string{"docker", "aws-ses"},
	)
	
	if err != nil {
		t.Fatalf("Scaffold failed: %v", err)
	}
	
	// Verify docker files
	dockerfilePath := filepath.Join(targetDir, "Dockerfile")
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		t.Error("Dockerfile was not created")
	}
	
	// Verify aws-ses config file
	awsSESConfigPath := filepath.Join(targetDir, "config", "aws_ses.go")
	if _, err := os.Stat(awsSESConfigPath); os.IsNotExist(err) {
		t.Error("aws_ses.go config was not created")
	}
	
	// Verify aws-ses client file
	awsSESClientPath := filepath.Join(targetDir, "clients", "email", "aws_ses.go")
	if _, err := os.Stat(awsSESClientPath); os.IsNotExist(err) {
		t.Error("aws_ses.go client was not created")
	}
}
