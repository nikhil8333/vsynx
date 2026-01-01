package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/secureopenvsx/internal/editor"
	"github.com/yourusername/secureopenvsx/internal/models"
	"github.com/yourusername/secureopenvsx/internal/validation"
)

var editorsCmd = &cobra.Command{
	Use:   "editors",
	Short: "Manage and inspect code editors",
	Long:  `Commands for listing, inspecting, and managing VS Code family editors and their extensions.`,
}

var editorsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all known editors",
	Long:  `Lists all known VS Code family editors with their default extension paths.`,
	Run: func(cmd *cobra.Command, args []string) {
		profiles := editor.GetDefaultEditorProfiles()

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(profiles, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n=== Known Editors (%d) ===\n\n", len(profiles))
		fmt.Printf("%-20s %-15s %-50s\n", "ID", "CLI Command", "Extensions Path")
		fmt.Println(strings.Repeat("-", 90))

		for _, p := range profiles {
			cli := p.CLICommand
			if cli == "" {
				cli = "(none)"
			}
			fmt.Printf("%-20s %-15s %-50s\n", p.ID, cli, p.ExtensionsDir)
		}
		fmt.Println()
	},
}

var editorsStatusCmd = &cobra.Command{
	Use:   "status [editor-id]",
	Short: "Check status of an editor",
	Long:  `Checks availability, extension count, and CLI status for a specific editor.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		profiles := editor.GetDefaultEditorProfiles()

		// If no editor specified, show all statuses
		if len(args) == 0 {
			statuses := make([]models.EditorStatus, 0, len(profiles))
			for _, p := range profiles {
				statuses = append(statuses, editor.CheckEditorStatus(p))
			}

			if outputFormat == "json" {
				data, _ := json.MarshalIndent(statuses, "", "  ")
				fmt.Println(string(data))
				return
			}

			fmt.Printf("\n=== Editor Status ===\n\n")
			fmt.Printf("%-20s %-12s %-10s %-12s %s\n", "Editor", "Available", "Exts", "CLI", "Path")
			fmt.Println(strings.Repeat("-", 100))

			for _, s := range statuses {
				available := "No"
				if s.IsAvailable {
					available = "Yes"
				}
				cli := "No"
				if s.CLIAvailable {
					cli = "Yes"
				}
				fmt.Printf("%-20s %-12s %-10d %-12s %s\n",
					s.Editor.ID, available, s.ExtensionCount, cli, s.Editor.ExtensionsDir)
			}
			fmt.Println()
			return
		}

		// Specific editor
		editorID := args[0]
		profile, err := editor.GetEditorProfile(models.EditorType(editorID))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown editor: %s\n", editorID)
			os.Exit(1)
		}

		status := editor.CheckEditorStatus(profile)

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(status, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n=== Editor Status: %s ===\n\n", status.Editor.Name)
		fmt.Printf("  ID:              %s\n", status.Editor.ID)
		fmt.Printf("  Extensions Dir:  %s\n", status.Editor.ExtensionsDir)
		fmt.Printf("  Dir Exists:      %v\n", status.DirExists)
		fmt.Printf("  Index Exists:    %v\n", status.IndexFileExists)
		fmt.Printf("  Extension Count: %d\n", status.ExtensionCount)
		fmt.Printf("  CLI Available:   %v\n", status.CLIAvailable)
		if status.CLIPath != "" {
			fmt.Printf("  CLI Path:        %s\n", status.CLIPath)
		}
		if status.DisabledReason != "" {
			fmt.Printf("  Disabled Reason: %s\n", status.DisabledReason)
		}
		fmt.Println()
	},
}

var editorsExtensionsCmd = &cobra.Command{
	Use:   "extensions [editor-id]",
	Short: "List extensions for an editor",
	Long:  `Lists all installed extensions for a specific editor.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		editorID := args[0]

		profile, err := editor.GetEditorProfile(models.EditorType(editorID))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown editor: %s\n", editorID)
			os.Exit(1)
		}

		scanner := validation.NewScanner()
		extensions, err := scanner.ScanInstalledExtensions(profile.ExtensionsDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning extensions: %v\n", err)
			os.Exit(1)
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(extensions, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n=== Extensions for %s (%d) ===\n\n", profile.Name, len(extensions))

		if len(extensions) == 0 {
			fmt.Println("No extensions found.")
			return
		}

		fmt.Printf("%-45s %-15s %-10s\n", "Extension ID", "Version", "Status")
		fmt.Println(strings.Repeat("-", 75))

		for _, ext := range extensions {
			status := "Enabled"
			if !ext.IsEnabled {
				status = "Disabled"
			}
			fmt.Printf("%-45s %-15s %-10s\n", ext.ID, ext.Version, status)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(editorsCmd)
	editorsCmd.AddCommand(editorsListCmd)
	editorsCmd.AddCommand(editorsStatusCmd)
	editorsCmd.AddCommand(editorsExtensionsCmd)
}
