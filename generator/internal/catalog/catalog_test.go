package catalog

import (
	"testing"
)

func TestNewCatalog(t *testing.T) {
	cat := NewCatalog("public")
	if cat == nil {
		t.Fatal("NewCatalog returned nil")
	}
	if cat.DefaultSchema != "public" {
		t.Errorf("Expected default schema 'public', got '%s'", cat.DefaultSchema)
	}
}

func TestCreateSchema(t *testing.T) {
	cat := NewCatalog("public")

	schema, err := cat.CreateSchema("test_schema")
	if err != nil {
		t.Fatalf("CreateSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("CreateSchema returned nil schema")
	}

	if schema.Name != "test_schema" {
		t.Errorf("Expected schema name 'test_schema', got '%s'", schema.Name)
	}
}

func TestCreateSchema_Duplicate(t *testing.T) {
	cat := NewCatalog("public")

	_, err := cat.CreateSchema("test_schema")
	if err != nil {
		t.Fatalf("First CreateSchema failed: %v", err)
	}

	_, err = cat.CreateSchema("test_schema")
	if err == nil {
		t.Error("Expected error for duplicate schema, got nil")
	}
}

func TestGetSchema(t *testing.T) {
	cat := NewCatalog("public")

	_, err := cat.CreateSchema("my_schema")
	if err != nil {
		t.Fatalf("CreateSchema failed: %v", err)
	}

	schema, err := cat.GetSchema("my_schema")
	if err != nil {
		t.Fatalf("GetSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("GetSchema returned nil")
	}

	if schema.Name != "my_schema" {
		t.Errorf("Expected schema name 'my_schema', got '%s'", schema.Name)
	}
}

func TestGetSchema_NotFound(t *testing.T) {
	cat := NewCatalog("public")

	_, err := cat.GetSchema("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent schema, got nil")
	}
}

func TestAddTable(t *testing.T) {
	cat := NewCatalog("public")

	table := NewTable("public", "test_table")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(NewColumn("name", "TEXT").SetNotNull())

	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	schema, err := cat.GetSchema("public")
	if err != nil {
		t.Fatalf("GetSchema failed: %v", err)
	}

	if len(schema.Tables) != 1 {
		t.Errorf("Expected 1 table, got %d", len(schema.Tables))
	}

	var testTable *Table
	for _, tbl := range schema.Tables {
		testTable = tbl
		break
	}
	if testTable.Name != "test_table" {
		t.Errorf("Expected table name 'test_table', got '%s'", testTable.Name)
	}
}

func TestAddTable_WithNonExistentSchema(t *testing.T) {
	cat := NewCatalog("public")

	table := NewTable("nonexistent", "test_table")

	err := cat.AddTable("nonexistent", table)
	if err == nil {
		t.Error("Expected error for non-existent schema, got nil")
	}
}

func TestGetTable(t *testing.T) {
	cat := NewCatalog("public")

	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(NewColumn("email", "TEXT").SetNotNull())

	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	retrievedTable, err := cat.GetTable("public", "users")
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}

	if retrievedTable == nil {
		t.Fatal("GetTable returned nil")
	}

	if retrievedTable.Name != "users" {
		t.Errorf("Expected table name 'users', got '%s'", retrievedTable.Name)
	}

	if len(retrievedTable.Columns) != 2 {
		t.Errorf("Expected 2 columns, got %d", len(retrievedTable.Columns))
	}
}

func TestGetTable_NotFound(t *testing.T) {
	cat := NewCatalog("public")

	_, err := cat.GetTable("public", "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent table, got nil")
	}
}

func TestDropTable(t *testing.T) {
	cat := NewCatalog("public")

	table := NewTable("public", "temp_table")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	err = cat.DropTable("public", "temp_table")
	if err != nil {
		t.Fatalf("DropTable failed: %v", err)
	}

	schema, err := cat.GetSchema("public")
	if err != nil {
		t.Fatalf("GetSchema failed: %v", err)
	}

	if len(schema.Tables) != 0 {
		t.Errorf("Expected 0 tables after drop, got %d", len(schema.Tables))
	}
}

func TestRenameTable(t *testing.T) {
	cat := NewCatalog("public")

	table := NewTable("public", "old_name")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	err = cat.RenameTable("public", "old_name", "new_name")
	if err != nil {
		t.Fatalf("RenameTable failed: %v", err)
	}

	retrievedTable, err := cat.GetTable("public", "new_name")
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}

	if retrievedTable.Name != "new_name" {
		t.Errorf("Expected table name 'new_name', got '%s'", retrievedTable.Name)
	}
}

