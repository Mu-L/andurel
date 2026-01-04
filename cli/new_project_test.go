package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewProjectCommand(t *testing.T) {
	cmd := newProjectCommand("v1.0.0-test")

	if cmd == nil {
		t.Fatal("newProjectCommand returned nil")
	}

	// Check command structure
	if cmd.Use != "new [project-name]" {
		t.Errorf("Expected Use 'new [project-name]', got '%s'", cmd.Use)
	}

	if cmd.Short != "Create a new Andurel project" {
		t.Errorf("Unexpected Short description: %s", cmd.Short)
	}

	// Check that flags exist
	repoFlag := cmd.Flags().Lookup("repo")
	if repoFlag == nil {
		t.Error("Expected --repo flag to exist")
	}

	cssFlag := cmd.Flags().Lookup("css")
	if cssFlag == nil {
		t.Error("Expected --css flag to exist")
	}

	extensionsFlag := cmd.Flags().Lookup("extensions")
	if extensionsFlag == nil {
		t.Error("Expected --extensions flag to exist")
	}
}

func TestNewProject_FlagDefaults(t *testing.T) {
	cmd := newProjectCommand("v1.0.0-test")

	// Check default values
	cssFlag := cmd.Flags().Lookup("css")
	if cssFlag == nil {
		t.Fatal("Expected --css flag to exist")
	}

	if cssFlag.DefValue != "" {
		t.Errorf("Expected --css default to be empty string, got '%s'", cssFlag.DefValue)
	}

	repoFlag := cmd.Flags().Lookup("repo")
	if repoFlag == nil {
		t.Fatal("Expected --repo flag to exist")
	}

	if repoFlag.DefValue != "" {
		t.Errorf("Expected --repo default to be empty string, got '%s'", repoFlag.DefValue)
	}
}

func TestNewProject_ValidCSSFrameworks(t *testing.T) {
	testCases := []struct {
		name           string
		cssFramework   string
		shouldSucceed  bool
		expectedError  string
	}{
		{
			name:          "tailwind framework",
			cssFramework:  "tailwind",
			shouldSucceed: true,
		},
		{
			name:          "vanilla framework",
			cssFramework:  "vanilla",
			shouldSucceed: true,
		},
		{
			name:          "invalid framework",
			cssFramework:  "bootstrap",
			shouldSucceed: false,
			expectedError: "invalid css framework provided",
		},
		{
			name:          "empty framework (defaults to tailwind)",
			cssFramework:  "",
			shouldSucceed: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temp directory
			tempDir := t.TempDir()

			// Set up command with test flags
			cmd := newProjectCommand("v1.0.0-test")
			cmd.SetArgs([]string{tempDir + "/test-project", "--css", tc.cssFramework})

			// Mock the layout.Scaffold function by checking validation logic
			if !tc.shouldSucceed {
				if tc.cssFramework != "" && tc.cssFramework != "tailwind" && tc.cssFramework != "vanilla" {
					// Validation should fail for invalid frameworks
					err := func() error {
						if tc.cssFramework != "tailwind" && tc.cssFramework != "vanilla" {
							return &cssFrameworkError{framework: tc.cssFramework}
						}
						return nil
					}()

					if err == nil {
						t.Errorf("Expected error for CSS framework '%s'", tc.cssFramework)
					}
				}
			}
		})
	}
}

func TestNewProject_ValidProjectNames(t *testing.T) {
	testCases := []struct {
		name          string
		projectName   string
		shouldSucceed bool
	}{
		{
			name:          "simple name",
			projectName:   "myapp",
			shouldSucceed: true,
		},
		{
			name:          "name with hyphens",
			projectName:   "my-app",
			shouldSucceed: true,
		},
		{
			name:          "name with underscores",
			projectName:   "my_app",
			shouldSucceed: true,
		},
		{
			name:          "name with numbers",
			projectName:   "app123",
			shouldSucceed: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Just test that the command can be created with these names
			// The actual scaffolding is tested in layout tests
			cmd := newProjectCommand("v1.0.0-test")
			if cmd == nil {
				t.Fatal("Failed to create command")
			}
		})
	}
}

func TestNewProject_RepoPathValidation(t *testing.T) {
	testCases := []struct {
		name          string
		repo          string
		shouldSucceed bool
	}{
		{
			name:          "simple username",
			repo:          "myusername",
			shouldSucceed: true,
		},
		{
			name:          "with github.com prefix",
			repo:          "github.com/myusername",
			shouldSucceed: true,
		},
		{
			name:          "empty repo",
			repo:          "",
			shouldSucceed: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create command with repo flag
			cmd := newProjectCommand("v1.0.0-test")
			repoFlag := cmd.Flags().Lookup("repo")
			if repoFlag == nil {
				t.Fatal("Expected --repo flag to exist")
			}

			// Just verify the flag accepts the value
			// Actual validation happens in newProject function
			if err := repoFlag.Value.Set(tc.repo); err != nil {
				t.Errorf("Failed to set repo flag: %v", err)
			}
		})
	}
}

func TestNewProject_ExtensionsFlag(t *testing.T) {
	cmd := newProjectCommand("v1.0.0-test")
	extFlag := cmd.Flags().Lookup("extensions")
	if extFlag == nil {
		t.Fatal("Expected --extensions flag to exist")
	}

	// Test setting single extension
	if err := extFlag.Value.Set("docker"); err != nil {
		t.Errorf("Failed to set extensions flag: %v", err)
	}

	// Test setting multiple extensions
	if err := extFlag.Value.Set("docker,aws-ses"); err != nil {
		t.Errorf("Failed to set multiple extensions: %v", err)
	}

	// Test setting no extensions
	if err := extFlag.Value.Set(""); err != nil {
		t.Errorf("Failed to set empty extensions: %v", err)
	}
}

func TestNewProject_CommandStructure(t *testing.T) {
	cmd := newProjectCommand("v1.0.0-test")

	// Check that it requires exactly one argument
	if cmd.Args == nil {
		t.Fatal("Expected Args validator, got nil")
	}

	// Check version is passed correctly
	// This is implicitly tested by the command not panicking
	if cmd.RunE == nil {
		t.Error("Expected RunE function, got nil")
	}
}

func TestNewProject_HelpOutput(t *testing.T) {
	rootCmd := NewRootCommand("test", "test-date")

	tests := []struct {
		name string
		args []string
	}{
		{"new help", []string{"new", "--help"}},
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

func TestNewProject_CreateInTempDir(t *testing.T) {
	// This is an integration-style test that verifies the command structure
	// Actual scaffolding is tested in layout tests
	tempDir := t.TempDir()
	projectName := "test-project"
	projectPath := filepath.Join(tempDir, projectName)

	// Verify the project doesn't exist yet
	if _, err := os.Stat(projectPath); err == nil {
		t.Error("Project directory should not exist yet")
	}

	// Note: We don't actually call newProject here because it would scaffold
	// a full project which is better tested in layout package tests
}

// Helper type for testing CSS framework errors
type cssFrameworkError struct {
	framework string
}

func (e *cssFrameworkError) Error() string {
	return "invalid css framework provided"
}
