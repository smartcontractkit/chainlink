package changeset

import (
	"context"
	"errors"
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
)

func TestChangeSetLegacyFunction_PassingCase(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	executedCs := false
	executedValidator := false

	csv2 := cldf.CreateChangeSet(
		func(e deployment.Environment, config uint32) (deployment.ChangesetOutput, error) {
			executedCs = true
			return deployment.ChangesetOutput{AddressBook: deployment.NewMemoryAddressBook()}, nil
		},
		func(e deployment.Environment, config uint32) error {
			executedValidator = true
			return nil
		},
	)
	assert.False(t, executedCs, "Not expected to have executed the changeset yet")
	assert.False(t, executedValidator, "Not expected to have executed the validator yet")
	_, err := Apply(t, e, nil, Configure(csv2, 1))
	assert.True(t, executedCs, "Validator should have returned nil, allowing changeset execution")
	assert.True(t, executedValidator, "Not expected to have executed the validator yet")
	assert.NoError(t, err)
}

func TestChangeSetLegacyFunction_ErrorCase(t *testing.T) {
	t.Parallel()
	e := NewNoopEnvironment(t)

	executedCs := false
	executedValidator := false

	csv2 := cldf.CreateChangeSet(
		func(e deployment.Environment, config uint32) (deployment.ChangesetOutput, error) {
			executedCs = true
			return deployment.ChangesetOutput{AddressBook: deployment.NewMemoryAddressBook()}, nil
		},
		func(e deployment.Environment, config uint32) error {
			executedValidator = true
			return errors.New("you shall not pass")
		},
	)
	assert.False(t, executedCs, "Not expected to have executed the changeset yet")
	assert.False(t, executedValidator, "Not expected to have executed the validator yet")
	_, err := Apply(t, e, nil, Configure(csv2, 1))
	assert.False(t, executedCs, "Validator should have fired, preventing changeset execution")
	assert.True(t, executedValidator, "Not expected to have executed the validator yet")
	assert.Equal(t, "failed to apply changeset at index 0: you shall not pass", err.Error())
}

func NewNoopEnvironment(t *testing.T) deployment.Environment {
	return *deployment.NewEnvironment(
		"noop",
		logger.TestLogger(t),
		deployment.NewMemoryAddressBook(),
		datastore.NewMemoryDataStore[
			datastore.DefaultMetadata,
			datastore.DefaultMetadata,
		]().Seal(),
		map[uint64]deployment.Chain{},
		map[uint64]deployment.SolChain{},
		map[uint64]deployment.AptosChain{},
		[]string{},
		nil,
		func() context.Context { return tests.Context(t) },
		cldf.XXXGenerateTestOCRSecrets(),
	)
}

func TestApplyChangesetsHelpers(t *testing.T) {
	t.Parallel()

	changesets := []ConfiguredChangeSet{
		Configure(cldf.CreateChangeSet(
			func(e deployment.Environment, config uint32) (deployment.ChangesetOutput, error) {
				ds := datastore.NewMemoryDataStore[
					datastore.DefaultMetadata,
					datastore.DefaultMetadata,
				]()

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
					return deployment.ChangesetOutput{}, err
				}

				// Add ContractMetadata
				err := ds.ContractMetadataStore.Upsert(datastore.ContractMetadata[datastore.DefaultMetadata]{
					ChainSelector: 1,
					Address:       "0x1234567890abcdef",
					Metadata:      datastore.DefaultMetadata{Data: "test"},
				})
				if err != nil {
					return deployment.ChangesetOutput{}, err
				}

				return deployment.ChangesetOutput{
					AddressBook: deployment.NewMemoryAddressBook(),
					DataStore:   ds,
				}, nil
			},
			func(e deployment.Environment, config uint32) error {
				return nil
			},
		), 1),
	}

	csTests := []struct {
		name                   string
		changesets             []ConfiguredChangeSet
		validate               func(t *testing.T, e deployment.Environment)
		wantError              bool
		changesetApplyFunction string
	}{
		{
			name:                   "ApplyChangesets validates datastore is merged after apply",
			changesets:             changesets,
			changesetApplyFunction: "V2",
			validate: func(t *testing.T, e deployment.Environment) {
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
				assert.Equal(t, "0x1234567890abcdef", record.Address)

				// Check metadata was stored correctly
				metadata, err := e.DataStore.ContractMetadata().Get(
					datastore.NewContractMetadataKey(1, "0x1234567890abcdef"),
				)
				require.NoError(t, err)
				assert.Equal(t, "test", metadata.Metadata.Data)
			},
			wantError: false,
		},
		{
			name:                   "ApplyChangesetsV2 validates datastore is merged after apply",
			changesets:             changesets,
			changesetApplyFunction: "V1",
			validate: func(t *testing.T, e deployment.Environment) {
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
				assert.Equal(t, "0x1234567890abcdef", record.Address)

				// Check metadata was stored correctly
				metadata, err := e.DataStore.ContractMetadata().Get(
					datastore.NewContractMetadataKey(1, "0x1234567890abcdef"),
				)
				require.NoError(t, err)
				assert.Equal(t, "test", metadata.Metadata.Data)
			},
			wantError: false,
		},
	}

	for _, tt := range csTests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.changesetApplyFunction {
			case "V2":
				e := NewNoopEnvironment(t)
				e, _, err := ApplyChangesetsV2(t, e, tt.changesets)
				if tt.wantError {
					require.Error(t, err)
					return
				}
				require.NoError(t, err)
				tt.validate(t, e)
			case "V1":
				e := NewNoopEnvironment(t)
				e, err := ApplyChangesets(t, e, nil, tt.changesets)
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
