package adapters

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	suistate "github.com/smartcontractkit/chainlink-sui/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
)

func TestInitializeSetsCCIPState(t *testing.T) {
	selector := uint64(123)

	ab := stubAddressBook{
		data: map[uint64]map[string]cldf.TypeAndVersion{
			selector: {
				"0xCCIP_PKG":     cldf.NewTypeAndVersion(suistate.SuiCCIPType, *semver.MustParse("1.0.0")),
				"0xCCIP_REF":     cldf.NewTypeAndVersion(suistate.SuiCCIPObjectRefType, *semver.MustParse("1.0.0")),
				"0xCCIP_CAP":     cldf.NewTypeAndVersion(suistate.SuiCCIPOwnerCapObjectIDType, *semver.MustParse("1.0.0")),
				"0xROUTER":       cldf.NewTypeAndVersion(suistate.SuiCCIPRouterType, *semver.MustParse("1.0.0")),
				"0xROUTER_STATE": cldf.NewTypeAndVersion(suistate.SuiCCIPRouterStateObjectType, *semver.MustParse("1.0.0")),
			},
		},
	}

	chain := cldf_sui.Chain{ChainMetadata: cldf_sui.ChainMetadata{Selector: selector}}
	env := cldf.Environment{
		ExistingAddresses: ab,
		BlockChains:       cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{chain}),
	}

	adapter := &CurseAdapter{}
	err := adapter.Initialize(env, selector)
	require.NoError(t, err)
	require.Equal(t, "0xCCIP_PKG", adapter.CCIPAddress)
	require.Equal(t, "0xCCIP_REF", adapter.CCIPObjectRef)
	require.Equal(t, "0xCCIP_CAP", adapter.CCIPOwnerCapObjectID)
	require.Equal(t, "0xROUTER", adapter.RouterAddress)
	require.Equal(t, "0xROUTER_STATE", adapter.RouterStateObjectID)
}

func TestSelectorSubjectConversions(t *testing.T) {
	adapter := &CurseAdapter{}
	selector := uint64(789)
	subject := adapter.SelectorToSubject(selector)
	outSelector, err := adapter.SubjectToSelector(subject)
	require.NoError(t, err)
	require.Equal(t, selector, outSelector)
	require.Equal(t, globals.FamilyAwareSelectorToSubject(selector, chainsel.FamilySui), subject)
}

type stubAddressBook struct {
	data map[uint64]map[string]cldf.TypeAndVersion
}

func (s stubAddressBook) Save(uint64, string, cldf.TypeAndVersion) error { return nil }
func (s stubAddressBook) Addresses() (map[uint64]map[string]cldf.TypeAndVersion, error) {
	return s.data, nil
}
func (s stubAddressBook) AddressesForChain(chain uint64) (map[string]cldf.TypeAndVersion, error) {
	return s.data[chain], nil
}
func (s stubAddressBook) Merge(cldf.AddressBook) error  { return nil }
func (s stubAddressBook) Remove(cldf.AddressBook) error { return nil }
