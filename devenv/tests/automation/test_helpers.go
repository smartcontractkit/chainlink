package automation

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog"
	ocr3 "github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	ocr2keepers30config "github.com/smartcontractkit/chainlink-automation/pkg/v3/config"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	ctf_concurrency "github.com/smartcontractkit/chainlink/devenv/products/automation/concurrency"

	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/contracts/ethereum"
	"github.com/smartcontractkit/chainlink/devenv/products/automation"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/automation_registrar_wrapper2_1"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/keeper_registrar_wrapper2_0"
)

type AutomationTest struct {
	ChainClient *seth.Client

	Config automation.Automation

	LinkToken   contracts.LinkToken
	Transcoder  contracts.UpkeepTranscoder
	LINKETHFeed contracts.MockLINKETHFeed
	ETHUSDFeed  contracts.MockETHUSDFeed
	LINKUSDFeed contracts.MockETHUSDFeed
	WETHToken   contracts.WETHToken
	GasFeed     contracts.MockGasFeed
	Registry    contracts.KeeperRegistry
	Registrar   contracts.KeeperRegistrar

	RegistrySettings       contracts.KeeperRegistrySettings
	RegistrarSettings      contracts.KeeperRegistrarSettings
	PluginConfig           ocr2keepers30config.OffchainConfig
	PublicConfig           ocr3.PublicConfig
	UpkeepPrivilegeManager common.Address
	UpkeepIDs              []*big.Int

	ChainlinkNodes []*clclient.ChainlinkClient

	TransmitterKeyIndex int

	Logger zerolog.Logger
}

type UpkeepConfig struct {
	UpkeepName     string
	EncryptedEmail []byte
	UpkeepContract common.Address
	GasLimit       uint32
	AdminAddress   common.Address
	TriggerType    uint8
	CheckData      []byte
	TriggerConfig  []byte
	OffchainConfig []byte
	FundingAmount  *big.Int
}

func NewAutomationTest(
	l zerolog.Logger,
	chainClient *seth.Client,
	chainlinkNodes []*clclient.ChainlinkClient,
	config automation.Automation,
) *AutomationTest {
	return &AutomationTest{
		ChainClient:            chainClient,
		Config:                 config,
		ChainlinkNodes:         chainlinkNodes,
		TransmitterKeyIndex:    0,
		UpkeepPrivilegeManager: chainClient.MustGetRootKeyAddress(),
		Logger:                 l,
	}
}

