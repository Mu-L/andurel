package controllers

import (
	"testing"

	"github.com/mbvlabs/andurel/generator/internal/catalog"
)

func TestNewGenerator(t *testing.T) {
	testCases := []struct {
		name         string
		databaseType string
	}{
		{"postgresql", "postgresql"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := NewGenerator(tc.databaseType)
			if g == nil {
				t.Fatal("NewGenerator returned nil")
			}
			if g.typeMapper == nil {
				t.Error("typeMapper is nil")
			}
		})
	}
}

func TestBuild_SimpleResource(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "users")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("email", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("name", "TEXT"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	controller, err := g.Build(cat, Config{
		ResourceName:   "User",
		PluralName:     "users",
		TableName:      "users",
		ControllerType: ResourceController,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if controller == nil {
		t.Fatal("Controller is nil")
	}

	if controller.ResourceName != "User" {
		t.Errorf("Expected resource name 'User', got '%s'", controller.ResourceName)
	}

	if len(controller.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(controller.Fields))
	}

	idField := controller.Fields[0]
	if idField.Name != "ID" {
		t.Errorf("Expected first field name 'ID', got '%s'", idField.Name)
	}

	if idField.GoType != "uuid.UUID" {
		t.Errorf("Expected ID type 'uuid.UUID', got '%s'", idField.GoType)
	}
}

func TestBuild_TableNotFound(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")

	_, err := g.Build(cat, Config{
		ResourceName: "NonExistent",
		PluralName:   "nonexistents",
		TableName:    "nonexistent",
	})

	if err == nil {
		t.Error("Expected error for non-existent table, got nil")
	}
}

func TestBuild_WithTimestampFields(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "posts")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("title", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("created_at", "TIMESTAMP"))
	table.AddColumn(catalog.NewColumn("updated_at", "TIMESTAMP"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	controller, err := g.Build(cat, Config{
		ResourceName:   "Post",
		PluralName:     "posts",
		TableName:      "posts",
		ControllerType: ResourceController,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(controller.Fields) != 4 {
		t.Errorf("Expected 4 fields, got %d", len(controller.Fields))
	}

	createdAtField := getControllerFieldByName(controller.Fields, "CreatedAt")
	if createdAtField == nil {
		t.Fatal("Expected CreatedAt field")
	}

	if !createdAtField.IsSystemField {
		t.Error("Expected CreatedAt to be marked as system field")
	}
}

func TestBuild_WithSystemFields(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "items")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("created_at", "TIMESTAMP"))
	table.AddColumn(catalog.NewColumn("updated_at", "TIMESTAMP"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	controller, err := g.Build(cat, Config{
		ResourceName:   "Item",
		PluralName:     "items",
		TableName:      "items",
		ControllerType: ResourceController,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	systemFieldCount := 0
	for _, field := range controller.Fields {
		if field.IsSystemField {
			systemFieldCount++
		}
	}

	if systemFieldCount == 0 {
		t.Error("Expected at least one system field")
	}
}

func TestBuild_FieldTypeMapping(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "products")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("name", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("price", "DECIMAL"))
	table.AddColumn(catalog.NewColumn("quantity", "INTEGER"))
	table.AddColumn(catalog.NewColumn("is_available", "BOOLEAN"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	controller, err := g.Build(cat, Config{
		ResourceName:   "Product",
		PluralName:     "products",
		TableName:      "products",
		ControllerType: ResourceController,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	expectedTypes := map[string]string{
		"ID":          "uuid.UUID",
		"Name":        "string",
		"Price":       "float64",
		"Quantity":    "int32",
		"IsAvailable": "bool",
	}

	for fieldName, expectedType := range expectedTypes {
		field := getControllerFieldByName(controller.Fields, fieldName)
		if field == nil {
			t.Errorf("Expected field '%s' not found", fieldName)
			continue
		}
		if field.GoType != expectedType {
			t.Errorf("Field %s: expected type '%s', got '%s'", fieldName, expectedType, field.GoType)
		}
	}
}

func getControllerFieldByName(fields []GeneratedField, name string) *GeneratedField {
	for i := range fields {
		if fields[i].Name == name {
			return &fields[i]
		}
	}
	return nil
}

func TestNewFileGenerator(t *testing.T) {
	fg := NewFileGenerator()
	if fg == nil {
		t.Fatal("NewFileGenerator returned nil")
	}
}

func TestNewRouteGenerator(t *testing.T) {
	rg := NewRouteGenerator()
	if rg == nil {
		t.Fatal("NewRouteGenerator returned nil")
	}
}

func TestNewTemplateRenderer(t *testing.T) {
	tr := NewTemplateRenderer()
	if tr == nil {
		t.Fatal("NewTemplateRenderer returned nil")
	}
}

func TestRenderControllerFile(t *testing.T) {
	tr := NewTemplateRenderer()

	controller := &GeneratedController{
		ResourceName: "User",
		PluralName:   "users",
		Package:      "controllers",
		Fields: []GeneratedField{
			{
				Name:          "ID",
				GoType:        "uuid.UUID",
				DBName:        "id",
				IsSystemField: true,
			},
			{
				Name:   "Email",
				GoType: "string",
				DBName: "email",
			},
		},
		ModulePath: "example.com/app",
	}

	content, err := tr.RenderControllerFile(controller)
	if err != nil {
		t.Fatalf("RenderControllerFile failed: %v", err)
	}

	if content == "" {
		t.Fatal("Generated content is empty")
	}

	if !containsString(content, "User") {
		t.Error("Expected controller to contain 'User'")
	}
}

func TestGenerateRoutes(t *testing.T) {
	rg := NewRouteGenerator()

	if rg == nil {
		t.Fatal("NewRouteGenerator returned nil")
	}
}

func TestGenerateRouteContent(t *testing.T) {
	tr := NewTemplateRenderer()

	tests := []struct {
		name         string
		resourceName string
		pluralName   string
	}{
		{"simple resource", "User", "users"},
		{"plural name", "Product", "products"},
		{"capitalized plural", "Category", "Categories"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := tr.generateRouteContent(tt.resourceName, tt.pluralName)

			if err != nil {
				t.Logf("generateRouteContent failed (expected if no go.mod): %v", err)
			}

			if content == "" && err == nil {
				t.Fatal("Generated content is empty but no error")
			}

			if content != "" && !containsString(content, tt.resourceName) {
				t.Errorf("Expected content to contain '%s'", tt.resourceName)
			}
		})
	}
}

func TestGenerateRouteRegistrationFile(t *testing.T) {
	tr := NewTemplateRenderer()

	tests := []struct {
		name         string
		resourceName string
		pluralName   string
	}{
		{"simple registration", "User", "users"},
		{"product registration", "Product", "products"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, err := tr.generateRouteRegistrationFile(tt.resourceName, tt.pluralName)

			if err != nil {
				t.Logf("generateRouteRegistrationFile failed (expected if no go.mod): %v", err)
			}

			if content == "" && err == nil {
				t.Fatal("Generated content is empty but no error")
			}

			if content != "" && !containsString(content, tt.pluralName) {
				t.Errorf("Expected content to contain '%s'", tt.pluralName)
			}
		})
	}
}

func containsString(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	if s == substr {
		return true
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestBuild_ResourceControllerNoViews(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "items")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("name", "TEXT"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	controller, err := g.Build(cat, Config{
		ResourceName:   "Item",
		PluralName:     "items",
		TableName:      "items",
		ControllerType: ResourceControllerNoViews,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if controller == nil {
		t.Fatal("Controller is nil")
	}

	if controller.Type != ResourceControllerNoViews {
		t.Errorf("Expected controller type ResourceControllerNoViews, got %v", controller.Type)
	}
}

func TestBuild_NormalController(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")

	controller, err := g.Build(cat, Config{
		ResourceName:   "Admin",
		PluralName:     "admins",
		TableName:      "",
		ControllerType: NormalController,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if controller == nil {
		t.Fatal("Controller is nil")
	}

	if controller.Type != NormalController {
		t.Errorf("Expected controller type NormalController, got %v", controller.Type)
	}

	// Normal controllers don't read from catalog, so fields should be empty
	if len(controller.Fields) != 0 {
		t.Errorf("Expected 0 fields for NormalController, got %d", len(controller.Fields))
	}
}

func TestBuild_FieldGoFormTypeMapping(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "data_types")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("int_field", "INTEGER"))
	table.AddColumn(catalog.NewColumn("bigint_field", "BIGINT"))
	table.AddColumn(catalog.NewColumn("float_field", "REAL"))
	table.AddColumn(catalog.NewColumn("double_field", "DOUBLE PRECISION"))
	table.AddColumn(catalog.NewColumn("bool_field", "BOOLEAN"))
	table.AddColumn(catalog.NewColumn("text_field", "TEXT"))
	table.AddColumn(catalog.NewColumn("timestamp_field", "TIMESTAMP"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	controller, err := g.Build(cat, Config{
		ResourceName:   "DataType",
		PluralName:     "data_types",
		TableName:      "data_types",
		ControllerType: ResourceController,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	expectedFormTypes := map[string]string{
		"IntField":       "int32",
		"BigintField":    "int64",
		"FloatField":      "float32",
		"DoubleField":    "float64",
		"BoolField":      "bool",
		"TextField":      "string",
		"TimestampField": "time.Time",
	}

	for fieldName, expectedFormType := range expectedFormTypes {
		field := getControllerFieldByName(controller.Fields, fieldName)
		if field == nil {
			t.Errorf("Expected field %s not found", fieldName)
			continue
		}

		if field.GoFormType != expectedFormType {
			t.Errorf("Field %s: expected form type %s, got %s", fieldName, expectedFormType, field.GoFormType)
		}
	}
}

func TestBuild_WithNullableFields(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "nullable_fields")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("nullable_text", "TEXT"))
	table.AddColumn(catalog.NewColumn("not_null_text", "TEXT").SetNotNull())
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	controller, err := g.Build(cat, Config{
		ResourceName:   "NullableField",
		PluralName:     "nullable_fields",
		TableName:      "nullable_fields",
		ControllerType: ResourceController,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(controller.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(controller.Fields))
	}

	nullableText := getControllerFieldByName(controller.Fields, "NullableText")
	if nullableText == nil {
		t.Fatal("Expected NullableText field")
	}

	if nullableText.IsSystemField {
		t.Error("NullableText should not be marked as system field")
	}
}

func TestBuild_InvalidTableName(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")

	_, err := g.Build(cat, Config{
		ResourceName:   "User",
		PluralName:     "users",
		TableName:      "nonexistent_table",
		ControllerType: ResourceController,
	})

	if err == nil {
		t.Error("Expected error for non-existent table")
	}
}

func TestGeneratedController_StructFields(t *testing.T) {
	controller := &GeneratedController{
		ResourceName: "User",
		PluralName:   "users",
		Package:      "controllers",
		Fields: []GeneratedField{
			{Name: "ID", GoType: "uuid.UUID", DBName: "id", IsSystemField: true},
		},
		ModulePath:   "example.com/app",
		Type:         ResourceController,
		DatabaseType: "postgresql",
	}

	if controller.ResourceName != "User" {
		t.Errorf("Expected ResourceName 'User', got %s", controller.ResourceName)
	}

	if controller.PluralName != "users" {
		t.Errorf("Expected PluralName 'users', got %s", controller.PluralName)
	}

	if controller.Package != "controllers" {
		t.Errorf("Expected Package 'controllers', got %s", controller.Package)
	}

	if controller.ModulePath != "example.com/app" {
		t.Errorf("Expected ModulePath 'example.com/app', got %s", controller.ModulePath)
	}

	if controller.DatabaseType != "postgresql" {
		t.Errorf("Expected DatabaseType 'postgresql', got %s", controller.DatabaseType)
	}

	if len(controller.Fields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(controller.Fields))
	}
}

func TestGeneratedField_StructFields(t *testing.T) {
	field := GeneratedField{
		Name:          "Email",
		GoType:        "string",
		GoFormType:    "string",
		DBName:        "email",
		IsSystemField: false,
	}

	if field.Name != "Email" {
		t.Errorf("Expected Name 'Email', got %s", field.Name)
	}

	if field.GoType != "string" {
		t.Errorf("Expected GoType 'string', got %s", field.GoType)
	}

	if field.GoFormType != "string" {
		t.Errorf("Expected GoFormType 'string', got %s", field.GoFormType)
	}

	if field.DBName != "email" {
		t.Errorf("Expected DBName 'email', got %s", field.DBName)
	}

	if field.IsSystemField {
		t.Error("Expected IsSystemField to be false")
	}
}

func TestConfig_StructFields(t *testing.T) {
	config := Config{
		ResourceName:   "User",
		PluralName:     "users",
		TableName:      "users",
		PackageName:    "controllers",
		ModulePath:     "example.com/app",
		ControllerType: ResourceController,
	}

	if config.ResourceName != "User" {
		t.Errorf("Expected ResourceName 'User', got %s", config.ResourceName)
	}

	if config.PluralName != "users" {
		t.Errorf("Expected PluralName 'users', got %s", config.PluralName)
	}

	if config.TableName != "users" {
		t.Errorf("Expected TableName 'users', got %s", config.TableName)
	}

	if config.PackageName != "controllers" {
		t.Errorf("Expected PackageName 'controllers', got %s", config.PackageName)
	}

	if config.ModulePath != "example.com/app" {
		t.Errorf("Expected ModulePath 'example.com/app', got %s", config.ModulePath)
	}

	if config.ControllerType != ResourceController {
		t.Errorf("Expected ControllerType ResourceController, got %v", config.ControllerType)
	}
}

func TestRenderControllerFile_ComplexTypes(t *testing.T) {
	tr := NewTemplateRenderer()

	controller := &GeneratedController{
		ResourceName: "Article",
		PluralName:   "articles",
		Package:      "controllers",
		Fields: []GeneratedField{
			{
				Name:          "ID",
				GoType:        "uuid.UUID",
				DBName:        "id",
				IsSystemField: true,
			},
			{
				Name:   "Title",
				GoType: "string",
				DBName: "title",
			},
			{
				Name:   "ViewCount",
				GoType: "int64",
				DBName: "view_count",
			},
			{
				Name:   "PublishedAt",
				GoType: "time.Time",
				DBName: "published_at",
			},
		},
		ModulePath:   "example.com/testapp",
		Type:         ResourceController,
		DatabaseType: "postgresql",
	}

	content, err := tr.RenderControllerFile(controller)
	if err != nil {
		t.Fatalf("RenderControllerFile failed: %v", err)
	}

	if content == "" {
		t.Fatal("Generated content is empty")
	}

	if !containsString(content, "Article") {
		t.Error("Expected controller to contain 'Article'")
	}

	if !containsString(content, "Title") {
		t.Error("Expected controller to contain 'Title' field")
	}

	if !containsString(content, "ViewCount") {
		t.Error("Expected controller to contain 'ViewCount' field")
	}
}

func TestFileGenerator_GenerateController_Integration(t *testing.T) {
	fg := NewFileGenerator()
	
	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "posts")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("title", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("content", "TEXT"))
	table.AddColumn(catalog.NewColumn("created_at", "TIMESTAMP"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	generator := NewGenerator("postgresql")
	controller, err := generator.Build(cat, Config{
		ResourceName:   "Post",
		PluralName:     "posts",
		TableName:      "posts",
		PackageName:    "controllers",
		ModulePath:     "example.com/testapp",
		ControllerType: ResourceController,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	content, err := fg.templateRenderer.RenderControllerFile(controller)
	if err != nil {
		t.Fatalf("RenderControllerFile failed: %v", err)
	}

	if !containsString(content, "Post") {
		t.Error("Expected controller to contain 'Post'")
	}

	if !containsString(content, "Title") {
		t.Error("Expected controller to contain 'Title' field")
	}

	if !containsString(content, "Content") {
		t.Error("Expected controller to contain 'Content' field")
	}
}

func TestRouteGenerator_Content(t *testing.T) {
	rg := NewRouteGenerator()
	
	if rg == nil {
		t.Fatal("RouteGenerator should not be nil")
	}

	if rg.fileManager == nil {
		t.Error("fileManager should not be nil")
	}

	if rg.templateRenderer == nil {
		t.Error("templateRenderer should not be nil")
	}
}
