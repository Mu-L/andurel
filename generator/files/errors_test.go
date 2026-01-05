package files

import (
	"errors"
	"os"
	"testing"
)

func TestFileOperationError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *FileOperationError
		expected string
	}{
		{
			name: "error without output",
			err: &FileOperationError{
				Operation: "write",
				Path:      "/tmp/test.txt",
				Err:       os.ErrPermission,
			},
			expected: "file operation 'write' failed for path '/tmp/test.txt': permission denied",
		},
		{
			name: "error with output",
			err: &FileOperationError{
				Operation: "sqlc_compile",
				Path:      "/project/root",
				Err:       errors.New("compilation error"),
				Output:    "syntax error on line 5",
			},
			expected: "file operation 'sqlc_compile' failed for path '/project/root': compilation error\nOutput: syntax error on line 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.expected {
				t.Errorf("Error() = %q, expected %q", got, tt.expected)
			}
		})
	}
}

func TestFileOperationError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	err := &FileOperationError{
		Operation: "test",
		Path:      "/test",
		Err:       originalErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("Unwrap() = %v, expected %v", unwrapped, originalErr)
	}
}

func TestFileOperationError_ErrorsIs(t *testing.T) {
	err := &FileOperationError{
		Operation: "read",
		Path:      "/test",
		Err:       os.ErrNotExist,
	}

	if !errors.Is(err, os.ErrNotExist) {
		t.Error("errors.Is should return true for wrapped os.ErrNotExist")
	}
}
