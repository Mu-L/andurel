package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewDoctorCommand(t *testing.T) {
	cmd := newDoctorCommand()

	if cmd == nil {
		t.Fatal("newDoctorCommand returned nil")
	}

	if cmd.Use != "doctor" {
		t.Errorf("Expected Use 'doctor', got '%s'", cmd.Use)
	}

	if cmd.Short != "Run diagnostic checks on your Andurel project" {
		t.Errorf("Unexpected Short description: %s", cmd.Short)
	}

	// Check that verbose flag exists
	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("Expected --verbose flag to exist")
	}

	if verboseFlag != nil && verboseFlag.DefValue != "false" {
		t.Errorf("Expected --verbose default to be 'false', got '%s'", verboseFlag.DefValue)
	}
}

func TestCheckGoVersion(t *testing.T) {
	result := checkGoVersion()

	if result.name != "Go version" {
		t.Errorf("checkGoVersion() name = %s, want 'Go version'", result.name)
	}

	if result.status != statusPass {
		t.Errorf("checkGoVersion() status = %d, want %d", result.status, statusPass)
	}

	if result.message == "" {
		t.Error("checkGoVersion() message should not be empty")
	}

	// Message should contain something like "go1.x.x"
	if !strings.Contains(result.message, "go") {
		t.Errorf("checkGoVersion() message = %s, should contain 'go'", result.message)
	}
}

func TestCheckInAndurelProject(t *testing.T) {
	// Test in current directory (should be an Andurel project)
	result := checkInAndurelProject()

	if result.name != "Andurel project" {
		t.Errorf("checkInAndurelProject() name = %s, want 'Andurel project'", result.name)
	}

	// In the andurel repo itself, this should pass
	// (assuming we're running tests in the andurel directory)
}

func TestCheckInAndurelProject_NotInProject(t *testing.T) {
	// Skip this test if we can't create an independent temp directory
	// because findGoModRoot will search parent directories
	t.Skip("Skipping test due to go.mod file detection in parent directories")

	// Create a temp directory that is NOT an Andurel project
	tempDir := t.TempDir()

	// Save current directory and change to temp dir
	originalCwd, _ := os.Getwd()
	defer os.Chdir(originalCwd)
	os.Chdir(tempDir)

	result := checkInAndurelProject()

	if result.name != "Andurel project" {
		t.Errorf("checkInAndurelProject() name = %s, want 'Andurel project'", result.name)
	}

	if result.status != statusFail {
		t.Errorf("checkInAndurelProject() status = %d, want %d", result.status, statusFail)
	}

	if result.message == "" {
		t.Error("checkInAndurelProject() message should not be empty when not in project")
	}
}

func TestCheckLockFile_Missing(t *testing.T) {
	tempDir := t.TempDir()

	result := checkLockFile(tempDir)

	if result.name != "andurel.lock" {
		t.Errorf("checkLockFile() name = %s, want 'andurel.lock'", result.name)
	}

	if result.status != statusFail {
		t.Errorf("checkLockFile() status = %d, want %d", result.status, statusFail)
	}

	if !contains(result.message, "file not found") {
		t.Errorf("checkLockFile() message = %s, should contain 'file not found'", result.message)
	}
}

func TestCheckLockFile_Invalid(t *testing.T) {
	tempDir := t.TempDir()
	lockPath := filepath.Join(tempDir, "andurel.lock")

	// Create an invalid lock file
	if err := os.WriteFile(lockPath, []byte("invalid json"), 0644); err != nil {
		t.Fatalf("Failed to create invalid lock file: %v", err)
	}

	result := checkLockFile(tempDir)

	if result.name != "andurel.lock" {
		t.Errorf("checkLockFile() name = %s, want 'andurel.lock'", result.name)
	}

	if result.status != statusFail {
		t.Errorf("checkLockFile() status = %d, want %d", result.status, statusFail)
	}

	if !contains(result.message, "invalid format") {
		t.Errorf("checkLockFile() message = %s, should contain 'invalid format'", result.message)
	}
}

func TestPrintResults(t *testing.T) {
	tests := []struct {
		name     string
		results  []checkResult
		verbose  bool
	}{
		{
			name:    "pass result",
			results: []checkResult{
				{name: "test check", status: statusPass, message: "all good"},
			},
			verbose: false,
		},
		{
			name:    "warn result",
			results: []checkResult{
				{name: "test check", status: statusWarn, message: "warning"},
			},
			verbose: false,
		},
		{
			name:    "fail result",
			results: []checkResult{
				{name: "test check", status: statusFail, message: "error"},
			},
			verbose: false,
		},
		{
			name: "pass result with details and verbose",
			results: []checkResult{
				{
					name:    "test check",
					status:  statusPass,
					message: "all good",
					details: []string{"detail 1", "detail 2"},
				},
			},
			verbose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Just verify that printResults doesn't panic
			printResults(tt.results, tt.verbose)
		})
	}
}