func TestListTables(t *testing.T) {
	cat := NewCatalog("public")

	table1 := NewTable("public", "users")
	table1.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	err := cat.AddTable("public", table1)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	table2 := NewTable("public", "posts")
	table2.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	err = cat.AddTable("public", table2)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	tables, err := cat.ListTables("public")
	if err != nil {
		t.Fatalf("ListTables failed: %v", err)
	}

	if len(tables) != 2 {
		t.Errorf("Expected 2 tables, got %d", len(tables))
	}
}

func TestAddEnum(t *testing.T) {
	cat := NewCatalog("public")

	enum := &Enum{
		Name: "status_type",
		Values: []string{"active", "inactive", "pending"},
	}

	err := cat.AddEnum("public", enum)
	if err != nil {
		t.Fatalf("AddEnum failed: %v", err)
	}

	schema, err := cat.GetSchema("public")
	if err != nil {
		t.Fatalf("GetSchema failed: %v", err)
	}

	if len(schema.Enums) != 1 {
		t.Errorf("Expected 1 enum, got %d", len(schema.Enums))
	}

	var testEnum *Enum
	for _, en := range schema.Enums {
		testEnum = en
		break
	}
	if testEnum.Name != "status_type" {
		t.Errorf("Expected enum name 'status_type', got '%s'", testEnum.Name)
	}
}

func TestAlterTable(t *testing.T) {
	cat := NewCatalog("public")

	table := NewTable("public", "test_table")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	err := cat.AddTable("public", table)
	if err != nil {
		t.Fatalf("AddTable failed: %v", err)
	}

	newColumn := NewColumn("new_field", "TEXT")
	alteration := TableAlteration{
		Type:   AddColumn,
		Column: newColumn,
	}

	err = cat.AlterTable("public", "test_table", alteration)
	if err != nil {
		t.Fatalf("AlterTable failed: %v", err)
	}

	retrievedTable, err := cat.GetTable("public", "test_table")
	if err != nil {
		t.Fatalf("GetTable failed: %v", err)
	}

	foundNewColumn := false
	for _, col := range retrievedTable.Columns {
		if col.Name == "new_field" {
			foundNewColumn = true
			break
		}
	}

	if !foundNewColumn {
		t.Error("Expected to find new_field column")
	}
}

func TestTable_Clone(t *testing.T) {
	original := NewTable("public", "users")
	original.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	original.AddColumn(NewColumn("email", "TEXT").SetNotNull())
	original.AddColumn(NewColumn("name", "TEXT"))

	cloned := original.Clone()

	if cloned.Name != original.Name {
		t.Errorf("Cloned table name mismatch: expected '%s', got '%s'", original.Name, cloned.Name)
	}

	if len(cloned.Columns) != len(original.Columns) {
		t.Errorf("Cloned columns count mismatch: expected %d, got %d", len(original.Columns), len(cloned.Columns))
	}

	// Verify the columns have same names
	for i := range original.Columns {
		if cloned.Columns[i].Name != original.Columns[i].Name {
			t.Errorf("Column %d name mismatch: expected '%s', got '%s'", i, original.Columns[i].Name, cloned.Columns[i].Name)
		}
	}
}

func TestTable_AddColumn(t *testing.T) {
	table := NewTable("public", "test_table")

	col := NewColumn("id", "UUID").SetPrimaryKey()
	err := table.AddColumn(col)
	if err != nil {
		t.Fatalf("AddColumn failed: %v", err)
	}

	if len(table.Columns) != 1 {
		t.Errorf("Expected 1 column, got %d", len(table.Columns))
	}

	if table.Columns[0].Name != "id" {
		t.Errorf("Expected column name 'id', got '%s'", table.Columns[0].Name)
	}
}

func TestTable_AddColumn_Duplicate(t *testing.T) {
	table := NewTable("public", "test_table")

	col := NewColumn("id", "UUID").SetPrimaryKey()
	err := table.AddColumn(col)
	if err != nil {
		t.Fatalf("First AddColumn failed: %v", err)
	}

	err = table.AddColumn(col)
	if err == nil {
		t.Error("Expected error for duplicate column, got nil")
	}
}

func TestTable_GetColumn(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(NewColumn("email", "TEXT").SetNotNull())
	table.AddColumn(NewColumn("name", "TEXT"))

	col, err := table.GetColumn("email")
	if err != nil {
		t.Fatalf("GetColumn failed: %v", err)
	}

	if col == nil {
		t.Fatal("GetColumn returned nil")
	}

	if col.Name != "email" {
		t.Errorf("Expected column name 'email', got '%s'", col.Name)
	}
}

