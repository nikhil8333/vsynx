package validation

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yourusername/secureopenvsx/internal/models"
)

// Scanner handles scanning for installed extensions
type Scanner struct {
	validator *Validator
}

// NewScanner creates a new scanner instance
func NewScanner() *Scanner {
	return &Scanner{
		validator: NewValidator(),
	}
}

// packageJSON represents the structure of an extension's package.json
type packageJSON struct {
	Publisher   string `json:"publisher"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
	Repository  struct {
		URL string `json:"url"`
	} `json:"repository"`
}

// GetExtensionsPath returns the path to the VS Code extensions directory
func GetExtensionsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	log.Printf("[Scanner] Home directory: %s", homeDir)

	// Default VS Code extensions path
	var extensionsPath string

	// Check for different editors and locations
	paths := []string{
		// Linux/WSL standard locations
		filepath.Join(homeDir, ".vscode", "extensions"),
		filepath.Join(homeDir, ".vscode-server", "extensions"),
		filepath.Join(homeDir, ".vscode-oss", "extensions"),
		filepath.Join(homeDir, ".vscodium", "extensions"),
		filepath.Join(homeDir, ".config", "Code", "User", "extensions"),
		// Windows paths (native)
		filepath.Join(homeDir, "AppData", "Local", "Programs", "Microsoft VS Code", "resources", "app", "extensions"),
		// Windows user extensions
		filepath.Join(homeDir, ".vscode", "extensions"),
		// macOS
		filepath.Join(homeDir, "Library", "Application Support", "Code", "User", "extensions"),
		// WSL accessing Windows paths
		"/mnt/c/Users/" + getWindowsUsername() + "/.vscode/extensions",
		"/mnt/c/Users/" + getWindowsUsername() + "/AppData/Local/Programs/Microsoft VS Code/resources/app/extensions",
	}

	log.Printf("[Scanner] Checking %d possible paths", len(paths))
	for i, path := range paths {
		log.Printf("[Scanner] Checking path %d: %s", i+1, path)
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			log.Printf("[Scanner] Found extensions directory: %s", path)
			extensionsPath = path
			break
		} else if err != nil {
			log.Printf("[Scanner] Path does not exist: %s", path)
		}
	}

	if extensionsPath == "" {
		log.Printf("[Scanner] Could not find VS Code extensions directory in any standard location")
		return "", fmt.Errorf("could not find VS Code extensions directory. Please use 'Change Path' to select manually. Checked locations: %s, %s, %s",
			filepath.Join(homeDir, ".vscode", "extensions"),
			filepath.Join(homeDir, ".vscode-server", "extensions"),
			filepath.Join(homeDir, ".config", "Code", "User", "extensions"))
	}

	return extensionsPath, nil
}

// getWindowsUsername extracts Windows username from WSL home path or returns current user
func getWindowsUsername() string {
	// Try to get from environment or current user
	if username := os.Getenv("USER"); username != "" {
		return username
	}
	if homeDir, err := os.UserHomeDir(); err == nil {
		// Extract username from /home/username
		parts := strings.Split(homeDir, string(filepath.Separator))
		if len(parts) > 2 {
			return parts[2]
		}
	}
	return "*"
}

// ScanInstalledExtensions scans for all installed extensions
func (s *Scanner) ScanInstalledExtensions(extensionsPath string) ([]models.InstalledExtension, error) {
	log.Printf("[Scanner] Scanning extensions at: %s", extensionsPath)
	if extensionsPath == "" {
		var err error
		extensionsPath, err = GetExtensionsPath()
		if err != nil {
			return nil, err
		}
	}

	// Expand ~ to home directory
	if strings.HasPrefix(extensionsPath, "~") {
		homeDir, err := os.UserHomeDir()
		if err == nil {
			extensionsPath = filepath.Join(homeDir, extensionsPath[1:])
			log.Printf("[Scanner] Expanded path to: %s", extensionsPath)
		}
	}

	entries, err := os.ReadDir(extensionsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read extensions directory: %w", err)
	}

	var extensions []models.InstalledExtension
	log.Printf("[Scanner] Found %d entries in directory", len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip system extensions and directories starting with '.'
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		extPath := filepath.Join(extensionsPath, entry.Name())
		packageJSONPath := filepath.Join(extPath, "package.json")

		// Check if package.json exists
		if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
			continue
		}

		// Read package.json
		data, err := os.ReadFile(packageJSONPath)
		if err != nil {
			continue
		}

		var pkg packageJSON
		if err := json.Unmarshal(data, &pkg); err != nil {
			continue
		}

		info, _ := entry.Info()
		var lastModified time.Time
		if info != nil {
			lastModified = info.ModTime()
		}

		extension := models.InstalledExtension{
			ID:           fmt.Sprintf("%s.%s", pkg.Publisher, pkg.Name),
			Path:         extPath,
			Publisher:    pkg.Publisher,
			Name:         pkg.Name,
			Version:      pkg.Version,
			IsEnabled:    true, // Assume enabled if installed
			LastModified: lastModified,
		}

		extensions = append(extensions, extension)
	}

	log.Printf("[Scanner] Successfully scanned %d extensions", len(extensions))
	return extensions, nil
}

// AuditExtensions performs a full audit of all installed extensions
func (s *Scanner) AuditExtensions(extensionsPath string) (*models.AuditReport, error) {
	log.Printf("[Scanner] Starting audit of extensions at: %s", extensionsPath)
	extensions, err := s.ScanInstalledExtensions(extensionsPath)
	if err != nil {
		return nil, err
	}

	report := &models.AuditReport{
		TotalExtensions: len(extensions),
		Results:         make([]models.ValidationResult, 0, len(extensions)),
		AuditTime:       time.Now(),
	}

	for _, ext := range extensions {
		result, err := s.validator.ValidateExtension(ext.ID)
		if err != nil {
			// Still add to results with error
			result = &models.ValidationResult{
				ExtensionID:    ext.ID,
				TrustLevel:     models.TrustLevelUnknown,
				ValidationTime: time.Now(),
				Error:          err.Error(),
			}
		}

		report.Results = append(report.Results, *result)

		// Update counters
		switch result.TrustLevel {
		case models.TrustLevelLegitimate:
			report.LegitimateCount++
		case models.TrustLevelSuspicious:
			report.SuspiciousCount++
		case models.TrustLevelMalicious:
			report.MaliciousCount++
		default:
			report.UnknownCount++
		}
	}

	log.Printf("[Scanner] Audit complete: %d total, %d legitimate, %d suspicious, %d malicious, %d unknown",
		report.TotalExtensions, report.LegitimateCount, report.SuspiciousCount, report.MaliciousCount, report.UnknownCount)
	return report, nil
}
