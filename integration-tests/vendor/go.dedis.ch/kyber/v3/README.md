[![Go test](https://github.com/dedis/kyber/actions/workflows/go_tests.yml/badge.svg)](https://github.com/dedis/kyber/actions/workflows/go_tests.yml)
[![Coverage Status](https://coveralls.io/repos/github/dedis/kyber/badge.svg?branch=master)](https://coveralls.io/github/dedis/kyber?branch=master)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=dedis_kyber&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=dedis_kyber)
[![Go Reference](https://pkg.go.dev/badge/github.com/dedis/kyber.svg)](https://pkg.go.dev/github.com/dedis/kyber)

DEDIS Advanced Crypto Library for Go
====================================

This package provides a toolbox of advanced cryptographic primitives for Go,
targeting applications like [Cothority](https://go.dedis.ch/cothority)
that need more than straightforward signing and encryption.
Please see the
[Godoc documentation for this package](https://godoc.org/go.dedis.ch/kyber)
for details on the library's purpose and API functionality.

This package includes a mix of variable time and constant time
implementations. If your application is sensitive to timing-based attacks
and you need to constrain Kyber to offering only constant time implementations,
you should use the [suites.RequireConstantTime()](https://godoc.org/go.dedis.ch/kyber/suites#RequireConstantTime)
function in the `init()` function of your `main` package.

Versioning - Development
------------------------

We use the following versioning model:

* crypto.v0 was the first semi-stable version. See [migration notes](https://github.com/dedis/kyber/wiki/Migration-from-gopkg.in-dedis-crypto.v0).
* kyber.v1 never existed, in order to keep kyber, onet and cothorithy versions linked
* gopkg.in/dedis/kyber.v2 was the last stable version
* Starting with v3.0.0, kyber is a Go module, and we respect [semantic versioning](https://golang.org/cmd/go/#hdr-Module_compatibility_and_semantic_versioning).

So if you depend on the master branch, you can expect breakages from time
to time. If you need something that doesn't change in a backward-compatible
way you should use have a `go.mod` file in the directory where your
main package is.

Using the module
----------------

Kyber supports Go modules, and currently has a major version of 3, which means that
the import path is: `go.dedis.ch/kyber/v3`.

Here is a basic example of getting started using it:
1. Make a new directory called “ex". Change directory to “ex" and put this in main.go:
```go
package main

import (
    "fmt"
    "go.dedis.ch/kyber/v3/suites"
)

func main() {
    s := suites.MustFind("Ed25519")
    x := s.Scalar().Zero()
    fmt.Println(x)
}
```
2. Type “go mod init example.com/ex”. The resulting go.mod file will have no dependencies listed yet.
3. Type “go build”. The go tool will fill in the new dependencies that it find for you, i.e. "require go.dedis.ch/kyber/v3 v3.0.13”.
4. Running `./ex` will print `0000000000000000000000000000000000000000000000000000000000000000`.

A note on deriving shared secrets
---------------------------------

Traditionally, ECDH (Elliptic curve Diffie-Hellman) derives the shared secret
from the x point only. In this framework, you can either manually retrieve the
value or use the MarshalBinary method to take the combined (x, y) value as the
shared secret. We recommend the latter process for new softare/protocols using
this framework as it is cleaner and generalizes across different types of groups
(e.g., both integer and elliptic curves), although it will likely be
incompatible with other implementations of ECDH. See [the Wikipedia
page](http://en.wikipedia.org/wiki/Elliptic_curve_Diffie%E2%80%93Hellman) on
ECDH.

Reporting security problems
---------------------------

This library is offered as-is, and without a guarantee. It will need an
independent security review before it should be considered ready for use in
security-critical applications. If you integrate Kyber into your application it
is YOUR RESPONSIBILITY to arrange for that audit.

If you notice a possible security problem, please report it
to dedis-security@epfl.ch.
