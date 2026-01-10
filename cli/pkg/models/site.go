package models

import "time"

// Database holds database connection info for a site
type Database struct {
	Name string `yaml:"name" validate:"required"`
	User string `yaml:"user" validate:"required"`
	Host string `yaml:"host" validate:"required"`
}

// Metadata holds additional site information
type Metadata struct {
	FreeSite      bool       `yaml:"free_site"`
	BackupEnabled bool       `yaml:"backup_enabled"`
	LastBackup    *time.Time `yaml:"last_backup,omitempty"`
}

// Site represents a WordPress site on a server
type Site struct {
	SystemName    string    `yaml:"system_name" validate:"required,alphanum"`
	PrimaryDomain string    `yaml:"primary_domain" validate:"required,fqdn"`
	CreatedAt     time.Time `yaml:"created_at"`
	AdminUser     string    `yaml:"admin_user" validate:"required"`
	AdminEmail    string    `yaml:"admin_email" validate:"required,email"`
	Domains       []Domain  `yaml:"domains"`
	Database      Database  `yaml:"database"`
	PHPVersion    string    `yaml:"php_version"`
	Metadata      Metadata  `yaml:"metadata"`
	Notes         string    `yaml:"notes,omitempty"`
}
