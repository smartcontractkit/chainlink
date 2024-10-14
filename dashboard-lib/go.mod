module github.com/smartcontractkit/chainlink/dashboard-lib

go 1.22.8

require (
	github.com/K-Phoen/grabana v0.22.1
	github.com/grafana/grafana-foundation-sdk/go v0.0.0-20240326122733-6f96a993222b
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.33.0
)

require (
	github.com/K-Phoen/sdk v0.12.4 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/gosimple/slug v1.13.1 // indirect
	github.com/gosimple/unidecode v1.0.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.59.1 // indirect
	golang.org/x/sys v0.25.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)

replace (
	github.com/grafana/grafana-foundation-sdk/go => github.com/grafana/grafana-foundation-sdk/go v0.0.0-20240314112857-a7c9c6d0044c
	github.com/sourcegraph/sourcegraph/lib => github.com/sourcegraph/sourcegraph-public-snapshot/lib v0.0.0-20240822153003-c864f15af264
)
