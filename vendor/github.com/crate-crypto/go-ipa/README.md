# go-ipa

> go-ipa is a library of cryptographic primitives for Verkle Trees.

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/crate-crypto/go-ipa/blob/main/LICENSE)
[![Go Version](https://img.shields.io/badge/go-v1.18-green.svg)](https://golang.org/dl/)

## Table of Contents

- [go-ipa](#go-ipa)
  - [Table of Contents](#table-of-contents)
  - [Description](#description)
  - [Usage in Verkle Tree client libraries](#usage-in-verkle-tree-client-libraries)
  - [Test \& Benchmarks](#test--benchmarks)
  - [Security](#security)
  - [LICENSE](#license)

## Description

go-ipa implements the [Verkle Tree cryptography spec](https://github.com/crate-crypto/verkle-trie-ref) with extra optimizations.

The includes:
- Implementation of the Bandersnatch curve, and Banderwagon prime-order group.
- Pedersen Commitment for vector commitments using precomputed tables.
- Inner Product Argument prover and verifier implementations for polynomials in evaluation form.
- Multiproof prover and verifier implementations.

## Usage in Verkle Tree client libraries

It's extremely important that clients using this library for Verkle Tree implementations only use the following packages:
- `common` for general utility functions.
- `banderwagon` for the prime-order group.
- `ipa` for proof generation and verification.

**Do not** use the `bandersnatch` package directly nor use unsafe functions to get into `banderwagon` internals. Doing so can create a security vulnerability in your implementation.

## Test & Benchmarks

To run the tests and benchmarks, run the following commands:
```bash
$ go test ./...
```

To run the benchmarks:
```bash
go test ./... -bench=. -run=none -benchmem
```

## Security

If you find any security vulnerability, please don't open a GH issue and contact repo owners directly.


## LICENSE

[MIT](LICENSE-MIT) and [Apache 2.0](LICENSE-APACHE).