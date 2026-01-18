# WordSail

Automated WordPress hosting on Ubuntu servers. One command to provision, one command to deploy.

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/shariffff/wordsail/main/install.sh | bash
```

Then initialize:
```bash
wordsail init
```

## Quick Start

```bash
# 1. Add your server
wordsail server add

# 2. Provision it (installs Nginx, PHP, MariaDB, Redis, SSL)
wordsail server provision myserver

# 3. Create a WordPress site
wordsail site create

# 4. Issue SSL certificate
wordsail domain ssl
```

## What It Does

**Server provisioning:**
- Nginx from official repo
- PHP 8.3 with optimized FPM pools
- MariaDB with secure defaults
- Redis for object caching
- Let's Encrypt SSL via Certbot
- UFW firewall + Fail2ban

**Site isolation:**
- Each site runs as its own Linux user
- Dedicated PHP-FPM pool per site
- Isolated file permissions

## Commands

```bash
wordsail server add          # Add a server
wordsail server provision    # Provision server with LEMP stack
wordsail server list         # List servers

wordsail site create         # Create WordPress site
wordsail site list           # List sites
wordsail site delete         # Delete site

wordsail domain add          # Add domain to site
wordsail domain ssl          # Issue SSL certificate

wordsail config show         # Show configuration
```

All commands support `--help` for details.

## Requirements

- Ansible 2.14+ on your local machine
- Ubuntu 24.04 target server with root SSH access

## Documentation

- [CLI Reference](cli/README.md)
- [Ansible Playbooks](ansible/README.md)

## Development

```bash
git clone https://github.com/shariffff/wordsail.git
cd wordsail
make build    # Build CLI
make test     # Run tests
```

## License

MIT
