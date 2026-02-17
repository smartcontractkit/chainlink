package keepers

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"text/template"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
)

var ZeroAddress = common.Address{}

func deployContracts(l zerolog.Logger, chainClient *seth.Client, config *Keepers) error {
	if config.DeployedContracts.LinkToken == "" {
		linkToken, err := contracts.DeployLinkTokenContract(l, chainClient)
		if err != nil {
			return fmt.Errorf("error deploying link token contract: %w", err)
		}
		config.DeployedContracts.LinkToken = linkToken.Address()
	}

	if config.DeployedContracts.LinkEthFeed == "" {
		ethLinkFeed, err := contracts.DeployMockLINKETHFeed(chainClient, big.NewInt(2e18))
		if err != nil {
			return fmt.Errorf("error deploying mock eth feed contract: %w", err)
		}
		config.DeployedContracts.LinkEthFeed = ethLinkFeed.Address()
	}

	if config.DeployedContracts.EthGasFeed == "" {
		gasFeed, err := contracts.DeployMockGASFeed(chainClient, big.NewInt(2e11))
		if err != nil {
			return fmt.Errorf("error deploying mock gas feed contract: %w", err)
		}
		config.DeployedContracts.EthGasFeed = gasFeed.Address()
	}

	if config.DeployedContracts.Transcoder == "" {
		transcoder, err := contracts.DeployUpkeepTranscoder(chainClient)
		if err != nil {
			return fmt.Errorf("error deploying transcoder contract: %w", err)
		}
		config.DeployedContracts.Transcoder = transcoder.Address()
	}

	if config.DeployedContracts.Registry == "" {
		registry, err := contracts.DeployKeeperRegistry(chainClient, &contracts.KeeperRegistryOpts{
			RegistryVersion: config.MustGetRegistryVersion(),
			LinkAddr:        config.DeployedContracts.LinkToken,
			ETHFeedAddr:     config.DeployedContracts.LinkEthFeed,
			GasFeedAddr:     config.DeployedContracts.EthGasFeed,
			TranscoderAddr:  config.DeployedContracts.Transcoder,
			RegistrarAddr:   ZeroAddress.Hex(),
			Settings:        config.GetRegistryConfig(),
		})
		if err != nil {
			return fmt.Errorf("error deploying registry contract: %w", err)
		}
		config.DeployedContracts.Registry = registry.Address()
	}

	if config.DeployedContracts.Registrar == "" {
		registrarSettings := contracts.KeeperRegistrarSettings{
			AutoApproveConfigType: 2,
			AutoApproveMaxAllowed: math.MaxUint16,
			RegistryAddr:          config.DeployedContracts.Registry,
			MinLinkJuels:          big.NewInt(0),
		}

		registryInstance, err := contracts.LoadKeeperRegistry(l, chainClient, common.HexToAddress(config.DeployedContracts.Registry), config.MustGetRegistryVersion(), ZeroAddress)
		if err != nil {
			return fmt.Errorf("error loading registry contract: %w", err)
		}

		linkInstance, err := contracts.LoadLinkTokenContract(l, chainClient, common.HexToAddress(config.DeployedContracts.LinkToken))
		if err != nil {
			return fmt.Errorf("error loading link token contract: %w", err)
		}

		registrar, err := DeployKeeperRegistrar(chainClient, config.MustGetRegistryVersion(), linkInstance, registrarSettings, registryInstance)
		if err != nil {
			return fmt.Errorf("error deploying registrar contract: %w", err)
		}
		config.DeployedContracts.Registrar = registrar.Address()
	}

	return nil
}

func DeployKeeperRegistrar(
	client *seth.Client,
	registryVersion contracts.KeeperRegistryVersion,
	linkToken contracts.LinkToken,
	registrarSettings contracts.KeeperRegistrarSettings,
	registry contracts.KeeperRegistry,
) (contracts.KeeperRegistrar, error) {
	registrar, err := contracts.DeployKeeperRegistrar(client, registryVersion, linkToken.Address(), registrarSettings)
	if err != nil {
		return nil, err
	}
	if registryVersion != contracts.RegistryVersion_2_0 {
		err = registry.SetRegistrar(registrar.Address())
		if err != nil {
			return nil, err
		}
	}

	return registrar, nil
}

func createJobs(
	l zerolog.Logger,
	chainlinkNodes []*clclient.ChainlinkClient,
	config *Keepers,
	evmChainID string,
) error {
	for _, chainlinkNode := range chainlinkNodes {
		chainlinkNodeAddress, err := chainlinkNode.PrimaryEthAddress()
		if err != nil {
			l.Error().Err(err).Msg("Error retrieving chainlink node address")
			return err
		}
		_, err = chainlinkNode.MustCreateJob(&KeeperJobSpec{
			Name:                     "keeper-test-" + config.DeployedContracts.Registry,
			ContractAddress:          config.DeployedContracts.Registry,
			FromAddress:              chainlinkNodeAddress,
			EVMChainID:               evmChainID,
			MinIncomingConfirmations: 1,
		})
		if err != nil {
			l.Error().Err(err).Msg("Creating KeeperV2 Job shouldn't fail")
			return err
		}
	}
	return nil
}

// KeeperJobSpec represents a V2 keeper spec
type KeeperJobSpec struct {
	Name                     string `toml:"name"`
	ContractAddress          string `toml:"contractAddress"`
	FromAddress              string `toml:"fromAddress"` // Hex representation of the from address
	EVMChainID               string `toml:"evmChainID"`  // Not optional
	MinIncomingConfirmations int    `toml:"minIncomingConfirmations"`
}

// Type returns the type of the job
func (k *KeeperJobSpec) Type() string { return "keeper" }

// String representation of the job
func (k *KeeperJobSpec) String() (string, error) {
	keeperTemplateString := `
type                     = "keeper"
schemaVersion            = 1
name                     = "{{.Name}}"
contractAddress          = "{{.ContractAddress}}"
fromAddress              = "{{.FromAddress}}"
evmChainID		 		 = "{{.EVMChainID}}"
minIncomingConfirmations = {{.MinIncomingConfirmations}}
`
	var buf bytes.Buffer
	tmpl, err := template.New("Keeper Job").Parse(keeperTemplateString)
	if err != nil {
		return "", err
	}
	err = tmpl.Execute(&buf, *k)
	if err != nil {
		return "", err
	}
	return buf.String(), err
}