func TestTable_GetColumn_NotFound(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())

	_, err := table.GetColumn("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent column, got nil")
	}
}

func TestTable_DropColumn(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(NewColumn("email", "TEXT").SetNotNull())
	table.AddColumn(NewColumn("name", "TEXT"))

	err := table.DropColumn("name")
	if err != nil {
		t.Fatalf("DropColumn failed: %v", err)
	}

	if len(table.Columns) != 2 {
		t.Errorf("Expected 2 columns after drop, got %d", len(table.Columns))
	}

	_, err = table.GetColumn("name")
	if err == nil {
		t.Error("Expected error when getting dropped column, got nil")
	}
}

func TestTable_ModifyColumn(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(NewColumn("name", "TEXT"))

	newColumn := NewColumn("name", "TEXT").SetNotNull().SetDefault("unknown")
	err := table.ModifyColumn("name", newColumn)
	if err != nil {
		t.Fatalf("ModifyColumn failed: %v", err)
	}

	col, err := table.GetColumn("name")
	if err != nil {
		t.Fatalf("GetColumn failed: %v", err)
	}

	if col.IsNullable {
		t.Error("Expected column to be NOT NULL after modification")
	}

	if col.DefaultVal == nil || *col.DefaultVal != "unknown" {
		t.Errorf("Expected default value 'unknown', got %v", col.DefaultVal)
	}
}

func TestTable_RenameColumn(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(NewColumn("old_name", "TEXT"))

	err := table.RenameColumn("old_name", "new_name")
	if err != nil {
		t.Fatalf("RenameColumn failed: %v", err)
	}

	_, err = table.GetColumn("old_name")
	if err == nil {
		t.Error("Expected error when getting old column name, got nil")
	}

	col, err := table.GetColumn("new_name")
	if err != nil {
		t.Fatalf("GetColumn failed: %v", err)
	}

	if col.Name != "new_name" {
		t.Errorf("Expected column name 'new_name', got '%s'", col.Name)
	}
}

func TestTable_AddIndex(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(NewColumn("email", "TEXT").SetNotNull())

	index := &Index{
		Name:     "idx_users_email",
		Columns:  []string{"email"},
		IsUnique: true,
	}

	err := table.AddIndex(index)
	if err != nil {
		t.Fatalf("AddIndex failed: %v", err)
	}

	if len(table.Indexes) != 1 {
		t.Errorf("Expected 1 index, got %d", len(table.Indexes))
	}

	if table.Indexes[0].Name != "idx_users_email" {
		t.Errorf("Expected index name 'idx_users_email', got '%s'", table.Indexes[0].Name)
	}
}

func TestTable_DropIndex(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())

	index := &Index{
		Name:    "idx_users_name",
		Columns: []string{"id"},
	}
	err := table.AddIndex(index)
	if err != nil {
		t.Fatalf("AddIndex failed: %v", err)
	}

	err = table.DropIndex("idx_users_name")
	if err != nil {
		t.Fatalf("DropIndex failed: %v", err)
	}

	if len(table.Indexes) != 0 {
		t.Errorf("Expected 0 indexes after drop, got %d", len(table.Indexes))
	}
}

func TestTable_GetPrimaryKeyColumns(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())
	table.AddColumn(NewColumn("email", "TEXT").SetNotNull())
	table.AddColumn(NewColumn("created_at", "TIMESTAMP").SetDefault("now()"))

	pkColumns := table.GetPrimaryKeyColumns()

	if len(pkColumns) != 1 {
		t.Errorf("Expected 1 primary key column, got %d", len(pkColumns))
	}

	if pkColumns[0].Name != "id" {
		t.Errorf("Expected primary key column name 'id', got '%s'", pkColumns[0].Name)
	}

	if !pkColumns[0].IsPrimaryKey {
		t.Error("Expected primary key column to be marked as IsPrimaryKey")
	}
}

func TestTable_SetCreatedBy(t *testing.T) {
	table := NewTable("public", "users")
	table.AddColumn(NewColumn("id", "UUID").SetPrimaryKey())

	result := table.SetCreatedBy("001_create_users.sql")

	if result == nil {
		t.Fatal("SetCreatedBy returned nil")
	}

	if table.CreatedBy != "001_create_users.sql" {
		t.Errorf("Expected CreatedBy '001_create_users.sql', got '%s'", table.CreatedBy)
	}

	if result.CreatedBy != "001_create_users.sql" {
		t.Errorf("Expected returned table CreatedBy '001_create_users.sql', got '%s'", result.CreatedBy)
	}
}