func TestCheckResult_Type(t *testing.T) {
	// Test that checkResult is properly defined
	result := checkResult{
		name:    "test",
		status:  statusPass,
		message: "test message",
		details: []string{"detail 1"},
	}

	if result.name != "test" {
		t.Errorf("checkResult.name = %s, want 'test'", result.name)
	}

	if result.status != statusPass {
		t.Errorf("checkResult.status = %d, want %d", result.status, statusPass)
	}

	if result.message != "test message" {
		t.Errorf("checkResult.message = %s, want 'test message'", result.message)
	}

	if len(result.details) != 1 {
		t.Errorf("checkResult.details length = %d, want 1", len(result.details))
	}
}

func TestCheckStatus_Values(t *testing.T) {
	// Test that checkStatus constants have expected values
	if statusPass != 0 {
		t.Errorf("statusPass = %d, want 0", statusPass)
	}

	if statusWarn != 1 {
		t.Errorf("statusWarn = %d, want 1", statusWarn)
	}

	if statusFail != 2 {
		t.Errorf("statusFail = %d, want 2", statusFail)
	}
}

func TestDoctor_CommandStructure(t *testing.T) {
	cmd := newDoctorCommand()

	// Check that it's a valid cobra command
	if cmd.RunE == nil {
		t.Error("Expected RunE function, got nil")
	}

	// Check that verbose flag works
	if cmd.Flags().Lookup("verbose") == nil {
		t.Error("Expected --verbose flag")
	}
}

func TestDoctor_HelpOutput(t *testing.T) {
	rootCmd := NewRootCommand("test", "test-date")

	tests := []struct {
		name string
		args []string
	}{
		{"doctor help", []string{"doctor", "--help"}},
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

func TestCheckGoVet(t *testing.T) {
	rootDir, err := findGoModRoot()
	if err != nil {
		t.Skipf("Cannot find go.mod root: %v", err)
	}

	tests := []struct {
		name    string
		verbose bool
	}{
		{"verbose true", true},
		{"verbose false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkGoVet(rootDir, tt.verbose)

			if result.name != "go vet" {
				t.Errorf("checkGoVet() name = %s, want 'go vet'", result.name)
			}

			if result.status < statusPass || result.status > statusFail {
				t.Errorf("checkGoVet() status = %d, want value between %d and %d", result.status, statusPass, statusFail)
			}

			if result.message == "" {
				t.Error("checkGoVet() message should not be empty")
			}
		})
	}
}

func TestCheckGoVet_InvalidDirectory(t *testing.T) {
	tempDir := t.TempDir()

	result := checkGoVet(tempDir, false)

	if result.name != "go vet" {
		t.Errorf("checkGoVet() name = %s, want 'go vet'", result.name)
	}

	if result.status != statusFail {
		t.Errorf("checkGoVet() status = %d, want %d for invalid directory", result.status, statusFail)
	}
}

func TestCheckGoModTidy(t *testing.T) {
	rootDir, err := findGoModRoot()
	if err != nil {
		t.Skipf("Cannot find go.mod root: %v", err)
	}

	tests := []struct {
		name    string
		verbose bool
	}{
		{"verbose true", true},
		{"verbose false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := checkGoModTidy(rootDir, tt.verbose)

			if result.name != "go mod tidy" {
				t.Errorf("checkGoModTidy() name = %s, want 'go mod tidy'", result.name)
			}

			if result.status < statusPass || result.status > statusFail {
				t.Errorf("checkGoModTidy() status = %d, want value between %d and %d", result.status, statusPass, statusFail)
			}

			if result.message == "" {
				t.Error("checkGoModTidy() message should not be empty")
			}
		})
	}
}

func TestCheckGoModTidy_InvalidDirectory(t *testing.T) {
	tempDir := t.TempDir()

	result := checkGoModTidy(tempDir, false)

	if result.name != "go mod tidy" {
		t.Errorf("checkGoModTidy() name = %s, want 'go mod tidy'", result.name)
	}

	if result.status != statusFail {
		t.Errorf("checkGoModTidy() status = %d, want %d for invalid directory", result.status, statusFail)
	}

	if !contains(result.message, "cannot read go.mod") {
		t.Errorf("checkGoModTidy() message = %s, should contain 'cannot read go.mod'", result.message)
	}
}

