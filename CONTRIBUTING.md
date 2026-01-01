# Contributing to WordSail Ansible

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## How to Contribute

### Reporting Issues

- Check existing issues before creating a new one
- Use a clear, descriptive title
- Include steps to reproduce the problem
- Specify your environment (OS, Ansible version, target server OS)

### Submitting Changes

1. **Fork the repository** and create a new branch from `main`
2. **Make your changes** following the guidelines below
3. **Test your changes** on a fresh Ubuntu server
4. **Submit a pull request** with a clear description

### Development Guidelines

#### Ansible Best Practices

- Use fully qualified collection names (e.g., `ansible.builtin.file`)
- Make tasks idempotent - running twice should produce the same result
- Use meaningful task names that describe what the task does
- Add `tags` to tasks for selective execution
- Use handlers for service restarts

#### Code Style

- Use 2-space indentation in YAML files
- Keep lines under 120 characters when possible
- Use lowercase for variable names with underscores (e.g., `site_user`)
- Quote strings that contain special characters or start with `{`

#### Variables

- Define defaults in `defaults/main.yml` for roles
- Document required variables in role README or comments
- Use descriptive variable names
- Avoid hardcoding values - use variables instead

#### Security

- Never commit real credentials, SSH keys, or secrets
- Use `ansible-vault` for sensitive data in examples
- Follow least-privilege principle for file permissions
- Test security configurations (UFW, SSH hardening)

### Testing

Before submitting a PR, please test:

1. **Syntax check**: `ansible-playbook --syntax-check provision.yml`
2. **Lint check**: `ansible-lint provision.yml` (if installed)
3. **Fresh server test**: Run against a fresh Ubuntu 22.04 server

### Commit Messages

- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Fix bug" not "Fixes bug")
- Keep the first line under 72 characters
- Reference issues when applicable (e.g., "Fix #123")

Example:
```
Add support for Ubuntu 24.04

- Update PHP PPA configuration for noble
- Adjust package names for compatibility
- Update documentation with new OS support

Fixes #42
```

### Pull Request Process

1. Update documentation if adding new features
2. Update `CHANGELOG.md` with your changes
3. Ensure all tests pass
4. Request review from maintainers

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow

## Questions?

Open an issue with the "question" label if you need help.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
