package ccipdata_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	lpmocks "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store_helper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp_1_0_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type offRampReaderTH struct {
	user   *bind.TransactOpts
	reader ccipdata.OffRampReader
}

func TestOffRampFilters(t *testing.T) {
	assertFilterRegistration(t, new(lpmocks.LogPoller), func(lp *lpmocks.LogPoller, addr common.Address) ccipdata.Closer {
		c, err := ccipdata.NewOffRampV1_0_0(logger.TestLogger(t), addr, new(mocks.Client), lp, nil)
		require.NoError(t, err)
		require.NoError(t, c.RegisterFilters())
		return c
	}, 3)
	assertFilterRegistration(t, new(lpmocks.LogPoller), func(lp *lpmocks.LogPoller, addr common.Address) ccipdata.Closer {
		c, err := ccipdata.NewOffRampV1_2_0(logger.TestLogger(t), addr, new(mocks.Client), lp, nil)
		require.NoError(t, err)
		require.NoError(t, c.RegisterFilters())
		return c
	}, 3)
}

func TestExecOffchainConfig_Encoding(t *testing.T) {
	tests := map[string]struct {
		want      ccipdata.ExecOffchainConfig
		expectErr bool
	}{
		"encodes and decodes config with all fields set": {
			want: ccipdata.ExecOffchainConfig{
				SourceFinalityDepth:         3,
				DestOptimisticConfirmations: 6,
				DestFinalityDepth:           3,
				BatchGasLimit:               5_000_000,
				RelativeBoostPerWaitHour:    0.07,
				MaxGasPrice:                 200e9,
				InflightCacheExpiry:         models.MustMakeDuration(64 * time.Second),
				RootSnoozeTime:              models.MustMakeDuration(128 * time.Minute),
			},
		},
		"fails decoding when all fields present but with 0 values": {
			want: ccipdata.ExecOffchainConfig{
				SourceFinalityDepth:         0,
				DestFinalityDepth:           0,
				DestOptimisticConfirmations: 0,
				BatchGasLimit:               0,
				RelativeBoostPerWaitHour:    0,
				MaxGasPrice:                 0,
				InflightCacheExpiry:         models.MustMakeDuration(0),
				RootSnoozeTime:              models.MustMakeDuration(0),
			},
			expectErr: true,
		},
		"fails decoding when all fields are missing": {
			want:      ccipdata.ExecOffchainConfig{},
			expectErr: true,
		},
		"fails decoding when some fields are missing": {
			want: ccipdata.ExecOffchainConfig{
				SourceFinalityDepth: 99999999,
				InflightCacheExpiry: models.MustMakeDuration(64 * time.Second),
			},
			expectErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			exp := tc.want
			encode, err := ccipconfig.EncodeOffchainConfig(&exp)
			require.NoError(t, err)
			got, err := ccipconfig.DecodeOffchainConfig[ccipdata.ExecOffchainConfig](encode)

			if tc.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.want, got)
			}
		})
	}
}