func TestCheckTemplGenerate(t *testing.T) {
	tempDir := t.TempDir()

	result := checkTemplGenerate(tempDir, false)

	if result.name != "templ generate" {
		t.Errorf("checkTemplGenerate() name = %s, want 'templ generate'", result.name)
	}

	if result.status != statusWarn {
		t.Errorf("checkTemplGenerate() status = %d, want %d when binary not found", result.status, statusWarn)
	}

	if !contains(result.message, "templ binary not found") {
		t.Errorf("checkTemplGenerate() message = %s, should contain 'templ binary not found'", result.message)
	}
}

func TestCheckTemplGenerate_WithBinary(t *testing.T) {
	tempDir := t.TempDir()
	binDir := filepath.Join(tempDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	templPath := filepath.Join(binDir, "templ")
	if err := os.WriteFile(templPath, []byte("#!/bin/sh\necho 'templ version'"), 0755); err != nil {
		t.Fatalf("Failed to create fake templ binary: %v", err)
	}

	result := checkTemplGenerate(tempDir, false)

	if result.name != "templ generate" {
		t.Errorf("checkTemplGenerate() name = %s, want 'templ generate'", result.name)
	}

	if result.status != statusFail && result.status != statusPass {
		t.Errorf("checkTemplGenerate() status = %d, want %d or %d", result.status, statusFail, statusPass)
	}
}

func TestCheckSqlcGenerate(t *testing.T) {
	tempDir := t.TempDir()

	result := checkSqlcGenerate(tempDir, false)

	if result.name != "sqlc compile" {
		t.Errorf("checkSqlcGenerate() name = %s, want 'sqlc compile'", result.name)
	}

	if result.status != statusWarn {
		t.Errorf("checkSqlcGenerate() status = %d, want %d when binary not found", result.status, statusWarn)
	}

	if !contains(result.message, "sqlc binary not found") {
		t.Errorf("checkSqlcGenerate() message = %s, should contain 'sqlc binary not found'", result.message)
	}
}

func TestCheckSqlcGenerate_WithBinary(t *testing.T) {
	tempDir := t.TempDir()
	binDir := filepath.Join(tempDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	sqlcPath := filepath.Join(binDir, "sqlc")
	if err := os.WriteFile(sqlcPath, []byte("#!/bin/sh\necho 'sqlc version'"), 0755); err != nil {
		t.Fatalf("Failed to create fake sqlc binary: %v", err)
	}

	result := checkSqlcGenerate(tempDir, false)

	if result.name != "sqlc compile" {
		t.Errorf("checkSqlcGenerate() name = %s, want 'sqlc compile'", result.name)
	}

	if result.status != statusWarn {
		t.Errorf("checkSqlcGenerate() status = %d, want %d when config not found", result.status, statusWarn)
	}

	if !contains(result.message, "database/sqlc.yaml not found") {
		t.Errorf("checkSqlcGenerate() message = %s, should contain 'database/sqlc.yaml not found'", result.message)
	}
}

func TestCheckSqlcGenerate_WithConfig(t *testing.T) {
	tempDir := t.TempDir()
	binDir := filepath.Join(tempDir, "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		t.Fatalf("Failed to create bin directory: %v", err)
	}

	sqlcPath := filepath.Join(binDir, "sqlc")
	if err := os.WriteFile(sqlcPath, []byte("#!/bin/sh\necho 'sqlc version'"), 0755); err != nil {
		t.Fatalf("Failed to create fake sqlc binary: %v", err)
	}

	dbDir := filepath.Join(tempDir, "database")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		t.Fatalf("Failed to create database directory: %v", err)
	}

	sqlcConfigPath := filepath.Join(dbDir, "sqlc.yaml")
	if err := os.WriteFile(sqlcConfigPath, []byte("version: 2"), 0644); err != nil {
		t.Fatalf("Failed to create sqlc.yaml: %v", err)
	}

	result := checkSqlcGenerate(tempDir, false)

	if result.name != "sqlc compile" {
		t.Errorf("checkSqlcGenerate() name = %s, want 'sqlc compile'", result.name)
	}

	if result.status != statusFail && result.status != statusPass {
		t.Errorf("checkSqlcGenerate() status = %d, want %d or %d", result.status, statusFail, statusPass)
	}
}

