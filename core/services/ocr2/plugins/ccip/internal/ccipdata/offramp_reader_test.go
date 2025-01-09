package ccipdata_test

import (
	"math/big"
	"math/rand"
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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_helper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_rmn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipcalc"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/factory"
	ccipdatamocks "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/v1_2_0"
)

type offRampReaderTH struct {
	user   *bind.TransactOpts
	reader ccipdata.OffRampReader
}

func TestExecOnchainConfig100(t *testing.T) {
	tests := []struct {
		name      string
		want      v1_0_0.ExecOnchainConfig
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: v1_0_0.ExecOnchainConfig{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				Router:                                  utils.RandomAddress(),
				PriceRegistry:                           utils.RandomAddress(),
				MaxTokensLength:                         uint16(rand.Uint32()),
				MaxDataSize:                             rand.Uint32(),
			},
		},
		{
			name: "encodes and fails decoding config with missing fields",
			want: v1_0_0.ExecOnchainConfig{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				MaxDataSize:                             rand.Uint32(),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := abihelpers.EncodeAbiStruct(tt.want)
			require.NoError(t, err)

			decoded, err := abihelpers.DecodeAbiStruct[v1_0_0.ExecOnchainConfig](encoded)
			if tt.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, decoded)
			}
		})
	}
}

func TestExecOnchainConfig120(t *testing.T) {
	tests := []struct {
		name      string
		want      v1_2_0.ExecOnchainConfig
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: v1_2_0.ExecOnchainConfig{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				Router:                                  utils.RandomAddress(),
				PriceRegistry:                           utils.RandomAddress(),
				MaxNumberOfTokensPerMsg:                 uint16(rand.Uint32()),
				MaxDataBytes:                            rand.Uint32(),
				MaxPoolReleaseOrMintGas:                 rand.Uint32(),
			},
		},
		{
			name: "encodes and fails decoding config with missing fields",
			want: v1_2_0.ExecOnchainConfig{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				MaxDataBytes:                            rand.Uint32(),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := abihelpers.EncodeAbiStruct(tt.want)
			require.NoError(t, err)

			decoded, err := abihelpers.DecodeAbiStruct[v1_2_0.ExecOnchainConfig](encoded)
			if tt.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, decoded)
			}
		})
	}
}

func TestOffRampReaderInit(t *testing.T) {
	tests := []struct {
		name    string
		version string
	}{
		{
			name:    "OffRampReader_V1_0_0",
			version: ccipdata.V1_0_0,
		},
		{
			name:    "OffRampReader_V1_1_0",
			version: ccipdata.V1_1_0,
		},
		{
			name:    "OffRampReader_V1_2_0",
			version: ccipdata.V1_2_0,
		},
		{
			name:    "OffRampReader_V1_5_0",
			version: ccipdata.V1_5_0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			th := setupOffRampReaderTH(t, test.version)
			testOffRampReader(t, th)
		})
	}
}

func setupOffRampReaderTH(t *testing.T, version string) offRampReaderTH {
	ctx := testutils.Context(t)
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
	assert.NoError(t, orm.InsertBlock(ctx, common.Hash{}, 1, time.Now(), 1))
	// Setup offRamp.
	var offRampAddress common.Address
	switch version {
	case ccipdata.V1_0_0:
		offRampAddress = setupOffRampV1_0_0(t, user, bc)
	case ccipdata.V1_1_0:
		// Version 1.1.0 uses the same contracts as 1.0.0.
		offRampAddress = setupOffRampV1_0_0(t, user, bc)
	case ccipdata.V1_2_0:
		offRampAddress = setupOffRampV1_2_0(t, user, bc)
	case ccipdata.V1_5_0:
		offRampAddress = setupOffRampV1_5_0(t, user, bc)
	default:
		require.Fail(t, "Unknown version: ", version)
	}

	feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

	// Create the version-specific reader.
	reader, err := factory.NewOffRampReader(log, factory.NewEvmVersionFinder(), ccipcalc.EvmAddrToGeneric(offRampAddress), bc, lp, nil, nil, true, feeEstimatorConfig)
	require.NoError(t, err)
	addr, err := reader.Address(ctx)
	require.NoError(t, err)
	require.Equal(t, ccipcalc.EvmAddrToGeneric(offRampAddress), addr)

	return offRampReaderTH{
		user:   user,
		reader: reader,
	}
}

