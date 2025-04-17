package crecli

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

const (
	CRECLISettingsFileName     = "cre.yaml"
	CRECLIWorkflowSettingsFile = "workflow.yaml"
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
	FeedID          string `json:"feed_id"`
	URL             string `json:"url"`
	ConsumerAddress string `json:"consumer_address"`
	WriteTargetName string `json:"write_target_name"`
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

func PrepareCRECLISettingsFile(profile string, workflowOwner, capRegAddr, workflowRegistryAddr common.Address, dataFeedsCacheAddress *common.Address, donID uint32, chainSelector uint64, rpcHTTPURL string) (*os.File, error) {
	settingsFile, err := os.Create(CRECLISettingsFileName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create CRE CLI settings file")
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
					Name:          "CapabilitiesRegistry",
					Address:       capRegAddr.Hex(),
					ChainSelector: chainSelector,
				},
				{
					Name:          "WorkflowRegistry",
					Address:       workflowRegistryAddr.Hex(),
					ChainSelector: chainSelector,
				},
			},
		},
		Rpcs: []RPC{
			{
				ChainSelector: chainSelector,
				URL:           rpcHTTPURL,
			},
		},
		WorkflowStorage: WorkflowStorage{
			Gist: Gist{
				GithubToken: `${CRE_GITHUB_API_TOKEN}`,
			},
		},
	}

	if dataFeedsCacheAddress != nil {
		profileSettings.Contracts.DataFeeds = []ContractRegistry{
			{
				Name:          "DataFeedsCache",
				Address:       dataFeedsCacheAddress.Hex(),
				ChainSelector: chainSelector,
			},
		}
	}

	settings, err := setProfile(profile, profileSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set profile")
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

	settings, err := setProfile(profile, profileSettings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set profile")
	}

	settingsMarshalled, err := yaml.Marshal(settings)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal CRE CLI settings")
	}

	_, err = settingsFile.Write(settingsMarshalled)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to write %s settings file", CRECLIWorkflowSettingsFile)
	}

	return settingsFile, nil
}
