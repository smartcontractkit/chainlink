package contracts

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/contracts/ethereum"
	ocrConfigHelper "github.com/smartcontractkit/libocr/offchainreporting/confighelper"
)

// ContractDeployer is an interface for abstracting the contract deployment methods across network implementations
type ContractDeployer interface {
	DeployStorageContract() (Storage, error)
	DeployAPIConsumer(linkAddr string) (APIConsumer, error)
	DeployOracle(linkAddr string) (Oracle, error)
	DeployReadAccessController() (ReadAccessController, error)
	DeployFlags(rac string) (Flags, error)
	DeployDeviationFlaggingValidator(
		flags string,
		flaggingThreshold *big.Int,
	) (DeviationFlaggingValidator, error)
	DeployFluxAggregatorContract(linkAddr string, fluxOptions FluxAggregatorOptions) (FluxAggregator, error)
	DeployLinkTokenContract() (LinkToken, error)
	DeployOffChainAggregator(linkAddr string, offchainOptions OffchainOptions) (OffchainAggregator, error)
	DeployVRFContract() (VRF, error)
	DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error)
	DeployMockGasFeed(answer *big.Int) (MockGasFeed, error)
	DeployKeeperRegistrar(linkAddr string, registrarSettings KeeperRegistrarSettings) (KeeperRegistrar, error)
	DeployUpkeepTranscoder() (UpkeepTranscoder, error)
	DeployKeeperRegistry(opts *KeeperRegistryOpts) (KeeperRegistry, error)
	DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error)
	DeployKeeperConsumerPerformance(
		testBlockRange,
		averageCadence,
		checkGasToBurn,
		performGasToBurn *big.Int,
	) (KeeperConsumerPerformance, error)
	DeployUpkeepCounter(testRange *big.Int, interval *big.Int) (UpkeepCounter, error)
	DeployUpkeepPerformCounterRestrictive(testRange *big.Int, averageEligibilityCadence *big.Int) (UpkeepPerformCounterRestrictive, error)
	DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error)
	DeployVRFConsumerV2(linkAddr string, coordinatorAddr string) (VRFConsumerV2, error)
	DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error)
	DeployVRFCoordinatorV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (VRFCoordinatorV2, error)
	DeployBlockhashStore() (BlockHashStore, error)
}

// NewContractDeployer returns an instance of a contract deployer based on the client type
func NewContractDeployer(bcClient blockchain.EVMClient) (ContractDeployer, error) {
	switch clientImpl := bcClient.Get().(type) {
	case *blockchain.EthereumClient:
		return NewEthereumContractDeployer(clientImpl), nil
	case *blockchain.KlaytnClient:
		return &KlaytnContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	case *blockchain.MetisClient:
		return &MetisContractDeployer{NewEthereumContractDeployer(clientImpl)}, nil
	}
	return nil, errors.New("unknown blockchain client implementation for contract deployer. Register blockchain client in NewContractDeployer")
}

// EthereumContractDeployer provides the implementations for deploying ETH (EVM) based contracts
type EthereumContractDeployer struct {
	client blockchain.EVMClient
}

// KlaytnContractDeployer wraps ethereum contract deployments for Klaytn
type KlaytnContractDeployer struct {
	*EthereumContractDeployer
}

// MetisContractDeployer wraps ethereum contract deployments for Metis
type MetisContractDeployer struct {
	*EthereumContractDeployer
}

// NewEthereumContractDeployer returns an instantiated instance of the ETH contract deployer
func NewEthereumContractDeployer(ethClient blockchain.EVMClient) *EthereumContractDeployer {
	return &EthereumContractDeployer{
		client: ethClient,
	}
}

// DefaultFluxAggregatorOptions produces some basic defaults for a flux aggregator contract
func DefaultFluxAggregatorOptions() FluxAggregatorOptions {
	return FluxAggregatorOptions{
		PaymentAmount: big.NewInt(1),
		Timeout:       uint32(30),
		MinSubValue:   big.NewInt(0),
		MaxSubValue:   big.NewInt(1000000000000),
		Decimals:      uint8(0),
		Description:   "Test Flux Aggregator",
	}
}

