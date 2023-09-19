# LOOP Plugins

:warning: Experimental :warning:

This directory supports Local-Out-Of-Process (LOOP) Plugins, an alternative node runtime where some systems execute in 
separate processes, plug-in via [github.com/hashicorp/go-plugin](https://github.com/hashicorp/go-plugin), and 
communicate via [GRPC](https://grpc.io).

There are currently two kinds of plugins, and one implementation of each: a Solana Relayer plugin, and a Median Reporting
plugin. The [cmd](cmd) directory contains their `package main`s for now. These can be built via `make install-solana` and 
`make install-median`.

## How to use

[chainlink.Dockerfile](chainlink.Dockerfile) extends the regular [core/chainlink.Dockerfile](../core/chainlink.Dockerfile)
to include the plugin binaries, and enables support by setting `CL_SOLANA_CMD`, `CL_STARKNET_CMD`, and `CL_MEDIAN_CMD`. 
Either plugin can be disabled by un-setting the environment variable, which will revert to the original in-process runtime. 
Images built from this Dockerfile can otherwise be used normally, provided that the [pre-requisites](#pre-requisites) have been met.

### Pre-requisites

#### Timeouts

LOOPPs communicate over GRPC, which always includes a `context.Context` and requires realistic timeouts. Placeholder/dummy
values (e.g. `MaxDurationQuery = 0`) will not work and must be updated to realistic values. In lieu of reconfiguring already
deployed contracts on Solana, the environment variable `CL_MIN_OCR2_MAX_DURATION_QUERY` can be set establish a new minimum
via libocr's [LocalConfig.MinOCR2MaxDureationQuery](https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/types#LocalConfig).
If left unset, the default value is `100ms`.

#### Prometheus

TODO how to preserve metrics https://smartcontract-it.atlassian.net/browse/BCF-2202
