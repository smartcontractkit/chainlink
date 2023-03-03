# go-ws-transport

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io/)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23libp2p)
[![Discourse posts](https://img.shields.io/discourse/https/discuss.libp2p.io/posts.svg)](https://discuss.libp2p.io)
[![GoDoc](https://godoc.org/github.com/libp2p/go-ws-transport?status.svg)](https://godoc.org/github.com/libp2p/go-ws-transport)
[![Coverage Status](https://coveralls.io/repos/github/libp2p/go-ws-transport/badge.svg?branch=master)](https://coveralls.io/github/libp2p/go-ws-transport?branch=master)
[![Build Status](https://travis-ci.org/libp2p/go-ws-transport.svg?branch=master)](https://travis-ci.org/libp2p/go-ws-transport)

> A libp2p transport implementation using WebSockets

`go-ws-transport` is an implementation of the [libp2p transport
interface][concept-transport] that streams data over
[WebSockets][spec-websockets], which are themselves layered over TCP/IP. It is
included by default in the main [`go-libp2p`][go-libp2p] "entry point" module.

## Table of Contents

- [go-ws-transport](#go-ws-transport)
    - [Table of Contents](#table-of-contents)
    - [Install](#install)
    - [Usage](#usage)
    - [Addresses](#addresses)
    - [Security and Multiplexing](#security-and-multiplexing)
    - [Contribute](#contribute)
        - [Want to hack on IPFS?](#want-to-hack-on-ipfs)
    - [License](#license)

## Install

`go-ws-transport` is included as a dependency of `go-libp2p`, which is the most
common libp2p entry point. If you depend on `go-libp2p`, there is generally no
need to explicitly depend on this module.

`go-ws-transport` is a standard Go module which can be installed with:

```sh
> go get github.com/libp2p/go-ws-transport
```

This repo is [gomod](https://github.com/golang/go/wiki/Modules)-compatible, and users of
go 1.11 and later with modules enabled will automatically pull the latest tagged release
by referencing this package. Upgrades to future releases can be managed using `go get`,
or by editing your `go.mod` file as [described by the gomod documentation](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies).

## Usage

WebSockets are one of the default transports enabled when constructing a standard libp2p
Host, along with [TCP](https://github.com/libp2p/go-tcp-transport).

Calling [`libp2p.New`][godoc-libp2p-new] to construct a libp2p Host will enable
the WebSocket transport, unless you override the default transports by passing in
`Options` to `libp2p.New`.

To explicitly enable the WebSocket transport while constructing a host, use the
`libp2p.Transport` option, passing in the `ws.New` constructor function:

``` go

import (
    "context"

    libp2p "github.com/libp2p/go-libp2p"
    ws "github.com/libp2p/go-ws-transport"
)

ctx := context.Background()

// WebSockets only:
h, err := libp2p.New(ctx,
    libp2p.Transport(ws.New)
)
```

The example above will replace the default transports with a single WebSocket
transport. To add multiple tranports, use `ChainOptions`:

``` go
// WebSockets and QUIC:
h, err := libp2p.New(ctx,
    libp2p.ChainOptions(
        libp2p.Transport(ws.New),
        libp2p.Transport(quic.NewTransport)) // see https://github.com/libp2p/go-libp2p-quic-transport
)
```

## Addresses

The WebSocket transport supports [multiaddrs][multiaddr] that contain a `ws`
component, which is encapsulated within (or layered onto) another valid TCP
multiaddr.

Examples:

| addr                          | description                                        |
|-------------------------------|----------------------------------------------------|
| `/ip4/1.2.3.4/tcp/1234/ws`    | IPv4: 1.2.3.4, TCP port 1234                       |
| `/ip6/::1/tcp/1234/ws`        | IPv6 loopback, TCP port 1234                       |
| `/dns4/example.com/tcp/80/ws` | DNS over IPv4, hostname `example.com`, TCP port 80 |

Notice that the `/ws` multiaddr component contextualizes an existing TCP/IP
multiaddr and does not require any additional addressing information.

## Security and Multiplexing

While the WebSocket spec defines a `wss` URI scheme for encrypted WebSocket
connections, support for `wss` URIs relies on TLS, which wraps the WebSocket
connection in a similar manner to TLS-protected HTTP traffic.

As libp2p does not integrate with the TLS Certificate Authority infrastructure
by design, security for WebSockets is provided by a [transport
upgrader][transport-upgrader]. The transport upgrader negotiates transport
security for each connection according to the protocols supported by each party.

The transport upgrader also negotiates a stream multiplexing protocol to allow
many bidirectional streams to coexist on a single WebSocket connection.

## Contribute

Feel free to join in. All welcome. Open an [issue](https://github.com/libp2p/go-ws-transport/issues)!

This repository falls under the libp2p [Code of Conduct](https://github.com/libp2p/community/blob/master/code-of-conduct.md).

### Want to hack on libp2p?

[![](https://cdn.rawgit.com/libp2p/community/master/img/contribute.gif)](https://github.com/libp2p/community/blob/master/CONTRIBUTE.md)

## License

MIT

---

The last gx published version of this module was: 2.0.27: QmaSWc4ox6SZQF6DHZvDuM9sP1syNajkKuPXmKR1t5BAz5

<!-- reference links -->
[go-libp2p]: https://github.com/libp2p/go-libp2p
[concept-transport]: https://docs.libp2p.io/concepts/transport/
[interface-host]: https://github.com/libp2p/go-libp2p-core/blob/master/host/host.go
[godoc-libp2p-new]: https://godoc.org/github.com/libp2p/go-libp2p#New
[transport-upgrader]: https://github.com/libp2p/go-libp2p-transport-upgrader
[multiaddr]: https://github.com/multiformats/multiaddr
[spec-websockets]: https://tools.ietf.org/html/rfc6455
