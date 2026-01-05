package migrations

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFilename(t *testing.T) {
	tests := []struct {
		name            string
		filename        string
		expectedSeq     int
		expectedName    string
		expectError     bool
	}{
		{
			name:         "valid filename",
			filename:     "001_create_users.sql",
			expectedSeq:  1,
			expectedName: "create_users",
			expectError:  false,
		},
		{
			name:         "valid filename with double digit",
			filename:     "10_add_column.sql",
			expectedSeq:  10,
			expectedName: "add_column",
			expectError:  false,
		},
		{
			name:         "valid filename with multiple parts",
			filename:     "123_rename_table_and_add_indexes.sql",
			expectedSeq:  123,
			expectedName: "rename_table_and_add_indexes",
			expectError:  false,
		},
		{
			name:        "invalid filename - no number",
			filename:    "create_users.sql",
			expectError: true,
		},
		{
			name:        "invalid filename - no underscore",
			filename:    "001.sql",
			expectError: true,
		},
		{
			name:        "invalid filename - no extension",
			filename:    "001_create_users",
			expectError: true,
		},
		{
			name:        "invalid filename - wrong extension",
			filename:    "001_create_users.txt",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, name, err := parseFilename(tt.filename)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for filename %s, got nil", tt.filename)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for filename %s: %v", tt.filename, err)
				return
			}

			if seq != tt.expectedSeq {
				t.Errorf("Expected sequence %d, got %d", tt.expectedSeq, seq)
			}

			if name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, name)
			}
		})
	}
}

func TestIsDownMigration(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "down migration with .down. suffix",
			filename: "001_create_users.down.sql",
			expected: true,
		},
		{
			name:     "down migration with .down.sql suffix",
			filename: "001_create_users.down.sql",
			expected: true,
		},
		{
			name:     "up migration",
			filename: "001_create_users.sql",
			expected: false,
		},
		{
			name:     "migration containing 'down' in name",
			filename: "001_add_down_button.sql",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsDownMigration(tt.filename)
			if result != tt.expected {
				t.Errorf("IsDownMigration(%s) = %v, expected %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestDetectMigrationFormat(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected MigrationFormat
	}{
		{
			name:     "goose up marker",
			content:  "-- +goose Up\nCREATE TABLE users (id INT);",
			expected: Goose,
		},
		{
			name:     "goose down marker",
			content:  "-- +goose Down\nDROP TABLE users;",
			expected: Goose,
		},
		{
			name:     "both goose markers",
			content:  "-- +goose Up\nCREATE TABLE users (id INT);\n-- +goose Down\nDROP TABLE users;",
			expected: Goose,
		},
		{
			name:     "plain SQL without markers",
			content:  "CREATE TABLE users (id INT);",
			expected: Goose,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectMigrationFormat(tt.content)
			if result != tt.expected {
				t.Errorf("detectMigrationFormat() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestRemoveRollbackStatements(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		format   MigrationFormat
		expected string
	}{
		{
			name: "goose migration with up and down",
			content: `-- +goose Up
CREATE TABLE users (id INT);
-- +goose Down
DROP TABLE users;`,
			format:   Goose,
			expected: "CREATE TABLE users (id INT);",
		},
		{
			name:     "plain SQL",
			content:  "CREATE TABLE users (id INT);",
			format:   Goose,
			expected: "",
		},
		{
			name: "goose with statement markers",
			content: `-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (id INT);
-- +goose StatementEnd`,
			format:   Goose,
			expected: "CREATE TABLE users (id INT);",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveRollbackStatements(tt.content, tt.format)
			if result != tt.expected {
				t.Errorf("RemoveRollbackStatements() =\n%s\nexpected:\n%s", result, tt.expected)
			}
		})
	}
}

func TestExtractUpSQLGoose(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "simple goose up",
			content: `-- +goose Up
CREATE TABLE users (id INT);`,
			expected: "CREATE TABLE users (id INT);",
		},
		{
			name: "goose up with down",
			content: `-- +goose Up
CREATE TABLE users (id INT);
-- +goose Down
DROP TABLE users;`,
			expected: "CREATE TABLE users (id INT);",
		},
		{
			name: "goose up with statement markers",
			content: `-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (id INT);
-- +goose StatementEnd`,
			expected: "CREATE TABLE users (id INT);",
		},
		{
			name: "multiple statements",
			content: `-- +goose Up
CREATE TABLE users (id INT);
CREATE TABLE posts (id INT);`,
			expected: `CREATE TABLE users (id INT);
CREATE TABLE posts (id INT);`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractUpSQLGoose(tt.content)
			if result != tt.expected {
				t.Errorf("extractUpSQLGoose() =\n%s\nexpected:\n%s", result, tt.expected)
			}
		})
	}
}

