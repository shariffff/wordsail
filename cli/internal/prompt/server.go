package prompt

import (
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/wordsail/cli/pkg/models"
)

// ServerInput holds the input for server creation
type ServerInput struct {
	Name     string
	Hostname string
	IP       string
	SSHUser  string
	SSHPort  int
	SSHKey   string
}

// PromptServerAdd prompts for server details
func PromptServerAdd() (*ServerInput, error) {
	input := &ServerInput{}

	// Server name
	namePrompt := &survey.Input{
		Message: "Server name:",
		Help:    "A friendly name to identify this server (e.g., production-1, staging-1)",
	}
	if err := survey.AskOne(namePrompt, &input.Name, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// Hostname
	hostnamePrompt := &survey.Input{
		Message: "Hostname or IP address:",
		Help:    "The server's hostname (e.g., server.example.com) or IP address",
	}
	if err := survey.AskOne(hostnamePrompt, &input.Hostname, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// IP address (with smart default if hostname looks like an IP)
	defaultIP := input.Hostname
	if net.ParseIP(defaultIP) == nil {
		defaultIP = ""
	}

	ipPrompt := &survey.Input{
		Message: "IP address:",
		Help:    "The server's IP address",
		Default: defaultIP,
	}
	if err := survey.AskOne(ipPrompt, &input.IP, survey.WithValidator(survey.Required), survey.WithValidator(validateIP)); err != nil {
		return nil, err
	}

	// SSH user
	userPrompt := &survey.Input{
		Message: "SSH user:",
		Default: "wordsail",
		Help:    "The SSH user to connect as",
	}
	if err := survey.AskOne(userPrompt, &input.SSHUser, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// SSH port
	portPrompt := &survey.Input{
		Message: "SSH port:",
		Default: "22",
		Help:    "The SSH port number",
	}
	var portStr string
	if err := survey.AskOne(portPrompt, &portStr, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}
	fmt.Sscanf(portStr, "%d", &input.SSHPort)
	if input.SSHPort == 0 {
		input.SSHPort = 22
	}

	// SSH key file
	homeDir, _ := os.UserHomeDir()
	defaultKeyPath := filepath.Join(homeDir, ".ssh", "wordsail_rsa")

	keyPrompt := &survey.Input{
		Message: "SSH private key file:",
		Default: defaultKeyPath,
		Help:    "Path to the SSH private key for authentication",
	}
	if err := survey.AskOne(keyPrompt, &input.SSHKey, survey.WithValidator(survey.Required)); err != nil {
		return nil, err
	}

	// Confirmation
	if err := confirmServerAdd(input); err != nil {
		return nil, err
	}

	return input, nil
}

// ToServer converts ServerInput to models.Server
func (si *ServerInput) ToServer() models.Server {
	return models.Server{
		Name:     si.Name,
		Hostname: si.Hostname,
		IP:       si.IP,
		SSH: models.SSHConfig{
			User:    si.SSHUser,
			Port:    si.SSHPort,
			KeyFile: si.SSHKey,
		},
		Status: "unprovisioned",
		Sites:  []models.Site{},
	}
}

func confirmServerAdd(input *ServerInput) error {
	fmt.Println("\nServer Configuration:")
	fmt.Printf("  Name:     %s\n", input.Name)
	fmt.Printf("  Hostname: %s\n", input.Hostname)
	fmt.Printf("  IP:       %s\n", input.IP)
	fmt.Printf("  SSH User: %s\n", input.SSHUser)
	fmt.Printf("  SSH Port: %d\n", input.SSHPort)
	fmt.Printf("  SSH Key:  %s\n", input.SSHKey)

	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: "Add this server?",
		Default: true,
	}

	if err := survey.AskOne(confirmPrompt, &confirm); err != nil {
		return err
	}

	if !confirm {
		return fmt.Errorf("server addition cancelled")
	}

	return nil
}

func validateIP(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("invalid type")
	}

	if net.ParseIP(str) == nil {
		return fmt.Errorf("invalid IP address format")
	}

	return nil
}
