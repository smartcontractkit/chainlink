module github.com/smartcontractkit/chainlink/system-tests/tests/smoke/cre/solana/solwrite

go 1.25.7

require (
	github.com/gagliardetto/binary v0.8.0
	github.com/gagliardetto/solana-go v1.14.0
	github.com/smartcontractkit/chain-selectors v1.0.97
	github.com/smartcontractkit/cre-sdk-go v1.5.0
	github.com/smartcontractkit/cre-sdk-go/capabilities/blockchain/solana v0.1.0
	github.com/smartcontractkit/cre-sdk-go/capabilities/scheduler/cron v1.3.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.1.1 // indirect
	github.com/blendle/zapdriver v1.3.1 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/fatih/color v1.18.0 // indirect
	github.com/gagliardetto/treeout v0.1.4 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.18.2 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/logrusorgru/aurora v2.0.3+incompatible // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/zstdpool-freelist v0.0.0-20201229113212-927304c0c3b1 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/smartcontractkit/chainlink-protos/cre/go v0.0.0-20260226130359-963f935e0396 // indirect
	github.com/streamingfast/logging v0.0.0-20230608130331-f22c91403091 // indirect
	github.com/stretchr/testify v1.11.1 // indirect
	go.mongodb.org/mongo-driver v1.17.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.1 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/term v0.40.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)

replace github.com/mattn/go-isatty => github.com/Unheilbar/go-isatty v0.0.2 // original isatty doesn't support wasim
