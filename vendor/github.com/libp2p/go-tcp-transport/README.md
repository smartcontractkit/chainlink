go-tcp-transport
==================

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io/)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23libp2p)
[![Discourse posts](https://img.shields.io/discourse/https/discuss.libp2p.io/posts.svg)](https://discuss.libp2p.io)
[![Coverage Status](https://coveralls.io/repos/github/libp2p/go-tcp-transport/badge.svg?branch=master)](https://coveralls.io/github/libp2p/go-tcp-transport?branch=master)
[![Travis CI](https://travis-ci.com/libp2p/go-tcp-transport.svg?branch=master)](https://travis-ci.com/libp2p/go-tcp-transport)

> A libp2p transport implementation for tcp, including reuseport socket options.

`go-tcp-transport` is an implementation of the [libp2p transport
interface][concept-transport] that streams data over TCP/IP sockets. It is
included by default in the main [`go-libp2p`][go-libp2p] "entry point" module.

## Table of Contents

- [go-tcp-transport](#go-tcp-transport)
    - [Table of Contents](#table-of-contents)
    - [Install](#install)
    - [Usage](#usage)
    - [Security and Multiplexing](#security-and-multiplexing)
    - [reuseport](#reuseport)
    - [Contribute](#contribute)
    - [License](#license)

## Install

`go-tcp-transport` is included as a dependency of `go-libp2p`, which is the most
common libp2p entry point. If you depend on `go-libp2p`, there is generally no
need to explicitly depend on this module.

`go-tcp-transport` is a standard Go module which can be installed with:

``` sh
go get github.com/libp2p/go-tcp-transport
```


This repo is [gomod](https://github.com/golang/go/wiki/Modules)-compatible, and users of
go 1.11 and later with modules enabled will automatically pull the latest tagged release
by referencing this package. Upgrades to future releases can be managed using `go get`,
or by editing your `go.mod` file as [described by the gomod documentation](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies).

## Usage

TCP is one of the default transports enabled when constructing a standard libp2p
Host, along with [WebSockets](https://github.com/libp2p/go-ws-transport).

Calling [`libp2p.New`][godoc-libp2p-new] to construct a libp2p Host will enable
the TCP transport, unless you override the default transports by passing in
`Options` to `libp2p.New`.

To explicitly enable the TCP transport while constructing a host, use the
`libp2p.Transport` option, passing in the `NewTCPTransport` constructor function:

``` go

import (
    "context"

    libp2p "github.com/libp2p/go-libp2p"
    tcp "github.com/libp2p/go-tcp-transport"
)

ctx := context.Background()

// TCP only:
h, err := libp2p.New(ctx,
    libp2p.Transport(tcp.NewTCPTransport)
)
```

The example above will replace the default transports with a single TCP
transport. To add multiple tranports, use `ChainOptions`:

``` go
// TCP and QUIC:
h, err := libp2p.New(ctx,
    libp2p.ChainOptions(
        libp2p.Transport(tcp.NewTCPTransport),
        libp2p.Transport(quic.NewTransport)) // see https://github.com/libp2p/go-libp2p-quic-transport
)
```

## Addresses

The TCP transport supports [multiaddrs][multiaddr] that contain a `tcp`
component, provided that there is sufficient addressing information for the IP
layer of the connection.

Examples:

| addr                       | description                                        |
|----------------------------|----------------------------------------------------|
| `/ip4/1.2.3.4/tcp/1234`    | IPv4: 1.2.3.4, TCP port 1234                       |
| `/ip6/::1/tcp/1234`        | IPv6 loopback, TCP port 1234                       |
| `/dns4/example.com/tcp/80` | DNS over IPv4, hostname `example.com`, TCP port 80 |


Support for IP layer protocols is provided by the
[go-multiaddr-net](https://github.com/multiformats/go-multiaddr-net) module.

## Security and Multiplexing

Because TCP lacks native connection security and stream multiplexing facilities,
the TCP transport uses a [transport upgrader][transport-upgrader] to provide
those features. The transport upgrader negotiates transport security and
multiplexing for each connection according to the protocols supported by each
party.

## reuseport

The [`SO_REUSEPORT`][explain-reuseport] socket option allows multiple processes
or threads to bind to the same TCP port, provided that all of them set the
socket option. This has some performance benefits, and it can potentially assist
in NAT traversal by only requiring one port to be accessible for many
connections.

The reuseport functionality is provided by a seperate module,
[go-reuseport-transport](https://github.com/libp2p/go-reuseport-transport). It
is enabled by default, but can be disabled at runtime by setting the
`LIBP2P_TCP_REUSEPORT` environment variable to `false` or `0`.

## Contribute

PRs are welcome!

Small note: If editing the Readme, please conform to the [standard-readme](https://github.com/RichardLitt/standard-readme) specification.

## License

MIT © Jeromy Johnson

---

The last gx published version of this module was: 2.0.28: QmTGiDkw4eeKq31wwpQRk5GwWiReaxrcTQLuCCLWgfKo5M

<!-- reference links -->
[go-libp2p]: https://github.com/libp2p/go-libp2p
[concept-transport]: https://docs.libp2p.io/concepts/transport/
[interface-host]: https://github.com/libp2p/go-libp2p-core/blob/master/host/host.go
[godoc-libp2p-new]: https://godoc.org/github.com/libp2p/go-libp2p#New
[transport-upgrader]: https://github.com/libp2p/go-libp2p-transport-upgrader
[explain-reuseport]: https://lwn.net/Articles/542629/
[multiaddr]: https://github.com/multiformats/multiaddr
