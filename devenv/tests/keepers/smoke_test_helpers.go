package keepers

import (
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products/keepers"
)

type Test struct {
	ChainClient *seth.Client

	Config *keepers.Keepers

	LinkToken   contracts.LinkToken
	Transcoder  contracts.UpkeepTranscoder
	LINKETHFeed contracts.MockLINKETHFeed
	ETHUSDFeed  contracts.MockETHUSDFeed
	LINKUSDFeed contracts.MockETHUSDFeed
	WETHToken   contracts.WETHToken
	GasFeed     contracts.MockGasFeed
	Registry    contracts.KeeperRegistry
	Registrar   contracts.KeeperRegistrar

	RegistrySettings contracts.KeeperRegistrySettings

	Logger zerolog.Logger
}

func NewTest(
	l zerolog.Logger,
	chainClient *seth.Client,
	config *keepers.Keepers,
) (*Test, error) {
	t := &Test{
		ChainClient:      chainClient,
		Config:           config,
		RegistrySettings: config.GetRegistryConfig(),
		Logger:           l,
	}

	if err := t.LoadContracts(); err != nil {
		return nil, err
	}

	return t, nil
}

func (a *Test) LoadContracts() error {
	if err := a.LoadLINK(a.Config.DeployedContracts.LinkToken); err != nil {
		return fmt.Errorf("error loading link token contract: %w", err)
	}

	if err := a.LoadLinkEthFeed(a.Config.DeployedContracts.LinkEthFeed); err != nil {
		return fmt.Errorf("error loading link eth feed contract: %w", err)
	}

	if err := a.LoadEthGasFeed(a.Config.DeployedContracts.EthGasFeed); err != nil {
		return fmt.Errorf("error loading gas feed contract: %w", err)
	}

	if err := a.LoadTranscoder(a.Config.DeployedContracts.Transcoder); err != nil {
		return fmt.Errorf("error loading transcoder contract: %w", err)
	}

	if err := a.LoadRegistry(a.Config.DeployedContracts.Registry); err != nil {
		return fmt.Errorf("error loading registry contract: %w", err)
	}

	if a.Registry.RegistryOwnerAddress().String() != a.ChainClient.MustGetRootKeyAddress().String() {
		return errors.New("registry owner address is not the root key address")
	}

	if err := a.LoadRegistrar(a.Config.DeployedContracts.Registrar); err != nil {
		return fmt.Errorf("error loading registrar contract: %w", err)
	}

	return nil
}

func (a *Test) LoadLINK(address string) error {
	linkToken, err := contracts.LoadLinkTokenContract(a.Logger, a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LinkToken = linkToken
	a.Logger.Info().Str("LINK Token Address", a.LinkToken.Address()).Msg("Successfully loaded LINK Token")
	return nil
}

func (a *Test) LoadTranscoder(address string) error {
	transcoder, err := contracts.LoadUpkeepTranscoder(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.Transcoder = transcoder
	a.Logger.Info().Str("Transcoder Address", a.Transcoder.Address()).Msg("Successfully loaded Transcoder")
	return nil
}

func (a *Test) LoadLinkEthFeed(address string) error {
	ethLinkFeed, err := contracts.LoadMockLINKETHFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LINKETHFeed = ethLinkFeed
	a.Logger.Info().Str("LINK/ETH Feed Address", a.LINKETHFeed.Address()).Msg("Successfully loaded LINK/ETH Feed")
	return nil
}

func (a *Test) LoadEthUSDFeed(address string) error {
	ethUSDFeed, err := contracts.LoadMockETHUSDFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.ETHUSDFeed = ethUSDFeed
	a.Logger.Info().Str("ETH/USD Feed Address", a.ETHUSDFeed.Address()).Msg("Successfully loaded ETH/USD Feed")
	return nil
}

func (a *Test) LoadLinkUSDFeed(address string) error {
	linkUSDFeed, err := contracts.LoadMockETHUSDFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LINKUSDFeed = linkUSDFeed
	a.Logger.Info().Str("LINK/USD Feed Address", a.LINKUSDFeed.Address()).Msg("Successfully loaded LINK/USD Feed")
	return nil
}

func (a *Test) LoadWETH(address string) error {
	wethToken, err := contracts.LoadWETHTokenContract(a.Logger, a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.WETHToken = wethToken
	a.Logger.Info().Str("WETH Token Address", a.WETHToken.Address()).Msg("Successfully loaded WETH Token")
	return nil
}

func (a *Test) LoadEthGasFeed(address string) error {
	gasFeed, err := contracts.LoadMockGASFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.GasFeed = gasFeed
	a.Logger.Info().Str("Gas Feed Address", a.GasFeed.Address()).Msg("Successfully loaded Gas Feed")
	return nil
}

func (a *Test) LoadRegistry(registryAddress string) error {
	registry, err := contracts.LoadKeeperRegistry(a.Logger, a.ChainClient, common.HexToAddress(registryAddress), a.RegistrySettings.RegistryVersion, common.Address{})
	if err != nil {
		return err
	}
	a.Registry = registry
	a.Logger.Info().Str("Registry Address", a.Registry.Address()).Msg("Successfully loaded Registry")
	return nil
}

func (a *Test) LoadRegistrar(address string) error {
	if a.Registry == nil {
		return errors.New("registry must be deployed or loaded before registrar")
	}
	registrar, err := contracts.LoadKeeperRegistrar(a.ChainClient, common.HexToAddress(address), a.RegistrySettings.RegistryVersion)
	if err != nil {
		return err
	}
	a.Logger.Info().Str("Registrar Address", registrar.Address()).Msg("Successfully loaded Registrar")
	a.Registrar = registrar
	return nil
}
