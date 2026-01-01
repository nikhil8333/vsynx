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

var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit all installed extensions",
	Long: `Scans all installed VS Code extensions and validates each one.
Provides a summary report showing trust levels and any issues found.`,
	Run: func(cmd *cobra.Command, args []string) {
		scanner := validation.NewScanner()
		
		report, err := scanner.AuditExtensions(extensionsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error auditing extensions: %v\n", err)
			os.Exit(1)
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(report, "", "  ")
			fmt.Println(string(data))
		} else {
			printAuditReport(report)
		}

		// Exit with non-zero if any malicious extensions found
		if report.MaliciousCount > 0 {
			os.Exit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(auditCmd)
}

func printAuditReport(report *models.AuditReport) {
	fmt.Printf("\n=== Extension Audit Report ===\n\n")
	fmt.Printf("Audit Time: %s\n", report.AuditTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Total Extensions: %d\n\n", report.TotalExtensions)

	// Summary
	fmt.Println("Summary:")
	fmt.Printf("  %sLegitimate: %d%s\n", colorGreen, report.LegitimateCount, colorReset)
	fmt.Printf("  %sSuspicious: %d%s\n", colorYellow, report.SuspiciousCount, colorReset)
	fmt.Printf("  %sMalicious: %d%s\n", colorRed, report.MaliciousCount, colorReset)
	fmt.Printf("  Unknown: %d\n\n", report.UnknownCount)

	// Show details for suspicious and malicious extensions
	if report.SuspiciousCount > 0 || report.MaliciousCount > 0 {
		fmt.Println("Extensions Requiring Attention:")
		fmt.Println(strings.Repeat("-", 80))

		for _, result := range report.Results {
			if result.TrustLevel == models.TrustLevelSuspicious || result.TrustLevel == models.TrustLevelMalicious {
				trustColor := getTrustColor(result.TrustLevel)
				fmt.Printf("\n%s [%s%s%s]\n", result.ExtensionID, trustColor, result.TrustLevel, colorReset)
				
				if len(result.Differences) > 0 {
					fmt.Println("  Issues:")
					for _, diff := range result.Differences {
						fmt.Printf("    - %s\n", diff)
					}
				}
				
				fmt.Printf("  Recommendation: %s\n", result.Recommendation)
			}
		}
		fmt.Println()
	}

	// All extensions list
	if report.LegitimateCount > 0 {
		fmt.Printf("\nLegitimate Extensions (%d):\n", report.LegitimateCount)
		for _, result := range report.Results {
			if result.TrustLevel == models.TrustLevelLegitimate {
				fmt.Printf("  ✓ %s\n", result.ExtensionID)
			}
		}
	}
}
