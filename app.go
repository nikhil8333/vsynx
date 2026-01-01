package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/yourusername/secureopenvsx/internal/editor"
	"github.com/yourusername/secureopenvsx/internal/marketplace"
	"github.com/yourusername/secureopenvsx/internal/models"
	"github.com/yourusername/secureopenvsx/internal/validation"
)

// App struct holds the application state
type App struct {
	ctx       context.Context
	validator *validation.Validator
	scanner   *validation.Scanner
}

// NewApp creates a new App instance
func NewApp() *App {
	return &App{
		validator: validation.NewValidator(),
		scanner:   validation.NewScanner(),
	}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ValidateExtension validates a single extension
func (a *App) ValidateExtension(extensionID string) (*models.ValidationResult, error) {
	log.Printf("[App] ValidateExtension called for: %s", extensionID)
	result, err := a.validator.ValidateExtension(extensionID)
	if err != nil {
		log.Printf("[App] ValidateExtension error: %v", err)
	}
	return result, err
}

// GetInstalledExtensions returns all installed extensions
func (a *App) GetInstalledExtensions(path string) ([]models.InstalledExtension, error) {
	log.Printf("[App] GetInstalledExtensions called with path: %s", path)
	exts, err := a.scanner.ScanInstalledExtensions(path)
	if err != nil {
		log.Printf("[App] GetInstalledExtensions error: %v", err)
	} else {
		log.Printf("[App] Found %d extensions", len(exts))
	}
	return exts, err
}

// AuditAllExtensions performs a full audit of all installed extensions
func (a *App) AuditAllExtensions(path string) (*models.AuditReport, error) {
	log.Printf("[App] AuditAllExtensions called with path: %s", path)
	report, err := a.scanner.AuditExtensions(path)
	if err != nil {
		log.Printf("[App] AuditAllExtensions error: %v", err)
	}
	return report, err
}

// DownloadOfficialExtension downloads the official version from Microsoft Marketplace
func (a *App) DownloadOfficialExtension(extensionID string) (string, string, error) {
	log.Printf("[App] DownloadOfficialExtension called for: %s", extensionID)
	data, hash, err := a.validator.DownloadOfficialExtension(extensionID)
	if err != nil {
		return "", "", err
	}

	// Save to temp file
	tempFile, err := os.CreateTemp("", fmt.Sprintf("%s-*.vsix", extensionID))
	if err != nil {
		return "", "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	if _, err := tempFile.Write(data); err != nil {
		return "", "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tempFile.Name(), hash, nil
}

// GetDefaultExtensionsPath returns the default VS Code extensions path
func (a *App) GetDefaultExtensionsPath() (string, error) {
	log.Println("[App] GetDefaultExtensionsPath called")
	path, err := validation.GetExtensionsPath()
	if err != nil {
		log.Printf("[App] GetDefaultExtensionsPath error: %v", err)
	} else {
		log.Printf("[App] Default extensions path: %s", path)
	}
	return path, err
}

// SearchMarketplace searches for extensions using keywords or wildcards
func (a *App) SearchMarketplace(searchTerm string) ([]*models.ExtensionMetadata, error) {
	log.Printf("[App] SearchMarketplace called for: %s", searchTerm)
	// Create marketplace client for search
	client := marketplace.NewClient()
	return client.SearchExtensions(searchTerm)
}

// SearchMarketplaceExtension searches for an extension in marketplace by ID and validates it
func (a *App) SearchMarketplaceExtension(extensionID string) (*models.ValidationResult, error) {
	log.Printf("[App] SearchMarketplaceExtension called for: %s", extensionID)
	return a.validator.ValidateExtension(extensionID)
}

// SelectDirectory opens a directory selection dialog
func (a *App) SelectDirectory() (string, error) {
	path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Extensions Directory",
	})
	if err != nil {
		return "", err
	}
	return path, nil
}

// ShowMessageDialog displays a message to the user
func (a *App) ShowMessageDialog(title, message string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})
}

// ShowErrorDialog displays an error message to the user
func (a *App) ShowErrorDialog(title, message string) {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   title,
		Message: message,
	})
}

// ShowQuestionDialog displays a yes/no question to the user
func (a *App) ShowQuestionDialog(title, message string) (string, error) {
	return runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.QuestionDialog,
		Title:   title,
		Message: message,
		Buttons: []string{"Yes", "No"},
	})
}

// ========== Editor Profile APIs ==========

// GetEditorProfiles returns all available editor profiles with default paths
func (a *App) GetEditorProfiles() []models.EditorProfile {
	log.Println("[App] GetEditorProfiles called")
	return editor.GetDefaultEditorProfiles()
}

