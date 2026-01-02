package editor

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/secureopenvsx/internal/models"
)

func TestGetDefaultEditorProfiles(t *testing.T) {
	profiles := GetDefaultEditorProfiles()

	if len(profiles) == 0 {
		t.Fatal("GetDefaultEditorProfiles returned empty list")
	}

	// Check that we have all expected editors
	expectedEditors := []models.EditorType{
		models.EditorVSCode,
		models.EditorVSCodeInsiders,
		models.EditorVSCodium,
		models.EditorWindsurf,
		models.EditorCursor,
		models.EditorKiro,
	}

	for _, expected := range expectedEditors {
		found := false
		for _, p := range profiles {
			if p.ID == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected editor %s not found in profiles", expected)
		}
	}
}

func TestGetDefaultEditorProfilesFields(t *testing.T) {
	profiles := GetDefaultEditorProfiles()

	for _, p := range profiles {
		t.Run(string(p.ID), func(t *testing.T) {
			if p.ID == "" {
				t.Error("Profile ID is empty")
			}
			if p.Name == "" {
				t.Error("Profile Name is empty")
			}
			if p.ExtensionsDir == "" {
				t.Error("Profile ExtensionsDir is empty")
			}
			if p.IndexFile == "" {
				t.Error("Profile IndexFile is empty")
			}
		})
	}
}

func TestGetEditorProfile(t *testing.T) {
	tests := []struct {
		name        string
		editorType  models.EditorType
		expectError bool
	}{
		{"VS Code", models.EditorVSCode, false},
		{"VS Code Insiders", models.EditorVSCodeInsiders, false},
		{"VSCodium", models.EditorVSCodium, false},
		{"Windsurf", models.EditorWindsurf, false},
		{"Cursor", models.EditorCursor, false},
		{"Kiro", models.EditorKiro, false},
		{"Unknown", models.EditorType("unknown"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := GetEditorProfile(tt.editorType)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error for unknown editor type")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if profile.ID != tt.editorType {
					t.Errorf("Profile ID = %s, want %s", profile.ID, tt.editorType)
				}
			}
		})
	}
}

func TestCheckEditorStatus(t *testing.T) {
	// Create a temp directory to simulate extensions dir
	tmpDir, err := os.MkdirTemp("", "editor_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create extensions directory
	extensionsDir := filepath.Join(tmpDir, "extensions")
	if err := os.MkdirAll(extensionsDir, 0755); err != nil {
		t.Fatalf("Failed to create extensions dir: %v", err)
	}

	// Create some fake extension directories
	for i := 0; i < 3; i++ {
		extDir := filepath.Join(extensionsDir, "test.extension"+string(rune('a'+i))+"-1.0.0")
		if err := os.MkdirAll(extDir, 0755); err != nil {
			t.Fatalf("Failed to create extension dir: %v", err)
		}
	}

	profile := models.EditorProfile{
		ID:            models.EditorVSCode,
		Name:          "Test Editor",
		ExtensionsDir: extensionsDir,
		IndexFile:     "extensions.json",
	}

	status := CheckEditorStatus(profile)

	if !status.IsAvailable {
		t.Error("Expected editor to be available")
	}
	if !status.DirExists {
		t.Error("Expected DirExists to be true")
	}
	if status.ExtensionCount != 3 {
		t.Errorf("ExtensionCount = %d, want 3", status.ExtensionCount)
	}
}

func TestCheckEditorStatusNonExistent(t *testing.T) {
	profile := models.EditorProfile{
		ID:            models.EditorVSCode,
		Name:          "Test Editor",
		ExtensionsDir: "/nonexistent/path/that/does/not/exist",
		IndexFile:     "extensions.json",
	}

	status := CheckEditorStatus(profile)

	if status.IsAvailable {
		t.Error("Expected editor to not be available")
	}
	if status.DirExists {
		t.Error("Expected DirExists to be false")
	}
	if status.DisabledReason == "" {
		t.Error("Expected DisabledReason to be set")
	}
}

func TestGetCLIStatus(t *testing.T) {
	status := GetCLIStatus()

	// Just verify it doesn't panic and returns a valid struct
	// Actual availability depends on the test environment
	if status.VSCodeAvailable && status.VSCodePath == "" {
		t.Error("VSCodePath should be set when VSCodeAvailable is true")
	}
	if status.InsidersAvailable && status.InsidersPath == "" {
		t.Error("InsidersPath should be set when InsidersAvailable is true")
	}
	if status.CodiumAvailable && status.CodiumPath == "" {
		t.Error("CodiumPath should be set when CodiumAvailable is true")
	}
}

func TestReadExtensionsIndexNonExistent(t *testing.T) {
	_, err := ReadExtensionsIndex("/nonexistent/path")
	if err == nil {
		t.Error("Expected error for non-existent path")
	}
}

func TestReadExtensionsIndex(t *testing.T) {
	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "editor_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create extensions.json
	indexContent := `[
		{
			"identifier": {"id": "test.extension1"},
			"version": "1.0.0",
			"location": {"path": "/path/to/ext1"}
		},
		{
			"identifier": {"id": "test.extension2"},
			"version": "2.0.0",
			"location": {"path": "/path/to/ext2"}
		}
	]`

	indexPath := filepath.Join(tmpDir, "extensions.json")
	if err := os.WriteFile(indexPath, []byte(indexContent), 0644); err != nil {
		t.Fatalf("Failed to write index file: %v", err)
	}

	entries, err := ReadExtensionsIndex(tmpDir)
	if err != nil {
		t.Fatalf("ReadExtensionsIndex failed: %v", err)
	}

	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}
}

