package gorules

import "github.com/quasilyte/go-ruleguard/dsl"

// NoAdhocGoBuild detects direct "go build" command executions.
func NoAdhocGoBuild(m dsl.Matcher) {
	m.Match(
		`exec.Command("go", "build", $*_)`,
		`exec.CommandContext($_, "go", "build", $*_)`,
	).Report("Do not invoke 'go build' directly for WASM builds. Use github.com/smartcontractkit/chainlink-common/pkg/wasmbuild instead.")
}

// NoWasmTargetEnv detects hardcoded WASM compilation target environment variables.
func NoWasmTargetEnv(m dsl.Matcher) {
	m.Match(
		`$x = append($_, $*_, $s, $*_)`,
		`append($_, $*_, $s, $*_)`,
		`os.Setenv($s, $_)`,
		`os.Setenv($_, $s)`,
	).Where(m["s"].Text.Matches(`.*GOARCH=wasm.*`) || m["s"].Text.Matches(`.*GOOS=wasip1.*`)).
		Report("Do not set GOARCH=wasm or GOOS=wasip1 directly. Use github.com/smartcontractkit/chainlink-common/pkg/wasmbuild instead.")
}
