package ddl

import (
	"testing"

	"github.com/mbvlabs/andurel/generator/internal/validation"
)

func TestValidatePrimaryKeyDatatype(t *testing.T) {
	testCases := []struct {
		name         string
		dataType     string
		databaseType string
		expectError  bool
	}{
		{"postgresql_uuid_valid", "UUID", "postgresql", false},
		{"postgresql_uuid_lowercase", "uuid", "postgresql", false},
		{"postgresql_text_invalid", "TEXT", "postgresql", true},
		{"postgresql_integer_invalid", "INTEGER", "postgresql", true},
		{"postgresql_varchar_invalid", "VARCHAR", "postgresql", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validation.ValidatePrimaryKeyDatatype(tc.dataType, tc.databaseType, "test.sql", "id")
			if tc.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidatePrimaryKeyDatatype_ErrorMessages(t *testing.T) {
	testCases := []struct {
		name           string
		dataType       string
		databaseType   string
		columnName     string
		migrationFile  string
		expectedSubstr string
	}{
		{
			name:           "postgresql_text_error_message",
			dataType:       "TEXT",
			databaseType:   "postgresql",
			columnName:     "id",
			migrationFile:  "/path/to/001_create_users.sql",
			expectedSubstr: "primary keys must use 'uuid'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validation.ValidatePrimaryKeyDatatype(
				tc.dataType,
				tc.databaseType,
				tc.migrationFile,
				tc.columnName,
			)
			if err == nil {
				t.Fatal("Expected error but got none")
			}

			errorMsg := err.Error()
			if !containsString(errorMsg, tc.expectedSubstr) {
				t.Errorf(
					"Expected error message to contain '%s', but got: %s",
					tc.expectedSubstr,
					errorMsg,
				)
			}

			if !containsString(errorMsg, tc.columnName) {
				t.Errorf(
					"Expected error message to contain column name '%s', but got: %s",
					tc.columnName,
					errorMsg,
				)
			}

			if !containsString(errorMsg, "001_create_users.sql") {
				t.Errorf(
					"Expected error message to contain migration file name, but got: %s",
					errorMsg,
				)
			}
		})
	}
}

func TestValidatePrimaryKeyDatatype_UnsupportedDatabase(t *testing.T) {
	// For unsupported database types, validation should return an error
	err := validation.ValidatePrimaryKeyDatatype("INTEGER", "mysql", "test.sql", "id")
	if err == nil {
		t.Error("Expected an error for unsupported database type, but got none")
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestParseColumnDefinitions_PrimaryKeyValidation(t *testing.T) {
	testCases := []struct {
		name         string
		columnDefs   string
		databaseType string
		expectError  bool
		errorSubstr  string
	}{
		{
			name:         "postgresql_valid_uuid_primary_key",
			columnDefs:   "id UUID PRIMARY KEY, name TEXT NOT NULL",
			databaseType: "postgresql",
			expectError:  false,
		},
		{
			name:         "postgresql_invalid_text_primary_key",
			columnDefs:   "id TEXT PRIMARY KEY, name TEXT NOT NULL",
			databaseType: "postgresql",
			expectError:  true,
			errorSubstr:  "primary keys must use 'uuid'",
		},
		{
			name:         "postgresql_separate_primary_key_constraint_valid",
			columnDefs:   "id UUID NOT NULL, name TEXT, PRIMARY KEY (id)",
			databaseType: "postgresql",
			expectError:  false,
		},
		{
			name:         "postgresql_separate_primary_key_constraint_invalid",
			columnDefs:   "id INTEGER NOT NULL, name TEXT, PRIMARY KEY (id)",
			databaseType: "postgresql",
			expectError:  true,
			errorSubstr:  "primary keys must use 'uuid'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewCreateTableParser()
			columns, err := parser.parseColumnDefinitions(tc.columnDefs, "test.sql", tc.databaseType)

			if tc.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				if tc.errorSubstr != "" && !containsSubstring(err.Error(), tc.errorSubstr) {
					t.Errorf(
						"Expected error to contain '%s', but got: %s",
						tc.errorSubstr,
						err.Error(),
					)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}

				// Verify that we got the expected columns
				if len(columns) < 2 {
					t.Errorf("Expected at least 2 columns, got %d", len(columns))
				}

				// Find the primary key column and verify it's marked correctly
				var foundPK bool
				for _, col := range columns {
					if col.IsPrimaryKey {
						foundPK = true
						if col.Name != "id" {
							t.Errorf("Expected primary key column to be 'id', got '%s'", col.Name)
						}
					}
				}

				if !foundPK {
					t.Error("Expected to find a primary key column but didn't")
				}
			}
		})
	}
}

func TestParseCreateTable_PrimaryKeyValidation(t *testing.T) {
	testCases := []struct {
		name         string
		sql          string
		databaseType string
		expectError  bool
		errorSubstr  string
	}{
		{
			name:         "postgresql_valid_create_table",
			sql:          "CREATE TABLE users (id UUID PRIMARY KEY, email TEXT NOT NULL)",
			databaseType: "postgresql",
			expectError:  false,
		},
		{
			name:         "postgresql_invalid_create_table",
			sql:          "CREATE TABLE users (id TEXT PRIMARY KEY, email TEXT NOT NULL)",
			databaseType: "postgresql",
			expectError:  true,
			errorSubstr:  "primary keys must use 'uuid'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			parser := NewCreateTableParser()
			stmt, err := parser.Parse(tc.sql, "test.sql", tc.databaseType)

			if tc.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				if tc.errorSubstr != "" && !containsSubstring(err.Error(), tc.errorSubstr) {
					t.Errorf(
						"Expected error to contain '%s', but got: %s",
						tc.errorSubstr,
						err.Error(),
					)
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error but got: %v", err)
				}

				if stmt == nil {
					t.Fatal("Expected statement but got nil")
				}

				if stmt.GetType() != CreateTable {
					t.Errorf("Expected CREATE TABLE statement type, got %v", stmt.GetType())
				}
			}
		})
	}
}

