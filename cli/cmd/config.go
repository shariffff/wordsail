package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wordsail/cli/internal/config"
	"gopkg.in/yaml.v3"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage wordsail configuration",
	Long:  `Initialize, display, and validate the wordsail configuration file.`,
}

// configInitCmd represents the config init command
var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize wordsail configuration",
	Long:  `Create a new configuration file at ~/.wordsail/servers.yaml with default values.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if mgr.ConfigExists() {
			color.Yellow("Configuration file already exists at: %s", mgr.GetConfigPath())
			fmt.Println("To reinitialize, please delete the existing file first.")
			os.Exit(1)
		}

		if err := mgr.Initialize(); err != nil {
			color.Red("Error: Failed to initialize configuration: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Created %s directory", mgr.GetConfigDir())
		color.Green("✓ Generated default servers.yaml")
		color.Green("✓ Configuration initialized successfully")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Edit ~/.wordsail/servers.yaml with your settings")
		fmt.Println("  2. Add your first server: wordsail server add")
		fmt.Println("  3. Provision the server: wordsail server provision <name>")
	},
}

// configShowCmd represents the config show command
var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Display current configuration",
	Long:  `Display the contents of the wordsail configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found at: %s", mgr.GetConfigPath())
			fmt.Println("Run 'wordsail config init' to create it.")
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		// Marshal to YAML for pretty display
		data, err := yaml.Marshal(cfg)
		if err != nil {
			color.Red("Error: Failed to marshal configuration: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration file: %s\n\n", mgr.GetConfigPath())
		fmt.Println(string(data))
	},
}

// configValidateCmd represents the config validate command
var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long:  `Validate the wordsail configuration file for correctness and consistency.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found at: %s", mgr.GetConfigPath())
			fmt.Println("Run 'wordsail config init' to create it.")
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		validator := config.NewValidator()

		// Validate struct
		fmt.Println("Validating configuration structure...")
		if err := validator.ValidateStruct(cfg); err != nil {
			color.Red("✗ Structure validation failed: %v", err)
			os.Exit(1)
		}
		color.Green("✓ Structure validation passed")

		// Validate business rules
		fmt.Println("Validating business rules...")
		if err := validator.ValidateBusinessRules(cfg); err != nil {
			color.Red("✗ Business rules validation failed: %v", err)
			os.Exit(1)
		}
		color.Green("✓ Business rules validation passed")

		// Validate Ansible environment
		fmt.Println("Validating Ansible environment...")
		if err := validator.ValidateAnsibleEnvironment(cfg); err != nil {
			color.Red("✗ Ansible environment validation failed: %v", err)
			os.Exit(1)
		}
		color.Green("✓ Ansible environment validation passed")

		fmt.Println()
		color.Green("✓ Configuration is valid")
		fmt.Printf("  Servers: %d\n", len(cfg.Servers))

		totalSites := 0
		for _, server := range cfg.Servers {
			totalSites += len(server.Sites)
		}
		fmt.Printf("  Sites: %d\n", totalSites)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configValidateCmd)
}
