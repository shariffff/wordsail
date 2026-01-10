package prompt

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/wordsail/cli/pkg/models"
)

// SiteInput holds the input for site creation
type SiteInput struct {
	ServerName    string
	Domain        string
	SystemName    string
	AdminUser     string
	AdminEmail    string
	AdminPassword string
	FreeSite      bool
}

// PromptSiteCreate prompts for site creation details
func PromptSiteCreate(servers []models.Server) (*SiteInput, error) {
	input := &SiteInput{}

	if len(servers) == 0 {
		return nil, fmt.Errorf("no servers available. Add a server first with: wordsail server add")
	}

	// Filter only provisioned servers
	provisionedServers := make([]models.Server, 0)
	for _, s := range servers {
		if s.Status == "provisioned" {
			provisionedServers = append(provisionedServers, s)
		}
	}

	if len(provisionedServers) == 0 {
		return nil, fmt.Errorf("no provisioned servers available. Provision a server first with: wordsail server provision <name>")
	}

	// 1. Select server
	serverOptions := make([]string, len(provisionedServers))
	for i, s := range provisionedServers {
		serverOptions[i] = fmt.Sprintf("%s (%s) - %d sites", s.Name, s.IP, len(s.Sites))
	}

	var serverIndex int
	serverPrompt := &survey.Select{
		Message: "Select target server:",
		Options: serverOptions,
		Help:    "Choose a provisioned server to host this WordPress site",
	}
	if err := survey.AskOne(serverPrompt, &serverIndex); err != nil {
		return nil, err
	}
	input.ServerName = provisionedServers[serverIndex].Name

	// 2. Domain name
	domainPrompt := &survey.Input{
		Message: "Primary domain name:",
		Help:    "The main domain for this WordPress site (e.g., example.com)",
	}
	if err := survey.AskOne(domainPrompt, &input.Domain, survey.WithValidator(survey.Required), survey.WithValidator(validateDomain)); err != nil {
		return nil, err
	}

	// 3. System name (with smart default)
	defaultSystemName := generateSystemName(input.Domain)
	systemNamePrompt := &survey.Input{
		Message: "System username:",
		Help:    "Linux user for this site (alphanumeric only, 3-16 chars)",
		Default: defaultSystemName,
	}
	if err := survey.AskOne(systemNamePrompt, &input.SystemName, survey.WithValidator(survey.Required), survey.WithValidator(validateSystemName)); err != nil {
		return nil, err
	}

	// Check if system name already exists on this server
	selectedServer := provisionedServers[serverIndex]
	for _, site := range selectedServer.Sites {
		if site.SystemName == input.SystemName {
			return nil, fmt.Errorf("system name '%s' already exists on server '%s'", input.SystemName, input.ServerName)
		}
	}

	// 4. WordPress admin user
	adminUserPrompt := &survey.Input{
		Message: "WordPress admin username:",
		Default: "admin",
		Help:    "Username for WordPress admin account",
	}
	if err := survey.AskOne(adminUserPrompt, &input.AdminUser, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// 5. WordPress admin email
	adminEmailPrompt := &survey.Input{
		Message: "WordPress admin email:",
		Help:    "Email address for WordPress admin account",
	}
	if err := survey.AskOne(adminEmailPrompt, &input.AdminEmail, survey.WithValidator(survey.Required), survey.WithValidator(validateEmail)); err != nil {
		return nil, err
	}

	// 6. WordPress admin password (with option to generate)
	var useGeneratedPassword bool
	generatePrompt := &survey.Confirm{
		Message: "Generate secure password?",
		Default: true,
		Help:    "Auto-generate a strong password or enter your own",
	}
	if err := survey.AskOne(generatePrompt, &useGeneratedPassword); err != nil {
		return nil, err
	}

	if useGeneratedPassword {
		input.AdminPassword = generateSecurePassword(20)
		fmt.Printf("\n")
		fmt.Printf("Generated password: %s\n", input.AdminPassword)
		fmt.Printf("⚠️  IMPORTANT: Save this password securely!\n")
		fmt.Printf("\n")

		var acknowledged bool
		ackPrompt := &survey.Confirm{
			Message: "Have you saved the password?",
			Default: false,
		}
		if err := survey.AskOne(ackPrompt, &acknowledged); err != nil {
			return nil, err
		}
		if !acknowledged {
			return nil, fmt.Errorf("please save the password before continuing")
		}
	} else {
		passwordPrompt := &survey.Password{
			Message: "WordPress admin password:",
			Help:    "Minimum 12 characters recommended",
		}
		if err := survey.AskOne(passwordPrompt, &input.AdminPassword, survey.WithValidator(survey.Required), survey.WithValidator(validatePasswordStrength)); err != nil {
			return nil, err
		}
	}

	// 7. Free site flag
	freeSitePrompt := &survey.Confirm{
		Message: "Mark as free site?",
		Default: false,
		Help:    "Free sites may have different resource limits or management",
	}
	if err := survey.AskOne(freeSitePrompt, &input.FreeSite); err != nil {
		return nil, err
	}

	// 8. Confirmation
	if err := confirmSiteCreation(input); err != nil {
		return nil, err
	}

	return input, nil
}

// confirmSiteCreation shows a summary and asks for confirmation
func confirmSiteCreation(input *SiteInput) error {
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println("Site Configuration Summary:")
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Printf("  Server:       %s\n", input.ServerName)
	fmt.Printf("  Domain:       %s\n", input.Domain)
	fmt.Printf("  System Name:  %s\n", input.SystemName)
	fmt.Printf("  Admin User:   %s\n", input.AdminUser)
	fmt.Printf("  Admin Email:  %s\n", input.AdminEmail)
	fmt.Printf("  Free Site:    %v\n", input.FreeSite)
	fmt.Println("═══════════════════════════════════════════════════")
	fmt.Println()

	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "Create this WordPress site?",
		Default: true,
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("site creation cancelled")
	}

	return nil
}

// Validators

func validateDomain(val interface{}) error {
	domain, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid domain type")
	}

	// Basic domain validation
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)
	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain format (e.g., example.com)")
	}

	return nil
}

func validateSystemName(val interface{}) error {
	name, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid system name type")
	}

	// Alphanumeric only, 3-16 characters
	if len(name) < 3 || len(name) > 16 {
		return fmt.Errorf("system name must be 3-16 characters")
	}

	alphanumRegex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumRegex.MatchString(name) {
		return fmt.Errorf("system name must be alphanumeric only")
	}

	return nil
}

func validateEmail(val interface{}) error {
	email, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid email type")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

func validatePasswordStrength(val interface{}) error {
	password, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid password type")
	}

	if len(password) < 12 {
		return fmt.Errorf("password must be at least 12 characters")
	}

	return nil
}

// Helper functions

func generateSystemName(domain string) string {
	// Remove common TLDs
	tlds := []string{".com", ".net", ".org", ".io", ".co", ".dev", ".app"}
	name := domain
	for _, tld := range tlds {
		name = strings.TrimSuffix(name, tld)
	}

	// Remove all non-alphanumeric characters
	alphanumRegex := regexp.MustCompile(`[^a-zA-Z0-9]`)
	name = alphanumRegex.ReplaceAllString(name, "")

	// Limit to 16 characters
	if len(name) > 16 {
		name = name[:16]
	}

	// Ensure at least 3 characters
	if len(name) < 3 {
		name = "site" + name
	}

	return strings.ToLower(name)
}

func generateSecurePassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// Fallback to a simple method if crypto/rand fails
			password[i] = charset[i%len(charset)]
		} else {
			password[i] = charset[num.Int64()]
		}
	}

	return string(password)
}
