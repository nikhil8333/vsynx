package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yourusername/secureopenvsx/internal/editor"
)

var (
	installEditor string
	installCLI    string
	installFile   string
)

var installCmd = &cobra.Command{
	Use:   "install [extension-id...]",
	Short: "Install extensions via VS Code CLI",
	Long: `Installs extensions using the VS Code family CLI tools.
Supports installing to VS Code, VS Code Insiders, or VSCodium.

Note: This does NOT install from OpenVSX. Extensions are installed from
the Microsoft Marketplace via the official VS Code CLI.`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Determine which CLI to use
		cliCommand := installCLI
		if cliCommand == "" {
			// Map editor flag to CLI command
			switch installEditor {
			case "vscode", "":
				cliCommand = "code"
			case "vscode-insiders":
				cliCommand = "code-insiders"
			case "vscodium":
				cliCommand = "codium"
			default:
				fmt.Fprintf(os.Stderr, "Unknown editor: %s\n", installEditor)
				fmt.Fprintln(os.Stderr, "Supported: vscode, vscode-insiders, vscodium")
				os.Exit(1)
			}
		}

		// Check CLI availability
		cliStatus := editor.GetCLIStatus()
		cliAvailable := false
		switch cliCommand {
		case "code":
			cliAvailable = cliStatus.VSCodeAvailable
		case "code-insiders":
			cliAvailable = cliStatus.InsidersAvailable
		case "codium":
			cliAvailable = cliStatus.CodiumAvailable
		}

		if !cliAvailable {
			fmt.Fprintf(os.Stderr, "CLI command '%s' not found.\n", cliCommand)
			fmt.Fprintln(os.Stderr, "Make sure the editor is installed and its CLI is in your PATH.")
			os.Exit(1)
		}

		// Collect extension IDs
		var extensionIDs []string

		// From command line args
		extensionIDs = append(extensionIDs, args...)

		// From file
		if installFile != "" {
			file, err := os.Open(installFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
				os.Exit(1)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				// Skip empty lines and comments
				if line != "" && !strings.HasPrefix(line, "#") {
					extensionIDs = append(extensionIDs, line)
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
				os.Exit(1)
			}
		}

		if len(extensionIDs) == 0 {
			fmt.Fprintln(os.Stderr, "No extensions specified.")
			fmt.Fprintln(os.Stderr, "Usage: vsynx install <extension-id> [extension-id...]")
			fmt.Fprintln(os.Stderr, "   or: vsynx install --file extensions.txt")
			os.Exit(1)
		}

		// Install extensions
		report := editor.InstallMultipleExtensionsViaCLI(cliCommand, extensionIDs)

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(report, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Printf("\n=== Install Report ===\n")
		fmt.Printf("CLI: %s\n", cliCommand)
		fmt.Printf("Extensions: %d\n\n", len(extensionIDs))

		for _, result := range report.Results {
			if result.Success {
				fmt.Printf("%s✓%s %s\n", colorGreen, colorReset, result.ExtensionID)
			} else {
				fmt.Printf("%s✗%s %s - %s\n", colorRed, colorReset, result.ExtensionID, result.Error)
			}
		}

		fmt.Printf("\nSuccess: %d, Failed: %d\n", report.TotalSuccess, report.TotalFailed)

		if report.TotalFailed > 0 {
			os.Exit(1)
		}
	},
}

var installCLICmd = &cobra.Command{
	Use:   "install-cli",
	Short: "Install the vsynx CLI to your PATH",
	Long: `Installs the vsynx CLI tool to your system PATH so it can be run from any terminal.

On Windows: Adds the vsynx directory to your user PATH.
On macOS/Linux: Creates a symlink in /usr/local/bin or ~/.local/bin.`,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := installCLITool()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error installing CLI: %v\n", err)
			os.Exit(1)
		}

		if outputFormat == "json" {
			data, _ := json.MarshalIndent(result, "", "  ")
			fmt.Println(string(data))
			return
		}

		fmt.Println(result.Message)
		if result.NeedsRestart {
			fmt.Println("\nPlease restart your terminal for changes to take effect.")
		}
	},
}

type CLIInstallResult struct {
	Success      bool   `json:"success"`
	Message      string `json:"message"`
	Path         string `json:"path,omitempty"`
	NeedsRestart bool   `json:"needsRestart"`
}

func installCLITool() (*CLIInstallResult, error) {
	// This is a placeholder - actual implementation depends on OS
	// The real implementation would be in the backend and called from the GUI
	return &CLIInstallResult{
		Success:      true,
		Message:      "CLI installation should be done via Vsynx Manager GUI or platform installer.",
		NeedsRestart: false,
	}, nil
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(installCLICmd)

	installCmd.Flags().StringVar(&installEditor, "editor", "vscode", "Target editor (vscode, vscode-insiders, vscodium)")
	installCmd.Flags().StringVar(&installCLI, "cli", "", "CLI command to use directly (code, code-insiders, codium)")
	installCmd.Flags().StringVarP(&installFile, "file", "f", "", "File containing extension IDs (one per line)")
}
