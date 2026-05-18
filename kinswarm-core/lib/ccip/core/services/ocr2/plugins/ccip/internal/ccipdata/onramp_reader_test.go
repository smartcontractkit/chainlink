package ccipdata_test

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccip"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	evmclientmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_1_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
)

type onRampReaderTH struct {
	user   *bind.TransactOpts
	reader ccipdata.OnRampReader
}

func TestNewOnRampReader_noContractAtAddress(t *testing.T) {
	_, bc := ccipdata.NewSimulation(t)
	addr := ccipcalc.EvmAddrToGeneric(utils.RandomAddress())
	_, err := factory.NewOnRampReader(logger.TestLogger(t), factory.NewEvmVersionFinder(), testutils.SimulatedChainID.Uint64(), testutils.SimulatedChainID.Uint64(), addr, lpmocks.NewLogPoller(t), bc)
	assert.EqualError(t, err, fmt.Sprintf("unable to read type and version: error calling typeAndVersion on addr: %s no contract code at given address", addr))
}

func TestOnRampReaderInit(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "OnRampReader_V1_0_0",
			version: ccipdata.V1_0_0,
		},
		{
			name:    "OnRampReader_V1_1_0",
			version: ccipdata.V1_1_0,
		},
		{
			name:    "OnRampReader_V1_2_0",
			version: ccipdata.V1_2_0,
		},
		{
			name:    "OnRampReader_V1_5_0",
			version: ccipdata.V1_5_0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			th := setupOnRampReaderTH(t, test.version)
			testVersionSpecificOnRampReader(t, th, test.version)
		})
	}
}

func setupOnRampReaderTH(t *testing.T, version string) onRampReaderTH {
	user, bc := ccipdata.NewSimulation(t)
	log := logger.TestLogger(t)
	orm := logpoller.NewORM(testutils.SimulatedChainID, pgtest.NewSqlxDB(t), log)
	lpOpts := logpoller.Opts{
		PollPeriod:               100 * time.Millisecond,
		FinalityDepth:            2,
		BackfillBatchSize:        3,
		RpcBatchSize:             2,
		KeepFinalizedBlocksDepth: 1000,
	}
	headTracker := headtracker.NewSimulatedHeadTracker(bc, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	if lpOpts.PollPeriod == 0 {
		lpOpts.PollPeriod = 1 * time.Hour
	}
	lp := logpoller.NewLogPoller(
		orm,
		bc,
		log,
		headTracker,
		lpOpts)

	// Setup onRamp.
	var onRampAddress common.Address
	switch version {
	case ccipdata.V1_0_0:
		onRampAddress = setupOnRampV1_0_0(t, user, bc)
	case ccipdata.V1_1_0:
		onRampAddress = setupOnRampV1_1_0(t, user, bc)
	case ccipdata.V1_2_0:
		onRampAddress = setupOnRampV1_2_0(t, user, bc)
	case ccipdata.V1_5_0:
		onRampAddress = setupOnRampV1_5_0(t, user, bc)
	default:
		require.Fail(t, "Unknown version: ", version)
	}

	// Create the version-specific reader.
	reader, err := factory.NewOnRampReader(log, factory.NewEvmVersionFinder(), testutils.SimulatedChainID.Uint64(), testutils.SimulatedChainID.Uint64(), ccipcalc.EvmAddrToGeneric(onRampAddress), lp, bc)
	require.NoError(t, err)

	return onRampReaderTH{
		user:   user,
		reader: reader,
	}
}

func setupOnRampV1_0_0(t *testing.T, user *bind.TransactOpts, bc *client.SimulatedBackendClient) common.Address {
	linkTokenAddress := common.HexToAddress("0x000011")
	staticConfig := evm_2_evm_onramp_1_0_0.EVM2EVMOnRampStaticConfig{
		LinkToken:         linkTokenAddress,
		ChainSelector:     testutils.SimulatedChainID.Uint64(),
		DestChainSelector: testutils.SimulatedChainID.Uint64(),
		DefaultTxGasLimit: 30000,
		MaxNopFeesJuels:   big.NewInt(1000000),
		PrevOnRamp:        common.Address{},
		ArmProxy:          utils.RandomAddress(),
	}
	dynamicConfig := evm_2_evm_onramp_1_0_0.EVM2EVMOnRampDynamicConfig{
		Router:          common.HexToAddress("0x000100"),
		MaxTokensLength: 4,
		PriceRegistry:   utils.RandomAddress(),
		MaxDataSize:     100000,
		MaxGasLimit:     100000,
	}
	rateLimiterConfig := evm_2_evm_onramp_1_0_0.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(5),
		Rate:      big.NewInt(5),
	}
	allowList := []common.Address{user.From}
	feeTokenConfigs := []evm_2_evm_onramp_1_0_0.EVM2EVMOnRampFeeTokenConfigArgs{
		{
			Token:                 linkTokenAddress,
			GasMultiplier:         1,
			NetworkFeeAmountUSD:   big.NewInt(0),
			DestGasOverhead:       50,
			DestGasPerPayloadByte: 60,
			Enabled:               false,
		},
	}
	tokenTransferConfigArgs := []evm_2_evm_onramp_1_0_0.EVM2EVMOnRampTokenTransferFeeConfigArgs{
		{
			Token:  utils.RandomAddress(),
			MinFee: 10,
			MaxFee: 1000,
			Ratio:  1,
		},
	}
	nopsAndWeights := []evm_2_evm_onramp_1_0_0.EVM2EVMOnRampNopAndWeight{
		{
			Nop:    utils.RandomAddress(),
			Weight: 1,
		},
	}
	tokenAndPool := []evm_2_evm_onramp_1_0_0.InternalPoolUpdate{}
	onRampAddress, transaction, _, err := evm_2_evm_onramp_1_0_0.DeployEVM2EVMOnRamp(
		user,
		bc,
		staticConfig,
		dynamicConfig,
		tokenAndPool,
		allowList,
		rateLimiterConfig,
		feeTokenConfigs,
		tokenTransferConfigArgs,
		nopsAndWeights,
	)
	bc.Commit()
	require.NoError(t, err)
	ccipdata.AssertNonRevert(t, transaction, bc, user)
	return onRampAddress
}

