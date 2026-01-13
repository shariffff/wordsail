package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

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

// siteCmd represents the site command
var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Manage WordPress sites",
	Long:  `Create, list, and delete WordPress sites on provisioned servers.`,
}

// siteCreateCmd represents the site create command
var siteCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"add"},
	Short:   "Create a new WordPress site",
	Long:    `Interactively create a new WordPress site on a provisioned server.`,
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

		// Check for non-interactive mode
		nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
		var input *prompt.SiteInput

		if nonInteractive {
			// Get values from flags
			serverName, _ := cmd.Flags().GetString("server")
			domain, _ := cmd.Flags().GetString("domain")
			systemName, _ := cmd.Flags().GetString("system-name")
			adminUser, _ := cmd.Flags().GetString("admin-user")
			adminEmail, _ := cmd.Flags().GetString("admin-email")
			adminPassword, _ := cmd.Flags().GetString("admin-password")

			if serverName == "" || domain == "" || systemName == "" || adminUser == "" || adminEmail == "" || adminPassword == "" {
				color.Red("Error: In non-interactive mode, all flags are required")
				fmt.Println("Required flags: --server, --domain, --system-name, --admin-user, --admin-email, --admin-password")
				os.Exit(1)
			}

			input = &prompt.SiteInput{
				ServerName:    serverName,
				Domain:        domain,
				SystemName:    systemName,
				AdminUser:     adminUser,
				AdminEmail:    adminEmail,
				AdminPassword: adminPassword,
			}
		} else {
			// Interactive prompts
			input, err = prompt.PromptSiteCreate(cfg.Servers)
			if err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}
		}

		// Find the target server
		var targetServer *models.Server
		for i := range cfg.Servers {
			if cfg.Servers[i].Name == input.ServerName {
				targetServer = &cfg.Servers[i]
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", input.ServerName)
			os.Exit(1)
		}

		if targetServer.Status != "provisioned" {
			color.Red("Error: Server '%s' is not provisioned", input.ServerName)
			fmt.Println("Provision the server first: wordsail server provision", input.ServerName)
			os.Exit(1)
		}

		// Prepare extra vars for Ansible
		extraVars := map[string]interface{}{
			"domain":            input.Domain,
			"system_name":       input.SystemName,
			"wp_admin_user":     input.AdminUser,
			"wp_admin_email":    input.AdminEmail,
			"wp_admin_password": input.AdminPassword,
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute website.yml playbook
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Creating WordPress site: %s", input.Domain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		if err := executor.ExecutePlaybook("website.yml", *targetServer, extraVars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Site creation failed: %v", err)
			os.Exit(1)
		}

		// Create site record
		now := time.Now()
		newSite := models.Site{
			SystemName:    input.SystemName,
			PrimaryDomain: input.Domain,
			CreatedAt:     now,
			AdminUser:     input.AdminUser,
			AdminEmail:    input.AdminEmail,
			Domains: []models.Domain{
				{
					Domain:     input.Domain,
					SSLEnabled: false, // SSL will be added separately
				},
			},
			Database: models.Database{
				Name: input.SystemName,
				User: input.SystemName,
				Host: "localhost",
			},
			PHPVersion: "8.3",
			Metadata: models.Metadata{
				BackupEnabled: false,
			},
		}

		// Add site to server configuration
		stateMgr := state.NewManager(mgr)
		if err := stateMgr.AddSiteToServer(input.ServerName, newSite); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		fmt.Println()
		color.Green("═══════════════════════════════════════════════════════")
		color.Green("  ✓ WordPress site created successfully!")
		color.Green("═══════════════════════════════════════════════════════")
		fmt.Println()
		fmt.Printf("Site URL:      http://%s\n", input.Domain)
		fmt.Printf("Admin URL:     http://%s/wp-admin\n", input.Domain)
		fmt.Printf("Admin User:    %s\n", input.AdminUser)
		fmt.Printf("Admin Email:   %s\n", input.AdminEmail)
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  1. Add www subdomain: wordsail domain add\n")
		fmt.Printf("  2. Issue SSL certificate: wordsail domain ssl\n")
	},
}

// SiteWithServer represents a site with its server name for JSON output
type SiteWithServer struct {
	ServerName string       `json:"server_name"`
	Site       models.Site  `json:"site"`
}

