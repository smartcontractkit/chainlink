package v1_5

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestOnRampView(t *testing.T) {
	e := memory.NewMemoryEnvironment(t, logger.TestLogger(t), zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		Chains: 1,
	})
	chain := e.Chains[e.AllChainSelectors()[0]]
	_, tx, c, err := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
		chain.DeployerKey, chain.Client,
		evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
			LinkToken:          common.HexToAddress("0x1"),
			ChainSelector:      chain.Selector,
			DestChainSelector:  100,
			DefaultTxGasLimit:  10,
			MaxNopFeesJuels:    big.NewInt(10),
			PrevOnRamp:         common.Address{},
			RmnProxy:           common.HexToAddress("0x2"),
			TokenAdminRegistry: common.HexToAddress("0x3"),
		},
		evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                            common.HexToAddress("0x4"),
			MaxNumberOfTokensPerMsg:           0,
			DestGasOverhead:                   0,
			DestGasPerPayloadByte:             0,
			DestDataAvailabilityOverheadGas:   0,
			DestGasPerDataAvailabilityByte:    0,
			DestDataAvailabilityMultiplierBps: 0,
			PriceRegistry:                     common.HexToAddress("0x5"),
			MaxDataBytes:                      0,
			MaxPerMsgGasLimit:                 0,
			DefaultTokenFeeUSDCents:           0,
			DefaultTokenDestGasOverhead:       0,
			EnforceOutOfOrder:                 false,
		},
		evm_2_evm_onramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  big.NewInt(100),
			Rate:      big.NewInt(10),
		},
		[]evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{},
		[]evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{},
		[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
	)
	_, err = deployment.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)
	v, err := GenerateOnRampView(c)
	require.NoError(t, err)
	// Check a few fields.
	assert.Equal(t, v.StaticConfig.ChainSelector, chain.Selector)
	assert.Equal(t, v.DynamicConfig.Router, common.HexToAddress("0x4"))
	assert.Equal(t, "EVM2EVMOnRamp 1.5.0", v.TypeAndVersion)
	_, err = json.MarshalIndent(v, "", "  ")
	require.NoError(t, err)
}
