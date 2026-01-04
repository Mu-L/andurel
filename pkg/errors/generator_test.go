package errors

import (
	"errors"
	"testing"
)

func TestGenerationError(t *testing.T) {
	t.Run("WithFile", func(t *testing.T) {
		cause := errors.New("write failed")
		err := &GenerationError{
			Operation: "create",
			Resource:  "model",
			File:      "user.go",
			Cause:     cause,
		}

		expected := "failed to create model in user.go: write failed"
		if err.Error() != expected {
			t.Errorf("Expected %q, got %q", expected, err.Error())
		}

		if err.Unwrap() != cause {
			t.Error("Unwrap should return the cause")
		}
	})

	t.Run("WithoutFile", func(t *testing.T) {
		cause := errors.New("template not found")
		err := &GenerationError{
			Operation: "render",
			Resource:  "view",
			File:      "",
			Cause:     cause,
		}

		expected := "failed to render view: template not found"
		if err.Error() != expected {
			t.Errorf("Expected %q, got %q", expected, err.Error())
		}
	})
}

func TestNewGeneratorError(t *testing.T) {
	cause := errors.New("test error")
	err := NewGeneratorError("generate", "controller", cause)

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	genErr, ok := err.(*GenerationError)
	if !ok {
		t.Fatal("Expected GenerationError type")
	}

	if genErr.Operation != "generate" {
		t.Error("Operation should be set")
	}

	if genErr.Resource != "controller" {
		t.Error("Resource should be set")
	}

	if genErr.Cause != cause {
		t.Error("Cause should be set")
	}

	if !errors.Is(err, cause) {
		t.Error("Should unwrap to cause")
	}
}

func TestNewFileOperationError(t *testing.T) {
	cause := errors.New("permission denied")
	err := NewFileOperationError("/path/to/file.go", "write", cause)

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	genErr, ok := err.(*GenerationError)
	if !ok {
		t.Fatal("Expected GenerationError type")
	}

	if genErr.Operation != "write" {
		t.Error("Operation should be set")
	}

	if genErr.Resource != "file" {
		t.Error("Resource should be 'file'")
	}

	if genErr.File != "/path/to/file.go" {
		t.Error("File path should be set")
	}

	if genErr.Cause != cause {
		t.Error("Cause should be set")
	}

	errorStr := err.Error()
	if !contains(errorStr, "/path/to/file.go") {
		t.Error("Error should contain file path")
	}
}

func TestFileOperationError(t *testing.T) {
	cause := errors.New("disk full")
	err := &FileOperationError{
		Path:      "/var/log/app.log",
		Operation: "append",
		Cause:     cause,
	}

	expected := "failed to append file /var/log/app.log: disk full"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap should return the cause")
	}
}

func TestNewSpecificFileOperationError(t *testing.T) {
	cause := errors.New("file not found")
	err := NewSpecificFileOperationError("/etc/config.yaml", "read", cause)

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	fileErr, ok := err.(*FileOperationError)
	if !ok {
		t.Fatal("Expected FileOperationError type")
	}

	if fileErr.Path != "/etc/config.yaml" {
		t.Error("Path should be set")
	}

	if fileErr.Operation != "read" {
		t.Error("Operation should be set")
	}

	if fileErr.Cause != cause {
		t.Error("Cause should be set")
	}

	if !errors.Is(err, cause) {
		t.Error("Should unwrap to cause")
	}
}

func TestTemplateError(t *testing.T) {
	cause := errors.New("syntax error")
	err := &TemplateError{
		TemplateName: "user/show.templ",
		Operation:    "parse",
		Cause:        cause,
	}

	expected := "failed to parse template user/show.templ: syntax error"
	if err.Error() != expected {
		t.Errorf("Expected %q, got %q", expected, err.Error())
	}

	if err.Unwrap() != cause {
		t.Error("Unwrap should return the cause")
	}
}

func TestNewTemplateError(t *testing.T) {
	cause := errors.New("missing variable")
	err := NewTemplateError("layout.templ", "render", cause)

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	tmplErr, ok := err.(*TemplateError)
	if !ok {
		t.Fatal("Expected TemplateError type")
	}

	if tmplErr.TemplateName != "layout.templ" {
		t.Error("TemplateName should be set")
	}

	if tmplErr.Operation != "render" {
		t.Error("Operation should be set")
	}

	if tmplErr.Cause != cause {
		t.Error("Cause should be set")
	}

	if !errors.Is(err, cause) {
		t.Error("Should unwrap to cause")
	}
}

func TestValidationError(t *testing.T) {
	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("regex match failed")
		err := &ValidationError{
			Field:  "email",
			Value:  "invalid-email",
			Reason: "must be valid email format",
			Cause:  cause,
		}

		expected := "validation failed for email 'invalid-email': must be valid email format (regex match failed)"
		if err.Error() != expected {
			t.Errorf("Expected %q, got %q", expected, err.Error())
		}

		if err.Unwrap() != cause {
			t.Error("Unwrap should return the cause")
		}
	})

	t.Run("WithoutCause", func(t *testing.T) {
		err := &ValidationError{
			Field:  "age",
			Value:  "-5",
			Reason: "must be positive",
			Cause:  nil,
		}

		expected := "validation failed for age '-5': must be positive"
		if err.Error() != expected {
			t.Errorf("Expected %q, got %q", expected, err.Error())
		}

		if err.Unwrap() != nil {
			t.Error("Unwrap should return nil")
		}
	})
}

