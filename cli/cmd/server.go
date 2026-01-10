package cmd

import (
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wordsail/cli/internal/ansible"
	"github.com/wordsail/cli/internal/config"
	"github.com/wordsail/cli/internal/prompt"
	"github.com/wordsail/cli/internal/state"
	"github.com/wordsail/cli/internal/utils"
	"github.com/wordsail/cli/pkg/models"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Manage servers",
	Long:  `Add, list, remove, and provision servers.`,
}

// serverAddCmd represents the server add command
var serverAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new server",
	Long:  `Interactively add a new server to the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found. Run 'wordsail config init' first.")
			os.Exit(1)
		}

		// Load existing config
		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		// Prompt for server details
		input, err := prompt.PromptServerAdd()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		// Check for duplicate server name
		for _, server := range cfg.Servers {
			if server.Name == input.Name {
				color.Red("Error: Server with name '%s' already exists", input.Name)
				os.Exit(1)
			}
		}

		// Add server to config
		newServer := input.ToServer()
		cfg.Servers = append(cfg.Servers, newServer)

		// Save config
		if err := mgr.Save(cfg); err != nil {
			color.Red("Error: Failed to save configuration: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Server '%s' added successfully", input.Name)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  Provision the server: wordsail server provision %s\n", input.Name)
	},
}

// serverListCmd represents the server list command
var serverListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all servers",
	Long:  `Display all servers in the configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found. Run 'wordsail config init' first.")
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		if len(cfg.Servers) == 0 {
			fmt.Println("No servers configured.")
			fmt.Println("Add a server with: wordsail server add")
			return
		}

		fmt.Printf("\nServers (%d total):\n\n", len(cfg.Servers))

		// Prepare table data
		headers := []string{"NAME", "HOSTNAME", "IP", "SSH USER", "STATUS", "SITES"}
		colWidths := []int{18, 28, 15, 12, 15, 6}
		rows := make([][]string, 0)

		for _, server := range cfg.Servers {
			statusStr := ""
			switch server.Status {
			case "provisioned":
				statusStr = color.GreenString(server.Status)
			case "unprovisioned":
				statusStr = color.YellowString(server.Status)
			case "error":
				statusStr = color.RedString(server.Status)
			default:
				statusStr = server.Status
			}

			row := []string{
				server.Name,
				server.Hostname,
				server.IP,
				server.SSH.User,
				statusStr,
				fmt.Sprintf("%d", len(server.Sites)),
			}
			rows = append(rows, row)
		}

		utils.PrintTableWithBorders(headers, rows, colWidths)
	},
}

