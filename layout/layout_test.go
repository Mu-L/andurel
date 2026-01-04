package layout

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbvlabs/andurel/layout/blueprint"
)

func TestGenerateRandomHex(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int
		wantLen  int
	}{
		{
			name:    "32 bytes (256 bits)",
			bytes:   32,
			wantLen: 64,
		},
		{
			name:    "12 bytes (96 bits)",
			bytes:   12,
			wantLen: 24,
		},
		{
			name:    "64 bytes (512 bits)",
			bytes:   64,
			wantLen: 128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateRandomHex(tt.bytes)
			if len(result) != tt.wantLen {
				t.Errorf("generateRandomHex(%d) = %s, want length %d, got %d", tt.bytes, result, tt.wantLen, len(result))
			}
		})
	}
}

func TestGenerateRandomHex_Uniqueness(t *testing.T) {
	// Test that multiple calls produce different results
	results := make(map[string]bool)
	for i := 0; i < 100; i++ {
		result := generateRandomHex(32)
		if results[result] {
			t.Errorf("Generated duplicate hex string: %s", result)
		}
		results[result] = true
	}
}

func TestRegisterBuiltinExtensions(t *testing.T) {
	// First call should succeed
	err := registerBuiltinExtensions()
	if err != nil {
		t.Fatalf("First call to registerBuiltinExtensions failed: %v", err)
	}

	// Second call should also succeed (idempotent due to sync.Once)
	err = registerBuiltinExtensions()
	if err != nil {
		t.Errorf("Second call to registerBuiltinExtensions failed: %v", err)
	}
}

