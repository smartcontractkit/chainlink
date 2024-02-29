module github.com/smartcontractkit/chainlink/charts/chainlink-cluster/dashboard

go 1.21

require (
	github.com/K-Phoen/grabana v0.22.1
	github.com/smartcontractkit/wasp v0.4.6
)

require (
	github.com/K-Phoen/sdk v0.12.4 // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/prometheus/common v0.45.0 // indirect
)

replace (
	github.com/go-kit/log => github.com/go-kit/log v0.2.1

	// replicating the replace directive on cosmos SDK
	github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

	// until merged upstream: https://github.com/hashicorp/go-plugin/pull/257
	github.com/hashicorp/go-plugin => github.com/smartcontractkit/go-plugin v0.0.0-20240208201424-b3b91517de16

	// until merged upstream: https://github.com/mwitkow/grpc-proxy/pull/69
	github.com/mwitkow/grpc-proxy => github.com/smartcontractkit/grpc-proxy v0.0.0-20230731113816-f1be6620749f

	github.com/sercand/kuberesolver/v4 => github.com/sercand/kuberesolver/v5 v5.1.1
)
