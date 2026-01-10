package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wordsail/cli/internal/config"
	"github.com/wordsail/cli/internal/installer"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize WordSail environment",
	Long: `Initialize WordSail by setting up configuration and copying Ansible playbooks.

This command will:
  1. Create ~/.wordsail/ directory structure
  2. Copy Ansible playbooks from the repository to ~/.wordsail/ansible/
  3. Create initial configuration file (servers.yaml)
  4. Validate the installation

Run this command once after installing WordSail.`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  WordSail Initialization")
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		// Check if already initialized
		if installer.IsInitialized() {
			color.Yellow("⚠️  WordSail is already initialized")
			fmt.Printf("    Ansible directory: %s\n", installer.GetAnsibleDir())
			fmt.Printf("    Config directory: %s\n", installer.GetWordsailDir())
			fmt.Println()

			// Check if config exists
			mgr, _ := config.NewManager()
			if mgr.ConfigExists() {
				fmt.Println("✓ Configuration file exists")
			} else {
				color.Yellow("⚠️  Configuration file not found")
				fmt.Println("    Creating default configuration...")
				if err := createDefaultConfig(mgr); err != nil {
					color.Red("✗ Failed to create configuration: %v", err)
					os.Exit(1)
				}
				color.Green("✓ Configuration file created")
			}

			fmt.Println()
			fmt.Println("To reinitialize, manually remove:", installer.GetWordsailDir())
			return
		}

		// Initialize
		fmt.Println("Initializing WordSail...")
		fmt.Println()

		// Step 1: Copy Ansible files
		fmt.Print("→ Copying Ansible playbooks... ")
		if err := installer.Initialize(); err != nil {
			color.Red("✗")
			color.Red("\nError: %v", err)
			os.Exit(1)
		}
		color.Green("✓")

		// Step 2: Create configuration
		fmt.Print("→ Creating configuration file... ")
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("✗")
			color.Red("\nError: %v", err)
			os.Exit(1)
		}

		if err := createDefaultConfig(mgr); err != nil {
			color.Red("✗")
			color.Red("\nError: %v", err)
			os.Exit(1)
		}
		color.Green("✓")

		// Step 3: Validate installation
		fmt.Print("→ Validating installation... ")
		if err := validateInstallation(); err != nil {
			color.Red("✗")
			color.Red("\nWarning: %v", err)
		} else {
			color.Green("✓")
		}

		// Success message
		fmt.Println()
		color.Green("═══════════════════════════════════════════════════════")
		color.Green("  ✓ WordSail initialized successfully!")
		color.Green("═══════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Installation paths:")
		fmt.Printf("  • Ansible:       %s\n", installer.GetAnsibleDir())
		fmt.Printf("  • Config:        %s\n", installer.GetWordsailDir())
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Edit ~/.wordsail/servers.yaml to configure global variables")
		fmt.Println("  2. Add a server:    wordsail server add")
		fmt.Println("  3. Provision:       wordsail server provision <name>")
		fmt.Println("  4. Create site:     wordsail site create")
		fmt.Println()
	},
}

func createDefaultConfig(mgr *config.Manager) error {
	ansiblePath := installer.GetAnsibleDir()

	cfg := config.DefaultConfig()
	cfg.Ansible.Path = ansiblePath

	return mgr.Save(cfg)
}

func validateInstallation() error {
	// Check if ansible directory has required files
	ansiblePath := installer.GetAnsibleDir()

	requiredFiles := []string{
		"provision.yml",
		"website.yml",
		"playbooks/domain_management.yml",
		"playbooks/delete_site.yml",
		"ansible.cfg",
	}

	for _, file := range requiredFiles {
		fullPath := fmt.Sprintf("%s/%s", ansiblePath, file)
		if _, err := os.Stat(fullPath); err != nil {
			return fmt.Errorf("missing required file: %s", file)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
