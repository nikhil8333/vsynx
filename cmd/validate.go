package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/secureopenvsx/internal/models"
	"github.com/yourusername/secureopenvsx/internal/validation"
)

var validateCmd = &cobra.Command{
	Use:   "validate [extension-id]",
	Short: "Validate an extension by ID",
	Long: `Validate an extension by querying the Microsoft Marketplace and OpenVSX registry.
Compares metadata such as publisher, version, repository URL, and hash.
Classifies the extension as Legitimate, Suspicious, or Malicious.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		extensionID := args[0]

		validator := validation.NewValidator()
		result, err := validator.ValidateExtension(extensionID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error validating extension: %v\n", err)
			os.Exit(1)
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
		} else {
			printValidationResult(result)
		}

		// Exit with non-zero if malicious
		if result.TrustLevel == models.TrustLevelMalicious {
			os.Exit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func printValidationResult(result *models.ValidationResult) {
	fmt.Printf("\n=== Extension Validation: %s ===\n\n", result.ExtensionID)

	// Trust level with color
	trustColor := getTrustColor(result.TrustLevel)
	fmt.Printf("Trust Level: %s%s%s\n\n", trustColor, result.TrustLevel, colorReset)

	if result.Error != "" {
		fmt.Printf("Error: %s\n\n", result.Error)
	}

	// Marketplace data
	if result.MarketplaceData != nil {
		fmt.Println("Microsoft Marketplace:")
		fmt.Printf("  Publisher: %s\n", result.MarketplaceData.Publisher)
		fmt.Printf("  Name: %s\n", result.MarketplaceData.Name)
		fmt.Printf("  Version: %s\n", result.MarketplaceData.Version)
		fmt.Printf("  Repository: %s\n", result.MarketplaceData.RepositoryURL)
		fmt.Printf("  Download: %s\n\n", result.MarketplaceData.DownloadURL)
	}

	// OpenVSX data
	if result.OpenVSXData != nil {
		fmt.Println("OpenVSX Registry:")
		fmt.Printf("  Publisher: %s\n", result.OpenVSXData.Publisher)
		fmt.Printf("  Name: %s\n", result.OpenVSXData.Name)
		fmt.Printf("  Version: %s\n", result.OpenVSXData.Version)
		fmt.Printf("  Repository: %s\n", result.OpenVSXData.RepositoryURL)
		fmt.Printf("  Download: %s\n\n", result.OpenVSXData.DownloadURL)
	}

	// Differences
	if len(result.Differences) > 0 {
		fmt.Println("Differences Found:")
		for _, diff := range result.Differences {
			fmt.Printf("  - %s\n", diff)
		}
		fmt.Println()
	}

	// Recommendation
	fmt.Printf("Recommendation: %s\n\n", result.Recommendation)
}

// Color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorGreen  = "\033[32m"
)

func getTrustColor(trustLevel models.TrustLevel) string {
	switch trustLevel {
	case models.TrustLevelLegitimate:
		return colorGreen
	case models.TrustLevelSuspicious:
		return colorYellow
	case models.TrustLevelMalicious:
		return colorRed
	default:
		return ""
	}
}
