# WordSail CLI

A command-line interface tool for managing WordPress hosting infrastructure using Ansible. WordSail provides both an intuitive, interactive mode for manual operations and a script mode with flags for automation and CI/CD pipelines.

## Features

- **Dual Operating Modes**: Interactive prompts for manual use, script mode with flags for automation
- **YAML-based State Management**: All configuration stored in `~/.wordsail/servers.yaml`
- **Ansible Integration**: Seamlessly executes existing Ansible playbooks
- **Server Management**: Add, list, remove, and provision servers
- **Site Management**: Create, list, and delete WordPress sites
- **Domain Management**: Add domains and manage SSL certificates

## Installation

### Prerequisites

- Go 1.21 or higher
- Ansible installed and configured
- SSH access to target servers

### Build from Source

```bash
# Clone or navigate to the repository
cd /path/to/ansible/cli

# Build the binary
make build

# Install to /usr/local/bin (requires sudo)
make install

# Or install to ~/bin (no sudo required)
make install-user
```

### Verify Installation

```bash
wordsail version
wordsail --help
```

## Quick Start

### 1. Initialize Configuration

```bash
wordsail config init
```

This creates `~/.wordsail/servers.yaml` with default settings.

### 2. Configure Ansible Path

Edit `~/.wordsail/servers.yaml` and set the correct Ansible project path:

```yaml
ansible:
  path: '/Users/yourname/Projects/ansible' # Update this path
  roles_path: './roles'
  inventory_path: '/tmp/wordsail-inventory-{timestamp}.ini'
  python_interpreter: '/usr/bin/python3'
```

### 3. Add a Server

```bash
wordsail server add
```

Follow the interactive prompts to add server details:

- Server name (e.g., production-1)
- Hostname or IP address
- SSH user and port
- SSH private key file

### 4. List Servers

```bash
wordsail server list
```

### 5. Validate Configuration

```bash
wordsail config validate
```

## Operating Modes

WordSail CLI supports two modes of operation:

### Interactive Mode (Default)

When you run commands without flags, the CLI guides you through the process with interactive prompts.

```bash
wordsail site create
# Prompts you for: server, domain, system name, admin credentials
```

**Use interactive mode when:**
- Learning the tool
- Performing manual operations
- You want validation and helpful hints
- Exploring available options

### Script Mode (Non-Interactive)

Provide all parameters as command-line flags for fully automated operations.

```bash
wordsail site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --system-name examplecom \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!
```

**Use script mode when:**
- Automating deployments
- Running in CI/CD pipelines
- Scripting repetitive tasks
- No user interaction is possible

**Common flags for script mode:**
- `--non-interactive`: Required flag to enable script mode
- `--force`: Skip confirmation prompts
- `--skip-ssh-check`: Skip SSH connectivity validation

## Commands

### Configuration Management

```bash
# Initialize configuration
wordsail config init

# Show current configuration
wordsail config show

# Validate configuration
wordsail config validate
```

### Server Management

```bash
# Add a new server
wordsail server add

# List all servers
wordsail server list

# Remove a server
wordsail server remove <name>

# Provision a server
wordsail server provision <name>

# Provision with options
wordsail server provision <name> --force              # Skip confirmation
wordsail server provision <name> --skip-ssh-check     # Skip SSH connectivity test
```

### Site Management

```bash
# Create a new WordPress site (interactive)
wordsail site create

# Create a site non-interactively
wordsail site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --system-name examplecom \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!

# List all sites
wordsail site list

# List sites on a specific server
wordsail site list --server production-1

# Delete a site (interactive selection)
wordsail site delete

# Delete a specific site
wordsail site delete --server production-1 --site examplecom

# Force delete without confirmation
wordsail site delete --server production-1 --site examplecom --force
```

### Domain Management

```bash
# Add a domain to a site (interactive)
wordsail domain add

# Add domain with automatic SSL
# (prompts will ask if you want to issue SSL)

# Remove a domain (interactive selection)
wordsail domain remove

# Force remove without confirmation
wordsail domain remove --force

# Issue SSL certificate for a domain (interactive)
wordsail domain ssl

# The CLI will:
# - Show only domains without SSL
# - Prompt for Let's Encrypt email
# - Obtain and configure SSL certificate
# - Update Nginx to use HTTPS
# - Track SSL expiration in configuration
```

## Configuration File

The configuration file is located at `~/.wordsail/servers.yaml`. Here's an example structure:

```yaml
version: '1.0'

ansible:
  path: '/Users/sharif/Projects/ansible'
  roles_path: './roles'
  inventory_path: '/tmp/wordsail-inventory-{timestamp}.ini'
  python_interpreter: '/usr/bin/python3'

global_vars:
  certbot_email: 'admin@example.com'
  mysql_wordsailbot_password: '${MYSQL_WORDSAILBOT_PASSWORD}'
  wordsail_ssh_key: '~/.ssh/wordsail_rsa.pub'

servers:
  - name: 'production-1'
    hostname: 'prod1.example.com'
    ip: '203.0.113.10'
    ssh:
      user: 'wordsail'
      port: 22
      key_file: '~/.ssh/wordsail_rsa'
    status: 'unprovisioned'
    sites: []
```

## Development

### Build

```bash
make build
```

### Test

```bash
make test
```

### Format Code

```bash
make fmt
```

### Clean Build Artifacts

```bash
make clean
```

## Project Structure

```
cli/
├── cmd/                  # Command definitions
│   ├── root.go          # Root command
│   ├── config.go        # Config commands
│   ├── server.go        # Server commands
│   └── version.go       # Version command
├── internal/
│   ├── config/          # Configuration management
│   ├── ansible/         # Ansible integration (coming soon)
│   ├── state/           # State management (coming soon)
│   ├── prompt/          # Interactive prompts
│   └── utils/           # Utilities (coming soon)
├── pkg/
│   └── models/          # Data models
├── templates/           # Templates (inventory, etc.)
├── main.go             # Entry point
├── Makefile            # Build automation
└── README.md           # This file
```

## Roadmap

- [ ] Shell completion scripts
- [ ] Comprehensive error handling
- [ ] Installation script
- [ ] Release automation
