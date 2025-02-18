package changeset

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"

	owner_helpers "github.com/smartcontractkit/ccip-owner-contracts/pkg/gethwrappers"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/proposalutils"
	"github.com/smartcontractkit/chainlink/deployment/common/types"

	"github.com/stretchr/testify/require"
)

func TestAddressesContainBundle(t *testing.T) {
	v100 := semver.MustParse("1.0.0")

	tests := []struct {
		name      string
		addrs     map[string]deployment.TypeAndVersion
		wantTypes []deployment.TypeAndVersion
		want      bool
		wantErr   string
	}{
		{
			name: "Contains all types exactly once",
			addrs: map[string]deployment.TypeAndVersion{
				"addr1": {Type: "type1", Version: *v100},
				"addr2": {Type: "type2", Version: *v100, Labels: deployment.NewLabelSet("label1")},
			},
			wantTypes: []deployment.TypeAndVersion{
				{Type: "type1", Version: *v100},
				{Type: "type2", Version: *v100, Labels: deployment.NewLabelSet("label1")},
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "Missing a type",
			addrs: map[string]deployment.TypeAndVersion{
				"addr1": {Type: "type1", Version: *v100},
			},
			wantTypes: []deployment.TypeAndVersion{
				{Type: "type1", Version: *v100},
				{Type: "type2", Version: *v100},
			},
			want:    false,
			wantErr: "missing contracts",
		},
		{
			name: "More than one instance of a type",
			addrs: map[string]deployment.TypeAndVersion{
				"addr1": {Type: "type1", Version: *v100},
				"addr2": {Type: "type1", Version: *v100},
			},
			wantTypes: []deployment.TypeAndVersion{
				{Type: "type1", Version: *v100},
			},
			want:    false,
			wantErr: "found more than one instance of contract type1 1.0.0 (labels=)",
		},
		{
			name:      "Empty inputs",
			addrs:     map[string]deployment.TypeAndVersion{},
			wantTypes: []deployment.TypeAndVersion{},
			want:      true,
			wantErr:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deployment.AddressesContainBundle(tt.addrs, tt.wantTypes)

			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}

func TestMaybeLoadMCMSWithTimelockChainState(t *testing.T) {
	type testCase struct {
		name      string
		chain     deployment.Chain
		addresses map[string]deployment.TypeAndVersion
		wantState *MCMSWithTimelockState // Expected state
		wantErr   string
		skip      bool
	}

	defaultChain := deployment.Chain{
		Selector: chainsel.ETHEREUM_TESTNET_SEPOLIA.Selector,
	}
	tests := []testCase{
		{
			name:  "Valid contracts",
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
			name:  "Missing contract",
			chain: defaultChain,
			addresses: map[string]deployment.TypeAndVersion{
				"0x123": deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0),
				"0x456": deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0),
				"0x789": deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0),
				"0xabc": deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0),
				// "0xdef": deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0), // Missing
			},
			wantState: nil,
			wantErr:   "missing a contract",
		},
		{
			name: "labelled multichain contract proposer role overrides existing proposer",
			// TODO(mstreet3): Skipping because this behavior is undefined.  The addresses are a map and there is no guarantee that
			// a labeled contract will override a named contract type.  Leaving this test case in until this behavior is resolved.
			skip:  true,
			chain: defaultChain,

			// addresses has two proposer contract candidates, which should make it into the final state?
			addresses: map[string]deployment.TypeAndVersion{
				"0x123": deployment.NewTypeAndVersion(types.RBACTimelock, deployment.Version1_0_0),
				"0x456": deployment.NewTypeAndVersion(types.CallProxy, deployment.Version1_0_0),
				"0x789": deployment.NewTypeAndVersion(types.ProposerManyChainMultisig, deployment.Version1_0_0),
				"0xabc": deployment.NewTypeAndVersion(types.CancellerManyChainMultisig, deployment.Version1_0_0),
				"0xdef": deployment.NewTypeAndVersion(types.BypasserManyChainMultisig, deployment.Version1_0_0),
				"0xaaa": func() deployment.TypeAndVersion {
					tv := deployment.NewTypeAndVersion(types.ManyChainMultisig, deployment.Version1_0_0)
					tv.Labels.Add(types.ProposerRole.String())
					return tv
				}(),
			},
		},
	}

	for _, tt := range tests {
		if tt.skip {
			continue
		}

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