func TestDropTableParser(t *testing.T) {
	parser := NewDropTableParser()

	tests := []struct {
		name          string
		sql           string
		expectError   bool
		expectedTable string
		expectedIfs   bool
	}{
		{"simple drop", "DROP TABLE users", false, "users", false},
		{"drop with if exists", "DROP TABLE IF EXISTS users", false, "users", true},
		{"drop with schema", "DROP TABLE public.users", false, "users", false},
		{"drop schema with if exists", "DROP TABLE IF EXISTS public.users", false, "users", true},
		{"uppercase", "DROP TABLE USERS", false, "USERS", false},
		{"mixed case", "Drop Table If Exists users", false, "users", true},
		{"invalid syntax", "DROP TABLE", true, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.sql)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement but got nil")
			}

			if stmt.TableName != tt.expectedTable {
				t.Errorf("TableName = %s, want %s", stmt.TableName, tt.expectedTable)
			}

			if stmt.IfExists != tt.expectedIfs {
				t.Errorf("IfExists = %v, want %v", stmt.IfExists, tt.expectedIfs)
			}
		})
	}
}

func TestCreateSchemaParser(t *testing.T) {
	parser := NewCreateSchemaParser()

	tests := []struct {
		name          string
		sql           string
		expectError   bool
		expectedName  string
	}{
		{"simple create schema", "CREATE SCHEMA myschema", false, "myschema"},
		{"create if not exists", "CREATE SCHEMA IF NOT EXISTS myschema", false, "myschema"},
		{"uppercase", "CREATE SCHEMA MYSCHEMA", false, "MYSCHEMA"},
		{"no schema name", "CREATE SCHEMA", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.sql)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement but got nil")
			}

			if stmt.SchemaName != tt.expectedName {
				t.Errorf("SchemaName = %s, want %s", stmt.SchemaName, tt.expectedName)
			}
		})
	}
}

