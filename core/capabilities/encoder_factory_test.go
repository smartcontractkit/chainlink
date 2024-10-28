package capabilities

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"
	"github.com/smartcontractkit/chainlink-common/pkg/values"

	"github.com/smartcontractkit/chainlink/v2/core/logger"

	"github.com/stretchr/testify/require"
)

func Test_NewEncoder(t *testing.T) {
	t.Parallel()
	t.Run("All ocr3 encoder types return a factory", func(t *testing.T) {
		evmEncoding, err := values.NewMap(map[string]any{"abi": "bytes[] Full_reports"})
		require.NoError(t, err)

		config := map[ocr3cap.Encoder]*values.Map{ocr3cap.EncoderEVM: evmEncoding}

		for _, tt := range ocr3cap.Encoders() {
			encoder, err2 := NewEncoder(string(tt), config[tt], logger.NullLogger)
			require.NoError(t, err2)
			require.NotNil(t, encoder)
		}
	})

	t.Run("Invalid encoder returns an error", func(t *testing.T) {
		_, err2 := NewEncoder("NotReal", values.EmptyMap(), logger.NullLogger)
		require.Error(t, err2)
	})
}