// siteListCmd represents the site list command
var siteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all WordPress sites",
	Long:  `Display all WordPress sites across all servers.`,
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

		// Filter by server if specified
		filterServer, _ := cmd.Flags().GetString("server")

		// Check for JSON output
		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			sites := make([]SiteWithServer, 0)
			for _, server := range cfg.Servers {
				if filterServer != "" && server.Name != filterServer {
					continue
				}
				for _, site := range server.Sites {
					sites = append(sites, SiteWithServer{
						ServerName: server.Name,
						Site:       site,
					})
				}
			}
			output, err := json.MarshalIndent(sites, "", "  ")
			if err != nil {
				color.Red("Error: Failed to marshal JSON: %v", err)
				os.Exit(1)
			}
			fmt.Println(string(output))
			return
		}

		// Count total sites
		totalSites := 0
		for _, server := range cfg.Servers {
			if filterServer != "" && server.Name != filterServer {
				continue
			}
			totalSites += len(server.Sites)
		}

		if totalSites == 0 {
			if filterServer != "" {
				fmt.Printf("No sites found on server '%s'\n", filterServer)
			} else {
				fmt.Println("No sites configured.")
				fmt.Println("Create a site with: wordsail site create")
			}
			return
		}

		// Display sites
		if filterServer != "" {
			fmt.Printf("\nSites on server '%s' (%d total):\n\n", filterServer, totalSites)
		} else {
			fmt.Printf("\nAll sites (%d total):\n\n", totalSites)
		}

		// Prepare table data
		headers := []string{"SERVER", "DOMAIN", "SYSTEM NAME", "NOTES"}
		colWidths := []int{20, 35, 20, 40}
		rows := make([][]string, 0)

		for _, server := range cfg.Servers {
			if filterServer != "" && server.Name != filterServer {
				continue
			}

			for _, site := range server.Sites {
				// Get notes (truncate if too long for display)
				notesStr := site.Notes
				if len(notesStr) > 38 {
					notesStr = notesStr[:35] + "..."
				}

				row := []string{
					server.Name,
					site.PrimaryDomain,
					site.SystemName,
					notesStr,
				}
				rows = append(rows, row)
			}
		}

		utils.PrintTableWithBorders(headers, rows, colWidths)
		fmt.Println()
	},
}

