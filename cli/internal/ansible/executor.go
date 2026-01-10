package ansible

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/wordsail/cli/pkg/models"
)

// Executor handles Ansible playbook execution
type Executor struct {
	ansiblePath string
	invGenerator *InventoryGenerator
}

// NewExecutor creates a new Ansible executor
func NewExecutor(ansiblePath string) *Executor {
	return &Executor{
		ansiblePath:  ansiblePath,
		invGenerator: NewInventoryGenerator(),
	}
}

// ExecutePlaybook runs an ansible-playbook command with the given parameters
func (e *Executor) ExecutePlaybook(playbookName string, server models.Server, extraVars map[string]interface{}, globalVars map[string]interface{}) error {
	// Expand home directory in ansible path if needed
	ansiblePath := e.ansiblePath
	if strings.HasPrefix(ansiblePath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to expand home directory: %w", err)
		}
		ansiblePath = filepath.Join(homeDir, ansiblePath[1:])
	}

	// Generate inventory
	inventoryPath, err := e.invGenerator.Generate(server, fmt.Sprintf("wordsail %s", playbookName), globalVars)
	if err != nil {
		return fmt.Errorf("failed to generate inventory: %w", err)
	}
	defer e.invGenerator.Cleanup(inventoryPath)

	// Build playbook path
	playbookPath := filepath.Join(ansiblePath, playbookName)

	// Check if playbook exists
	if _, err := os.Stat(playbookPath); os.IsNotExist(err) {
		return fmt.Errorf("playbook not found: %s", playbookPath)
	}

	// Build command arguments
	args := []string{
		playbookPath,
		"-i", inventoryPath,
	}

	// Add extra vars if provided
	if len(extraVars) > 0 {
		varsJSON, err := json.Marshal(extraVars)
		if err != nil {
			return fmt.Errorf("failed to marshal extra vars: %w", err)
		}
		args = append(args, "--extra-vars", string(varsJSON))
	}

	// Create command
	cmd := exec.Command("ansible-playbook", args...)
	cmd.Dir = ansiblePath

	// Set environment variables
	cmd.Env = os.Environ()

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	fmt.Printf("\n")
	color.Cyan("Running: ansible-playbook %s", strings.Join(args, " "))
	fmt.Printf("\n")

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ansible-playbook: %w", err)
	}

	// Stream output in real-time
	done := make(chan bool)

	// Stream stdout
	go func() {
		e.streamOutput(stdout, false)
		done <- true
	}()

	// Stream stderr
	go func() {
		e.streamOutput(stderr, true)
	}()

	// Wait for stdout streaming to complete
	<-done

	// Wait for command to finish
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ansible-playbook failed: %w", err)
	}

	return nil
}

// streamOutput reads and prints output with color coding
func (e *Executor) streamOutput(reader io.Reader, isError bool) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()

		// Color code based on content
		if isError {
			color.Red(line)
		} else if strings.Contains(line, "FAILED") || strings.Contains(line, "fatal:") {
			color.Red(line)
		} else if strings.Contains(line, "ok:") || strings.Contains(line, "skipping:") {
			color.Green(line)
		} else if strings.Contains(line, "changed:") {
			color.Yellow(line)
		} else if strings.Contains(line, "PLAY [") || strings.Contains(line, "TASK [") {
			color.Cyan(line)
		} else if strings.Contains(line, "PLAY RECAP") {
			color.Magenta(line)
		} else {
			fmt.Println(line)
		}
	}
}
