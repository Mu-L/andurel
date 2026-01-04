package views

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

	view, err := g.Build(cat, Config{
		ResourceName: "User",
		PluralName:   "users",
		TableName:    "users",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if view == nil {
		t.Fatal("View is nil")
	}

	if view.ResourceName != "User" {
		t.Errorf("Expected resource name 'User', got '%s'", view.ResourceName)
	}

	if view.PluralName != "users" {
		t.Errorf("Expected plural name 'users', got '%s'", view.PluralName)
	}

	if view.ModulePath != "example.com/app" {
		t.Errorf("Expected module path 'example.com/app', got '%s'", view.ModulePath)
	}

	// Note: views skip the id column
	if len(view.Fields) != 2 {
		t.Errorf("Expected 2 fields (id is skipped), got %d", len(view.Fields))
	}

	emailField := view.Fields[0]
	if emailField.Name != "Email" {
		t.Errorf("Expected first field name 'Email', got '%s'", emailField.Name)
	}

	if emailField.GoType != "string" {
		t.Errorf("Expected Email type 'string', got '%s'", emailField.GoType)
	}
}

func TestBuild_TableNotFound(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")

	_, err := g.Build(cat, Config{
		ResourceName: "NonExistent",
		PluralName:   "nonexistents",
		TableName:    "nonexistent",
		ModulePath:   "example.com/app",
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
	table.AddColumn(catalog.NewColumn("content", "TEXT"))
	table.AddColumn(catalog.NewColumn("created_at", "TIMESTAMP"))
	table.AddColumn(catalog.NewColumn("updated_at", "TIMESTAMP"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	view, err := g.Build(cat, Config{
		ResourceName: "Post",
		PluralName:   "posts",
		TableName:    "posts",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	// id is skipped, so 4 fields remain
	if len(view.Fields) != 4 {
		t.Errorf("Expected 4 fields (id is skipped), got %d", len(view.Fields))
	}

	createdAtField := getViewFieldByName(view.Fields, "CreatedAt")
	if createdAtField == nil {
		t.Fatal("Expected CreatedAt field")
	}

	if !createdAtField.IsTimestamp {
		t.Error("Expected CreatedAt to be marked as timestamp")
	}

	if !createdAtField.IsSystemField {
		t.Error("Expected CreatedAt to be a system field")
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

	view, err := g.Build(cat, Config{
		ResourceName: "Item",
		PluralName:   "items",
		TableName:    "items",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	systemFieldCount := 0
	for _, field := range view.Fields {
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

	view, err := g.Build(cat, Config{
		ResourceName: "Product",
		PluralName:   "products",
		TableName:    "products",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	expectedTypes := map[string]string{
		"Name":        "string",
		"Price":       "float64",
		"Quantity":    "int32",
		"IsAvailable": "bool",
	}

	for fieldName, expectedType := range expectedTypes {
		field := getViewFieldByName(view.Fields, fieldName)
		if field == nil {
			t.Errorf("Expected field '%s' not found", fieldName)
			continue
		}
		if field.GoType != expectedType {
			t.Errorf("Field %s: expected type '%s', got '%s'", fieldName, expectedType, field.GoType)
		}
	}
}

func TestBuild_DisplayNameGeneration(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "profiles")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("first_name", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("last_name", "TEXT"))
	table.AddColumn(catalog.NewColumn("email_address", "TEXT"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	view, err := g.Build(cat, Config{
		ResourceName: "Profile",
		PluralName:   "profiles",
		TableName:    "profiles",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	firstNameField := getViewFieldByName(view.Fields, "FirstName")
	if firstNameField == nil {
		t.Fatal("Expected FirstName field")
	}

	if firstNameField.DisplayName != "First Name" {
		t.Errorf("Expected display name 'First Name', got '%s'", firstNameField.DisplayName)
	}

	lastNameField := getViewFieldByName(view.Fields, "LastName")
	if lastNameField == nil {
		t.Fatal("Expected LastName field")
	}

	if lastNameField.DisplayName != "Last Name" {
		t.Errorf("Expected display name 'Last Name', got '%s'", lastNameField.DisplayName)
	}
}

func TestBuild_InputTypeGeneration(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "users")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("email", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("age", "INTEGER"))
	table.AddColumn(catalog.NewColumn("bio", "TEXT"))
	table.AddColumn(catalog.NewColumn("is_active", "BOOLEAN"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	view, err := g.Build(cat, Config{
		ResourceName: "User",
		PluralName:   "users",
		TableName:    "users",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	emailField := getViewFieldByName(view.Fields, "Email")
	if emailField != nil && emailField.InputType == "" {
		t.Error("Expected Email field to have an input type")
	}

	bioField := getViewFieldByName(view.Fields, "Bio")
	if bioField != nil {
		// TEXT fields use "text" input type in views
		if bioField.InputType != "text" {
			t.Errorf("Expected Bio field to have text input type, got '%s'", bioField.InputType)
		}
	}
}

func TestGenerateViewFile(t *testing.T) {
	g := NewGenerator("postgresql")

	view := &GeneratedView{
		ResourceName: "User",
		PluralName:   "users",
		ModulePath:   "example.com/app",
		Fields: []ViewField{
			{
				Name:          "ID",
				GoType:        "uuid.UUID",
				DisplayName:   "ID",
				DBName:        "id",
				CamelCase:     "id",
				InputType:     "text",
				IsSystemField: true,
			},
			{
				Name:         "Email",
				GoType:       "string",
				DisplayName:  "Email",
				DBName:       "email",
				CamelCase:    "email",
				InputType:    "email",
			},
		},
	}

	content, err := g.GenerateViewFile(view, false)
	if err != nil {
		t.Fatalf("GenerateViewFile failed: %v", err)
	}

	if content == "" {
		t.Fatal("Generated content is empty")
	}

	if !containsString(content, "User") {
		t.Error("Expected view to contain 'User'")
	}

	if !containsString(content, "Email") {
		t.Error("Expected view to contain 'Email'")
	}
}

func TestBuild_CamelCaseGeneration(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "user_profiles")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("full_name", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("address_line_1", "TEXT"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	view, err := g.Build(cat, Config{
		ResourceName: "UserProfile",
		PluralName:   "user_profiles",
		TableName:    "user_profiles",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	fullNameField := getViewFieldByName(view.Fields, "FullName")
	if fullNameField == nil {
		t.Fatal("Expected FullName field")
	}

	if fullNameField.CamelCase != "fullName" {
		t.Errorf("Expected camelCase 'fullName', got '%s'", fullNameField.CamelCase)
	}

	addressField := getViewFieldByName(view.Fields, "AddressLine1")
	if addressField == nil {
		t.Fatal("Expected AddressLine1 field")
	}

	if addressField.CamelCase != "addressLine1" {
		t.Errorf("Expected camelCase 'addressLine1', got '%s'", addressField.CamelCase)
	}
}

func getViewFieldByName(fields []ViewField, name string) *ViewField {
	for i := range fields {
		if fields[i].Name == name {
			return &fields[i]
		}
	}
	return nil
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
