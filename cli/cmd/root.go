package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "dev"
	CommitSHA = "unknown"
	BuildDate = "unknown"

	// Global flags
	Verbose bool
	DryRun  bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "wordsail",
	Short: "WordSail - Ansible wrapper for WordPress hosting management",
	Long: `WordSail is a CLI tool that simplifies WordPress hosting management
by wrapping Ansible playbooks with an intuitive, interactive interface.

Manage servers, sites, and domains with ease while maintaining full
visibility into your infrastructure state via ~/.wordsail/servers.yaml

Examples:
  # Initialize configuration
  wordsail config init

  # Add and provision a new server
  wordsail server provision

  # Create a WordPress site
  wordsail site create

  # Add a domain with SSL
  wordsail domain add

  # List all servers
  wordsail server list --json`,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&DryRun, "dry-run", false, "Show what would be done without making changes")
}
