package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCommand(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("successful command", func(t *testing.T) {
		err := RunCommand(t, "echo", tempDir, nil, "test")
		if err != nil {
			t.Errorf("RunCommand should succeed for 'echo test', got: %v", err)
		}
	})

	t.Run("command with args", func(t *testing.T) {
		err := RunCommand(t, "echo", tempDir, nil, "arg1", "arg2", "arg3")
		if err != nil {
			t.Errorf("RunCommand should succeed for 'echo' with multiple args, got: %v", err)
		}
	})

	t.Run("command in specific directory", func(t *testing.T) {
		testDir := filepath.Join(tempDir, "testdir")
		if err := os.Mkdir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Create a file in testdir
		testFile := filepath.Join(testDir, "test.txt")
		if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Run a command that checks if the file exists in the working directory
		// Using 'ls' which should succeed in a directory with files
		err := RunCommand(t, "ls", testDir, nil)
		if err != nil {
			t.Errorf("RunCommand should succeed in test directory, got: %v", err)
		}
	})

	t.Run("command with custom env", func(t *testing.T) {
		testFile := filepath.Join(tempDir, "envtest.txt")

		// On Unix-like systems, create a script that checks env var
		script := `#!/bin/sh
if [ -n "$TEST_VAR" ]; then
  echo "$TEST_VAR" > ` + testFile + `
fi
`

		scriptPath := filepath.Join(tempDir, "test.sh")
		if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
			t.Fatalf("Failed to create test script: %v", err)
		}

		env := []string{"TEST_VAR=test_value"}
		err := RunCommand(t, scriptPath, tempDir, env)
		if err != nil {
			t.Errorf("RunCommand should succeed with custom env, got: %v", err)
		}

		// Check that the env var was set
		if _, err := os.Stat(testFile); err != nil {
			t.Error("Custom environment variable was not set")
		}
	})
}

func TestRunCommandOutput(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("successful command with output", func(t *testing.T) {
		output, err := RunCommandOutput(t, "echo", tempDir, nil, "test output")
		if err != nil {
			t.Errorf("RunCommandOutput should succeed, got: %v", err)
		}

		if !strings.Contains(output, "test output") {
			t.Errorf("Expected output to contain 'test output', got: %s", output)
		}
	})

	t.Run("command returns output with newlines", func(t *testing.T) {
		output, err := RunCommandOutput(t, "echo", tempDir, nil, "line1\nline2\nline3")
		if err != nil {
			t.Errorf("RunCommandOutput should succeed, got: %v", err)
		}

		if !strings.Contains(output, "line1") {
			t.Error("Expected output to contain 'line1'")
		}
	})

	t.Run("command with empty args", func(t *testing.T) {
		output, err := RunCommandOutput(t, "echo", tempDir, nil)
		if err != nil {
			t.Errorf("RunCommandOutput should succeed with no args, got: %v", err)
		}

		// echo with no args outputs a newline
		if output != "" && !strings.HasSuffix(output, "\n") {
			t.Errorf("Expected output to be empty or end with newline, got: %s", output)
		}
	})
}

func TestRunCLI(t *testing.T) {
	tempDir := t.TempDir()
	andurelBinary := "echo"

	t.Run("CLI command with args", func(t *testing.T) {
		env := []string{"TEST_VAR=value"}
		err := RunCLI(t, andurelBinary, tempDir, env, "test", "args")
		if err != nil {
			t.Errorf("RunCLI should succeed, got: %v", err)
		}
	})

	t.Run("CLI command without env", func(t *testing.T) {
		err := RunCLI(t, andurelBinary, tempDir, nil, "test")
		if err != nil {
			t.Errorf("RunCLI should succeed without custom env, got: %v", err)
		}
	})
}

func TestRunCommand_InvalidCommand(t *testing.T) {
	tempDir := t.TempDir()

	err := RunCommand(t, "nonexistent_command_xyz", tempDir, nil, "arg")
	if err == nil {
		t.Error("RunCommand should fail for nonexistent command")
	}
}

func TestRunCommandOutput_InvalidCommand(t *testing.T) {
	tempDir := t.TempDir()

	_, err := RunCommandOutput(t, "nonexistent_command_xyz", tempDir, nil, "arg")
	if err == nil {
		t.Error("RunCommandOutput should fail for nonexistent command")
	}
}