func TestDropSchemaParser(t *testing.T) {
	parser := NewDropSchemaParser()

	tests := []struct {
		name          string
		sql           string
		expectError   bool
		expectedName  string
	}{
		{"simple drop schema", "DROP SCHEMA myschema", false, "myschema"},
		{"drop if exists", "DROP SCHEMA IF EXISTS myschema", false, "myschema"},
		{"uppercase", "DROP SCHEMA MYSCHEMA", false, "MYSCHEMA"},
		{"no schema name", "DROP SCHEMA", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.sql)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement but got nil")
			}

			if stmt.SchemaName != tt.expectedName {
				t.Errorf("SchemaName = %s, want %s", stmt.SchemaName, tt.expectedName)
			}
		})
	}
}

func TestCreateEnumParser(t *testing.T) {
	parser := NewCreateEnumParser()

	tests := []struct {
		name          string
		sql           string
		expectError   bool
		expectedName  string
	}{
		{"simple create enum", "CREATE TYPE status AS ENUM", false, "status"},
		{"create enum with schema", "CREATE TYPE public.status AS ENUM", false, "status"},
		{"uppercase", "CREATE TYPE STATUS AS ENUM", false, "STATUS"},
		{"no enum name", "CREATE TYPE AS ENUM", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.sql)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement but got nil")
			}

			if stmt.EnumName != tt.expectedName {
				t.Errorf("EnumName = %s, want %s", stmt.EnumName, tt.expectedName)
			}
		})
	}
}

func TestDropEnumParser(t *testing.T) {
	parser := NewDropEnumParser()

	tests := []struct {
		name          string
		sql           string
		expectError   bool
		expectedName  string
	}{
		{"simple drop enum", "DROP TYPE status", false, "status"},
		{"drop enum with schema", "DROP TYPE public.status", false, "status"},
		{"drop if exists", "DROP TYPE IF EXISTS status", false, "status"},
		{"uppercase", "DROP TYPE STATUS", false, "STATUS"},
		{"no enum name", "DROP TYPE", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.sql)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement but got nil")
			}

			if stmt.EnumName != tt.expectedName {
				t.Errorf("EnumName = %s, want %s", stmt.EnumName, tt.expectedName)
			}
		})
	}
}

func TestCreateIndexParser(t *testing.T) {
	parser := NewCreateIndexParser()

	tests := []struct {
		name        string
		sql         string
		expectError bool
	}{
		{"simple create index", "CREATE INDEX idx_name ON users (email)", false},
		{"unique index", "CREATE UNIQUE INDEX idx_name ON users (email)", false},
		{"uppercase", "CREATE INDEX IDX_NAME ON USERS (EMAIL)", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.sql)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement but got nil")
			}
		})
	}
}

func TestDropIndexParser(t *testing.T) {
	parser := NewDropIndexParser()

	tests := []struct {
		name        string
		sql         string
		expectError bool
	}{
		{"simple drop index", "DROP INDEX idx_name", false},
		{"drop if exists", "DROP INDEX IF EXISTS idx_name", false},
		{"uppercase", "DROP INDEX IDX_NAME", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmt, err := parser.Parse(tt.sql)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error but got: %v", err)
			}

			if stmt == nil {
				t.Fatal("Expected statement but got nil")
			}
		})
	}
}

