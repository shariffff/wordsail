package config

import (
	"github.com/wordsail/cli/internal/installer"
	"github.com/wordsail/cli/pkg/models"
)

// AnsibleConfig holds Ansible-specific configuration
type AnsibleConfig struct {
	Path              string `yaml:"path" validate:"required"`
	RolesPath         string `yaml:"roles_path"`
	InventoryPath     string `yaml:"inventory_path"`
	PythonInterpreter string `yaml:"python_interpreter"`
}

// BackupConfig holds backup configuration (future use)
type BackupConfig struct {
	Enabled       bool   `yaml:"enabled"`
	Schedule      string `yaml:"schedule,omitempty"`
	RetentionDays int    `yaml:"retention_days,omitempty"`
	Destination   string `yaml:"destination,omitempty"`
}

// Config represents the main configuration file structure
type Config struct {
	Version         string                 `yaml:"version" validate:"required"`
	Ansible         AnsibleConfig          `yaml:"ansible"`
	GlobalVars      map[string]interface{} `yaml:"global_vars"`
	Servers         []models.Server        `yaml:"servers"`
	Backup          BackupConfig           `yaml:"backup,omitempty"`
	PreferredEditor string                 `yaml:"preferred_editor,omitempty"`
}

// DefaultConfig returns a new Config with sensible defaults
func DefaultConfig() *Config {
	// Detect ansible path dynamically
	ansiblePath, err := installer.GetAnsiblePath()
	if err != nil {
		// Fallback to user's wordsail directory
		ansiblePath = installer.GetAnsibleDir()
	}

	return &Config{
		Version: "1.0",
		Ansible: AnsibleConfig{
			Path:              ansiblePath,
			RolesPath:         "./roles",
			InventoryPath:     "/tmp/wordsail-inventory-{timestamp}.ini",
			PythonInterpreter: "/usr/bin/python3",
		},
		GlobalVars: map[string]interface{}{
			"certbot_email":    "admin@example.com",
			"wordsail_ssh_key": "~/.ssh/wordsail_rsa.pub",
		},
		Servers: []models.Server{},
		Backup: BackupConfig{
			Enabled: false,
		},
	}
}
