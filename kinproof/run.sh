#!/usr/bin/env bash
set -euo pipefail
go mod tidy
go run ./cmd/kinproof
