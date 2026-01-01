package editor

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/yourusername/secureopenvsx/internal/models"
)

// GetDefaultEditorProfiles returns all known editor profiles with default paths
func GetDefaultEditorProfiles() []models.EditorProfile {
	homeDir, _ := os.UserHomeDir()

	profiles := []models.EditorProfile{
		{
			ID:             models.EditorVSCode,
			Name:           "VS Code",
			ExtensionsDir:  getDefaultExtensionsDir(homeDir, ".vscode"),
			IndexFile:      "extensions.json",
			CLICommand:     "code",
			IsVSCodeFamily: true,
			IsCustom:       false,
		},
		{
			ID:             models.EditorVSCodeInsiders,
			Name:           "VS Code Insiders",
			ExtensionsDir:  getDefaultExtensionsDir(homeDir, ".vscode-insiders"),
			IndexFile:      "extensions.json",
			CLICommand:     "code-insiders",
			IsVSCodeFamily: true,
			IsCustom:       false,
		},
		{
			ID:             models.EditorVSCodium,
			Name:           "VSCodium",
			ExtensionsDir:  getDefaultExtensionsDir(homeDir, ".vscode-oss"),
			IndexFile:      "extensions.json",
			CLICommand:     "codium",
			IsVSCodeFamily: true,
			IsCustom:       false,
		},
		{
			ID:             models.EditorWindsurf,
			Name:           "Windsurf",
			ExtensionsDir:  getDefaultExtensionsDir(homeDir, ".windsurf"),
			IndexFile:      "extensions.json",
			CLICommand:     "",
			IsVSCodeFamily: false,
			IsCustom:       false,
		},
		{
			ID:             models.EditorCursor,
			Name:           "Cursor",
			ExtensionsDir:  getDefaultExtensionsDir(homeDir, ".cursor"),
			IndexFile:      "extensions.json",
			CLICommand:     "",
			IsVSCodeFamily: false,
			IsCustom:       false,
		},
		{
			ID:             models.EditorKiro,
			Name:           "Kiro",
			ExtensionsDir:  getDefaultExtensionsDir(homeDir, ".kiro"),
			IndexFile:      "extensions.json",
			CLICommand:     "",
			IsVSCodeFamily: false,
			IsCustom:       false,
		},
	}

	return profiles
}

// getDefaultExtensionsDir returns the default extensions directory for an editor
func getDefaultExtensionsDir(homeDir, editorFolder string) string {
	return filepath.Join(homeDir, editorFolder, "extensions")
}

// GetEditorProfile returns a specific editor profile by ID
func GetEditorProfile(editorType models.EditorType) (models.EditorProfile, error) {
	profiles := GetDefaultEditorProfiles()
	for _, p := range profiles {
		if p.ID == editorType {
			return p, nil
		}
	}
	return models.EditorProfile{}, fmt.Errorf("unknown editor type: %s", editorType)
}

// CheckEditorStatus checks the availability and status of an editor
func CheckEditorStatus(profile models.EditorProfile) models.EditorStatus {
	status := models.EditorStatus{
		Editor:      profile,
		IsAvailable: false,
	}

	// Check if extensions directory exists
	if info, err := os.Stat(profile.ExtensionsDir); err == nil && info.IsDir() {
		status.DirExists = true
	} else {
		status.DisabledReason = fmt.Sprintf("Extensions directory not found: %s", profile.ExtensionsDir)
		return status
	}

	// Check if index file exists
	indexPath := filepath.Join(profile.ExtensionsDir, profile.IndexFile)
	if _, err := os.Stat(indexPath); err == nil {
		status.IndexFileExists = true
	} else {
		// Try fallback to extension.json (singular)
		fallbackPath := filepath.Join(profile.ExtensionsDir, "extension.json")
		if _, err := os.Stat(fallbackPath); err == nil {
			status.IndexFileExists = true
		}
	}

	// Check CLI availability for VS Code family
	if profile.IsVSCodeFamily && profile.CLICommand != "" {
		cliPath, err := findCLI(profile.CLICommand)
		if err == nil {
			status.CLIAvailable = true
			status.CLIPath = cliPath
		}
	}

	// Count extensions
	entries, err := os.ReadDir(profile.ExtensionsDir)
	if err == nil {
		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				status.ExtensionCount++
			}
		}
	}

	status.IsAvailable = status.DirExists
	return status
}

