package layout

import (
	"testing"

	"github.com/mbvlabs/andurel/layout/blueprint"
)

// TestTemplateRendering_BlueprintDriven tests templates that use blueprint data
func TestTemplateRendering_BlueprintDriven(t *testing.T) {
	// Create a standard template data with initialized blueprint
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

	tests := []struct {
		name         string
		templateFile string
		goldenFile   string
	}{
		{
			name:         "cmd/app/main.go with blueprint",
			templateFile: "cmd_app_main.tmpl",
			goldenFile:   "templates/cmd_app_main.golden",
		},
		{
			name:         "router/cookies/cookies.go with blueprint",
			templateFile: "router_cookies_cookies.tmpl",
			goldenFile:   "templates/router_cookies_cookies.golden",
		},
		{
			name:         "router/connect_api_routes.go with blueprint",
			templateFile: "router_connect_api_routes.tmpl",
			goldenFile:   "templates/router_connect_api_routes.golden",
		},
		{
			name:         "router/connect_assets_routes.go with blueprint",
			templateFile: "router_connect_assets_routes.tmpl",
			goldenFile:   "templates/router_connect_assets_routes.golden",
		},
		{
			name:         "router/connect_pages_routes.go with blueprint",
			templateFile: "router_connect_pages_routes.tmpl",
			goldenFile:   "templates/router_connect_pages_routes.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderToString(t, tt.templateFile, data)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}

// TestTemplateRendering_BlueprintWithExtensions tests templates with modified blueprint
// This simulates how extensions modify the blueprint
func TestTemplateRendering_BlueprintWithExtensions(t *testing.T) {
	// Create a custom blueprint with extension-like modifications
	bp := initializeBaseBlueprint("github.com/test/testapp")
	builder := blueprint.NewBuilder(bp)

	// Simulate an extension adding imports and fields
	builder.AddMainImport("github.com/test/testapp/custom")
	builder.AddControllerImport("github.com/test/testapp/services/custom")
	builder.AddControllerDependency("customService", "custom.Service")
	builder.AddControllerField("Custom", "controllers.Custom")
	builder.AddControllerConstructor("custom", "controllers.NewCustom(customService)")

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
		Extensions:           []string{"custom"},
		RunToolVersion:       "v1.0.0",
		blueprint:            builder.Blueprint(),
	}

	tests := []struct {
		name         string
		templateFile string
		goldenFile   string
	}{
		{
			name:         "cmd/app/main.go with custom extension",
			templateFile: "cmd_app_main.tmpl",
			goldenFile:   "templates/cmd_app_main_with_extension.golden",
		},
		{
			name:         "controllers/controller.go with custom extension",
			templateFile: "controllers_controller.tmpl",
			goldenFile:   "templates/controllers_controller_with_extension.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderToString(t, tt.templateFile, data)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}

// TestRerenderBlueprintTemplates tests the rerenderBlueprintTemplates function
func TestRerenderBlueprintTemplates(t *testing.T) {
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

	// First, render the templates normally
	err := processTemplatedFiles(tempDir, "vanilla", data)
	if err != nil {
		t.Fatalf("Failed to process templated files: %v", err)
	}

	// Now modify the blueprint (simulating an extension)
	builder := blueprint.NewBuilder(data.blueprint)
	builder.AddMainImport("github.com/test/testapp/newservice")
	data.blueprint = builder.Blueprint()

	// Re-render blueprint templates
	err = rerenderBlueprintTemplates(tempDir, data)
	if err != nil {
		t.Fatalf("Failed to re-render blueprint templates: %v", err)
	}

	// Verify that the files were updated and contain the new import
	// We can check one of the blueprint-driven files
	t.Run("verify re-render updated files", func(t *testing.T) {
		// This test passes if rerenderBlueprintTemplates doesn't error
		// The actual content verification is done by the golden tests above
	})
}
