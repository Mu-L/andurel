package layout

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestMigrationGeneration tests the generation of migration files
func TestMigrationGeneration(t *testing.T) {
	// Set test mode to get deterministic timestamps
	os.Setenv("ANDUREL_TEST_MODE", "true")
	defer os.Unsetenv("ANDUREL_TEST_MODE")

	tempDir := t.TempDir()

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

	lastTime, err := processMigrations(tempDir, data)
	if err != nil {
		t.Fatalf("processMigrations() error = %v", err)
	}

	// Verify lastTime is not zero
	if lastTime.IsZero() {
		t.Error("processMigrations() returned zero time")
	}

	// Read all migration files and compare them to golden files
	migrationsDir := filepath.Join(tempDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Expected migrations
	expectedMigrations := []string{
		"create_river_migration_table",
		"create_river_job_and_leader_tables",
		"alter_river_job_tags",
		"alter_river_job_args_metadata_add_queue",
		"add_river_job_unique_key_and_clients",
		"add_river_job_unique_states",
		"create_users_table",
		"create_tokens_table",
	}

	if len(entries) != len(expectedMigrations) {
		t.Errorf("Expected %d migration files, got %d", len(expectedMigrations), len(entries))
	}

	// Test each migration file
	for i, entry := range entries {
		if entry.IsDir() {
			continue
		}

		t.Run(entry.Name(), func(t *testing.T) {
			// Read the migration file
			content, err := os.ReadFile(filepath.Join(migrationsDir, entry.Name()))
			if err != nil {
				t.Fatalf("Failed to read migration file: %v", err)
			}

			// Extract the migration name (everything after the timestamp)
			parts := strings.SplitN(entry.Name(), "_", 2)
			if len(parts) != 2 {
				t.Fatalf("Invalid migration filename format: %s", entry.Name())
			}
			migrationName := parts[1] // includes .sql extension

			// Compare with golden file
			goldenFile := "migrations/" + strings.TrimSuffix(migrationName, ".sql") + ".golden"
			compareGolden(t, goldenFile, string(content))
		})

		// Verify migration names match expected order
		if i < len(expectedMigrations) {
			if !strings.Contains(entry.Name(), expectedMigrations[i]) {
				t.Errorf("Migration %d: expected name to contain %s, got %s",
					i, expectedMigrations[i], entry.Name())
			}
		}
	}
}

// TestMigrationTimestamps tests that migrations have sequential timestamps
func TestMigrationTimestamps(t *testing.T) {
	// Set test mode to get deterministic timestamps
	os.Setenv("ANDUREL_TEST_MODE", "true")
	defer os.Unsetenv("ANDUREL_TEST_MODE")

	tempDir := t.TempDir()

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

	_, err := processMigrations(tempDir, data)
	if err != nil {
		t.Fatalf("processMigrations() error = %v", err)
	}

	migrationsDir := filepath.Join(tempDir, "database", "migrations")
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		t.Fatalf("Failed to read migrations directory: %v", err)
	}

	// Verify timestamps are sequential
	var lastTimestamp time.Time
	for i, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Extract timestamp from filename (first 14 characters)
		timestampStr := entry.Name()[:14]
		timestamp, err := time.Parse("20060102150405", timestampStr)
		if err != nil {
			t.Errorf("Migration %d (%s): invalid timestamp format: %v", i, entry.Name(), err)
			continue
		}

		// Verify timestamp is after the last one
		if i > 0 && !timestamp.After(lastTimestamp) {
			t.Errorf("Migration %d (%s): timestamp %s is not after previous timestamp %s",
				i, entry.Name(), timestamp, lastTimestamp)
		}

		lastTimestamp = timestamp
	}
}

// TestProcessMigrationsReturnsCorrectTime tests the return value of processMigrations
func TestProcessMigrationsReturnsCorrectTime(t *testing.T) {
	// Set test mode to get deterministic timestamps
	os.Setenv("ANDUREL_TEST_MODE", "true")
	defer os.Unsetenv("ANDUREL_TEST_MODE")

	tempDir := t.TempDir()

	data := &TemplateData{
		AppName:     "testapp",
		ModuleName:  "testapp",
		Database:    "postgresql",
		blueprint:   initializeBaseBlueprint("testapp"),
		RunToolVersion: "v1.0.0",
	}

	lastTime, err := processMigrations(tempDir, data)
	if err != nil {
		t.Fatalf("processMigrations() error = %v", err)
	}

	// Verify lastTime is 1 second after the last migration
	// In test mode, base time is 2025-01-01 00:00:00
	// Last migration (create_tokens_table) is at base + 7 seconds
	// So lastTime should be base + 8 seconds
	expectedTime := time.Date(2025, 1, 1, 0, 0, 8, 0, time.UTC)
	if !lastTime.Equal(expectedTime) {
		t.Errorf("processMigrations() returned time %s, want %s", lastTime, expectedTime)
	}
}

