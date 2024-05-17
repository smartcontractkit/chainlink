module github.com/smartcontractkit/chainlink-common

go 1.21

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.3.0
	github.com/dominikbraun/graph v0.23.0
	github.com/fxamacker/cbor/v2 v2.5.0
	github.com/go-json-experiment/json v0.0.0-20231102232822-2e55bd4e08b0
	github.com/google/go-cmp v0.6.0
	github.com/google/uuid v1.3.1
	github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus v1.0.0
	github.com/hashicorp/consul/sdk v0.16.0
	github.com/hashicorp/go-hclog v1.5.0
	github.com/hashicorp/go-plugin v1.6.0
	github.com/invopop/jsonschema v0.12.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/jonboulle/clockwork v0.4.0
	github.com/jpillora/backoff v1.0.0
	github.com/lib/pq v1.2.0
	github.com/linkedin/goavro/v2 v2.12.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/mwitkow/grpc-proxy v0.0.0-20230212185441-f345521cb9c9
	github.com/pelletier/go-toml/v2 v2.1.1
	github.com/prometheus/client_golang v1.17.0
	github.com/riferrei/srclient v0.5.4
	github.com/shopspring/decimal v1.3.1
	github.com/smartcontractkit/libocr v0.0.0-20240326191951-2bbe9382d052
	github.com/stretchr/testify v1.9.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.45.0
	go.opentelemetry.io/otel v1.19.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.18.0
	go.opentelemetry.io/otel/sdk v1.19.0
	go.opentelemetry.io/otel/trace v1.19.0
	go.uber.org/goleak v1.2.1
	go.uber.org/zap v1.26.0
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa
	golang.org/x/mod v0.14.0
	google.golang.org/grpc v1.58.3
	google.golang.org/protobuf v1.31.0
	sigs.k8s.io/yaml v1.4.0
)

require (
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.14.1 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.5-0.20220116011046-fa5810519dcb // indirect; indirec
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.0.0-rc.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.17 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.2.0
	github.com/stretchr/objx v0.5.2 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.18.0 // indirect
	go.opentelemetry.io/otel/metric v1.19.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.18.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sync v0.6.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230711160842-782d3b101e98 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230711160842-782d3b101e98 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	// until merged upstream: https://github.com/hashicorp/go-plugin/pull/257
	github.com/hashicorp/go-plugin => github.com/smartcontractkit/go-plugin v0.0.0-20240208201424-b3b91517de16

	// until merged upstream: https://github.com/mitchellh/mapstructure/pull/343
	github.com/mitchellh/mapstructure v1.5.0 => github.com/nolag/mapstructure v1.5.1

	// until merged upstream: https://github.com/mwitkow/grpc-proxy/pull/69
	github.com/mwitkow/grpc-proxy => github.com/smartcontractkit/grpc-proxy v0.0.0-20230731113816-f1be6620749f
)
