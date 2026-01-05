package upgrade

import (
	"strings"
	"testing"

	"github.com/mbvlabs/andurel/layout"
	"github.com/mbvlabs/andurel/layout/templates"
)

func TestNewTemplateGenerator(t *testing.T) {
	gen := NewTemplateGenerator("v1.0.0")
	if gen == nil {
		t.Fatal("expected non-nil generator")
	}
	if gen.targetVersion != "v1.0.0" {
		t.Errorf("expected targetVersion to be v1.0.0, got %s", gen.targetVersion)
	}
}

func TestGetFrameworkTemplates(t *testing.T) {
	templates := GetFrameworkTemplates()

	if len(templates) == 0 {
		t.Fatal("expected non-empty template list")
	}

	// Verify expected templates exist
	expectedTemplates := map[string]string{
		"framework_elements_renderer_render.tmpl":  "internal/renderer/render.go",
		"framework_elements_routing_definitions.tmpl": "internal/routing/definitions.go",
		"framework_elements_server_server.tmpl":    "internal/server/server.go",
		"cmd_run_main.tmpl":                        "cmd/run/main.go",
	}

	templateMap := make(map[string]string)
	for _, ft := range templates {
		templateMap[ft.TemplateName] = ft.TargetPath
	}

	for name, expectedPath := range expectedTemplates {
		if path, ok := templateMap[name]; !ok {
			t.Errorf("expected template %s not found", name)
		} else if path != expectedPath {
			t.Errorf("template %s has path %s, expected %s", name, path, expectedPath)
		}
	}
}

func TestBuildTemplateData(t *testing.T) {
	gen := NewTemplateGenerator("v1.0.0")

	tests := []struct {
		name   string
		config layout.ScaffoldConfig
		checkFn func(*testing.T, *layout.TemplateData)
	}{
		{
			name: "basic config",
			config: layout.ScaffoldConfig{
				ProjectName:  "myapp",
				Repository:   "",
				Database:     "postgresql",
				CSSFramework: "tailwind",
			},
			checkFn: func(t *testing.T, data *layout.TemplateData) {
				if data.AppName != "myapp" {
					t.Errorf("expected AppName to be myapp, got %s", data.AppName)
				}
				if data.ModuleName != "myapp" {
					t.Errorf("expected ModuleName to be myapp, got %s", data.ModuleName)
				}
				if data.Database != "postgresql" {
					t.Errorf("expected Database to be postgresql, got %s", data.Database)
				}
			},
		},
		{
			name: "config with repository",
			config: layout.ScaffoldConfig{
				ProjectName:  "myapp",
				Repository:   "github.com/user",
				Database:     "postgresql",
				CSSFramework: "vanilla",
			},
			checkFn: func(t *testing.T, data *layout.TemplateData) {
				if data.ModuleName != "github.com/user/myapp" {
					t.Errorf("expected ModuleName to be github.com/user/myapp, got %s", data.ModuleName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := gen.buildTemplateData(tt.config)
			tt.checkFn(t, data)
		})
	}
}

func TestRenderFrameworkTemplates(t *testing.T) {
	gen := NewTemplateGenerator("v1.0.0")

	config := layout.ScaffoldConfig{
		ProjectName:  "testapp",
		Repository:   "github.com/test",
		Database:     "postgresql",
		CSSFramework: "tailwind",
		Extensions:   []string{},
	}

	result, err := gen.RenderFrameworkTemplates(config)
	if err != nil {
		t.Fatalf("RenderFrameworkTemplates failed: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("expected non-empty result")
	}

	// Verify some expected files are present
	expectedFiles := []string{
		"internal/renderer/render.go",
		"internal/routing/definitions.go",
		"internal/server/server.go",
		"cmd/run/main.go",
	}

	for _, file := range expectedFiles {
		if content, ok := result[file]; !ok {
			t.Errorf("expected file %s not found in result", file)
		} else if len(content) == 0 {
			t.Errorf("file %s has empty content", file)
		}
	}
}

func TestRenderTemplateToBytes(t *testing.T) {
	// Create simple test data
	data := &layout.TemplateData{
		AppName:     "testapp",
		ProjectName: "testapp",
		ModuleName:  "github.com/test/testapp",
		Database:    "postgresql",
	}

	// Render a real template
	content, err := renderTemplateToBytes("cmd_run_main.tmpl", templates.Files, data)
	if err != nil {
		t.Fatalf("renderTemplateToBytes failed: %v", err)
	}

	if len(content) == 0 {
		t.Error("expected non-empty content")
	}

	// Verify the content contains expected strings
	contentStr := string(content)
	if !strings.Contains(contentStr, "package main") {
		t.Error("expected content to contain 'package main'")
	}
}

func TestRenderTemplateToBytes_InvalidTemplate(t *testing.T) {
	data := &layout.TemplateData{
		AppName: "testapp",
	}

	_, err := renderTemplateToBytes("nonexistent.tmpl", templates.Files, data)
	if err == nil {
		t.Error("expected error for nonexistent template")
	}
}
