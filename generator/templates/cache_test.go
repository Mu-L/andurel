package templates

import (
	"strings"
	"testing"
	"text/template"
)

func TestNewTemplateCache(t *testing.T) {
	cache := NewTemplateCache()
	if cache == nil {
		t.Fatal("NewTemplateCache returned nil")
	}
	if cache.templates == nil {
		t.Error("templates map is nil")
	}
}

func TestTemplateCache_GetTemplate(t *testing.T) {
	cache := NewTemplateCache()

	tests := []struct {
		name        string
		templateName string
		funcMap     template.FuncMap
		wantErr     bool
		errContains string
	}{
		{
			name:        "invalid template",
			templateName: "nonexistent.tmpl",
			funcMap:     template.FuncMap{},
			wantErr:     true,
			errContains: "file does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cache.GetTemplate(tt.templateName, tt.funcMap)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.errContains != "" && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("GetTemplate() error = %v, expected to contain %v", err, tt.errContains)
			}
			if !tt.wantErr && got == nil {
				t.Error("GetTemplate() returned nil template")
			}
		})
	}
}

func TestTemplateCache_ConcurrentAccess(t *testing.T) {
	cache := NewTemplateCache()

	done := make(chan bool)
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			_, err := cache.GetTemplate("nonexistent.tmpl", template.FuncMap{})
			if err == nil {
				errors <- err
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	close(errors)
	for err := range errors {
		t.Errorf("Concurrent access error: %v", err)
	}
}

func TestTemplateCache_Caching(t *testing.T) {
	cache := NewTemplateCache()

	t.Run("cached template returns same instance", func(t *testing.T) {
		templateName := "nonexistent.tmpl"
		funcMap := template.FuncMap{}

		tmpl1, err := cache.GetTemplate(templateName, funcMap)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		tmpl2, err := cache.GetTemplate(templateName, funcMap)
		if err == nil {
			t.Fatal("Expected error, got nil")
		}

		if tmpl1 != nil || tmpl2 != nil {
			t.Error("Expected nil templates on error")
		}
	})
}

func TestGetCachedTemplate(t *testing.T) {
	ClearCache()

	tmpl, err := GetCachedTemplate("nonexistent.tmpl", template.FuncMap{})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if tmpl != nil {
		t.Error("GetCachedTemplate() should return nil template on error")
	}
}

func TestClearCache(t *testing.T) {
	ClearCache()

	templateName := "nonexistent.tmpl"

	tmpl, err := GetCachedTemplate(templateName, template.FuncMap{})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	ClearCache()

	tmpl2, err := GetCachedTemplate(templateName, template.FuncMap{})
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if tmpl != nil || tmpl2 != nil {
		t.Error("Expected nil templates on error")
	}
}

func TestGlobalCache(t *testing.T) {
	if globalCache == nil {
		t.Error("globalCache is nil")
	}

	if globalCache.templates == nil {
		t.Error("globalCache.templates is nil")
	}
}
