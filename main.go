package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/yourusername/secureopenvsx/cmd"
)

func main() {
	// Check executable name to determine mode
	execPath, _ := os.Executable()
	baseName := strings.ToLower(filepath.Base(execPath))
	// If named "vsynx" (not vsynx-manager), treat as CLI
	isCLI := strings.HasPrefix(baseName, "vsynx") && !strings.Contains(baseName, "manager")

	// On Windows, if running as CLI, attach to the parent console to show output
	if isCLI && runtime.GOOS == "windows" {
		modkernel32 := syscall.NewLazyDLL("kernel32.dll")
		procAttachConsole := modkernel32.NewProc("AttachConsole")
		const ATTACH_PARENT_PROCESS = ^uintptr(0) // -1
		procAttachConsole.Call(ATTACH_PARENT_PROCESS)
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
