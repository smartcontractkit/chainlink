module github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/llo_consumer

go 1.25.3

require (
	github.com/shopspring/decimal v1.4.0
	github.com/smartcontractkit/chainlink-common v0.9.6-0.20260113161422-e653fe07e0d5
	github.com/smartcontractkit/cre-sdk-go v0.10.0
	google.golang.org/protobuf v1.36.10
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/bahlo/generic-list-go v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/buger/jsonparser v1.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/invopop/jsonschema v0.13.0 // indirect
	github.com/mailru/easyjson v0.9.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/prometheus/client_golang v1.22.0 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.65.0 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/santhosh-tekuri/jsonschema/v5 v5.3.1 // indirect
	github.com/smartcontractkit/chainlink-protos/cre/go v0.0.0-20251124151448-0448aefdaab9 // indirect
	github.com/smartcontractkit/libocr v0.0.0-20250912173940-f3ab0246e23d // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	github.com/wk8/go-ordered-map/v2 v2.1.8 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250825161204-c5933d9347a5 // indirect
	google.golang.org/grpc v1.75.0 // indirect
)

// Use local chainlink-common when it has ReportFormat on streams.Report. Remove before merge when chainlink-common is published with that change.
replace github.com/smartcontractkit/chainlink-common => ../../../../../../chainlink-common
