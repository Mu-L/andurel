package types

import (
	"testing"

	"github.com/mbvlabs/andurel/generator/internal/catalog"
)

func TestNewPostgreSQLTypeMapper(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()
	if mapper == nil {
		t.Fatal("NewPostgreSQLTypeMapper returned nil")
	}
	if mapper.GetDatabaseType() != "postgresql" {
		t.Errorf("GetDatabaseType() = %s, want postgresql", mapper.GetDatabaseType())
	}
}

func TestPostgreSQLMapper_MapToGoType_UUID(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	column := catalog.NewColumn("id", "uuid")
	column.SetNotNull()

	result, err := mapper.MapToGoType(column)
	if err != nil {
		t.Fatalf("MapToGoType failed: %v", err)
	}

	if result.GoType != "uuid.UUID" {
		t.Errorf("GoType = %s, want uuid.UUID", result.GoType)
	}
	if result.SQLCType != "uuid.UUID" {
		t.Errorf("SQLCType = %s, want uuid.UUID", result.SQLCType)
	}
	if result.PackageName != "github.com/google/uuid" {
		t.Errorf("PackageName = %s, want github.com/google/uuid", result.PackageName)
	}
}

func TestPostgreSQLMapper_MapToGoType_NonNullableTypes(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name        string
		dataType    string
		expectedGo  string
		expectedSQL string
		expectedPkg string
	}{
		{"varchar", "varchar", "string", "string", ""},
		{"text", "text", "string", "string", ""},
		{"char", "char", "string", "string", ""},
		{"boolean", "boolean", "bool", "bool", ""},
		{"bool", "bool", "bool", "bool", ""},
		{"integer", "integer", "int32", "int32", ""},
		{"int", "int", "int32", "int32", ""},
		{"int4", "int4", "int32", "int32", ""},
		{"serial", "serial", "int32", "int32", ""},
		{"bigint", "bigint", "int64", "int64", ""},
		{"int8", "int8", "int64", "int64", ""},
		{"bigserial", "bigserial", "int64", "int64", ""},
		{"smallint", "smallint", "int16", "int16", ""},
		{"int2", "int2", "int16", "int16", ""},
		{"smallserial", "smallserial", "int16", "int16", ""},
		{"real", "real", "float32", "float32", ""},
		{"float4", "float4", "float32", "float32", ""},
		{"double precision", "double precision", "float64", "float64", ""},
		{"float8", "float8", "float64", "float64", ""},
		{"bytea", "bytea", "[]byte", "[]byte", ""},
		{"jsonb", "jsonb", "[]byte", "pgtype.JSONB", "github.com/jackc/pgx/v5/pgtype"},
		{"json", "json", "[]byte", "pgtype.JSON", "github.com/jackc/pgx/v5/pgtype"},
		{"inet", "inet", "string", "pgtype.Inet", "github.com/jackc/pgx/v5/pgtype"},
		{"cidr", "cidr", "string", "pgtype.CIDR", "github.com/jackc/pgx/v5/pgtype"},
		{"macaddr", "macaddr", "string", "pgtype.Macaddr", "github.com/jackc/pgx/v5/pgtype"},
		{"macaddr8", "macaddr8", "string", "pgtype.Macaddr8", "github.com/jackc/pgx/v5/pgtype"},
		{"point", "point", "string", "pgtype.point", "github.com/jackc/pgx/v5/pgtype"},
		{"money", "money", "string", "pgtype.Money", "github.com/jackc/pgx/v5/pgtype"},
		{"xml", "xml", "string", "string", ""},
		{"tsvector", "tsvector", "string", "string", ""},
		{"tsquery", "tsquery", "string", "string", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := catalog.NewColumn("test_col", tt.dataType).SetNotNull()

			result, err := mapper.MapToGoType(column)
			if err != nil {
				t.Fatalf("MapToGoType failed for %s: %v", tt.dataType, err)
			}

			if result.GoType != tt.expectedGo {
				t.Errorf("GoType = %s, want %s", result.GoType, tt.expectedGo)
			}
			if result.SQLCType != tt.expectedSQL {
				t.Errorf("SQLCType = %s, want %s", result.SQLCType, tt.expectedSQL)
			}
			if result.PackageName != tt.expectedPkg {
				t.Errorf("PackageName = %s, want %s", result.PackageName, tt.expectedPkg)
			}
		})
	}
}

