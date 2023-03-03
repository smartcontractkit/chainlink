# go-reuseport-transport

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io/)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23libp2p)
[![Discourse posts](https://img.shields.io/discourse/https/discuss.libp2p.io/posts.svg)](https://discuss.libp2p.io)
[![GoDoc](https://godoc.org/github.com/libp2p/go-reuseport-transport?status.svg)](https://godoc.org/github.com/libp2p/go-reuseport-transport)
[![Build Status](https://travis-ci.org/libp2p/go-reuseport-transport.svg?branch=master)](https://travis-ci.org/libp2p/go-reuseport-transport)

> Basic reuseport TCP transport

This package provides a basic transport for automatically (and intelligently) reusing TCP ports.

To use, construct a new `Transport` (the zero value is safe to use) and configure any listeners (`tr.Listen(...)`).

Then, when dialing (`tr.Dial(...)`), the transport will attempt to reuse the ports it's currently listening on, choosing the best one depending on the destination address.


NOTE: Currently, we don't make any attempts to prevent two reusport transports from interfering with each other (reusing each other's ports). However, we reserve the right to fix this in the future.

## Install

`go-reuseport-transport` is a standard Go module which can be installed with:

```sh
go get github.com/libp2p/go-reuseport-transport
```

This repo is [gomod](https://github.com/golang/go/wiki/Modules)-compatible, and users of
go 1.11 and later with modules enabled will automatically pull the latest tagged release
by referencing this package. Upgrades to future releases can be managed using `go get`,
or by editing your `go.mod` file as [described by the gomod documentation](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies).

## Current use in libp2p

This package is *currently* used by the [go-tcp-transport](https://github.com/libp2p/go-tcp-transport) libp2p transport and will likely be used by more libp2p transports in the future.

## Contribute

Feel free to join in. All welcome. Open an [issue](https://github.com/libp2p/go-reuseport-transport/issues)!

This repository falls under the libp2p [Code of Conduct](https://github.com/libp2p/community/blob/master/code-of-conduct.md).

### Want to hack on libp2p?

[![](https://cdn.rawgit.com/libp2p/community/master/img/contribute.gif)](https://github.com/libp2p/community/blob/master/CONTRIBUTE.md)

## License

MIT

---

The last gx published version of this module was: 0.2.3: QmTmbamDjDWgHe8qeMt7ZpaeNpTj349JpFKuwTF321XavT
