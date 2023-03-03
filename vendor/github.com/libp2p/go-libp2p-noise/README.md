# go-libp2p-noise

[![](https://img.shields.io/badge/made%20by-ETHBerlinZwei-blue.svg?style=flat-square)](https://ethberlinzwei.com)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io/)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23libp2p)
[![Discourse posts](https://img.shields.io/discourse/https/discuss.libp2p.io/posts.svg)](https://discuss.libp2p.io)
[![GoDoc](https://godoc.org/github.com/libp2p/go-libp2p-noise?status.svg)](https://godoc.org/github.com/libp2p/go-libp2p-noise)
[![Build Status](https://travis-ci.com/libp2p/go-libp2p-noise.svg?branch=master)](https://travis-ci.com/libp2p/go-libp2p-noise)

> go-libp2p's noise encrypted transport

`go-libp2p-noise` is a component of the [libp2p project](https://libp2p.io), a
modular networking stack for developing peer-to-peer applications. It provides a
secure transport channel for [`go-libp2p`][go-libp2p] based on the 
[Noise Protocol Framework](https://noiseprotocol.org). Following an initial
plaintext handshake, all data exchanged between peers using `go-libp2p-noise` is
encrypted and protected from eavesdropping.

libp2p supports multiple [transport protocols][docs-transport], many of which
lack native channel security. `go-libp2p-noise` is designed to work with
go-libp2p's ["transport upgrader"][transport-upgrader], which applies security
modules (like `go-libp2p-noise`) to an insecure channel. `go-libp2p-noise`
implements the [`SecureTransport` interface][godoc-securetransport], which
allows the upgrader to secure any underlying connection.

More detail on the handshake protocol and wire format used is available in the
[noise-libp2p specification][noise-libp2p-spec]. Details about security protocol
negotiation in libp2p can be found in the [connection establishment spec][conn-spec].

## Status

This implementation is currently considered "feature complete," but it has not yet
been widely tested in a production environment. 

## Install

As `go-libp2p-noise` is still in development, it is not included as a default dependency of `go-libp2p`.

`go-libp2p-noise` is a standard Go module which can be installed with:

```sh
go get github.com/libp2p/go-libp2p-noise
```

This repo is [gomod](https://github.com/golang/go/wiki/Modules) compatible, and users of
go 1.12 and later with modules enabled will automatically pull the latest tagged release
by referencing this package. Upgrades to future releases can be managed using `go get`,
or by editing your `go.mod` file as [described by the gomod documentation](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies).

## Usage

`go-libp2p-noise` is not currently enabled by default when constructing a new libp2p
[Host][godoc-host], so you will need to explicitly enable it in order to use it.

To do so, you can pass `noise.New` as an argument to a `libp2p.Security` `Option` when
constructing a libp2p `Host` with `libp2p.New`:

```go
import (
  libp2p "github.com/libp2p/go-libp2p"
  noise "github.com/libp2p/go-libp2p-noise"
)

// wherever you create your libp2p instance:
host := libp2p.New(
  libp2p.Security(noise.ID, noise.New)
)
```

Note that the above snippet will _replace_ the default security protocols. To add Noise
as an additional protocol, chain it to the default options instead:

```go
libp2p.ChainOptions(libp2p.DefaultSecurity, libp2p.Security(noise.ID, noise.New))
```

## Contribute

Feel free to join in. All welcome. Open an [issue](https://github.com/libp2p/go-libp2p-noise/issues)!

This repository falls under the libp2p [Code of Conduct](https://github.com/libp2p/community/blob/master/code-of-conduct.md).

### Want to hack on libp2p?

[![](https://cdn.rawgit.com/libp2p/community/master/img/contribute.gif)](https://github.com/libp2p/community/blob/master/CONTRIBUTE.md)

## License

MIT

---

[go-libp2p]: https://github.com/libp2p/go-libp2p
[noise-libp2p-spec]: https://github.com/libp2p/specs/blob/master/noise/README.md
[conn-spec]: https://github.com/libp2p/specs/blob/master/connections/README.md
[docs-transport]: https://docs.libp2p.io/concepts/transport
[transport-upgrader]: https://github.com/libp2p/go-libp2p-transport-upgrader
[godoc-host]: https://godoc.org/github.com/libp2p/go-libp2p-core/host#Host
[godoc-option]: https://godoc.org/github.com/libp2p/go-libp2p#Option
[godoc-go-libp2p-pkg-vars]: https://godoc.org/github.com/libp2p/go-libp2p#pkg-variables 
[godoc-security-option]: https://godoc.org/github.com/libp2p/go-libp2p#Security
[godoc-securetransport]: https://godoc.org/github.com/libp2p/go-libp2p-core/sec#SecureTransport

