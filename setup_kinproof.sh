#!/usr/bin/env bash
set -euo pipefail

echo "=== KinProof-Go Sovereign x402 Identity Layer Bootstrap ==="

# Clean previous state if exists
rm -rf kinproof-go

mkdir -p kinproof-go/{cmd/kinproof,internal/{hd,identity,intent,proof,nonce,attribution,vault,sei},proofs,receipts,keys,scripts,docs}

cd kinproof-go

# go.mod
cat <<'GO_MOD' > go.mod
module kinproof-go

go 1.23

require (
	github.com/btcsuite/btcd/btcutil v1.1.5
	github.com/btcsuite/btcd/chaincfg v1.3.0
	github.com/ethereum/go-ethereum v1.14.12
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/tyler-smith/go-bip39 v1.1.0
)
GO_MOD

# README.md
cat <<'README' > README.md
# KinProof-Go — Sovereign Intent Layer for x402 + Sei KinSigil

Proof-bound rotating identity with persistent sovereign root, replay-resistant execution envelopes, and attribution lineage.

This implementation provides the missing high-value identity and intent continuity layer above x402 payment protocols.

Author: The Keeper (Angel)
Repository: https://github.com/Pray4Love1/kinproof-go
