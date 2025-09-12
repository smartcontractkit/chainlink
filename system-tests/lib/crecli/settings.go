/*
We will keep this file for now, because we want to be able to create the `cre.yaml` file used by the CRE CLI v0.2.x,
when local CRE is started or when sandboxes are created.
*/
package crecli

import (
	"os"

	"github.com/google/uuid"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/components/s3provider"
	df_changeset "github.com/smartcontractkit/chainlink/deployment/data-feeds/changeset"
	keystone_changeset "github.com/smartcontractkit/chainlink/deployment/keystone/changeset"
	"github.com/smartcontractkit/chainlink/system-tests/lib/cre/contracts"
)

const (
	CRECLISettingsFileName = "cre.yaml"
	CRECLIProfile          = "test"
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
	DonID uint64 `yaml:"don-id,omitempty"`
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
	Gist  Gist                 `yaml:"gist"`
	Minio MinioStorageSettings `yaml:"minio,omitempty"` // Optional, if not provided, Gist will be used
}

type Gist struct {
	GithubToken string `yaml:"github_token"`
}

type MinioStorageSettings struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	SessionToken    string `yaml:"session_token"`
	UseSSL          bool   `yaml:"use_ssl"`
	Region          string `yaml:"region"`
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
func PrepareCRECLISettingsFile(
	profile string,
	workflowOwner common.Address,
	addressBook cldf.AddressBook,
	donID uint64,
	homeChainSelector uint64,
	rpcs map[uint64]string,
	s3ProviderOutput *s3provider.Output,
) (*os.File, error) {
	settingsFile, err := os.Create(CRECLISettingsFileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CRE CLI settings file")
	}

	capRegAddr, _, capRegErr := contracts.FindAddressesForChain(addressBook, homeChainSelector, keystone_changeset.CapabilitiesRegistry.String())
	if capRegErr != nil {
		return nil, errors.Wrapf(capRegErr, "failed to get capabilities registry address for chain %d", homeChainSelector)
	}

	workflowRegistryAddr, _, workflowRegistryErr := contracts.FindAddressesForChain(addressBook, homeChainSelector, keystone_changeset.WorkflowRegistry.String())
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
	}

	if s3ProviderOutput != nil {
		profileSettings.WorkflowStorage.Minio = MinioStorageSettings{
			Endpoint:        s3ProviderOutput.Endpoint,
			AccessKeyID:     s3ProviderOutput.AccessKey,
			SecretAccessKey: s3ProviderOutput.SecretKey,
			SessionToken:    uuid.NewString(),
			UseSSL:          false,
			Region:          s3ProviderOutput.Region,
		}
	}

	profileSettings.WorkflowStorage.Gist = Gist{
		GithubToken: `${CRE_GITHUB_API_TOKEN}`,
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
		dfAddr, _, dfErr := contracts.FindAddressesForChain(addressBook, chainSelector, df_changeset.DataFeedsCache.String())
		if dfErr == nil {
			profileSettings.Contracts.DataFeeds = append(profileSettings.Contracts.DataFeeds, ContractRegistry{
				Name:          df_changeset.DataFeedsCache.String(),
				Address:       dfAddr.Hex(),
				ChainSelector: chainSelector,
			})
		}
		// it is okay if there's no data feeds cache address for a chain

		forwaderAddr, _, forwaderErr := contracts.FindAddressesForChain(addressBook, chainSelector, string(keystone_changeset.KeystoneForwarder))
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
