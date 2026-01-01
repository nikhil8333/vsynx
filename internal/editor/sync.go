package editor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/secureopenvsx/internal/models"
)

// SyncExtensions syncs extensions from source editor to target editors
func SyncExtensions(request models.SyncRequest) (*models.SyncReport, error) {
	// Get source editor profile
	sourceProfile, err := GetEditorProfile(request.SourceEditor)
	if err != nil {
		return nil, fmt.Errorf("invalid source editor: %w", err)
	}

	// Check source is available
	sourceStatus := CheckEditorStatus(sourceProfile)
	if !sourceStatus.IsAvailable {
		return nil, fmt.Errorf("source editor not available: %s", sourceStatus.DisabledReason)
	}

	// Read source extensions index
	sourceIndex, err := ReadExtensionsIndex(sourceProfile.ExtensionsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read source extensions index: %w", err)
	}

	report := &models.SyncReport{
		SourceEditor: request.SourceEditor,
		Results:      make([]models.SyncResult, 0, len(request.TargetEditors)),
	}

	// Sync to each target
	for _, targetType := range request.TargetEditors {
		result := syncToTarget(sourceProfile, sourceIndex, targetType, request.ExtensionIDs, request.OverwriteConflicts)
		report.Results = append(report.Results, result)
		report.TotalCopied += result.CopiedCount
		report.TotalSkipped += result.SkippedCount
		report.TotalErrors += len(result.Errors)
	}

	return report, nil
}

// syncToTarget syncs extensions to a single target editor
func syncToTarget(sourceProfile models.EditorProfile, sourceIndex []models.ExtensionIndexEntry, targetType models.EditorType, extensionIDs []string, overwriteConflicts bool) models.SyncResult {
	result := models.SyncResult{
		TargetEditor: targetType,
		Conflicts:    []string{},
		Errors:       []string{},
	}

	// Get target profile
	targetProfile, err := GetEditorProfile(targetType)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid target editor: %s", err))
		return result
	}

	// Check target is available
	targetStatus := CheckEditorStatus(targetProfile)
	if !targetStatus.DirExists {
		// Create the extensions directory if it doesn't exist
		if err := os.MkdirAll(targetProfile.ExtensionsDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to create target directory: %s", err))
			return result
		}
	}

	// Read target extensions index (may not exist)
	targetIndex, _ := ReadExtensionsIndex(targetProfile.ExtensionsDir)

	// Build map of target extensions for conflict detection
	targetExtMap := make(map[string]*models.ExtensionIndexEntry)
	for i, entry := range targetIndex {
		targetExtMap[strings.ToLower(entry.Identifier.ID)] = &targetIndex[i]
	}

	// Track entries to add to target index
	entriesToAdd := []models.ExtensionIndexEntry{}

	// Process each extension
	for _, extID := range extensionIDs {
		extIDLower := strings.ToLower(extID)

		// Find in source index
		sourceEntry := FindExtensionEntry(sourceIndex, extID)
		if sourceEntry == nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Extension %s not found in source index", extID))
			continue
		}

		// Check for conflict
		if existingEntry, exists := targetExtMap[extIDLower]; exists {
			result.Conflicts = append(result.Conflicts, extID)
			if !overwriteConflicts {
				result.SkippedCount++
				continue
			}
			// Remove old folder if overwriting
			oldFolderPath := filepath.Join(targetProfile.ExtensionsDir, existingEntry.RelativeLocation)
			if err := os.RemoveAll(oldFolderPath); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Failed to remove existing extension %s: %s", extID, err))
				continue
			}
			result.OverwrittenCount++
		}

		// Copy extension folder
		sourceFolderPath := filepath.Join(sourceProfile.ExtensionsDir, sourceEntry.RelativeLocation)
		targetFolderPath := filepath.Join(targetProfile.ExtensionsDir, sourceEntry.RelativeLocation)

		if err := copyDir(sourceFolderPath, targetFolderPath); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to copy extension %s: %s", extID, err))
			continue
		}

		// Create new index entry for target
		newEntry := models.ExtensionIndexEntry{
			Identifier: models.ExtensionIdentifier{
				ID:   sourceEntry.Identifier.ID,
				UUID: sourceEntry.Identifier.UUID,
			},
			Version:          sourceEntry.Version,
			RelativeLocation: sourceEntry.RelativeLocation,
			Location: models.ExtensionLocation{
				Mid:    1,
				Path:   "/" + strings.ReplaceAll(targetFolderPath, "\\", "/"),
				Scheme: "file",
			},
			Metadata: sourceEntry.Metadata,
		}
		entriesToAdd = append(entriesToAdd, newEntry)
		result.CopiedCount++
	}

	// Update target index
	if len(entriesToAdd) > 0 {
		// Merge with existing entries
		newIndex := make([]models.ExtensionIndexEntry, 0, len(targetIndex)+len(entriesToAdd))
		
		// Add existing entries (excluding overwritten ones)
		for _, entry := range targetIndex {
			entryID := strings.ToLower(entry.Identifier.ID)
			overwritten := false
			for _, added := range entriesToAdd {
				if strings.ToLower(added.Identifier.ID) == entryID {
					overwritten = true
					break
				}
			}
			if !overwritten {
				newIndex = append(newIndex, entry)
			}
		}
		
		// Add new entries
		newIndex = append(newIndex, entriesToAdd...)

		if err := WriteExtensionsIndex(targetProfile.ExtensionsDir, newIndex); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("Failed to update index: %s", err))
		} else {
			result.IndexUpdated = true
		}
	}

	result.Success = len(result.Errors) == 0
	return result
}

// copyDir copies a directory recursively
func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("source directory not found: %w", err)
	}

	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile copies a single file
func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// DetectConflicts checks which extensions would conflict if synced
func DetectConflicts(sourceEditor models.EditorType, targetEditor models.EditorType, extensionIDs []string) ([]string, error) {
	targetProfile, err := GetEditorProfile(targetEditor)
	if err != nil {
		return nil, err
	}

	targetIndex, err := ReadExtensionsIndex(targetProfile.ExtensionsDir)
	if err != nil {
		// No index = no conflicts
		return []string{}, nil
	}

	conflicts := []string{}
	for _, extID := range extensionIDs {
		if FindExtensionEntry(targetIndex, extID) != nil {
			conflicts = append(conflicts, extID)
		}
	}

	return conflicts, nil
}
