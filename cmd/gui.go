package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

var guiCmd = &cobra.Command{
	Use:   "gui",
	Short: "Launch the Vsynx Manager GUI",
	Long:  `Launches the graphical user interface (Vsynx Manager) for managing VS Code extensions.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the path to the current executable
		execPath, err := os.Executable()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding executable: %v\n", err)
			os.Exit(1)
		}

		// On Windows, we need to launch vsynx-manager.exe
		// On other platforms, launch vsynx-manager
		var guiPath string
		if runtime.GOOS == "windows" {
			guiPath = execPath[:len(execPath)-len("vsynx.exe")] + "vsynx-manager.exe"
		} else {
			guiPath = execPath[:len(execPath)-len("vsynx")] + "vsynx-manager"
		}

		// Check if GUI exists
		if _, err := os.Stat(guiPath); os.IsNotExist(err) {
			fmt.Println("Vsynx Manager GUI not found.")
			fmt.Println("Please install the full Vsynx Manager application.")
			os.Exit(1)
		}

		// Launch GUI
		guiCmd := exec.Command(guiPath)
		if err := guiCmd.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Error launching GUI: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Vsynx Manager launched.")
	},
}

func init() {
	rootCmd.AddCommand(guiCmd)
}