// siteDeleteCmd represents the site delete command
var siteDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove"},
	Short:   "Delete a WordPress site",
	Long:    `Delete a WordPress site and all its associated files and databases.`,
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

		// Get server and site from flags
		serverName, _ := cmd.Flags().GetString("server")
		siteName, _ := cmd.Flags().GetString("site")

		// If not provided, prompt interactively
		if serverName == "" || siteName == "" {
			// Build list of all sites
			type SiteOption struct {
				ServerName string
				Site       models.Site
			}

			var siteOptions []SiteOption
			for _, server := range cfg.Servers {
				for _, site := range server.Sites {
					siteOptions = append(siteOptions, SiteOption{
						ServerName: server.Name,
						Site:       site,
					})
				}
			}

			if len(siteOptions) == 0 {
				fmt.Println("No sites available to delete.")
				return
			}

			// Create selection options
			optionStrings := make([]string, len(siteOptions))
			for i, opt := range siteOptions {
				optionStrings[i] = fmt.Sprintf("%s on %s (%s)",
					opt.Site.PrimaryDomain, opt.ServerName, opt.Site.SystemName)
			}

			var selectedIndex int
			selectPrompt := &survey.Select{
				Message: "Select site to delete:",
				Options: optionStrings,
			}
			if err := survey.AskOne(selectPrompt, &selectedIndex); err != nil {
				color.Red("Error: %v", err)
				os.Exit(1)
			}

			serverName = siteOptions[selectedIndex].ServerName
			siteName = siteOptions[selectedIndex].Site.SystemName
		}

		// Find the server and site
		var targetServer *models.Server
		var targetSite *models.Site

		for i := range cfg.Servers {
			if cfg.Servers[i].Name == serverName {
				targetServer = &cfg.Servers[i]
				for j := range cfg.Servers[i].Sites {
					if cfg.Servers[i].Sites[j].SystemName == siteName {
						targetSite = &cfg.Servers[i].Sites[j]
						break
					}
				}
				break
			}
		}

		if targetServer == nil {
			color.Red("Error: Server '%s' not found", serverName)
			os.Exit(1)
		}

		if targetSite == nil {
			color.Red("Error: Site '%s' not found on server '%s'", siteName, serverName)
			os.Exit(1)
		}

		// Show warning and confirm
		color.Yellow("⚠️  WARNING: This will permanently delete:")
		fmt.Printf("  - Site: %s (%s)\n", targetSite.PrimaryDomain, targetSite.SystemName)
		fmt.Printf("  - Server: %s\n", serverName)
		fmt.Printf("  - All files in /home/%s\n", targetSite.SystemName)
		fmt.Printf("  - Database: %s\n", targetSite.Database.Name)
		fmt.Printf("  - Nginx configuration\n")
		fmt.Printf("  - PHP-FPM pool\n")
		fmt.Println()

		force, _ := cmd.Flags().GetBool("force")
		if !force {
			var confirm bool
			if err := survey.AskOne(&survey.Confirm{
				Message: "Are you absolutely sure you want to delete this site?",
				Default: false,
			}, &confirm); err != nil {
				os.Exit(1)
			}

			if !confirm {
				fmt.Println("Site deletion cancelled")
				return
			}

			// Double confirmation for safety
			var doubleConfirm string
			doublePrompt := &survey.Input{
				Message: fmt.Sprintf("Type '%s' to confirm deletion:", targetSite.SystemName),
			}
			if err := survey.AskOne(doublePrompt, &doubleConfirm); err != nil {
				os.Exit(1)
			}

			if doubleConfirm != targetSite.SystemName {
				color.Red("Confirmation failed. Site deletion cancelled.")
				return
			}
		}

		// Prepare extra vars for delete operation
		extraVars := map[string]interface{}{
			"system_name": targetSite.SystemName,
			"site_domain": targetSite.PrimaryDomain,
			"db_host":     targetSite.Database.Host,
		}

		// Create Ansible executor
		executor := ansible.NewExecutor(cfg.Ansible.Path)
		executor.SetVerbose(Verbose)
		executor.SetDryRun(DryRun)

		// Execute delete_site tasks
		fmt.Println()
		color.Cyan("═══════════════════════════════════════════════════════")
		color.Cyan("  Deleting site: %s", targetSite.PrimaryDomain)
		color.Cyan("═══════════════════════════════════════════════════════")
		fmt.Println()

		// Note: We need to create a playbook that includes the delete_site role
		// For now, we'll use a direct approach
		if err := executor.ExecutePlaybook("playbooks/delete_site.yml", *targetServer, extraVars, cfg.GlobalVars); err != nil {
			color.Red("\n✗ Site deletion failed: %v", err)
			color.Yellow("Note: You may need to manually clean up resources on the server")
			os.Exit(1)
		}

		// Remove site from configuration
		stateMgr := state.NewManager(mgr)
		if err := stateMgr.RemoveSiteFromServer(serverName, siteName); err != nil {
			color.Red("Warning: Failed to update configuration: %v", err)
		}

		fmt.Println()
		color.Green("✓ Site '%s' deleted successfully", targetSite.PrimaryDomain)
	},
}

func init() {
	rootCmd.AddCommand(siteCmd)
	siteCmd.AddCommand(siteCreateCmd)
	siteCmd.AddCommand(siteListCmd)
	siteCmd.AddCommand(siteDeleteCmd)

	// site create flags
	siteCreateCmd.Flags().Bool("non-interactive", false, "Use flags instead of interactive prompts")
	siteCreateCmd.Flags().String("server", "", "Target server name")
	siteCreateCmd.Flags().String("domain", "", "Primary domain")
	siteCreateCmd.Flags().String("system-name", "", "System username")
	siteCreateCmd.Flags().String("admin-user", "", "WordPress admin username")
	siteCreateCmd.Flags().String("admin-email", "", "WordPress admin email")
	siteCreateCmd.Flags().String("admin-password", "", "WordPress admin password")

	// site create json flag
	siteCreateCmd.Flags().Bool("json", false, "Output in JSON format")

	// site list flags
	siteListCmd.Flags().String("server", "", "Filter by server name")
	siteListCmd.Flags().Bool("json", false, "Output in JSON format")

	// site delete flags
	siteDeleteCmd.Flags().String("server", "", "Server name")
	siteDeleteCmd.Flags().String("site", "", "Site system name")
	siteDeleteCmd.Flags().BoolP("force", "f", false, "Force deletion without confirmation")
	siteDeleteCmd.Flags().Bool("json", false, "Output in JSON format")
}
