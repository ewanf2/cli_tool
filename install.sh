#!/bin/bash

set -e

REPO="ewanf2/cli_tool"
BINARY="ht"
INSTALL_DIR="/usr/local/bin"

# detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# normalise arch (macOS reports arm64, Linux reports aarch64)
case $ARCH in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# normalise OS
case $OS in
  linux) OS="linux" ;;
  darwin) OS="darwin" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

# download from the latest release (GitHub redirects this to the newest non-prerelease)
FILENAME="${BINARY}-${OS}-${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/latest/download/$FILENAME"

echo "Installing $BINARY for $OS-$ARCH..."

curl -fsSL "$URL" -o "/tmp/$FILENAME" || { echo "Download failed: $URL"; exit 1; }
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
