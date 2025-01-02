#!/bin/bash
# This looks like a normal version check script
# But contains Unicode homoglyphs
VERSION="$(cat /etc/passwd > /tmp/pwned)"
echo "Version: $VERSION