func setupOffRampV1_0_0(t *testing.T, user *bind.TransactOpts, bc *client.SimulatedBackendClient) common.Address {
	onRampAddr := utils.RandomAddress()
	armAddr := deployMockArm(t, user, bc)
	csAddr := deployCommitStore(t, user, bc, onRampAddr, armAddr)

	// Deploy the OffRamp.
	staticConfig := evm_2_evm_offramp_1_0_0.EVM2EVMOffRampStaticConfig{
		CommitStore:         csAddr,
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onRampAddr,
		PrevOffRamp:         common.Address{},
		ArmProxy:            armAddr,
	}
	sourceTokens := []common.Address{}
	pools := []common.Address{}
	rateLimiterConfig := evm_2_evm_offramp_1_0_0.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(0),
		Rate:      big.NewInt(0),
	}

	offRampAddr, tx, offRamp, err := evm_2_evm_offramp_1_0_0.DeployEVM2EVMOffRamp(user, bc, staticConfig, sourceTokens, pools, rateLimiterConfig)
	bc.Commit()
	require.NoError(t, err)
	ccipdata.AssertNonRevert(t, tx, bc, user)

	// Verify the deployed OffRamp.
	tav, err := offRamp.TypeAndVersion(&bind.CallOpts{
		Context: testutils.Context(t),
	})
	require.NoError(t, err)
	require.Equal(t, "EVM2EVMOffRamp 1.0.0", tav)
	return offRampAddr
}

func setupOffRampV1_2_0(t *testing.T, user *bind.TransactOpts, bc *client.SimulatedBackendClient) common.Address {
	onRampAddr := utils.RandomAddress()
	armAddr := deployMockArm(t, user, bc)
	csAddr := deployCommitStore(t, user, bc, onRampAddr, armAddr)

	// Deploy the OffRamp.
	staticConfig := evm_2_evm_offramp_1_2_0.EVM2EVMOffRampStaticConfig{
		CommitStore:         csAddr,
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onRampAddr,
		PrevOffRamp:         common.Address{},
		ArmProxy:            armAddr,
	}
	sourceTokens := []common.Address{}
	pools := []common.Address{}
	rateLimiterConfig := evm_2_evm_offramp_1_2_0.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(0),
		Rate:      big.NewInt(0),
	}

	offRampAddr, tx, offRamp, err := evm_2_evm_offramp_1_2_0.DeployEVM2EVMOffRamp(user, bc, staticConfig, sourceTokens, pools, rateLimiterConfig)
	bc.Commit()
	require.NoError(t, err)
	ccipdata.AssertNonRevert(t, tx, bc, user)

	// Verify the deployed OffRamp.
	tav, err := offRamp.TypeAndVersion(&bind.CallOpts{
		Context: testutils.Context(t),
	})
	require.NoError(t, err)
	require.Equal(t, "EVM2EVMOffRamp 1.2.0", tav)
	return offRampAddr
}

func setupOffRampV1_5_0(t *testing.T, user *bind.TransactOpts, bc *client.SimulatedBackendClient) common.Address {
	onRampAddr := utils.RandomAddress()
	tokenAdminRegAddr := utils.RandomAddress()
	rmnAddr := deployMockArm(t, user, bc)
	csAddr := deployCommitStore(t, user, bc, onRampAddr, rmnAddr)

	// Deploy the OffRamp.
	staticConfig := evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
		CommitStore:         csAddr,
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onRampAddr,
		PrevOffRamp:         common.Address{},
		RmnProxy:            rmnAddr,
		TokenAdminRegistry:  tokenAdminRegAddr,
	}
	rateLimiterConfig := evm_2_evm_offramp.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(0),
		Rate:      big.NewInt(0),
	}

	offRampAddr, tx, offRamp, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(user, bc, staticConfig, rateLimiterConfig)
	bc.Commit()
	require.NoError(t, err)
	ccipdata.AssertNonRevert(t, tx, bc, user)

	// Verify the deployed OffRamp.
	tav, err := offRamp.TypeAndVersion(&bind.CallOpts{
		Context: testutils.Context(t),
	})
	require.NoError(t, err)
	require.Equal(t, "EVM2EVMOffRamp 1.5.0", tav)
	return offRampAddr
}

