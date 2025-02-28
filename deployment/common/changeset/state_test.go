package changeset

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"

	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/stretchr/testify/require"
)

func TestMaybeLoadMCMSWithTimelockChainState(t *testing.T) {
	type testCase struct {
		name      string
		chain     deployment.Chain
		addresses map[string]deployment.TypeAndVersion
		wantState *MCMSWithTimelockState // Expected state
		wantErr   string
	}

	defaultChain := deployment.Chain{
		Selector: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
	}
	tests := []testCase{
		{
			name:  "OK_load all contracts from address book",
			chain: defaultChain,
			addresses: map[string]deployment.TypeAndVersion{
				"0x123": deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0),
				"0x456": deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0),
				"0x789": deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0),
				"0xabc": deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0),
				"0xdef": deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0),
			},
			wantState: &MCMSWithTimelockState{
				MCMSWithTimelockContracts: &proposalutils.MCMSWithTimelockContracts{
					Timelock: func() *owner_helpers.RBACTimelock {
						tl, err := owner_helpers.NewRBACTimelock(common.HexToAddress("0x123"), nil)
						require.NoError(t, err)
						return tl
					}(),
					CallProxy: func() *owner_helpers.CallProxy {
						cp, err := owner_helpers.NewCallProxy(common.HexToAddress("0x456"), nil)
						require.NoError(t, err)
						return cp
					}(),
					ProposerMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0x789"), nil)
						require.NoError(t, err)
						return mcms
					}(),
					CancellerMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0xabc"), nil)
						require.NoError(t, err)
						return mcms
					}(),
					BypasserMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0xdef"), nil)
						require.NoError(t, err)
						return mcms
					}(),
				},
			},
			wantErr: "",
		},
		{
			name:  "OK_labelled multichain contract is selected if only option",
			chain: defaultChain,
			addresses: map[string]deployment.TypeAndVersion{
				"0x123": deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0),
				"0x456": deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0),
				"0xabc": deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0),
				"0xdef": deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0),
				"0xaaa": func() deployment.TypeAndVersion {
					tv := deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
					tv.Labels.Add(types.ProposerRole.String())
					return tv
				}(),
			},
			wantState: &MCMSWithTimelockState{
				MCMSWithTimelockContracts: &proposalutils.MCMSWithTimelockContracts{
					Timelock: func() *owner_helpers.RBACTimelock {
						tl, err := owner_helpers.NewRBACTimelock(common.HexToAddress("0x123"), nil)
						require.NoError(t, err)
						return tl
					}(),
					CallProxy: func() *owner_helpers.CallProxy {
						cp, err := owner_helpers.NewCallProxy(common.HexToAddress("0x456"), nil)
						require.NoError(t, err)
						return cp
					}(),
					CancellerMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0xabc"), nil)
						require.NoError(t, err)
						return mcms
					}(),
					BypasserMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0xdef"), nil)
						require.NoError(t, err)
						return mcms
					}(),
					ProposerMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0xaaa"), nil)
						require.NoError(t, err)
						return mcms
					}(),
				},
			},
		},
		{
			// A ManyChainMultiSig contract can be labeled or it can be typed.  A typed contract should be selected
			// over a labeled contract.
			name:  "OK_labelled multichain contract has lower selection precedence",
			chain: defaultChain,
			addresses: map[string]deployment.TypeAndVersion{
				"0x123": deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0),
				"0x456": deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0),
				"0xabc": deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0),
				"0xdef": deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0),

				// Two contracts can satisfy the ProposerMcm field.  This test ensures that the typed contract is
				// returned over the labeled contract.
				"0x789": deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0),
				"0xaaa": func() deployment.TypeAndVersion {
					tv := deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
					tv.Labels.Add(types.ProposerRole.String())
					return tv
				}(),
			},
			wantState: &MCMSWithTimelockState{
				MCMSWithTimelockContracts: &proposalutils.MCMSWithTimelockContracts{
					Timelock: func() *owner_helpers.RBACTimelock {
						tl, err := owner_helpers.NewRBACTimelock(common.HexToAddress("0x123"), nil)
						require.NoError(t, err)
						return tl
					}(),
					CallProxy: func() *owner_helpers.CallProxy {
						cp, err := owner_helpers.NewCallProxy(common.HexToAddress("0x456"), nil)
						require.NoError(t, err)
						return cp
					}(),
					CancellerMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0xabc"), nil)
						require.NoError(t, err)
						return mcms
					}(),
					BypasserMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0xdef"), nil)
						require.NoError(t, err)
						return mcms
					}(),
					ProposerMcm: func() *owner_helpers.ManyChainMultiSig {
						mcms, err := owner_helpers.NewManyChainMultiSig(common.HexToAddress("0x789"), nil)
						require.NoError(t, err)
						return mcms
					}(),
				},
			},
		},
		{
			name:  "NOK_multiple contracts of same type",
			chain: defaultChain,
			addresses: map[string]deployment.TypeAndVersion{
				"0x123": deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0),
				"0x456": deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0),
				"0xabc": deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0),
				"0xdef": deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0),

				"0x789": deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0),
				"0xaaa": deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0), // duplicate
			},
			wantErr: "unable to check MCMS contracts",
		},
		{
			name:  "NOK_multiple generic MCMS contracts with same label",
			chain: defaultChain,
			addresses: map[string]deployment.TypeAndVersion{
				"0x123": deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0),
				"0x456": deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0),
				"0xabc": deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0),
				"0xdef": deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0),

				"0x789": deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0),
				"0xaaa": func() deployment.TypeAndVersion {
					tv := deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
					tv.Labels.Add(types.ProposerRole.String())
					return tv
				}(),
				"0xbbb": func() deployment.TypeAndVersion {
					tv := deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
					tv.Labels.Add(types.ProposerRole.String())
					return tv
				}(),
			},
			wantErr: "unable to check MCMS contracts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotState, err := MaybeLoadMCMSWithTimelockChainState(tt.chain, tt.addresses)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantState.MCMSWithTimelockContracts.Timelock.Address(), gotState.MCMSWithTimelockContracts.Timelock.Address())
			require.Equal(t, tt.wantState.MCMSWithTimelockContracts.CallProxy.Address(), gotState.MCMSWithTimelockContracts.CallProxy.Address())
			require.Equal(t, tt.wantState.MCMSWithTimelockContracts.ProposerMcm.Address(), gotState.MCMSWithTimelockContracts.ProposerMcm.Address())
			require.Equal(t, tt.wantState.MCMSWithTimelockContracts.CancellerMcm.Address(), gotState.MCMSWithTimelockContracts.CancellerMcm.Address())
			require.Equal(t, tt.wantState.MCMSWithTimelockContracts.BypasserMcm.Address(), gotState.MCMSWithTimelockContracts.BypasserMcm.Address())
		})
	}
}
