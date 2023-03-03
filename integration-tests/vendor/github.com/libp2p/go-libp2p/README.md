
<h1 align="center">
  <a href="libp2p.io"><img width="250" src="https://github.com/libp2p/libp2p/blob/master/logo/black-bg-2.png?raw=true" alt="libp2p hex logo" /></a>
</h1>

<h3 align="center">The Go implementation of the libp2p Networking Stack.</h3>

<p align="center">
  <a href="http://protocol.ai"><img src="https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square" /></a>
  <a href="http://libp2p.io/"><img src="https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square" /></a>
  <a href="http://webchat.freenode.net/?channels=%23libp2p"><img src="https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square" /></a>
  <a href="https://waffle.io/libp2p/libp2p"><img src="https://img.shields.io/badge/pm-waffle-yellow.svg?style=flat-square" /></a>
  <a href="https://discuss.libp2p.io"><img src="https://img.shields.io/discourse/https/discuss.libp2p.io/posts.svg"/></a>
</p>

<p align="center">
  <a href="https://travis-ci.com/libp2p/go-libp2p"><img src="https://travis-ci.com/libp2p/go-libp2p.svg?branch=master" /></a>
  <br>
  <a href="https://github.com/RichardLitt/standard-readme"><img src="https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square" /></a>
  <a href="https://godoc.org/github.com/libp2p/go-libp2p"><img src="https://godoc.org/github.com/libp2p/go-libp2p?status.svg" /></a>
  <a href=""><img src="https://img.shields.io/badge/golang-%3E%3D1.8.0-orange.svg?style=flat-square" /></a>
  <br>
</p>

# Table of Contents

