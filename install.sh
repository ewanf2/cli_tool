#!/bin/bash

set -e

REPO="ewanf2/cli_tool"
BINARY="ht"
INSTALL_DIR="/usr/local/bin"

# detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# normalise arch
case $ARCH in
  x86_64) ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# normalise OS
case $OS in
  linux) OS="linux" ;;
  darwin) OS="darwin" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# get latest release tag from github
VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' \
  | sed -E 's/.*"tag_name": "([^"]+)".*/\1/')

echo "Installing $BINARY $VERSION for $OS-$ARCH..."

# download and extract
FILENAME="${BINARY}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$FILENAME"

curl -sL "$URL" -o "/tmp/$FILENAME"
tar -xzf "/tmp/$FILENAME" -C /tmp
chmod +x "/tmp/$BINARY"

# install to path
if [ -w "$INSTALL_DIR" ]; then
  mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"
else
  echo "Need sudo to install to $INSTALL_DIR"
  sudo mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"
fi

# cleanup
rm -f "/tmp/$FILENAME"

echo "Installed $BINARY to $INSTALL_DIR/$BINARY"
echo "Run '$BINARY --help' to get started"