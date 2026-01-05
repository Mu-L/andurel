package cmds

import (
	"path/filepath"
	"testing"
)

func TestCapitalize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "Hello"},
		{"HELLO", "HELLO"},
		{"h", "H"},
		{"darwin", "Darwin"},
		{"linux", "Linux"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := capitalize(tt.input)
			if got != tt.expected {
				t.Errorf("capitalize(%q) = %q, expected %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestMapArch(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"amd64", "x86_64"},
		{"arm64", "arm64"},
		{"386", "i386"},
		{"other", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := mapArch(tt.input)
			if got != tt.expected {
				t.Errorf("mapArch(%q) = %q, expected %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestExtractGitHubRepo(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"github.com/a-h/templ/cmd/templ", "a-h/templ"},
		{"github.com/sqlc-dev/sqlc", "sqlc-dev/sqlc"},
		{"a-h/templ", "a-h/templ"},
		{"simple", "simple"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := extractGitHubRepo(tt.input)
			if got != tt.expected {
				t.Errorf("extractGitHubRepo(%q) = %q, expected %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGetPlatform(t *testing.T) {
	goos, goarch := GetPlatform()
	
	if goos == "" {
		t.Error("GetPlatform returned empty GOOS")
	}
	if goarch == "" {
		t.Error("GetPlatform returned empty GOARCH")
	}

	// Should match runtime values
	validGOOS := map[string]bool{
		"darwin": true, "linux": true, "windows": true,
		"freebsd": true, "openbsd": true, "netbsd": true,
	}
	validGOARCH := map[string]bool{
		"amd64": true, "arm64": true, "386": true, "arm": true,
	}

	if !validGOOS[goos] && !validGOARCH[goarch] {
		t.Logf("Got platform: %s/%s", goos, goarch)
	}
}

func TestToolDownloader_GetReleaseURL(t *testing.T) {
	tests := []struct {
		name        string
		downloader  *ToolDownloader
		goos        string
		goarch      string
		wantArchive string
		wantErr     bool
	}{
		{
			name: "templ linux amd64",
			downloader: &ToolDownloader{
				Name:    "templ",
				Module:  "github.com/a-h/templ/cmd/templ",
				Version: "v0.2.543",
			},
			goos:        "linux",
			goarch:      "amd64",
			wantArchive: "tar.gz",
			wantErr:     false,
		},
		{
			name: "templ windows amd64",
			downloader: &ToolDownloader{
				Name:    "templ",
				Module:  "github.com/a-h/templ/cmd/templ",
				Version: "v0.2.543",
			},
			goos:        "windows",
			goarch:      "amd64",
			wantArchive: "zip",
			wantErr:     false,
		},
		{
			name: "sqlc",
			downloader: &ToolDownloader{
				Name:    "sqlc",
				Module:  "github.com/sqlc-dev/sqlc/cmd/sqlc",
				Version: "v1.25.0",
			},
			goos:        "linux",
			goarch:      "amd64",
			wantArchive: "tar.gz",
			wantErr:     false,
		},
		{
			name: "goose binary",
			downloader: &ToolDownloader{
				Name:    "goose",
				Module:  "github.com/pressly/goose/v3/cmd/goose",
				Version: "v3.18.0",
			},
			goos:        "linux",
			goarch:      "amd64",
			wantArchive: "binary",
			wantErr:     false,
		},
		{
			name: "air",
			downloader: &ToolDownloader{
				Name:    "air",
				Module:  "github.com/cosmtrek/air",
				Version: "v1.49.0",
			},
			goos:        "linux",
			goarch:      "amd64",
			wantArchive: "tar.gz",
			wantErr:     false,
		},
		{
			name: "mailpit",
			downloader: &ToolDownloader{
				Name:    "mailpit",
				Module:  "github.com/axllent/mailpit",
				Version: "v1.13.3",
			},
			goos:        "linux",
			goarch:      "amd64",
			wantArchive: "tar.gz",
			wantErr:     false,
		},
		{
			name: "usql",
			downloader: &ToolDownloader{
				Name:    "usql",
				Module:  "github.com/xo/usql",
				Version: "v0.14.10",
			},
			goos:        "linux",
			goarch:      "amd64",
			wantArchive: "tar.bz2",
			wantErr:     false,
		},
		{
			name: "dblab",
			downloader: &ToolDownloader{
				Name:    "dblab",
				Module:  "github.com/danvergara/dblab",
				Version: "v0.22.0",
			},
			goos:        "linux",
			goarch:      "amd64",
			wantArchive: "tar.gz",
			wantErr:     false,
		},
		{
			name: "unknown tool",
			downloader: &ToolDownloader{
				Name:    "unknown",
				Module:  "github.com/example/unknown",
				Version: "v1.0.0",
			},
			goos:        "linux",
			goarch:      "amd64",
			wantArchive: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, archiveType, err := tt.downloader.getReleaseURL(tt.goos, tt.goarch)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("getReleaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !tt.wantErr {
				if url == "" {
					t.Error("getReleaseURL() returned empty URL")
				}
				if archiveType != tt.wantArchive {
					t.Errorf("getReleaseURL() archiveType = %v, want %v", archiveType, tt.wantArchive)
				}
			}
		})
	}
}

func TestExtractBinaryUnsupportedType(t *testing.T) {
	err := extractBinary("test.xyz", "binary", "dest", "unsupported")
	if err == nil {
		t.Error("expected error for unsupported archive type")
	}
}

func TestCopyFileNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	
	err := copyFile("/nonexistent/file", filepath.Join(tmpDir, "dest"))
	if err == nil {
		t.Error("expected error for nonexistent source file")
	}
}