- [Background](#background)
- [Usage](#usage)
  - [API](#api)
  - [Examples](#examples)
- [Development](#development)
  - [Using the go-libp2p Workspace](#using-the-go-libp2p-workspace)
  - [About gx](#about-gx)
  - [Tests](#tests)
  - [Packages](#packages)
- [Contribute](#contribute)

## Background

[libp2p](https://github.com/libp2p/specs) is a networking stack and library modularized out of [The IPFS Project](https://github.com/ipfs/ipfs), and bundled separately for other tools to use.
>
libp2p is the product of a long, and arduous quest of understanding -- a deep dive into the internet's network stack, and plentiful peer-to-peer protocols from the past. Building large-scale peer-to-peer systems has been complex and difficult in the last 15 years, and libp2p is a way to fix that. It is a "network stack" -- a protocol suite -- that cleanly separates concerns, and enables sophisticated applications to only use the protocols they absolutely need, without giving up interoperability and upgradeability. libp2p grew out of IPFS, but it is built so that lots of people can use it, for lots of different projects.
>
> We will be writing a set of docs, posts, tutorials, and talks to explain what p2p is, why it is tremendously useful, and how it can help your existing and new projects. But in the meantime, check out
>
> - [**Our developing collection of docs**](https://docs.libp2p.io)
> - [**Our community discussion forums**](https://discuss.libp2p.io)
> - [**The libp2p Specification**](https://github.com/libp2p/specs)
> - [**go-libp2p implementation**](https://github.com/libp2p/go-libp2p)
> - [**js-libp2p implementation**](https://github.com/libp2p/js-libp2p)
> - [**rust-libp2p implementation**](https://github.com/libp2p/rust-libp2p)

## Usage

This repository (`go-libp2p`) serves as the entrypoint to the universe of modules that compose the Go implementation of the libp2p stack. Libp2p requires go 1.12+.

We mainly use [Go modules](https://github.com/golang/go/wiki/Modules) for our dependency and release management (and thus require go >= 1.12+). In order to get the best developer experience, we recommend you do too. Otherwise, you may ocassionally encounter a breaking build as you'll be running off master (which, by definition, is not guaranteed to be stable).

You can start using go-libp2p in your Go application simply by adding imports from our repos, e.g.:

```go
import "github.com/libp2p/go-libp2p"
```

Run `go get` or `go build`, excluding the libp2p repos from Go modules proxy usage. You only need to do this the first time you import go-libp2p to make sure you latch onto the correct version lineage (see [golang/go#34189](https://github.com/golang/go/issues/34189) for context):

```sh
$ GOPRIVATE='github.com/libp2p/*' go get ./...
```

The Go build tools will look for [available releases](https://github.com/libp2p/go-libp2p/releases), and will pick the highest available one.

As new releases of go-libp2p are made available, you can upgrade your application by manually editing your `go.mod` file, or using the [Go tools](https://golang.org/cmd/go/#hdr-Maintaining_module_requirements) to maintain module requirements.

### API

[![GoDoc](https://godoc.org/github.com/libp2p/go-libp2p?status.svg)](https://godoc.org/github.com/libp2p/go-libp2p)

### Examples

Examples can be found in the [examples repo](https://github.com/libp2p/go-libp2p-examples).

## Development

### Using the go-libp2p Workspace

While developing, you may need to make changes to several modules at once, or you may want changes made locally in one module to be available for import by another.

The [go-libp2p workspace](https://github.com/libp2p/workspace-go-libp2p) provides a developer-oriented view of the modules that comprise go-libp2p.

Using the tooling in the workspace repository, you can checkout all of go-libp2p's module repos and enter "local mode", which adds [replace directives](https://github.com/golang/go/wiki/Modules#gomod) to the go.mod files in each local working copy. When you build locally, the libp2p depdendencies will be resolved from your local working copies.

Once you've committed your changes, you can switch back to "remote mode", which removes the replace directives and pulls imports from the main go module cache.

See the [workspace repo](https://github.com/libp2p/workspace-go-libp2p) for more information.

### About gx

Before adopting gomod, libp2p used [gx](https://github.com/whyrusleeping/gx) to manage dependencies using [IPFS](https://ipfs.io).

Due to the difficulties in keeping both dependency management solutions up-to-date, gx support was ended in April 2019.

Ending gx support does not mean that existing gx builds will break. Because gx references dependencies by their immutable IPFS hash, any currently working gx builds will continue to work for as long as the dependencies are resolvable in IPFS.

However, new changes to go-libp2p will not be published via gx, and users are encouraged to adopt gomod to stay up-to-date.

If you experience any issues migrating from gx to gomod, please [join the discussion at the libp2p forums](https://discuss.libp2p.io/t/gomod-and-go-libp2p/44).

### Tests

`go test ./...` will run all tests in the repo.

### Packages

> This table is generated using the module [`package-table`](https://github.com/ipfs-shipyard/package-table) with `package-table --data=package-list.json`.

List of packages currently in existence for libp2p:

| Name | CI/Travis | Coverage | Description |
| ---------|---------|---------|--------- |
| **Libp2p** |
| [`go-libp2p`](//github.com/libp2p/go-libp2p) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p) | go-libp2p entry point |
| [`go-libp2p-core`](//github.com/libp2p/go-libp2p-core) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-core.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-core) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-core/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-core) | core interfaces, types, and abstractions |
| [`go-libp2p-blankhost`](//github.com/libp2p/go-libp2p-blankhost) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-blankhost.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-blankhost) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-blankhost/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-blankhost) | minimal implementation of the "host" interface |
| **Network** |
| [`go-libp2p-swarm`](//github.com/libp2p/go-libp2p-swarm) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-swarm.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-swarm) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-swarm/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-swarm) | reference implementation of network state machine |
| **Transport** |
| [`go-ws-transport`](//github.com/libp2p/go-ws-transport) | [![Travis CI](https://travis-ci.com/libp2p/go-ws-transport.svg?branch=master)](https://travis-ci.com/libp2p/go-ws-transport) | [![codecov](https://codecov.io/gh/libp2p/go-ws-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-ws-transport) | WebSocket transport |
| [`go-tcp-transport`](//github.com/libp2p/go-tcp-transport) | [![Travis CI](https://travis-ci.com/libp2p/go-tcp-transport.svg?branch=master)](https://travis-ci.com/libp2p/go-tcp-transport) | [![codecov](https://codecov.io/gh/libp2p/go-tcp-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-tcp-transport) | TCP transport |
| [`go-libp2p-quic-transport`](//github.com/libp2p/go-libp2p-quic-transport) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-quic-transport.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-quic-transport) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-quic-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-quic-transport) | QUIC transport |
| [`go-udp-transport`](//github.com/libp2p/go-udp-transport) | [![Travis CI](https://travis-ci.com/libp2p/go-udp-transport.svg?branch=master)](https://travis-ci.com/libp2p/go-udp-transport) | [![codecov](https://codecov.io/gh/libp2p/go-udp-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-udp-transport) | UDP transport |
| [`go-utp-transport`](//github.com/libp2p/go-utp-transport) | [![Travis CI](https://travis-ci.com/libp2p/go-utp-transport.svg?branch=master)](https://travis-ci.com/libp2p/go-utp-transport) | [![codecov](https://codecov.io/gh/libp2p/go-utp-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-utp-transport) | uTorrent transport (UTP) |
| [`go-libp2p-circuit`](//github.com/libp2p/go-libp2p-circuit) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-circuit.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-circuit) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-circuit/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-circuit) | relay transport |
| [`go-libp2p-transport-upgrader`](//github.com/libp2p/go-libp2p-transport-upgrader) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-transport-upgrader.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-transport-upgrader) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-transport-upgrader/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-transport-upgrader) | upgrades multiaddr-net connections into full libp2p transports |
| [`go-libp2p-reuseport-transport`](//github.com/libp2p/go-libp2p-reuseport-transport) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-reuseport-transport.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-reuseport-transport) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-reuseport-transport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-reuseport-transport) | partial transport for building transports that reuse ports |
| **Encrypted Channels** |
| [`go-libp2p-secio`](//github.com/libp2p/go-libp2p-secio) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-secio.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-secio) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-secio/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-secio) | SecIO crypto channel |
| [`go-libp2p-tls`](//github.com/libp2p/go-libp2p-tls) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-tls.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-tls) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-tls/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-tls) | TLS 1.3+ crypto channel |
| [`go-conn-security-multistream`](//github.com/libp2p/go-conn-security-multistream) | [![Travis CI](https://travis-ci.com/libp2p/go-conn-security-multistream.svg?branch=master)](https://travis-ci.com/libp2p/go-conn-security-multistream) | [![codecov](https://codecov.io/gh/libp2p/go-conn-security-multistream/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-conn-security-multistream) | multistream multiplexed meta crypto channel |
| **Private Network** |
| [`go-libp2p-pnet`](//github.com/libp2p/go-libp2p-pnet) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-pnet.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-pnet) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-pnet/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-pnet) | reference private networking implementation |
| **Stream Muxers** |
| [`go-libp2p-yamux`](//github.com/libp2p/go-libp2p-yamux) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-yamux.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-yamux) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-yamux/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-yamux) | YAMUX stream multiplexer |
| [`go-libp2p-mplex`](//github.com/libp2p/go-libp2p-mplex) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-mplex.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-mplex) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-mplex/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-mplex) | MPLEX stream multiplexer |
| **NAT Traversal** |
| [`go-libp2p-nat`](//github.com/libp2p/go-libp2p-nat) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-nat.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-nat) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-nat/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-nat) |  |
| **Peerstore** |
| [`go-libp2p-peerstore`](//github.com/libp2p/go-libp2p-peerstore) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-peerstore.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-peerstore) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-peerstore/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-peerstore) | reference implementation of peer metadata storage component |
| **Connection Manager** |
| [`go-libp2p-connmgr`](//github.com/libp2p/go-libp2p-connmgr) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-connmgr.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-connmgr) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-connmgr/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-connmgr) | reference implementation of connection manager |
| **Routing** |
| [`go-libp2p-record`](//github.com/libp2p/go-libp2p-record) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-record.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-record) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-record/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-record) | record type and validator logic |
| [`go-libp2p-kad-dht`](//github.com/libp2p/go-libp2p-kad-dht) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-kad-dht.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-kad-dht) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-kad-dht/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-kad-dht) | Kademlia-like router |
| [`go-libp2p-kbucket`](//github.com/libp2p/go-libp2p-kbucket) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-kbucket.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-kbucket) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-kbucket/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-kbucket) | Kademlia routing table helper types |
| [`go-libp2p-coral-dht`](//github.com/libp2p/go-libp2p-coral-dht) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-coral-dht.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-coral-dht) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-coral-dht/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-coral-dht) | Router based on Coral DHT |
| [`go-libp2p-pubsub-router`](//github.com/libp2p/go-libp2p-pubsub-router) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-pubsub-router.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-pubsub-router) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-pubsub-router/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-pubsub-router) | record-store over pubsub adapter |
| **Consensus** |
| [`go-libp2p-consensus`](//github.com/libp2p/go-libp2p-consensus) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-consensus.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-consensus) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-consensus/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-consensus) | consensus protocols interfaces |
| [`go-libp2p-raft`](//github.com/libp2p/go-libp2p-raft) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-raft.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-raft) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-raft/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-raft) | consensus implementation over raft |
| **Pubsub** |
| [`go-libp2p-pubsub`](//github.com/libp2p/go-libp2p-pubsub) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-pubsub.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-pubsub) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-pubsub/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-pubsub) | multiple pubsub over libp2p implementations |
| **RPC** |
| [`go-libp2p-gorpc`](//github.com/libp2p/go-libp2p-gorpc) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-gorpc.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-gorpc) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-gorpc/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-gorpc) | a simple RPC library for libp2p |
| **Utilities/miscellaneous** |
| [`go-libp2p-loggables`](//github.com/libp2p/go-libp2p-loggables) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-loggables.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-loggables) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-loggables/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-loggables) | logging helpers |
| [`go-maddr-filter`](//github.com/libp2p/go-maddr-filter) | [![Travis CI](https://travis-ci.com/libp2p/go-maddr-filter.svg?branch=master)](https://travis-ci.com/libp2p/go-maddr-filter) | [![codecov](https://codecov.io/gh/libp2p/go-maddr-filter/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-maddr-filter) | multiaddr filtering helpers |
| [`go-libp2p-netutil`](//github.com/libp2p/go-libp2p-netutil) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-netutil.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-netutil) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-netutil/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-netutil) | misc utilities |
| [`go-msgio`](//github.com/libp2p/go-msgio) | [![Travis CI](https://travis-ci.com/libp2p/go-msgio.svg?branch=master)](https://travis-ci.com/libp2p/go-msgio) | [![codecov](https://codecov.io/gh/libp2p/go-msgio/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-msgio) | length prefixed data channel |
| [`go-addr-util`](//github.com/libp2p/go-addr-util) | [![Travis CI](https://travis-ci.com/libp2p/go-addr-util.svg?branch=master)](https://travis-ci.com/libp2p/go-addr-util) | [![codecov](https://codecov.io/gh/libp2p/go-addr-util/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-addr-util) | address utilities for libp2p swarm |
| [`go-buffer-pool`](//github.com/libp2p/go-buffer-pool) | [![Travis CI](https://travis-ci.com/libp2p/go-buffer-pool.svg?branch=master)](https://travis-ci.com/libp2p/go-buffer-pool) | [![codecov](https://codecov.io/gh/libp2p/go-buffer-pool/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-buffer-pool) | a variable size buffer pool for go |
| [`go-libp2p-routing-helpers`](//github.com/libp2p/go-libp2p-routing-helpers) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-routing-helpers.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-routing-helpers) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-routing-helpers/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-routing-helpers) | routing helpers |
| [`go-reuseport`](//github.com/libp2p/go-reuseport) | [![Travis CI](https://travis-ci.com/libp2p/go-reuseport.svg?branch=master)](https://travis-ci.com/libp2p/go-reuseport) | [![codecov](https://codecov.io/gh/libp2p/go-reuseport/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-reuseport) | enables reuse of addresses |
| [`go-sockaddr`](//github.com/libp2p/go-sockaddr) | [![Travis CI](https://travis-ci.com/libp2p/go-sockaddr.svg?branch=master)](https://travis-ci.com/libp2p/go-sockaddr) | [![codecov](https://codecov.io/gh/libp2p/go-sockaddr/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-sockaddr) | utils for sockaddr conversions |
| [`go-flow-metrics`](//github.com/libp2p/go-flow-metrics) | [![Travis CI](https://travis-ci.com/libp2p/go-flow-metrics.svg?branch=master)](https://travis-ci.com/libp2p/go-flow-metrics) | [![codecov](https://codecov.io/gh/libp2p/go-flow-metrics/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-flow-metrics) | metrics library |
| [`go-libp2p-gostream`](//github.com/libp2p/go-libp2p-gostream) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-gostream.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-gostream) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-gostream/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-gostream) | Go 'net' wrappers for libp2p |
| [`go-libp2p-http`](//github.com/libp2p/go-libp2p-http) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-http.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-http) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-http/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-http) | HTTP on top of libp2p streams |
| **Testing and examples** |
| [`go-libp2p-testing`](//github.com/libp2p/go-libp2p-testing) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-testing.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-testing) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-testing/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-testing) | a collection of testing utilities for libp2p |
| [`go-libp2p-examples`](//github.com/libp2p/go-libp2p-examples) | [![Travis CI](https://travis-ci.com/libp2p/go-libp2p-examples.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-examples) | [![codecov](https://codecov.io/gh/libp2p/go-libp2p-examples/branch/master/graph/badge.svg)](https://codecov.io/gh/libp2p/go-libp2p-examples) | go-libp2p examples and tutorials |

# Contribute

go-libp2p is part of [The IPFS Project](https://github.com/ipfs/ipfs), and is MIT-licensed open source software. We welcome contributions big and small! Take a look at the [community contributing notes](https://github.com/ipfs/community/blob/master/CONTRIBUTING.md). Please make sure to check the [issues](https://github.com/ipfs/go-libp2p/issues). Search the closed ones before reporting things, and help us with the open ones.

Guidelines:

- read the [libp2p spec](https://github.com/libp2p/specs)
- please make branches + pull-request, even if working on the main repository
- ask questions or talk about things in [issues](https://github.com/libp2p/go-libp2p/issues), our [discussion forums](https://discuss.libp2p.io), or #libp2p or #ipfs on freenode.
- ensure you are able to contribute (no legal issues please -- we use the DCO)
- run `go fmt` before pushing any code
- run `golint` and `go vet` too -- some things (like protobuf files) are expected to fail.
- get in touch with @raulk and @mgoelzer about how best to contribute
- have fun!

There's a few things you can do right now to help out:
 - Go through the modules below and **check out existing issues**. This would be especially useful for modules in active development. Some knowledge of IPFS/libp2p may be required, as well as the infrasture behind it - for instance, you may need to read up on p2p and more complex operations like muxing to be able to help technically.
 - **Perform code reviews**.
 - **Add tests**. There can never be enough tests.

---

The last gx published version of this module was: 6.0.41: QmTRN7hRxvGkxKxDdeudty7sRet4L7ZKZCqKsXHa79wmAc
