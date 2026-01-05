package validation

import (
	"testing"
)

func TestValidatePrimaryKeyDatatype(t *testing.T) {
	tests := []struct {
		name          string
		dataType      string
		databaseType  string
		migrationFile string
		columnName    string
		expectError   bool
	}{
		{
			name:          "valid UUID primary key",
			dataType:      "uuid",
			databaseType:  "postgresql",
			migrationFile: "001_create_users.sql",
			columnName:    "id",
			expectError:   false,
		},
		{
			name:          "valid UUID uppercase",
			dataType:      "UUID",
			databaseType:  "postgresql",
			migrationFile: "001_create_users.sql",
			columnName:    "id",
			expectError:   false,
		},
		{
			name:          "valid UUID mixed case",
			dataType:      "Uuid",
			databaseType:  "postgresql",
			migrationFile: "001_create_users.sql",
			columnName:    "id",
			expectError:   false,
		},
		{
			name:          "invalid integer primary key",
			dataType:      "integer",
			databaseType:  "postgresql",
			migrationFile: "001_create_users.sql",
			columnName:    "id",
			expectError:   true,
		},
		{
			name:          "invalid serial primary key",
			dataType:      "serial",
			databaseType:  "postgresql",
			migrationFile: "001_create_posts.sql",
			columnName:    "post_id",
			expectError:   true,
		},
		{
			name:          "invalid bigint primary key",
			dataType:      "bigint",
			databaseType:  "postgresql",
			migrationFile: "001_create_comments.sql",
			columnName:    "comment_id",
			expectError:   true,
		},
		{
			name:          "invalid varchar primary key",
			dataType:      "varchar(255)",
			databaseType:  "postgresql",
			migrationFile: "001_create_sessions.sql",
			columnName:    "session_id",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePrimaryKeyDatatype(tt.dataType, tt.databaseType, tt.migrationFile, tt.columnName)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for data type %s", tt.dataType)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error for data type %s: %v", tt.dataType, err)
			}

			if err != nil && tt.expectError {
				// Verify error message contains expected information
				errorMsg := err.Error()
				if !contains(errorMsg, tt.migrationFile) {
					t.Error("Error message should contain migration file name")
				}
				if !contains(errorMsg, tt.columnName) {
					t.Error("Error message should contain column name")
				}
				if !contains(errorMsg, tt.dataType) {
					t.Error("Error message should contain data type")
				}
				if !contains(errorMsg, "uuid") {
					t.Error("Error message should mention uuid")
				}
			}
		})
	}
}

func TestValidatePrimaryKeyDatatype_ErrorFormat(t *testing.T) {
	dataType := "integer"
	databaseType := "postgresql"
	migrationFile := "database/migrations/001_create_users.sql"
	columnName := "user_id"

	err := ValidatePrimaryKeyDatatype(dataType, databaseType, migrationFile, columnName)
	if err == nil {
		t.Fatal("Expected error for non-uuid data type")
	}

	errorMsg := err.Error()

	// Check that error message contains all required information
	expectedParts := []string{
		"001_create_users.sql",
		"user_id",
		"integer",
		"uuid",
		"PRIMARY KEY",
	}

	for _, part := range expectedParts {
		if !contains(errorMsg, part) {
			t.Errorf("Error message should contain %q, got:\n%s", part, errorMsg)
		}
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
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
