package aptos

import (
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/aptos/config"
	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

func TestDeployAptosChainImp_VerifyPreconditions(t *testing.T) {
	tests := []struct {
		name      string
		env       deployment.Environment
		config    config.DeployAptosChainConfig
		wantErrRe string
		wantErr   bool
	}{
		{
			name: "success - valid configs",
			env: deployment.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				AptosChains: map[uint64]deployment.AptosChain{
					743186221051783445:  {},
					4457093679053095497: {},
				},
				ExistingAddresses: deployment.NewMemoryAddressBook(),
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
			env: deployment.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				AptosChains: map[uint64]deployment.AptosChain{
					743186221051783445:  {},
					4457093679053095497: {},
				},
				ExistingAddresses: getTestAddressBook(
					map[uint64]map[string]deployment.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: changeset.AptosMCMSType},
						},
						743186221051783445: {
							mockMCMSAddress: {Type: changeset.AptosMCMSType},
						},
					},
				),
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
			name: "error - chain has no env",
			env: deployment.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				AptosChains: map[uint64]deployment.AptosChain{
					4457093679053095497: {},
				},
				ExistingAddresses: getTestAddressBook(
					map[uint64]map[string]deployment.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: changeset.AptosMCMSType},
						},
						743186221051783445: {
							mockMCMSAddress: {Type: changeset.AptosMCMSType},
						},
					},
				),
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
			env: deployment.Environment{
				Name:              "test",
				Logger:            logger.TestLogger(t),
				ExistingAddresses: deployment.NewMemoryAddressBook(),
				AptosChains:       map[uint64]deployment.AptosChain{},
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
			env: deployment.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				AptosChains: map[uint64]deployment.AptosChain{
					4457093679053095497: {},
				},
				ExistingAddresses: getTestAddressBook(
					map[uint64]map[string]deployment.TypeAndVersion{
						4457093679053095497: {}, // No MCMS address in state
					},
				),
			},
			config: config.DeployAptosChainConfig{
				ContractParamsPerChain: map[uint64]config.ChainContractParams{
					4457093679053095497: GetMockChainContractParams(t, 4457093679053095497),
				},
				// MCMSDeployConfigPerChain is missing needed configs
			},
			wantErrRe: `invalid mcms configs for chain 4457093679053095497`,
			wantErr:   true,
		},
		{
			name: "error - invalid config for chain",
			env: deployment.Environment{
				Name:   "test",
				Logger: logger.TestLogger(t),
				AptosChains: map[uint64]deployment.AptosChain{
					4457093679053095497: {},
				},
				ExistingAddresses: getTestAddressBook(
					map[uint64]map[string]deployment.TypeAndVersion{
						4457093679053095497: {
							mockMCMSAddress: {Type: changeset.AptosMCMSType}, // MCMS already deployed
						},
					},
				),
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
			wantErrRe: `invalid config for chain 4457093679053095497`,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := DeployAptosChain{}
			err := cs.VerifyPreconditions(tt.env, tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				errStr := err.Error()
				assert.Regexp(t, tt.wantErrRe, errStr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDeployAptosChain_Apply(t *testing.T) {
	lggr := logger.TestLogger(t)

	// Setup memory environment with 1 Aptos chain
	env := memory.NewMemoryEnvironment(t, lggr, zapcore.InfoLevel, memory.MemoryEnvironmentConfig{
		AptosChains: 1,
	})

	// Get chain selectors
	aptosChainSelectors := env.AllChainSelectorsAptos()
	require.Equal(t, 1, len(aptosChainSelectors), "Expected exactly 1 Aptos chain")
	chainSelector := aptosChainSelectors[0]
	t.Log("Deployer: ", env.AptosChains[chainSelector].DeployerSigner)

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
				TimelockMinDelay: big.NewInt(0),
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
	env, _, err := commonchangeset.ApplyChangesetsV2(t, env, []commonchangeset.ConfiguredChangeSet{
		commonchangeset.Configure(DeployAptosChain{}, ccipConfig),
	})
	require.NoError(t, err)

	// Verify CCIP deployment state by binding ccip contract and checking if it's deployed
	state, err := changeset.LoadOnchainStateAptos(env)
	require.NoError(t, err)
	require.NotNil(t, state[chainSelector], "No state found for chain")

	ccipAddr := state[chainSelector].CCIPAddress
	require.NotEmpty(t, ccipAddr, "CCIP address should not be empty")

	// Bind CCIP contract
	offrampBind := ccip_offramp.Bind(ccipAddr, env.AptosChains[chainSelector].Client)
	offRampSourceConfig, err := offrampBind.Offramp().GetSourceChainConfig(nil, mockCCIPParams.OffRampParams.SourceChainSelectors[0])
	require.NoError(t, err)
	require.Equal(t, true, offRampSourceConfig.IsEnabled, "contracts were not initialized correctly")
}
