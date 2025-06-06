package seqs

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations/optest"
)

func Test_PersistAddress(t *testing.T) {
	t.Parallel()

	var (
		chainSelector   = chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector
		addr            = "0xeC91988D7dD84d8adE801b739172ad15c860A700"
		contractType    = "SomeContract"
		contractVersion = "1.0.0"
		qualifier       = "test"
		labels          = []string{"label1", "label2"}
	)

	tests := []struct {
		name       string
		beforeFunc func(*testing.T, datastore.MutableDataStore)
		give       SeqPersistAddressInput
		want       SeqPersistAddressOutput
		wantErr    string
	}{
		{
			name: "adds to address book and data store",
			give: SeqPersistAddressInput{
				ChainSelector: chainSelector,
				Address:       addr,
				Type:          contractType,
				Version:       contractVersion,
				Qualifier:     qualifier,
				Labels:        labels,
			},
			want: SeqPersistAddressOutput{},
		},
		{
			name: "err: cannot save to address book",
			give: SeqPersistAddressInput{
				ChainSelector: 1, // invalid chain selector
				Address:       addr,
				Type:          contractType,
				Version:       contractVersion,
				Qualifier:     qualifier,
				Labels:        labels,
			},
			wantErr: "invalid chain selector",
		},
		{
			name: "err: cannot save to data store",
			beforeFunc: func(t *testing.T, ds datastore.MutableDataStore) {
				// Pre-populate the datastore with an existing record
				err := ds.Addresses().Add(datastore.AddressRef{
					ChainSelector: chainSelector,
					Address:       addr,
					Type:          datastore.ContractType(contractType),
					Version:       semver.MustParse(contractVersion),
					Qualifier:     qualifier,
					Labels:        datastore.NewLabelSet(labels...),
				})
				require.NoError(t, err)
			},
			give: SeqPersistAddressInput{
				ChainSelector: chainSelector,
				Address:       addr,
				Type:          contractType,
				Version:       contractVersion,
				Qualifier:     qualifier,
				Labels:        labels,
			},
			wantErr: "an address ref with the supplied key already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				addrBook = cldf.NewMemoryAddressBook()
				ds       = datastore.NewMemoryDataStore()
				deps     = SeqPersistAddressDeps{
					AddrBook:  addrBook,
					Datastore: ds,
				}
			)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t, ds)
			}

			got, err := operations.ExecuteSequence(
				optest.NewBundle(t), SeqPersistAddress, deps, tt.give,
			)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got.Output)

				// Check that the address book has the link token contract for each chain
				addrBookByChain, err := addrBook.AddressesForChain(tt.give.ChainSelector)
				require.NoError(t, err)
				require.NotEmpty(t, addrBookByChain)
				require.Len(t, addrBookByChain, 1)

				// Check that the address reference is in the datastore
				addrRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(
					tt.give.ChainSelector,
					datastore.ContractType(tt.give.Type),
					semver.MustParse(tt.give.Version),
					tt.give.Qualifier,
				))
				require.NoError(t, err)
				require.Equal(t, tt.give.Address, addrRef.Address)
			}
		})
	}
}
