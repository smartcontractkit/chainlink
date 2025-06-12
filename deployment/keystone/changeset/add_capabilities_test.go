package changeset_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	kcr "github.com/smartcontractkit/chainlink-evm/gethwrappers/keystone/generated/capabilities_registry_1_1_0"

	commonchangeset "github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/deployment/keystone/changeset/test"
)

func TestAddCapabilities(t *testing.T) {
	t.Parallel()

	capabilitiesToAdd := []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:   "test-cap",
			Version:        "0.0.1",
			CapabilityType: 1,
		},
		{
			LabelledName:   "test-cap-2",
			Version:        "0.0.1",
			CapabilityType: 1,
		},
	}
	t.Run("no mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
		})

		csOut, err := changeset.AddCapabilities(te.Env, &changeset.AddCapabilitiesRequest{
			RegistryChainSel: te.RegistrySelector,
			Capabilities:     capabilitiesToAdd,
			RegistryRef:      te.CapabilityRegistryAddressRef(),
		})
		require.NoError(t, err)
		require.Empty(t, csOut.MCMSTimelockProposals)
		require.Nil(t, csOut.AddressBook)
		assertCapabilitiesExist(t, te.CapabilitiesRegistry(), capabilitiesToAdd...)
	})

	t.Run("with mcms", func(t *testing.T) {
		te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
			WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
			AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
			WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
			NumChains:       1,
			UseMCMS:         true,
		})

		req := &changeset.AddCapabilitiesRequest{
			RegistryChainSel: te.RegistrySelector,
			Capabilities:     capabilitiesToAdd,
			MCMSConfig:       &changeset.MCMSConfig{MinDuration: 0},
			RegistryRef:      te.CapabilityRegistryAddressRef(),
		}
		csOut, err := changeset.AddCapabilities(te.Env, req)
		require.NoError(t, err)
		require.Len(t, csOut.MCMSTimelockProposals, 1)
		require.Nil(t, csOut.AddressBook)

		// now apply the changeset such that the proposal is signed and execed
		err = applyProposal(t, te, commonchangeset.Configure(cldf.CreateLegacyChangeSet(changeset.AddCapabilities), req))
		require.NoError(t, err)

		assertCapabilitiesExist(t, te.CapabilitiesRegistry(), capabilitiesToAdd...)
	})
}

func TestAddCapabilitiesRequest_Validate_WriterCapability(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		req           func(wrapper test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error)
		expectedError error
	}{
		{
			name: "valid request with chain ID on capability name",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				chainID, err := chainselectors.GetChainIDFromSelector(chainselectors.TEST_90000001.Selector)
				if err != nil {
					return nil, err
				}
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: fmt.Sprintf("%s%s", changeset.CapabilityTypeTargetNamePrefix, chainID), Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: nil,
		},
		// Cannot test this since `chainselectors.ChainIdFromName()` uses the chains from the `.yaml` files,
		// and the chain name is not set in the `test_selectors.yaml` file.
		// {
		//	name: "valid request with chain name on capability name",
		//	req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
		//		chain := te.Env.BlockChains.EVMChains()[chainselectors.TEST_90000001.Selector]
		//		return &changeset.AddCapabilitiesRequest{
		//			RegistryChainSel: te.RegistrySelector,
		//			Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: fmt.Sprintf("%s%s", changeset.CapabilityTypeTargetNamePrefix, chain.Name()), Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
		//			RegistryRef:      te.CapabilityRegistryAddressRef(),
		//		}, nil
		//	},
		//	expectError: false,
		// },
		{
			name: "empty capability name",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: "", Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: changeset.ErrEmptyWriteCapName,
		},
		{
			name: "only has prefix on capability name",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: changeset.CapabilityTypeTargetNamePrefix, Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: changeset.ErrEmptyTrimmedWriteCapName,
		},
		{
			name: "missing prefix on capability name",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: "test-cap", Version: "1.0.0", CapabilityType: 3}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: changeset.ErrInvalidWriteCapName,
		},
		{
			name: "invalid chain name on capability name",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: changeset.CapabilityTypeTargetNamePrefix + "test-cap", Version: "1.0.0", CapabilityType: 3}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: changeset.ErrInvalidWriteCapNameFormat,
		},
		{
			name: "invalid chain ID on capability name",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: changeset.CapabilityTypeTargetNamePrefix + "12345", Version: "1.0.0", CapabilityType: 3}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: changeset.ErrInvalidWriteCapNameFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			te := test.SetupContractTestEnv(t, test.EnvWrapperConfig{
				WFDonConfig:     test.DonConfig{Name: "wfDon", N: 4},
				AssetDonConfig:  test.DonConfig{Name: "assetDon", N: 4},
				WriterDonConfig: test.DonConfig{Name: "writerDon", N: 4},
				NumChains:       1,
				UseMCMS:         true,
			})

			req, err := tt.req(te)
			require.NoError(t, err)
			err = req.Validate(te.Env)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func assertCapabilitiesExist(t *testing.T, registry *kcr.CapabilitiesRegistry, capabilities ...kcr.CapabilitiesRegistryCapability) {
	for _, capability := range capabilities {
		wantID, err := registry.GetHashedCapabilityId(nil, capability.LabelledName, capability.Version)
		require.NoError(t, err)
		got, err := registry.GetCapability(nil, wantID)
		require.NoError(t, err)
		require.NotEmpty(t, got)
		assert.Equal(t, capability.CapabilityType, got.CapabilityType)
		assert.Equal(t, capability.LabelledName, got.LabelledName)
		assert.Equal(t, capability.Version, got.Version)
	}
}
