# SGX Development

Working on the SGX build of chainlink is a little bit more involved than the
regular chainlink, especially if you're on Mac OS X. The recommended way is to
do your compilation and testing from within our builder image. This has the
exact set of tools installed to build chainlink with.

To enter the Docker image you call `tools/bin/sgx-env`.

From there you can run `make install` to build the `chainlink` executable with
SGX adapters.

If you don't have an SGX chip in your PC, you should set `SGX_SIMULATION=yes`.

## Layout

The SGX code in chainlink primarily sits within adapters. Adapters are
structures that implement a `Perform` function. This part is implemented in Go.

Each adapter will eventually have an SGX implementation, the SGX version
includes an `sgx` suffix. For example: `http.go` has an SGX version
`http_sgx.go`. The one used is toggled by the `sgx_enabled` build tag.

The untrusted side lives in `sgx/libadapters` and is linked in via LDFLAGS at
the beginning of each `_sgx.go` file.

The enclave code lives in `sgx/enclave`. This code is built into a signed
shared library called `enclave.signed.so` and it must live in the same
directory as the chainlink executable at run time.

## Docker Image

If you want to build the SGX docker image, you can use:

```bash
SGX_ENABLED=yes make docker
```