func setupOnRampV1_1_0(t *testing.T, user *bind.TransactOpts, bc *client.SimulatedBackendClient) common.Address {
	linkTokenAddress := common.HexToAddress("0x000011")
	staticConfig := evm_2_evm_onramp_1_1_0.EVM2EVMOnRampStaticConfig{
		LinkToken:         linkTokenAddress,
		ChainSelector:     testutils.SimulatedChainID.Uint64(),
		DestChainSelector: testutils.SimulatedChainID.Uint64(),
		DefaultTxGasLimit: 30000,
		MaxNopFeesJuels:   big.NewInt(1000000),
		PrevOnRamp:        common.Address{},
		ArmProxy:          utils.RandomAddress(),
	}
	dynamicConfig := evm_2_evm_onramp_1_1_0.EVM2EVMOnRampDynamicConfig{
		Router:          common.HexToAddress("0x000110"),
		MaxTokensLength: 4,
		PriceRegistry:   common.HexToAddress("0x000066"),
		MaxDataSize:     100000,
		MaxGasLimit:     100000,
	}
	rateLimiterConfig := evm_2_evm_onramp_1_1_0.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(5),
		Rate:      big.NewInt(5),
	}
	allowList := []common.Address{user.From}
	feeTokenConfigs := []evm_2_evm_onramp_1_1_0.EVM2EVMOnRampFeeTokenConfigArgs{
		{
			Token:                  linkTokenAddress,
			NetworkFeeUSD:          0,
			MinTokenTransferFeeUSD: 0,
			MaxTokenTransferFeeUSD: 0,
			GasMultiplier:          0,
			PremiumMultiplier:      0,
			Enabled:                false,
		},
	}
	tokenTransferConfigArgs := []evm_2_evm_onramp_1_1_0.EVM2EVMOnRampTokenTransferFeeConfigArgs{
		{
			Token:           linkTokenAddress,
			Ratio:           0,
			DestGasOverhead: 0,
		},
	}
	nopsAndWeights := []evm_2_evm_onramp_1_1_0.EVM2EVMOnRampNopAndWeight{
		{
			Nop:    common.HexToAddress("0x222222222"),
			Weight: 1,
		},
	}
	tokenAndPool := []evm_2_evm_onramp_1_1_0.InternalPoolUpdate{}
	onRampAddress, transaction, _, err := evm_2_evm_onramp_1_1_0.DeployEVM2EVMOnRamp(
		user,
		bc,
		staticConfig,
		dynamicConfig,
		tokenAndPool,
		allowList,
		rateLimiterConfig,
		feeTokenConfigs,
		tokenTransferConfigArgs,
		nopsAndWeights,
	)
	bc.Commit()
	require.NoError(t, err)
	ccipdata.AssertNonRevert(t, transaction, bc, user)
	return onRampAddress
}

