# Builder Docker Image

This directory contains a docker image for building chainlink which includes experimental support for SGX.

To enable SGX support in the chainlink docker image, build it like so:

```bash
SGX_ENABLED=yes make docker
```

NOTE: With SGX enabled the HTTP Adapter operates from within an SGX enclave, it currently is a no-op.
