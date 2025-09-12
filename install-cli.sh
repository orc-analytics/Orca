#!/bin/bash

set -e

REPO="orc-analytics/orca"
INSTALL_NAME="orca"

# Disallow root user
if [ "$EUID" -eq 0 ]; then
  echo "Do not run this script as root. Please run as a regular user."
  exit 1
fi

# Detect OS type
detect_os() {
  UNAME="$(uname -s)"
  ARCH="$(uname -m)"
  case "$UNAME" in
    Darwin)
      if [ "$ARCH" = "x86_64" ]; then
        OS="mac-intel"
      elif [ "$ARCH" = "arm64" ]; then
        OS="mac-arm"
      else
        echo "Unsupported Mac architecture: $ARCH"
        exit 1
      fi
      ;;
    Linux)
      OS="linux"
      ;;
    MINGW*|MSYS*|CYGWIN*|Windows_NT)
      OS="windows"
      ;;
    *)
      echo "Unsupported OS: $UNAME"
      exit 1
      ;;
  esac
}

# Get latest release version from GitHub API
get_latest_version() {
  echo "Fetching latest Orca CLI version..."
  LATEST_VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep -oP '"tag_name":\s*"\K[^"]+')
  if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to retrieve latest version"
    exit 1
  fi
  echo "Latest version: $LATEST_VERSION"
}

# Download the appropriate binary
download_binary() {
  BINARY_NAME="orca-cli-${OS}"
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}"
  TMP_FILE="$(mktemp)"
  echo "Downloading $DOWNLOAD_URL"
  curl -L "$DOWNLOAD_URL" -o "$TMP_FILE"
  chmod +x "$TMP_FILE"
}

# Determine writable install directories
find_install_dirs() {
  # Create directories if they don't exist
  mkdir -p "$HOME/.local/bin" "$HOME/.local/share"
  
  SHARE_CANDIDATES=("$HOME/.local/share" "$HOME/share" "/usr/local/share")
  BIN_CANDIDATES=("$HOME/.local/bin" "$HOME/bin" "/usr/local/bin")

  for dir in "${SHARE_CANDIDATES[@]}"; do
    if [ -d "$dir" ] && [ -w "$dir" ]; then
      SHARE_DIR="$dir/orc_a"
      mkdir -p "$SHARE_DIR"
      break
    fi
  done

  for dir in "${BIN_CANDIDATES[@]}"; do
    if [ -d "$dir" ] && [ -w "$dir" ]; then
      BIN_DIR="$dir"
      break
    fi
  done

  if [ -z "$SHARE_DIR" ] || [ -z "$BIN_DIR" ]; then
    echo "No writable share/bin directory found. Please add one or run with elevated permissions."
    exit 1
  fi
}

# Install binary and manage symlink
install_binary() {
  FINAL_BINARY="$SHARE_DIR/$INSTALL_NAME"
  SYMLINK_PATH="$BIN_DIR/$INSTALL_NAME"

  rm -f "$SYMLINK_PATH"

  mv "$TMP_FILE" "$FINAL_BINARY"
  chmod +x "$FINAL_BINARY"
  ln -sf "$FINAL_BINARY" "$SYMLINK_PATH"

  echo ""
  echo "✅ Orca CLI installed to: $FINAL_BINARY"
  echo "✅ Symlink created at: $SYMLINK_PATH"
  echo "🔗 To get started, visit: https://github.com/orc-analytics/Orca#readme"
}

# Run install steps
detect_os
get_latest_version
download_binary
find_install_dirs
install_binary

