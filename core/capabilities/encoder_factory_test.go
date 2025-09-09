package capabilities

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/capabilities/consensus/ocr3/ocr3cap"
	"github.com/smartcontractkit/chainlink-protos/cre/go/values"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/stretchr/testify/require"
)

func Test_NewEncoder(t *testing.T) {
	t.Parallel()
	t.Run("All ocr3 encoder types return a factory", func(t *testing.T) {
		evmEncoding, err := values.NewMap(map[string]any{"abi": "bytes[] Full_reports"})
		require.NoError(t, err)
		solanaEncoding := getSolCfg(t)

		config := map[ocr3cap.Encoder]*values.Map{
			ocr3cap.EncoderEVM:   evmEncoding,
			ocr3cap.EncoderBorsh: solanaEncoding,
		}

		for _, tt := range ocr3cap.Encoders() {
			encoder, err2 := NewEncoder(string(tt), config[tt], logger.Nop())
			require.NoError(t, err2)
			require.NotNil(t, encoder)
		}
	})

	t.Run("Invalid encoder returns an error", func(t *testing.T) {
		_, err2 := NewEncoder("NotReal", values.EmptyMap(), logger.Nop())
		require.Error(t, err2)
	})
}

func getSolCfg(t *testing.T) *values.Map {
	cfg := map[string]any{
		"report_schema": `{
      "kind": "struct",
      "fields": [
        { "name": "payload", "type": { "vec": { "defined": "DecimalReport" } } }
      ]
    }`,
		"defined_types": `
		[
      {
        "name":"DecimalReport",
         "type":{
          "kind":"struct",
          "fields":[
            { "name":"timestamp", "type":"u32" },
            { "name":"answer",    "type":"u128" },
            { "name": "dataId",   "type": {"array": ["u8",16]}}
          ]
        }
      }
    ]`,
	}

	mcfg, err := values.NewMap(cfg)
	require.NoError(t, err)
	return mcfg
}
