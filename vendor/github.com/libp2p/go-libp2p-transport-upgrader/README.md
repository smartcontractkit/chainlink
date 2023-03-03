# go-libp2p-transport-upgrader

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io/)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23libp2p)
[![GoDoc](https://godoc.org/github.com/libp2p/go-libp2p-transport-upgrader?status.svg)](https://godoc.org/github.com/libp2p/go-libp2p-transport-upgrader)
[![Build Status](https://travis-ci.org/libp2p/go-libp2p-transport-upgrader.svg?branch=master)](https://travis-ci.org/libp2p/go-libp2p-transport-upgrader)
[![Discourse posts](https://img.shields.io/discourse/https/discuss.libp2p.io/posts.svg)](https://discuss.libp2p.io)

> Add encryption and multiplexing capabilities to libp2p transport connections

This package is a component of [libp2p](https://libp2p.io), a modular networking
stack for building peer-to-peer applications.

For two libp2p peers to communicate, the connection between them must be secure,
and each peer must be able to open multiple independent streams of communication
over a single channel. We call connections with these features "capable"
connections.

Many of the underlying [transport protocols][docs-transport] that are used by
libp2p do not provide the required capabilities "out of the box."
`go-libp2p-transport-upgrader` provides the necessary logic to upgrade
connections and listeners into fully capable connections and connection
listeners.

In order to be upgraded, the underlying connection or listener must be a
[`multiaddr-net`][manet] [`Conn`][manet-conn] or [`Listener`][manet-listener].
The `multiaddr-net` types integrate the Go standard library connection types
with [`multiaddr`][multiaddr], an extensible addressing format used throughout
libp2p.

As well as the mandatory capabilities of security and multiplexing, the upgrader
can optionally apply a `Protector` for [private networking][pnet], as well as an
[address filter][maddr-filter] to prevent connections to specific addresses.

## Install

Most people building applications with libp2p will have no need to install
`go-libp2p-transport-upgrader` directly. It is included as a dependency of the
main [`go-libp2p`][go-libp2p] "entry point" module and is integrated into the
libp2p `Host`.

For users who do not depend on `go-libp2p` and are managing their libp2p module
dependencies in a more manual fashion, `go-libp2p-transport-upgrader` is a
standard Go module which can be installed with:

```sh
go get github.com/libp2p/go-libp2p-transport-upgrader
```

This repo is [gomod](https://github.com/golang/go/wiki/Modules)-compatible, and users of
go 1.11 and later with modules enabled will automatically pull the latest tagged release
by referencing this package. Upgrades to future releases can be managed using `go get`,
or by editing your `go.mod` file as [described by the gomod documentation](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies).

## Usage

To use, construct a new `Upgrader` with:

* An optional [pnet][pnet] `Protector`.
* An optional [go-maddr-filter][maddr-filter] address `Filter`.
* A mandatory [stream security transport][ss].
* A mandatory [stream multiplexer][smux].

In practice, most users will not need to construct an `Upgrader` directly.
Instead, when constructing a libp2p [`Host`][godoc-host], you can pass in some
combination of the [`PrivateNetwork`][godoc-pnet-option],
[`Filters`][godoc-filters-option], [`Security`][godoc-security-option], and
[`Muxer`][godoc-muxer-option] `Option`s. This will configure the `Upgrader` that
is created and used by the `Host` internally.

## Example

Below is a simplified TCP transport implementation using the transport upgrader.
In practice, you'll want to use
[go-tcp-transport](https://github.com/libp2p/go-tcp-transport), which is
optimized for production usage. 

```go
package tcptransport

import (
	"context"

	tptu "github.com/libp2p/go-libp2p-transport-upgrader"

	ma "github.com/multiformats/go-multiaddr"
	mafmt "github.com/multiformats/go-multiaddr-fmt"
	manet "github.com/multiformats/go-multiaddr-net"
	tpt "github.com/libp2p/go-libp2p-core/transport"
	peer "github.com/libp2p/go-libp2p-core/peer"
)

// TcpTransport is a simple TCP transport.
type TcpTransport struct {
	// Connection upgrader for upgrading insecure stream connections to
	// secure multiplex connections.
	Upgrader *tptu.Upgrader
}

var _ tpt.Transport = &TcpTransport{}

// NewTCPTransport creates a new TCP transport instance.
func NewTCPTransport(upgrader *tptu.Upgrader) *TcpTransport {
	return &TcpTransport{Upgrader: upgrader}
}

// CanDial returns true if this transport believes it can dial the given
// multiaddr.
func (t *TcpTransport) CanDial(addr ma.Multiaddr) bool {
	return mafmt.TCP.Matches(addr)
}

// Dial dials the peer at the remote address.
func (t *TcpTransport) Dial(ctx context.Context, raddr ma.Multiaddr, p peer.ID) (tpt.CapableConn, error) {
	var dialer manet.Dialer
	conn, err := dialer.DialContext(ctx, raddr)
	if err != nil {
		return nil, err
	}
	return t.Upgrader.UpgradeOutbound(ctx, t, conn, p)
}

// Listen listens on the given multiaddr.
func (t *TcpTransport) Listen(laddr ma.Multiaddr) (tpt.Listener, error) {
	list, err := manet.Listen(laddr)
	if err != nil {
		return nil, err
	}
	return t.Upgrader.UpgradeListener(t, list), nil
}

// Protocols returns the list of terminal protocols this transport can dial.
func (t *TcpTransport) Protocols() []int {
	return []int{ma.P_TCP}
}

// Proxy always returns false for the TCP transport.
func (t *TcpTransport) Proxy() bool {
	return false
}

```

## Contribute

Feel free to join in. All welcome. Open an [issue](https://github.com/libp2p/go-libp2p-transport-upgrader/issues)!

This repository falls under the libp2p [Code of Conduct](https://github.com/libp2p/community/blob/master/code-of-conduct.md).

### Want to hack on libp2p?

[![](https://cdn.rawgit.com/libp2p/community/master/img/contribute.gif)](https://github.com/libp2p/community/blob/master/CONTRIBUTE.md)

## License

MIT

---

The last gx published version of this module was: 0.1.28: QmeqC5shQjEBRG9B8roZqQCJ9xb7Pq6AbWxJFMyLgqBBWh

[tpt]: https://godoc.org/github.com/libp2p/go-libp2p-core/transport
[manet]: https://github.com/multiformats/go-multiaddr-net
[ss]: https://godoc.org/github.com/libp2p/go-libp2p-core/sec
[smux]: https://godoc.org/github.com/libp2p/go-libp2p-core/mux
[pnet]: https://godoc.org/github.com/libp2p/go-libp2p-core/pnet
[manet-conn]: https://godoc.org/github.com/multiformats/go-multiaddr-net#Conn
[manet-listener]: https://godoc.org/github.com/multiformats/go-multiaddr-net#Listener
[maddr-filter]: https://github.com/libp2p/go-maddr-filter
[docs-transport]: https://docs.libp2p.io/concepts/transport
[multiaddr]: https://github.com/multiformats/multiaddr
[go-libp2p]: https://github.com/lib2p2/go-libp2p
[godoc-host]: https://godoc.org/github.com/libp2p/go-libp2p-core/host#Host
[godoc-pnet-option]: https://godoc.org/github.com/libp2p/go-libp2p#PrivateNetwork
[godoc-filters-option]: https://godoc.org/github.com/libp2p/go-libp2p#Filters
[godoc-security-option]: https://godoc.org/github.com/libp2p/go-libp2p#Security
[godoc-muxer-option]: https://godoc.org/github.com/libp2p/go-libp2p#Muxer
