package errors

import (
	"errors"
	"testing"
)

func TestErrorContext(t *testing.T) {
	t.Run("NewErrorContext", func(t *testing.T) {
		ctx := NewErrorContext("test", "resource", "file.txt")

		if ctx.Operation != "test" {
			t.Errorf("Expected operation 'test', got '%s'", ctx.Operation)
		}
		if ctx.Resource != "resource" {
			t.Errorf("Expected resource 'resource', got '%s'", ctx.Resource)
		}
		if ctx.File != "file.txt" {
			t.Errorf("Expected file 'file.txt', got '%s'", ctx.File)
		}
	})

	t.Run("WithDetail", func(t *testing.T) {
		ctx := NewErrorContext("test", "resource", "file.txt")
		ctx.WithDetail("key", "value")

		if ctx.Details["key"] != "value" {
			t.Errorf("Expected detail 'key' to be 'value', got '%v'", ctx.Details["key"])
		}
	})
}

func TestContextualError(t *testing.T) {
	t.Run("WrapError", func(t *testing.T) {
		originalErr := errors.New("original error")
		ctx := NewErrorContext("operation", "resource", "file.txt")

		wrappedErr := WrapError(originalErr, *ctx)

		if wrappedErr == nil {
			t.Fatal("Expected non-nil error")
		}

		if !errors.Is(wrappedErr, originalErr) {
			t.Error("Wrapped error should wrap the original error")
		}

		errorStr := wrappedErr.Error()
		expected := "operation: operation, resource: resource, file: file.txt: original error"
		if errorStr != expected {
			t.Errorf("Expected error string '%s', got '%s'", expected, errorStr)
		}
	})

	t.Run("WrapErrorWithNil", func(t *testing.T) {
		ctx := NewErrorContext("operation", "resource", "file.txt")
		wrappedErr := WrapError(nil, *ctx)

		if wrappedErr != nil {
			t.Error("Wrapping nil error should return nil")
		}
	})
}

func TestErrorBuilder(t *testing.T) {
	t.Run("BuildAndWrap", func(t *testing.T) {
		originalErr := errors.New("test error")

		wrappedErr := NewErrorBuilder().
			Operation("test-op").
			Resource("test-resource").
			File("test.txt").
			Detail("extra", "info").
			Wrap(originalErr)

		if wrappedErr == nil {
			t.Fatal("Expected non-nil error")
		}

		errorStr := wrappedErr.Error()
		if !contains(errorStr, "operation: test-op") {
			t.Error("Error should contain operation")
		}
		if !contains(errorStr, "resource: test-resource") {
			t.Error("Error should contain resource")
		}
		if !contains(errorStr, "file: test.txt") {
			t.Error("Error should contain file")
		}
		if !contains(errorStr, "extra: info") {
			t.Error("Error should contain details")
		}
	})

	t.Run("BuildAndNew", func(t *testing.T) {
		newErr := NewErrorBuilder().
			Operation("create").
			Resource("user").
			New("failed to create user")

		if newErr == nil {
			t.Fatal("Expected non-nil error")
		}

		errorStr := newErr.Error()
		if !contains(errorStr, "operation: create") {
			t.Error("Error should contain operation")
		}
		if !contains(errorStr, "resource: user") {
			t.Error("Error should contain resource")
		}
		if !contains(errorStr, "failed to create user") {
			t.Error("Error should contain message")
		}
	})
}

func TestConvenienceFunctions(t *testing.T) {
	t.Run("WrapFileError", func(t *testing.T) {
		originalErr := errors.New("file not found")
		wrappedErr := WrapFileError(originalErr, "read", "/path/to/file.txt")

		if wrappedErr == nil {
			t.Fatal("Expected non-nil error")
		}

		errorStr := wrappedErr.Error()
		if !contains(errorStr, "operation: read") {
			t.Error("Error should contain operation")
		}
		if !contains(errorStr, "resource: file") {
			t.Error("Error should contain resource")
		}
		if !contains(errorStr, "file: /path/to/file.txt") {
			t.Error("Error should contain file")
		}
	})

	t.Run("WrapTemplateError", func(t *testing.T) {
		originalErr := errors.New("template syntax error")
		wrappedErr := WrapTemplateError(originalErr, "parse", "template.tmpl")

		if wrappedErr == nil {
			t.Fatal("Expected non-nil error")
		}

		errorStr := wrappedErr.Error()
		if !contains(errorStr, "operation: parse") {
			t.Error("Error should contain operation")
		}
		if !contains(errorStr, "resource: template") {
			t.Error("Error should contain resource")
		}
		if !contains(errorStr, "template_name: template.tmpl") {
			t.Error("Error should contain template name")
		}
	})
}

