.PHONY: all build install install-user test test-cli test-ansible clean fmt lint help

# Version from version.txt
VERSION=$(shell cat version.txt)

# Default target
all: build

# Build CLI
build:
	@echo "Building WordSail CLI (v$(VERSION))..."
	@cd cli && make build
	@echo "✓ Build complete: cli/wordsail"

# Install CLI to /usr/local/bin (requires sudo)
install: build
	@echo "Installing WordSail CLI to /usr/local/bin..."
	@cd cli && make install
	@echo ""
	@echo "✓ Installation complete!"
	@echo ""
	@echo "Next step: Run 'wordsail init' to set up your environment"

# Install CLI to ~/bin (no sudo required)
install-user: build
	@echo "Installing WordSail CLI to ~/bin..."
	@cd cli && make install-user
	@echo ""
	@echo "✓ Installation complete!"
	@echo ""
	@echo "Next step: Run 'wordsail init' to set up your environment"

# Run all tests
test: test-cli test-ansible
	@echo "✓ All tests passed"

# Run CLI tests
test-cli:
	@echo "Running CLI tests..."
	@cd cli && make test

# Test Ansible playbooks (syntax check)
test-ansible:
	@echo "Validating Ansible playbooks..."
	@cd ansible && ansible-playbook --syntax-check provision.yml
	@cd ansible && ansible-playbook --syntax-check website.yml
	@cd ansible && ansible-playbook --syntax-check playbooks/domain_management.yml
	@cd ansible && ansible-playbook --syntax-check playbooks/delete_site.yml
	@echo "✓ Ansible syntax validation passed"

# Format Go code
fmt:
	@echo "Formatting Go code..."
	@cd cli && make fmt

# Lint Go code
lint:
	@echo "Linting Go code..."
	@cd cli && make lint

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@cd cli && make clean
	@rm -f wordsail
	@echo "✓ Clean complete"

# Development: quick run
run: build
	@./cli/wordsail

# Show help
help:
	@echo "WordSail Build System"
	@echo "Version: $(VERSION)"
	@echo ""
	@echo "Available targets:"
	@echo "  build        - Build the CLI binary"
	@echo "  install      - Install CLI to /usr/local/bin (requires sudo)"
	@echo "  install-user - Install CLI to ~/bin (no sudo required)"
	@echo "  test         - Run all tests (CLI + Ansible)"
	@echo "  test-cli     - Run CLI tests only"
	@echo "  test-ansible - Validate Ansible playbook syntax"
	@echo "  fmt          - Format Go code"
	@echo "  lint         - Lint Go code"
	@echo "  clean        - Remove build artifacts"
	@echo "  run          - Build and run CLI"
	@echo "  help         - Show this help message"
	@echo ""
	@echo "Quick start:"
	@echo "  make build && ./cli/wordsail init"
