package templates

import (
	"strings"
	"testing"
	"text/template"
)

func TestNewTemplateService(t *testing.T) {
	service := NewTemplateService()
	if service == nil {
		t.Fatal("NewTemplateService returned nil")
	}
	if service.cache == nil {
		t.Error("service.cache is nil")
	}
	if service.functions == nil {
		t.Error("service.functions is nil")
	}
}

func TestTemplateService_RenderTemplate(t *testing.T) {
	service := NewTemplateService()

	tests := []struct {
		name        string
		templateName string
		data        any
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid template",
			templateName: "nonexistent.tmpl",
			data:        nil,
			wantErr:     true,
		},
		{
			name:        "valid template with nil data",
			templateName: "model.tmpl",
			data:        nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.RenderTemplate(tt.templateName, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errContains != "" && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("RenderTemplate() error = %v, expected to contain %v", err, tt.errContains)
			}
			if !tt.wantErr && got == "" {
				t.Error("RenderTemplate() returned empty string")
			}
		})
	}
}

func TestTemplateService_RenderTemplateWithCustomFunctions(t *testing.T) {
	service := NewTemplateService()

	customFuncs := template.FuncMap{
		"customFunc": func(s string) string {
			return "custom:" + s
		},
	}

	got, err := service.RenderTemplateWithCustomFunctions("model.tmpl", nil, customFuncs)
	if err == nil {
		t.Error("Expected error with nil data, got nil")
	}
	if got != "" {
		t.Error("Expected empty string with error, got non-empty")
	}
}

func TestNewTemplateBuilder(t *testing.T) {
	service := NewTemplateService()
	builder := NewTemplateBuilder(service)

	if builder == nil {
		t.Fatal("NewTemplateBuilder returned nil")
	}
	if builder.service != service {
		t.Error("builder.service does not match provided service")
	}
	if builder.data == nil {
		t.Error("builder.data is nil")
	}
	if builder.data.Custom == nil {
		t.Error("builder.data.Custom is nil")
	}
}

func TestTemplateBuilder_WithResource(t *testing.T) {
	service := NewTemplateService()
	builder := NewTemplateBuilder(service)

	result := builder.WithResource("User", "Users", "github.com/test/app", "model", "postgresql", []string{"field1", "field2"})
	if result == nil {
		t.Fatal("WithResource returned nil")
	}
	if result.data.Resource.Name != "User" {
		t.Errorf("Resource.Name = %v, want User", result.data.Resource.Name)
	}
	if result.data.Resource.PluralName != "Users" {
		t.Errorf("Resource.PluralName = %v, want Users", result.data.Resource.PluralName)
	}
	if result.data.Resource.ModulePath != "github.com/test/app" {
		t.Errorf("Resource.ModulePath = %v, want github.com/test/app", result.data.Resource.ModulePath)
	}
}

func TestTemplateBuilder_WithDatabase(t *testing.T) {
	service := NewTemplateService()
	builder := NewTemplateBuilder(service)

	result := builder.WithDatabase("postgresql", "DB", "pgx")
	if result == nil {
		t.Fatal("WithDatabase returned nil")
	}
	if result.data.Database.Type != "postgresql" {
		t.Errorf("Database.Type = %v, want postgresql", result.data.Database.Type)
	}
	if result.data.Database.Method != "DB" {
		t.Errorf("Database.Method = %v, want DB", result.data.Database.Method)
	}
	if result.data.Database.Driver != "pgx" {
		t.Errorf("Database.Driver = %v, want pgx", result.data.Database.Driver)
	}
}

func TestTemplateBuilder_WithProject(t *testing.T) {
	service := NewTemplateService()
	builder := NewTemplateBuilder(service)

	result := builder.WithProject("github.com/test/app", "testapp")
	if result == nil {
		t.Fatal("WithProject returned nil")
	}
	if result.data.Project.ModulePath != "github.com/test/app" {
		t.Errorf("Project.ModulePath = %v, want github.com/test/app", result.data.Project.ModulePath)
	}
	if result.data.Project.Name != "testapp" {
		t.Errorf("Project.Name = %v, want testapp", result.data.Project.Name)
	}
}

func TestTemplateBuilder_WithCustom(t *testing.T) {
	service := NewTemplateService()
	builder := NewTemplateBuilder(service)

	result := builder.WithCustom("key1", "value1")
	if result == nil {
		t.Fatal("WithCustom returned nil")
	}
	if result.data.Custom["key1"] != "value1" {
		t.Errorf("Custom['key1'] = %v, want value1", result.data.Custom["key1"])
	}
}