func TestErrorRecovery(t *testing.T) {
	recovery := &DefaultErrorRecovery{}

	t.Run("RecoverableError", func(t *testing.T) {
		err := errors.New("operation timeout")

		if !recovery.CanRecover(err) {
			t.Error("Timeout error should be recoverable")
		}

		recoveredErr := recovery.Recover(err)
		if recoveredErr != nil {
			t.Error("Recovery should return nil for recoverable error")
		}
	})

	t.Run("NonRecoverableError", func(t *testing.T) {
		err := errors.New("syntax error")

		if recovery.CanRecover(err) {
			t.Error("Syntax error should not be recoverable")
		}

		recoveredErr := recovery.Recover(err)
		if recoveredErr == nil {
			t.Error("Recovery should return error for non-recoverable error")
		}
		if recoveredErr != err {
			t.Error("Recovery should return original error for non-recoverable error")
		}
	})
}

func TestIsRecoverable(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"NilError", nil, true},
		{"TimeoutError", errors.New("operation timeout"), true},
		{"TemporaryError", errors.New("temporary failure"), true},
		{"ConnectionRefused", errors.New("connection refused"), true},
		{"SyntaxError", errors.New("syntax error"), false},
		{"ValidationError", errors.New("validation failed"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsRecoverable(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for error: %v", tt.expected, result, tt.err)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestErrorContext_WithCaller(t *testing.T) {
	ctx := NewErrorContext("test", "resource", "file.txt")
	ctx.WithCaller(0)

	if ctx.Details["caller"] == nil {
		t.Error("Expected caller detail to be set")
	}

	caller, ok := ctx.Details["caller"].(string)
	if !ok {
		t.Fatal("Expected caller to be a string")
	}

	if !contains(caller, "context_test.go") {
		t.Errorf("Expected caller to contain file name, got %s", caller)
	}
}

func TestContextualError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	ctx := NewErrorContext("operation", "resource", "file.txt")
	wrappedErr := WrapError(originalErr, *ctx)

	ctxErr, ok := wrappedErr.(*ContextualError)
	if !ok {
		t.Fatal("Expected ContextualError type")
	}

	unwrapped := ctxErr.Unwrap()
	if unwrapped != originalErr {
		t.Error("Unwrap should return the original error")
	}

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is should match the original error")
	}
}