func TestWriteExtensionsIndex(t *testing.T) {
	// Create a temp directory
	tmpDir, err := os.MkdirTemp("", "editor_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	entries := []models.ExtensionIndexEntry{
		{
			Identifier: models.ExtensionIdentifier{ID: "test.extension1"},
			Version:    "1.0.0",
		},
		{
			Identifier: models.ExtensionIdentifier{ID: "test.extension2"},
			Version:    "2.0.0",
		},
	}

	err = WriteExtensionsIndex(tmpDir, entries)
	if err != nil {
		t.Fatalf("WriteExtensionsIndex failed: %v", err)
	}

	// Verify the file was created
	indexPath := filepath.Join(tmpDir, "extensions.json")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("extensions.json was not created")
	}

	// Read it back and verify
	readEntries, err := ReadExtensionsIndex(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read back index: %v", err)
	}

	if len(readEntries) != len(entries) {
		t.Errorf("Read %d entries, want %d", len(readEntries), len(entries))
	}
}

func TestFindExtensionEntry(t *testing.T) {
	entries := []models.ExtensionIndexEntry{
		{Identifier: models.ExtensionIdentifier{ID: "test.extension1"}, Version: "1.0.0"},
		{Identifier: models.ExtensionIdentifier{ID: "test.extension2"}, Version: "2.0.0"},
		{Identifier: models.ExtensionIdentifier{ID: "Test.Extension3"}, Version: "3.0.0"},
	}

	tests := []struct {
		name        string
		searchID    string
		expectFound bool
		expectVer   string
	}{
		{"Exact match", "test.extension1", true, "1.0.0"},
		{"Case insensitive", "TEST.EXTENSION2", true, "2.0.0"},
		{"Mixed case in list", "test.extension3", true, "3.0.0"},
		{"Not found", "nonexistent.extension", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FindExtensionEntry(entries, tt.searchID)

			if tt.expectFound {
				if result == nil {
					t.Error("Expected to find entry")
				} else if result.Version != tt.expectVer {
					t.Errorf("Version = %s, want %s", result.Version, tt.expectVer)
				}
			} else {
				if result != nil {
					t.Error("Expected not to find entry")
				}
			}
		})
	}
}

func TestInstallExtensionViaCLINoCommand(t *testing.T) {
	err := InstallExtensionViaCLI("nonexistent-cli-command-xyz", "test.extension")
	if err == nil {
		t.Error("Expected error for non-existent CLI command")
	}
}

func TestInstallMultipleExtensionsViaCLI(t *testing.T) {
	// This will fail because the CLI doesn't exist, but tests the structure
	report := InstallMultipleExtensionsViaCLI("nonexistent-cli-xyz", []string{"ext1", "ext2"})

	if report.CLIUsed != "nonexistent-cli-xyz" {
		t.Errorf("CLIUsed = %s, want nonexistent-cli-xyz", report.CLIUsed)
	}

	if len(report.Results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(report.Results))
	}

	// All should fail because CLI doesn't exist
	if report.TotalSuccess != 0 {
		t.Errorf("TotalSuccess = %d, want 0", report.TotalSuccess)
	}
	if report.TotalFailed != 2 {
		t.Errorf("TotalFailed = %d, want 2", report.TotalFailed)
	}
}

func TestGetDefaultExtensionsDir(t *testing.T) {
	homeDir := "/home/testuser"
	editorFolder := ".vscode"

	result := getDefaultExtensionsDir(homeDir, editorFolder)
	expected := filepath.Join(homeDir, editorFolder, "extensions")

	if result != expected {
		t.Errorf("getDefaultExtensionsDir() = %s, want %s", result, expected)
	}
}
