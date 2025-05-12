package crecli

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

const (
	CRECLISettingsFileName     = "cre.yaml"
	CRECLIWorkflowSettingsFile = "workflow.yaml"
	CRECLIProfile              = "test"
)

type Profiles struct {
	Test               Settings `yaml:"test,omitempty"`
	Staging            Settings `yaml:"staging,omitempty"`
	ProductionTestinet Settings `yaml:"production-testnet,omitempty"`
	Production         Settings `yaml:"production,omitempty"`
}

type Settings struct {
	DevPlatform     DevPlatform     `yaml:"dev-platform,omitempty"`
	UserWorkflow    UserWorkflow    `yaml:"user-workflow,omitempty"`
	Logging         Logging         `yaml:"logging,omitempty"`
	McmsConfig      McmsConfig      `yaml:"mcms-config,omitempty"`
	Contracts       Contracts       `yaml:"contracts,omitempty"`
	Rpcs            []RPC           `yaml:"rpcs,omitempty"`
	WorkflowStorage WorkflowStorage `yaml:"workflow_storage,omitempty"`
}

type DevPlatform struct {
	DonID uint32 `yaml:"don-id,omitempty"`
}

type UserWorkflow struct {
	WorkflowOwnerAddress string `yaml:"workflow-owner-address,omitempty"`
	WorkflowName         string `yaml:"workflow-name,omitempty"`
}

type Logging struct {
	SethConfigPath string `yaml:"seth-config-path,omitempty"`
}

type McmsConfig struct {
	ProposalsDirectory string `yaml:"proposals-directory,omitempty"`
}

type Contracts struct {
	ContractRegistry []ContractRegistry `yaml:"registries,omitempty"`
	DataFeeds        []ContractRegistry `yaml:"data-feeds,omitempty"`
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

type WorkflowStorage struct {
	Gist Gist `yaml:"gist"`
}

type Gist struct {
	GithubToken string `yaml:"github_token"`
}

type PoRWorkflowConfig struct {
	FeedID            string  `json:"feed_id"`
	URL               string  `json:"url"`
	ConsumerAddress   string  `json:"consumer_address"`
	WriteTargetName   string  `json:"write_target_name"`
	AuthKeySecretName *string `json:"auth_key_secret_name,omitempty"`
}

func setProfile(profile string, settings Settings) (Profiles, error) {
	var profiles Profiles

	switch profile {
	case "test":
		profiles = Profiles{Test: settings}
	case "staging":
		profiles = Profiles{Staging: settings}
	case "production-testnet":
		profiles = Profiles{ProductionTestinet: settings}
	case "production":
		profiles = Profiles{Production: settings}
	default:
		return Profiles{}, errors.Errorf("invalid profile: %s", profile)
	}

	return profiles, nil
}

// rpcs: chainSelector -> url
func PrepareCRECLISettingsFile(profile string, workflowOwner common.Address, addressBook cldf.AddressBook, donID uint32, homeChainSelector uint64, rpcs map[uint64]string) (*os.File, error) {
	settingsFile, err := os.Create(CRECLISettingsFileName)
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

	profileSettings := Settings{
		DevPlatform: DevPlatform{
			DonID: donID,
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
		WorkflowStorage: WorkflowStorage{
			Gist: Gist{
				GithubToken: `${CRE_GITHUB_API_TOKEN}`,
			},
		},
	}

	for chainSelector, rpc := range rpcs {
		profileSettings.Rpcs = append(profileSettings.Rpcs, RPC{
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
			profileSettings.Contracts.DataFeeds = append(profileSettings.Contracts.DataFeeds, ContractRegistry{
				Name:          df_changeset.DataFeedsCache.String(),
				Address:       dfAddr.Hex(),
				ChainSelector: chainSelector,
			})
		}
		// it is okay if there's no data feeds cache address for a chain

		forwaderAddr, forwaderErr := contracts.FindAddressesForChain(addressBook, chainSelector, string(keystone_changeset.KeystoneForwarder))
		if forwaderErr == nil {
			profileSettings.Contracts.Keystone = append(profileSettings.Contracts.Keystone, ContractRegistry{
				Name:          keystone_changeset.KeystoneForwarder.String(),
				Address:       forwaderAddr.Hex(),
				ChainSelector: chainSelector,
			})
		}
		// it is okay if there's no keystone forwarder address for a chain
	}

	settings, settingsErr := setProfile(profile, profileSettings)
	if settingsErr != nil {
		return nil, errors.Wrap(settingsErr, "failed to set profile")
	}

	settingsMarshalled, settingsMarshalledErr := yaml.Marshal(settings)
	if settingsMarshalledErr != nil {
		return nil, errors.Wrap(settingsMarshalledErr, "failed to marshal CRE CLI settings")
	}

	_, writeErr := settingsFile.Write(settingsMarshalled)
	if writeErr != nil {
		return nil, errors.Wrapf(writeErr, "failed to write %s settings file", CRECLISettingsFileName)
	}

	return settingsFile, nil
}

func PrepareCRECLIWorkflowSettingsFile(profile string, workflowOwner common.Address, workflowName string) (*os.File, error) {
	settingsFile, err := os.CreateTemp("", CRECLIWorkflowSettingsFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CRE CLI workflow settings file")
	}

	profileSettings := Settings{
		UserWorkflow: UserWorkflow{
			WorkflowOwnerAddress: workflowOwner.Hex(),
			WorkflowName:         workflowName,
		},
	}

	settings, settingsErr := setProfile(profile, profileSettings)
	if settingsErr != nil {
		return nil, errors.Wrap(settingsErr, "failed to set profile")
	}

	settingsMarshalled, marshallErr := yaml.Marshal(settings)
	if marshallErr != nil {
		return nil, errors.Wrap(marshallErr, "failed to marshal CRE CLI settings")
	}

	_, writeErr := settingsFile.Write(settingsMarshalled)
	if writeErr != nil {
		return nil, errors.Wrapf(writeErr, "failed to write %s settings file", CRECLIWorkflowSettingsFile)
	}

	return settingsFile, nil
}