func TestTemplateBuilder_Chaining(t *testing.T) {
	service := NewTemplateService()
	builder := NewTemplateBuilder(service)

	result := builder.
		WithResource("User", "Users", "github.com/test/app", "model", "postgresql", nil).
		WithDatabase("postgresql", "DB", "pgx").
		WithProject("github.com/test/app", "testapp").
		WithCustom("key1", "value1")

	if result == nil {
		t.Fatal("Method chaining returned nil")
	}
	if result.data.Resource.Name != "User" {
		t.Errorf("Resource.Name = %v, want User", result.data.Resource.Name)
	}
	if result.data.Database.Type != "postgresql" {
		t.Errorf("Database.Type = %v, want postgresql", result.data.Database.Type)
	}
	if result.data.Project.Name != "testapp" {
		t.Errorf("Project.Name = %v, want testapp", result.data.Project.Name)
	}
	if result.data.Custom["key1"] != "value1" {
		t.Errorf("Custom['key1'] = %v, want value1", result.data.Custom["key1"])
	}
}

func TestTemplateData_Structures(t *testing.T) {
	data := &TemplateData{
		Resource: ResourceData{
			Name:         "User",
			PluralName:   "Users",
			ModulePath:   "github.com/test/app",
			Type:         "model",
			DatabaseType: "postgresql",
		},
		Database: DatabaseData{
			Type:   "postgresql",
			Method: "DB",
			Driver: "pgx",
		},
		Project: ProjectData{
			ModulePath: "github.com/test/app",
			Name:       "testapp",
		},
		Custom: map[string]any{
			"key1": "value1",
		},
	}

	if data.Resource.Name != "User" {
		t.Errorf("Resource.Name = %v, want User", data.Resource.Name)
	}
	if data.Database.Type != "postgresql" {
		t.Errorf("Database.Type = %v, want postgresql", data.Database.Type)
	}
	if data.Project.Name != "testapp" {
		t.Errorf("Project.Name = %v, want testapp", data.Project.Name)
	}
	if data.Custom["key1"] != "value1" {
		t.Errorf("Custom['key1'] = %v, want value1", data.Custom["key1"])
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"hello_world", "helloWorld"},
		{"user_id", "userId"},
		{"simple", "simple"},
		{"multi_part_string", "multiPartString"},
		{"", ""},
		{"_leading", "Leading"},
		{"trailing_", "trailing"},
		{"a_b_c", "aBC"},
		{"HTTP_request", "httpRequest"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toCamelCase(tt.input)
			if got != tt.want {
				t.Errorf("toCamelCase(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestToLowerCamelCase(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"Hello", "hello"},
		{"HelloWorld", "helloWorld"},
		{"ABC", "aBC"},
		{"a", "a"},
		{"", ""},
		{"HTTPResponse", "hTTPResponse"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := toLowerCamelCase(tt.input)
			if got != tt.want {
				t.Errorf("toLowerCamelCase(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetDefaultTemplateFunctions(t *testing.T) {
	funcs := getDefaultTemplateFunctions()

	if funcs == nil {
		t.Fatal("getDefaultTemplateFunctions returned nil")
	}

	requiredFuncs := []string{
		"ToLower",
		"ToUpper",
		"ToSnakeCase",
		"ToCamelCase",
		"ToLowerCamelCase",
		"DeriveTableName",
		"DatabaseType",
		"DatabaseMethod",
		"uuidParam",
		"toLowerCamelCase",
		"toCamelCase",
	}

	for _, name := range requiredFuncs {
		if _, ok := funcs[name]; !ok {
			t.Errorf("Function %s not found in func map", name)
		}
	}
}

func TestGetGlobalTemplateService(t *testing.T) {
	service := GetGlobalTemplateService()

	if service == nil {
		t.Fatal("GetGlobalTemplateService returned nil")
	}
	if service.cache == nil {
		t.Error("service.cache is nil")
	}
}

func TestRenderTemplateUsingGlobal(t *testing.T) {
	got, err := RenderTemplateUsingGlobal("model.tmpl", nil)
	if err == nil {
		t.Error("Expected error with nil data, got nil")
	}
	if got != "" {
		t.Error("Expected empty string with error, got non-empty")
	}
}

func TestNewTemplateBuilderUsingGlobal(t *testing.T) {
	builder := NewTemplateBuilderUsingGlobal()

	if builder == nil {
		t.Fatal("NewTemplateBuilderUsingGlobal returned nil")
	}
	if builder.service == nil {
		t.Error("builder.service is nil")
	}
}
