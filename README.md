# WordSail Ansible

Automated WordPress hosting infrastructure using Ansible.

## TL;DR

```bash
# Provision a fresh Ubuntu server
ansible-playbook provision.yml -i "SERVER_IP," -u root

# Create a WordPress site
ansible-playbook website.yml -i "SERVER_IP," -u wordsail \
  --extra-vars "domain=example.com system_name=examplecom wp_admin_user=admin wp_admin_email=admin@example.com wp_admin_password=SecurePass123"
```

## What It Does

| Playbook                          | Purpose                                               |
| --------------------------------- | ----------------------------------------------------- |
| `provision.yml`                   | Full server setup (Nginx, PHP 8.3, MariaDB, security) |
| `website.yml`                     | Create isolated WordPress sites with SSL              |
| `playbooks/domain_management.yml` | Add/remove domains, issue SSL certs                   |

## Roles Overview

| Role           | What It Does                                                                                  |
| -------------- | --------------------------------------------------------------------------------------------- |
| **bootstrap**  | Creates `wordsail` user, installs base packages, sets up certbot, configures fail2ban & redis |
| **database**   | Installs MariaDB, creates admin user, secures installation                                    |
| **nginx**      | Installs Nginx from official repo, configures global settings, generates default SSL          |
| **php**        | Installs PHP 8.3 + extensions, Composer, WP-CLI, security hardening                           |
| **security**   | UFW firewall (ports 22/80/443), SSH hardening                                                 |
| **website**    | Creates site user, database, PHP-FPM pool, Nginx vhost, installs WordPress                    |
| **libs**       | Reusable tasks: add_domain, remove_domain, issue_ssl                                          |
| **operations** | Server ops: delete_site                                                                       |

xw

## Domain Management

```bash
# Add domain to Nginx
ansible-playbook playbooks/domain_management.yml -i "IP," -u wordsail \
  --extra-vars "operation=add_domain domain=newdomain.com system_name=sitename"

# Remove domain
ansible-playbook playbooks/domain_management.yml -i "IP," -u wordsail \
  --extra-vars "operation=remove_domain domain=olddomain.com"

# Issue SSL certificate
ansible-playbook playbooks/domain_management.yml -i "IP," -u wordsail \
  --extra-vars "operation=issue_ssl domain=example.com certbot_email=admin@example.com"
```

## Directory Structure

```
ansible/
├── provision.yml          # Server provisioning playbook
├── website.yml            # WordPress site creation playbook
├── playbooks/             # Additional playbooks
│   └── domain_management.yml
├── roles/
│   ├── bootstrap/         # Base server setup
│   ├── database/          # MariaDB installation
│   ├── nginx/             # Web server config
│   ├── php/               # PHP 8.3 + extensions
│   ├── security/          # UFW + SSH hardening
│   ├── website/           # WordPress deployment
│   ├── libs/              # Reusable task libraries
│   └── operations/        # Server operation tasks
├── group_vars/
│   └── all.yml            # Global variables
├── inventory/             # Inventory files
└── requirements.yml       # Ansible Galaxy dependencies
```

## Server Layout

After provisioning, servers have:

```
/sites/
└── example.com/
    ├── public/            # WordPress files
    ├── logs/              # Site logs

/etc/nginx/sites-available/example.com/
├── example.com            # Main config
├── server/                # Server-level includes
├── location/              # Location blocks
├── before/                # Pre-processing rules
└── after/                 # Post-processing (redirects)
```

## Requirements

```bash
# Install Ansible
pip install ansible

# Install required collections
ansible-galaxy install -r requirements.yml
```

## Required Variables

Set in `group_vars/all.yml` or pass via `--extra-vars`:

| Variable                     | Description                      |
| ---------------------------- | -------------------------------- |
| `wordsail_ssh_key`           | SSH public key for wordsail user |
| `mysql_wordsailbot_password` | MySQL admin password             |
| `certbot_email`              | Email for Let's Encrypt          |

## Stack

- **OS**: Ubuntu 24.04
- **Web**: Nginx (official repo)
- **PHP**: 8.3 (ondrej/php PPA)
- **DB**: MariaDB
- **Cache**: Redis
- **SSL**: Let's Encrypt (Certbot)
- **Security**: UFW, Fail2ban

## Upcoming

- [ ] Multi-PHP version management
- [ ] Backup/restore automation