func TestExtractDownSQLGoose(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name: "simple goose down",
			content: `-- +goose Up
CREATE TABLE users (id INT);
-- +goose Down
DROP TABLE users;`,
			expected: "DROP TABLE users;",
		},
		{
			name: "goose down only",
			content: `-- +goose Down
DROP TABLE users;`,
			expected: "DROP TABLE users;",
		},
		{
			name: "goose down with statement markers",
			content: `-- +goose Up
CREATE TABLE users (id INT);
-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd`,
			expected: "DROP TABLE users;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDownSQLGoose(tt.content)
			if result != tt.expected {
				t.Errorf("extractDownSQLGoose() =\n%s\nexpected:\n%s", result, tt.expected)
			}
		})
	}
}

func TestExtractDownSQL(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		format   MigrationFormat
		expected string
	}{
		{
			name: "goose format",
			content: `-- +goose Up
CREATE TABLE users (id INT);
-- +goose Down
DROP TABLE users;`,
			format:   Goose,
			expected: "DROP TABLE users;",
		},
		{
			name:     "goose format without down",
			content:  "-- +goose Up\nCREATE TABLE users (id INT);",
			format:   Goose,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDownSQL(tt.content, tt.format)
			if result != tt.expected {
				t.Errorf("extractDownSQL() = %q, expected %q", result, tt.expected)
			}
		})
	}
}

func TestParseStatements(t *testing.T) {
	tests := []struct {
		name     string
		sql       string
		expected  []string
	}{
		{
			name:     "single statement",
			sql:       "CREATE TABLE users (id INT);",
			expected:  []string{"CREATE TABLE users (id INT);"},
		},
		{
			name: "multiple statements",
			sql: `CREATE TABLE users (id INT);
CREATE TABLE posts (id INT);`,
			expected: []string{
				"CREATE TABLE users (id INT);",
				"CREATE TABLE posts (id INT);",
			},
		},
		{
			name: "statements with comments",
			sql: `-- This is a comment
CREATE TABLE users (id INT);
-- Another comment
DROP TABLE posts;`,
			expected: []string{
				"CREATE TABLE users (id INT);",
				"DROP TABLE posts;",
			},
		},
		{
			name: "statements with blank lines",
			sql: `CREATE TABLE users (id INT);

CREATE TABLE posts (id INT);`,
			expected: []string{
				"CREATE TABLE users (id INT);",
				"CREATE TABLE posts (id INT);",
			},
		},
		{
			name:     "statement without semicolon",
			sql:       "CREATE TABLE users (id INT)",
			expected:  []string{"CREATE TABLE users (id INT)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseStatements(tt.sql)
			if len(result) != len(tt.expected) {
				t.Errorf("parseStatements() returned %d statements, expected %d", len(result), len(tt.expected))
				t.Logf("Got: %v", result)
				t.Logf("Expected: %v", tt.expected)
				return
			}

			for i, stmt := range result {
				if stmt != tt.expected[i] {
					t.Errorf("Statement %d: got %q, expected %q", i, stmt, tt.expected[i])
				}
			}
		})
	}
}

func TestParseMigration(t *testing.T) {
	t.Run("valid migration file", func(t *testing.T) {
		tempDir := t.TempDir()
		migrationFile := filepath.Join(tempDir, "001_create_users.sql")
		content := "-- +goose Up\nCREATE TABLE users (id INT);\n-- +goose Down\nDROP TABLE users;"
		if err := os.WriteFile(migrationFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create migration file: %v", err)
		}

		migration, err := ParseMigration(migrationFile)
		if err != nil {
			t.Fatalf("ParseMigration failed: %v", err)
		}

		if migration.Sequence != 1 {
			t.Errorf("Expected sequence 1, got %d", migration.Sequence)
		}

		if migration.Name != "create_users" {
			t.Errorf("Expected name 'create_users', got %s", migration.Name)
		}

		if migration.Format != Goose {
			t.Errorf("Expected format Goose, got %v", migration.Format)
		}

		if !strings.Contains(migration.UpSQL, "CREATE TABLE") {
			t.Errorf("Expected UpSQL to contain CREATE TABLE")
		}

		if !strings.Contains(migration.DownSQL, "DROP TABLE") {
			t.Errorf("Expected DownSQL to contain DROP TABLE")
		}

		if len(migration.Statements) != 1 {
			t.Errorf("Expected 1 statement, got %d", len(migration.Statements))
		}
	})

	t.Run("nonexistent file", func(t *testing.T) {
		_, err := ParseMigration("/nonexistent/migration.sql")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
	})

	t.Run("invalid filename", func(t *testing.T) {
		tempDir := t.TempDir()
		migrationFile := filepath.Join(tempDir, "invalid.sql")
		content := "CREATE TABLE users (id INT);"
		if err := os.WriteFile(migrationFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create migration file: %v", err)
		}

		_, err := ParseMigration(migrationFile)
		if err == nil {
			t.Error("Expected error for invalid filename")
		}
	})
}