// findCLI finds the CLI command in PATH
func findCLI(command string) (string, error) {
	// On Windows, also try .cmd extension
	if runtime.GOOS == "windows" {
		if path, err := exec.LookPath(command + ".cmd"); err == nil {
			return path, nil
		}
	}
	return exec.LookPath(command)
}

// GetCLIStatus checks the availability of all VS Code family CLI tools
func GetCLIStatus() models.CLIStatus {
	status := models.CLIStatus{}

	// Check VS Code
	if path, err := findCLI("code"); err == nil {
		status.VSCodeAvailable = true
		status.VSCodePath = path
		status.AnyAvailable = true
		if status.PreferredCLI == "" {
			status.PreferredCLI = "code"
		}
	}

	// Check VS Code Insiders
	if path, err := findCLI("code-insiders"); err == nil {
		status.InsidersAvailable = true
		status.InsidersPath = path
		status.AnyAvailable = true
		if status.PreferredCLI == "" {
			status.PreferredCLI = "code-insiders"
		}
	}

	// Check VSCodium
	if path, err := findCLI("codium"); err == nil {
		status.CodiumAvailable = true
		status.CodiumPath = path
		status.AnyAvailable = true
		if status.PreferredCLI == "" {
			status.PreferredCLI = "codium"
		}
	}

	return status
}

// ReadExtensionsIndex reads and parses the extensions.json file
func ReadExtensionsIndex(extensionsDir string) ([]models.ExtensionIndexEntry, error) {
	// Try extensions.json first
	indexPath := filepath.Join(extensionsDir, "extensions.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		// Try fallback to extension.json
		indexPath = filepath.Join(extensionsDir, "extension.json")
		data, err = os.ReadFile(indexPath)
		if err != nil {
			return nil, fmt.Errorf("no extensions index file found in %s", extensionsDir)
		}
	}

	var entries []models.ExtensionIndexEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse extensions index: %w", err)
	}

	return entries, nil
}

// WriteExtensionsIndex writes the extensions.json file
func WriteExtensionsIndex(extensionsDir string, entries []models.ExtensionIndexEntry) error {
	indexPath := filepath.Join(extensionsDir, "extensions.json")
	
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal extensions index: %w", err)
	}

	// Compact to single line like VS Code does (optional, but matches real format)
	var compactData []byte
	compactData, err = json.Marshal(entries)
	if err != nil {
		compactData = data // fallback to indented
	}

	if err := os.WriteFile(indexPath, compactData, 0644); err != nil {
		return fmt.Errorf("failed to write extensions index: %w", err)
	}

	return nil
}

// FindExtensionEntry finds an extension entry by ID in the index
func FindExtensionEntry(entries []models.ExtensionIndexEntry, extensionID string) *models.ExtensionIndexEntry {
	extensionID = strings.ToLower(extensionID)
	for i, entry := range entries {
		if strings.ToLower(entry.Identifier.ID) == extensionID {
			return &entries[i]
		}
	}
	return nil
}

// InstallExtensionViaCLI installs an extension using the VS Code CLI
func InstallExtensionViaCLI(cliCommand string, extensionID string) error {
	cliPath, err := findCLI(cliCommand)
	if err != nil {
		return fmt.Errorf("CLI command '%s' not found: %w", cliCommand, err)
	}

	cmd := exec.Command(cliPath, "--install-extension", extensionID)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install extension: %s - %w", string(output), err)
	}

	return nil
}

// InstallMultipleExtensionsViaCLI installs multiple extensions using the VS Code CLI
func InstallMultipleExtensionsViaCLI(cliCommand string, extensionIDs []string) models.InstallReport {
	report := models.InstallReport{
		CLIUsed: cliCommand,
		Results: make([]models.InstallResult, 0, len(extensionIDs)),
	}

	for _, extID := range extensionIDs {
		result := models.InstallResult{
			ExtensionID: extID,
		}

		err := InstallExtensionViaCLI(cliCommand, extID)
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			report.TotalFailed++
		} else {
			result.Success = true
			result.Message = "Successfully installed"
			report.TotalSuccess++
		}

		report.Results = append(report.Results, result)
	}

	return report
}