func TestPostgreSQLMapper_MapToGoType_NullableTypes(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name        string
		dataType    string
		expectedGo  string
		expectedSQL string
		expectedPkg string
	}{
		{"varchar nullable", "varchar", "string", "pgtype.Text", "github.com/jackc/pgx/v5/pgtype"},
		{"text nullable", "text", "string", "pgtype.Text", "github.com/jackc/pgx/v5/pgtype"},
		{"char nullable", "char", "string", "pgtype.Text", "github.com/jackc/pgx/v5/pgtype"},
		{"boolean nullable", "boolean", "bool", "pgtype.Bool", "github.com/jackc/pgx/v5/pgtype"},
		{"bool nullable", "bool", "bool", "pgtype.Bool", "github.com/jackc/pgx/v5/pgtype"},
		{"integer nullable", "integer", "int32", "pgtype.Int4", "github.com/jackc/pgx/v5/pgtype"},
		{"bigint nullable", "bigint", "int64", "pgtype.Int8", "github.com/jackc/pgx/v5/pgtype"},
		{"smallint nullable", "smallint", "int16", "pgtype.Int2", "github.com/jackc/pgx/v5/pgtype"},
		{"real nullable", "real", "float32", "pgtype.Float4", "github.com/jackc/pgx/v5/pgtype"},
		{"double precision nullable", "double precision", "float64", "pgtype.Float8", "github.com/jackc/pgx/v5/pgtype"},
		{"decimal nullable", "decimal", "float64", "pgtype.Numeric", "github.com/jackc/pgx/v5/pgtype"},
		{"numeric nullable", "numeric", "float64", "pgtype.Numeric", "github.com/jackc/pgx/v5/pgtype"},
		{"timestamp nullable", "timestamp", "time.Time", "pgtype.Timestamp", "github.com/jackc/pgx/v5/pgtype"},
		{"timestamptz nullable", "timestamptz", "time.Time", "pgtype.Timestamptz", "github.com/jackc/pgx/v5/pgtype"},
		{"date nullable", "date", "time.Time", "pgtype.Date", "github.com/jackc/pgx/v5/pgtype"},
		{"time nullable", "time", "time.Time", "pgtype.Time", "github.com/jackc/pgx/v5/pgtype"},
		{"timetz nullable", "timetz", "time.Time", "pgtype.Timetz", "github.com/jackc/pgx/v5/pgtype"},
		{"interval nullable", "interval", "string", "pgtype.Interval", "github.com/jackc/pgx/v5/pgtype"},
		{"jsonb nullable", "jsonb", "[]byte", "pgtype.JSONB", "github.com/jackc/pgx/v5/pgtype"},
		{"json nullable", "json", "[]byte", "pgtype.JSON", "github.com/jackc/pgx/v5/pgtype"},
		{"bytea nullable", "bytea", "[]byte", "pgtype.Bytea", "github.com/jackc/pgx/v5/pgtype"},
		{"inet nullable", "inet", "string", "pgtype.Inet", "github.com/jackc/pgx/v5/pgtype"},
		{"money nullable", "money", "string", "pgtype.Money", "github.com/jackc/pgx/v5/pgtype"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := catalog.NewColumn("test_col", tt.dataType)

			result, err := mapper.MapToGoType(column)
			if err != nil {
				t.Fatalf("MapToGoType failed for %s: %v", tt.dataType, err)
			}

			if result.GoType != tt.expectedGo {
				t.Errorf("GoType = %s, want %s", result.GoType, tt.expectedGo)
			}
			if result.SQLCType != tt.expectedSQL {
				t.Errorf("SQLCType = %s, want %s", result.SQLCType, tt.expectedSQL)
			}
			if result.PackageName != tt.expectedPkg {
				t.Errorf("PackageName = %s, want %s", result.PackageName, tt.expectedPkg)
			}
		})
	}
}

func TestPostgreSQLMapper_MapToGoType_ArrayTypes(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name        string
		dataType    string
		expectedGo  string
		expectedSQL string
		expectedPkg string
	}{
		{"_integer array", "_integer", "[]int32", "pgtype.Array[int32]", "github.com/jackc/pgx/v5/pgtype"},
		{"_integer nullable array", "_integer", "[]int32", "pgtype.Array[int32]", "github.com/jackc/pgx/v5/pgtype"},
		{"_text array", "_text", "[]string", "pgtype.Array[string]", "github.com/jackc/pgx/v5/pgtype"},
		{"_text nullable array", "_text", "[]string", "pgtype.Array[string]", "github.com/jackc/pgx/v5/pgtype"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := catalog.NewColumn("test_col", tt.dataType)

			result, err := mapper.MapToGoType(column)
			if err != nil {
				t.Fatalf("MapToGoType failed for %s: %v", tt.dataType, err)
			}

			if result.GoType != tt.expectedGo {
				t.Errorf("GoType = %s, want %s", result.GoType, tt.expectedGo)
			}
			if result.SQLCType != tt.expectedSQL {
				t.Errorf("SQLCType = %s, want %s", result.SQLCType, tt.expectedSQL)
			}
			if result.PackageName != tt.expectedPkg {
				t.Errorf("PackageName = %s, want %s", result.PackageName, tt.expectedPkg)
			}
		})
	}
}

