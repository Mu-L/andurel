package catalog

import (
	"testing"
)

func TestColumn_SetUnique(t *testing.T) {
	col := &Column{Name: "test"}
	result := col.SetUnique()
	
	if !result.IsUnique {
		t.Error("expected column to be unique")
	}
	if result != col {
		t.Error("expected method chaining to return same column")
	}
}

func TestColumn_SetLength(t *testing.T) {
	col := &Column{Name: "test"}
	length := int32(255)
	result := col.SetLength(length)
	
	if result.Length == nil {
		t.Fatal("expected length to be set")
	}
	if *result.Length != length {
		t.Errorf("expected length %d, got %d", length, *result.Length)
	}
}

func TestColumn_SetPrecisionScale(t *testing.T) {
	col := &Column{Name: "test"}
	precision := int32(10)
	scale := int32(2)
	result := col.SetPrecisionScale(precision, scale)
	
	if result.Precision == nil || result.Scale == nil {
		t.Fatal("expected precision and scale to be set")
	}
	if *result.Precision != precision {
		t.Errorf("expected precision %d, got %d", precision, *result.Precision)
	}
	if *result.Scale != scale {
		t.Errorf("expected scale %d, got %d", scale, *result.Scale)
	}
}

func TestColumn_SetArray(t *testing.T) {
	col := &Column{Name: "test"}
	result := col.SetArray()
	
	if !result.IsArray {
		t.Error("expected column to be array")
	}
}

func TestColumn_SetCreatedBy(t *testing.T) {
	col := &Column{Name: "test"}
	migrationFile := "001_create_users.sql"
	result := col.SetCreatedBy(migrationFile)
	
	if result.CreatedBy != migrationFile {
		t.Errorf("expected CreatedBy to be %s, got %s", migrationFile, result.CreatedBy)
	}
}

func TestColumn_SetModifiedBy(t *testing.T) {
	col := &Column{Name: "test"}
	migrationFile := "002_alter_users.sql"
	result := col.SetModifiedBy(migrationFile)
	
	if result.ModifiedBy != migrationFile {
		t.Errorf("expected ModifiedBy to be %s, got %s", migrationFile, result.ModifiedBy)
	}
}

func TestColumn_Clone(t *testing.T) {
	length := int32(255)
	precision := int32(10)
	scale := int32(2)
	defaultVal := "test_default"
	
	original := &Column{
		Name:         "test_column",
		DataType:     "VARCHAR",
		IsNullable:   false,
		IsArray:      true,
		CreatedBy:    "001_create.sql",
		ModifiedBy:   "002_alter.sql",
		IsPrimaryKey: true,
		IsUnique:     true,
		DefaultVal:   &defaultVal,
		Length:       &length,
		Precision:    &precision,
		Scale:        &scale,
	}
	
	clone := original.Clone()
	
	// Verify all fields are copied
	if clone.Name != original.Name {
		t.Errorf("expected Name %s, got %s", original.Name, clone.Name)
	}
	if clone.DataType != original.DataType {
		t.Errorf("expected DataType %s, got %s", original.DataType, clone.DataType)
	}
	if clone.IsNullable != original.IsNullable {
		t.Error("IsNullable not copied correctly")
	}
	if clone.IsArray != original.IsArray {
		t.Error("IsArray not copied correctly")
	}
	if clone.IsPrimaryKey != original.IsPrimaryKey {
		t.Error("IsPrimaryKey not copied correctly")
	}
	if clone.IsUnique != original.IsUnique {
		t.Error("IsUnique not copied correctly")
	}
	
	// Verify pointer fields are deep copied
	if clone.DefaultVal == nil || *clone.DefaultVal != *original.DefaultVal {
		t.Error("DefaultVal not copied correctly")
	}
	if clone.Length == nil || *clone.Length != *original.Length {
		t.Error("Length not copied correctly")
	}
	if clone.Precision == nil || *clone.Precision != *original.Precision {
		t.Error("Precision not copied correctly")
	}
	if clone.Scale == nil || *clone.Scale != *original.Scale {
		t.Error("Scale not copied correctly")
	}
	
	// Verify it's a deep copy (changing clone doesn't affect original)
	clone.Name = "modified"
	if original.Name == "modified" {
		t.Error("modifying clone affected original")
	}
}

func TestColumn_CloneNilPointers(t *testing.T) {
	original := &Column{
		Name:     "test",
		DataType: "INT",
	}
	
	clone := original.Clone()
	
	if clone.DefaultVal != nil {
		t.Error("expected DefaultVal to be nil")
	}
	if clone.Length != nil {
		t.Error("expected Length to be nil")
	}
	if clone.Precision != nil {
		t.Error("expected Precision to be nil")
	}
	if clone.Scale != nil {
		t.Error("expected Scale to be nil")
	}
}