func TestExecOnchainConfig100(t *testing.T) {
	tests := []struct {
		name      string
		want      ccipdata.ExecOnchainConfigV1_0_0
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: ccipdata.ExecOnchainConfigV1_0_0{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				Router:                                  utils.RandomAddress(),
				PriceRegistry:                           utils.RandomAddress(),
				MaxTokensLength:                         uint16(rand.Uint32()),
				MaxDataSize:                             rand.Uint32(),
			},
		},
		{
			name: "encodes and fails decoding config with missing fields",
			want: ccipdata.ExecOnchainConfigV1_0_0{
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

			decoded, err := abihelpers.DecodeAbiStruct[ccipdata.ExecOnchainConfigV1_0_0](encoded)
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
		want      ccipdata.ExecOnchainConfigV1_2_0
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: ccipdata.ExecOnchainConfigV1_2_0{
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
			want: ccipdata.ExecOnchainConfigV1_2_0{
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

			decoded, err := abihelpers.DecodeAbiStruct[ccipdata.ExecOnchainConfigV1_2_0](encoded)
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			th := setupOffRampReaderTH(t, test.version)
			testOffRampReader(t, th)
		})
	}
}

func setupOffRampReaderTH(t *testing.T, version string) offRampReaderTH {
	user, bc := newSimulation(t)
	log := logger.TestLogger(t)
	orm := logpoller.NewORM(testutils.SimulatedChainID, pgtest.NewSqlxDB(t), log, pgtest.NewQConfig(true))
	lp := logpoller.NewLogPoller(
		orm,
		bc,
		log,
		100*time.Millisecond, false, 2, 3, 2, 1000)

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
	default:
		require.Fail(t, "Unknown version: ", version)
	}

	// Create the version-specific reader.
	reader, err := ccipdata.NewOffRampReader(log, offRampAddress, bc, lp, nil)
	require.NoError(t, err)
	require.Equal(t, offRampAddress, reader.Address())

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
	assertNonRevert(t, tx, bc, user)

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
	staticConfig := evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
		CommitStore:         csAddr,
		ChainSelector:       testutils.SimulatedChainID.Uint64(),
		SourceChainSelector: testutils.SimulatedChainID.Uint64(),
		OnRamp:              onRampAddr,
		PrevOffRamp:         common.Address{},
		ArmProxy:            armAddr,
	}
	sourceTokens := []common.Address{}
	pools := []common.Address{}
	rateLimiterConfig := evm_2_evm_offramp.RateLimiterConfig{
		IsEnabled: false,
		Capacity:  big.NewInt(0),
		Rate:      big.NewInt(0),
	}

	offRampAddr, tx, offRamp, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(user, bc, staticConfig, sourceTokens, pools, rateLimiterConfig)
	bc.Commit()
	require.NoError(t, err)
	assertNonRevert(t, tx, bc, user)

	// Verify the deployed OffRamp.
	tav, err := offRamp.TypeAndVersion(&bind.CallOpts{
		Context: testutils.Context(t),
	})
	require.NoError(t, err)
	require.Equal(t, "EVM2EVMOffRamp 1.2.0", tav)
	return offRampAddr
}

func deployMockArm(
	t *testing.T,
	user *bind.TransactOpts,
	bc *client.SimulatedBackendClient,
) common.Address {
	armAddr, tx, _, err := mock_arm_contract.DeployMockARMContract(user, bc)
	require.NoError(t, err)
	bc.Commit()
	assertNonRevert(t, tx, bc, user)
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
		ArmProxy:            armAddress,
	})
	require.NoError(t, err)
	bc.Commit()
	assertNonRevert(t, tx, bc, user)

	// Test the deployed CommitStore.
	callOpts := &bind.CallOpts{
		Context: testutils.Context(t),
	}
	tav, err := cs.TypeAndVersion(callOpts)
	require.NoError(t, err)
	require.Equal(t, "CommitStore 1.2.0", tav)
	return csAddr
}

func testOffRampReader(t *testing.T, th offRampReaderTH) {
	ctx := th.user.Context
	addresses, err := th.reader.GetDestinationTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, []common.Address{}, addresses)

	tokens, err := th.reader.GetSupportedTokens(ctx)
	require.NoError(t, err)
	require.Equal(t, []common.Address{}, tokens)

	events, err := th.reader.GetExecutionStateChangesBetweenSeqNums(ctx, 0, 10, 0)
	require.NoError(t, err)
	require.Equal(t, []ccipdata.Event[ccipdata.ExecutionStateChanged]{}, events)

	destTokens, err := th.reader.GetDestinationTokensFromSourceTokens(ctx, tokens)
	require.NoError(t, err)
	require.Empty(t, destTokens)

	rateLimits, err := th.reader.GetTokenPoolsRateLimits(ctx, []common.Address{})
	require.NoError(t, err)
	require.Empty(t, rateLimits)
}
