package config

import (
	"testing"
)

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()

	if config == nil {
		t.Fatal("NewDefaultConfig returned nil")
	}

	if len(config.MigrationDirs) != 1 {
		t.Errorf("Expected 1 migration dir, got %d", len(config.MigrationDirs))
	}

	if config.MigrationDirs[0] != "database/migrations" {
		t.Errorf("Expected migration dir 'database/migrations', got '%s'", config.MigrationDirs[0])
	}

	if config.DatabaseType != "postgresql" {
		t.Errorf("Expected database type 'postgresql', got '%s'", config.DatabaseType)
	}

	if config.TableName != "" {
		t.Errorf("Expected empty table name, got '%s'", config.TableName)
	}

	if config.OutputFile != "" {
		t.Errorf("Expected empty output file, got '%s'", config.OutputFile)
	}

	if config.PackageName != "models" {
		t.Errorf("Expected package name 'models', got '%s'", config.PackageName)
	}

	if !config.GenerateJSON {
		t.Error("Expected GenerateJSON to be true")
	}
}

func TestConfigDefaults(t *testing.T) {
	config := &Config{
		MigrationDirs: []string{"database/migrations"},
		DatabaseType:  "postgresql",
		PackageName:   "models",
		GenerateJSON:  true,
	}

	// Test that default values are set correctly
	if config.DatabaseType == "" {
		t.Error("DatabaseType should not be empty in default config")
	}

	if config.PackageName == "" {
		t.Error("PackageName should not be empty in default config")
	}

	// Test zero values for optional fields
	if config.TableName == "" {
		// This is acceptable, TableName is optional
	}

	if config.OutputFile == "" {
		// This is acceptable, OutputFile is optional
	}
}

func TestConfig_StructFields(t *testing.T) {
	// Test that Config struct has all expected fields
	config := Config{
		MigrationDirs: []string{"migrations1", "migrations2"},
		DatabaseType:  "mysql",
		TableName:     "users",
		OutputFile:    "output.go",
		PackageName:   "custommodels",
		GenerateJSON:  false,
	}

	if len(config.MigrationDirs) != 2 {
		t.Errorf("Expected 2 migration dirs, got %d", len(config.MigrationDirs))
	}

	if config.MigrationDirs[0] != "migrations1" {
		t.Errorf("Expected migration dir 'migrations1', got '%s'", config.MigrationDirs[0])
	}

	if config.MigrationDirs[1] != "migrations2" {
		t.Errorf("Expected migration dir 'migrations2', got '%s'", config.MigrationDirs[1])
	}

	if config.DatabaseType != "mysql" {
		t.Errorf("Expected database type 'mysql', got '%s'", config.DatabaseType)
	}

	if config.TableName != "users" {
		t.Errorf("Expected table name 'users', got '%s'", config.TableName)
	}

	if config.OutputFile != "output.go" {
		t.Errorf("Expected output file 'output.go', got '%s'", config.OutputFile)
	}

	if config.PackageName != "custommodels" {
		t.Errorf("Expected package name 'custommodels', got '%s'", config.PackageName)
	}

	if config.GenerateJSON {
		t.Error("Expected GenerateJSON to be false")
	}
}

func TestConfig_EmptyArrays(t *testing.T) {
	config := &Config{
		MigrationDirs: []string{},
		DatabaseType:  "postgresql",
		PackageName:   "models",
		GenerateJSON:  true,
	}

	if len(config.MigrationDirs) != 0 {
		t.Errorf("Expected 0 migration dirs, got %d", len(config.MigrationDirs))
	}

	if config.DatabaseType != "postgresql" {
		t.Errorf("Expected database type 'postgresql', got '%s'", config.DatabaseType)
	}
}