func TestPostgreSQLMapper_MapToGoType_UnknownTypes(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	column := catalog.NewColumn("unknown_col", "custom_type")

	result, err := mapper.MapToGoType(column)
	if err != nil {
		t.Fatalf("MapToGoType failed for unknown type: %v", err)
	}

	if result.GoType != "interface{}" {
		t.Errorf("Unknown type GoType = %s, want interface{}", result.GoType)
	}
	if result.SQLCType != "interface{}" {
		t.Errorf("Unknown type SQLCType = %s, want interface{}", result.SQLCType)
	}
	if result.PackageName != "" {
		t.Errorf("Unknown type PackageName = %s, want empty string", result.PackageName)
	}
}

func TestPostgreSQLMapper_MapToSQLType(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name         string
		goType       string
		nullable     bool
		expectedSQL  string
		expectErr    bool
	}{
		{"string non-nullable", "string", false, "string", false},
		{"string nullable", "string", true, "pgtype.Text", false},
		{"int32 non-nullable", "int32", false, "int32", false},
		{"int32 nullable", "int32", true, "pgtype.Int4", false},
		{"int64 non-nullable", "int64", false, "int64", false},
		{"int64 nullable", "int64", true, "pgtype.Int8", false},
		{"bool non-nullable", "bool", false, "bool", false},
		{"bool nullable", "bool", true, "pgtype.Bool", false},
		{"float32 non-nullable", "float32", false, "float32", false},
		{"float32 nullable", "float32", true, "pgtype.Float4", false},
		{"float64 non-nullable", "float64", false, "float64", false},
		{"float64 nullable", "float64", true, "pgtype.Float8", false},
		{"time.Time non-nullable", "time.Time", false, "pgtype.Timestamptz", false},
		{"time.Time nullable", "time.Time", true, "pgtype.Timestamptz", false},
		{"[]byte non-nullable", "[]byte", false, "[]byte", false},
		{"[]byte nullable", "[]byte", true, "pgtype.Bytea", false},
		{"unknown type", "custom.Type", false, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mapper.MapToSQLType(tt.goType, tt.nullable)

			if tt.expectErr {
				if err == nil {
					t.Errorf("Expected error for %s, got nil", tt.goType)
				}
				return
			}

			if err != nil {
				t.Fatalf("MapToSQLType failed for %s: %v", tt.goType, err)
			}

			if result != tt.expectedSQL {
				t.Errorf("Result = %s, want %s", result, tt.expectedSQL)
			}
		})
	}
}

func TestPostgreSQLMapper_GenerateConversionFromDB(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name         string
		fieldName    string
		sqlcType     string
		goType       string
		expectedCode string
	}{
		{"pgtype.Text", "Name", "pgtype.Text", "string", "row.Name.String"},
		{"pgtype.Bool", "IsActive", "pgtype.Bool", "bool", "row.IsActive.Bool"},
		{"pgtype.Int2", "SmallAge", "pgtype.Int2", "int16", "row.SmallAge.Int16"},
		{"pgtype.Int4", "Age", "pgtype.Int4", "int32", "row.Age.Int32"},
		{"pgtype.Int8", "BigAge", "pgtype.Int8", "int64", "row.BigAge.Int64"},
		{"pgtype.Float4", "SmallPrice", "pgtype.Float4", "float32", "row.SmallPrice.Float32"},
		{"pgtype.Float8", "Price", "pgtype.Float8", "float64", "row.Price.Float64"},
		{"pgtype.Timestamptz", "CreatedAt", "pgtype.Timestamptz", "time.Time", "row.CreatedAt.Time"},
		{"pgtype.Timestamp", "UpdatedAt", "pgtype.Timestamp", "time.Time", "row.UpdatedAt.Time"},
		{"pgtype.Date", "BirthDate", "pgtype.Date", "time.Time", "row.BirthDate.Time"},
		{"pgtype.Time", "StartTime", "pgtype.Time", "time.Time", "row.StartTime.Time"},
		{"pgtype.Timetz", "EndTimetz", "pgtype.Timetz", "time.Time", "row.EndTimetz.Time"},
		{"pgtype.Interval", "Duration", "pgtype.Interval", "string", "row.Duration.Microseconds"},
		{"pgtype.JSONB", "Metadata", "pgtype.JSONB", "[]byte", "row.Metadata.Bytes"},
		{"pgtype.JSON", "Data", "pgtype.JSON", "[]byte", "row.Data.Bytes"},
		{"pgtype.Inet", "IPAddress", "pgtype.Inet", "string", "row.IPAddress.IPNet.String()"},
		{"pgtype.CIDR", "Network", "pgtype.CIDR", "string", "row.Network.IPNet.String()"},
		{"pgtype.Money", "Amount", "pgtype.Money", "string", "row.Amount.String"},
		{"pgtype.Bit", "BitField", "pgtype.Bit", "string", "string(row.BitField.Bytes)"},
		{"pgtype.Array[int32]", "Tags", "pgtype.Array[int32]", "[]int32", "row.Tags.Elements"},
		{"pgtype.Array[string]", "Names", "pgtype.Array[string]", "[]string", "row.Names.Elements"},
		{"direct field", "ID", "string", "string", "row.ID"},
		{"pgtype.Point", "Location", "pgtype.Point", "string", "string(row.Location.Bytes)"},
		{"pgtype.Int4range", "AgeRange", "pgtype.Int4range", "string", "string(row.AgeRange.Bytes)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.GenerateConversionFromDB(tt.fieldName, tt.sqlcType, tt.goType)
			if result != tt.expectedCode {
				t.Errorf("GenerateConversionFromDB(%s, %s, %s) = %s, want %s",
					tt.fieldName, tt.sqlcType, tt.goType, result, tt.expectedCode)
			}
		})
	}
}

