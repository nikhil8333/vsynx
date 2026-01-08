package editor

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// CLIInstallStatus represents the current CLI installation status
type CLIInstallStatus struct {
	Installed    bool   `json:"installed"`
	Path         string `json:"path,omitempty"`
	Version      string `json:"version,omitempty"`
	CanInstall   bool   `json:"canInstall"`
	Instructions string `json:"instructions,omitempty"`
}

// CLIInstallResult represents the result of a CLI installation attempt
type CLIInstallResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Path         string `json:"path,omitempty"`
	NeedsRestart bool   `json:"needsRestart"`
}

// CheckVsynxCLIStatus checks if the vsynx CLI is installed and accessible
func CheckVsynxCLIStatus() CLIInstallStatus {
	status := CLIInstallStatus{
		CanInstall: true,
	}

	// Check if running in dev mode
	execPath, _ := os.Executable()
	isDevMode := strings.Contains(strings.ToLower(execPath), "wails") ||
		strings.Contains(strings.ToLower(execPath), "tmp") ||
		strings.Contains(strings.ToLower(execPath), "temp")

	if isDevMode {
		status.Installed = false
		status.CanInstall = false
		status.Instructions = "Development mode: Use 'go run . <command>' to test CLI commands.\n\nExamples:\n• go run . validate ms-python.python\n• go run . editors status\n• go run . marketplace search python\n• go run . sync preview --from vscode --to windsurf --all\n\nCLI installation is available only in production builds."
		return status
	}

	// Try to find vsynx in PATH
	var cliName string
	if runtime.GOOS == "windows" {
		cliName = "vsynx.exe"
	} else {
		cliName = "vsynx"
	}

	path, err := exec.LookPath(cliName)
	if err == nil {
		status.Installed = true
		status.Path = path

		// Try to get version
		cmd := exec.Command(path, "--version")
		output, err := cmd.Output()
		if err == nil {
			status.Version = strings.TrimSpace(string(output))
		}
	} else {
		// If vsynx.exe is missing but we are running the app, we can likely install it by copying ourselves.
		// So we report CanInstall = true, but Installed = false
		status.Installed = false
		status.CanInstall = true
		status.Instructions = "Click 'Install CLI' to configure the vsynx command."
	}

	return status
}

// InstallVsynxCLI installs the vsynx CLI to the system PATH
func InstallVsynxCLI() CLIInstallResult {
	switch runtime.GOOS {
	case "windows":
		return installCLIWindows()
	case "darwin":
		return installCLIDarwin()
	default:
		return installCLILinux()
	}
}

// UninstallVsynxCLI removes the vsynx CLI from the system PATH
func UninstallVsynxCLI() CLIInstallResult {
	switch runtime.GOOS {
	case "windows":
		return uninstallCLIWindows()
	case "darwin":
		return uninstallCLIDarwin()
	default:
		return uninstallCLILinux()
	}
}

