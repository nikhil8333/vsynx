package main

import (
	"os"

	"github.com/yourusername/secureopenvsx/cmd"
)

func main() {
	// Check if running as CLI or GUI
	// If no arguments (or only "gui" argument), start GUI
	// Otherwise, run CLI commands
	
	if len(os.Args) == 1 || (len(os.Args) == 2 && os.Args[1] == "gui") {
		// Start GUI
		startGUI()
	} else {
		// Run CLI
		cmd.Execute()
	}
}