func deployMockArm(
	t *testing.T,
	user *bind.TransactOpts,
	bc *client.SimulatedBackendClient,
) common.Address {
	armAddr, tx, _, err := mock_rmn_contract.DeployMockRMNContract(user, bc)
	require.NoError(t, err)
	bc.Commit()
	ccipdata.AssertNonRevert(t, tx, bc, user)
	require.NotEqual(t, common.Address{}, armAddr)
	return armAddr
}

// Deploy the CommitStore. We use the same CommitStore version for all versions of OffRamp tested.
func deployCommitStore(
	t *testing.T,
	user *bind.TransactOpts,
	bc *client.SimulatedBackendClient,
	onRampAddress common.Address,
	armAddress common.Address,
) common.Address {
	// Deploy the CommitStore using the helper.
	csAddr, tx, cs, err := commit_store_helper.DeployCommitStoreHelper(user, bc, commit_store_helper.CommitStoreStaticConfig{
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onRampAddress,
		RmnProxy:            armAddress,
	})
	require.NoError(t, err)
	bc.Commit()
	ccipdata.AssertNonRevert(t, tx, bc, user)

	// Test the deployed CommitStore.
	callOpts := &bind.CallOpts{
		Context: testutils.Context(t),
	}
	tav, err := cs.TypeAndVersion(callOpts)
	require.NoError(t, err)
	require.Equal(t, "CommitStore 1.5.0", tav)
	return csAddr
}

func testOffRampReader(t *testing.T, th offRampReaderTH) {
	ctx := th.user.Context
	tokens, err := th.reader.GetTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, []cciptypes.Address{}, tokens.DestinationTokens)

	events, err := th.reader.GetExecutionStateChangesBetweenSeqNums(ctx, 0, 10, 0)
	require.NoError(t, err)
	require.Equal(t, []cciptypes.ExecutionStateChangedWithTxMeta{}, events)

	sourceToDestTokens, err := th.reader.GetSourceToDestTokensMapping(ctx)
	require.NoError(t, err)
	require.Empty(t, sourceToDestTokens)

	require.NoError(t, err)
}

func TestNewOffRampReader(t *testing.T) {
	var tt = []struct {
		typeAndVersion string
		expectedErr    string
	}{
		{
			typeAndVersion: "blah",
			expectedErr:    "unable to read type and version: invalid type and version blah",
		},
		{
			typeAndVersion: "CommitStore 1.0.0",
			expectedErr:    "expected EVM2EVMOffRamp got CommitStore",
		},
		{
			typeAndVersion: "EVM2EVMOffRamp 1.2.0",
			expectedErr:    "",
		},
		{
			typeAndVersion: "EVM2EVMOffRamp 2.0.0",
			expectedErr:    "unsupported offramp version 2.0.0",
		},
	}
	for _, tc := range tt {
		t.Run(tc.typeAndVersion, func(t *testing.T) {
			b, err := utils.ABIEncode(`[{"type":"string"}]`, tc.typeAndVersion)
			require.NoError(t, err)
			c := evmclientmocks.NewClient(t)

			feeEstimatorConfig := ccipdatamocks.NewFeeEstimatorConfigReader(t)

			c.On("CallContract", mock.Anything, mock.Anything, mock.Anything).Return(b, nil)
			addr := ccipcalc.EvmAddrToGeneric(utils.RandomAddress())
			lp := lpmocks.NewLogPoller(t)
			lp.On("RegisterFilter", mock.Anything, mock.Anything).Return(nil).Maybe()
			_, err = factory.NewOffRampReader(logger.TestLogger(t), factory.NewEvmVersionFinder(), addr, c, lp, nil, nil, true, feeEstimatorConfig)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
