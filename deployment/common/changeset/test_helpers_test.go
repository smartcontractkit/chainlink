package changeset

import (
	"errors"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestChangeSetLegacyFunction_PassingCase(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	executedCs := false
	executedValidator := false

	csv2 := cldf.CreateChangeSet(
		func(e cldf.Environment, config uint32) (cldf.ChangesetOutput, error) {
			executedCs = true
			return cldf.ChangesetOutput{AddressBook: cldf.NewMemoryAddressBook()}, nil
		},
		func(e cldf.Environment, config uint32) error {
			executedValidator = true
			return nil
		},
	)
	require.False(t, executedCs, "Not expected to have executed the changeset yet")
	require.False(t, executedValidator, "Not expected to have executed the validator yet")
	_, err := Apply(t, e, Configure(csv2, 1))
	require.True(t, executedCs, "Validator should have returned nil, allowing changeset execution")
	require.True(t, executedValidator, "Not expected to have executed the validator yet")
	require.NoError(t, err)
}

func TestChangeSetLegacyFunction_ErrorCase(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	executedCs := false
	executedValidator := false

	csv2 := cldf.CreateChangeSet(
		func(e cldf.Environment, config uint32) (cldf.ChangesetOutput, error) {
			executedCs = true
			return cldf.ChangesetOutput{AddressBook: cldf.NewMemoryAddressBook()}, nil
		},
		func(e cldf.Environment, config uint32) error {
			executedValidator = true
			return errors.New("you shall not pass")
		},
	)
	require.False(t, executedCs, "Not expected to have executed the changeset yet")
	require.False(t, executedValidator, "Not expected to have executed the validator yet")
	_, err := Apply(t, e, Configure(csv2, 1))
	require.False(t, executedCs, "Validator should have fired, preventing changeset execution")
	require.True(t, executedValidator, "Not expected to have executed the validator yet")
	require.Equal(t, "failed to apply changeset at index 0: you shall not pass", err.Error())
}

func NewNoopEnvironment(t *testing.T) cldf.Environment {
	return *cldf.NewEnvironment(
		"noop",
		logger.TestLogger(t),
		cldf.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore().Seal(),
		[]string{},
		nil,
		t.Context,
		cldf.XXXGenerateTestOCRSecrets(),
		chain.NewBlockChains(map[uint64]chain.BlockChain{}),
	)
}

func TestApplyChangesetsHelpers(t *testing.T) {
	t.Parallel()

	changesets := []ConfiguredChangeSet{
		Configure(cldf.CreateChangeSet(
			func(e cldf.Environment, config uint32) (cldf.ChangesetOutput, error) {
				ds := datastore.NewMemoryDataStore()

				// Store Address
				if err := ds.Addresses().Add(
					datastore.AddressRef{
						ChainSelector: 1,
						Address:       "0x1234567890abcdef",
						Type:          "TEST_CONTRACT",
						Version:       semver.MustParse("1.0.0"),
						Qualifier:     "qualifier1",
					},
				); err != nil {
					return cldf.ChangesetOutput{}, err
				}

				// Add ContractMetadata
				err := ds.ContractMetadataStore.Upsert(datastore.ContractMetadata{
					ChainSelector: 1,
					Address:       "0x1234567890abcdef",
					Metadata:      testMetadata{Data: "test"},
				})
				if err != nil {
					return cldf.ChangesetOutput{}, err
				}

				return cldf.ChangesetOutput{
					AddressBook: cldf.NewMemoryAddressBook(),
					DataStore:   ds,
				}, nil
			},
			func(e cldf.Environment, config uint32) error {
				return nil
			},
		), 1),
	}

	csTests := []struct {
		name                   string
		changesets             []ConfiguredChangeSet
		validate               func(t *testing.T, e cldf.Environment)
		wantError              bool
		changesetApplyFunction string
	}{
		{
			name:                   "ApplyChangesets validates datastore is merged after apply",
			changesets:             changesets,
			changesetApplyFunction: "V2",
			validate: func(t *testing.T, e cldf.Environment) {
				// Check address was stored correctly
				record, err := e.DataStore.Addresses().Get(
					datastore.NewAddressRefKey(
						1,
						"TEST_CONTRACT",
						semver.MustParse("1.0.0"),
						"qualifier1",
					),
				)
				require.NoError(t, err)
				require.Equal(t, "0x1234567890abcdef", record.Address)

				// Check metadata was stored correctly
				metadata, err := e.DataStore.ContractMetadata().Get(
					datastore.NewContractMetadataKey(1, "0x1234567890abcdef"),
				)
				require.NoError(t, err)
				concrete, err := datastore.As[testMetadata](metadata.Metadata)
				require.NoError(t, err)
				require.Equal(t, "test", concrete.Data)
			},
			wantError: false,
		},
		{
			name:                   "ApplyChangesets validates datastore is merged after apply",
			changesets:             changesets,
			changesetApplyFunction: "V1",
			validate: func(t *testing.T, e cldf.Environment) {
				// Check address was stored correctly
				record, err := e.DataStore.Addresses().Get(
					datastore.NewAddressRefKey(
						1,
						"TEST_CONTRACT",
						semver.MustParse("1.0.0"),
						"qualifier1",
					),
				)
				require.NoError(t, err)
				require.Equal(t, "0x1234567890abcdef", record.Address)

				// Check metadata was stored correctly
				metadata, err := e.DataStore.ContractMetadata().Get(
					datastore.NewContractMetadataKey(1, "0x1234567890abcdef"),
				)
				require.NoError(t, err)
				concrete, err := datastore.As[testMetadata](metadata.Metadata)
				require.NoError(t, err)
				require.Equal(t, "test", concrete.Data)
			},
			wantError: false,
		},
	}

	for _, tt := range csTests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.changesetApplyFunction {
			case "V2":
				e := NewNoopEnvironment(t)
				e, _, err := ApplyChangesets(t, e, tt.changesets)
				if tt.wantError {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				tt.validate(t, e)
			case "V1":
				e := NewNoopEnvironment(t)
				e, _, err := ApplyChangesets(t, e, tt.changesets)
				if tt.wantError {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				tt.validate(t, e)
			default:
				t.Fatalf("unknown changeset apply function: %s", tt.changesetApplyFunction)
			}
		})
	}
}
