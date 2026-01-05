package layout

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mbvlabs/andurel/layout/templates"
)

var updateGolden = flag.Bool("update", false, "update golden files")

// goldenFilePath returns the path to a golden file
func goldenFilePath(name string) string {
	return filepath.Join("testdata", "golden", name)
}

// compareGolden compares content against a golden file, updating it if -update flag is set
func compareGolden(t *testing.T, name string, got string) {
	t.Helper()

	goldenPath := goldenFilePath(name)

	if *updateGolden {
		// Create directory if it doesn't exist
		dir := filepath.Dir(goldenPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create golden file directory: %v", err)
		}

		// Write the golden file
		if err := os.WriteFile(goldenPath, []byte(got), 0644); err != nil {
			t.Fatalf("Failed to write golden file: %v", err)
		}
		t.Logf("Updated golden file: %s", goldenPath)
		return
	}

	// Read the golden file
	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("Failed to read golden file %s: %v\nRun with -update flag to create it", goldenPath, err)
	}

	// Normalize line endings for comparison
	gotNormalized := normalizeLineEndings(got)
	wantNormalized := normalizeLineEndings(string(want))

	if gotNormalized != wantNormalized {
		t.Errorf("Content does not match golden file: %s\n\nGot:\n%s\n\nWant:\n%s\n\nRun with -update flag to update golden file",
			goldenPath, got, string(want))
	}
}

// normalizeLineEndings converts all line endings to \n
func normalizeLineEndings(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return s
}

// renderToString renders a template and returns the result as a string
func renderToString(t *testing.T, templateFile string, data *TemplateData) string {
	t.Helper()

	tempDir := t.TempDir()
	targetPath := filepath.Base(templateFile)
	targetPath = strings.TrimSuffix(targetPath, ".tmpl")

	err := renderTemplate(tempDir, templateFile, targetPath, templates.Files, data)
	if err != nil {
		t.Fatalf("Failed to render template %s: %v", templateFile, err)
	}

	content, err := os.ReadFile(filepath.Join(tempDir, targetPath))
	if err != nil {
		t.Fatalf("Failed to read rendered file: %v", err)
	}

	return string(content)
}
