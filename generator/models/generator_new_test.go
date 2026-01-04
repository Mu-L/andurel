package models

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mbvlabs/andurel/generator/internal/catalog"
	"github.com/mbvlabs/andurel/generator/internal/types"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator("postgresql")
	if g == nil {
		t.Fatal("NewGenerator returned nil")
	}
	if g.databaseType != "postgresql" {
		t.Errorf("Expected database type 'postgresql', got %s", g.databaseType)
	}
	if g.typeMapper == nil {
		t.Error("typeMapper is nil")
	}
	if g.typeMapper.GetDatabaseType() != "postgresql" {
		t.Errorf("TypeMapper database type mismatch: expected 'postgresql', got %s",
			g.typeMapper.GetDatabaseType())
	}
}

func TestBuild_SimpleTable(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "users")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("email", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("created_at", "TIMESTAMP").SetDefault("now()"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	model, err := g.Build(cat, Config{
		TableName:    "users",
		ResourceName: "User",
		PackageName:  "models",
		DatabaseType: "postgresql",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if model == nil {
		t.Fatal("Model is nil")
	}

	if model.Name != "User" {
		t.Errorf("Expected model name 'User', got '%s'", model.Name)
	}

	if model.Package != "models" {
		t.Errorf("Expected package 'models', got '%s'", model.Package)
	}

	if model.TableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", model.TableName)
	}

	if len(model.Fields) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(model.Fields))
	}

	idField := model.Fields[0]
	if idField.Name != "ID" {
		t.Errorf("Expected first field name 'ID', got '%s'", idField.Name)
	}
	if idField.Type != "uuid.UUID" {
		t.Errorf("Expected ID type 'uuid.UUID', got '%s'", idField.Type)
	}

	foundUUIDImport := false
	for _, imp := range model.Imports {
		if imp == "github.com/google/uuid" {
			foundUUIDImport = true
		}
	}
	if !foundUUIDImport {
		t.Error("Expected uuid import")
	}
}

func TestBuild_WithNullableFields(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "posts")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("title", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("content", "TEXT"))
	table.AddColumn(catalog.NewColumn("author_id", "UUID"))
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	model, err := g.Build(cat, Config{
		TableName:    "posts",
		ResourceName: "Post",
		PackageName:  "models",
		DatabaseType: "postgresql",
		ModulePath:   "example.com/app",
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(model.Fields) != 4 {
		t.Errorf("Expected 4 fields, got %d", len(model.Fields))
	}

	contentField := model.Fields[2]
	if contentField.Name != "Content" {
		t.Errorf("Expected field name 'Content', got '%s'", contentField.Name)
	}
}

func TestBuild_WithCustomTypes(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")
	table := catalog.NewTable("public", "products")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("name", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("price", "DECIMAL").SetNotNull())
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("Failed to add table: %v", err)
	}

	customTypes := []types.TypeOverride{
		{DatabaseType: "postgresql", GoType: "decimal.Decimal", Package: "github.com/shopspring/decimal"},
	}

	model, err := g.Build(cat, Config{
		TableName:    "products",
		ResourceName: "Product",
		PackageName:  "models",
		DatabaseType: "postgresql",
		ModulePath:   "example.com/app",
		CustomTypes:  customTypes,
	})

	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	priceField := model.Fields[2]
	if priceField.Name == "Price" || priceField.Name == "price" {
		if priceField.Type == "decimal.Decimal" {
			foundDecimalImport := false
			for _, imp := range model.Imports {
				if imp == "github.com/shopspring/decimal" {
					foundDecimalImport = true
				}
			}
			if !foundDecimalImport {
				t.Error("Expected decimal import for custom type")
			}
		}
	}
}

func TestBuild_TableNotFound(t *testing.T) {
	g := NewGenerator("postgresql")

	cat := catalog.NewCatalog("public")

	_, err := g.Build(cat, Config{
		TableName:    "nonexistent",
		ResourceName: "NonExistent",
		PackageName:  "models",
		DatabaseType: "postgresql",
		ModulePath:   "example.com/app",
	})

	if err == nil {
		t.Error("Expected error for non-existent table, got nil")
	}
}

