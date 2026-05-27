package changeset

import (
	"math/big"
	"testing"

	"github.com/Masterminds/semver/v3"
	chainselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldfproposalutils "github.com/smartcontractkit/chainlink-deployments-framework/engine/cld/mcms/proposalutils"
	"github.com/smartcontractkit/chainlink-deployments-framework/engine/test/environment"

	commontypes "github.com/smartcontractkit/chainlink/deployment/common/types"
	"github.com/smartcontractkit/chainlink/deployment/vault/changeset/types"
)

const (
	testCustomTimelockQualifier = "vault-test-timelock-qualifier"

	testTimelockAddrEmptyQualifier  = "0xa111111111111111111111111111111111111111"
	testTimelockAddrCustomQualifier = "0xa222222222222222222222222222222222222222"
	testProposerAddrEmptyQualifier  = "0xa333333333333333333333333333333333333333"
	testProposerAddrCustomQualifier = "0xa444444444444444444444444444444444444444"
)

var testMCMSContractVersion = semver.MustParse("1.0.0")

func TestTransferERC20Validation(t *testing.T) {
	t.Parallel()

	env, err := environment.New(t.Context())
	require.NoError(t, err)

	tests := []struct {
		name      string
		config    types.TransferERC20Config
		wantError bool
		errorMsg  string
	}{
		{
			name: "empty transfers_by_chain",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{},
				MCMSConfig:       &cldfproposalutils.TimelockConfig{},
			},
			wantError: true,
			errorMsg:  "transfers_by_chain must not be empty",
		},
		{
			name: "missing MCMS config",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					testChainID: {
						{
							Payee:  testAddr1,
							Token:  testAddr2,
							Amount: big.NewInt(1),
						},
					},
				},
			},
			wantError: true,
			errorMsg:  "MCMSConfig is required",
		},
		{
			name: "payee not whitelisted",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					testChainID: {
						{
							Payee:  testAddr1,
							Token:  testAddr2,
							Amount: big.NewInt(1),
						},
					},
				},
				MCMSConfig: &cldfproposalutils.TimelockConfig{},
			},
			wantError: true,
			errorMsg:  "is not whitelisted",
		},
		{
			name: "zero amount",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					testChainID: {
						{
							Payee:  testAddr1,
							Token:  testAddr2,
							Amount: big.NewInt(0),
						},
					},
				},
				MCMSConfig: &cldfproposalutils.TimelockConfig{},
			},
			wantError: true,
			errorMsg:  "amount must be positive",
		},
		{
			name: "chain with no transfers",
			config: types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					testChainID: {},
				},
				MCMSConfig: &cldfproposalutils.TimelockConfig{},
			},
			wantError: true,
			errorMsg:  "has no transfers",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateTransferERC20Config(env.GetContext(), *env, tt.config)

			if tt.wantError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					require.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestTransferERC20TimelockQualifierPerChain_lookup(t *testing.T) {
	t.Parallel()

	chain := chainselectors.TEST_90000001.Selector

	tests := []struct {
		name         string
		qualifierMap map[uint64]string
		want         string
	}{
		{
			name:         "nil map defaults to empty qualifier",
			qualifierMap: nil,
			want:         "",
		},
		{
			name:         "empty map defaults to empty qualifier",
			qualifierMap: map[uint64]string{},
			want:         "",
		},
		{
			name: "explicit empty qualifier for chain",
			qualifierMap: map[uint64]string{
				chain: "",
			},
			want: "",
		},
		{
			name: "custom qualifier for chain",
			qualifierMap: map[uint64]string{
				chain: testCustomTimelockQualifier,
			},
			want: testCustomTimelockQualifier,
		},
		{
			name: "other chain in map does not affect lookup",
			qualifierMap: map[uint64]string{
				chain + 1: testCustomTimelockQualifier,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.qualifierMap[chain]
			require.Equal(t, tt.want, got)
		})
	}
}

func TestTransferERC20GetContractAddressWithQualifier(t *testing.T) {
	t.Parallel()

	chain := chainselectors.TEST_90000001.Selector
	ds := datastore.NewMemoryDataStore()
	require.NoError(t, addMCMSContractRefs(ds, chain, "", testTimelockAddrEmptyQualifier, testProposerAddrEmptyQualifier))
	require.NoError(t, addMCMSContractRefs(ds, chain, testCustomTimelockQualifier, testTimelockAddrCustomQualifier, testProposerAddrCustomQualifier))

	env, err := environment.New(t.Context(),
		environment.WithEVMSimulated(t, []uint64{chain}),
		environment.WithDatastore(ds.Seal()),
	)
	require.NoError(t, err)

	tests := []struct {
		name      string
		qualifier string
		wantAddr  string
		wantErr   bool
	}{
		{
			name:      "empty qualifier resolves default timelock",
			qualifier: "",
			wantAddr:  testTimelockAddrEmptyQualifier,
		},
		{
			name:      "custom qualifier resolves custom timelock",
			qualifier: testCustomTimelockQualifier,
			wantAddr:  testTimelockAddrCustomQualifier,
		},
		{
			name:      "unknown qualifier not found",
			qualifier: "does-not-exist",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			addr, err := GetContractAddressWithQualifier(env.DataStore, chain, commontypes.RBACTimelock, tt.qualifier)
			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.qualifier)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.wantAddr, addr)

			proposerAddr, err := GetContractAddressWithQualifier(env.DataStore, chain, commontypes.ProposerManyChainMultisig, tt.qualifier)
			require.NoError(t, err)
			if tt.qualifier == "" {
				require.Equal(t, testProposerAddrEmptyQualifier, proposerAddr)
			} else {
				require.Equal(t, testProposerAddrCustomQualifier, proposerAddr)
			}
		})
	}
}

