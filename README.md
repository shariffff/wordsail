# WordSail

**Automated WordPress hosting infrastructure and management CLI**

WordSail provides a complete solution for deploying and managing WordPress sites on Ubuntu servers. It combines powerful Ansible playbooks with an intuitive CLI tool for streamlined operations.

## Features

- 🚀 **One-command server provisioning** - Full LEMP stack setup (Nginx, PHP 8.3, MariaDB)
- 🔒 **Security hardened** - UFW firewall, Fail2ban, SSH hardening, automatic SSL with Let's Encrypt
- 🎯 **Isolated WordPress sites** - Each site runs under its own user with dedicated PHP-FPM pool
- 🛠️ **Dual-mode CLI** - Interactive prompts for manual operations, script mode with flags for automation
- 🔄 **Infrastructure as code** - All configuration reproducible via Ansible

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/shariffff/wordsail.git
cd wordsail

# Build and install the CLI
make install

# Initialize WordSail (copies Ansible playbooks to ~/.wordsail/)
wordsail init
```

### Usage

**Interactive Mode** (guided prompts):
```bash
# Add and provision a server
wordsail server add
wordsail server provision production-1

# Create a WordPress site
wordsail site create

# Manage domains
wordsail domain add
wordsail domain ssl

# List everything
wordsail server list
wordsail site list
```

**Script Mode** (non-interactive with flags):
```bash
# Create a site with all parameters
wordsail site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --system-name examplecom \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!

# Delete a site without confirmation
wordsail site delete --server production-1 --site examplecom --force

# Provision server with options
wordsail server provision production-1 --force --skip-ssh-check
```

## Project Structure

```
wordsail/
├── cli/           # Go CLI tool (see cli/README.md)
├── ansible/       # Ansible playbooks and roles (see ansible/README.md)
├── docs/          # Additional documentation
├── Makefile       # Build automation
└── version.txt    # Version tracking
```

## Documentation

- **[CLI Documentation](cli/README.md)** - CLI installation, commands, and development
- **[Ansible Documentation](ansible/README.md)** - Playbooks, roles, and direct Ansible usage
- **[CLAUDE.md](CLAUDE.md)** - Development guide for AI assistants

### CLI

- **Language**: Go 1.21+
- **Framework**: Cobra (commands), Survey (interactive prompts)
- **Config**: YAML-based state management
- **Operating Modes**: Interactive (prompts) and Script (flags for automation)

### Infrastructure

- **Target OS**: Ubuntu 24.04
- **Web Server**: Nginx (official repo)
- **PHP**: 8.3 (ondrej/php PPA)
- **Database**: MariaDB
- **Cache**: Redis
- **SSL**: Let's Encrypt (Certbot)
- **Security**: UFW, Fail2ban

## Development

### Build Commands

```bash
# Build CLI
make build

# Run tests
make test

# Format code
make fmt

# Lint code
make lint

# Clean artifacts
make clean
```

### Testing Ansible

```bash
# Validate Ansible playbook syntax
make test-ansible

# Or run directly
cd ansible
ansible-playbook --syntax-check provision.yml
```

## Installation Methods

### From Source (Current)

```bash
git clone https://github.com/your-org/wordsail.git
cd wordsail
make install
wordsail init
```

### Future: Package Managers

```bash
# Coming soon
brew install wordsail
apt install wordsail
```

## Requirements

### CLI Usage

- Go 1.21+ (for building from source)
- Ansible 2.14+
- SSH access to target servers

### Target Servers

- Ubuntu 24.04 LTS
- Fresh server (recommended)
- Root SSH access for provisioning

## How It Works

1. **CLI manages state** - Server and site configuration stored in `~/.wordsail/wordsail.yaml`
2. **CLI executes Ansible** - Generates dynamic inventory and runs playbooks
3. **Ansible configures servers** - Idempotent playbooks ensure consistent state
4. **CLI updates state** - After successful operations, configuration is updated

## CLI Modes

WordSail CLI supports two operating modes to fit different workflows:

### Interactive Mode

Perfect for manual operations and learning the tool. The CLI prompts you for all required information with helpful hints and validation.

```bash
wordsail site create
# Prompts for: server, domain, system name, admin credentials
```

**Benefits:**
- User-friendly with guided prompts
- Built-in validation and helpful hints
- Easy discovery of available options
- Great for manual operations

### Script Mode (Non-Interactive)

Ideal for automation, CI/CD pipelines, and scripting. Provide all parameters as command-line flags.

```bash
wordsail site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --system-name examplecom \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password SecurePass123!
```

**Benefits:**
- Fully scriptable and automatable
- No user interaction required
- Perfect for CI/CD pipelines
- Consistent and reproducible operations

**Available flags:**
- `--non-interactive`: Disable all prompts (required for script mode)
- `--force`: Skip confirmation prompts
- `--skip-ssh-check`: Skip SSH connectivity validation

## Common Workflows

### New Server Setup

```bash
wordsail server add           # Add server details
wordsail server provision     # Install LEMP stack
wordsail site create          # Create first WordPress site
wordsail domain ssl           # Issue SSL certificate
```

### Adding Sites to Existing Server

```bash
wordsail site create          # Interactive site creation
wordsail domain add           # Add www subdomain
wordsail domain ssl           # Issue SSL certificate
```

### Managing Existing Sites

```bash
wordsail site list            # View all sites
wordsail domain add           # Add domain to site
wordsail site delete          # Remove site completely
```

### Automated Deployment (Script Mode)

```bash
# Fully automated site deployment
wordsail site create --non-interactive \
  --server production-1 \
  --domain example.com \
  --system-name examplecom \
  --admin-user admin \
  --admin-email admin@example.com \
  --admin-password "${WP_ADMIN_PASSWORD}"

# Automated site cleanup
wordsail site delete \
  --server production-1 \
  --site examplecom \
  --force
```

## Roadmap

- [ ] SSO Login
- [ ] Site backup/restore/download
- [ ] Multi-PHP version support
- [ ] Site cloning / Staging
- [ ] Resource monitoring
- [ ] Debug Helper
- [ ] Homebrew formula
- [ ] APT repository
- [ ] Auto-updates
- [ ] Shell completions

- **CLI Help**: `wordsail --help` or `wordsail <command> --help`
- **Documentation**: See [cli/README.md](cli/README.md) and [ansible/README.md](ansible/README.md)
- **Issues**: GitHub Issues

---

Built with ❤️ for WordPress hosting automation
