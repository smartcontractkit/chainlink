# Distributed Tracing

As part of the LOOP plugin effort, we've added distributed tracing to the core node. This is helpful for initial development and maintenance of LOOPs, but will also empower product teams building on top of core. 

## Dev environment

One way to generate traces locally today is with the OCR2 basic smoke test. 

1. navigate to `.github/tracing/` and then run `docker compose --file local-smoke-docker-compose.yaml up`
2. setup a local docker registry at `127.0.0.1:5000` (https://www.docker.com/blog/how-to-use-your-own-registry-2/)
3. run `make build_push_plugin_docker_image` in `chainlink/integration-tests/Makefile`
4. run `SELECTED_NETWORKS=SIMULATED CHAINLINK_IMAGE="127.0.0.1:5000/chainlink" CHAINLINK_VERSION="develop" go test -run TestOCRv2Basic ./smoke/ocr2_test.go`
5. navigate to `localhost:3000/explore` in a web browser to query for traces

Core and the median plugins are instrumented with open telemetry traces, which are sent to the OTEL collector and forwarded to the Tempo backend. The grafana UI can then read the trace data from the Tempo backend.



## CI environment

Another way to generate traces is by enabling traces for PRs. This will instrument traces for `TestOCRv2Basic` in the CI run. 

1. Cut a PR in the core repo
2. Add the `enable tracing` label to the PR
3. Navigate to `Integration Tests / ETH Smoke Tests ocr2-plugins (pull_request)` details
4. Navigate to the summary of the integration tests
5. After the test completes, the generated trace data will be saved as an artifact, currently called `trace-data`
6. Download the artifact to this directory (`chainlink/.github/tracing`)
7. `docker compose --file local-smoke-docker-compose.yaml up`
8. Run `sh replay.sh` to replay those traces to the otel-collector container that was spun up in the last step. 
9. navigate to `localhost:3000/explore` in a web browser to query for traces

The artifact is not json encoded - each individual line is a well formed and complete json object. 

## Configuration
This folder contains the following config files:
* otel-collector-ci.yaml
* otel-collector-dev.yaml
* tempo.yaml
* grafana-datasources.yaml

These config files are for an OTEL collector, grafana Tempo, and a grafana UI instance to run as containers on the same network.
`otel-collector-dev.yaml` is the configuration for dev (i.e. your local machine) environments, and forwards traces from the otel collector to the grafana tempo instance on the same network. 
`otel-collector-ci.yaml` is the configuration for the CI runs, and exports the trace data to the artifact from the github run. 