package naming

import "testing"

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "single word", input: "User", expected: "user"},
		{name: "multi word", input: "CompanyAccount", expected: "company_account"},
		{name: "acronym handling", input: "APIKey", expected: "api_key"},
		{name: "long phrase", input: "CompanyIntelligenceReport", expected: "company_intelligence_report"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSnakeCase(tt.input)
			if got != tt.expected {
				t.Fatalf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDeriveTableName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "singular", input: "User", expected: "users"},
		{name: "already plural", input: "CompanyAccounts", expected: "company_accounts"},
		{name: "complex singular", input: "CompanyIntelligenceReport", expected: "company_intelligence_reports"},
		{name: "complex plural", input: "CompanyIntelligenceReports", expected: "company_intelligence_reports"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeriveTableName(tt.input)
			if got != tt.expected {
				t.Fatalf("DeriveTableName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "single word", input: "users", expected: "users"},
		{name: "two words", input: "admin_users", expected: "adminUsers"},
		{name: "three words", input: "product_categories", expected: "productCategories"},
		{name: "many words", input: "user_profile_settings", expected: "userProfileSettings"},
		{name: "empty string", input: "", expected: ""},
		{name: "already camelCase", input: "adminUsers", expected: "adminusers"},
		{name: "single char parts", input: "a_b_c", expected: "aBC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			if got != tt.expected {
				t.Fatalf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToLowerCamelCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "single word", input: "User", expected: "user"},
		{name: "two words", input: "NewUser", expected: "newUser"},
		{name: "three words", input: "AdminUser", expected: "adminUser"},
		{name: "many words", input: "UserProfileSettings", expected: "userProfileSettings"},
		{name: "empty string", input: "", expected: ""},
		{name: "already camelCase", input: "user", expected: "user"},
		{name: "single char", input: "U", expected: "u"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToLowerCamelCase(tt.input)
			if got != tt.expected {
				t.Fatalf("ToLowerCamelCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "lowercase word", input: "hello", expected: "Hello"},
		{name: "already capitalized", input: "Hello", expected: "Hello"},
		{name: "single lowercase char", input: "h", expected: "H"},
		{name: "single uppercase char", input: "H", expected: "H"},
		{name: "empty string", input: "", expected: ""},
		{name: "number first", input: "1hello", expected: "1hello"},
		{name: "special char first", input: "_hello", expected: "_hello"},
		{name: "mixed case", input: "hELLO", expected: "HELLO"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Capitalize(tt.input)
			if got != tt.expected {
				t.Fatalf("Capitalize(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToSnakeCase_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "empty string", input: "", expected: ""},
		{name: "single lowercase letter", input: "a", expected: "a"},
		{name: "single uppercase letter", input: "A", expected: "a"},
		{name: "numbers", input: "Test123", expected: "test123"},
		{name: "numbers between words", input: "Test123Value", expected: "test123_value"},
		{name: "consecutive uppercase", input: "HTTPServer", expected: "http_server"},
		{name: "ending with uppercase", input: "UserAPI", expected: "user_api"},
		{name: "all uppercase", input: "API", expected: "api"},
		{name: "mixed numbers and uppercase", input: "HTTP2Server", expected: "http2_server"},
		{name: "underscore in input", input: "User_Name", expected: "user_name"},
		{name: "special characters", input: "User-Name", expected: "user-name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToSnakeCase(tt.input)
			if got != tt.expected {
				t.Fatalf("ToSnakeCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestToCamelCase_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "trailing underscore", input: "user_", expected: "user"},
		{name: "leading underscore", input: "_user", expected: "User"},
		{name: "double underscore", input: "user__name", expected: "userName"},
		{name: "multiple underscores", input: "user___name", expected: "userName"},
		{name: "numbers in parts", input: "user_123_name", expected: "user123Name"},
		{name: "uppercase parts", input: "USER_NAME", expected: "userName"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToCamelCase(tt.input)
			if got != tt.expected {
				t.Fatalf("ToCamelCase(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestDeriveTableName_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "irregular plural - person", input: "Person", expected: "people"},
		{name: "irregular plural - child", input: "Child", expected: "children"},
		{name: "word ending in y", input: "Category", expected: "categories"},
		{name: "word ending in s", input: "Status", expected: "statuses"},
		{name: "word ending in x", input: "Index", expected: "indices"},
		{name: "acronym", input: "API", expected: "apis"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DeriveTableName(tt.input)
			if got != tt.expected {
				t.Fatalf("DeriveTableName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
