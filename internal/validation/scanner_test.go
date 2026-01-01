package validation

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestNewScanner(t *testing.T) {
	scanner := NewScanner()
	if scanner == nil {
		t.Fatal("NewScanner returned nil")
	}
	if scanner.validator == nil {
		t.Error("validator is nil")
	}
}

func TestGetExtensionsPath(t *testing.T) {
	path, err := GetExtensionsPath()
	
	if err != nil {
		t.Logf("GetExtensionsPath returned error (expected if VS Code not installed): %v", err)
		return
	}

	if path == "" {
		t.Error("GetExtensionsPath returned empty path")
	}
}

func TestScanInstalledExtensions(t *testing.T) {
	scanner := NewScanner()
	tempDir := t.TempDir()

	// Create mock extension
	extDir := filepath.Join(tempDir, "publisher.name-1.0.0")
	err := os.MkdirAll(extDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create extension directory: %v", err)
	}

	packageJSON := map[string]interface{}{
		"publisher": "publisher",
		"name":      "name",
		"version":   "1.0.0",
	}

	data, _ := json.MarshalIndent(packageJSON, "", "  ")
	os.WriteFile(filepath.Join(extDir, "package.json"), data, 0644)

	extensions, err := scanner.ScanInstalledExtensions(tempDir)
	if err != nil {
		t.Fatalf("ScanInstalledExtensions failed: %v", err)
	}

	if len(extensions) != 1 {
		t.Errorf("Expected 1 extension, got %d", len(extensions))
	}

	if len(extensions) > 0 {
		ext := extensions[0]
		if ext.ID != "publisher.name" {
			t.Errorf("ID = %s, want publisher.name", ext.ID)
		}
	}
}

func TestScanInstalledExtensionsEmptyDirectory(t *testing.T) {
	scanner := NewScanner()
	tempDir := t.TempDir()

	extensions, err := scanner.ScanInstalledExtensions(tempDir)
	if err != nil {
		t.Fatalf("Failed: %v", err)
	}

	if len(extensions) != 0 {
		t.Errorf("Expected 0 extensions, got %d", len(extensions))
	}
}

func TestAuditExtensionsEmptyDirectory(t *testing.T) {
	scanner := NewScanner()
	tempDir := t.TempDir()

	report, err := scanner.AuditExtensions(tempDir)
	if err != nil {
		t.Fatalf("AuditExtensions failed: %v", err)
	}

	if report.TotalExtensions != 0 {
		t.Errorf("TotalExtensions = %d, want 0", report.TotalExtensions)
	}

	if report.AuditTime.IsZero() {
		t.Error("AuditTime should be set")
	}
}
