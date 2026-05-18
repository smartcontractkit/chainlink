# LOOP Plugins

:warning: Experimental :warning:

This directory supports Local-Out-Of-Process (LOOP) Plugins, an alternative node runtime where some systems execute in 
separate processes, plug-in via [github.com/hashicorp/go-plugin](https://github.com/hashicorp/go-plugin), and 
communicate via [GRPC](https://grpc.io).

There are currently two kinds of plugins: Relayer plugins, and a Median product plugin. The [cmd](cmd) directory contains
some `package main`s while we transition, and they can be built via `make install-<plugin>`. Solana & Starknet has been 
moved to their respective repos, and all must be moved out of this module eventually.

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
via libocr's [LocalConfig.MinOCR2MaxDurationQuery](https://pkg.go.dev/github.com/smartcontractkit/libocr/offchainreporting2plus/types#LocalConfig).
If left unset, the default value is `100ms`.

#### Prometheus


LOOPPs are dynamic, and so must be monitoring. 
We use Plugin discovery to dynamically determine what to monitor based on what plugins are running
and we route external prom scraping to the plugins without exposing them directly

The endpoints are

`/discovery` : HTTP Service Discovery [https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config]
Prometheus server is configured to poll this url to discover new endpoints to monitor. The node serves the response based on what plugins are running,

`/plugins/<name>/metrics`: The node acts as very thin middleware to route from Prometheus server scrape requests to individual plugin /metrics endpoint
Once a plugin is discovered via the discovery mechanism above, the Prometheus service calls the target endpoint at the scrape interval
The node acts as middleware to route the request to the /metrics endpoint of the requested plugin

The simplest change to monitor LOOPPs is to add a service discovery to the scrape configuration
- job_name: 'chainlink_node'
  ...
+  http_sd_configs:
+      - url: "http://127.0.0.1:6688/discovery"
+        refresh_interval: 30s


See the Prometheus documentation for full details [https://prometheus.io/docs/prometheus/latest/configuration/configuration/#http_sd_config]