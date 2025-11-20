#!/usr/bin/env bash

# Free up disk space on GitHub ubuntu-latest runners before running CCIP tests
# This script removes unnecessary pre-installed software to prevent "no space left on device" errors

set -e

echo "🧹 Freeing up disk space for CCIP tests..."
df -h

# Remove large pre-installed software (~7GB)
echo "Removing large pre-installed software..."
sudo rm -rf /usr/share/dotnet /usr/local/lib/android /opt/ghc /usr/local/.ghcup /usr/share/swift /usr/local/lib/node_modules || true

# Clean package cache (~1GB)
echo "Cleaning package cache..."
sudo apt-get autoremove -y && sudo apt-get autoclean -y && sudo apt-get clean || true

# Remove unnecessary packages (~2GB)
echo "Removing unnecessary packages..."
sudo apt-get remove -y '^aspnetcore-.*' '^dotnet-.*' azure-cli google-cloud-sdk hhvm google-chrome-stable firefox powershell mono-devel || true

# Docker cleanup (~1GB)
echo "Cleaning Docker..."
docker system prune -af --volumes || true

echo "✅ Cleanup completed. Available space:"
df -h
echo ""