func TestNewValidationError(t *testing.T) {
	t.Run("WithCause", func(t *testing.T) {
		cause := errors.New("length check failed")
		err := NewValidationError("username", "ab", "must be at least 3 characters", cause)

		if err == nil {
			t.Fatal("Expected non-nil error")
		}

		valErr, ok := err.(*ValidationError)
		if !ok {
			t.Fatal("Expected ValidationError type")
		}

		if valErr.Field != "username" {
			t.Error("Field should be set")
		}

		if valErr.Value != "ab" {
			t.Error("Value should be set")
		}

		if valErr.Reason != "must be at least 3 characters" {
			t.Error("Reason should be set")
		}

		if valErr.Cause != cause {
			t.Error("Cause should be set")
		}

		if !errors.Is(err, cause) {
			t.Error("Should unwrap to cause")
		}
	})

	t.Run("WithoutCause", func(t *testing.T) {
		err := NewValidationError("status", "invalid", "must be active or inactive", nil)

		if err == nil {
			t.Fatal("Expected non-nil error")
		}

		valErr, ok := err.(*ValidationError)
		if !ok {
			t.Fatal("Expected ValidationError type")
		}

		if valErr.Cause != nil {
			t.Error("Cause should be nil")
		}
	})
}

func TestDatabaseError(t *testing.T) {
	t.Run("WithTable", func(t *testing.T) {
		cause := errors.New("duplicate key")
		err := &DatabaseError{
			Operation: "insert",
			Table:     "users",
			Cause:     cause,
		}

		expected := "database insert failed for table users: duplicate key"
		if err.Error() != expected {
			t.Errorf("Expected %q, got %q", expected, err.Error())
		}

		if err.Unwrap() != cause {
			t.Error("Unwrap should return the cause")
		}
	})

	t.Run("WithoutTable", func(t *testing.T) {
		cause := errors.New("connection timeout")
		err := &DatabaseError{
			Operation: "connect",
			Table:     "",
			Cause:     cause,
		}

		expected := "database connect failed: connection timeout"
		if err.Error() != expected {
			t.Errorf("Expected %q, got %q", expected, err.Error())
		}
	})
}

func TestNewDatabaseError(t *testing.T) {
	cause := errors.New("foreign key constraint")
	err := NewDatabaseError("delete", "posts", cause)

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	dbErr, ok := err.(*DatabaseError)
	if !ok {
		t.Fatal("Expected DatabaseError type")
	}

	if dbErr.Operation != "delete" {
		t.Error("Operation should be set")
	}

	if dbErr.Table != "posts" {
		t.Error("Table should be set")
	}

	if dbErr.Cause != cause {
		t.Error("Cause should be set")
	}

	if !errors.Is(err, cause) {
		t.Error("Should unwrap to cause")
	}
}

func TestErrorTypes_Unwrapping(t *testing.T) {
	t.Run("GenerationError unwraps", func(t *testing.T) {
		cause := errors.New("original")
		err := NewGeneratorError("test", "resource", cause)

		var genErr *GenerationError
		if !errors.As(err, &genErr) {
			t.Error("Should be able to extract GenerationError")
		}

		if !errors.Is(err, cause) {
			t.Error("Should wrap the original error")
		}
	})

	t.Run("FileOperationError unwraps", func(t *testing.T) {
		cause := errors.New("original")
		err := NewSpecificFileOperationError("test.go", "read", cause)

		var fileErr *FileOperationError
		if !errors.As(err, &fileErr) {
			t.Error("Should be able to extract FileOperationError")
		}

		if !errors.Is(err, cause) {
			t.Error("Should wrap the original error")
		}
	})

	t.Run("TemplateError unwraps", func(t *testing.T) {
		cause := errors.New("original")
		err := NewTemplateError("test.templ", "parse", cause)

		var tmplErr *TemplateError
		if !errors.As(err, &tmplErr) {
			t.Error("Should be able to extract TemplateError")
		}

		if !errors.Is(err, cause) {
			t.Error("Should wrap the original error")
		}
	})

	t.Run("ValidationError unwraps", func(t *testing.T) {
		cause := errors.New("original")
		err := NewValidationError("field", "value", "reason", cause)

		var valErr *ValidationError
		if !errors.As(err, &valErr) {
			t.Error("Should be able to extract ValidationError")
		}

		if !errors.Is(err, cause) {
			t.Error("Should wrap the original error")
		}
	})

	t.Run("DatabaseError unwraps", func(t *testing.T) {
		cause := errors.New("original")
		err := NewDatabaseError("query", "table", cause)

		var dbErr *DatabaseError
		if !errors.As(err, &dbErr) {
			t.Error("Should be able to extract DatabaseError")
		}

		if !errors.Is(err, cause) {
			t.Error("Should wrap the original error")
		}
	})
}
