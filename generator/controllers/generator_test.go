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
