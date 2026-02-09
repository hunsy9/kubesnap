#!/bin/sh
set -e

# Configuration
BINARY="ks"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.kubesnap"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m'

log() {
    echo "${GREEN}[kubesnap-uninstall]${NC} $1"
}

warn() {
    echo "${YELLOW}[kubesnap-uninstall] Warning:${NC} $1"
}

error() {
    echo "${RED}[kubesnap-uninstall] Error:${NC} $1"
    exit 1
}

# 1. Remove Binary
if [ -f "$INSTALL_DIR/$BINARY" ]; then
    log "Removing binary from $INSTALL_DIR/$BINARY..."
    if [ -w "$INSTALL_DIR" ]; then
        rm "$INSTALL_DIR/$BINARY"
    else
        log "Sudo permission required to remove binary"
        sudo rm "$INSTALL_DIR/$BINARY"
    fi
else
    warn "Binary not found in $INSTALL_DIR/$BINARY"
fi

# 2. Ask to remove configuration
printf "${YELLOW}[kubesnap-uninstall]${NC} Do you want to remove configuration and cache at $CONFIG_DIR? (y/N): "
read -r CONFIRM
if [ "$CONFIRM" = "y" ] || [ "$CONFIRM" = "Y" ]; then
    if [ -d "$CONFIG_DIR" ]; then
        log "Removing configuration directory $CONFIG_DIR..."
        rm -rf "$CONFIG_DIR"
    else
        warn "Configuration directory not found."
    fi
else
    log "Skipping configuration removal."
fi

log "Uninstallation complete!"
