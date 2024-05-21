# Distributed Tracing

As part of the LOOP plugin effort, we've added distributed tracing to the core node. This is helpful for initial development and maintenance of LOOPs, but will also empower product teams building on top of core. 

## Dev environment

One way to generate traces locally today is with the OCR2 basic smoke test. 

1. navigate to `.github/tracing/` and then run `docker compose --file local-smoke-docker-compose.yaml up`
2. setup a local docker registry at `127.0.0.1:5000` (https://www.docker.com/blog/how-to-use-your-own-registry-2/)
3. run `make build_push_plugin_docker_image` in `chainlink/integration-tests/Makefile`
4. preapre your `overrides.toml` file with selected network and CL image name and version and place it anywhere
inside `integration-tests` directory. Sample `overrides.toml` file:
```toml
[ChainlinkImage]
image="127.0.0.1:5000/chainlink"
version="develop"

[Network]
selected_networks=["simulated"]
```
5. run `go test -run TestOCRv2Basic ./smoke/ocr2_test.go`
6. navigate to `localhost:3000/explore` in a web browser to query for traces

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

## Adding Traces to Plugins and to core

Adding traces requires identifying an observability gap in a related group of code executions or a critical path in your application. This is intuitive for the developer:

- "What's the flow of component interaction in this distributed system?"
- "What's the behavior of the JobProcessorOne component when jobs with [x, y, z] attributes are processed?"
- "Is this critical path workflow behaving the way we expect?"

The developer will measure a flow of execution from end to end in one trace. Each logically separate measure of this flow is called a span. Spans have either one or no parent span and multiple children span. The relationship between parent and child spans in agreggate will form a directed acyclic graph. The trace begins at the root of this graph.

The most trivial application of a span is measuring top level performance in one critical path. There is much more you can do, including creating human readable and timestamped events within a span (useful for monitoring concurrent access to resources), recording errors, linking parent and children spans through large parts of an application, and even extending a span beyond a single process.

Spans are created by `tracers` and passed through go applications by `Context`s. A tracer must be initialized first. Both core and plugin developers will initialize a tracer from the globally registered trace provider:

```
tracer := otel.GetTracerProvider().Tracer("example.com/foo")
```

The globally registered tracer provider is available for plugins after they are initialized, and available in core after configuration is processed (`initGlobals`).

Add spans by:
```
  func interestingFunc() {
    // Assuming there is an appropriate parentContext
	ctx, span := tracer.Start(parentContext, "hello-span")
	defer span.End()

	// do some work to track with hello-span
  }
```
As implied by the example, `span` is a child of its parent span captured by `parentContext`.


Note that in certain situations, there are 3rd party libraries that will setup spans. For instance:

```
import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

router := gin.Default()
router.Use(otelgin.Middleware("service-name"))
```

The developer aligns with best practices when they:
- Start with critical paths
- Measure paths from end to end (Context is wired all the way through)
- Emphasize broadness of measurement over depth
- Use automatic instrumentation if possible