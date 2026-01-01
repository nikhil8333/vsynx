package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yourusername/secureopenvsx/internal/validation"
)

var (
	outputDir string
)

var downloadCmd = &cobra.Command{
	Use:   "download [extension-id]",
	Short: "Download official extension from Microsoft Marketplace",
	Long: `Downloads the official VSIX package for an extension from the Microsoft Marketplace.
The package is verified and can be used to replace a suspicious or malicious extension.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		extensionID := args[0]

		validator := validation.NewValidator()

		fmt.Printf("Downloading official extension: %s\n", extensionID)

		data, hash, err := validator.DownloadOfficialExtension(extensionID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error downloading extension: %v\n", err)
			os.Exit(1)
		}

		// Determine output path
		if outputDir == "" {
			outputDir = "."
		}

		filename := fmt.Sprintf("%s.vsix", extensionID)
		outputPath := filepath.Join(outputDir, filename)

		// Write to file
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\n✓ Successfully downloaded extension\n")
		fmt.Printf("  File: %s\n", outputPath)
		fmt.Printf("  Size: %d bytes\n", len(data))
		fmt.Printf("  SHA256: %s\n\n", hash)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVarP(&outputDir, "output-dir", "d", ".", "Output directory for downloaded extension")
}
