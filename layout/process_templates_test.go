package layout

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbvlabs/andurel/layout/templates"
)

// TestProcessTemplatedFiles tests the processTemplatedFiles function
func TestProcessTemplatedFiles(t *testing.T) {
	tests := []struct {
		name         string
		cssFramework string
		expectFiles  []string
	}{
		{
			name:         "vanilla CSS framework",
			cssFramework: "vanilla",
			expectFiles: []string{
				".env.example",
				".gitignore",
				"README.md",
				"assets/css/normalize.css",
				"views/layout.templ",
			},
		},
		{
			name:         "tailwind CSS framework",
			cssFramework: "tailwind",
			expectFiles: []string{
				".env.example",
				".gitignore",
				"README.md",
				"css/theme.css",
				"css/base.css",
				"views/layout.templ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			data := &TemplateData{
				AppName:              "testapp",
				ProjectName:          "testapp",
				ModuleName:           "github.com/test/testapp",
				Database:             "postgresql",
				CSSFramework:         tt.cssFramework,
				GoVersion:            "1.25.0",
				SessionKey:           "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
				SessionEncryptionKey: "0123456789abcdef0123456789abcdef",
				TokenSigningKey:      "0123456789abcdef0123456789abcdef",
				Pepper:               "0123456789abcdef01234567",
				Extensions:           []string{},
				RunToolVersion:       "v1.0.0",
				blueprint:            initializeBaseBlueprint("github.com/test/testapp"),
			}

			err := processTemplatedFiles(tempDir, tt.cssFramework, data)
			if err != nil {
				t.Fatalf("processTemplatedFiles() error = %v", err)
			}

			// Verify expected files exist
			for _, file := range tt.expectFiles {
				fullPath := filepath.Join(tempDir, file)
				if _, err := os.Stat(fullPath); os.IsNotExist(err) {
					t.Errorf("Expected file does not exist: %s", file)
				}
			}
		})
	}
}

// TestProcessTemplatedFiles_InvalidData tests error handling
func TestProcessTemplatedFiles_InvalidData(t *testing.T) {
	tempDir := t.TempDir()

	// This should work even with minimal data
	data := &TemplateData{
		AppName:    "testapp",
		ModuleName: "testapp",
		blueprint:  initializeBaseBlueprint("testapp"),
	}

	err := processTemplatedFiles(tempDir, "vanilla", data)
	if err != nil {
		t.Errorf("processTemplatedFiles() should handle minimal data, got error: %v", err)
	}
}

// TestRerenderBlueprintTemplates_Integration tests the full re-rendering flow
func TestRerenderBlueprintTemplates_Integration(t *testing.T) {
	tempDir := t.TempDir()

	data := &TemplateData{
		AppName:              "testapp",
		ProjectName:          "testapp",
		ModuleName:           "github.com/test/testapp",
		Database:             "postgresql",
		CSSFramework:         "vanilla",
		GoVersion:            "1.25.0",
		SessionKey:           "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		SessionEncryptionKey: "0123456789abcdef0123456789abcdef",
		TokenSigningKey:      "0123456789abcdef0123456789abcdef",
		Pepper:               "0123456789abcdef01234567",
		Extensions:           []string{},
		RunToolVersion:       "v1.0.0",
		blueprint:            initializeBaseBlueprint("github.com/test/testapp"),
	}

	// First render
	err := processTemplatedFiles(tempDir, "vanilla", data)
	if err != nil {
		t.Fatalf("processTemplatedFiles() error = %v", err)
	}

	// Get initial file modification time
	mainGoPath := filepath.Join(tempDir, "cmd", "app", "main.go")
	initialInfo, err := os.Stat(mainGoPath)
	if err != nil {
		t.Fatalf("Failed to stat initial main.go: %v", err)
	}

	// Re-render
	err = rerenderBlueprintTemplates(tempDir, data)
	if err != nil {
		t.Fatalf("rerenderBlueprintTemplates() error = %v", err)
	}

	// Verify file was updated
	updatedInfo, err := os.Stat(mainGoPath)
	if err != nil {
		t.Fatalf("Failed to stat updated main.go: %v", err)
	}

	// File should exist and have been modified
	if updatedInfo.ModTime().Before(initialInfo.ModTime()) {
		t.Error("File modification time should be updated after re-render")
	}
}

// TestRerenderBlueprintTemplates_NilData tests error handling with nil data
func TestRerenderBlueprintTemplates_NilData(t *testing.T) {
	tempDir := t.TempDir()

	err := rerenderBlueprintTemplates(tempDir, nil)
	if err == nil {
		t.Error("rerenderBlueprintTemplates() should return error for nil data")
	}
	if err != nil && err.Error() != "template data is nil" {
		t.Errorf("Expected 'template data is nil' error, got: %v", err)
	}
}

// TestRenderTemplate_ErrorPaths tests error handling in renderTemplate
func TestRenderTemplate_ErrorPaths(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("invalid template file", func(t *testing.T) {
		data := &TemplateData{
			AppName: "testapp",
		}

		err := renderTemplate(tempDir, "nonexistent.tmpl", "output.txt", templates.Files, data)
		if err == nil {
			t.Error("renderTemplate() should return error for nonexistent template")
		}
	})

	t.Run("nil data", func(t *testing.T) {
		// renderTemplate should handle nil data by creating an empty TemplateData
		err := renderTemplate(tempDir, "env.tmpl", "test.env", templates.Files, nil)
		// This might succeed or fail depending on the template requirements
		// We're just checking it doesn't panic
		_ = err
	})
}

// TestCreateGoMod tests the createGoMod function
func TestCreateGoMod(t *testing.T) {
	tempDir := t.TempDir()

	data := &TemplateData{
		ModuleName: "github.com/test/testapp",
		GoVersion:  "1.25.0",
		blueprint:  initializeBaseBlueprint("github.com/test/testapp"),
	}

	err := createGoMod(tempDir, data)
	if err != nil {
		t.Fatalf("createGoMod() error = %v", err)
	}

	// Verify go.mod was created
	goModPath := filepath.Join(tempDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Fatalf("Failed to read go.mod: %v", err)
	}

	// Verify it contains the module name
	if !contains(string(content), "github.com/test/testapp") {
		t.Error("go.mod should contain the module name")
	}
}

// TestCreateGoMod_NilData tests error handling with nil data
func TestCreateGoMod_NilData(t *testing.T) {
	tempDir := t.TempDir()

	err := createGoMod(tempDir, nil)
	if err == nil {
		t.Error("createGoMod() should return error for nil data")
	}
}

// TestCopyFile_Integration tests the copyFile function
func TestCopyFile_Integration(t *testing.T) {
	tempDir := t.TempDir()

	// Test copying a real file from the templates
	err := copyFile(tempDir, "assets_js_datastar.tmpl", "assets/js/datastar.min.js", templates.Files)
	if err != nil {
		t.Fatalf("copyFile() error = %v", err)
	}

	// Verify the file was copied
	targetPath := filepath.Join(tempDir, "assets", "js", "datastar.min.js")
	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		t.Error("copyFile() should create the target file")
	}
}