func TestGenerateSQLContent(t *testing.T) {
	g := NewGenerator("postgresql")

	table := catalog.NewTable("public", "users")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("email", "TEXT").SetNotNull())
	table.AddColumn(catalog.NewColumn("created_at", "TIMESTAMP").SetDefault("now()"))
	table.AddColumn(catalog.NewColumn("updated_at", "TIMESTAMP").SetDefault("now()"))

	sqlContent, err := g.GenerateSQLContent("User", "users", table)
	if err != nil {
		t.Fatalf("GenerateSQLContent failed: %v", err)
	}

	if sqlContent == "" {
		t.Fatal("SQL content is empty")
	}

	expectedQueries := []string{
		"-- name: QueryUserByID",
		"-- name: QueryUsers",
		"-- name: InsertUser",
		"-- name: UpdateUser",
		"-- name: DeleteUser",
	}

	for _, expected := range expectedQueries {
		if !containsString(sqlContent, expected) {
			t.Errorf("Expected SQL to contain '%s'", expected)
		}
	}

	if !containsString(sqlContent, "$1") {
		t.Error("Expected PostgreSQL placeholders ($1, $2, etc.)")
	}

	if !containsString(sqlContent, "now()") {
		t.Error("Expected now() function for timestamps")
	}
}

func TestGenerateSQLContent_WithoutTimestamps(t *testing.T) {
	g := NewGenerator("postgresql")

	table := catalog.NewTable("public", "categories")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("name", "TEXT").SetNotNull())

	sqlContent, err := g.GenerateSQLContent("Category", "categories", table)
	if err != nil {
		t.Fatalf("GenerateSQLContent failed: %v", err)
	}

	if sqlContent == "" {
		t.Fatal("SQL content is empty")
	}

	if containsString(sqlContent, "now()") {
		t.Error("Did not expect now() for table without timestamps")
	}
}

