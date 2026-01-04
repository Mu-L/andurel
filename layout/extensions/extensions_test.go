package extensions

import (
	"testing"
)

func TestRegister(t *testing.T) {
	tests := []struct {
		name    string
		ext     Extension
		wantErr bool
	}{
		{
			name:    "register valid extension",
			ext:     Docker{},
			wantErr: false,
		},
		{
			name:    "register nil extension",
			ext:     nil,
			wantErr: true,
		},
		{
			name: "register extension with empty name",
			ext: &testEmptyNameExtension{},
			wantErr: true,
		},
		{
			name:    "register duplicate extension",
			ext:     Docker{},
			wantErr: false, // duplicate registration should succeed (idempotent)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Register(tt.ext)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Register() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Register() unexpected error: %v", err)
				}
			}
		})
	}
}

func TestGet(t *testing.T) {
	// Ensure docker is registered
	_ = Register(Docker{})

	tests := []struct {
		name      string
		extName   string
		wantFound bool
	}{
		{
			name:      "get existing extension",
			extName:   "docker",
			wantFound: true,
		},
		{
			name:      "get non-existing extension",
			extName:   "nonexistent",
			wantFound: false,
		},
		{
			name:      "get empty name",
			extName:   "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ext, found := Get(tt.extName)

			if found != tt.wantFound {
				t.Errorf("Get() found = %v, want %v", found, tt.wantFound)
			}

			if found && ext == nil {
				t.Error("Get() found extension but returned nil")
			}
		})
	}
}

func TestNames(t *testing.T) {
	// Register test extensions
	_ = Register(Docker{})
	_ = Register(AwsSes{})
	_ = Register(Workflows{})
	_ = Register(Paddle{})

	names := Names()

	if len(names) == 0 {
		t.Error("Names() returned empty slice")
	}

	// Check that expected names are present
	expectedNames := map[string]bool{
		"docker":  true,
		"aws-ses": true,
		"workflows": true,
		"paddle": true,
	}

	foundNames := make(map[string]bool)
	for _, name := range names {
		foundNames[name] = true
	}

	for expected := range expectedNames {
		if !foundNames[expected] {
			t.Errorf("Names() missing expected extension: %s", expected)
		}
	}
}

func TestDocker_Name(t *testing.T) {
	ext := Docker{}
	name := ext.Name()

	if name != "docker" {
		t.Errorf("Docker.Name() = %s, want 'docker'", name)
	}
}

func TestDocker_Dependencies(t *testing.T) {
	ext := Docker{}
	deps := ext.Dependencies()

	if deps != nil {
		t.Errorf("Docker.Dependencies() = %v, want nil", deps)
	}
}

func TestDocker_Apply_NilContext(t *testing.T) {
	ext := Docker{}
	err := ext.Apply(nil)

	if err == nil {
		t.Error("Docker.Apply() with nil context should return error")
	}

	if err != nil && !contains(err.Error(), "context or data is nil") {
		t.Errorf("Docker.Apply() error = %v, want error containing 'context or data is nil'", err)
	}
}

func TestAwsSes_Name(t *testing.T) {
	ext := AwsSes{}
	name := ext.Name()

	if name != "aws-ses" {
		t.Errorf("AwsSes.Name() = %s, want 'aws-ses'", name)
	}
}

func TestAwsSes_Dependencies(t *testing.T) {
	ext := AwsSes{}
	deps := ext.Dependencies()

	if deps != nil {
		t.Errorf("AwsSes.Dependencies() = %v, want nil", deps)
	}
}

func TestAwsSes_Apply_NilContext(t *testing.T) {
	ext := AwsSes{}
	err := ext.Apply(nil)

	if err == nil {
		t.Error("AwsSes.Apply() with nil context should return error")
	}

	if err != nil && !contains(err.Error(), "context or data is nil") {
		t.Errorf("AwsSes.Apply() error = %v, want error containing 'context or data is nil'", err)
	}
}

func TestWorkflows_Name(t *testing.T) {
	ext := Workflows{}
	name := ext.Name()

	if name != "workflows" {
		t.Errorf("Workflows.Name() = %s, want 'workflows'", name)
	}
}

func TestWorkflows_Dependencies(t *testing.T) {
	ext := Workflows{}
	deps := ext.Dependencies()

	if deps != nil {
		t.Errorf("Workflows.Dependencies() = %v, want nil", deps)
	}
}

func TestWorkflows_Apply_NilContext(t *testing.T) {
	ext := Workflows{}
	err := ext.Apply(nil)

	if err == nil {
		t.Error("Workflows.Apply() with nil context should return error")
	}

	if err != nil && !contains(err.Error(), "context or data is nil") {
		t.Errorf("Workflows.Apply() error = %v, want error containing 'context or data is nil'", err)
	}
}

func TestPaddle_Name(t *testing.T) {
	ext := Paddle{}
	name := ext.Name()

	if name != "paddle" {
		t.Errorf("Paddle.Name() = %s, want 'paddle'", name)
	}
}

func TestPaddle_Dependencies(t *testing.T) {
	ext := Paddle{}
	deps := ext.Dependencies()

	if deps != nil {
		t.Errorf("Paddle.Dependencies() = %v, want nil", deps)
	}
}

func TestPaddle_Apply_NilContext(t *testing.T) {
	ext := Paddle{}
	err := ext.Apply(nil)

	if err == nil {
		t.Error("Paddle.Apply() with nil context should return error")
	}

	if err != nil && !contains(err.Error(), "context or data is nil") {
		t.Errorf("Paddle.Apply() error = %v, want error containing 'context or data is nil'", err)
	}
}

func TestContext_Builder_NilContext(t *testing.T) {
	var ctx *Context

	builder := ctx.Builder()
	if builder != nil {
		t.Error("Context.Builder() with nil context should return nil")
	}
}

func TestContext_Builder_NilData(t *testing.T) {
	ctx := &Context{
		Data: nil,
	}

	builder := ctx.Builder()
	if builder != nil {
		t.Error("Context.Builder() with nil data should return nil")
	}
}

// Test helper types

type testEmptyNameExtension struct{}

func (t *testEmptyNameExtension) Name() string {
	return ""
}

func (t *testEmptyNameExtension) Apply(ctx *Context) error {
	return nil
}

func (t *testEmptyNameExtension) Dependencies() []string {
	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
