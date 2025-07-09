#!/bin/bash

LINK_NAME="document-tracking-api"

# Remove old link if it exists
[ -L "$LINK_NAME" ] && rm "$LINK_NAME"

# Detect OS and set target
if [[ "$1" == "--linux-arm64" ]]; then
    TARGET="builds/doc-tracker-linux-arm64"
elif [[ "$1" == "--linux-amd64" ]]; then
    TARGET="builds/doc-tracker-linux-amd64"
elif [[ "$1" == "--mac" ]]; then
    TARGET="builds/doc-tracker-darwin-amd64"
elif [[ "$1" == "--windows" ]]; then
    TARGET="builds/doc-tracker-windows-amd64"
elif [[ "$(uname)" == "Linux" ]]; then
    TARGET="builds/doc-tracker-linux-amd64"
else
    TARGET="builds/doc-tracker-darwin-amd64"
fi

# Create new symlink
ln -s "$TARGET" "$LINK_NAME"