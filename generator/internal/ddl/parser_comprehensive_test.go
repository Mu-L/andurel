package ddl

import (
	"testing"

	"github.com/mbvlabs/andurel/generator/internal/catalog"
)

func TestNewDDLParser(t *testing.T) {
	parser := NewDDLParser()
	
	if parser == nil {
		t.Fatal("expected non-nil parser")
	}
	
	if parser.createTableParser == nil {
		t.Error("expected non-nil createTableParser")
	}
	if parser.alterTableParser == nil {
		t.Error("expected non-nil alterTableParser")
	}
}

func TestDDLParser_Parse_CreateTable(t *testing.T) {
	parser := NewDDLParser()
	
	ddl := `CREATE TABLE users (
		id UUID PRIMARY KEY,
		email TEXT NOT NULL
	);`
	
	stmt, err := parser.Parse(ddl, "001_create_users.sql", "postgresql")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	
	if stmt.GetType() != CreateTable {
		t.Errorf("expected CreateTable type, got %v", stmt.GetType())
	}
}

func TestDDLParser_Parse_AlterTable(t *testing.T) {
	parser := NewDDLParser()
	
	ddl := `ALTER TABLE users ADD COLUMN age INTEGER;`
	
	stmt, err := parser.Parse(ddl, "002_alter_users.sql", "postgresql")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	
	if stmt.GetType() != AlterTable {
		t.Errorf("expected AlterTable type, got %v", stmt.GetType())
	}
}

func TestDDLParser_Parse_DropTable(t *testing.T) {
	parser := NewDDLParser()
	
	ddl := `DROP TABLE users;`
	
	stmt, err := parser.Parse(ddl, "003_drop_users.sql", "postgresql")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	
	if stmt.GetType() != DropTable {
		t.Errorf("expected DropTable type, got %v", stmt.GetType())
	}
}

func TestDDLParser_Parse_CreateIndex(t *testing.T) {
	parser := NewDDLParser()
	
	ddl := `CREATE INDEX idx_users_email ON users(email);`
	
	stmt, err := parser.Parse(ddl, "004_create_index.sql", "postgresql")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	
	if stmt.GetType() != CreateIndex {
		t.Errorf("expected CreateIndex type, got %v", stmt.GetType())
	}
}

func TestDDLParser_Parse_DropIndex(t *testing.T) {
	parser := NewDDLParser()
	
	ddl := `DROP INDEX idx_users_email;`
	
	stmt, err := parser.Parse(ddl, "005_drop_index.sql", "postgresql")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	
	if stmt.GetType() != DropIndex {
		t.Errorf("expected DropIndex type, got %v", stmt.GetType())
	}
}

func TestDDLParser_Parse_CreateSchema(t *testing.T) {
	parser := NewDDLParser()
	
	ddl := `CREATE SCHEMA custom;`
	
	stmt, err := parser.Parse(ddl, "006_create_schema.sql", "postgresql")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	
	if stmt.GetType() != CreateSchema {
		t.Errorf("expected CreateSchema type, got %v", stmt.GetType())
	}
}

func TestDDLParser_Parse_EmptyString(t *testing.T) {
	parser := NewDDLParser()
	
	stmt, err := parser.Parse("", "test.sql", "postgresql")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if stmt != nil {
		t.Error("expected nil statement for empty string")
	}
}

func TestDDLParser_Parse_Unknown(t *testing.T) {
	parser := NewDDLParser()
	
	ddl := `SELECT * FROM users;`
	
	stmt, err := parser.Parse(ddl, "query.sql", "postgresql")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	
	if stmt == nil {
		t.Fatal("expected non-nil statement")
	}
	
	if stmt.GetType() != Unknown {
		t.Errorf("expected Unknown statement type, got %v", stmt.GetType())
	}
}

func TestApplyDDL_CreateTable(t *testing.T) {
	cat := catalog.NewCatalog("public")
	
	ddl := `CREATE TABLE users (
		id UUID PRIMARY KEY,
		email TEXT NOT NULL
	);`
	
	err := ApplyDDL(cat, ddl, "001_create_users.sql", "postgresql")
	if err != nil {
		t.Fatalf("ApplyDDL failed: %v", err)
	}
	
	// Verify table was created
	table, err := cat.GetTable("public", "users")
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}
	
	if table == nil {
		t.Fatal("expected table to exist")
	}
}

func TestApplyDDL_EmptySQL(t *testing.T) {
	cat := catalog.NewCatalog("public")
	
	err := ApplyDDL(cat, "", "test.sql", "postgresql")
	if err != nil {
		t.Fatalf("ApplyDDL should not error on empty SQL: %v", err)
	}
}

func TestApplyDDL_UnknownStatement(t *testing.T) {
	cat := catalog.NewCatalog("public")
	
	// Unknown statement should be logged but not error
	err := ApplyDDL(cat, "SELECT * FROM users;", "test.sql", "postgresql")
	if err != nil {
		t.Fatalf("ApplyDDL should not error on unknown statement: %v", err)
	}
}

func TestApplyDDL_CreateIndex(t *testing.T) {
	cat := catalog.NewCatalog("public")
	
	// Create index statements should not error (they're skipped)
	err := ApplyDDL(cat, "CREATE INDEX idx_users_email ON users(email);", "test.sql", "postgresql")
	if err != nil {
		t.Fatalf("ApplyDDL should not error on CREATE INDEX: %v", err)
	}
}

func TestNewCatalogVisitor(t *testing.T) {
	cat := catalog.NewCatalog("public")
	visitor := NewCatalogVisitor(cat, "001_test.sql", "postgresql")
	
	if visitor == nil {
		t.Fatal("expected non-nil visitor")
	}
}

func TestNewAlterTableParser(t *testing.T) {
	parser := NewAlterTableParser()
	
	if parser == nil {
		t.Fatal("expected non-nil parser")
	}
}

func TestParseDataType_SimpleTypes(t *testing.T) {
	tests := []struct{
		input    string
		expected string
	}{
		{"INTEGER", "integer"},
		{"TEXT", "text"},
		{"BOOLEAN", "boolean"},
		{"UUID", "uuid"},
		{"TIMESTAMP", "timestamp"},
		{"DATE", "date"},
		{"JSONB", "jsonb"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			dataType, _, _, _ := ParseDataType(tt.input)
			if dataType != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, dataType)
			}
		})
	}
}

func TestParseDataType_WithLength(t *testing.T) {
	dataType, length, _, _ := ParseDataType("VARCHAR(255)")
	
	if dataType != "varchar" {
		t.Errorf("expected varchar, got %s", dataType)
	}
	
	if length == nil {
		t.Fatal("expected non-nil length")
	}
	
	if *length != 255 {
		t.Errorf("expected length 255, got %d", *length)
	}
}

func TestParseDataType_WithPrecisionScale(t *testing.T) {
	dataType, _, precision, scale := ParseDataType("DECIMAL(10,2)")
	
	if dataType != "decimal" {
		t.Errorf("expected decimal, got %s", dataType)
	}
	
	if precision == nil || scale == nil {
		t.Fatal("expected non-nil precision and scale")
	}
	
	if *precision != 10 {
		t.Errorf("expected precision 10, got %d", *precision)
	}
	
	if *scale != 2 {
		t.Errorf("expected scale 2, got %d", *scale)
	}
}

func TestParseDataType_TimestampWithTimeZone(t *testing.T) {
	dataType, _, _, _ := ParseDataType("TIMESTAMP WITH TIME ZONE")
	
	if dataType != "timestamp with time zone" {
		t.Errorf("expected 'timestamp with time zone', got %s", dataType)
	}
}
