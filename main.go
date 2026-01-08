package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/yourusername/secureopenvsx/cmd"
)

func main() {
	// Check executable name to determine mode
	execPath, _ := os.Executable()
	baseName := strings.ToLower(filepath.Base(execPath))
	// If named "vsynx" (not vsynx-manager), treat as CLI
	isCLI := strings.HasPrefix(baseName, "vsynx") && !strings.Contains(baseName, "manager")

	// If running as CLI, attach to the parent console if necessary (Windows only)
	if isCLI {
		attachConsole()
	}

	// Routing Logic:
	// 1. If explicit "gui" argument, start GUI
	// 2. If NOT CLI mode (i.e. named vsynx-manager) AND no args, start GUI
	// 3. Otherwise, run CLI
	if (len(os.Args) > 1 && os.Args[1] == "gui") || (!isCLI && len(os.Args) == 1) {
		startGUI()
	} else {
		cmd.Execute()
	}
}
