package adapters

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/stretchr/testify/require"

	chainsel "github.com/smartcontractkit/chain-selectors"

	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_aptos "github.com/smartcontractkit/chainlink-deployments-framework/chain/aptos"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/globals"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"
)

func TestInitializeSetsCCIPAddress(t *testing.T) {
	selector := uint64(123)
	ccipAddrStr := "0x1"
	ccipAddr := aptos.AccountAddress{}
	require.NoError(t, ccipAddr.ParseStringRelaxed(ccipAddrStr))

	ab := stubAddressBook{
		data: map[uint64]map[string]cldf.TypeAndVersion{
			selector: {
				ccipAddrStr: cldf.NewTypeAndVersion(shared.AptosCCIPType, *semver.MustParse("1.0.0")),
			},
		},
	}

	chain := cldf_aptos.Chain{Selector: selector}
	env := cldf.Environment{
		ExistingAddresses: ab,
		BlockChains:       cldf_chain.NewBlockChainsFromSlice([]cldf_chain.BlockChain{chain}),
	}

	adapter := &CurseAdapter{}
	err := adapter.Initialize(env, selector)
	require.NoError(t, err)
	require.Equal(t, ccipAddr, adapter.CCIPAddress)
}

func TestSelectorSubjectConversions(t *testing.T) {
	adapter := &CurseAdapter{}
	selector := uint64(789)
	subject := adapter.SelectorToSubject(selector)
	outSelector, err := adapter.SubjectToSelector(subject)
	require.NoError(t, err)
	require.Equal(t, selector, outSelector)
	require.Equal(t, globals.FamilyAwareSelectorToSubject(selector, chainsel.FamilyAptos), subject)
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