func TestTransferERC20Validation_timelockQualifier(t *testing.T) {
	t.Parallel()

	chain := chainselectors.TEST_90000001.Selector

	baseTransfer := types.ERC20Transfer{
		Payee:  testAddr1,
		Token:  testAddr2,
		Amount: big.NewInt(1),
	}

	tests := []struct {
		name         string
		seedEmpty    bool
		seedCustom   bool
		qualifierMap map[uint64]string
		wantError    bool
		errorMsg     string
	}{
		{
			name:         "nil qualifier map uses empty qualifier when contracts use empty qualifier",
			seedEmpty:    true,
			qualifierMap: nil,
		},
		{
			name:         "empty qualifier map uses empty qualifier when contracts use empty qualifier",
			seedEmpty:    true,
			qualifierMap: map[uint64]string{},
		},
		{
			name:      "explicit empty qualifier for chain",
			seedEmpty: true,
			qualifierMap: map[uint64]string{
				chain: "",
			},
		},
		{
			name:         "custom qualifier resolves timelock and proposer",
			seedCustom:   true,
			qualifierMap: map[uint64]string{chain: testCustomTimelockQualifier},
		},
		{
			name:         "nil qualifier map fails when only custom qualifier contracts exist",
			seedCustom:   true,
			qualifierMap: nil,
			wantError:    true,
			errorMsg:     "timelock not found",
		},
		{
			name:       "wrong qualifier fails timelock lookup",
			seedCustom: true,
			qualifierMap: map[uint64]string{
				chain: "",
			},
			wantError: true,
			errorMsg:  "timelock not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ds := datastore.NewMemoryDataStore()
			if tt.seedEmpty {
				require.NoError(t, addMCMSContractRefs(ds, chain, "", testTimelockAddrEmptyQualifier, testProposerAddrEmptyQualifier))
			}
			if tt.seedCustom {
				require.NoError(t, addMCMSContractRefs(ds, chain, testCustomTimelockQualifier, testTimelockAddrCustomQualifier, testProposerAddrCustomQualifier))
			}
			require.NoError(t, addWhitelistForChain(ds, chain, testAddr1))

			env, err := environment.New(t.Context(),
				environment.WithEVMSimulated(t, []uint64{chain}),
				environment.WithDatastore(ds.Seal()),
			)
			require.NoError(t, err)

			cfg := types.TransferERC20Config{
				TransfersByChain: map[uint64][]types.ERC20Transfer{
					chain: {baseTransfer},
				},
				MCMSConfig: &cldfproposalutils.TimelockConfig{
					TimelockQualifierPerChain: tt.qualifierMap,
				},
			}

			err = ValidateTransferERC20Config(env.GetContext(), *env, cfg)
			if tt.wantError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorMsg)
				return
			}
			require.NoError(t, err)
		})
	}
}

func addMCMSContractRefs(ds *datastore.MemoryDataStore, chainSelector uint64, qualifier, timelockAddr, proposerAddr string) error {
	if err := ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Address:       timelockAddr,
		Type:          datastore.ContractType(commontypes.RBACTimelock),
		Version:       testMCMSContractVersion,
		Qualifier:     qualifier,
	}); err != nil {
		return err
	}
	return ds.Addresses().Add(datastore.AddressRef{
		ChainSelector: chainSelector,
		Address:       proposerAddr,
		Type:          datastore.ContractType(commontypes.ProposerManyChainMultisig),
		Version:       testMCMSContractVersion,
		Qualifier:     qualifier,
	})
}

func addWhitelistForChain(ds *datastore.MemoryDataStore, chainSelector uint64, addresses ...string) error {
	entries := make([]types.WhitelistAddress, len(addresses))
	for i, addr := range addresses {
		entries[i] = types.WhitelistAddress{
			Address:     addr,
			Description: "test whitelist",
		}
	}
	return ds.ChainMetadata().Upsert(datastore.ChainMetadata{
		ChainSelector: chainSelector,
		Metadata: types.WhitelistMetadata{
			Addresses: entries,
		},
	})
}