func TestGenerateSQLFile(t *testing.T) {
	g := NewGenerator("postgresql")

	table := catalog.NewTable("public", "test_table")
	table.AddColumn(catalog.NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(catalog.NewColumn("name", "TEXT").SetNotNull())

	tempDir := t.TempDir()
	sqlPath := filepath.Join(tempDir, "test_table.sql")

	err := g.GenerateSQLFile("TestTable", "test_table", table, sqlPath)
	if err != nil {
		t.Fatalf("GenerateSQLFile failed: %v", err)
	}

	if _, err := os.Stat(sqlPath); os.IsNotExist(err) {
		t.Error("SQL file was not created")
	}

	content, err := os.ReadFile(sqlPath)
	if err != nil {
		t.Fatalf("Failed to read SQL file: %v", err)
	}

	if len(content) == 0 {
		t.Error("SQL file is empty")
	}
}

func TestGenerateModelFile(t *testing.T) {
	g := NewGenerator("postgresql")

	model := &GeneratedModel{
		Name:      "User",
		Package:   "models",
		TableName: "users",
		Fields: []GeneratedField{
			{
				Name:             "ID",
				Type:             "uuid.UUID",
				SQLCType:         "uuid.UUID",
				ConversionFromDB: "row.ID",
				ConversionToDB:   "data.ID",
				ZeroCheck:        "data.ID != uuid.Nil",
			},
			{
				Name:             "Email",
				Type:             "string",
				SQLCType:         "string",
				ConversionFromDB: "row.Email",
				ConversionToDB:   "data.Email",
				ZeroCheck:        "data.Email != \"\"",
			},
		},
		Imports:           []string{"github.com/google/uuid"},
		StandardImports:   []string{},
		ExternalImports:   []string{"github.com/google/uuid"},
		DatabaseType:      "postgresql",
	}

	templateStr := `package {{.Package}}

type {{.Name}} struct {
{{- range .Fields}}
	{{.Name}} {{.Type}}
{{- end}}
}
`

	content, err := g.GenerateModelFile(model, templateStr)
	if err != nil {
		t.Fatalf("GenerateModelFile failed: %v", err)
	}

	if content == "" {
		t.Fatal("Generated content is empty")
	}

	if !containsString(content, "type User struct") {
		t.Error("Expected 'type User struct' in generated code")
	}

	if !containsString(content, "ID uuid.UUID") {
		t.Error("Expected 'ID uuid.UUID' field")
	}

	if !containsString(content, "Email string") {
		t.Error("Expected 'Email string' field")
	}
}

func TestGenerateModelFile_InvalidTemplate(t *testing.T) {
	g := NewGenerator("postgresql")

	model := &GeneratedModel{
		Name:      "User",
		Package:   "models",
		TableName: "users",
		Fields:    []GeneratedField{},
	}

	invalidTemplate := `package {{.Package}}

type {{.Name}} struct {
	{{.UnknownField}
}
`

	_, err := g.GenerateModelFile(model, invalidTemplate)
	if err == nil {
		t.Error("Expected error for invalid template, got nil")
	}
}

func TestGenerateConstructorFile(t *testing.T) {
	g := NewGenerator("postgresql")

	model := &GeneratedModel{
		Name:      "User",
		Package:   "db",
		TableName: "users",
		Fields: []GeneratedField{
			{Name: "ID", Type: "uuid.UUID", SQLCType: "uuid.UUID"},
			{Name: "Email", Type: "string", SQLCType: "string"},
		},
		Imports:      []string{"github.com/google/uuid"},
		DatabaseType: "postgresql",
	}

	templateStr := `package {{.Package}}

import (
	{{range .Imports}}"{{.}}"{{end}}
)

func NewUser(id uuid.UUID, email string) *User {
	return &User{
		ID:    id,
		Email: email,
	}
}
`

	content, err := g.GenerateConstructorFile(model, templateStr)
	if err != nil {
		t.Fatalf("GenerateConstructorFile failed: %v", err)
	}

	if content == "" {
		t.Fatal("Generated content is empty")
	}

	if !containsString(content, "package db") {
		t.Error("Expected 'package db'")
	}

	if !containsString(content, "func NewUser") {
		t.Error("Expected 'func NewUser' constructor")
	}
}

func TestExtractGeneratedParts(t *testing.T) {
	g := NewGenerator("postgresql")

	content := `package models

type User struct {
	ID    uuid.UUID
	Email string
}

type CreateUserData struct {
	Email string
}

func FindUser(id uuid.UUID) (*User, error) {
	return &User{}, nil
}
`

	parts := g.extractGeneratedParts(content, "User")
	if len(parts) == 0 {
		t.Error("Expected to extract generated parts")
	}

	expectedSignatures := []string{
		"type User struct",
		"type CreateUserData struct",
		"func FindUser(",
	}

	for _, expected := range expectedSignatures {
		if _, ok := parts[expected]; !ok {
			t.Errorf("Expected to extract part with signature '%s'", expected)
		}
	}
}

func TestExtractPartBySignature(t *testing.T) {
	g := NewGenerator("postgresql")

	lines := []string{
		"package models",
		"",
		"type User struct {",
		"	ID    uuid.UUID",
		"	Email string",
		"}",
		"",
		"func NewUser() *User {",
		"	return &User{}",
		"}",
	}

	structPart := g.extractPartBySignature(lines, "type User struct")
	if structPart == "" {
		t.Error("Failed to extract struct part")
	}

	if !containsString(structPart, "ID    uuid.UUID") {
		t.Error("Expected struct to contain ID field")
	}

	funcPart := g.extractPartBySignature(lines, "func NewUser(")
	if funcPart == "" {
		t.Error("Failed to extract function part")
	}
}

func TestReplacePartBySignature(t *testing.T) {
	g := NewGenerator("postgresql")

	content := `package models

type User struct {
	ID    uuid.UUID
	Email string
}

func FindUser(id uuid.UUID) (*User, error) {
	return nil, nil
}
`

	newStruct := `type User struct {
	ID      uuid.UUID
	Email   string
	Name    string
	CreatedAt time.Time
}`

	updated := g.replacePartBySignature(content, "type User struct", newStruct)

	if !containsString(updated, "Name    string") {
		t.Error("Expected updated struct to contain Name field")
	}

	if !containsString(updated, "CreatedAt time.Time") {
		t.Error("Expected updated struct to contain CreatedAt field")
	}

	if !containsString(updated, "func FindUser") {
		t.Error("Expected other functions to be preserved")
	}
}

func TestExtractGeneratedSQLQueries(t *testing.T) {
	g := NewGenerator("postgresql")

	content := `-- name: QueryUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: InsertUser :one
INSERT INTO users (email) VALUES ($1) RETURNING *;

-- name: UpdateUser :one
UPDATE users SET email = $2 WHERE id = $1 RETURNING *;
`

	queries := g.extractGeneratedSQLQueries(content, "User")

	expectedQueries := []string{
		"QueryUserByID",
		"InsertUser",
		"UpdateUser",
	}

	for _, expected := range expectedQueries {
		if _, ok := queries[expected]; !ok {
			t.Errorf("Expected to extract query '%s'", expected)
		}
	}
}

func TestExtractSQLQueryByName(t *testing.T) {
	g := NewGenerator("postgresql")

	lines := []string{
		"-- name: QueryUserByID :one",
		"SELECT * FROM users WHERE id = $1 LIMIT 1;",
		"",
		"-- name: InsertUser :one",
		"INSERT INTO users (email) VALUES ($1) RETURNING *;",
	}

	query := g.extractSQLQueryByName(lines, "QueryUserByID")
	if query == "" {
		t.Fatal("Failed to extract query")
	}

	if !containsString(query, "SELECT * FROM users") {
		t.Error("Expected query to contain SELECT statement")
	}
}

func TestReplaceSQLQueryByName(t *testing.T) {
	g := NewGenerator("postgresql")

	content := `-- name: QueryUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: InsertUser :one
INSERT INTO users (email) VALUES ($1) RETURNING *;
`

	newQuery := `-- name: QueryUserByID :one
SELECT id, email, created_at FROM users WHERE id = $1 LIMIT 1;`

	updated := g.replaceSQLQueryByName(content, "QueryUserByID", newQuery)

	if !containsString(updated, "id, email, created_at") {
		t.Error("Expected updated query to contain new columns")
	}

	if !containsString(updated, "-- name: InsertUser") {
		t.Error("Expected other queries to be preserved")
	}
}

func TestQueryExistsInContent(t *testing.T) {
	g := NewGenerator("postgresql")

	content := `-- name: QueryUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: InsertUser :one
INSERT INTO users (email) VALUES ($1);
`

	if !g.queryExistsInContent(content, "QueryUserByID") {
		t.Error("Expected QueryUserByID to exist")
	}

	if !g.queryExistsInContent(content, "InsertUser") {
		t.Error("Expected InsertUser to exist")
	}

	if g.queryExistsInContent(content, "NonExistent") {
		t.Error("Expected NonExistent query to not exist")
	}
}

func TestReplaceGeneratedSQLQueries(t *testing.T) {
	g := NewGenerator("postgresql")

	existingContent := `-- name: QueryUserByID :one
SELECT id FROM users WHERE id = $1;
`

	newContent := `-- name: QueryUserByID :one
SELECT id, email, created_at FROM users WHERE id = $1;

-- name: InsertUser :one
INSERT INTO users (email) VALUES ($1) RETURNING *;
`

	updated := g.replaceGeneratedSQLQueries(existingContent, newContent, "User")

	if !containsString(updated, "id, email, created_at") {
		t.Error("Expected query to be updated")
	}

	if !containsString(updated, "-- name: InsertUser") {
		t.Error("Expected new query to be added")
	}
}

func TestAddModelTypeImports(t *testing.T) {
	g := NewGenerator("postgresql")

	testCases := []struct {
		goType         string
		expectedImport string
	}{
		{"time.Time", "time"},
		{"[]time.Time", "time"},
		{"uuid.UUID", "github.com/google/uuid"},
		{"[]uuid.UUID", "github.com/google/uuid"},
		{"string", ""},
		{"int", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.goType, func(t *testing.T) {
			imports := g.addModelTypeImports(tc.goType)
			if tc.expectedImport != "" {
				if !imports[tc.expectedImport] {
					t.Errorf("Expected import '%s' for type '%s'", tc.expectedImport, tc.goType)
				}
			} else {
				if len(imports) != 0 {
					t.Errorf("Expected no imports for type '%s', got %v", tc.goType, imports)
				}
			}
		})
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
