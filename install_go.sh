#!/bin/bash

set -e

GO_VERSION="1.24.3"

# Detect OS
UNAME=$(uname -s)
case "$UNAME" in
    Linux*)     OS="linux";;
    Darwin*)    OS="darwin";;
    *)          echo "Unsupported OS: $UNAME"; exit 1;;
esac

# Detect ARCH
ARCH_UNAME=$(uname -m)
case "$ARCH_UNAME" in
    x86_64)    ARCH="amd64";;
    arm64)     ARCH="arm64";;
    aarch64)   ARCH="arm64";;
    *)         echo "Unsupported architecture: $ARCH_UNAME"; exit 1;;
esac

echo "[+] Downloading Go ${GO_VERSION} for ${OS}-${ARCH}..."
wget https://go.dev/dl/go${GO_VERSION}.${OS}-${ARCH}.tar.gz

echo "[+] Removing any existing Go installation..."
sudo rm -rf /usr/local/go

echo "[+] Extracting Go to /usr/local..."
sudo tar -C /usr/local -xzf go${GO_VERSION}.${OS}-${ARCH}.tar.gz

echo "[+] Cleaning up archive..."
rm go${GO_VERSION}.${OS}-${ARCH}.tar.gz

echo "[+] Setting up PATH..."
SHELL_RC="$HOME/.bashrc"
if [[ $SHELL == *"zsh"* ]]; then
    SHELL_RC="$HOME/.zshrc"
fi

# Hapus PATH lama Go dulu
sed -i.bak '/\/usr\/local\/go\/bin/d' $SHELL_RC

echo "export PATH=\$PATH:/usr/local/go/bin" >> $SHELL_RC
source $SHELL_RC

echo "[✓] Go ${GO_VERSION} installed successfully."
go version
