package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewToolCommand(t *testing.T) {
	cmd := newToolCommand()

	if cmd == nil {
		t.Fatal("newToolCommand returned nil")
	}

	if cmd.Use != "tool" {
		t.Errorf("Expected Use 'tool', got %q", cmd.Use)
	}

	if len(cmd.Commands()) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(cmd.Commands()))
	}

	for _, subcmd := range cmd.Commands() {
		t.Logf("Found subcommand: %q", subcmd.Use)
	}

	subcommands := make(map[string]bool)
	for _, subcmd := range cmd.Commands() {
		subcommands[subcmd.Use] = true
	}

	if !subcommands["sync"] {
		t.Error("Expected subcommand 'sync' not found")
	}

	setVersionFound := false
	for name := range subcommands {
		if strings.Contains(name, "set-version") {
			setVersionFound = true
			break
		}
	}

	if !setVersionFound {
		t.Error("Expected subcommand 'set-version' not found")
	}
}

func TestNewSyncCommand(t *testing.T) {
	cmd := newSyncCommand()

	if cmd == nil {
		t.Fatal("newSyncCommand returned nil")
	}

	if cmd.Use != "sync" {
		t.Errorf("Expected Use 'sync', got %q", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Args should not be nil")
	}

	if cmd.RunE == nil {
		t.Error("RunE should not be nil")
	}
}

func TestSyncBinaries_NoLockFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	err = syncBinaries(tmpDir)
	if err == nil {
		t.Error("Expected error when andurel.lock not found")
	}
}

func TestNewSetVersionCommand(t *testing.T) {
	cmd := newSetVersionCommand()

	if cmd == nil {
		t.Fatal("newSetVersionCommand returned nil")
	}

	if !strings.Contains(cmd.Use, "set-version") {
		t.Errorf("Expected Use to contain 'set-version', got %q", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Args should not be nil")
	}

	if cmd.RunE == nil {
		t.Error("RunE should not be nil")
	}
}

func TestSetVersion_UnknownTool(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	lockPath := filepath.Join(tmpDir, "andurel.lock")
	lockContent := `{
  "tools": {}
}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0644); err != nil {
		t.Fatalf("Failed to create andurel.lock: %v", err)
	}

	err = setVersion(tmpDir, "unknown_tool", "1.0.0")
	if err == nil {
		t.Error("Expected error for unknown tool")
	}
}

func TestSetVersion_EmptyVersion(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	lockPath := filepath.Join(tmpDir, "andurel.lock")
	lockContent := `{
  "tools": {}
}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0644); err != nil {
		t.Fatalf("Failed to create andurel.lock: %v", err)
	}

	err = setVersion(tmpDir, "templ", "")
	if err == nil {
		t.Error("Expected error for empty version")
	}
}

func TestSetVersion_NoLockFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	err = setVersion(tmpDir, "templ", "1.0.0")
	if err == nil {
		t.Error("Expected error when andurel.lock not found")
	}
}

func TestExtractModulePath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "github.com/a-h/templ/cmd/templ",
			expected: "github.com/a-h/templ",
		},
		{
			input:    "github.com/sqlc-dev/sqlc/cmd/sqlc",
			expected: "github.com/sqlc-dev/sqlc",
		},
		{
			input:    "short/path",
			expected: "short/path",
		},
		{
			input:    "single",
			expected: "single",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractModulePath(tt.input)
			if got != tt.expected {
				t.Errorf("extractModulePath(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestSetVersion_VersionPrefixHandling(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	lockPath := filepath.Join(tmpDir, "andurel.lock")
	lockContent := `{
  "tools": {}
}`
	if err := os.WriteFile(lockPath, []byte(lockContent), 0644); err != nil {
		t.Fatalf("Failed to create andurel.lock: %v", err)
	}

	tests := []struct {
		name    string
		version string
		wantErr bool
	}{
		{
			name:    "version with v prefix",
			version: "v1.0.0",
			wantErr: true,
		},
		{
			name:    "version without prefix",
			version: "1.0.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := setVersion(tmpDir, "templ", tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("setVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestToolCommandStructure(t *testing.T) {
	rootCmd := NewRootCommand("test", "test-date")

	var toolCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "tool" {
			toolCmd = cmd
			break
		}
	}

	if toolCmd == nil {
		t.Fatal("tool command not found")
	}

	subcommands := make(map[string]bool)
	for _, cmd := range toolCmd.Commands() {
		subcommands[cmd.Use] = true
	}

	if !subcommands["sync"] {
		t.Error("tool subcommand 'sync' not found")
	}

	setVersionFound := false
	for name := range subcommands {
		if strings.Contains(name, "set-version") {
			setVersionFound = true
			break
		}
	}

	if !setVersionFound {
		t.Error("tool subcommand 'set-version' not found")
	}
}