func TestContextualError_EmptyFields(t *testing.T) {
	t.Run("EmptyOperation", func(t *testing.T) {
		err := &ContextualError{
			Context: &ErrorContext{
				Operation: "",
				Resource:  "resource",
				File:      "file.txt",
				Details:   make(map[string]interface{}),
			},
			Cause: errors.New("test error"),
		}

		errorStr := err.Error()
		if contains(errorStr, "operation:") {
			t.Error("Error should not contain empty operation")
		}
		if !contains(errorStr, "resource: resource") {
			t.Error("Error should contain resource")
		}
	})

	t.Run("EmptyResource", func(t *testing.T) {
		err := &ContextualError{
			Context: &ErrorContext{
				Operation: "operation",
				Resource:  "",
				File:      "file.txt",
				Details:   make(map[string]interface{}),
			},
			Cause: errors.New("test error"),
		}

		errorStr := err.Error()
		if contains(errorStr, "resource:") {
			t.Error("Error should not contain empty resource")
		}
		if !contains(errorStr, "operation: operation") {
			t.Error("Error should contain operation")
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		err := &ContextualError{
			Context: &ErrorContext{
				Operation: "operation",
				Resource:  "resource",
				File:      "",
				Details:   make(map[string]interface{}),
			},
			Cause: errors.New("test error"),
		}

		errorStr := err.Error()
		if contains(errorStr, "file:") {
			t.Error("Error should not contain empty file")
		}
	})

	t.Run("NoCause", func(t *testing.T) {
		err := &ContextualError{
			Context: &ErrorContext{
				Operation: "operation",
				Resource:  "resource",
				File:      "file.txt",
				Details:   make(map[string]interface{}),
			},
			Cause: nil,
		}

		errorStr := err.Error()
		if !contains(errorStr, "operation: operation") {
			t.Error("Error should contain operation")
		}
		// Should not panic, should just return context
	})
}

func TestWrapErrorWithCaller(t *testing.T) {
	originalErr := errors.New("test error")
	ctx := ErrorContext{
		Operation: "test",
		Resource:  "resource",
		File:      "file.txt",
		Details:   make(map[string]interface{}),
	}

	wrappedErr := WrapErrorWithCaller(originalErr, ctx)

	if wrappedErr == nil {
		t.Fatal("Expected non-nil error")
	}

	ctxErr, ok := wrappedErr.(*ContextualError)
	if !ok {
		t.Fatal("Expected ContextualError type")
	}

	if ctxErr.Context.Details["caller"] == nil {
		t.Error("Expected caller detail to be set")
	}

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should wrap the original error")
	}
}

func TestWrapErrorWithCaller_Nil(t *testing.T) {
	ctx := ErrorContext{
		Operation: "test",
		Resource:  "resource",
		File:      "file.txt",
		Details:   make(map[string]interface{}),
	}

	wrappedErr := WrapErrorWithCaller(nil, ctx)

	if wrappedErr != nil {
		t.Error("Wrapping nil error should return nil")
	}
}

func TestNewContextualError(t *testing.T) {
	originalErr := errors.New("test error")
	err := NewContextualError("operation", "resource", "file.txt", originalErr)

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	ctxErr, ok := err.(*ContextualError)
	if !ok {
		t.Fatal("Expected ContextualError type")
	}

	if ctxErr.Context.Operation != "operation" {
		t.Error("Operation should be set")
	}

	if ctxErr.Context.Resource != "resource" {
		t.Error("Resource should be set")
	}

	if ctxErr.Context.File != "file.txt" {
		t.Error("File should be set")
	}

	if ctxErr.Context.Details["caller"] == nil {
		t.Error("Caller should be set automatically")
	}

	if ctxErr.Cause != originalErr {
		t.Error("Cause should be the original error")
	}
}

func TestErrorBuilder_WrapNil(t *testing.T) {
	wrappedErr := NewErrorBuilder().
		Operation("test").
		Wrap(nil)

	if wrappedErr != nil {
		t.Error("Wrapping nil error should return nil")
	}
}

func TestErrorBuilder_ChainedCalls(t *testing.T) {
	err := NewErrorBuilder().
		Operation("create").
		Resource("user").
		File("user.go").
		Detail("field1", "value1").
		Detail("field2", 123).
		Detail("field3", true).
		New("creation failed")

	if err == nil {
		t.Fatal("Expected non-nil error")
	}

	ctxErr, ok := err.(*ContextualError)
	if !ok {
		t.Fatal("Expected ContextualError type")
	}

	if ctxErr.Context.Operation != "create" {
		t.Error("Operation should be set")
	}

	if ctxErr.Context.Resource != "user" {
		t.Error("Resource should be set")
	}

	if ctxErr.Context.File != "user.go" {
		t.Error("File should be set")
	}

	if ctxErr.Context.Details["field1"] != "value1" {
		t.Error("Detail field1 should be set")
	}

	if ctxErr.Context.Details["field2"] != 123 {
		t.Error("Detail field2 should be set")
	}

	if ctxErr.Context.Details["field3"] != true {
		t.Error("Detail field3 should be set")
	}

	errorStr := err.Error()
	if !contains(errorStr, "creation failed") {
		t.Error("Error should contain message")
	}
}