// Windows implementation
func installCLIWindows() CLIInstallResult {
	// Get current executable directory
	execPath, err := os.Executable()
	if err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to get executable path: %v", err),
		}
	}
	execDir := filepath.Dir(execPath)

	// Check if vsynx.exe exists in the same directory
	cliPath := filepath.Join(execDir, "vsynx.exe")
	if _, err := os.Stat(cliPath); os.IsNotExist(err) {
		// If vsynx.exe is missing, we check if we (the running process) can act as the CLI.
		// Since main.go supports CLI args, we can just copy ourselves to "vsynx.exe"

		// Check if running in dev mode
		if strings.Contains(strings.ToLower(execPath), "wails") ||
			strings.Contains(strings.ToLower(execPath), "tmp") ||
			strings.Contains(strings.ToLower(execPath), "temp") {
			return CLIInstallResult{
				Success: false,
				Message: "Development mode detected. CLI installation is only available in production builds.",
			}
		}

		// Self-heal: Copy current executable to vsynx.exe
		if err := copyFile(execPath, cliPath); err != nil {
			return CLIInstallResult{
				Success: false,
				Message: fmt.Sprintf("vsynx.exe not found and failed to create copy from main binary: %v", err),
			}
		}
	}

	// Add to user PATH using setx
	// First get current user PATH
	cmd := exec.Command("cmd", "/c", "echo", "%PATH%")
	output, _ := cmd.Output()
	currentPath := strings.TrimSpace(string(output))

	// Check if already in PATH
	if strings.Contains(strings.ToLower(currentPath), strings.ToLower(execDir)) {
		return CLIInstallResult{
			Success: true,
			Message: "vsynx CLI is already in PATH.",
			Path:    cliPath,
		}
	}

	// Add to user PATH using PowerShell (more reliable than setx for PATH)
	psCmd := fmt.Sprintf(`[Environment]::SetEnvironmentVariable("Path", [Environment]::GetEnvironmentVariable("Path", "User") + ";%s", "User")`, execDir)
	cmd = exec.Command("powershell", "-Command", psCmd)
	if err := cmd.Run(); err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to update PATH: %v", err),
		}
	}

	return CLIInstallResult{
		Success:      true,
		Message:      fmt.Sprintf("vsynx CLI installed successfully.\nPath: %s\nPlease restart your terminal to use 'vsynx' command.", cliPath),
		Path:         cliPath,
		NeedsRestart: true,
	}
}

func uninstallCLIWindows() CLIInstallResult {
	execPath, err := os.Executable()
	if err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to get executable path: %v", err),
		}
	}
	execDir := filepath.Dir(execPath)

	// Remove vsynx.exe file with retry strategy
	cliPath := filepath.Join(execDir, "vsynx.exe")
	if _, err := os.Stat(cliPath); err == nil {
		// Try to verify it's not the running executable (just in case)
		if strings.EqualFold(execPath, cliPath) {
			// This should happen only if user renamed vsynx-manager.exe to vsynx.exe and ran it
			// We cannot delete ourselves.
			return CLIInstallResult{
				Success: false,
				Message: "Cannot uninstall CLI because the application is running as 'vsynx.exe'. Please rename the executable to 'vsynx-manager.exe' and try again.",
			}
		}

		// Attempt removal
		err := os.Remove(cliPath)
		if err != nil {
			// Simple retry logic for Windows file locking
			// Wait 500ms and try again
			time.Sleep(500 * time.Millisecond)
			if err := os.Remove(cliPath); err != nil {
				return CLIInstallResult{
					Success: false,
					Message: fmt.Sprintf("Failed to remove vsynx.exe: %v. Please ensure no terminal windows are using 'vsynx' and try again, or delete '%s' manually.", err, cliPath),
				}
			}
		}
	}

	// Remove from user PATH using PowerShell
	psCmd := fmt.Sprintf(`$path = [Environment]::GetEnvironmentVariable("Path", "User"); $path = ($path.Split(";") | Where-Object { $_ -ne "%s" }) -join ";"; [Environment]::SetEnvironmentVariable("Path", $path, "User")`, execDir)
	cmd := exec.Command("powershell", "-Command", psCmd)
	if err := cmd.Run(); err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to update PATH: %v", err),
		}
	}

	return CLIInstallResult{
		Success:      true,
		Message:      "vsynx CLI removed from PATH. Please restart your terminal.",
		NeedsRestart: true,
	}
}

