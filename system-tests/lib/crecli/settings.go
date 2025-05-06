package crecli

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/smartcontractkit/chainlink/deployment"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

const (
	CRECLISettingsFileName = ".cre-cli-settings.yaml"
)

type Settings struct {
	DevPlatform  DevPlatform  `yaml:"dev-platform"`
	UserWorkflow UserWorkflow `yaml:"user-workflow"`
	Logging      Logging      `yaml:"logging"`
	McmsConfig   McmsConfig   `yaml:"mcms-config"`
	Contracts    Contracts    `yaml:"contracts"`
	Rpcs         []RPC        `yaml:"rpcs"`
}

type DevPlatform struct {
	CapabilitiesRegistryAddress string `yaml:"capabilities-registry-contract-address"`
	DonID                       uint32 `yaml:"don-id"`
	WorkflowRegistryAddress     string `yaml:"workflow-registry-contract-address"`
}

type UserWorkflow struct {
	WorkflowOwnerAddress string `yaml:"workflow-owner-address"`
}

type Logging struct {
	SethConfigPath string `yaml:"seth-config-path"`
}

type McmsConfig struct {
	ProposalsDirectory string `yaml:"proposals-directory"`
}

type Contracts struct {
	ContractRegistry []ContractRegistry `yaml:"registries"`
	DataFeeds        []ContractRegistry `yaml:"data-feeds"`
	Keystone         []ContractRegistry `yaml:"keystone"`
}

type ContractRegistry struct {
	Name          string `yaml:"name"`
	Address       string `yaml:"address"`
	ChainSelector uint64 `yaml:"chain-selector"`
}

type RPC struct {
	ChainSelector uint64 `yaml:"chain-selector"`
	URL           string `yaml:"url"`
}

type PoRWorkflowConfig struct {
	FeedID            string  `json:"feed_id"`
	URL               string  `json:"url"`
	ConsumerAddress   string  `json:"consumer_address"`
	WriteTargetName   string  `json:"write_target_name"`
	AuthKeySecretName *string `json:"auth_key_secret_name,omitempty"`
}

// rpcs: chainSelector -> url
func PrepareCRECLISettingsFile(workflowOwner common.Address, addressBook deployment.AddressBook, donID uint32, homeChainSelector uint64, rpcs map[uint64]string) (*os.File, error) {
	settingsFile, err := os.CreateTemp("", CRECLISettingsFileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CRE CLI settings file")
	}

	capRegAddr, capRegErr := contracts.FindAddressesForChain(addressBook, homeChainSelector, keystone_changeset.CapabilitiesRegistry.String())
	if capRegErr != nil {
		return nil, errors.Wrapf(capRegErr, "failed to get capabilities registry address for chain %d", homeChainSelector)
	}

	workflowRegistryAddr, workflowRegistryErr := contracts.FindAddressesForChain(addressBook, homeChainSelector, keystone_changeset.WorkflowRegistry.String())
	if workflowRegistryErr != nil {
		return nil, errors.Wrapf(workflowRegistryErr, "failed to get workflow registry address for chain %d", homeChainSelector)
	}

	settings := Settings{
		DevPlatform: DevPlatform{
			CapabilitiesRegistryAddress: capRegAddr.Hex(),
			DonID:                       donID,
			WorkflowRegistryAddress:     workflowRegistryAddr.Hex(),
		},
		UserWorkflow: UserWorkflow{
			WorkflowOwnerAddress: workflowOwner.Hex(),
		},
		Logging: Logging{},
		McmsConfig: McmsConfig{
			ProposalsDirectory: "./",
		},
		Contracts: Contracts{
			ContractRegistry: []ContractRegistry{
				{
					Name:          keystone_changeset.CapabilitiesRegistry.String(),
					Address:       capRegAddr.Hex(),
					ChainSelector: homeChainSelector,
				},
				{
					Name:          keystone_changeset.WorkflowRegistry.String(),
					Address:       workflowRegistryAddr.Hex(),
					ChainSelector: homeChainSelector,
				},
			},
		},
	}

	for chainSelector, rpc := range rpcs {
		settings.Rpcs = append(settings.Rpcs, RPC{
			ChainSelector: chainSelector,
			URL:           rpc,
		})
	}

	addresses, addrErr := addressBook.Addresses()
	if addrErr != nil {
		return nil, errors.Wrap(addrErr, "failed to get address book addresses")
	}

	for chainSelector := range addresses {
		dfAddr, dfErr := contracts.FindAddressesForChain(addressBook, chainSelector, df_changeset.DataFeedsCache.String())
		if dfErr == nil {
			settings.Contracts.DataFeeds = append(settings.Contracts.DataFeeds, ContractRegistry{
				Name:          df_changeset.DataFeedsCache.String(),
				Address:       dfAddr.Hex(),
				ChainSelector: chainSelector,
			})
		}
		// it is okay if there's no data feeds cache address for a chain

		forwaderAddr, forwaderErr := contracts.FindAddressesForChain(addressBook, chainSelector, string(keystone_changeset.KeystoneForwarder))
		if forwaderErr == nil {
			settings.Contracts.Keystone = append(settings.Contracts.Keystone, ContractRegistry{
				Name:          keystone_changeset.KeystoneForwarder.String(),
				Address:       forwaderAddr.Hex(),
				ChainSelector: chainSelector,
			})
		}
		// it is okay if there's no keystone forwarder address for a chain
	}

	settingsMarshalled, err := yaml.Marshal(settings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal CRE CLI settings")
	}

	_, err = settingsFile.Write(settingsMarshalled)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write %s settings file", CRECLISettingsFileName)
	}

	return settingsFile, nil
}
