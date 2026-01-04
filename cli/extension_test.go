package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestNewExtensionCommand(t *testing.T) {
	cmd := newExtensionCommand()

	if cmd == nil {
		t.Fatal("newExtensionCommand returned nil")
	}

	if cmd.Use != "extension" {
		t.Errorf("Expected Use 'extension', got %q", cmd.Use)
	}

	expectedAliases := []string{"ext", "e"}
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

func TestNewExtensionAddCommand(t *testing.T) {
	cmd := newExtensionAddCommand()

	if cmd == nil {
		t.Fatal("newExtensionAddCommand returned nil")
	}

	if cmd.Use != "add [extension-name]" {
		t.Errorf("Expected Use 'add [extension-name]', got %q", cmd.Use)
	}

	if cmd.Args == nil {
		t.Error("Args should not be nil")
	}

	if cmd.RunE == nil {
		t.Error("RunE should not be nil")
	}
}

func TestNewExtensionListCommand(t *testing.T) {
	cmd := newExtensionListCommand()

	if cmd == nil {
		t.Fatal("newExtensionListCommand returned nil")
	}

	if cmd.Use != "list" {
		t.Errorf("Expected Use 'list', got %q", cmd.Use)
	}

	expectedAliases := []string{"ls"}
	if len(cmd.Aliases) != len(expectedAliases) {
		t.Errorf("Expected %d aliases, got %d", len(expectedAliases), len(cmd.Aliases))
	}

	if cmd.Args == nil {
		t.Error("Args should not be nil")
	}

	if cmd.RunE == nil {
		t.Error("RunE should not be nil")
	}
}

func TestExtensionCommandStructure(t *testing.T) {
	rootCmd := NewRootCommand("test", "test-date")

	var extCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "extension" {
			extCmd = cmd
			break
		}
	}

	if extCmd == nil {
		t.Fatal("extension command not found")
	}

	for _, cmd := range extCmd.Commands() {
		t.Logf("Found extension subcommand: %q", cmd.Use)
	}

	subcommands := make(map[string]bool)
	for _, cmd := range extCmd.Commands() {
		subcommands[cmd.Use] = true
	}

	expectedSubcommands := []string{"add [extension-name]", "list"}

	for _, expected := range expectedSubcommands {
		if !subcommands[expected] {
			t.Errorf("extension subcommand %q not found", expected)
		}
	}
}

func TestExtensionAdd_NoGoMod(t *testing.T) {
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

	cmd := newExtensionAddCommand()
	args := []string{"test-ext"}

	cmd.SetArgs(args)
	err = cmd.RunE(cmd, args)

	if err == nil {
		t.Error("Expected error when no go.mod found")
	}
}

func TestExtensionList_NoGoMod(t *testing.T) {
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

	cmd := newExtensionListCommand()
	args := []string{}

	cmd.SetArgs(args)
	err = cmd.RunE(cmd, args)

	if err == nil {
		t.Error("Expected error when no go.mod found")
	}
}

func TestExtensionAdd_NoLockFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	cmd := newExtensionAddCommand()
	args := []string{"test-ext"}

	cmd.SetArgs(args)
	err = cmd.RunE(cmd, args)

	if err == nil {
		t.Error("Expected error when andurel.lock not found")
	}
}

func TestExtensionList_NoLockFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "andurel-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	goModPath := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goModPath, []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	origDir, _ := os.Getwd()
	defer os.Chdir(origDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to chdir: %v", err)
	}

	cmd := newExtensionListCommand()
	args := []string{}

	cmd.SetArgs(args)
	err = cmd.RunE(cmd, args)

	if err == nil {
		t.Error("Expected error when andurel.lock not found")
	}
}
