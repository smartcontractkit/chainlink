package deployment

import (
	"math/big"
	"sync"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/ethereum/go-ethereum/common"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddressBook_Save(t *testing.T) {
	ab := NewMemoryAddressBook()
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()

	err := ab.Save(chainsel.TEST_90000001.Selector, addr1, onRamp100)
	require.NoError(t, err)

	// Invalid address
	err = ab.Save(chainsel.TEST_90000001.Selector, "asdlfkj", onRamp100)
	require.ErrorIs(t, err, ErrInvalidAddress)

	// Valid chain but not present.
	_, err = ab.AddressesForChain(chainsel.TEST_90000002.Selector)
	require.ErrorIs(t, err, ErrChainNotFound)

	// Invalid selector
	err = ab.Save(0, addr1, onRamp100)
	require.ErrorIs(t, err, ErrInvalidChainSelector)

	// Duplicate
	err = ab.Save(chainsel.TEST_90000001.Selector, addr1, onRamp100)
	require.Error(t, err)

	// Zero address
	err = ab.Save(chainsel.TEST_90000001.Selector, common.HexToAddress("0x0").Hex(), onRamp100)
	require.Error(t, err)

	// Zero address but non evm chain
	err = NewMemoryAddressBook().Save(chainsel.APTOS_MAINNET.Selector, common.HexToAddress("0x0").Hex(), onRamp100)
	require.NoError(t, err)

	// Distinct address same TV will not
	err = ab.Save(chainsel.TEST_90000001.Selector, addr2, onRamp100)
	require.NoError(t, err)
	// Same address different chain will not error
	err = ab.Save(chainsel.TEST_90000002.Selector, addr1, onRamp100)
	require.NoError(t, err)
	// We can save different versions of the same contract
	err = ab.Save(chainsel.TEST_90000002.Selector, addr2, onRamp110)
	require.NoError(t, err)

	addresses, err := ab.Addresses()
	require.NoError(t, err)
	assert.Equal(t, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp100,
			addr2: onRamp110,
		},
	}, addresses)
}

func TestAddressBook_Merge(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	a1 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
	})
	a2 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	})
	require.NoError(t, a1.Merge(a2))

	addresses, err := a1.Addresses()
	require.NoError(t, err)
	assert.Equal(t, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	}, addresses)

	// Merge with conflicting addresses should error
	a3 := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
	})
	require.Error(t, a1.Merge(a3))
	// a1 should not have changed
	addresses, err = a1.Addresses()
	require.NoError(t, err)
	assert.Equal(t, map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
		},
	}, addresses)
}

func TestAddressBook_Remove(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	addr3 := common.HexToAddress("0x3").String()

	baseAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr2: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110,
			addr3: onRamp110,
		},
	})

	copyOfBaseAB := NewMemoryAddressBookFromMap(baseAB.cloneAddresses(baseAB.addressesByChain))

	// this address book shouldn't be removed (state of baseAB not changed, error thrown)
	failAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
			addr3: onRamp100, // doesn't exist in TEST_90000001.Selector
		},
	})
	require.Error(t, baseAB.Remove(failAB))
	require.EqualValues(t, baseAB, copyOfBaseAB)

	// this Address book should be removed without error
	successAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000002.Selector: {
			addr3: onRamp100,
		},
		chainsel.TEST_90000001.Selector: {
			addr2: onRamp100,
		},
	})

	expectingAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			addr1: onRamp100,
		},
		chainsel.TEST_90000002.Selector: {
			addr1: onRamp110},
	})

	require.NoError(t, baseAB.Remove(successAB))
	require.EqualValues(t, baseAB, expectingAB)
}

