package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/secureopenvsx/internal/editor"
	"github.com/yourusername/secureopenvsx/internal/models"
)

var (
	syncFrom      string
	syncTo        string
	syncExts      string
	syncAll       bool
	syncOverwrite bool
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync extensions between editors",
	Long: `Commands for syncing VS Code extensions between different editors.
Supports copying extensions from a source editor (e.g., vscode) to target editors
(e.g., windsurf, cursor, kiro).`,
}

var syncRunCmd = &cobra.Command{
	Use:   "run",
	Short: "Execute extension sync",
	Long: `Syncs selected extensions from the source editor to one or more target editors.
Use --overwrite to replace existing extensions in target editors.`,
	Run: func(cmd *cobra.Command, args []string) {
		if syncFrom == "" {
			fmt.Fprintln(os.Stderr, "Error: --from is required")
			os.Exit(1)
		}
		if syncTo == "" {
			fmt.Fprintln(os.Stderr, "Error: --to is required")
			os.Exit(1)
		}
		if syncExts == "" && !syncAll {
			fmt.Fprintln(os.Stderr, "Error: --ext or --all is required")
			os.Exit(1)
		}

		// Parse target editors
		targets := strings.Split(syncTo, ",")
		targetTypes := make([]models.EditorType, len(targets))
		for i, t := range targets {
			targetTypes[i] = models.EditorType(strings.TrimSpace(t))
		}

		// Get extension IDs
		var extensionIDs []string
		if syncAll {
			// Get all extensions from source
			profile, err := editor.GetEditorProfile(models.EditorType(syncFrom))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid source editor: %v\n", err)
				os.Exit(1)
			}
			entries, err := editor.ReadExtensionsIndex(profile.ExtensionsDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading source extensions: %v\n", err)
				os.Exit(1)
			}
			for _, e := range entries {
				extensionIDs = append(extensionIDs, e.Identifier.ID)
			}
		} else {
			extensionIDs = strings.Split(syncExts, ",")
			for i, id := range extensionIDs {
				extensionIDs[i] = strings.TrimSpace(id)
			}
		}

		request := models.SyncRequest{
			SourceEditor:       models.EditorType(syncFrom),
			TargetEditors:      targetTypes,
			ExtensionIDs:       extensionIDs,
			OverwriteConflicts: syncOverwrite,
		}

		report, err := editor.SyncExtensions(request)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error syncing extensions: %v\n", err)
			os.Exit(1)
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(report, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n=== Sync Report ===\n")
		fmt.Printf("Source: %s\n", report.SourceEditor)
		fmt.Printf("Extensions: %d\n\n", len(extensionIDs))

		hasConflicts := false
		for _, result := range report.Results {
			status := "✓ Success"
			if !result.Success {
				status = "✗ Failed"
			}
			fmt.Printf("Target: %s - %s\n", result.TargetEditor, status)
			fmt.Printf("  Copied: %d, Skipped: %d, Overwritten: %d\n",
				result.CopiedCount, result.SkippedCount, result.OverwrittenCount)

			if len(result.Conflicts) > 0 {
				hasConflicts = true
				fmt.Printf("  Conflicts: %s\n", strings.Join(result.Conflicts, ", "))
			}
			if len(result.Errors) > 0 {
				fmt.Printf("  Errors:\n")
				for _, e := range result.Errors {
					fmt.Printf("    - %s\n", e)
				}
			}
		}

		fmt.Printf("\nTotal: Copied %d, Skipped %d, Errors %d\n",
			report.TotalCopied, report.TotalSkipped, report.TotalErrors)

		if hasConflicts && !syncOverwrite {
			fmt.Println("\nNote: Use --overwrite to replace existing extensions.")
			os.Exit(3)
		}
	},
}

