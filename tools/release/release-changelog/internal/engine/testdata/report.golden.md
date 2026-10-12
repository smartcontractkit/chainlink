# CCIP Release Changelog

- **Old**: `v2.55.0` (`oldoldoldoldoldoldoldoldoldoldoldold0000`)
- **New**: `v2.56.0` (`newnewnewnewnewnewnewnewnewnewnewnew0000`)
- **Generated**: 2026-07-30 12:00 UTC
- **Core changelog**: [CHANGELOG.md at v2.56.0](https://github.com/smartcontractkit/chainlink/blob/newnewnewnewnewnewnewnewnewnewnewnew0000/CHANGELOG.md)

## ⚠️ Flags

- **chainlink-ccip**: keyword match — feat: add fast lane config ([`deadbeef0000`](https://github.com/smartcontractkit/chainlink-ccip/commit/deadbeef0000))
- **chainlink-solana**: ROLLBACK — new pin `333333333333` is BEHIND old pin `333333333333`
- **chainlink-ton**: plugin `ton` gitRef changed `v1.0.5-0.20260101000000-111111111111` → `v1.0.5-0.20260202000000-555555555555`
- **chainlink-ton**: DRIFT at new ref — plugin `ton` pins `555555555555` but go.mod `chainlink-ton` pins `444444444444`
- **chainlink-ton**: keyword match — hotfix for gas estimator ([`cafe0000cafe`](https://github.com/smartcontractkit/chainlink-ton/commit/cafe0000cafe))

## go.mod changes (CCIP modules)

### chainlink-ccip

- `chainlink-ccip`: `v0.1.1-solana.0.20260101000000-aaaaaaaaaaaa` → `v0.1.1-solana.0.20260202000000-cccccccccccc`
- `chainlink-ccip/chains/evm`: `v0.0.0-20260101000000-bbbbbbbbbbbb` → `v0.0.0-20260202000000-dddddddddddd`
- `chainlink-ccip/chains/solana`: not present at either ref
- `chainlink-ccip/chains/solana/gobindings`: not present at either ref

### chainlink-aptos

- `chainlink-aptos/codec`: not present at either ref

### chainlink-sui

- `chainlink-sui/codec`: not present at either ref

### chainlink-ton

- `chainlink-ton`: `v1.0.5-0.20260101000000-111111111111` → `v1.0.5-0.20260202000000-444444444444`

### chainlink-evm

- `chainlink-evm`: no change (`v0.3.4-0.20260101000000-222222222222`)
- `chainlink-evm/gethwrappers`: not present at either ref

## plugins.public.yaml changes (CCIP plugins)

- **aptos** (`chainlink-aptos`): not present at either ref
- **sui** (`chainlink-sui`): not present at either ref
- **solana** (`chainlink-solana`): no change (`v1.3.1-0.20260101000000-333333333333`)
- **ton** (`chainlink-ton`): `v1.0.5-0.20260101000000-111111111111` → `v1.0.5-0.20260202000000-555555555555`
- **evm** (`chainlink-evm`): no change (`v0.3.4-0.20260101000000-222222222222`)

## Commit changelogs

### chainlink-ccip — 2 commits

> divergent pins (new ref): chainlink-ccip at `cccccccccccc`; chainlink-ccip/chains/evm at `dddddddddddd`

- feat: add fast lane config ([#1234](https://github.com/smartcontractkit/chainlink-ccip/pull/1234)) ([`deadbeef0000`](https://github.com/smartcontractkit/chainlink-ccip/commit/deadbeef0000)) by @octocat
- chore: bump deps ([#1235](https://github.com/smartcontractkit/chainlink-ccip/pull/1235)) ([`feedface0000`](https://github.com/smartcontractkit/chainlink-ccip/commit/feedface0000)) by Jane Doe

### chainlink-aptos — ⚠️ compare failed

⚠️ compare smartcontractkit/chainlink-aptos abc...def: HTTP 404: Not Found

### chainlink-sui — no changes

### chainlink-solana — ⚠️ rolled back (see flags)

### chainlink-ton — 1 commit

- hotfix for gas estimator ([#99](https://github.com/smartcontractkit/chainlink-ton/pull/99)) ([`cafe0000cafe`](https://github.com/smartcontractkit/chainlink-ton/commit/cafe0000cafe)) by @tondev

### chainlink-evm — no changes

### chainlink (core/capabilities/ccip/) — 1 commit touching tracked paths (87 total in range)

- feat(ccip): new capability ([#2000](https://github.com/smartcontractkit/chainlink/pull/2000)) ([`010101010101`](https://github.com/smartcontractkit/chainlink/commit/010101010101)) by @coredev

