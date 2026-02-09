#!/bin/sh
set -e

# Configuration
OWNER="hunsy9"
REPO="kubesnap"
BINARY="ks"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

log() {
    echo "${GREEN}[kubesnap-installer]${NC} $1"
}

error() {
    echo "${RED}[kubesnap-installer] Error:${NC} $1"
    exit 1
}

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [ "$OS" = "linux" ]; then
    case "$ARCH" in
        x86_64) ARCH="x86_64" ;;
        aarch64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac
elif [ "$OS" = "darwin" ]; then
    case "$ARCH" in
        x86_64) ARCH="x86_64" ;;
        arm64) ARCH="arm64" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac
else
    error "Unsupported OS: $OS"
fi

log "Detected $OS/$ARCH"

# Get Latest Version
log "Checking latest version..."
LATEST_URL="https://api.github.com/repos/$OWNER/$REPO/releases/latest"

if command -v curl >/dev/null; then
    VERSION=$(curl -s $LATEST_URL | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
elif command -v wget >/dev/null; then
    VERSION=$(wget -qO- $LATEST_URL | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
else
    error "Neither curl nor wget found. Please install one."
fi

if [ -z "$VERSION" ]; then
    error "Failed to fetch latest version."
fi

log "Latest version is $VERSION"

# Construct Download URL
# Pattern from goreleaser: kubesnap_Linux_x86_64.tar.gz
# Note: OS name in artifact is capitalized (Linux, Darwin)
ARTIFACT_OS=$(echo "$OS" | awk '{print toupper(substr($0,1,1))substr($0,2)}')
ARTIFACT_NAME="${REPO}_${ARTIFACT_OS}_${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$OWNER/$REPO/releases/download/$VERSION/$ARTIFACT_NAME"

# Download
TMP_DIR=$(mktemp -d)
clean_up() {
    rm -rf "$TMP_DIR"
}
trap clean_up EXIT

log "Downloading $DOWNLOAD_URL ..."
if command -v curl >/dev/null; then
    curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/$ARTIFACT_NAME"
else
    wget -qO "$TMP_DIR/$ARTIFACT_NAME" "$DOWNLOAD_URL"
fi

# Extract
log "Extracting..."
tar -xzf "$TMP_DIR/$ARTIFACT_NAME" -C "$TMP_DIR"

# Install
log "Installing to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
else
    log "Sudo permission required to move binary to $INSTALL_DIR"
    sudo mv "$TMP_DIR/$BINARY" "$INSTALL_DIR/$BINARY"
fi

# Verify
if command -v $BINARY >/dev/null; then
    log "Installation successful!"
    log "Run '$BINARY --version' to get started."
else
    error "Installation failed."
fi
