package changesets

import (
	"fmt"
	"testing"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/tokens/internal/ops"
)

func Test_validateNoDupeSelectors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		give    []uint64
		wantErr string
	}{
		{
			name: "valid input",
			give: []uint64{
				1, 2, 3,
			},
		},
		{
			name: "duplicate input",
			give: []uint64{
				1, 2, 3, 1,
			},
			wantErr: fmt.Sprintf("duplicate chain selector found: %d", 1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateNoDupeSelectors(tt.give)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_validateNoExistingLinkToken(t *testing.T) {
	var (
		csel    = chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector
		ethAddr = "0xeC91988D7dD84d8adE801b739172ad15c860A700"
	)

	tests := []struct {
		name          string
		beforeFunc    func(*testing.T, *cldf.AddressBookMap, *datastore.MemoryDataStore)
		give          []uint64
		giveQualifier string
		wantErr       string
	}{
		{
			name: "no existing link token found",
			give: []uint64{csel},
		},
		{
			name: "no existing link token found with qualifier",
			beforeFunc: func(
				t *testing.T,
				_ *cldf.AddressBookMap,
				ds *datastore.MemoryDataStore,
			) {
				err := ds.Addresses().Add(datastore.AddressRef{
					ChainSelector: csel,
					Address:       ethAddr,
					Qualifier:     "test",
					Type:          datastore.ContractType(ops.LinkTokenTypeAndVersion1.Type.String()),
					Version:       &ops.LinkTokenTypeAndVersion1.Version,
				})
				require.NoError(t, err)
			},
			give: []uint64{csel},
		},
		{
			name: "existing in address book",
			beforeFunc: func(
				t *testing.T,
				ab *cldf.AddressBookMap,
				_ *datastore.MemoryDataStore,
			) {
				err := ab.Save(csel, ethAddr, ops.LinkTokenTypeAndVersion1)
				require.NoError(t, err)
			},
			give: []uint64{csel},
			wantErr: fmt.Sprintf(
				"link token contract already exists for chain selector %d in address book", csel,
			),
		},
		{
			name: "existing in datastore",
			beforeFunc: func(
				t *testing.T,
				_ *cldf.AddressBookMap,
				ds *datastore.MemoryDataStore,
			) {
				err := ds.Addresses().Add(datastore.AddressRef{
					ChainSelector: csel,
					Address:       ethAddr,
					Qualifier:     "test",
					Type:          datastore.ContractType(ops.LinkTokenTypeAndVersion1.Type.String()),
					Version:       &ops.LinkTokenTypeAndVersion1.Version,
				})
				require.NoError(t, err)
			},
			give:          []uint64{csel},
			giveQualifier: "test",
			wantErr: fmt.Sprintf(
				"link token contract already exists for chain selector %d in datastore", csel,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ab := cldf.NewMemoryAddressBook()
			ds := datastore.NewMemoryDataStore()

			if tt.beforeFunc != nil {
				tt.beforeFunc(t, ab, ds)
			}

			err := validateNoExistingLinkToken(tt.give, tt.giveQualifier, ab, ds.Seal())
			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func Test_validateChainSelectorsFamily(t *testing.T) {
	t.Parallel()

	var (
		csel = chain_selectors.ETHEREUM_TESTNET_SEPOLIA.Selector
	)

	tests := []struct {
		name    string
		give    []uint64
		giveFam string
		wantErr string
	}{
		{
			name:    "valid input",
			give:    []uint64{csel},
			giveFam: "evm",
		},
		{
			name:    "invalid input",
			give:    []uint64{csel},
			giveFam: "solana",
			wantErr: fmt.Sprintf("chain selector %d is not in the solana family", csel),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := validateChainSelectorsFamily(tt.give, tt.giveFam)

			if tt.wantErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
