package layout

import (
	"testing"
)

// TestTemplateRendering_BaseTemplates tests rendering of core base templates
func TestTemplateRendering_BaseTemplates(t *testing.T) {
	// Create a standard test template data
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
			name:         "go.mod",
			templateFile: "go_mod.tmpl",
			goldenFile:   "templates/go_mod.golden",
		},
		{
			name:         ".env.example",
			templateFile: "env.tmpl",
			goldenFile:   "templates/env.golden",
		},
		{
			name:         "README.md",
			templateFile: "readme.tmpl",
			goldenFile:   "templates/readme.golden",
		},
		{
			name:         ".gitignore",
			templateFile: "gitignore.tmpl",
			goldenFile:   "templates/gitignore.golden",
		},
		{
			name:         "config/app.go",
			templateFile: "config_app.tmpl",
			goldenFile:   "templates/config_app.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderToString(t, tt.templateFile, data)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}

// TestTemplateRendering_Tailwind tests Tailwind CSS specific templates
func TestTemplateRendering_Tailwind(t *testing.T) {
	data := &TemplateData{
		AppName:              "testapp",
		ProjectName:          "testapp",
		ModuleName:           "github.com/test/testapp",
		Database:             "postgresql",
		CSSFramework:         "tailwind",
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
			name:         "tailwind theme.css",
			templateFile: "tw_css_theme.tmpl",
			goldenFile:   "templates/tw_css_theme.golden",
		},
		{
			name:         "tailwind base.css",
			templateFile: "tw_css_base.tmpl",
			goldenFile:   "templates/tw_css_base.golden",
		},
		{
			name:         "tailwind layout.templ",
			templateFile: "tw_views_layout.tmpl",
			goldenFile:   "templates/tw_views_layout.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderToString(t, tt.templateFile, data)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}

// TestTemplateRendering_VanillaCSS tests Vanilla CSS specific templates
func TestTemplateRendering_VanillaCSS(t *testing.T) {
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
			name:         "vanilla normalize.css",
			templateFile: "assets_vanilla_css_normalize.tmpl",
			goldenFile:   "templates/vanilla_css_normalize.golden",
		},
		{
			name:         "vanilla layout.templ",
			templateFile: "vanilla_views_layout.tmpl",
			goldenFile:   "templates/vanilla_views_layout.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderToString(t, tt.templateFile, data)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}

// TestTemplateRendering_DatabaseConfig tests database-specific templates
func TestTemplateRendering_DatabaseConfig(t *testing.T) {
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
			name:         "config/config.go",
			templateFile: "config_config.tmpl",
			goldenFile:   "templates/config_config.golden",
		},
		{
			name:         "config/database.go",
			templateFile: "config_database.tmpl",
			goldenFile:   "templates/config_database.golden",
		},
		{
			name:         "database/database.go",
			templateFile: "psql_database.tmpl",
			goldenFile:   "templates/psql_database.golden",
		},
		{
			name:         "database/sqlc.yaml",
			templateFile: "psql_sqlc.tmpl",
			goldenFile:   "templates/psql_sqlc.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderToString(t, tt.templateFile, data)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}

// TestTemplateRendering_Controllers tests controller templates
func TestTemplateRendering_Controllers(t *testing.T) {
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
			name:         "controllers/controller.go",
			templateFile: "controllers_controller.tmpl",
			goldenFile:   "templates/controllers_controller.golden",
		},
		{
			name:         "controllers/pages.go",
			templateFile: "controllers_pages.tmpl",
			goldenFile:   "templates/controllers_pages.golden",
		},
		{
			name:         "controllers/sessions.go",
			templateFile: "controllers_sessions.tmpl",
			goldenFile:   "templates/controllers_sessions.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderToString(t, tt.templateFile, data)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}

// TestTemplateRendering_Auth tests authentication-related templates
func TestTemplateRendering_Auth(t *testing.T) {
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
			name:         "config/auth.go",
			templateFile: "config_auth.tmpl",
			goldenFile:   "templates/config_auth.golden",
		},
		{
			name:         "models/user.go",
			templateFile: "models_user.tmpl",
			goldenFile:   "templates/models_user.golden",
		},
		{
			name:         "services/authentication.go",
			templateFile: "services_authentication.tmpl",
			goldenFile:   "templates/services_authentication.golden",
		},
		{
			name:         "views/login.templ",
			templateFile: "views_login.tmpl",
			goldenFile:   "templates/views_login.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderToString(t, tt.templateFile, data)
			compareGolden(t, tt.goldenFile, got)
		})
	}
}
