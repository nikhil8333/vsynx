package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/secureopenvsx/internal/marketplace"
)

var (
	searchLimit int
)

var marketplaceCmd = &cobra.Command{
	Use:   "marketplace",
	Short: "Search and explore the Microsoft Marketplace",
	Long:  `Commands for searching and exploring extensions on the Microsoft Visual Studio Marketplace.`,
}

var marketplaceSearchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Search for extensions",
	Long:  `Searches the Microsoft Marketplace for extensions matching the given query.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := args[0]

		client := marketplace.NewClient()
		results, err := client.SearchExtensions(query)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching marketplace: %v\n", err)
			os.Exit(1)
		}

		// Apply limit
		if searchLimit > 0 && len(results) > searchLimit {
			results = results[:searchLimit]
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(results, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n=== Marketplace Search: \"%s\" (%d results) ===\n\n", query, len(results))

		if len(results) == 0 {
			fmt.Println("No extensions found.")
			return
		}

		fmt.Printf("%-45s %-25s %-12s %s\n", "Extension ID", "Publisher", "Version", "Verified")
		fmt.Println(strings.Repeat("-", 100))

		for _, ext := range results {
			verified := ""
			if ext.IsVerifiedPublisher {
				verified = "✓"
			}
			fmt.Printf("%-45s %-25s %-12s %s\n", ext.ID, ext.Publisher, ext.Version, verified)
		}
		fmt.Println()
	},
}

var marketplaceOpenCmd = &cobra.Command{
	Use:   "open <extension-id>",
	Short: "Get marketplace URL for an extension",
	Long:  `Prints the Microsoft Marketplace URL for an extension. Use --browser to open it directly.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		extensionID := args[0]
		url := fmt.Sprintf("https://marketplace.visualstudio.com/items?itemName=%s", extensionID)

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(map[string]string{
				"extensionId": extensionID,
				"url":         url,
			}, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Println(url)
	},
}

func init() {
	rootCmd.AddCommand(marketplaceCmd)
	marketplaceCmd.AddCommand(marketplaceSearchCmd)
	marketplaceCmd.AddCommand(marketplaceOpenCmd)

	marketplaceSearchCmd.Flags().IntVarP(&searchLimit, "limit", "l", 20, "Maximum number of results to display")
}
