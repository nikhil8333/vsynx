package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/secureopenvsx/internal/models"
	"github.com/yourusername/secureopenvsx/internal/validation"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed extensions",
	Long:  `Lists all installed VS Code extensions without performing validation.`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := validation.NewScanner()

		extensions, err := scanner.ScanInstalledExtensions(extensionsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning extensions: %v\n", err)
			os.Exit(1)
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(extensions, "", "  ")
			fmt.Println(string(data))
		} else {
			printExtensionsList(extensions)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func printExtensionsList(extensions []models.InstalledExtension) {
	fmt.Printf("\n=== Installed Extensions (%d) ===\n\n", len(extensions))

	if len(extensions) == 0 {
		fmt.Println("No extensions found.")
		return
	}

	// Print table header
	fmt.Printf("%-40s %-15s %-10s\n", "Extension ID", "Version", "Status")
	fmt.Println(strings.Repeat("-", 70))

	for _, ext := range extensions {
		status := "Enabled"
		if !ext.IsEnabled {
			status = "Disabled"
		}
		fmt.Printf("%-40s %-15s %-10s\n", ext.ID, ext.Version, status)
	}

	fmt.Println()
}
