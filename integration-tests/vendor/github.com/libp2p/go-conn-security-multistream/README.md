# go-conn-security-multistream

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](https://protocol.ai)
[![](https://img.shields.io/badge/project-libp2p-yellow.svg?style=flat-square)](https://libp2p.io/)
[![](https://img.shields.io/badge/freenode-%23libp2p-yellow.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23libp2p)
[![Discourse posts](https://img.shields.io/discourse/https/discuss.libp2p.io/posts.svg)](https://discuss.libp2p.io)
[![GoDoc](https://godoc.org/github.com/libp2p/go-conn-security-multistream?status.svg)](https://godoc.org/github.com/libp2p/go-conn-security-multistream)

> Connection security multistream multiplexer

This package provides a multistream multiplexed [security transport](https://github.com/libp2p/go-conn-security). It:

1. Selects a security transport using multistream-select.
2. Secures the stream using the selected transport.

Known libp2p security transports include:

* [go-libp2p-secio](https://github.com/libp2p/go-libp2p-secio)
* [go-libp2p-tls](https://github.com/libp2p/go-libp2p-tls)

## Install

`go-conn-security-multistream` is a standard Go module which can be installed with:

```sh
go get github.com/libp2p/go-conn-security-multistream
```

This repo is [gomod](https://github.com/golang/go/wiki/Modules)-compatible, and users of
go 1.11 and later with modules enabled will automatically pull the latest tagged release
by referencing this package. Upgrades to future releases can be managed using `go get`,
or by editing your `go.mod` file as [described by the gomod documentation](https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies).

## Usage

For more information about how `go-conn-security-multistream` is used in the libp2p context, you can see the [go-libp2p-conn](https://github.com/libp2p/go-libp2p-conn) module.

## Contribute

Feel free to join in. All welcome. Open an [issue](https://github.com/libp2p/go-conn-security-multistream/issues)!

This repository falls under the libp2p [Code of Conduct](https://github.com/libp2p/community/blob/master/code-of-conduct.md).

### Want to hack on libp2p?

[![](https://cdn.rawgit.com/libp2p/community/master/img/contribute.gif)](https://github.com/libp2p/community/blob/master/CONTRIBUTE.md)

## License

MIT

---

The last gx published version of this module was: 0.1.26: QmZWmFkMm28sWeDr5Xh1LexdKBGYGp946MNCfgtLqfX73z
