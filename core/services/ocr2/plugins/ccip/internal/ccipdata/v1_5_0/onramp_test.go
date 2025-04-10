package v1_5_0

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_5_0/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
)

func TestLogPollerClient_GetSendRequestsBetweenSeqNums1_4_0(t *testing.T) {
	onRampAddr := utils.RandomAddress()
	seqNum := uint64(100)
	limit := uint64(10)
	lggr := logger.TestLogger(t)

	tests := []struct {
		name          string
		finalized     bool
		confirmations evmtypes.Confirmations
	}{
		{"finalized", true, evmtypes.Finalized},
		{"unfinalized", false, evmtypes.Confirmations(0)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lp := mocks.NewLogPoller(t)
			onRampV2, err := NewOnRamp(lggr, 1, 1, onRampAddr, lp, nil)
			require.NoError(t, err)

			lp.On("LogsDataWordRange",
				mock.Anything,
				onRampV2.sendRequestedEventSig,
				onRampAddr,
				onRampV2.sendRequestedSeqNumberWord,
				abihelpers.EvmWord(seqNum),
				abihelpers.EvmWord(seqNum+limit),
				tt.confirmations,
			).Once().Return([]logpoller.Log{}, nil)

			events, err1 := onRampV2.GetSendRequestsBetweenSeqNums(context.Background(), seqNum, seqNum+limit, tt.finalized)
			assert.NoError(t, err1)
			assert.Empty(t, events)

			lp.AssertExpectations(t)
		})
	}
}

func Test_ProperlyRecognizesPerLaneCurses(t *testing.T) {
	user, bc := ccipdata.NewSimulation(t)
	ctx := testutils.Context(t)
	destChainSelector := uint64(100)
	sourceChainSelector := uint64(200)
	onRampAddress, mockRMN, mockRMNAddress := setupOnRampV1_5_0(t, user, bc)

	onRamp, err := NewOnRamp(logger.TestLogger(t), 1, destChainSelector, onRampAddress, mocks.NewLogPoller(t), bc)
	require.NoError(t, err)

	onRamp.cachedStaticConfig = func(ctx context.Context) (evm_2_evm_onramp.EVM2EVMOnRampStaticConfig, error) {
		return evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
			RmnProxy: mockRMNAddress,
		}, nil
	}

	// Lane is not cursed right after deployment
	isCursed, err := onRamp.IsSourceCursed(ctx)
	require.NoError(t, err)
	assert.False(t, isCursed)

	// Cursing different chain selector
	_, err = mockRMN.VoteToCurse0(user, [32]byte{}, ccipcommon.SelectorToBytes(sourceChainSelector))
	require.NoError(t, err)
	bc.Commit()

	isCursed, err = onRamp.IsSourceCursed(ctx)
	require.NoError(t, err)
	assert.False(t, isCursed)

	// Cursing the correct chain selector
	_, err = mockRMN.VoteToCurse0(user, [32]byte{}, ccipcommon.SelectorToBytes(destChainSelector))
	require.NoError(t, err)
	bc.Commit()

	isCursed, err = onRamp.IsSourceCursed(ctx)
	require.NoError(t, err)
	assert.True(t, isCursed)

	// Uncursing the chain selector
	_, err = mockRMN.OwnerUnvoteToCurse(user, []mock_rmn_contract.RMNUnvoteToCurseRecord{}, ccipcommon.SelectorToBytes(destChainSelector))
	require.NoError(t, err)
	bc.Commit()

	isCursed, err = onRamp.IsSourceCursed(ctx)
	require.NoError(t, err)
	assert.False(t, isCursed)
}

// This is written to benchmark before and after the caching of StaticConfig and RMNContract
func BenchmarkIsSourceCursedWithCache(b *testing.B) {
	user, bc := ccipdata.NewSimulation(b)
	ctx := testutils.Context(b)
	destChainSelector := uint64(100)
	onRampAddress, _, _ := setupOnRampV1_5_0(b, user, bc)

	onRamp, err := NewOnRamp(logger.TestLogger(b), 1, destChainSelector, onRampAddress, mocks.NewLogPoller(b), bc)
	require.NoError(b, err)

	for i := 0; i < b.N; i++ {
		_, _ = onRamp.IsSourceCursed(ctx)
	}
}

func setupOnRampV1_5_0(t testing.TB, user *bind.TransactOpts, bc *client.SimulatedBackendClient) (common.Address, *mock_rmn_contract.MockRMNContract, common.Address) {
	rmnAddress, transaction, rmnContract, err := mock_rmn_contract.DeployMockRMNContract(user, bc)
	bc.Commit()
	require.NoError(t, err)
	ccipdata.AssertNonRevert(t, transaction, bc, user)

	linkTokenAddress := common.HexToAddress("0x000011")
	staticConfig := evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
		LinkToken:          linkTokenAddress,
		ChainSelector:      testutils.SimulatedChainID.Uint64(),
		DestChainSelector:  testutils.SimulatedChainID.Uint64(),
		DefaultTxGasLimit:  30000,
		MaxNopFeesJuels:    big.NewInt(1000000),
		PrevOnRamp:         common.Address{},
		RmnProxy:           rmnAddress,
		TokenAdminRegistry: utils.RandomAddress(),
	}
	dynamicConfig := evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
		Router:                            common.HexToAddress("0x0000000000000000000000000000000000000150"),
		MaxNumberOfTokensPerMsg:           0,
		DestGasOverhead:                   0,
		DestGasPerPayloadByte:             0,
		DestDataAvailabilityOverheadGas:   0,
		DestGasPerDataAvailabilityByte:    0,
		DestDataAvailabilityMultiplierBps: 0,
		PriceRegistry:                     utils.RandomAddress(),
		MaxDataBytes:                      0,
		MaxPerMsgGasLimit:                 0,
		DefaultTokenFeeUSDCents:           50,
		DefaultTokenDestGasOverhead:       125_000,
	}
	rateLimiterConfig := evm_2_evm_onramp.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(5),
		Rate:      big.NewInt(5),
	}
	feeTokenConfigs := []evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
		{
			Token:                      linkTokenAddress,
			NetworkFeeUSDCents:         0,
			GasMultiplierWeiPerEth:     0,
			PremiumMultiplierWeiPerEth: 0,
			Enabled:                    false,
		},
	}
	tokenTransferConfigArgs := []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
		{
			Token:                     linkTokenAddress,
			MinFeeUSDCents:            0,
			MaxFeeUSDCents:            0,
			DeciBps:                   0,
			DestGasOverhead:           0,
			DestBytesOverhead:         32,
			AggregateRateLimitEnabled: true,
		},
	}
	nopsAndWeights := []evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{
		{
			Nop:    utils.RandomAddress(),
			Weight: 1,
		},
	}
	onRampAddress, transaction, _, err := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
		user,
		bc,
		staticConfig,
		dynamicConfig,
		rateLimiterConfig,
		feeTokenConfigs,
		tokenTransferConfigArgs,
		nopsAndWeights,
	)
	bc.Commit()
	require.NoError(t, err)
	ccipdata.AssertNonRevert(t, transaction, bc, user)

	return onRampAddress, rmnContract, rmnAddress
}