func TestPostgreSQLMapper_GenerateConversionToDB(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name         string
		sqlcType     string
		goType       string
		valueExpr    string
		expectedCode string
	}{
		{"pgtype.Text", "pgtype.Text", "string", "data.Name", "pgtype.Text{String: data.Name, Valid: true}"},
		{"pgtype.Bool", "pgtype.Bool", "bool", "data.IsActive", "pgtype.Bool{Bool: data.IsActive, Valid: true}"},
		{"pgtype.Int2", "pgtype.Int2", "int16", "data.SmallAge", "pgtype.Int2{Int16: data.SmallAge, Valid: true}"},
		{"pgtype.Int4", "pgtype.Int4", "int32", "data.Age", "pgtype.Int4{Int32: data.Age, Valid: true}"},
		{"pgtype.Int8", "pgtype.Int8", "int64", "data.BigAge", "pgtype.Int8{Int64: data.BigAge, Valid: true}"},
		{"pgtype.Float4", "pgtype.Float4", "float32", "data.SmallPrice", "pgtype.Float4{Float32: data.SmallPrice, Valid: true}"},
		{"pgtype.Float8", "pgtype.Float8", "float64", "data.Price", "pgtype.Float8{Float64: data.Price, Valid: true}"},
		{"pgtype.Timestamptz", "pgtype.Timestamptz", "time.Time", "data.CreatedAt", "pgtype.Timestamptz{Time: data.CreatedAt, Valid: true}"},
		{"pgtype.Timestamp", "pgtype.Timestamp", "time.Time", "data.UpdatedAt", "pgtype.Timestamp{Time: data.UpdatedAt, Valid: true}"},
		{"pgtype.Date", "pgtype.Date", "time.Time", "data.BirthDate", "pgtype.Date{Time: data.BirthDate, Valid: true}"},
		{"pgtype.Time", "pgtype.Time", "time.Time", "data.StartTime", "pgtype.Time{Time: data.StartTime, Valid: true}"},
		{"pgtype.Timetz", "pgtype.Timetz", "time.Time", "data.EndTimetz", "pgtype.Timetz{Time: data.EndTimetz, Valid: true}"},
		{"pgtype.Interval", "pgtype.Interval", "string", "data.Duration", "pgtype.Interval{Microseconds: data.Duration, Valid: true}"},
		{"pgtype.JSONB", "pgtype.JSONB", "[]byte", "data.Metadata", "pgtype.JSONB{Bytes: data.Metadata, Valid: true}"},
		{"pgtype.JSON", "pgtype.JSON", "[]byte", "data.Data", "pgtype.JSON{Bytes: data.Data, Valid: true}"},
		{"pgtype.Inet", "pgtype.Inet", "string", "data.IPAddress", "pgtype.Inet{IPNet: data.IPAddress, Valid: true}"},
		{"pgtype.Money", "pgtype.Money", "string", "data.Amount", "pgtype.Money{String: data.Amount, Valid: true}"},
		{"pgtype.Array[int32]", "pgtype.Array[int32]", "[]int32", "data.Tags", "pgtype.Array[int32]{Elements: data.Tags, Valid: true}"},
		{"pgtype.Array[string]", "pgtype.Array[string]", "[]string", "data.Names", "pgtype.Array[string]{Elements: data.Names, Valid: true}"},
		{"direct value", "string", "string", "data.Name", "data.Name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.GenerateConversionToDB(tt.sqlcType, tt.goType, tt.valueExpr)
			if result != tt.expectedCode {
				t.Errorf("GenerateConversionToDB(%s, %s, %s) = %s, want %s",
					tt.sqlcType, tt.goType, tt.valueExpr, result, tt.expectedCode)
			}
		})
	}
}

