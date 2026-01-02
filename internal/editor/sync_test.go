package editor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/secureopenvsx/internal/models"
)

func TestCopyDir(t *testing.T) {
	// Create temp directories for source and target
	srcDir, err := os.MkdirTemp("", "copydir_test_src")
	if err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	targetDir, err := os.MkdirTemp("", "copydir_test_target")
	if err != nil {
		t.Fatalf("Failed to create target dir: %v", err)
	}
	defer os.RemoveAll(targetDir)

	// Create some files in source
	testFile := filepath.Join(srcDir, "package.json")
	if err := os.WriteFile(testFile, []byte(`{"name": "test"}`), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create a subdirectory with a file
	subDir := filepath.Join(srcDir, "lib")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}
	subFile := filepath.Join(subDir, "index.js")
	if err := os.WriteFile(subFile, []byte("module.exports = {}"), 0644); err != nil {
		t.Fatalf("Failed to create sub file: %v", err)
	}

	// Test copying
	dstPath := filepath.Join(targetDir, "copied")
	err = copyDir(srcDir, dstPath)
	if err != nil {
		t.Errorf("copyDir failed: %v", err)
	}

	// Verify the files were copied
	copiedFile := filepath.Join(dstPath, "package.json")
	if _, err := os.Stat(copiedFile); os.IsNotExist(err) {
		t.Error("package.json was not copied")
	}

	copiedSubFile := filepath.Join(dstPath, "lib", "index.js")
	if _, err := os.Stat(copiedSubFile); os.IsNotExist(err) {
		t.Error("lib/index.js was not copied")
	}
}

func TestCopyDirNonExistent(t *testing.T) {
	err := copyDir("/nonexistent/source", "/nonexistent/target")
	if err == nil {
		t.Error("Expected error for non-existent source")
	}
}

func TestCopyFile(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "copyfile_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source file
	srcFile := filepath.Join(tmpDir, "source.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcFile, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy file
	dstFile := filepath.Join(tmpDir, "dest.txt")
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Errorf("copyFile failed: %v", err)
	}

	// Verify content
	copied, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(copied) != string(content) {
		t.Errorf("Content mismatch: got %s, want %s", string(copied), string(content))
	}
}

func TestDetectConflicts(t *testing.T) {
	// This test uses the actual DetectConflicts function
	// It will return empty conflicts for non-existent editors
	conflicts, err := DetectConflicts(models.EditorVSCode, models.EditorWindsurf, []string{"test.extension"})
	
	// If either editor doesn't exist, we expect either an error or empty conflicts
	if err != nil {
		// Expected if editor profiles don't exist
		t.Logf("DetectConflicts returned error (expected if editors not available): %v", err)
	} else {
		// Should return a slice (possibly empty)
		if conflicts == nil {
			t.Error("DetectConflicts returned nil slice")
		}
	}
}

func TestSyncExtensionsInvalidSource(t *testing.T) {
	request := models.SyncRequest{
		SourceEditor:  models.EditorType("nonexistent"),
		TargetEditors: []models.EditorType{models.EditorWindsurf},
		ExtensionIDs:  []string{"test.extension"},
	}

	_, err := SyncExtensions(request)
	if err == nil {
		t.Error("Expected error for invalid source editor")
	}
}

func TestSyncExtensionsEmptyTargets(t *testing.T) {
	request := models.SyncRequest{
		SourceEditor:  models.EditorVSCode,
		TargetEditors: []models.EditorType{},
		ExtensionIDs:  []string{"test.extension"},
	}

	// This might fail if VS Code is not installed, which is OK
	report, err := SyncExtensions(request)
	if err != nil {
		t.Logf("SyncExtensions returned error (expected if VS Code not available): %v", err)
		return
	}

	// With empty targets, should have empty results
	if len(report.Results) != 0 {
		t.Errorf("Expected 0 results for empty targets, got %d", len(report.Results))
	}
}

func TestSyncExtensionsEmptyExtensionIDs(t *testing.T) {
	request := models.SyncRequest{
		SourceEditor:  models.EditorVSCode,
		TargetEditors: []models.EditorType{models.EditorWindsurf},
		ExtensionIDs:  []string{},
	}

	// This might fail if VS Code is not installed, which is OK
	report, err := SyncExtensions(request)
	if err != nil {
		t.Logf("SyncExtensions returned error (expected if VS Code not available): %v", err)
		return
	}

	// With empty extension IDs, should have results but no copies
	if report.TotalCopied != 0 {
		t.Errorf("Expected 0 copies for empty extension IDs, got %d", report.TotalCopied)
	}
}