// macOS implementation
func installCLIDarwin() CLIInstallResult {
	execPath, err := os.Executable()
	if err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to get executable path: %v", err),
		}
	}

	// For macOS app bundle, the CLI should be bundled inside
	// Try to find it relative to the app
	var cliSource string
	if strings.Contains(execPath, ".app/") {
		// Inside app bundle
		appPath := execPath[:strings.Index(execPath, ".app/")+5]
		cliSource = filepath.Join(appPath, "Contents", "MacOS", "vsynx")
	} else {
		cliSource = filepath.Join(filepath.Dir(execPath), "vsynx")
	}

	if _, err := os.Stat(cliSource); os.IsNotExist(err) {
		return CLIInstallResult{
			Success: false,
			Message: "vsynx CLI not found in app bundle. Try: brew install vsynx",
		}
	}

	// Try /usr/local/bin first, fall back to ~/.local/bin
	targetDirs := []string{"/usr/local/bin", filepath.Join(os.Getenv("HOME"), ".local", "bin")}

	for _, targetDir := range targetDirs {
		targetPath := filepath.Join(targetDir, "vsynx")

		// Ensure target directory exists
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			continue
		}

		// Create symlink
		os.Remove(targetPath) // Remove existing
		if err := os.Symlink(cliSource, targetPath); err != nil {
			continue
		}

		return CLIInstallResult{
			Success:      true,
			Message:      fmt.Sprintf("vsynx CLI installed to %s", targetPath),
			Path:         targetPath,
			NeedsRestart: targetDir == filepath.Join(os.Getenv("HOME"), ".local", "bin"),
		}
	}

	return CLIInstallResult{
		Success: false,
		Message: "Failed to install CLI. Try: brew install vsynx",
	}
}

func uninstallCLIDarwin() CLIInstallResult {
	paths := []string{
		"/usr/local/bin/vsynx",
		filepath.Join(os.Getenv("HOME"), ".local", "bin", "vsynx"),
	}

	removed := false
	for _, p := range paths {
		if _, err := os.Lstat(p); err == nil {
			if err := os.Remove(p); err == nil {
				removed = true
			}
		}
	}

	if removed {
		return CLIInstallResult{
			Success: true,
			Message: "vsynx CLI removed.",
		}
	}

	return CLIInstallResult{
		Success: false,
		Message: "vsynx CLI not found in standard locations.",
	}
}

// Linux implementation
func installCLILinux() CLIInstallResult {
	execPath, err := os.Executable()
	if err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to get executable path: %v", err),
		}
	}

	cliSource := filepath.Join(filepath.Dir(execPath), "vsynx")
	if _, err := os.Stat(cliSource); os.IsNotExist(err) {
		return CLIInstallResult{
			Success: false,
			Message: "vsynx CLI not found alongside the application.",
		}
	}

	// Install to ~/.local/bin
	homeDir := os.Getenv("HOME")
	targetDir := filepath.Join(homeDir, ".local", "bin")
	targetPath := filepath.Join(targetDir, "vsynx")

	// Ensure directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to create directory: %v", err),
		}
	}

	// Create symlink
	os.Remove(targetPath) // Remove existing
	if err := os.Symlink(cliSource, targetPath); err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to create symlink: %v", err),
		}
	}

	// Check if ~/.local/bin is in PATH
	pathEnv := os.Getenv("PATH")
	needsPathUpdate := !strings.Contains(pathEnv, targetDir)

	message := fmt.Sprintf("vsynx CLI installed to %s", targetPath)
	if needsPathUpdate {
		message += "\n\nAdd to your shell profile:\nexport PATH=\"$HOME/.local/bin:$PATH\""
	}

	return CLIInstallResult{
		Success:      true,
		Message:      message,
		Path:         targetPath,
		NeedsRestart: needsPathUpdate,
	}
}

func uninstallCLILinux() CLIInstallResult {
	homeDir := os.Getenv("HOME")
	targetPath := filepath.Join(homeDir, ".local", "bin", "vsynx")

	if _, err := os.Lstat(targetPath); os.IsNotExist(err) {
		return CLIInstallResult{
			Success: false,
			Message: "vsynx CLI not found in ~/.local/bin",
		}
	}

	if err := os.Remove(targetPath); err != nil {
		return CLIInstallResult{
			Success: false,
			Message: fmt.Sprintf("Failed to remove CLI: %v", err),
		}
	}

	return CLIInstallResult{
		Success: true,
		Message: "vsynx CLI removed from ~/.local/bin",
	}
}
