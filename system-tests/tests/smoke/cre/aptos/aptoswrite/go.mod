module github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/aptos/aptoswrite

go 1.25.5

require (
	github.com/smartcontractkit/chainlink-protos/cre/go v0.0.0-20260227170625-e0e1c4094174
	github.com/smartcontractkit/cre-sdk-go v1.0.1-0.20251111122439-00032d582c18
	github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/aptos v0.0.0
	github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron v0.10.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace (
	github.com/smartcontractkit/cre-sdk-go => /Users/yashvardhan/cre-sdk-go
	github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/aptos => /Users/yashvardhan/cre-sdk-go/capabilities/blockchain/aptos
)