func (a *AutomationTest) LoadLINK(address string) error {
	linkToken, err := contracts.LoadLinkTokenContract(a.Logger, a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LinkToken = linkToken
	a.Logger.Info().Str("LINK Token Address", a.LinkToken.Address()).Msg("Successfully loaded LINK Token")
	return nil
}

func (a *AutomationTest) LoadTranscoder(address string) error {
	transcoder, err := contracts.LoadUpkeepTranscoder(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.Transcoder = transcoder
	a.Logger.Info().Str("Transcoder Address", a.Transcoder.Address()).Msg("Successfully loaded Transcoder")
	return nil
}

func (a *AutomationTest) LoadLinkEthFeed(address string) error {
	ethLinkFeed, err := contracts.LoadMockLINKETHFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LINKETHFeed = ethLinkFeed
	a.Logger.Info().Str("LINK/ETH Feed Address", a.LINKETHFeed.Address()).Msg("Successfully loaded LINK/ETH Feed")
	return nil
}

func (a *AutomationTest) LoadEthUSDFeed(address string) error {
	ethUSDFeed, err := contracts.LoadMockETHUSDFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.ETHUSDFeed = ethUSDFeed
	a.Logger.Info().Str("ETH/USD Feed Address", a.ETHUSDFeed.Address()).Msg("Successfully loaded ETH/USD Feed")
	return nil
}

func (a *AutomationTest) LoadLinkUSDFeed(address string) error {
	linkUSDFeed, err := contracts.LoadMockETHUSDFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.LINKUSDFeed = linkUSDFeed
	a.Logger.Info().Str("LINK/USD Feed Address", a.LINKUSDFeed.Address()).Msg("Successfully loaded LINK/USD Feed")
	return nil
}

func (a *AutomationTest) LoadWETH(address string) error {
	wethToken, err := contracts.LoadWETHTokenContract(a.Logger, a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.WETHToken = wethToken
	a.Logger.Info().Str("WETH Token Address", a.WETHToken.Address()).Msg("Successfully loaded WETH Token")
	return nil
}

func (a *AutomationTest) LoadEthGasFeed(address string) error {
	gasFeed, err := contracts.LoadMockGASFeed(a.ChainClient, common.HexToAddress(address))
	if err != nil {
		return err
	}
	a.GasFeed = gasFeed
	a.Logger.Info().Str("Gas Feed Address", a.GasFeed.Address()).Msg("Successfully loaded Gas Feed")
	return nil
}

func (a *AutomationTest) LoadRegistry(registryAddress, chainModuleAddress string) error {
	registry, err := contracts.LoadKeeperRegistry(a.Logger, a.ChainClient, common.HexToAddress(registryAddress), a.RegistrySettings.RegistryVersion, common.HexToAddress(chainModuleAddress))
	if err != nil {
		return err
	}
	a.Registry = registry
	a.Logger.Info().Str("ChainModule Address", chainModuleAddress).Str("Registry Address", a.Registry.Address()).Msg("Successfully loaded Registry")
	return nil
}

func (a *AutomationTest) LoadRegistrar(address string) error {
	if a.Registry == nil {
		return errors.New("registry must be deployed or loaded before registrar")
	}
	a.RegistrarSettings.RegistryAddr = a.Registry.Address()
	registrar, err := contracts.LoadKeeperRegistrar(a.ChainClient, common.HexToAddress(address), a.RegistrySettings.RegistryVersion)
	if err != nil {
		return err
	}
	a.Logger.Info().Str("Registrar Address", registrar.Address()).Msg("Successfully loaded Registrar")
	a.Registrar = registrar
	return nil
}

type registrationResult struct {
	txHash common.Hash
}

func (r registrationResult) GetResult() common.Hash {
	return r.txHash
}

func (a *AutomationTest) RegisterUpkeeps(upkeepConfigs []UpkeepConfig, maxConcurrency int) ([]common.Hash, error) {
	concurrency, err := automation.GetAndAssertCorrectConcurrency(a.ChainClient, 1)
	if err != nil {
		return nil, err
	}

	if concurrency > maxConcurrency {
		concurrency = maxConcurrency
		a.Logger.Debug().
			Msgf("Concurrency is higher than max concurrency, setting concurrency to %d", concurrency)
	}

	var registerUpkeep = func(resultCh chan registrationResult, errorCh chan error, executorNum int, upkeepConfig UpkeepConfig) {
		keyNum := executorNum + 1 // key 0 is the root key
		var registrationRequest []byte
		var registrarABI *abi.ABI
		var err error
		switch a.RegistrySettings.RegistryVersion {
		case ethereum.RegistryVersion_2_0:
			registrarABI, err = keeper_registrar_wrapper2_0.KeeperRegistrarMetaData.GetAbi()
			if err != nil {
				errorCh <- errors.Join(err, errors.New("failed to get registrar abi"))
				return
			}
			registrationRequest, err = registrarABI.Pack(
				"register",
				upkeepConfig.UpkeepName,
				upkeepConfig.EncryptedEmail,
				upkeepConfig.UpkeepContract,
				upkeepConfig.GasLimit,
				upkeepConfig.AdminAddress,
				upkeepConfig.CheckData,
				upkeepConfig.OffchainConfig,
				upkeepConfig.FundingAmount,
				a.ChainClient.Addresses[keyNum])
			if err != nil {
				errorCh <- errors.Join(err, errors.New("failed to pack registrar request"))
				return
			}
		case ethereum.RegistryVersion_2_1, ethereum.RegistryVersion_2_2: // 2.1 and 2.2 use the same registrar
			registrarABI, err = automation_registrar_wrapper2_1.AutomationRegistrarMetaData.GetAbi()
			if err != nil {
				errorCh <- errors.Join(err, errors.New("failed to get registrar abi"))
				return
			}
			registrationRequest, err = registrarABI.Pack(
				"register",
				upkeepConfig.UpkeepName,
				upkeepConfig.EncryptedEmail,
				upkeepConfig.UpkeepContract,
				upkeepConfig.GasLimit,
				upkeepConfig.AdminAddress,
				upkeepConfig.TriggerType,
				upkeepConfig.CheckData,
				upkeepConfig.TriggerConfig,
				upkeepConfig.OffchainConfig,
				upkeepConfig.FundingAmount,
				a.ChainClient.Addresses[keyNum])
			if err != nil {
				errorCh <- errors.Join(err, errors.New("failed to pack registrar request"))
				return
			}
		default:
			errorCh <- errors.New("v2.0, v2.1, and v2.2 are the only supported versions")
			return
		}

		tx, err := a.LinkToken.TransferAndCallFromKey(a.Registrar.Address(), upkeepConfig.FundingAmount, registrationRequest, keyNum)
		if err != nil {
			errorCh <- errors.Join(err, fmt.Errorf("client number %d failed to register upkeep %s", keyNum, upkeepConfig.UpkeepContract.Hex()))
			return
		}

		resultCh <- registrationResult{txHash: tx.Hash()}
	}

	executor := ctf_concurrency.NewConcurrentExecutor[common.Hash, registrationResult, UpkeepConfig](a.Logger)
	results, err := executor.Execute(concurrency, upkeepConfigs, registerUpkeep)
	if err != nil {
		return nil, err
	}

	if len(results) != len(upkeepConfigs) {
		return nil, fmt.Errorf("failed to register all upkeeps. Expected %d, got %d", len(upkeepConfigs), len(results))
	}

	a.Logger.Info().Msg("Successfully registered all upkeeps")

	return results, nil
}

type confirmationResult struct {
	upkeepID automation.UpkeepId
}

func (c confirmationResult) GetResult() automation.UpkeepId {
	return c.upkeepID
}

func (a *AutomationTest) ConfirmUpkeepsRegistered(registrationTxHashes []common.Hash, maxConcurrency int) ([]*big.Int, error) {
	concurrency, err := automation.GetAndAssertCorrectConcurrency(a.ChainClient, 1)
	if err != nil {
		return nil, err
	}

	if concurrency > maxConcurrency {
		concurrency = maxConcurrency
		a.Logger.Debug().
			Msgf("Concurrency is higher than max concurrency, setting concurrency to %d", concurrency)
	}

	var confirmUpkeep = func(resultCh chan confirmationResult, errorCh chan error, _ int, txHash common.Hash) {
		receipt, err := a.ChainClient.Client.TransactionReceipt(context.Background(), txHash)
		if err != nil {
			errorCh <- errors.Join(err, errors.New("failed to confirm upkeep registration"))
			return
		}

		var upkeepId *big.Int
		for _, rawLog := range receipt.Logs {
			parsedUpkeepId, err := a.Registry.ParseUpkeepIdFromRegisteredLog(rawLog)
			if err == nil {
				upkeepId = parsedUpkeepId
				break
			}
		}
		if upkeepId == nil {
			errorCh <- errors.New("failed to parse upkeep id from registration receipt")
			return
		}
		resultCh <- confirmationResult{upkeepID: upkeepId}
	}

	executor := ctf_concurrency.NewConcurrentExecutor[automation.UpkeepId, confirmationResult, common.Hash](a.Logger)
	results, err := executor.Execute(concurrency, registrationTxHashes, confirmUpkeep)

	if err != nil {
		return nil, fmt.Errorf("failed confirmations: %d | successful confirmations: %d", len(executor.GetErrors()), len(results))
	}

	if len(registrationTxHashes) != len(results) {
		return nil, fmt.Errorf("failed to confirm all upkeeps. Expected %d, got %d", len(registrationTxHashes), len(results))
	}

	seen := make(map[*big.Int]bool)
	for _, upkeepId := range results {
		if seen[upkeepId] {
			return nil, fmt.Errorf("duplicate upkeep id: %s. Something went wrong during upkeep confirmation. Please check the test code", upkeepId.String())
		}
		seen[upkeepId] = true
	}

	a.Logger.Info().Msg("Successfully confirmed all upkeeps")
	a.UpkeepIDs = results

	return results, nil
}

func (a *AutomationTest) LoadContracts() error {
	if err := a.LoadLINK(a.Config.DeployedContracts.LinkToken); err != nil {
		return fmt.Errorf("error loading link token contract: %w", err)
	}

	if err := a.LoadWETH(a.Config.DeployedContracts.Weth); err != nil {
		return fmt.Errorf("error loading weth token contract: %w", err)
	}

	if err := a.LoadLinkEthFeed(a.Config.DeployedContracts.LinkEthFeed); err != nil {
		return fmt.Errorf("error loading link eth feed contract: %w", err)
	}

	if err := a.LoadEthGasFeed(a.Config.DeployedContracts.EthGasFeed); err != nil {
		return fmt.Errorf("error loading gas feed contract: %w", err)
	}

	if err := a.LoadEthUSDFeed(a.Config.DeployedContracts.EthUSDFeed); err != nil {
		return fmt.Errorf("error loading eth usd feed contract: %w", err)
	}

	if err := a.LoadLinkUSDFeed(a.Config.DeployedContracts.LinkUSDFeed); err != nil {
		return fmt.Errorf("error loading link usd feed contract: %w", err)
	}

	if err := a.LoadTranscoder(a.Config.DeployedContracts.Transcoder); err != nil {
		return fmt.Errorf("error loading transcoder contract: %w", err)
	}

	if err := a.LoadRegistry(a.Config.DeployedContracts.Registry, a.Config.DeployedContracts.ChainModule); err != nil {
		return fmt.Errorf("error loading registry contract: %w", err)
	}

	if a.Registry.RegistryOwnerAddress().String() != a.ChainClient.MustGetRootKeyAddress().String() {
		return fmt.Errorf("registry owner address is not the root key address")
	}

	if err := a.LoadRegistrar(a.Config.DeployedContracts.Registrar); err != nil {
		return fmt.Errorf("error loading registrar contract: %w", err)
	}

	return nil
}
