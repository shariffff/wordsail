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
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "wordsail",
	Short: "WordSail - Ansible wrapper for WordPress hosting management",
	Long: `WordSail is a CLI tool that simplifies WordPress hosting management
by wrapping Ansible playbooks with an intuitive, interactive interface.

Manage servers, sites, and domains with ease while maintaining full
visibility into your infrastructure state via ~/.wordsail/servers.yaml`,
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.wordsail/servers.yaml)")
}