var syncPreviewCmd = &cobra.Command{
	Use:   "preview",
	Short: "Preview sync operation",
	Long:  `Shows what would happen during a sync without actually performing it.`,
	Run: func(cmd *cobra.Command, args []string) {
		if syncFrom == "" {
			fmt.Fprintln(os.Stderr, "Error: --from is required")
			os.Exit(1)
		}
		if syncTo == "" {
			fmt.Fprintln(os.Stderr, "Error: --to is required")
			os.Exit(1)
		}
		if syncExts == "" && !syncAll {
			fmt.Fprintln(os.Stderr, "Error: --ext or --all is required")
			os.Exit(1)
		}

		// Parse target editors
		targets := strings.Split(syncTo, ",")

		// Get extension IDs
		var extensionIDs []string
		if syncAll {
			profile, err := editor.GetEditorProfile(models.EditorType(syncFrom))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid source editor: %v\n", err)
				os.Exit(1)
			}
			entries, err := editor.ReadExtensionsIndex(profile.ExtensionsDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading source extensions: %v\n", err)
				os.Exit(1)
			}
			for _, e := range entries {
				extensionIDs = append(extensionIDs, e.Identifier.ID)
			}
		} else {
			extensionIDs = strings.Split(syncExts, ",")
			for i, id := range extensionIDs {
				extensionIDs[i] = strings.TrimSpace(id)
			}
		}

		type PreviewResult struct {
			Target    string   `json:"target"`
			NewCount  int      `json:"newCount"`
			Conflicts []string `json:"conflicts"`
			Overwrite int      `json:"overwriteCount"`
		}

		var results []PreviewResult

		for _, target := range targets {
			target = strings.TrimSpace(target)
			conflicts, err := editor.DetectConflicts(
				models.EditorType(syncFrom),
				models.EditorType(target),
				extensionIDs,
			)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error checking conflicts for %s: %v\n", target, err)
				continue
			}

			results = append(results, PreviewResult{
				Target:    target,
				NewCount:  len(extensionIDs) - len(conflicts),
				Conflicts: conflicts,
				Overwrite: len(conflicts),
			})
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n=== Sync Preview ===\n")
		fmt.Printf("Source: %s\n", syncFrom)
		fmt.Printf("Extensions to sync: %d\n\n", len(extensionIDs))

		for _, r := range results {
			fmt.Printf("Target: %s\n", r.Target)
			fmt.Printf("  New (to install): %d\n", r.NewCount)
			fmt.Printf("  Conflicts (overwrite): %d\n", r.Overwrite)
			if len(r.Conflicts) > 0 {
				fmt.Printf("  Conflicting IDs: %s\n", strings.Join(r.Conflicts, ", "))
			}
			fmt.Println()
		}
	},
}

var syncConflictsCmd = &cobra.Command{
	Use:   "conflicts",
	Short: "Check for sync conflicts",
	Long:  `Checks which extensions would conflict when syncing to a specific target editor.`,
	Run: func(cmd *cobra.Command, args []string) {
		if syncFrom == "" {
			fmt.Fprintln(os.Stderr, "Error: --from is required")
			os.Exit(1)
		}
		if syncTo == "" {
			fmt.Fprintln(os.Stderr, "Error: --to is required (single target)")
			os.Exit(1)
		}
		if syncExts == "" && !syncAll {
			fmt.Fprintln(os.Stderr, "Error: --ext or --all is required")
			os.Exit(1)
		}

		// Get extension IDs
		var extensionIDs []string
		if syncAll {
			profile, err := editor.GetEditorProfile(models.EditorType(syncFrom))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: invalid source editor: %v\n", err)
				os.Exit(1)
			}
			entries, err := editor.ReadExtensionsIndex(profile.ExtensionsDir)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading source extensions: %v\n", err)
				os.Exit(1)
			}
			for _, e := range entries {
				extensionIDs = append(extensionIDs, e.Identifier.ID)
			}
		} else {
			extensionIDs = strings.Split(syncExts, ",")
			for i, id := range extensionIDs {
				extensionIDs[i] = strings.TrimSpace(id)
			}
		}

		conflicts, err := editor.DetectConflicts(
			models.EditorType(syncFrom),
			models.EditorType(syncTo),
			extensionIDs,
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting conflicts: %v\n", err)
			os.Exit(1)
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(map[string]interface{}{
				"source":        syncFrom,
				"target":        syncTo,
				"totalChecked":  len(extensionIDs),
				"conflictCount": len(conflicts),
				"conflicts":     conflicts,
			}, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n=== Sync Conflicts ===\n")
		fmt.Printf("Source: %s → Target: %s\n", syncFrom, syncTo)
		fmt.Printf("Extensions checked: %d\n", len(extensionIDs))
		fmt.Printf("Conflicts found: %d\n\n", len(conflicts))

		if len(conflicts) == 0 {
			fmt.Println("No conflicts found. All extensions can be synced safely.")
		} else {
			fmt.Println("Conflicting extensions (already exist in target):")
			for _, c := range conflicts {
				fmt.Printf("  - %s\n", c)
			}
			os.Exit(3)
		}
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
	syncCmd.AddCommand(syncRunCmd)
	syncCmd.AddCommand(syncPreviewCmd)
	syncCmd.AddCommand(syncConflictsCmd)

	// Shared flags for all sync subcommands
	for _, cmd := range []*cobra.Command{syncRunCmd, syncPreviewCmd, syncConflictsCmd} {
		cmd.Flags().StringVar(&syncFrom, "from", "", "Source editor (e.g., vscode)")
		cmd.Flags().StringVar(&syncTo, "to", "", "Target editor(s), comma-separated (e.g., windsurf,cursor)")
		cmd.Flags().StringVar(&syncExts, "ext", "", "Extension IDs to sync, comma-separated")
		cmd.Flags().BoolVar(&syncAll, "all", false, "Sync all extensions from source")
	}

	syncRunCmd.Flags().BoolVar(&syncOverwrite, "overwrite", false, "Overwrite existing extensions in target")
}
