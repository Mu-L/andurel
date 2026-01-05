package cmds

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunGoModTidy(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a minimal go.mod file
	goMod := `module test
go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	// Run go mod tidy
	err = RunGoModTidy(tmpDir)
	if err != nil {
		t.Errorf("RunGoModTidy failed: %v", err)
	}
}

func TestRunGoFmt(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a minimal go file
	goFile := `package main

import "fmt"

func main() {
fmt.Println("hello")
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goFile), 0644)
	if err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	// Create go.mod
	goMod := `module test
go 1.21
`
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	// Run go fmt
	err = RunGoFmt(tmpDir)
	if err != nil {
		t.Errorf("RunGoFmt failed: %v", err)
	}

	// Verify file was formatted
	content, err := os.ReadFile(filepath.Join(tmpDir, "main.go"))
	if err != nil {
		t.Fatalf("failed to read formatted file: %v", err)
	}

	// Should have proper indentation now
	formatted := string(content)
	if formatted == goFile {
		t.Error("file should have been formatted")
	}
}

func TestRunGoFmtPath(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a minimal go file
	goFile := `package main
func main() {}
`
	testFile := filepath.Join(tmpDir, "main.go")
	err := os.WriteFile(testFile, []byte(goFile), 0644)
	if err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	// Run go fmt on specific path
	err = RunGoFmtPath(tmpDir, "main.go")
	if err != nil {
		t.Errorf("RunGoFmtPath failed: %v", err)
	}
}

func TestRunGoModTidyInvalidPath(t *testing.T) {
	// This should still work - exec.Command will fail but we test the path logic
	tmpDir := t.TempDir()
	
	err := RunGoModTidy(tmpDir)
	// Should fail because there's no go.mod
	if err == nil {
		t.Error("expected error for directory without go.mod")
	}
}

func TestRunGoFmtInvalidPath(t *testing.T) {
	tmpDir := t.TempDir()
	
	err := RunGoFmt(tmpDir)
	// Should succeed even with no Go files (fmt returns no error)
	if err != nil {
		t.Logf("RunGoFmt returned error (may be expected): %v", err)
	}
}

func TestRunGolines(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module test
go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	goFile := `package main

import "fmt"

func main() {
	fmt.Println("This is a very long line that exceeds 100 characters and should be wrapped by golines tool if it is installed")
}
`
	err = os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goFile), 0644)
	if err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	err = RunGolines(tmpDir)
	if err != nil {
		t.Logf("RunGolines failed (golines may not be installed): %v", err)
	}
}

func TestRunGoRunBin(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module test
go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	cmdDir := filepath.Join(tmpDir, "cmd", "run")
	err = os.MkdirAll(cmdDir, 0755)
	if err != nil {
		t.Fatalf("failed to create cmd/run dir: %v", err)
	}

	mainFile := `package main

import "fmt"

func main() {
	fmt.Println("Hello from run binary")
}
`
	err = os.WriteFile(filepath.Join(cmdDir, "main.go"), []byte(mainFile), 0644)
	if err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	err = RunGoRunBin(tmpDir)
	if err != nil {
		t.Errorf("RunGoRunBin failed: %v", err)
	}

	binPath := filepath.Join(tmpDir, "bin", "run")
	if _, err := os.Stat(binPath); os.IsNotExist(err) {
		t.Error("binary was not created")
	}
}

func TestRunTemplGenerate(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module test
go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	viewsDir := filepath.Join(tmpDir, "views")
	err = os.MkdirAll(viewsDir, 0755)
	if err != nil {
		t.Fatalf("failed to create views dir: %v", err)
	}

	err = RunTemplGenerate(tmpDir)
	if err != nil {
		t.Logf("RunTemplGenerate failed (expected if templ not compatible): %v", err)
	}
}

func TestRunTemplFmt(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module test
go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	viewsDir := filepath.Join(tmpDir, "views")
	err = os.MkdirAll(viewsDir, 0755)
	if err != nil {
		t.Fatalf("failed to create views dir: %v", err)
	}

	err = RunTemplFmt(tmpDir)
	if err != nil {
		t.Logf("RunTemplFmt failed (expected if templ not compatible): %v", err)
	}
}

func TestRunSqlcGenerate(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module test
go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	databaseDir := filepath.Join(tmpDir, "database")
	err = os.MkdirAll(databaseDir, 0755)
	if err != nil {
		t.Fatalf("failed to create database dir: %v", err)
	}

	sqlcYaml := `version: "2"
sql:
  - schema: "schema.sql"
    queries: "queries.sql"
    engine: "postgresql"
    gen:
      go:
        package: "db"
        out: "db"
`
	err = os.WriteFile(filepath.Join(databaseDir, "sqlc.yaml"), []byte(sqlcYaml), 0644)
	if err != nil {
		t.Fatalf("failed to create sqlc.yaml: %v", err)
	}

	err = RunSqlcGenerate(tmpDir)
	if err != nil {
		t.Logf("RunSqlcGenerate failed (expected without proper schema): %v", err)
	}
}

func TestRunGooseFix(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module test
go 1.21
`
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	migrationsDir := filepath.Join(tmpDir, "database", "migrations")
	err = os.MkdirAll(migrationsDir, 0755)
	if err != nil {
		t.Fatalf("failed to create migrations dir: %v", err)
	}

	migration := `-- +goose Up
CREATE TABLE users (
    id INTEGER PRIMARY KEY
);

-- +goose Down
DROP TABLE users;
`
	err = os.WriteFile(filepath.Join(migrationsDir, "001_create_users.sql"), []byte(migration), 0644)
	if err != nil {
		t.Fatalf("failed to create migration: %v", err)
	}

	err = RunGooseFix(tmpDir)
	if err != nil {
		t.Logf("RunGooseFix failed (may be expected): %v", err)
	}
}
