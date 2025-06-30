package aptos

import (
	"math/big"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip"
	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	aptoschain "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
	aptosstate "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/aptos"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestDeployAptosChainImp_VerifyPreconditions(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		env       cldf.Environment
		config    config.DeployAptosChainConfig
		wantErrRe string
		wantErr   bool
	}{
		{
			name: "success - valid configs",
			env: cldf.Environment{
				Name:              "test",
				Logger:            logger.TestLogger(t),
				ExistingAddresses: cldf.NewMemoryAddressBook(),
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						743186221051783445:  aptoschain.Chain{},
						4457093679053095497: aptoschain.Chain{},
					}),
			},
			config: config.DeployAptosChainConfig{
				ContractParamsPerChain: map[uint64]config.ChainContractParams{
					4457093679053095497: GetMockChainContractParams(t, 4457093679053095497),
					743186221051783445:  GetMockChainContractParams(t, 743186221051783445),
				},
				MCMSDeployConfigPerChain: map[uint64]types.MCMSWithTimelockConfigV2{
					4457093679053095497: getMockMCMSConfig(t),
					743186221051783445:  getMockMCMSConfig(t),
				},
			},
			wantErr: false,
		},
		{
			name: "success - valid config w MCMS deployed",
			env: cldf.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				ExistingAddresses: getTestAddressBook(
					t,
					map[uint64]map[string]cldf.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: shared.AptosMCMSType},
						},
						743186221051783445: {
							mockMCMSAddress: {Type: shared.AptosMCMSType},
						},
					},
				),
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						743186221051783445:  aptoschain.Chain{},
						4457093679053095497: aptoschain.Chain{},
					}),
			},
			config: config.DeployAptosChainConfig{
				ContractParamsPerChain: map[uint64]config.ChainContractParams{
					4457093679053095497: GetMockChainContractParams(t, 4457093679053095497),
					743186221051783445:  GetMockChainContractParams(t, 743186221051783445),
				},
			},
			wantErr: false,
		},
		{
			name: "error - invalid min delay",
			env: cldf.Environment{
				Name:              "test",
				Logger:            logger.TestLogger(t),
				ExistingAddresses: cldf.NewMemoryAddressBook(),
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						4457093679053095497: aptoschain.Chain{},
					}),
			},
			config: config.DeployAptosChainConfig{
				ContractParamsPerChain: map[uint64]config.ChainContractParams{
					4457093679053095497: GetMockChainContractParams(t, 4457093679053095497),
				},
				MCMSDeployConfigPerChain: map[uint64]types.MCMSWithTimelockConfigV2{
					4457093679053095497: {
						Canceller:        proposalutils.SingleGroupMCMSV2(t),
						Proposer:         proposalutils.SingleGroupMCMSV2(t),
						Bypasser:         proposalutils.SingleGroupMCMSV2(t),
						TimelockMinDelay: nil, // Invalid min delay
					},
				},
			},
			wantErr:   true,
			wantErrRe: `invalid MCMS timelock min delay for Aptos chain 4457093679053095497`,
		},
		{
			name: "error - chain has no env",
			env: cldf.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				ExistingAddresses: getTestAddressBook(
					t,
					map[uint64]map[string]cldf.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: shared.AptosMCMSType},
						},
						743186221051783445: {
							mockMCMSAddress: {Type: shared.AptosMCMSType},
						},
					},
				),
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						4457093679053095497: aptoschain.Chain{},
					}),
			},
			config: config.DeployAptosChainConfig{
				ContractParamsPerChain: map[uint64]config.ChainContractParams{
					4457093679053095497: GetMockChainContractParams(t, 4457093679053095497),
					743186221051783445:  GetMockChainContractParams(t, 743186221051783445),
				},
			},
			wantErrRe: `chain 743186221051783445 not found in env`,
			wantErr:   true,
		},
		{
			name: "error - invalid config - chainSelector",
			env: cldf.Environment{
				Name:              "test",
				Logger:            logger.TestLogger(t),
				ExistingAddresses: cldf.NewMemoryAddressBook(),
			},
			config: config.DeployAptosChainConfig{
				ContractParamsPerChain: map[uint64]config.ChainContractParams{
					1: {},
				},
			},
			wantErrRe: "invalid chain selector:",
			wantErr:   true,
		},
		{
			name: "error - missing MCMS config for chain without MCMS deployed",
			env: cldf.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				ExistingAddresses: getTestAddressBook(
					t,
					map[uint64]map[string]cldf.TypeAndVersion{
						4457093679053095497: {}, // No MCMS address in state
					},
				),
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						4457093679053095497: aptoschain.Chain{},
					}),
			},
			config: config.DeployAptosChainConfig{
				ContractParamsPerChain: map[uint64]config.ChainContractParams{
					4457093679053095497: GetMockChainContractParams(t, 4457093679053095497),
				},
				// MCMSDeployConfigPerChain is missing needed configs
			},
			wantErrRe: `invalid mcms configs for Aptos chain 4457093679053095497`,
			wantErr:   true,
		},
		{
			name: "error - invalid config for chain",
			env: cldf.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				ExistingAddresses: getTestAddressBook(
					t,
					map[uint64]map[string]cldf.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: shared.AptosMCMSType}, // MCMS already deployed
						},
					},
				),
				BlockChains: chain.NewBlockChains(
					map[uint64]chain.BlockChain{
						4457093679053095497: aptoschain.Chain{},
					}),
			},
			config: config.DeployAptosChainConfig{
				ContractParamsPerChain: map[uint64]config.ChainContractParams{
					4457093679053095497: {
						FeeQuoterParams: config.FeeQuoterParams{
							TokenPriceStalenessThreshold: 0, // Invalid gas limit (assuming 0 is invalid)
						},
					},
				},
			},
			wantErrRe: `invalid config for Aptos chain 4457093679053095497`,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			cs := DeployAptosChain{}
			err := cs.VerifyPreconditions(tt.env, tt.config)
			if tt.wantErr {
				require.Error(t, err)
				errStr := err.Error()
				assert.Regexp(t, tt.wantErrRe, errStr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeployAptosChain_Apply(t *testing.T) {
	t.Parallel()
	lggr := logger.TestLogger(t)

	// Setup memory environment with 1 Aptos chain
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		AptosChains: 1,
	})

	// Get chain selectors
	aptosChainSelectors := env.BlockChains.ListChainSelectors(chain.WithFamily(chain_selectors.FamilyAptos))
	require.Len(t, aptosChainSelectors, 1, "Expected exactly 1 Aptos chain")
	chainSelector := aptosChainSelectors[0]
	t.Log("Deployer: ", env.BlockChains.AptosChains()[chainSelector].DeployerSigner)

	// Deploy CCIP to Aptos chain
	mockCCIPParams := GetMockChainContractParams(t, chainSelector)
	ccipConfig := config.DeployAptosChainConfig{
		ContractParamsPerChain: map[uint64]config.ChainContractParams{
			chainSelector: mockCCIPParams,
		},
		MCMSDeployConfigPerChain: map[uint64]types.MCMSWithTimelockConfigV2{
			chainSelector: {
				Canceller:        proposalutils.SingleGroupMCMSV2(t),
				Proposer:         proposalutils.SingleGroupMCMSV2(t),
				Bypasser:         proposalutils.SingleGroupMCMSV2(t),
				TimelockMinDelay: big.NewInt(1),
			},
		},
		MCMSTimelockConfigPerChain: map[uint64]proposalutils.TimelockConfig{
			chainSelector: {
				MinDelay:     time.Duration(1) * time.Second,
				MCMSAction:   mcmstypes.TimelockActionSchedule,
				OverrideRoot: false,
			},
		},
	}
	env, _, err := commonchangeset.ApplyChangesets(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(DeployAptosChain{}, ccipConfig),
	})
	require.NoError(t, err)

	// Verify CCIP deployment state by binding ccip contract and checking if it's deployed
	states, err := aptosstate.LoadOnchainStateAptos(env)
	require.NoError(t, err)
	require.NotNil(t, states[chainSelector], "No state found for chain")
	state := states[chainSelector]

	ccipAddr := state.CCIPAddress
	require.NotEmpty(t, ccipAddr, "CCIP address should not be empty")

	// Bind CCIP offramp contract
	offrampBind := ccip_offramp.Bind(ccipAddr, env.BlockChains.AptosChains()[chainSelector].Client)
	offRampSourceConfig, err := offrampBind.Offramp().GetSourceChainConfig(nil, mockCCIPParams.OffRampParams.SourceChainSelectors[0])
	require.NoError(t, err)
	require.True(t, offRampSourceConfig.IsEnabled, "contracts were not initialized correctly")

	// Check premium multiplier
	ccipBind := ccip.Bind(ccipAddr, env.BlockChains.AptosChains()[chainSelector].Client)
	require.NotEqual(t, aptos.AccountAddress{}, state.LinkTokenAddress, "Link token address should not be empty")
	mult, err := ccipBind.FeeQuoter().GetPremiumMultiplierWeiPerEth(nil, state.LinkTokenAddress)
	require.NoError(t, err)
	require.Equal(t, mockCCIPParams.FeeQuoterParams.PremiumMultiplierWeiPerEthByFeeToken[shared.LinkSymbol], mult)

	aptTokenAdd := MustParseAddress(t, "0xa")
	mult, err = ccipBind.FeeQuoter().GetPremiumMultiplierWeiPerEth(nil, aptTokenAdd)
	require.NoError(t, err)
	require.Equal(t, mockCCIPParams.FeeQuoterParams.PremiumMultiplierWeiPerEthByFeeToken[shared.APTSymbol], mult)
}
