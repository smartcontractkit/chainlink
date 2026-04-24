package example_test

import (
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	chain_selectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/runtime"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset/example"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/common/types"
)

// setupLinkTransferRuntime deploys all required contracts on a simulated chain to run tests which
// operate on the link token and MCMS contracts.
//
// Returns the test runtime and the chain selector.
func setupLinkTransferRuntime(t *testing.T) (*runtime.Runtime, uint64) {
	t.Helper()

	selector := chain_selectors.TEST_90000001.Selector
	rt, err := runtime.New(t.Context(), runtime.WithEnvOpts(
		environment.WithEVMSimulated(t, []uint64{selector}),
	))
	require.NoError(t, err)

	// Deploy MCMS and Timelock
	config := proposalutils.SingleGroupMCMSV2(t)
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployLinkToken), []uint64{selector}),
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(changeset.DeployMCMSWithTimelockV2), map[uint64]types.MCMSWithTimelockConfigV2{
			selector: {
				Canceller:        config,
				Bypasser:         config,
				Proposer:         config,
				TimelockMinDelay: big.NewInt(0),
			},
		}),
	)
	require.NoError(t, err)

	return rt, selector
}

func TestValidate(t *testing.T) {
	rt, selector := setupLinkTransferRuntime(t)

	chain := rt.Environment().BlockChains.EVMChains()[selector]
	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)

	tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = linkState.LinkToken.Mint(chain.DeployerKey, chain.DeployerKey.From, big.NewInt(750))
	require.NoError(t, err)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tests := []struct {
		name     string
		cfg      example.LinkTransferConfig
		errorMsg string
	}{
		{
			name: "valid config",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {{To: mcmsState.Timelock.Address(), Value: big.NewInt(100)}}},
				From: chain.DeployerKey.From,
				McmsConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Hour,
				},
			},
		},
		{
			name: "valid non mcms config",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {{To: mcmsState.Timelock.Address(), Value: big.NewInt(100)}}},
				From: chain.DeployerKey.From,
			},
		},
		{
			name: "insufficient funds",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {
						{To: chain.DeployerKey.From, Value: big.NewInt(100)},
						{To: chain.DeployerKey.From, Value: big.NewInt(500)},
						{To: chain.DeployerKey.From, Value: big.NewInt(1250)},
					},
				},
				From: mcmsState.Timelock.Address(),
				McmsConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Hour,
				},
			},
			errorMsg: "sender does not have enough funds for transfers for chain selector 909606746561742123, required: 1850, available: 0",
		},
		{
			name:     "invalid config: empty transfers",
			cfg:      example.LinkTransferConfig{Transfers: map[uint64][]example.TransferConfig{}},
			errorMsg: "transfers map must have at least one chainSel",
		},
		{
			name: "invalid chain selector",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					1: {{To: common.Address{}, Value: big.NewInt(100)}}},
			},
			errorMsg: "invalid chain selector: unknown chain selector 1",
		},
		{
			name: "chain selector not found",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					chain_selectors.ETHEREUM_TESTNET_GOERLI_ARBITRUM_1.Selector: {{To: common.Address{}, Value: big.NewInt(100)}}},
			},
			errorMsg: "chain with selector 6101244977088475029 not found",
		},
		{
			name: "empty transfer list",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {},
				},
			},
			errorMsg: "transfers for chainSel 909606746561742123 must have at least one LinkTransfer",
		},
		{
			name: "empty value",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {
						{To: chain.DeployerKey.From, Value: nil},
					},
				},
			},
			errorMsg: "value for transfers must be set",
		},
		{
			name: "zero value",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {
						{To: chain.DeployerKey.From, Value: big.NewInt(0)},
					},
				},
			},
			errorMsg: "value for transfers must be non-zero",
		},
		{
			name: "negative value",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {
						{To: chain.DeployerKey.From, Value: big.NewInt(-5)},
					},
				},
			},
			errorMsg: "value for transfers must be positive",
		},
		{
			name: "non-evm-chain",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					chain_selectors.APTOS_MAINNET.Selector: {{To: mcmsState.Timelock.Address(), Value: big.NewInt(100)}}},
				From: chain.DeployerKey.From,
			},
			errorMsg: "chain selector 4741433654826277614 is not an EVM chain",
		},
		{
			name: "delay greater than max allowed",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {{To: mcmsState.Timelock.Address(), Value: big.NewInt(100)}}},
				From: chain.DeployerKey.From,
				McmsConfig: &proposalutils.TimelockConfig{
					MinDelay: time.Hour * 24 * 10,
				},
			},
			errorMsg: "minDelay must be less than 7 days",
		},
		{
			name: "invalid config: transfer to address missing",
			cfg: example.LinkTransferConfig{
				Transfers: map[uint64][]example.TransferConfig{
					selector: {{To: common.Address{}, Value: big.NewInt(100)}}},
			},
			errorMsg: "'to' address for transfers must be set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate(rt.Environment())
			if tt.errorMsg != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLinkTransferMCMSV2(t *testing.T) {
	t.Parallel()

	var (
		ctx          = t.Context()
		rt, selector = setupLinkTransferRuntime(t)
	)

	chain := rt.Environment().BlockChains.EVMChains()[selector]
	addrs, err := rt.State().AddressBook.AddressesForChain(selector)
	require.NoError(t, err)
	require.Len(t, addrs, 6)

	mcmsState, err := changeset.MaybeLoadMCMSWithTimelockChainState(chain, addrs)
	require.NoError(t, err)
	linkState, err := changeset.MaybeLoadLinkTokenChainState(chain, addrs)
	require.NoError(t, err)
	timelockAddress := mcmsState.Timelock.Address()

	// Mint some funds
	// grant minter permissions
	tx, err := linkState.LinkToken.GrantMintRole(chain.DeployerKey, chain.DeployerKey.From)
	require.NoError(t, err)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	tx, err = linkState.LinkToken.Mint(chain.DeployerKey, timelockAddress, big.NewInt(750))
	require.NoError(t, err)
	_, err = cldf.ConfirmIfNoError(chain, tx, err)
	require.NoError(t, err)

	// Apply the changeset
	err = rt.Exec(
		runtime.ChangesetTask(cldf.CreateLegacyChangeSet(example.LinkTransferV2), &example.LinkTransferConfig{
			From: timelockAddress,
			Transfers: map[uint64][]example.TransferConfig{
				selector: {{To: chain.DeployerKey.From, Value: big.NewInt(500)}},
			},
			McmsConfig: &proposalutils.TimelockConfig{
				MinDelay:     0,
				OverrideRoot: true,
			},
		}),
		runtime.SignAndExecuteProposalsTask([]*ecdsa.PrivateKey{proposalutils.TestXXXMCMSSigner}),
	)
	require.NoError(t, err)

	// Check new balances
	endBalance, err := linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, chain.DeployerKey.From)
	require.NoError(t, err)
	expectedBalance := big.NewInt(500)
	require.Equal(t, expectedBalance, endBalance)

	// check timelock balance
	endBalance, err = linkState.LinkToken.BalanceOf(&bind.CallOpts{Context: ctx}, timelockAddress)
	require.NoError(t, err)
	expectedBalance = big.NewInt(250)
	require.Equal(t, expectedBalance, endBalance)
}
