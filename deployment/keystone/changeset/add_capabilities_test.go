package changeset_test

import (
	"errors"
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
			name: "valid request with chain ID on capability name and `writer_` prefix",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				chainID, err := chainselectors.GetChainIDFromSelector(chainselectors.TEST_90000001.Selector)
				if err != nil {
					return nil, err
				}
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: fmt.Sprintf("%s%s", changeset.CapabilityTypeTargetNamePrefix1, chainID), Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "valid request with chain ID on capability name and `writer-` prefix",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				chainID, err := chainselectors.GetChainIDFromSelector(chainselectors.TEST_90000001.Selector)
				if err != nil {
					return nil, err
				}
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: fmt.Sprintf("%s%s", changeset.CapabilityTypeTargetNamePrefix2, chainID), Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "valid request with chain name on capability name and `writer_` prefix",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				chainName := "random-chain-name"
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: fmt.Sprintf("%s%s", changeset.CapabilityTypeTargetNamePrefix1, chainName), Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "valid request with chain name on capability name and `writer-` prefix",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				chainName := "random-chain-name-1"
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: fmt.Sprintf("%s%s", changeset.CapabilityTypeTargetNamePrefix2, chainName), Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "valid request with suffix like `:region` on capability name",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: changeset.CapabilityTypeTargetNamePrefix1 + "family:region", Version: "1.0.0", CapabilityType: 3}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: nil,
		},
		{
			name: "valid request with multiple capabilities with different nomenclatures",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities: []kcr.CapabilitiesRegistryCapability{
						{LabelledName: "write_aptos-mainnet", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_aptos-testnet:region_secondary", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_aptos-testnet", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_avalanche-mainnet", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_avalanche-testnet-fuji", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_binance_smart_chain-mainnet", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_binance_smart_chain-testnet", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_bsc-testnet", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_celo-testnet-alfajores", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-mainnet-arbitrum-1", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-mainnet-base-1", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-mainnet-optimism-1", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-mainnet", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-testnet-sepolia-arbitrum-1", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-testnet-sepolia-base-1", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-testnet-sepolia-linea-1", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-testnet-sepolia-optimism-1", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_ethereum-testnet-sepolia", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_polygon-mainnet", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write_polygon-testnet-amoy", Version: "1.0.0", CapabilityType: 3},
						{LabelledName: "write-evm-celo-testnet-44787", Version: "1.0.0", CapabilityType: 3},
					},
					RegistryRef: te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: nil,
		},
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
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: changeset.CapabilityTypeTargetNamePrefix1, Version: "1.0.0", CapabilityType: changeset.CapabilityTypeTarget}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: changeset.ErrInvalidWriteCapName,
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
			name: "mixed chars after prefix as chain family",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: changeset.CapabilityTypeTargetNamePrefix2 + "test23-test", Version: "1.0.0", CapabilityType: 3}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: errors.New("chain family name 'test23' is not valid"),
		},
		{
			name: "mixed chars after prefix as network name",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: changeset.CapabilityTypeTargetNamePrefix1 + "test-cap123", Version: "1.0.0", CapabilityType: 3}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: errors.New("network name or chain ID 'cap123' is not valid"),
		},
		{
			name: "with chain family but without network name or chain ID",
			req: func(te test.EnvWrapper) (*changeset.AddCapabilitiesRequest, error) {
				return &changeset.AddCapabilitiesRequest{
					RegistryChainSel: te.RegistrySelector,
					Capabilities:     []kcr.CapabilitiesRegistryCapability{{LabelledName: changeset.CapabilityTypeTargetNamePrefix1 + "family-", Version: "1.0.0", CapabilityType: 3}},
					RegistryRef:      te.CapabilityRegistryAddressRef(),
				}, nil
			},
			expectedError: changeset.ErrEmptyWriteCapNetworkNameOrChainID,
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
				assert.ErrorContains(t, err, tt.expectedError.Error())
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