func setupOnRampV1_2_0(t *testing.T, user *bind.TransactOpts, bc *client.SimulatedBackendClient) common.Address {
	linkTokenAddress := common.HexToAddress("0x000011")
	staticConfig := evm_2_evm_onramp_1_2_0.EVM2EVMOnRampStaticConfig{
		LinkToken:         linkTokenAddress,
		ChainSelector:     testutils.SimulatedChainID.Uint64(),
		DestChainSelector: testutils.SimulatedChainID.Uint64(),
		DefaultTxGasLimit: 30000,
		MaxNopFeesJuels:   big.NewInt(1000000),
		PrevOnRamp:        common.Address{},
		ArmProxy:          utils.RandomAddress(),
	}
	dynamicConfig := evm_2_evm_onramp_1_2_0.EVM2EVMOnRampDynamicConfig{
		Router:                            common.HexToAddress("0x0000000000000000000000000000000000000120"),
		MaxNumberOfTokensPerMsg:           0,
		DestGasOverhead:                   0,
		DestGasPerPayloadByte:             0,
		DestDataAvailabilityOverheadGas:   0,
		DestGasPerDataAvailabilityByte:    0,
		DestDataAvailabilityMultiplierBps: 0,
		PriceRegistry:                     utils.RandomAddress(),
		MaxDataBytes:                      0,
		MaxPerMsgGasLimit:                 0,
	}
	rateLimiterConfig := evm_2_evm_onramp_1_2_0.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(5),
		Rate:      big.NewInt(5),
	}
	feeTokenConfigs := []evm_2_evm_onramp_1_2_0.EVM2EVMOnRampFeeTokenConfigArgs{
		{
			Token:                      linkTokenAddress,
			NetworkFeeUSDCents:         0,
			GasMultiplierWeiPerEth:     0,
			PremiumMultiplierWeiPerEth: 0,
			Enabled:                    false,
		},
	}
	tokenTransferConfigArgs := []evm_2_evm_onramp_1_2_0.EVM2EVMOnRampTokenTransferFeeConfigArgs{
		{
			Token:             linkTokenAddress,
			MinFeeUSDCents:    0,
			MaxFeeUSDCents:    0,
			DeciBps:           0,
			DestGasOverhead:   0,
			DestBytesOverhead: 0,
		},
	}
	nopsAndWeights := []evm_2_evm_onramp_1_2_0.EVM2EVMOnRampNopAndWeight{
		{
			Nop:    utils.RandomAddress(),
			Weight: 1,
		},
	}
	tokenAndPool := []evm_2_evm_onramp_1_2_0.InternalPoolUpdate{}
	onRampAddress, transaction, _, err := evm_2_evm_onramp_1_2_0.DeployEVM2EVMOnRamp(
		user,
		bc,
		staticConfig,
		dynamicConfig,
		tokenAndPool,
		rateLimiterConfig,
		feeTokenConfigs,
		tokenTransferConfigArgs,
		nopsAndWeights,
	)
	bc.Commit()
	require.NoError(t, err)
	ccipdata.AssertNonRevert(t, transaction, bc, user)
	return onRampAddress
}

