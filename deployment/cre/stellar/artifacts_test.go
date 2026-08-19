package stellar

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArtifactServesEveryConstant(t *testing.T) {
	wasmMagic := []byte{0x00, 'a', 's', 'm'}

	for _, artifact := range []string{ForwarderWasm, ReceiverWasm, RejectingReceiverWasm, ReadFixtureWasm} {
		wasm, err := Artifact(artifact)
		require.NoErrorf(t, err, "embedded artifact %s", artifact)
		require.Greaterf(t, len(wasm), len(wasmMagic), "embedded artifact %s is truncated", artifact)
		require.Equalf(t, wasmMagic, wasm[:4], "embedded artifact %s is not a wasm module", artifact)
	}
}

func TestArtifactUnknownName(t *testing.T) {
	_, err := Artifact("bogus.wasm")
	require.Error(t, err)
}
