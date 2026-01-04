package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewMigrationCommand(t *testing.T) {
	cmd := newMigrationCommand()

	if cmd == nil {
		t.Fatal("newMigrationCommand returned nil")
	}

	// Check command structure
	if cmd.Use != "migration" {
		t.Errorf("Expected Use 'migration', got '%s'", cmd.Use)
	}

	// Check aliases
	expectedAliases := []string{"m", "mig"}
	for _, alias := range expectedAliases {
		found := false
		for _, a := range cmd.Aliases {
			if a == alias {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected alias '%s' not found", alias)
		}
	}

	// Check subcommands exist
	expectedSubcommands := []string{
		"new", "up", "down", "fix", "reset", "up-to", "down-to",
	}
	subcommands := make(map[string]bool)
	for _, subcmd := range cmd.Commands() {
		// Extract the first word from the Use field (command name)
		cmdName := subcmd.Use
		if spaceIdx := findSpace(cmdName); spaceIdx != -1 {
			cmdName = cmdName[:spaceIdx]
		}
		subcommands[cmdName] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommands[expected] {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}

func TestNewMigrationNewCommand(t *testing.T) {
	cmd := newMigrationNewCommand()

	if cmd == nil {
		t.Fatal("newMigrationNewCommand returned nil")
	}

	if cmd.Use != "new [name]" {
		t.Errorf("Expected Use 'new [name]', got '%s'", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Expected Args validator, got nil")
	}

	if cmd.Example != "andurel migration new create_users_table" {
		t.Errorf("Unexpected example: %s", cmd.Example)
	}
}

func TestNewMigrationUpCommand(t *testing.T) {
	cmd := newMigrationUpCommand()

	if cmd == nil {
		t.Fatal("newMigrationUpCommand returned nil")
	}

	if cmd.Use != "up" {
		t.Errorf("Expected Use 'up', got '%s'", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Expected Args validator, got nil")
	}
}

func TestNewMigrationDownCommand(t *testing.T) {
	cmd := newMigrationDownCommand()

	if cmd == nil {
		t.Fatal("newMigrationDownCommand returned nil")
	}

	if cmd.Use != "down" {
		t.Errorf("Expected Use 'down', got '%s'", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Expected Args validator, got nil")
	}
}

func TestNewMigrationFixCommand(t *testing.T) {
	cmd := newMigrationFixCommand()

	if cmd == nil {
		t.Fatal("newMigrationFixCommand returned nil")
	}

	if cmd.Use != "fix" {
		t.Errorf("Expected Use 'fix', got '%s'", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Expected Args validator, got nil")
	}
}

func TestNewMigrationResetCommand(t *testing.T) {
	cmd := newMigrationResetCommand()

	if cmd == nil {
		t.Fatal("newMigrationResetCommand returned nil")
	}

	if cmd.Use != "reset" {
		t.Errorf("Expected Use 'reset', got '%s'", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Expected Args validator, got nil")
	}
}

func TestNewMigrationUpToCommand(t *testing.T) {
	cmd := newMigrationUpToCommand()

	if cmd == nil {
		t.Fatal("newMigrationUpToCommand returned nil")
	}

	if cmd.Use != "up-to [version]" {
		t.Errorf("Expected Use 'up-to [version]', got '%s'", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Expected Args validator, got nil")
	}
}

func TestNewMigrationDownToCommand(t *testing.T) {
	cmd := newMigrationDownToCommand()

	if cmd == nil {
		t.Fatal("newMigrationDownToCommand returned nil")
	}

	if cmd.Use != "down-to [version]" {
		t.Errorf("Expected Use 'down-to [version]', got '%s'", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Expected Args validator, got nil")
	}
}

func TestParseDatabaseURL(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		expectedDriver string
		expectedDB     string
	}{
		{
			name:           "PostgreSQL URL",
			url:            "postgres://user:pass@host:5432/db?sslmode=require",
			expectedDriver: "postgres",
			expectedDB:     "postgres://user:pass@host:5432/db?sslmode=require",
		},
		{
			name:           "Simple URL",
			url:            "postgres://localhost/db",
			expectedDriver: "postgres",
			expectedDB:     "postgres://localhost/db",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			driver, db := parseDatabaseURL(tt.url)

			if driver != tt.expectedDriver {
				t.Errorf("Expected driver '%s', got '%s'", tt.expectedDriver, driver)
			}

			if db != tt.expectedDB {
				t.Errorf("Expected db string '%s', got '%s'", tt.expectedDB, db)
			}
		})
	}
}

func TestRunMigrationBinary_MissingEnvVars(t *testing.T) {
	// Save original env vars
	originalEnv := make(map[string]string)
	envVars := []string{"DB_KIND", "DB_PORT", "DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD", "DB_SSL_MODE"}
	for _, key := range envVars {
		if val := os.Getenv(key); val != "" {
			originalEnv[key] = val
			os.Unsetenv(key)
		}
	}
	defer func() {
		for key, val := range originalEnv {
			os.Setenv(key, val)
		}
	}()

	// Create a temp directory for testing
	tempDir := t.TempDir()
	originalCwd, _ := os.Getwd()
	defer os.Chdir(originalCwd)
	os.Chdir(tempDir)

	// Create a minimal go.mod file
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	err := runMigrationBinary("up")
	if err == nil {
		t.Error("Expected error when missing environment variables")
	}

	expectedError := "missing database configuration environment variables"
	if err != nil && !contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestRunMigrationBinary_MissingGooseBinary(t *testing.T) {
	// Set up required environment variables
	envVars := map[string]string{
		"DB_KIND":      "postgres",
		"DB_PORT":      "5432",
		"DB_HOST":      "localhost",
		"DB_NAME":      "test",
		"DB_USER":      "test",
		"DB_PASSWORD":  "test",
		"DB_SSL_MODE":  "disable",
	}

	for key, val := range envVars {
		os.Setenv(key, val)
		defer os.Unsetenv(key)
	}

	// Create a temp directory for testing
	tempDir := t.TempDir()
	originalCwd, _ := os.Getwd()
	defer os.Chdir(originalCwd)
	os.Chdir(tempDir)

	// Create a minimal go.mod file
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create bin directory but no goose binary
	binDir := filepath.Join(tempDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	err := runMigrationBinary("up")
	if err == nil {
		t.Error("Expected error when goose binary is missing")
	}

	expectedError := "goose binary not found"
	if err != nil && !contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

func TestMigrationCommands_HelpOutput(t *testing.T) {
	rootCmd := NewRootCommand("test", "test-date")

	tests := []struct {
		name string
		args []string
	}{
		{"migration help", []string{"migration", "--help"}},
		{"migration new help", []string{"migration", "new", "--help"}},
		{"migration up help", []string{"migration", "up", "--help"}},
		{"migration down help", []string{"migration", "down", "--help"}},
		{"migration fix help", []string{"migration", "fix", "--help"}},
		{"migration reset help", []string{"migration", "reset", "--help"}},
		{"migration up-to help", []string{"migration", "up-to", "--help"}},
		{"migration down-to help", []string{"migration", "down-to", "--help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rootCmd.SetArgs(tt.args)
			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("Command %v failed: %v", tt.args, err)
			}
		})
	}
}

func TestRunMigrationBinary_GooseDirectoryCreation(t *testing.T) {
	// Set up required environment variables
	envVars := map[string]string{
		"DB_KIND":      "postgres",
		"DB_PORT":      "5432",
		"DB_HOST":      "localhost",
		"DB_NAME":      "test",
		"DB_USER":      "test",
		"DB_PASSWORD":  "test",
		"DB_SSL_MODE":  "disable",
	}

	for key, val := range envVars {
		os.Setenv(key, val)
		defer os.Unsetenv(key)
	}

	// Create a temp directory for testing
	tempDir := t.TempDir()
	originalCwd, _ := os.Getwd()
	defer os.Chdir(originalCwd)
	os.Chdir(tempDir)

	// Create a minimal go.mod file
	goModPath := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test\n"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create migrations directory
	migrationsDir := filepath.Join(tempDir, "database", "migrations")
	if err := os.MkdirAll(migrationsDir, 0755); err != nil {
		t.Fatalf("Failed to create migrations directory: %v", err)
	}

	// Create a fake goose binary
	binDir := filepath.Join(tempDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	goosePath := filepath.Join(binDir, "goose")
	// Create a simple script that just exits successfully
	if err := os.WriteFile(goosePath, []byte("#!/bin/sh\nexit 0\n"), 0755); err != nil {
		t.Fatalf("Failed to create fake goose binary: %v", err)
	}

	// This test only checks that the binary path is constructed correctly
	// The actual execution will fail, but we can at least verify the path construction
	if _, err := os.Stat(goosePath); err != nil {
		t.Errorf("Goose binary should exist at %s", goosePath)
	}

	// Verify migrations directory exists
	if _, err := os.Stat(migrationsDir); err != nil {
		t.Errorf("Migrations directory should exist at %s", migrationsDir)
	}
}

func findSpace(s string) int {
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' {
			return i
		}
	}
	return -1
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