func TestAddressBook_ConcurrencyAndDeadlock(t *testing.T) {
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)

	baseAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
		chainsel.TEST_90000001.Selector: {
			common.BigToAddress(big.NewInt(1)).String(): onRamp100,
		},
	})

	// concurrent writes
	var i int64
	wg := sync.WaitGroup{}
	for i = 2; i < 1000; i++ {
		wg.Add(1)
		go func(input int64) {
			assert.NoError(t, baseAB.Save(
				chainsel.TEST_90000001.Selector,
				common.BigToAddress(big.NewInt(input)).String(),
				onRamp100,
			))
			wg.Done()
		}(i)
	}

	// concurrent reads
	for i = 0; i < 100; i++ {
		wg.Add(1)
		go func(input int64) {
			addresses, err := baseAB.Addresses()
			if !assert.NoError(t, err) {
				return
			}
			for chainSelector, chainAddresses := range addresses {
				// concurrent read chainAddresses from Addresses() method
				for address := range chainAddresses {
					addresses[chainSelector][address] = onRamp110
				}

				// concurrent read chainAddresses from AddressesForChain() method
				chainAddresses, err = baseAB.AddressesForChain(chainSelector)
				if assert.NoError(t, err) {
					for address := range chainAddresses {
						_ = addresses[chainSelector][address]
					}
				}
			}
			wg.Done()
		}(i)
	}

	// concurrent merges, starts from 1001 to avoid address conflicts
	for i = 1001; i < 1100; i++ {
		wg.Add(1)
		go func(input int64) {
			// concurrent merge
			additionalAB := NewMemoryAddressBookFromMap(map[uint64]map[string]TypeAndVersion{
				chainsel.TEST_90000002.Selector: {
					common.BigToAddress(big.NewInt(input)).String(): onRamp100,
				},
			})
			assert.NoError(t, baseAB.Merge(additionalAB))
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func TestAddressesContainBundle(t *testing.T) {
	t.Parallel()

	// Define some TypeAndVersion values
	onRamp100 := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp110 := NewTypeAndVersion("OnRamp", Version1_1_0)
	onRamp120 := NewTypeAndVersion("OnRamp", Version1_2_0)

	// Create one with labels
	onRamp100WithLabels := NewTypeAndVersion("OnRamp", Version1_0_0)
	onRamp100WithLabels.Labels.Add("sa")
	onRamp100WithLabels.Labels.Add("staging")

	addr1 := common.HexToAddress("0x1").String()
	addr2 := common.HexToAddress("0x2").String()
	addr3 := common.HexToAddress("0x3").String()

	tests := []struct {
		name       string
		addrs      map[string]TypeAndVersion // input address map
		wantTypes  []TypeAndVersion          // the “bundle” we want
		wantErr    bool
		wantErrMsg string
		wantResult bool // expected boolean return when no error
	}{
		{
			name: "More than one instance => error",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100,
				addr2: onRamp100, // duplicate
			},
			wantTypes: []TypeAndVersion{onRamp100},
			wantErr:   true,
			// an example substring check:
			wantErrMsg: "found more than one instance of contract",
		},
		{
			name: "No instance => result false, no error",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp110,
				addr2: onRamp110,
			},
			wantTypes:  []TypeAndVersion{onRamp100},
			wantErr:    false,
			wantResult: false,
		},
		{
			name: "2 elements => success",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100,
				addr2: onRamp110,
				addr3: onRamp120,
			},
			wantTypes:  []TypeAndVersion{onRamp100, onRamp110},
			wantErr:    false,
			wantResult: true,
		},
		{
			name: "Mismatched labels => false",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100, // no labels
			},
			wantTypes:  []TypeAndVersion{onRamp100WithLabels},
			wantErr:    false,
			wantResult: false, // label mismatch => not found
		},
		{
			name: "Exact label match => success",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100WithLabels,
			},
			wantTypes:  []TypeAndVersion{onRamp100WithLabels},
			wantErr:    false,
			wantResult: true,
		},
		{
			name: "Duplicate labeled => error",
			addrs: map[string]TypeAndVersion{
				addr1: onRamp100WithLabels,
				addr2: onRamp100WithLabels, // same type/version/labels => duplicate
			},
			wantTypes:  []TypeAndVersion{onRamp100WithLabels},
			wantErr:    true,
			wantErrMsg: "more than one instance of contract",
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotResult, gotErr := AddressesContainBundle(tt.addrs, tt.wantTypes)

			if tt.wantErr {
				require.Error(t, gotErr, "expected an error but got none")
				if tt.wantErrMsg != "" {
					require.Contains(t, gotErr.Error(), tt.wantErrMsg)
				}
				return
			}
			require.NoError(t, gotErr, "did not expect an error but got one")
			assert.Equal(t, tt.wantResult, gotResult,
				"expected result %v but got %v", tt.wantResult, gotResult)
		})
	}
}

func TestTypeAndVersionFromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		input              string
		wantErr            bool
		wantType           ContractType
		wantVersion        semver.Version
		wantLabels         LabelSet
		wantTypeAndVersion string
	}{
		{
			name:               "valid - no labels",
			input:              "CallProxy 1.0.0",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantLabels:         NewLabelSet(),
			wantTypeAndVersion: "CallProxy 1.0.0",
		},
		{
			name:               "valid - multiple labels, normal spacing",
			input:              "CallProxy 1.0.0 SA staging",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantLabels:         NewLabelSet("SA", "staging"),
			wantTypeAndVersion: "CallProxy 1.0.0 SA staging",
		},
		{
			name:               "valid - multiple labels, extra spacing",
			input:              "   CallProxy     1.0.0    SA    staging   ",
			wantErr:            false,
			wantType:           "CallProxy",
			wantVersion:        Version1_0_0,
			wantLabels:         NewLabelSet("SA", "staging"),
			wantTypeAndVersion: "CallProxy 1.0.0 SA staging",
		},
		{
			name:    "invalid - not enough parts",
			input:   "CallProxy",
			wantErr: true,
		},
		{
			name:    "invalid - version not parseable",
			input:   "CallProxy notASemver",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotTV, gotErr := TypeAndVersionFromString(tt.input)
			if tt.wantErr {
				require.Error(t, gotErr, "expected error but got none")
				return
			}
			require.NoError(t, gotErr, "did not expect an error but got one")

			// Check ContractType
			require.Equal(t, tt.wantType, gotTV.Type, "incorrect contract type")

			// Check Version
			require.Equal(t, tt.wantVersion.String(), gotTV.Version.String(), "incorrect version")

			// Check labels
			require.Equal(t, tt.wantLabels, gotTV.Labels, "labels mismatch")

			// Check type and version
			require.Equal(t, tt.wantTypeAndVersion, gotTV.String(), "type and version mismatch")
		})
	}
}

func TestTypeAndVersion_AddLabels(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		initialLabels []string
		toAdd         []string
		wantContains  []string
		wantLen       int
	}{
		{
			name:          "add single labels to empty set",
			initialLabels: nil,
			toAdd:         []string{"foo"},
			wantContains:  []string{"foo"},
			wantLen:       1,
		},
		{
			name:          "add multiple labels to existing set",
			initialLabels: []string{"alpha"},
			toAdd:         []string{"beta", "gamma"},
			wantContains:  []string{"alpha", "beta", "gamma"},
			wantLen:       3,
		},
		{
			name:          "add duplicate labels",
			initialLabels: []string{"dup"},
			toAdd:         []string{"dup", "dup", "new"},
			wantContains:  []string{"dup", "new"},
			wantLen:       2,
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Construct a TypeAndVersion with any initial labels
			tv := TypeAndVersion{
				Type:    "CallProxy",
				Version: Version1_0_0,
				Labels:  NewLabelSet(tt.initialLabels...),
			}

			// Call AddLabel for each item in toAdd
			for _, label := range tt.toAdd {
				tv.AddLabel(label)
			}

			// Check final labels length
			require.Len(t, tv.Labels, tt.wantLen, "labels size mismatch")

			// Check that expected labels is present
			for _, md := range tt.wantContains {
				require.True(t, tv.Labels.Contains(md),
					"expected labels %q was not found in tv.Labels", md)
			}
		})
	}
}
