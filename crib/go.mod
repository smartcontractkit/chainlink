module github.com/smartcontractkit/chainlink/charts/chainlink-cluster

go 1.21.7

require (
	github.com/K-Phoen/grabana v0.22.1
	github.com/smartcontractkit/chainlink/dashboard-lib v0.0.0-00010101000000-000000000000
	github.com/smartcontractkit/wasp v0.4.6
)

require (
	github.com/K-Phoen/sdk v0.12.4 // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/rs/zerolog v1.32.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
)

replace (
	github.com/go-kit/log => github.com/go-kit/log v0.2.1

	// replicating the replace directive on cosmos SDK
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1
	github.com/grafana/grafana-foundation-sdk/go => github.com/grafana/grafana-foundation-sdk/go v0.0.0-20240314112857-a7c9c6d0044c

	// until merged upstream: https://github.com/hashicorp/go-plugin/pull/257
	github.com/hashicorp/go-plugin => github.com/smartcontractkit/go-plugin v0.0.0-20240208201424-b3b91517de16

	// until merged upstream: https://github.com/mwitkow/grpc-proxy/pull/69
	github.com/mwitkow/grpc-proxy => github.com/smartcontractkit/grpc-proxy v0.0.0-20230731113816-f1be6620749f

	github.com/sercand/kuberesolver/v4 => github.com/sercand/kuberesolver/v5 v5.1.1
	github.com/smartcontractkit/chainlink/dashboard-lib => ./../dashboard-lib
)
