//go:build ignore

package wasmtest

//go:generate go run ./generator/main.go -pkg core/capabilities/compute/test/simple/cmd
//go:generate go run ./generator/main.go -pkg core/capabilities/compute/test/fetch/cmd
//go:generate go run ./generator/main.go -pkg core/services/workflows/cmd/cre/examples/legacy/data_feeds
//go:generate go run ./generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/http_read
//go:generate go run ./generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/simple_cron
//go:generate go run ./generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/simple_cron_with_config
//go:generate go run ./generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/simple_cron_with_secrets
//go:generate go run ./generator/main.go -pkg core/services/workflows/cmd/cre/examples/v2/empty
//go:generate go run ./generator/main.go -pkg core/services/workflows/test/zerotimeout/cmd
//go:generate go run ./generator/main.go -pkg core/services/workflows/test/wasm/legacy/cmd
//go:generate go run ./generator/main.go -pkg core/services/workflows/test/break/cmd
//go:generate go run ./generator/main.go -pkg core/services/workflows/test/wasm/v2/cmd/without_tee
//go:generate go run ./generator/main.go -pkg core/services/workflows/test/wasm/v2/cmd/with_tee
//go:generate go run ./generator/main.go -pkg core/services/workflows/test/wasm/v2/cmd
//go:generate go run ./generator/main.go -pkg core/services/workflows/test/wasm/v2/cmd/with_config
//go:generate go run ./generator/main.go -pkg core/services/workflows/test/wasm/v2/cmd/with_secrets
