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


## Production and NOPs environments

In a production environment, we suggest coupling the lifecycle of nodes and otel-collectors. A best practice is to deploy the otel-collector alongside your node, using infrastructure as code (IAC) to automate deployments and certificate lifecycles. While there are valid use cases for using `Tracing.Mode = unencrypted`, we have set the default encryption setting to `Tracing.Mode = tls`. Externally deployed otel-collectors can not be used with `Tracing.Mode = unencrypted`. i.e. If `Tracing.Mode = unencrypted` and an external URI is detected for `Tracing.CollectorTarget` node configuration will fail to validate and the node will not boot. The node requires a valid encryption mode and collector target to send traces.

Once traces reach the otel-collector, the rest of the observability pipeline is flexible. We recommend deploying (through automation) centrally managed Grafana Tempo and Grafana UI instances to receive from one or many otel-collector instances. Always use networking best practices and encrypt trace data, especially at network boundaries.

## Configuration
This folder contains the following config files:
* otel-collector-ci.yaml
* otel-collector-dev.yaml
* tempo.yaml
* grafana-datasources.yaml

These config files are for an OTEL collector, grafana Tempo, and a grafana UI instance to run as containers on the same network.
`otel-collector-dev.yaml` is the configuration for dev (i.e. your local machine) environments, and forwards traces from the otel collector to the grafana tempo instance on the same network. 
`otel-collector-ci.yaml` is the configuration for the CI runs, and exports the trace data to the artifact from the github run. 