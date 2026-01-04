package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewSqlcCommand(t *testing.T) {
	cmd := newSqlcCommand()

	if cmd == nil {
		t.Fatal("newSqlcCommand returned nil")
	}

	if cmd.Use != "sqlc" {
		t.Errorf("Expected Use 'sqlc', got %q", cmd.Use)
	}

	expectedAliases := []string{"s"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}
	for i, alias := range expectedAliases {
		if cmd.Aliases[i] != alias {
			t.Errorf("Expected alias %q, got %q", alias, cmd.Aliases[i])
		}
	}

	if len(cmd.Commands()) != 2 {
		t.Errorf("Expected 2 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestNewSqlcCompileCommand(t *testing.T) {
	cmd := newSqlcCompileCommand()

	if cmd == nil {
		t.Fatal("newSqlcCompileCommand returned nil")
	}

	if cmd.Use != "compile" {
		t.Errorf("Expected Use 'compile', got %q", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Args should not be nil")
	}

	if cmd.RunE == nil {
		t.Error("RunE should not be nil")
	}
}

func TestNewSqlcGenerateCommand(t *testing.T) {
	cmd := newSqlcGenerateCommand()

	if cmd == nil {
		t.Fatal("newSqlcGenerateCommand returned nil")
	}

	if cmd.Use != "generate" {
		t.Errorf("Expected Use 'generate', got %q", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Args should not be nil")
	}

	if cmd.RunE == nil {
		t.Error("RunE should not be nil")
	}
}

func TestNewTemplCommand(t *testing.T) {
	cmd := newTemplCommand()

	if cmd == nil {
		t.Fatal("newTemplCommand returned nil")
	}

	if cmd.Use != "templ" {
		t.Errorf("Expected Use 'templ', got %q", cmd.Use)
	}

	expectedAliases := []string{"t"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}
	for i, alias := range expectedAliases {
		if cmd.Aliases[i] != alias {
			t.Errorf("Expected alias %q, got %q", alias, cmd.Aliases[i])
		}
	}

	if len(cmd.Commands()) != 1 {
		t.Errorf("Expected 1 subcommand, got %d", len(cmd.Commands()))
	}
}

func TestNewTemplGenerateCommand(t *testing.T) {
	cmd := newTemplGenerateCommand()

	if cmd == nil {
		t.Fatal("newTemplGenerateCommand returned nil")
	}

	if cmd.Use != "generate" {
		t.Errorf("Expected Use 'generate', got %q", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Args should not be nil")
	}

	if cmd.RunE == nil {
		t.Error("RunE should not be nil")
	}
}

func TestRunSqlc_NoGoMod(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	err = runSqlc("compile")
	if err == nil {
		t.Error("Expected error when no go.mod found")
	}
}

func TestRunSqlc_NoConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	err = runSqlc("compile")
	if err == nil {
		t.Error("Expected error when sqlc config not found")
	}
}

func TestRunTempl_NoGoMod(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	err = runTempl("generate")
	if err == nil {
		t.Error("Expected error when no go.mod found")
	}
}

func TestSqlcAndTemplCommandStructure(t *testing.T) {
	rootCmd := NewRootCommand("test", "test-date")

	var sqlcCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "sqlc" {
			sqlcCmd = cmd
			break
		}
	}

	if sqlcCmd == nil {
		t.Fatal("sqlc command not found")
	}

	expectedSubcommands := []string{"compile", "generate"}
	subcommands := make(map[string]bool)
	for _, cmd := range sqlcCmd.Commands() {
		subcommands[cmd.Use] = true
	}

	for _, expected := range expectedSubcommands {
		if !subcommands[expected] {
			t.Errorf("sqlc subcommand %q not found", expected)
		}
	}

	var templCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "templ" {
			templCmd = cmd
			break
		}
	}

	if templCmd == nil {
		t.Fatal("templ command not found")
	}

	expectedTemplSubcommands := []string{"generate"}
	templSubcommands := make(map[string]bool)
	for _, cmd := range templCmd.Commands() {
		templSubcommands[cmd.Use] = true
	}

	for _, expected := range expectedTemplSubcommands {
		if !templSubcommands[expected] {
			t.Errorf("templ subcommand %q not found", expected)
		}
	}
}