// GetEditorStatus checks the availability and status of a specific editor
func (a *App) GetEditorStatus(editorType string) models.EditorStatus {
	log.Printf("[App] GetEditorStatus called for: %s", editorType)
	profile, err := editor.GetEditorProfile(models.EditorType(editorType))
	if err != nil {
		return models.EditorStatus{
			DisabledReason: err.Error(),
			IsAvailable:    false,
		}
	}
	return editor.CheckEditorStatus(profile)
}

// GetAllEditorStatuses returns the status of all known editors
func (a *App) GetAllEditorStatuses() []models.EditorStatus {
	log.Println("[App] GetAllEditorStatuses called")
	profiles := editor.GetDefaultEditorProfiles()
	statuses := make([]models.EditorStatus, 0, len(profiles))
	for _, profile := range profiles {
		statuses = append(statuses, editor.CheckEditorStatus(profile))
	}
	return statuses
}

// GetEditorExtensions returns all installed extensions for a specific editor
func (a *App) GetEditorExtensions(editorType string) ([]models.InstalledExtension, error) {
	log.Printf("[App] GetEditorExtensions called for: %s", editorType)
	profile, err := editor.GetEditorProfile(models.EditorType(editorType))
	if err != nil {
		return nil, err
	}
	return a.scanner.ScanInstalledExtensions(profile.ExtensionsDir)
}

// ========== CLI Status APIs ==========

// GetCLIStatus returns the availability of all VS Code family CLI tools
func (a *App) GetCLIStatus() models.CLIStatus {
	log.Println("[App] GetCLIStatus called")
	return editor.GetCLIStatus()
}

// InstallExtensionViaCLI installs an extension using the VS Code CLI
func (a *App) InstallExtensionViaCLI(cliCommand string, extensionID string) error {
	log.Printf("[App] InstallExtensionViaCLI called: cli=%s, ext=%s", cliCommand, extensionID)
	return editor.InstallExtensionViaCLI(cliCommand, extensionID)
}

// InstallExtensionsViaCLI installs multiple extensions using the VS Code CLI
func (a *App) InstallExtensionsViaCLI(cliCommand string, extensionIDs []string) models.InstallReport {
	log.Printf("[App] InstallExtensionsViaCLI called: cli=%s, exts=%v", cliCommand, extensionIDs)
	report := editor.InstallMultipleExtensionsViaCLI(cliCommand, extensionIDs)
	return report
}

// ========== Sync APIs ==========

// SyncExtensions syncs selected extensions from source to target editors
func (a *App) SyncExtensions(sourceEditor string, targetEditors []string, extensionIDs []string, overwriteConflicts bool) (*models.SyncReport, error) {
	log.Printf("[App] SyncExtensions called: source=%s, targets=%v, exts=%d, overwrite=%v",
		sourceEditor, targetEditors, len(extensionIDs), overwriteConflicts)

	targets := make([]models.EditorType, len(targetEditors))
	for i, t := range targetEditors {
		targets[i] = models.EditorType(t)
	}

	request := models.SyncRequest{
		SourceEditor:       models.EditorType(sourceEditor),
		TargetEditors:      targets,
		ExtensionIDs:       extensionIDs,
		OverwriteConflicts: overwriteConflicts,
	}

	return editor.SyncExtensions(request)
}

// DetectSyncConflicts checks which extensions would conflict if synced
func (a *App) DetectSyncConflicts(sourceEditor string, targetEditor string, extensionIDs []string) ([]string, error) {
	log.Printf("[App] DetectSyncConflicts called: source=%s, target=%s, exts=%d",
		sourceEditor, targetEditor, len(extensionIDs))
	return editor.DetectConflicts(models.EditorType(sourceEditor), models.EditorType(targetEditor), extensionIDs)
}

// GetExtensionsIndex reads the extensions.json for a specific editor
func (a *App) GetExtensionsIndex(editorType string) ([]models.ExtensionIndexEntry, error) {
	log.Printf("[App] GetExtensionsIndex called for: %s", editorType)
	profile, err := editor.GetEditorProfile(models.EditorType(editorType))
	if err != nil {
		return nil, err
	}
	return editor.ReadExtensionsIndex(profile.ExtensionsDir)
}

// ========== CLI Installation APIs ==========

// GetCLIInstallStatus checks if the vsynx CLI is installed and accessible
func (a *App) GetCLIInstallStatus() editor.CLIInstallStatus {
	log.Println("[App] GetCLIInstallStatus called")
	return editor.CheckVsynxCLIStatus()
}

// InstallCLI installs the vsynx CLI to the system PATH
func (a *App) InstallCLI() editor.CLIInstallResult {
	log.Println("[App] InstallCLI called")
	return editor.InstallVsynxCLI()
}

// UninstallCLI removes the vsynx CLI from the system PATH
func (a *App) UninstallCLI() editor.CLIInstallResult {
	log.Println("[App] UninstallCLI called")
	return editor.UninstallVsynxCLI()
}
