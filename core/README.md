<br/>
<p align="center">
<a href="https://chain.link" target="_blank">
<img src="https://raw.githubusercontent.com/smartcontractkit/chainlink/develop/docs/logo-chainlink-blue.svg" width="225" alt="Chainlink logo">
</a>
</p>
<br/>

[![Go Report Card](https://goreportcard.com/badge/github.com/smartcontractkit/chainlink)](https://goreportcard.com/report/github.com/smartcontractkit/chainlink)
[![GoDoc](https://godoc.org/github.com/smartcontractkit/chainlink?status.svg)](https://godoc.org/github.com/smartcontractkit/chainlink)

Chainlink Core is the API backend that Chainlink client contracts on Ethereum
make requests to. The backend utilizes Solidity contract ABIs to generate types
for interacting with Ethereum contracts.

## Features

* Headless API implementation
* CLI tool providing conveniance commands for node configuration, administration,
  and CRUD object operations (e.g. Jobs, Runs, and even the VRF)

## Installation

See the [root README](../README.md#install)
for instructions on how to build the full Chainlink node.

## Directory Structure

This directory contains the majority of the code for the backend of Chainlink.

Static assets are pulled in using Go's [`embed`](https://pkg.go.dev/embed) package
and included in the final binary.

## Common Commands

**Install:**

By default `go install` will install this directory under the name `core`.
You can instead, build it, and place it in your path as `chainlink`:

```sh
go build -o $GOPATH/bin/chainlink .
```

**Test:**

```sh
# A higher parallel number can speed up tests at the expense of more RAM.
go test -p 1 ./...
```

This excludes more extensive integration tests which require a bit more setup, head over to [./integration-tests]
(../integration-tests/README.md) for more details on running those.

The golang testsuite is almost entirely parallelizable, and so running the default
`go test ./...` will commonly peg your processor. Limit parallelization with the
`-p 2` or whatever best fits your computer: `go test -p 4 ./...`.