// TestMigrationContent_RiverQueue tests the content of River queue migrations
func TestMigrationContent_RiverQueue(t *testing.T) {
	os.Setenv("ANDUREL_TEST_MODE", "true")
	defer os.Unsetenv("ANDUREL_TEST_MODE")

	tempDir := t.TempDir()
	data := &TemplateData{
		AppName:     "testapp",
		ModuleName:  "testapp",
		Database:    "postgresql",
		blueprint:   initializeBaseBlueprint("testapp"),
		RunToolVersion: "v1.0.0",
	}

	_, err := processMigrations(tempDir, data)
	if err != nil {
		t.Fatalf("processMigrations() error = %v", err)
	}

	tests := []struct {
		name           string
		migrationName  string
		expectedInSQL  []string
	}{
		{
			name:          "river migration table",
			migrationName: "create_river_migration_table",
			expectedInSQL: []string{"CREATE TABLE", "river_migration", "version", "created_at"},
		},
		{
			name:          "river job table",
			migrationName: "create_river_job_and_leader_tables",
			expectedInSQL: []string{"CREATE TABLE", "river_job", "river_leader", "state"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrationsDir := filepath.Join(tempDir, "database", "migrations")
			entries, err := os.ReadDir(migrationsDir)
			if err != nil {
				t.Fatalf("Failed to read migrations directory: %v", err)
			}

			// Find the migration file
			var migrationFile string
			for _, entry := range entries {
				if strings.Contains(entry.Name(), tt.migrationName) {
					migrationFile = filepath.Join(migrationsDir, entry.Name())
					break
				}
			}

			if migrationFile == "" {
				t.Fatalf("Migration file not found for: %s", tt.migrationName)
			}

			content, err := os.ReadFile(migrationFile)
			if err != nil {
				t.Fatalf("Failed to read migration file: %v", err)
			}

			// Verify expected SQL statements are present
			contentStr := string(content)
			for _, expected := range tt.expectedInSQL {
				if !strings.Contains(contentStr, expected) {
					t.Errorf("Migration %s missing expected SQL: %s", tt.migrationName, expected)
				}
			}
		})
	}
}

// TestMigrationContent_Auth tests the content of auth migrations
func TestMigrationContent_Auth(t *testing.T) {
	os.Setenv("ANDUREL_TEST_MODE", "true")
	defer os.Unsetenv("ANDUREL_TEST_MODE")

	tempDir := t.TempDir()
	data := &TemplateData{
		AppName:     "testapp",
		ModuleName:  "testapp",
		Database:    "postgresql",
		blueprint:   initializeBaseBlueprint("testapp"),
		RunToolVersion: "v1.0.0",
	}

	_, err := processMigrations(tempDir, data)
	if err != nil {
		t.Fatalf("processMigrations() error = %v", err)
	}

	tests := []struct {
		name           string
		migrationName  string
		expectedInSQL  []string
	}{
		{
			name:          "users table",
			migrationName: "create_users_table",
			expectedInSQL: []string{"CREATE TABLE", "users", "email", "password", "is_admin"},
		},
		{
			name:          "tokens table",
			migrationName: "create_tokens_table",
			expectedInSQL: []string{"CREATE TABLE", "tokens", "expires_at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			migrationsDir := filepath.Join(tempDir, "database", "migrations")
			entries, err := os.ReadDir(migrationsDir)
			if err != nil {
				t.Fatalf("Failed to read migrations directory: %v", err)
			}

			// Find the migration file
			var migrationFile string
			for _, entry := range entries {
				if strings.Contains(entry.Name(), tt.migrationName) {
					migrationFile = filepath.Join(migrationsDir, entry.Name())
					break
				}
			}

			if migrationFile == "" {
				t.Fatalf("Migration file not found for: %s", tt.migrationName)
			}

			content, err := os.ReadFile(migrationFile)
			if err != nil {
				t.Fatalf("Failed to read migration file: %v", err)
			}

			// Verify expected SQL statements are present
			contentStr := string(content)
			for _, expected := range tt.expectedInSQL {
				if !strings.Contains(contentStr, expected) {
					t.Errorf("Migration %s missing expected SQL: %s", tt.migrationName, expected)
				}
			}
		})
	}
}
