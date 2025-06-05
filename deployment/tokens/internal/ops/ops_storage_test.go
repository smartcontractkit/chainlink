package ops

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

func Test_OpAddAddrBookRecord(t *testing.T) {
	t.Parallel()

	var (
		csel = chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector
		addr = "0x1234567890123456789012345678901234567890"
	)

	tests := []struct {
		name    string
		give    OpAddAddrBookRecordInput
		want    OpAddAddrBookRecordOutput
		wantErr string
	}{
		{
			name: "adds to the address book",
			give: OpAddAddrBookRecordInput{
				ChainSelector: csel,
				Address:       addr,
				Type:          "LinkToken",
				Version:       "1.0.0",
				Labels:        []string{"test", "label"},
			},
			want: OpAddAddrBookRecordOutput{
				ChainSelector:  csel,
				Address:        addr,
				TypeAndVersion: "LinkToken 1.0.0 label test",
			},
		},
		{
			name: "fail: could not save to address book",
			give: OpAddAddrBookRecordInput{
				ChainSelector: 1,
				Address:       addr,
				Type:          "LinkToken",
				Version:       "1.0.0",
				Labels:        []string{"test", "label"},
			},
			wantErr: "chain selector 1: invalid chain selector",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				addrBook = cldf.NewMemoryAddressBook()
				deps     = OpAddAddrBookRecordDeps{AddrBook: addrBook}
			)

			got, err := operations.ExecuteOperation(
				optest.NewBundle(t), OpAddAddrBookRecord, deps, tt.give,
			)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got.Output)

				// Check the address book for the record
				addresses, err := addrBook.AddressesForChain(tt.give.ChainSelector)
				require.NoError(t, err)

				gotVal, ok := addresses[tt.give.Address]
				require.True(t, ok)
				require.Equal(t, got.Output.TypeAndVersion, gotVal.String())
			}
		})
	}
}

func Test_OpAddDatastoreAddrRef(t *testing.T) {
	t.Parallel()

	var (
		csel = chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector
		addr = "0x1234567890123456789012345678901234567890"
	)

	tests := []struct {
		name       string
		beforeFunc func(*testing.T, datastore.MutableDataStore)
		give       OpAddDatastoreAddrRefInput
		want       OpAddDatastoreAddrRefOutput
		wantErr    string
	}{
		{
			name: "adds to the address book",
			give: OpAddDatastoreAddrRefInput{
				ChainSelector: csel,
				Address:       addr,
				Type:          "LinkToken",
				Version:       "1.0.0",
				Labels:        []string{"test", "label"},
				Qualifier:     "test",
			},
			want: OpAddDatastoreAddrRefOutput{
				ChainSelector: csel,
				Address:       addr,
			},
		},
		{
			name: "fail: could not save to address book due to duplicate record key",
			beforeFunc: func(t *testing.T, ds datastore.MutableDataStore) {
				// Pre-populate the datastore with an existing record
				err := ds.Addresses().Add(datastore.AddressRef{
					ChainSelector: csel,
					Address:       addr,
					Type:          "LinkToken",
					Version:       semver.MustParse("1.0.0"),
					Labels:        datastore.NewLabelSet("test", "label"),
					Qualifier:     "test",
				})
				require.NoError(t, err)
			},
			give: OpAddDatastoreAddrRefInput{
				ChainSelector: csel,
				Address:       addr,
				Type:          "LinkToken",
				Version:       "1.0.0",
				Labels:        []string{"test", "label"},
				Qualifier:     "test",
			},
			wantErr: "an address ref with the supplied key already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var (
				ds   = datastore.NewMemoryDataStore()
				deps = OpAddDatastoreAddrRefDeps{
					Datastore: ds,
				}
			)

			if tt.beforeFunc != nil {
				tt.beforeFunc(t, ds)
			}

			got, err := operations.ExecuteOperation(
				optest.NewBundle(t), OpAddDatastoreAddrRef, deps, tt.give,
			)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.EqualError(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, got.Output)

				// Check the datastore for the record
				addrRef, err := ds.Addresses().Get(datastore.NewAddressRefKey(
					tt.give.ChainSelector,
					datastore.ContractType(tt.give.Type),
					semver.MustParse(tt.give.Version),
					tt.give.Qualifier,
				))
				require.NoError(t, err)

				require.Equal(t, tt.give.ChainSelector, addrRef.ChainSelector)
				require.Equal(t, tt.give.Address, addrRef.Address)
				require.Equal(t, tt.give.Type, addrRef.Type.String())
				require.Equal(t, tt.give.Version, addrRef.Version.String())
				require.Equal(t, tt.give.Qualifier, addrRef.Qualifier)
				require.Equal(t, datastore.NewLabelSet(tt.give.Labels...), addrRef.Labels)
			}
		})
	}
}
