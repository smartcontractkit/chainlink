package cosmostxm

import (
	"os"
	"testing"

	"github.com/smartcontractkit/chainlink-cosmos/pkg/cosmos/params"
)

func TestMain(m *testing.M) {
	params.InitCosmosSdk(
		/* bech32Prefix= */ "wasm",
		/* token= */ "atom",
	)
	code := m.Run()
	os.Exit(code)
}
