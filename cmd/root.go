package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Flags
	extensionsPath string
	outputFormat   string
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "vsynx",
	Short: "vsynx - Secure VS Code extension manager",
	Long: `vsynx is a cross-platform CLI for securely managing VS Code extensions.
It validates extension identities, classifies trust levels, syncs extensions
between editors, and replaces unverified extensions with official ones.

Use 'vsynx gui' to launch the graphical interface (Vsynx Manager).`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&extensionsPath, "path", "p", "", "Path to extensions directory")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text, json)")
}
