#!/bin/bash

set -e

REPO="Predixus/Orca"
INSTALL_NAME="orca"

# Disallow root user
if [ "$EUID" -eq 0 ]; then
  echo "Do not run this script as root. Please run as a regular user."
  exit 1
fi

# Determine OS type
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

# Find latest release version
get_latest_version() {
  LATEST_VERSION=$(curl -s https://api.github.com/repos/${REPO}/releases/latest | grep -oP '"tag_name":\s*"\K[^"]+')
  if [ -z "$LATEST_VERSION" ]; then
    echo "Failed to retrieve latest version"
    exit 1
  fi
}

# Download binary
download_binary() {
  BINARY_NAME="orca-cli-${OS}"
  DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}"
  TMP_FILE="$(mktemp)"
  echo "Downloading $DOWNLOAD_URL"
  curl -L "$DOWNLOAD_URL" -o "$TMP_FILE"
  chmod +x "$TMP_FILE"
}

# Find share directory for actual binary and bin dir for symlink
find_install_dirs() {
  # Preferred share and bin directories
  SHARE_CANDIDATES=("$HOME/.local/share" "$HOME/share" "/usr/local/share")
  BIN_CANDIDATES=("$HOME/.local/bin" "$HOME/bin" "/usr/local/bin")

  for dir in "${SHARE_CANDIDATES[@]}"; do
    if [ -d "$dir" ] && [ -w "$dir" ]; then
      SHARE_DIR="$dir/predixus"
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

# Install binary safely and create symlink
install_binary() {
  FINAL_BINARY="$SHARE_DIR/$INSTALL_NAME"
  SYMLINK_PATH="$BIN_DIR/$INSTALL_NAME"

  if [ -e "$SYMLINK_PATH" ]; then
    read -p "A binary named '$INSTALL_NAME' already exists at $SYMLINK_PATH. Replace it? (y/N): " choice
    case "$choice" in
      y|Y ) echo "Replacing existing symlink...";;
      * ) echo "Installation aborted."; exit 1;;
    esac
    rm -f "$SYMLINK_PATH"
  fi

  mv "$TMP_FILE" "$FINAL_BINARY"
  chmod +x "$FINAL_BINARY"
  ln -s "$FINAL_BINARY" "$SYMLINK_PATH"
  echo "Binary installed to $FINAL_BINARY"
  echo "Symlink created at $SYMLINK_PATH"
  echo "To get started, visit the documentation at: https://github.com/Predixus/Orca#readme"
}

# Main
detect_os
get_latest_version
download_binary
find_install_dirs
install_binary