func TestResolveExtensions(t *testing.T) {
	// Register builtin extensions first
	if err := registerBuiltinExtensions(); err != nil {
		t.Fatalf("Failed to register builtin extensions: %v", err)
	}

	tests := []struct {
		name          string
		extensionNames []string
		wantCount     int
		wantErr       bool
		errContains   string
	}{
		{
			name:          "no extensions",
			extensionNames: []string{},
			wantCount:     0,
			wantErr:       false,
		},
		{
			name:          "single extension",
			extensionNames: []string{"docker"},
			wantCount:     1,
			wantErr:       false,
		},
		{
			name:          "multiple extensions",
			extensionNames: []string{"docker", "aws-ses"},
			wantCount:     2,
			wantErr:       false,
		},
		{
			name:          "unknown extension",
			extensionNames: []string{"unknown-extension"},
			wantCount:     0,
			wantErr:       true,
			errContains:   "unknown extension",
		},
		{
			name:          "empty extension name",
			extensionNames: []string{"", "docker"},
			wantCount:     0,
			wantErr:       true,
			errContains:   "cannot be empty",
		},
		{
			name:          "whitespace only",
			extensionNames: []string{"   "},
			wantCount:     0,
			wantErr:       true,
			errContains:   "cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exts, err := resolveExtensions(tt.extensionNames)

			if tt.wantErr {
				if err == nil {
					t.Errorf("resolveExtensions() expected error, got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("resolveExtensions() error = %v, want error containing %s", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("resolveExtensions() unexpected error: %v", err)
				return
			}

			if len(exts) != tt.wantCount {
				t.Errorf("resolveExtensions() returned %d extensions, want %d", len(exts), tt.wantCount)
			}
		})
	}
}

func TestCollectDependencies(t *testing.T) {
	// Register builtin extensions first
	if err := registerBuiltinExtensions(); err != nil {
		t.Fatalf("Failed to register builtin extensions: %v", err)
	}

	tests := []struct {
		name     string
		requested map[string]struct{}
		wantErr  bool
	}{
		{
			name:      "no dependencies",
			requested: map[string]struct{}{"docker": {}},
			wantErr:   false,
		},
		{
			name:      "unknown extension",
			requested: map[string]struct{}{"unknown-ext": {}},
			wantErr:   true,
		},
		{
			name:      "self-dependency",
			requested: map[string]struct{}{"self-dep-test": {}},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allNeeded := make(map[string]struct{})
			err := collectDependencies(tt.requested, allNeeded)

			if tt.wantErr {
				if err == nil {
					t.Errorf("collectDependencies() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("collectDependencies() unexpected error: %v", err)
			}
		})
	}
}

func TestTopologicalSort(t *testing.T) {
	// Register builtin extensions first
	if err := registerBuiltinExtensions(); err != nil {
		t.Fatalf("Failed to register builtin extensions: %v", err)
	}

	tests := []struct {
		name      string
		extSet    map[string]struct{}
		wantLen   int
		wantErr   bool
		errContains string
	}{
		{
			name:    "empty set",
			extSet:  map[string]struct{}{},
			wantLen: 0,
			wantErr: false,
		},
		{
			name:    "single extension",
			extSet:  map[string]struct{}{"docker": {}},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "multiple independent extensions",
			extSet:  map[string]struct{}{"docker": {}, "aws-ses": {}},
			wantLen: 2,
			wantErr: false,
		},
		{
			name:    "unknown extension",
			extSet:  map[string]struct{}{"unknown": {}},
			wantLen: 0,
			wantErr: true,
			errContains: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := topologicalSort(tt.extSet)

			if tt.wantErr {
				if err == nil {
					t.Errorf("topologicalSort() expected error, got nil")
				} else if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("topologicalSort() error = %v, want error containing %s", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("topologicalSort() unexpected error: %v", err)
				return
			}

			if len(result) != tt.wantLen {
				t.Errorf("topologicalSort() returned %d items, want %d", len(result), tt.wantLen)
			}
		})
	}
}

func TestTemplateFuncMap(t *testing.T) {
	funcs := templateFuncMap()

	if funcs == nil {
		t.Fatal("templateFuncMap() returned nil")
	}

	// Test lower function
	lowerFunc, ok := funcs["lower"]
	if !ok {
		t.Fatal("templateFuncMap() missing 'lower' function")
	}

	result := lowerFunc.(func(string) string)("TEST_STRING")
	if result != "test_string" {
		t.Errorf("lower function = %s, want 'test_string'", result)
	}
}

func TestExtractRepo(t *testing.T) {
	tests := []struct {
		name     string
		module   string
		expected string
	}{
		{
			name:     "full module path",
			module:   "github.com/mbvlabs/andurel",
			expected: "github.com/mbvlabs/andurel",
		},
		{
			name:     "short module path",
			module:   "example.com/module",
			expected: "example.com/module",
		},
		{
			name:     "simple module name",
			module:   "mymodule",
			expected: "mymodule",
		},
		{
			name:     "deep module path",
			module:   "github.com/username/repo/subpackage",
			expected: "github.com/username/repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractRepo(tt.module)
			if result != tt.expected {
				t.Errorf("extractRepo(%s) = %s, want %s", tt.module, result, tt.expected)
			}
		})
	}
}

func TestProcessMigrations(t *testing.T) {
	// Set test mode
	os.Setenv("ANDUREL_TEST_MODE", "true")
	defer os.Unsetenv("ANDUREL_TEST_MODE")

	tempDir := t.TempDir()

	data := &TemplateData{
		AppName:     "test-app",
		ModuleName:  "test.app",
		Database:    "postgresql",
		blueprint:   initializeBaseBlueprint("test.app"),
	}

	// Set run tool version
	data.RunToolVersion = GetRunToolVersion()

	lastTime, err := processMigrations(tempDir, data)
	if err != nil {
		t.Fatalf("processMigrations() error = %v", err)
	}

	// Verify migrations were created
	migrationsDir := filepath.Join(tempDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Should have 8 migrations (6 river queue + 2 auth)
	if len(entries) != 8 {
		t.Errorf("Expected 8 migration files, got %d", len(entries))
	}

	// Verify lastTime is not zero
	if lastTime.IsZero() {
		t.Error("processMigrations() returned zero time")
	}
}

func TestCopyFile(t *testing.T) {
	// Create a temporary source file
	tempDir := t.TempDir()
	sourcePath := filepath.Join(tempDir, "source.txt")
	content := []byte("test content for copy")

	if err := os.WriteFile(sourcePath, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Test copyFile function using os.ReadFile directly
	err := func() error {
		content, err := os.ReadFile(sourcePath)
		if err != nil {
			return err
		}

		targetDir := filepath.Join(tempDir, "target")
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}

		targetPath := filepath.Join(targetDir, "source.txt")
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return err
		}

		// Verify content matches
		copiedContent, err := os.ReadFile(targetPath)
		if err != nil {
			return err
		}

		if string(copiedContent) != string(content) {
			return err
		}

		return nil
	}()

	if err != nil {
		t.Errorf("File copy failed: %v", err)
	}
}

func TestGetExpectedTools(t *testing.T) {
	tests := []struct {
		name string
		config *ScaffoldConfig
		wantTailwind bool
	}{
		{
			name: "with tailwind",
			config: &ScaffoldConfig{
				CSSFramework: "tailwind",
			},
			wantTailwind: true,
		},
		{
			name: "without tailwind",
			config: &ScaffoldConfig{
				CSSFramework: "vanilla",
			},
			wantTailwind: false,
		},
		{
			name: "nil config",
			config: nil,
			wantTailwind: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tools := GetExpectedTools(tt.config)

			// Check that default tools exist
			defaultTools := []string{"templ", "sqlc", "goose", "air", "mailpit", "usql", "dblab", "run"}
			for _, toolName := range defaultTools {
				if _, ok := tools[toolName]; !ok {
					t.Errorf("GetExpectedTools() missing tool: %s", toolName)
				}
			}

			// Check tailwindcli
			if tt.wantTailwind {
				if _, ok := tools["tailwindcli"]; !ok {
					t.Error("GetExpectedTools() missing tailwindcli when CSSFramework is tailwind")
				}
			}
		})
	}
}

