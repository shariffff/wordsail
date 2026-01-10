package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"github.com/wordsail/cli/pkg/models"
)

// TestSSHConnection tests SSH connectivity to a server
func TestSSHConnection(server models.Server) error {
	// Expand home directory in key file path
	keyFile := server.SSH.KeyFile
	if strings.HasPrefix(keyFile, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		keyFile = filepath.Join(homeDir, keyFile[1:])
	}

	// Read SSH private key
	key, err := os.ReadFile(keyFile)
	if err != nil {
		return fmt.Errorf("failed to read SSH key file %s: %w", keyFile, err)
	}

	// Parse private key
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return fmt.Errorf("failed to parse SSH private key: %w", err)
	}

	// Configure SSH client
	config := &ssh.ClientConfig{
		User: server.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Implement proper host key verification
		Timeout:         10 * time.Second,
	}

	// Connect to server
	addr := fmt.Sprintf("%s:%d", server.IP, server.SSH.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("SSH connection failed to %s: %w", addr, err)
	}
	defer client.Close()

	// Create session
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	// Test command execution
	output, err := session.CombinedOutput("echo 'wordsail-test'")
	if err != nil {
		return fmt.Errorf("test command failed: %w", err)
	}

	if strings.TrimSpace(string(output)) != "wordsail-test" {
		return fmt.Errorf("unexpected test output: %s", output)
	}

	return nil
}
