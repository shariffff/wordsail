#!/bin/bash
set -e

# WordSail Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/shariffff/wordsail/main/install.sh | bash

REPO="shariffff/wordsail"
BINARY_NAME="wordsail"

# Install directory (like Bun's ~/.bun)
install_dir="${WORDSAIL_INSTALL:-$HOME/.wordsail}"
bin_dir="$install_dir/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
DIM='\033[2m'
NC='\033[0m'

info() { echo -e "${GREEN}$1${NC}"; }
warn() { echo -e "${YELLOW}$1${NC}"; }
error() { echo -e "${RED}error${NC}: $1" >&2; exit 1; }

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Linux*)  echo "linux" ;;
        Darwin*) echo "darwin" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) error "Unsupported operating system: $(uname -s)" ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) error "Unsupported architecture: $(uname -m)" ;;
    esac
}

# Get latest release version from GitHub
get_latest_version() {
    curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" 2>/dev/null | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Add to shell config
setup_shell() {
    local shell_name=$(basename "$SHELL")

    case $shell_name in
        fish)
            local config_file="$HOME/.config/fish/config.fish"
            local shell_export="set --export WORDSAIL_INSTALL \"$install_dir\"\nset --export PATH \$WORDSAIL_INSTALL/bin \$PATH"
            ;;
        zsh)
            local config_file="$HOME/.zshrc"
            local shell_export="export WORDSAIL_INSTALL=\"$install_dir\"\nexport PATH=\"\$WORDSAIL_INSTALL/bin:\$PATH\""
            ;;
        bash)
            # Prefer .bashrc on Linux, .bash_profile on macOS
            if [[ -f "$HOME/.bashrc" ]]; then
                local config_file="$HOME/.bashrc"
            else
                local config_file="$HOME/.bash_profile"
            fi
            local shell_export="export WORDSAIL_INSTALL=\"$install_dir\"\nexport PATH=\"\$WORDSAIL_INSTALL/bin:\$PATH\""
            ;;
        *)
            echo ""
            warn "Could not detect shell. Manually add to your shell config:"
            echo ""
            echo "  export WORDSAIL_INSTALL=\"$install_dir\""
            echo "  export PATH=\"\$WORDSAIL_INSTALL/bin:\$PATH\""
            echo ""
            return
            ;;
    esac

    # Check if already configured
    if [[ -f "$config_file" ]] && grep -q "WORDSAIL_INSTALL" "$config_file" 2>/dev/null; then
        return
    fi

    # Create config file if it doesn't exist
    if [[ ! -f "$config_file" ]]; then
        touch "$config_file"
    fi

    # Check if writable
    if [[ ! -w "$config_file" ]]; then
        warn "Could not write to $config_file. Manually add:"
        echo ""
        echo -e "  $shell_export"
        echo ""
        return
    fi

    # Append to config
    echo "" >> "$config_file"
    echo "# WordSail" >> "$config_file"
    echo -e "$shell_export" >> "$config_file"

    echo -e "${DIM}Added to $config_file${NC}"
}

# Download and install
install_wordsail() {
    local os=$(detect_os)
    local arch=$(detect_arch)
    local version="${WORDSAIL_VERSION:-$(get_latest_version)}"

    if [[ -z "$version" ]]; then
        error "Could not determine latest version. Set WORDSAIL_VERSION manually."
    fi

    # Remove 'v' prefix if present for filename
    local version_num="${version#v}"

    echo -e "${DIM}Installing WordSail ${version} (${os}/${arch})${NC}"

    # Construct download URL
    local filename="${BINARY_NAME}_${version_num}_${os}_${arch}"
    if [[ "$os" = "windows" ]]; then
        filename="${filename}.zip"
    else
        filename="${filename}.tar.gz"
    fi

    local url="https://github.com/${REPO}/releases/download/${version}/${filename}"

    # Create temp directory
    local tmp_dir=$(mktemp -d)
    trap "rm -rf ${tmp_dir}" EXIT

    # Download
    if ! curl -fsSL "$url" -o "${tmp_dir}/${filename}" 2>/dev/null; then
        error "Failed to download from ${url}"
    fi

    # Extract
    cd "${tmp_dir}"
    if [[ "$os" = "windows" ]]; then
        unzip -q "${filename}"
    else
        tar -xzf "${filename}"
    fi

    # Create install directory
    mkdir -p "${bin_dir}"

    # Install binary
    local binary="${BINARY_NAME}"
    if [[ "$os" = "windows" ]]; then
        binary="${BINARY_NAME}.exe"
    fi

    if [[ -f "${binary}" ]]; then
        mv "${binary}" "${bin_dir}/"
    elif [[ -f "${BINARY_NAME}/${binary}" ]]; then
        mv "${BINARY_NAME}/${binary}" "${bin_dir}/"
    else
        error "Binary not found in archive"
    fi

    chmod +x "${bin_dir}/${binary}"
}

# Main
main() {
    echo ""
    echo -e "${BOLD}WordSail${NC} Installer"
    echo ""

    # Check for required tools
    command -v curl >/dev/null 2>&1 || error "curl is required but not installed"
    command -v tar >/dev/null 2>&1 || error "tar is required but not installed"

    install_wordsail
    setup_shell

    echo ""
    echo -e "${GREEN}WordSail was installed successfully!${NC}"
    echo ""
    echo "Run the following to get started:"
    echo ""
    echo -e "  ${BOLD}wordsail init${NC}"
    echo ""
    echo -e "${DIM}If 'wordsail' is not found, restart your terminal.${NC}"
    echo ""
}

main "$@"