func TestPostgreSQLMapper_GenerateZeroCheck(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name         string
		goType       string
		valueExpr    string
		expectedCode string
	}{
		{"uuid.UUID", "uuid.UUID", "data.ID", "data.ID != uuid.Nil"},
		{"pgtype.Text", "pgtype.Text", "data.Name", "data.Name.Valid"},
		{"pgtype.Bool", "pgtype.Bool", "data.IsActive", "data.IsActive.Valid"},
		{"pgtype.Int2", "pgtype.Int2", "data.SmallAge", "data.SmallAge.Valid"},
		{"pgtype.Int4", "pgtype.Int4", "data.Age", "data.Age.Valid"},
		{"pgtype.Int8", "pgtype.Int8", "data.Count", "data.Count.Valid"},
		{"pgtype.Float4", "pgtype.Float4", "data.SmallPrice", "data.SmallPrice.Valid"},
		{"pgtype.Float8", "pgtype.Float8", "data.Amount", "data.Amount.Valid"},
		{"pgtype.Timestamptz", "pgtype.Timestamptz", "data.CreatedAt", "data.CreatedAt.Valid"},
		{"pgtype.JSONB", "pgtype.JSONB", "data.Metadata", "data.Metadata.Valid"},
		{"interface{}", "interface{}", "data.Unknown", "true"},
		{"string", "string", "data.Name", "true"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapper.GenerateZeroCheck(tt.goType, tt.valueExpr)
			if result != tt.expectedCode {
				t.Errorf("GenerateZeroCheck(%s, %s) = %s, want %s",
					tt.goType, tt.valueExpr, result, tt.expectedCode)
			}
		})
	}
}

func TestPostgreSQLMapper_TimestampTypes(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name        string
		dataType    string
		nullable    bool
		expectedGo  string
		expectedSQL string
	}{
		{"timestamp", "timestamp", false, "time.Time", "pgtype.Timestamp"},
		{"timestamp without time zone", "timestamp without time zone", false, "time.Time", "pgtype.Timestamp"},
		{"timestamptz", "timestamptz", false, "time.Time", "pgtype.Timestamptz"},
		{"timestamp with time zone", "timestamp with time zone", false, "time.Time", "pgtype.Timestamptz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := catalog.NewColumn("test_col", tt.dataType)
			if !tt.nullable {
				column.SetNotNull()
			}

			result, err := mapper.MapToGoType(column)
			if err != nil {
				t.Fatalf("MapToGoType failed for %s: %v", tt.dataType, err)
			}

			if result.GoType != tt.expectedGo {
				t.Errorf("GoType = %s, want %s", result.GoType, tt.expectedGo)
			}
			if result.SQLCType != tt.expectedSQL {
				t.Errorf("SQLCType = %s, want %s", result.SQLCType, tt.expectedSQL)
			}
		})
	}
}

func TestPostgreSQLMapper_NumericTypes(t *testing.T) {
	mapper := NewPostgreSQLTypeMapper()

	tests := []struct {
		name        string
		dataType    string
		nullable    bool
		expectedGo  string
		expectedSQL string
	}{
		{"decimal", "decimal", false, "float64", "pgtype.Numeric"},
		{"decimal nullable", "decimal", true, "float64", "pgtype.Numeric"},
		{"numeric", "numeric", false, "float64", "pgtype.Numeric"},
		{"numeric nullable", "numeric", true, "float64", "pgtype.Numeric"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			column := catalog.NewColumn("test_col", tt.dataType)
			if !tt.nullable {
				column.SetNotNull()
			}

			result, err := mapper.MapToGoType(column)
			if err != nil {
				t.Fatalf("MapToGoType failed for %s: %v", tt.dataType, err)
			}

			if result.GoType != tt.expectedGo {
				t.Errorf("GoType = %s, want %s", result.GoType, tt.expectedGo)
			}
			if result.SQLCType != tt.expectedSQL {
				t.Errorf("SQLCType = %s, want %s", result.SQLCType, tt.expectedSQL)
			}
		})
	}
}