func TestWrapValidationError(t *testing.T) {
	originalErr := errors.New("invalid format")
	wrappedErr := WrapValidationError(originalErr, "email", "not-an-email")

	if wrappedErr == nil {
		t.Fatal("Expected non-nil error")
	}

	errorStr := wrappedErr.Error()
	if !contains(errorStr, "operation: validation") {
		t.Error("Error should contain validation operation")
	}

	if !contains(errorStr, "resource: email") {
		t.Error("Error should contain field name as resource")
	}

	if !contains(errorStr, "value: not-an-email") {
		t.Error("Error should contain value in details")
	}

	if !errors.Is(wrappedErr, originalErr) {
		t.Error("Wrapped error should wrap the original error")
	}
}

func TestWrapDatabaseError(t *testing.T) {
	originalErr := errors.New("connection failed")
	wrappedErr := WrapDatabaseError(originalErr, "insert", "users")

	if wrappedErr == nil {
		t.Fatal("Expected non-nil error")
	}

	errorStr := wrappedErr.Error()
	if !contains(errorStr, "operation: insert") {
		t.Error("Error should contain operation")
	}

	if !contains(errorStr, "resource: database") {
		t.Error("Error should contain database as resource")
	}

	if !contains(errorStr, "table: users") {
		t.Error("Error should contain table in details")
	}
}

func TestWrapGenerationError(t *testing.T) {
	originalErr := errors.New("generation failed")
	wrappedErr := WrapGenerationError(originalErr, "generate", "model")

	if wrappedErr == nil {
		t.Fatal("Expected non-nil error")
	}

	errorStr := wrappedErr.Error()
	if !contains(errorStr, "operation: generate") {
		t.Error("Error should contain operation")
	}

	if !contains(errorStr, "resource: model") {
		t.Error("Error should contain resource")
	}
}

func TestIsRecoverable_ContextualError(t *testing.T) {
	t.Run("RecoverableWrapped", func(t *testing.T) {
		originalErr := errors.New("connection timeout")
		ctx := NewErrorContext("connect", "database", "")
		wrappedErr := WrapError(originalErr, *ctx)

		if !IsRecoverable(wrappedErr) {
			t.Error("Wrapped timeout error should be recoverable")
		}
	})

	t.Run("NonRecoverableWrapped", func(t *testing.T) {
		originalErr := errors.New("syntax error")
		ctx := NewErrorContext("parse", "file", "test.go")
		wrappedErr := WrapError(originalErr, *ctx)

		if IsRecoverable(wrappedErr) {
			t.Error("Wrapped syntax error should not be recoverable")
		}
	})
}

func TestErrorContext_MultipleDetails(t *testing.T) {
	ctx := NewErrorContext("test", "resource", "file.txt").
		WithDetail("detail1", "value1").
		WithDetail("detail2", 123).
		WithDetail("detail3", true)

	if len(ctx.Details) != 3 {
		t.Errorf("Expected 3 details, got %d", len(ctx.Details))
	}

	if ctx.Details["detail1"] != "value1" {
		t.Error("Detail1 should be set")
	}

	if ctx.Details["detail2"] != 123 {
		t.Error("Detail2 should be set")
	}

	if ctx.Details["detail3"] != true {
		t.Error("Detail3 should be set")
	}
}

func TestContextualError_MultipleDetails(t *testing.T) {
	ctx := NewErrorContext("operation", "resource", "file.txt").
		WithDetail("key1", "value1").
		WithDetail("key2", "value2").
		WithDetail("key3", "value3")

	err := &ContextualError{
		Context: ctx,
		Cause:   errors.New("test error"),
	}

	errorStr := err.Error()

	// All details should be in the error string
	if !contains(errorStr, "key1: value1") {
		t.Error("Error should contain key1")
	}
	if !contains(errorStr, "key2: value2") {
		t.Error("Error should contain key2")
	}
	if !contains(errorStr, "key3: value3") {
		t.Error("Error should contain key3")
	}
}
