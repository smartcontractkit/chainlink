package evm_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/testutils"
	evmrelaytypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func Test_ChainWriterConfig_MarshalUnmarshal(t *testing.T) {
	raw := evm.ChainWriterConfigRaw(
		testutils.NewAddress(),
		assets.NewWeiI(1),
		500_000, // commit gas limit
		500_000, // exec gas limit
	)
	encoded, err := json.Marshal(raw)
	require.NoError(t, err)
	var decoded evmrelaytypes.ChainWriterConfig
	err = json.Unmarshal(encoded, &decoded)
	require.NoError(t, err)
}