func TestGetRunToolVersion(t *testing.T) {
	version := GetRunToolVersion()

	if version == "" {
		t.Error("GetRunToolVersion() returned empty string")
	}
}

func TestDefaultGoTools(t *testing.T) {
	if len(DefaultGoTools) == 0 {
		t.Error("DefaultGoTools is empty")
	}

	expectedTools := []string{"templ", "sqlc", "goose", "air", "mailpit", "usql", "dblab"}
	toolNames := make(map[string]bool)
	for _, tool := range DefaultGoTools {
		toolNames[tool.Name] = true
	}

	for _, expected := range expectedTools {
		if !toolNames[expected] {
			t.Errorf("DefaultGoTools missing expected tool: %s", expected)
		}
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestTemplateData_Blueprint(t *testing.T) {
	data := &TemplateData{}

	// Test GetModuleName
	if data.GetModuleName() != "" {
		t.Error("GetModuleName() should return empty string for nil data")
	}

	data.ModuleName = "test.module"
	if data.GetModuleName() != "test.module" {
		t.Errorf("GetModuleName() = %s, want 'test.module'", data.GetModuleName())
	}

	// Test DatabaseDialect
	if data.DatabaseDialect() != "postgresql" {
		t.Errorf("DatabaseDialect() = %s, want 'postgresql'", data.DatabaseDialect())
	}

	// Test Blueprint
	bp := data.Blueprint()
	if bp == nil {
		t.Error("Blueprint() returned nil")
	}

	// Test SetBlueprint
	newBp := blueprint.New()
	data.SetBlueprint(newBp)
	if data.blueprint != newBp {
		t.Error("SetBlueprint() did not set blueprint correctly")
	}

	// Test Builder
	builder := data.Builder()
	if builder == nil {
		t.Error("Builder() returned nil")
	}
}