// DeployReadAccessController deploys read/write access controller contract
func (e *EthereumContractDeployer) DeployReadAccessController() (ReadAccessController, error) {
	address, _, instance, err := e.client.DeployContract("Read Access Controller", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeploySimpleReadAccessController(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumReadAccessController{
		client:  e.client,
		rac:     instance.(*ethereum.SimpleReadAccessController),
		address: address,
	}, nil
}

// DeployFlags deploys flags contract
func (e *EthereumContractDeployer) DeployFlags(
	rac string,
) (Flags, error) {
	address, _, instance, err := e.client.DeployContract("Flags", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		racAddr := common.HexToAddress(rac)
		return ethereum.DeployFlags(auth, backend, racAddr)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFlags{
		client:  e.client,
		flags:   instance.(*ethereum.Flags),
		address: address,
	}, nil
}

// DeployDeviationFlaggingValidator deploys deviation flagging validator contract
func (e *EthereumContractDeployer) DeployDeviationFlaggingValidator(
	flags string,
	flaggingThreshold *big.Int,
) (DeviationFlaggingValidator, error) {
	address, _, instance, err := e.client.DeployContract("Deviation flagging validator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		flagAddr := common.HexToAddress(flags)
		return ethereum.DeployDeviationFlaggingValidator(auth, backend, flagAddr, flaggingThreshold)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumDeviationFlaggingValidator{
		client:  e.client,
		dfv:     instance.(*ethereum.DeviationFlaggingValidator),
		address: address,
	}, nil
}

// DeployFluxAggregatorContract deploys the Flux Aggregator Contract on an EVM chain
func (e *EthereumContractDeployer) DeployFluxAggregatorContract(
	linkAddr string,
	fluxOptions FluxAggregatorOptions,
) (FluxAggregator, error) {
	address, _, instance, err := e.client.DeployContract("Flux Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return ethereum.DeployFluxAggregator(auth,
			backend,
			la,
			fluxOptions.PaymentAmount,
			fluxOptions.Timeout,
			fluxOptions.Validator,
			fluxOptions.MinSubValue,
			fluxOptions.MaxSubValue,
			fluxOptions.Decimals,
			fluxOptions.Description)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumFluxAggregator{
		client:         e.client,
		fluxAggregator: instance.(*ethereum.FluxAggregator),
		address:        address,
	}, nil
}

// DeployLinkTokenContract deploys a Link Token contract to an EVM chain
func (e *EthereumContractDeployer) DeployLinkTokenContract() (LinkToken, error) {
	linkTokenAddress, _, instance, err := e.client.DeployContract("LINK Token", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployLinkToken(auth, backend)
	})
	if err != nil {
		return nil, err
	}

	return &EthereumLinkToken{
		client:   e.client,
		instance: instance.(*ethereum.LinkToken),
		address:  *linkTokenAddress,
	}, err
}

// DefaultOffChainAggregatorOptions returns some base defaults for deploying an OCR contract
func DefaultOffChainAggregatorOptions() OffchainOptions {
	return OffchainOptions{
		MaximumGasPrice:         uint32(3000),
		ReasonableGasPrice:      uint32(10),
		MicroLinkPerEth:         uint32(500),
		LinkGweiPerObservation:  uint32(500),
		LinkGweiPerTransmission: uint32(500),
		MinimumAnswer:           big.NewInt(1),
		MaximumAnswer:           big.NewInt(50000000000000000),
		Decimals:                8,
		Description:             "Test OCR",
	}
}

// DefaultOffChainAggregatorConfig returns some base defaults for configuring an OCR contract
func DefaultOffChainAggregatorConfig(numberNodes int) OffChainAggregatorConfig {
	if numberNodes <= 4 {
		log.Err(fmt.Errorf("Insufficient number of nodes (%d) supplied for OCR, need at least 5", numberNodes)).
			Int("Number Chainlink Nodes", numberNodes).
			Msg("You likely need more chainlink nodes to properly configure OCR, try 5 or more.")
	}
	s := []int{1}
	// First node's stage already inputted as a 1 in line above, so numberNodes-1.
	for i := 0; i < numberNodes-1; i++ {
		s = append(s, 2)
	}
	return OffChainAggregatorConfig{
		AlphaPPB:         1,
		DeltaC:           time.Minute * 60,
		DeltaGrace:       time.Second * 12,
		DeltaProgress:    time.Second * 35,
		DeltaStage:       time.Second * 60,
		DeltaResend:      time.Second * 17,
		DeltaRound:       time.Second * 30,
		RMax:             6,
		S:                s,
		N:                numberNodes,
		F:                1,
		OracleIdentities: []ocrConfigHelper.OracleIdentityExtra{},
	}
}

// DeployOffChainAggregator deploys the offchain aggregation contract to the EVM chain
func (e *EthereumContractDeployer) DeployOffChainAggregator(
	linkAddr string,
	offchainOptions OffchainOptions,
) (OffchainAggregator, error) {
	address, _, instance, err := e.client.DeployContract("OffChain Aggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		la := common.HexToAddress(linkAddr)
		return ethereum.DeployOffchainAggregator(auth,
			backend,
			offchainOptions.MaximumGasPrice,
			offchainOptions.ReasonableGasPrice,
			offchainOptions.MicroLinkPerEth,
			offchainOptions.LinkGweiPerObservation,
			offchainOptions.LinkGweiPerTransmission,
			la,
			offchainOptions.MinimumAnswer,
			offchainOptions.MaximumAnswer,
			offchainOptions.BillingAccessController,
			offchainOptions.RequesterAccessController,
			offchainOptions.Decimals,
			offchainOptions.Description)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOffchainAggregator{
		client:  e.client,
		ocr:     instance.(*ethereum.OffchainAggregator),
		address: address,
	}, err
}

// DeployStorageContract deploys a vanilla storage contract that is a value store
func (e *EthereumContractDeployer) DeployStorageContract() (Storage, error) {
	_, _, instance, err := e.client.DeployContract("Storage", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumStorage{
		client: e.client,
		store:  instance.(*ethereum.Store),
	}, err
}

// DeployAPIConsumer deploys api consumer for oracle
func (e *EthereumContractDeployer) DeployAPIConsumer(linkAddr string) (APIConsumer, error) {
	addr, _, instance, err := e.client.DeployContract("APIConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployAPIConsumer(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumAPIConsumer{
		address:  addr,
		client:   e.client,
		consumer: instance.(*ethereum.APIConsumer),
	}, err
}

// DeployOracle deploys oracle for consumer test
func (e *EthereumContractDeployer) DeployOracle(linkAddr string) (Oracle, error) {
	addr, _, instance, err := e.client.DeployContract("Oracle", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployOracle(auth, backend, common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumOracle{
		address: addr,
		client:  e.client,
		oracle:  instance.(*ethereum.Oracle),
	}, err
}

// DeployVRFContract deploy VRF contract
func (e *EthereumContractDeployer) DeployVRFContract() (VRF, error) {
	address, _, instance, err := e.client.DeployContract("VRF", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRF(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRF{
		client:  e.client,
		vrf:     instance.(*ethereum.VRF),
		address: address,
	}, err
}

func (e *EthereumContractDeployer) DeployMockETHLINKFeed(answer *big.Int) (MockETHLINKFeed, error) {
	address, _, instance, err := e.client.DeployContract("MockETHLINKAggregator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployMockV3AggregatorContract(auth, backend, 18, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockETHLINKFeed{
		client:  e.client,
		feed:    instance.(*ethereum.MockV3AggregatorContract),
		address: address,
	}, err
}

func (e *EthereumContractDeployer) DeployMockGasFeed(answer *big.Int) (MockGasFeed, error) {
	address, _, instance, err := e.client.DeployContract("MockGasFeed", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployMockGASAggregator(auth, backend, answer)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumMockGASFeed{
		client:  e.client,
		feed:    instance.(*ethereum.MockGASAggregator),
		address: address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepTranscoder() (UpkeepTranscoder, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepTranscoder", func(
		opts *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployUpkeepTranscoder(opts, backend)
	})

	if err != nil {
		return nil, err
	}

	return &EthereumUpkeepTranscoder{
		client:     e.client,
		transcoder: instance.(*ethereum.UpkeepTranscoder),
		address:    address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperRegistrar(linkAddr string,
	registrarSettings KeeperRegistrarSettings) (KeeperRegistrar, error) {

	address, _, instance, err := e.client.DeployContract("KeeperRegistrar", func(
		opts *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployKeeperRegistrar(opts, backend, common.HexToAddress(linkAddr), registrarSettings.AutoApproveConfigType,
			registrarSettings.AutoApproveMaxAllowed, common.HexToAddress(registrarSettings.RegistryAddr), registrarSettings.MinLinkJuels)
	})

	if err != nil {
		return nil, err
	}

	return &EthereumKeeperRegistrar{
		client:    e.client,
		registrar: instance.(*ethereum.KeeperRegistrar),
		address:   address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperRegistry(
	opts *KeeperRegistryOpts,
) (KeeperRegistry, error) {
	switch opts.RegistryVersion {
	case ethereum.RegistryVersion_1_0, ethereum.RegistryVersion_1_1:
		address, _, instance, err := e.client.DeployContract("KeeperRegistry1_1", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return ethereum.DeployKeeperRegistry11(
				auth,
				backend,
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
				opts.Settings.PaymentPremiumPPB,
				opts.Settings.FlatFeeMicroLINK,
				opts.Settings.BlockCountPerTurn,
				opts.Settings.CheckGasLimit,
				opts.Settings.StalenessSeconds,
				opts.Settings.GasCeilingMultiplier,
				opts.Settings.FallbackGasPrice,
				opts.Settings.FallbackLinkPrice,
			)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			client:      e.client,
			version:     ethereum.RegistryVersion_1_1,
			registry1_1: instance.(*ethereum.KeeperRegistry11),
			registry1_2: nil,
			address:     address,
		}, err
	case ethereum.RegistryVersion_1_2:
		address, _, instance, err := e.client.DeployContract("KeeperRegistry", func(
			auth *bind.TransactOpts,
			backend bind.ContractBackend,
		) (common.Address, *types.Transaction, interface{}, error) {
			return ethereum.DeployKeeperRegistry(
				auth,
				backend,
				common.HexToAddress(opts.LinkAddr),
				common.HexToAddress(opts.ETHFeedAddr),
				common.HexToAddress(opts.GasFeedAddr),
				ethereum.Config{
					PaymentPremiumPPB:    opts.Settings.PaymentPremiumPPB,
					FlatFeeMicroLink:     opts.Settings.FlatFeeMicroLINK,
					BlockCountPerTurn:    opts.Settings.BlockCountPerTurn,
					CheckGasLimit:        opts.Settings.CheckGasLimit,
					StalenessSeconds:     opts.Settings.StalenessSeconds,
					GasCeilingMultiplier: opts.Settings.GasCeilingMultiplier,
					MinUpkeepSpend:       opts.Settings.MinUpkeepSpend,
					MaxPerformGas:        opts.Settings.MaxPerformGas,
					FallbackGasPrice:     opts.Settings.FallbackGasPrice,
					FallbackLinkPrice:    opts.Settings.FallbackLinkPrice,
					Transcoder:           common.HexToAddress(opts.TranscoderAddr),
					Registrar:            common.HexToAddress(opts.RegistrarAddr),
				},
			)
		})
		if err != nil {
			return nil, err
		}
		return &EthereumKeeperRegistry{
			client:      e.client,
			version:     ethereum.RegistryVersion_1_2,
			registry1_1: nil,
			registry1_2: instance.(*ethereum.KeeperRegistry),
			address:     address,
		}, err

	default:
		return nil, fmt.Errorf("keeper registry version %d is not supported", opts.RegistryVersion)
	}
}

func (e *EthereumContractDeployer) DeployKeeperConsumer(updateInterval *big.Int) (KeeperConsumer, error) {
	address, _, instance, err := e.client.DeployContract("KeeperConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployKeeperConsumer(auth, backend, updateInterval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumer{
		client:   e.client,
		consumer: instance.(*ethereum.KeeperConsumer),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepCounter(testRange *big.Int, interval *big.Int) (UpkeepCounter, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepCounter", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployUpkeepCounter(auth, backend, testRange, interval)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepCounter{
		client:   e.client,
		consumer: instance.(*ethereum.UpkeepCounter),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployUpkeepPerformCounterRestrictive(testRange *big.Int, averageEligibilityCadence *big.Int) (UpkeepPerformCounterRestrictive, error) {
	address, _, instance, err := e.client.DeployContract("UpkeepPerformCounterRestrictive", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployUpkeepPerformCounterRestrictive(auth, backend, testRange, averageEligibilityCadence)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumUpkeepPerformCounterRestrictive{
		client:   e.client,
		consumer: instance.(*ethereum.UpkeepPerformCounterRestrictive),
		address:  address,
	}, err
}

func (e *EthereumContractDeployer) DeployKeeperConsumerPerformance(
	testBlockRange,
	averageCadence,
	checkGasToBurn,
	performGasToBurn *big.Int,
) (KeeperConsumerPerformance, error) {
	address, _, instance, err := e.client.DeployContract("KeeperConsumerPerformance", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployKeeperConsumerPerformance(
			auth,
			backend,
			testBlockRange,
			averageCadence,
			checkGasToBurn,
			performGasToBurn,
		)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumKeeperConsumerPerformance{
		client:   e.client,
		consumer: instance.(*ethereum.KeeperConsumerPerformance),
		address:  address,
	}, err
}

// DeployBlockhashStore deploys blockhash store used with VRF contract
func (e *EthereumContractDeployer) DeployBlockhashStore() (BlockHashStore, error) {
	address, _, instance, err := e.client.DeployContract("BlockhashStore", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployBlockhashStore(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return &EthereumBlockhashStore{
		client:         e.client,
		blockHashStore: instance.(*ethereum.BlockhashStore),
		address:        address,
	}, err
}

// DeployVRFCoordinatorV2 deploys VRFV2 coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinatorV2(linkAddr string, bhsAddr string, linkEthFeedAddr string) (VRFCoordinatorV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinatorV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFCoordinatorV2(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr), common.HexToAddress(linkEthFeedAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinatorV2{
		client:      e.client,
		coordinator: instance.(*ethereum.VRFCoordinatorV2),
		address:     address,
	}, err
}

// DeployVRFCoordinator deploys VRF coordinator contract
func (e *EthereumContractDeployer) DeployVRFCoordinator(linkAddr string, bhsAddr string) (VRFCoordinator, error) {
	address, _, instance, err := e.client.DeployContract("VRFCoordinator", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFCoordinator(auth, backend, common.HexToAddress(linkAddr), common.HexToAddress(bhsAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFCoordinator{
		client:      e.client,
		coordinator: instance.(*ethereum.VRFCoordinator),
		address:     address,
	}, err
}

// DeployVRFConsumer deploys VRF consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumer(linkAddr string, coordinatorAddr string) (VRFConsumer, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumer", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFConsumer(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumer{
		client:   e.client,
		consumer: instance.(*ethereum.VRFConsumer),
		address:  address,
	}, err
}

// DeployVRFConsumerV2 deploys VRFv@ consumer contract
func (e *EthereumContractDeployer) DeployVRFConsumerV2(linkAddr string, coordinatorAddr string) (VRFConsumerV2, error) {
	address, _, instance, err := e.client.DeployContract("VRFConsumerV2", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return ethereum.DeployVRFConsumerV2(auth, backend, common.HexToAddress(coordinatorAddr), common.HexToAddress(linkAddr))
	})
	if err != nil {
		return nil, err
	}
	return &EthereumVRFConsumerV2{
		client:   e.client,
		consumer: instance.(*ethereum.VRFConsumerV2),
		address:  address,
	}, err
}
