# go-ipns

[![](https://img.shields.io/badge/made%20by-Protocol%20Labs-blue.svg?style=flat-square)](http://protocol.ai)
[![](https://img.shields.io/badge/project-IPFS-blue.svg?style=flat-square)](http://ipfs.io/)
[![](https://img.shields.io/badge/freenode-%23ipfs-blue.svg?style=flat-square)](http://webchat.freenode.net/?channels=%23ipfs)
[![standard-readme compliant](https://img.shields.io/badge/standard--readme-OK-green.svg?style=flat-square)](https://github.com/RichardLitt/standard-readme)
[![GoDoc](https://godoc.org/github.com/ipfs/go-datastore?status.svg)](https://godoc.org/github.com/ipfs/go-ipns)

> ipns record definitions

This package contains all of the components necessary to create, understand, and validate IPNS records. It does *not* publish or resolve those records. [`go-ipfs`](https://github.com/ipfs/go-ipfs) uses this package internally to manipulate records.

## Lead Maintainer

[Adin Schmahmann](https://github.com/aschmahmann)

## Usage

To create a new IPNS record:

```go
import (
	"time"

	ipns "github.com/ipfs/go-ipns"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

// Generate a private key to sign the IPNS record with. Most of the time, 
// however, you'll want to retrieve an already-existing key from IPFS using the
// go-ipfs/core/coreapi CoreAPI.KeyAPI() interface.
privateKey, publicKey, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
if err != nil {
  panic(err)
}

// Create an IPNS record that expires in one hour and points to the IPFS address
// /ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5
ipnsRecord, err := ipns.Create(privateKey, []byte("/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5"), 0, time.Now().Add(1*time.Hour))
if err != nil {
	panic(err)
}
```

Once you have the record, you’ll need to use IPFS to *publish* it.

There are several other major operations you can do with `go-ipns`. Check out the [API docs](https://godoc.org/github.com/ipfs/go-ipns) or look at the tests in this repo for examples.

## Documentation

https://godoc.org/github.com/ipfs/go-ipns

## Contribute

Feel free to join in. All welcome. Open an [issue](https://github.com/ipfs/go-ipns/issues)!

This repository falls under the IPFS [Code of Conduct](https://github.com/ipfs/community/blob/master/code-of-conduct.md).

### Want to hack on IPFS?

[![](https://cdn.rawgit.com/jbenet/contribute-ipfs-gif/master/img/contribute.gif)](https://github.com/ipfs/community/blob/master/CONTRIBUTING.md)

## License

Copyright (c) Protocol Labs, Inc. under the **MIT license**. See [LICENSE file](./LICENSE) for details.