func TestParseDataType(t *testing.T) {
	tests := []struct {
		name             string
		typeStr          string
		expectedType     string
		expectedLength   *int32
		expectedPrec     *int32
		expectedScale    *int32
	}{
		{"varchar with length", "VARCHAR(255)", "varchar", ptr(255), nil, nil},
		{"varchar lowercase", "varchar(100)", "varchar", ptr(100), nil, nil},
		{"char with length", "CHAR(10)", "char", ptr(10), nil, nil},
		{"decimal with precision", "DECIMAL(10,2)", "decimal", nil, ptr(10), ptr(2)},
		{"numeric with precision", "NUMERIC(15,4)", "numeric", nil, ptr(15), ptr(4)},
		{"timestamp with time zone", "TIMESTAMP WITH TIME ZONE", "timestamp with time zone", nil, nil, nil},
		{"timestamp without time zone", "TIMESTAMP WITHOUT TIME ZONE", "timestamp without time zone", nil, nil, nil},
		{"integer", "INTEGER", "integer", nil, nil, nil},
		{"bigint", "BIGINT", "bigint", nil, nil, nil},
		{"smallint", "SMALLINT", "smallint", nil, nil, nil},
		{"serial", "SERIAL", "serial", nil, nil, nil},
		{"bigserial", "BIGSERIAL", "bigserial", nil, nil, nil},
		{"text", "TEXT", "text", nil, nil, nil},
		{"boolean", "BOOLEAN", "boolean", nil, nil, nil},
		{"date", "DATE", "date", nil, nil, nil},
		{"time", "TIME", "time", nil, nil, nil},
		{"timestamp", "TIMESTAMP", "timestamp", nil, nil, nil},
		{"real", "REAL", "real", nil, nil, nil},
		{"double precision", "DOUBLE PRECISION", "double precision", nil, nil, nil},
		{"float4", "FLOAT4", "real", nil, nil, nil},
		{"float8", "FLOAT8", "double precision", nil, nil, nil},
		{"uuid", "UUID", "uuid", nil, nil, nil},
		{"json", "JSON", "json", nil, nil, nil},
		{"jsonb", "JSONB", "jsonb", nil, nil, nil},
		{"int", "INT", "integer", nil, nil, nil},
		{"int4", "INT4", "integer", nil, nil, nil},
		{"int8", "INT8", "bigint", nil, nil, nil},
		{"int2", "INT2", "smallint", nil, nil, nil},
		{"bool", "BOOL", "boolean", nil, nil, nil},
		{"unknown type", "CUSTOMTYPE", "CUSTOMTYPE", nil, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataType, length, precision, scale := ParseDataType(tt.typeStr)

			if dataType != tt.expectedType {
				t.Errorf("ParseDataType(%s) dataType = %s, want %s", tt.typeStr, dataType, tt.expectedType)
			}

			if tt.expectedLength != nil {
				if length == nil {
					t.Errorf("ParseDataType(%s) expected length %v, got nil", tt.typeStr, *tt.expectedLength)
				} else if *length != *tt.expectedLength {
					t.Errorf("ParseDataType(%s) length = %v, want %v", tt.typeStr, *length, *tt.expectedLength)
				}
			} else if length != nil {
				t.Errorf("ParseDataType(%s) expected nil length, got %v", tt.typeStr, *length)
			}

			if tt.expectedPrec != nil {
				if precision == nil {
					t.Errorf("ParseDataType(%s) expected precision %v, got nil", tt.typeStr, *tt.expectedPrec)
				} else if *precision != *tt.expectedPrec {
					t.Errorf("ParseDataType(%s) precision = %v, want %v", tt.typeStr, *precision, *tt.expectedPrec)
				}
			} else if precision != nil {
				t.Errorf("ParseDataType(%s) expected nil precision, got %v", tt.typeStr, *precision)
			}

			if tt.expectedScale != nil {
				if scale == nil {
					t.Errorf("ParseDataType(%s) expected scale %v, got nil", tt.typeStr, *tt.expectedScale)
				} else if *scale != *tt.expectedScale {
					t.Errorf("ParseDataType(%s) scale = %v, want %v", tt.typeStr, *scale, *tt.expectedScale)
				}
			} else if scale != nil {
				t.Errorf("ParseDataType(%s) expected nil scale, got %v", tt.typeStr, *scale)
			}
		})
	}
}

func TestParseDataType_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name         string
		typeStr      string
		expectedType string
	}{
		{"varchar lowercase", "varchar", "varchar"},
		{"VARCHAR uppercase", "VARCHAR", "VARCHAR"},
		{"Varchar Mixed", "Varchar", "Varchar"},
		{"integer lowercase", "integer", "integer"},
		{"INTEGER uppercase", "INTEGER", "integer"},
		{"Integer Mixed", "Integer", "integer"},
		{"timestamp with zone lowercase", "timestamp with time zone", "timestamp with time zone"},
		{"TIMESTAMP WITH TIME ZONE uppercase", "TIMESTAMP WITH TIME ZONE", "timestamp with time zone"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dataType, _, _, _ := ParseDataType(tt.typeStr)
			if dataType != tt.expectedType {
				t.Errorf("ParseDataType(%s) dataType = %s, want %s", tt.typeStr, dataType, tt.expectedType)
			}
		})
	}
}

func ptr(i int32) *int32 {
	return &i
}