func TestDiscoverMigrations(t *testing.T) {
	t.Run("empty directory", func(t *testing.T) {
		tempDir := t.TempDir()

		migrations, err := DiscoverMigrations([]string{tempDir})
		if err != nil {
			t.Errorf("DiscoverMigrations failed: %v", err)
		}

		if len(migrations) != 0 {
			t.Errorf("Expected 0 migrations, got %d", len(migrations))
		}
	})

	t.Run("single migration", func(t *testing.T) {
		tempDir := t.TempDir()
		migrationFile := filepath.Join(tempDir, "001_create_users.sql")
		content := "-- +goose Up\nCREATE TABLE users (id INT);\n-- +goose Down\nDROP TABLE users;"
		if err := os.WriteFile(migrationFile, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create migration file: %v", err)
		}

		migrations, err := DiscoverMigrations([]string{tempDir})
		if err != nil {
			t.Errorf("DiscoverMigrations failed: %v", err)
		}

		if len(migrations) != 1 {
			t.Errorf("Expected 1 migration, got %d", len(migrations))
		}

		if migrations[0].Sequence != 1 {
			t.Errorf("Expected sequence 1, got %d", migrations[0].Sequence)
		}
	})

	t.Run("multiple migrations", func(t *testing.T) {
		tempDir := t.TempDir()
		content := "-- +goose Up\nCREATE TABLE users (id INT);\n-- +goose Down\nDROP TABLE users;"

		for i := 1; i <= 3; i++ {
			filename := filepath.Join(tempDir, "00"+string(rune('0'+i))+"_test.sql")
			if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
				t.Fatalf("Failed to create migration file: %v", err)
			}
		}

		migrations, err := DiscoverMigrations([]string{tempDir})
		if err != nil {
			t.Errorf("DiscoverMigrations failed: %v", err)
		}

		if len(migrations) != 3 {
			t.Errorf("Expected 3 migrations, got %d", len(migrations))
		}

		for i, m := range migrations {
			if m.Sequence != i+1 {
				t.Errorf("Migration %d: expected sequence %d, got %d", i, i+1, m.Sequence)
			}
		}
	})

	t.Run("skips down migrations", func(t *testing.T) {
		tempDir := t.TempDir()
		upFile := filepath.Join(tempDir, "001_create_users.sql")
		downFile := filepath.Join(tempDir, "002_remove_users.down.sql")

		if err := os.WriteFile(upFile, []byte("CREATE TABLE users;"), 0644); err != nil {
			t.Fatalf("Failed to create up file: %v", err)
		}
		if err := os.WriteFile(downFile, []byte("DROP TABLE users;"), 0644); err != nil {
			t.Fatalf("Failed to create down file: %v", err)
		}

		migrations, err := DiscoverMigrations([]string{tempDir})
		if err != nil {
			t.Errorf("DiscoverMigrations failed: %v", err)
		}

		if len(migrations) != 1 {
			t.Errorf("Expected 1 migration (down should be skipped), got %d", len(migrations))
		}

		if migrations[0].Sequence != 1 {
			t.Errorf("Expected sequence 1, got %d", migrations[0].Sequence)
		}
	})

	t.Run("skips non-sql files", func(t *testing.T) {
		tempDir := t.TempDir()
		migrationFile := filepath.Join(tempDir, "001_create_users.sql")
		textFile := filepath.Join(tempDir, "002_readme.txt")

		if err := os.WriteFile(migrationFile, []byte("CREATE TABLE users;"), 0644); err != nil {
			t.Fatalf("Failed to create migration file: %v", err)
		}
		if err := os.WriteFile(textFile, []byte("readme"), 0644); err != nil {
			t.Fatalf("Failed to create text file: %v", err)
		}

		migrations, err := DiscoverMigrations([]string{tempDir})
		if err != nil {
			t.Errorf("DiscoverMigrations failed: %v", err)
		}

		if len(migrations) != 1 {
			t.Errorf("Expected 1 migration (non-SQL should be skipped), got %d", len(migrations))
		}
	})

	t.Run("nonexistent directory", func(t *testing.T) {
		_, err := DiscoverMigrations([]string{"/nonexistent/directory"})
		if err == nil {
			t.Error("Expected error for nonexistent directory")
		}
	})
}

func TestMigrationFormatString(t *testing.T) {
	tests := []struct {
		format   MigrationFormat
		expected string
	}{
		{Goose, "goose"},
		{MigrationFormat(99), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.format.String()
			if result != tt.expected {
				t.Errorf("MigrationFormat(%d).String() = %s, expected %s", tt.format, result, tt.expected)
			}
		})
	}
}