// serverRemoveCmd represents the server remove command
var serverRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a server",
	Long:  `Remove a server from the configuration by name.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found. Run 'wordsail config init' first.")
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		// Find and remove server
		found := false
		newServers := make([]models.Server, 0)
		var removedServer models.Server

		for _, server := range cfg.Servers {
			if server.Name == serverName {
				found = true
				removedServer = server
			} else {
				newServers = append(newServers, server)
			}
		}

		if !found {
			color.Red("Error: Server '%s' not found", serverName)
			os.Exit(1)
		}

		// Warn if server has sites
		if len(removedServer.Sites) > 0 {
			color.Yellow("Warning: Server '%s' has %d site(s)", serverName, len(removedServer.Sites))
			fmt.Println("Removing the server will also remove all site records.")

			force, _ := cmd.Flags().GetBool("force")
			if !force {
				var confirm bool
				if err := survey.AskOne(&survey.Confirm{
					Message: "Are you sure you want to remove this server?",
					Default: false,
				}, &confirm); err != nil {
					os.Exit(1)
				}

				if !confirm {
					fmt.Println("Server removal cancelled")
					return
				}
			}
		}

		cfg.Servers = newServers

		// Save config
		if err := mgr.Save(cfg); err != nil {
			color.Red("Error: Failed to save configuration: %v", err)
			os.Exit(1)
		}

		color.Green("✓ Server '%s' removed successfully", serverName)
	},
}

// serverProvisionCmd represents the server provision command
var serverProvisionCmd = &cobra.Command{
	Use:   "provision <name>",
	Short: "Provision a server",
	Long:  `Run the provision.yml playbook to set up a server with Nginx, PHP, MariaDB, and security hardening.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		serverName := args[0]

		mgr, err := config.NewManager()
		if err != nil {
			color.Red("Error: %v", err)
			os.Exit(1)
		}

		if !mgr.ConfigExists() {
			color.Red("Configuration file not found. Run 'wordsail config init' first.")
			os.Exit(1)
		}

		cfg, err := mgr.Load()
		if err != nil {
			color.Red("Error: Failed to load configuration: %v", err)
			os.Exit(1)
		}

		// Find the server
		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == serverName {
				targetServer = &cfg.Servers[i]
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", serverName)
			fmt.Println("Run 'wordsail server list' to see available servers.")
			os.Exit(1)
		}

		// Check if already provisioned
		if targetServer.Status == "provisioned" {
			color.Yellow("Warning: Server '%s' is already marked as provisioned", serverName)

			skipCheck, _ := cmd.Flags().GetBool("skip-check")
			if !skipCheck {
				var confirm bool
				if err := survey.AskOne(&survey.Confirm{
					Message: "Provision again anyway?",
					Default: false,
				}, &confirm); err != nil {
					os.Exit(1)
				}

				if !confirm {
					fmt.Println("Provisioning cancelled")
					return
				}
			}
		}

		// Pre-flight SSH check
		skipSSH, _ := cmd.Flags().GetBool("skip-ssh-check")
		if !skipSSH {
			fmt.Println("Checking SSH connectivity...")
			if err := utils.TestSSHConnection(*targetServer); err != nil {
				color.Red("✗ SSH connectivity check failed: %v", err)
				fmt.Println()
				fmt.Println("Please verify:")
				fmt.Println("  1. Server is reachable")
				fmt.Println("  2. SSH key file exists and has correct permissions")
				fmt.Println("  3. SSH user has access to the server")
				fmt.Println()
				fmt.Println("Use --skip-ssh-check to bypass this check (not recommended)")
				os.Exit(1)
			}
			color.Green("✓ SSH connectivity check passed")
			fmt.Println()
		}

		// Confirm provisioning
		color.Cyan("About to provision server: %s (%s)", targetServer.Name, targetServer.IP)
		fmt.Println("This will:")
		fmt.Println("  - Install Nginx, PHP 8.3, MariaDB")
		fmt.Println("  - Configure security (UFW, Fail2ban, SSH hardening)")
		fmt.Println("  - Set up Certbot for SSL certificates")
		fmt.Println("  - Create wordsail user and environment")
		fmt.Println()

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			var confirm bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Continue with provisioning?",
				Default: true,
			}, &confirm); err != nil {
				os.Exit(1)
			}

			if !confirm {
				fmt.Println("Provisioning cancelled")
				return
			}
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)

		// Execute provision.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Starting provisioning: %s", serverName)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("provision.yml", *targetServer, nil, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Provisioning failed: %v", err)

			// Mark server as error
			stateMgr := state.NewManager(mgr)
			stateMgr.MarkServerError(serverName)

			os.Exit(1)
		}

		// Update server status to provisioned
		stateMgr := state.NewManager(mgr)
		if err := stateMgr.MarkServerProvisioned(serverName); err != nil {
			color.Red("Warning: Failed to update server status: %v", err)
		}

		fmt.Println()
		color.Green("═══════════════════════════════════════════════════════")
		color.Green("  ✓ Server '%s' provisioned successfully!", serverName)
		color.Green("═══════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  Create a WordPress site: wordsail site create")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverAddCmd)
	serverCmd.AddCommand(serverListCmd)
	serverCmd.AddCommand(serverRemoveCmd)
	serverCmd.AddCommand(serverProvisionCmd)

	// Flags
	serverRemoveCmd.Flags().BoolP("force", "f", false, "Force removal without confirmation")
	serverProvisionCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
	serverProvisionCmd.Flags().Bool("skip-ssh-check", false, "Skip SSH connectivity check")
	serverProvisionCmd.Flags().Bool("skip-check", false, "Skip already-provisioned check")
}
