#!/bin/bash

set -e

GO_VERSION="1.24.3"
OS="linux"
ARCH="amd64"

echo "[+] Downloading Go ${GO_VERSION}..."
wget https://go.dev/dl/go${GO_VERSION}.${OS}-${ARCH}.tar.gz

echo "[+] Removing any existing Go installation..."
sudo rm -rf /usr/local/go

echo "[+] Extracting Go to /usr/local..."
sudo tar -C /usr/local -xzf go${GO_VERSION}.${OS}-${ARCH}.tar.gz

echo "[+] Cleaning up archive..."
rm go${GO_VERSION}.${OS}-${ARCH}.tar.gz

echo "[+] Setting up PATH..."

# Tambahkan ke ~/.bashrc atau ~/.zshrc tergantung shell
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
source ~/.bashrc

echo "[✓] Go ${GO_VERSION} installed successfully."
go version