func setupOnRampV1_5_0(t *testing.T, user *bind.TransactOpts, bc *client.SimulatedBackendClient) common.Address {
	linkTokenAddress := common.HexToAddress("0x000011")
	staticConfig := evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
		LinkToken:          linkTokenAddress,
		ChainSelector:      testutils.SimulatedChainID.Uint64(),
		DestChainSelector:  testutils.SimulatedChainID.Uint64(),
		DefaultTxGasLimit:  30000,
		MaxNopFeesJuels:    big.NewInt(1000000),
		PrevOnRamp:         common.Address{},
		RmnProxy:           utils.RandomAddress(),
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
			DestBytesOverhead:         64,
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
	return onRampAddress
}

func testVersionSpecificOnRampReader(t *testing.T, th onRampReaderTH, version string) {
	switch version {
	case ccipdata.V1_0_0:
		testOnRampReader(t, th, common.HexToAddress("0x0000000000000000000000000000000000000100"))
	case ccipdata.V1_1_0:
		testOnRampReader(t, th, common.HexToAddress("0x0000000000000000000000000000000000000110"))
	case ccipdata.V1_2_0:
		testOnRampReader(t, th, common.HexToAddress("0x0000000000000000000000000000000000000120"))
	case ccipdata.V1_5_0:
		testOnRampReader(t, th, common.HexToAddress("0x0000000000000000000000000000000000000150"))
	default:
		require.Fail(t, "Unknown version: ", version)
	}
}

func testOnRampReader(t *testing.T, th onRampReaderTH, expectedRouterAddress common.Address) {
	ctx := th.user.Context
	res, err := th.reader.RouterAddress(ctx)
	require.NoError(t, err)
	require.Equal(t, ccipcalc.EvmAddrToGeneric(expectedRouterAddress), res)

	msg, err := th.reader.GetSendRequestsBetweenSeqNums(ctx, 0, 10, true)
	require.NoError(t, err)
	require.NotNil(t, msg)
	require.Equal(t, []cciptypes.EVM2EVMMessageWithTxMeta{}, msg)

	address, err := th.reader.Address(ctx)
	require.NoError(t, err)
	require.NotNil(t, address)

	cfg, err := th.reader.GetDynamicConfig(ctx)
	require.NoError(t, err)
	require.NotNil(t, cfg)
	require.Equal(t, ccipcalc.EvmAddrToGeneric(expectedRouterAddress), cfg.Router)
}

func TestNewOnRampReader(t *testing.T) {
	var tt = []struct {
		typeAndVersion string
		expectedErr    string
	}{
		{
			typeAndVersion: "blah",
			expectedErr:    "unable to read type and version: invalid type and version blah",
		},
		{
			typeAndVersion: "EVM2EVMOffRamp 1.0.0",
			expectedErr:    "expected EVM2EVMOnRamp got EVM2EVMOffRamp",
		},
		{
			typeAndVersion: "EVM2EVMOnRamp 1.2.0",
			expectedErr:    "",
		},
		{
			typeAndVersion: "EVM2EVMOnRamp 2.0.0",
			expectedErr:    "unsupported onramp version 2.0.0",
		},
	}
	for _, tc := range tt {
		t.Run(tc.typeAndVersion, func(t *testing.T) {
			b, err := utils.ABIEncode(`[{"type":"string"}]`, tc.typeAndVersion)
			require.NoError(t, err)
			c := evmclientmocks.NewClient(t)
			c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(b, nil)
			addr := ccipcalc.EvmAddrToGeneric(utils.RandomAddress())
			lp := lpmocks.NewLogPoller(t)
			lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil).Maybe()
			_, err = factory.NewOnRampReader(logger.TestLogger(t), factory.NewEvmVersionFinder(), 1, 2, addr, lp, c)
			if tc.expectedErr != "" {
				require.EqualError(t, err, tc.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
