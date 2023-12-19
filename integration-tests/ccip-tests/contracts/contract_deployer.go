package contracts

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/curve25519"

	ocrconfighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	"github.com/smartcontractkit/chainlink/integration-tests/client"
	"github.com/smartcontractkit/chainlink/integration-tests/contracts"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/erc20"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

// CCIPContractsDeployer provides the implementations for deploying CCIP ETH contracts
type CCIPContractsDeployer struct {
	evmClient   blockchain.EVMClient
	EthDeployer *contracts.EthereumContractDeployer
}

// NewCCIPContractsDeployer returns an instance of a contract deployer for CCIP
func NewCCIPContractsDeployer(logger zerolog.Logger, bcClient blockchain.EVMClient) (*CCIPContractsDeployer, error) {
	return &CCIPContractsDeployer{
		evmClient:   bcClient,
		EthDeployer: contracts.NewEthereumContractDeployer(bcClient, logger),
	}, nil
}

func (e *CCIPContractsDeployer) Client() blockchain.EVMClient {
	return e.evmClient
}

func (e *CCIPContractsDeployer) DeployMultiCallContract() (common.Address, error) {
	multiCallABI, err := abi.JSON(strings.NewReader(MultiCallABI))
	if err != nil {
		return common.Address{}, err
	}
	address, tx, _, err := e.evmClient.DeployContract("MultiCall Contract", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		address, tx, contract, err := bind.DeployContract(auth, multiCallABI, common.FromHex(MultiCallBIN), backend)
		if err != nil {
			return common.Address{}, nil, nil, err
		}
		return address, tx, contract, err
	})
	if err != nil {
		return common.Address{}, err
	}
	r, err := bind.WaitMined(context.Background(), e.evmClient.DeployBackend(), tx)
	if err != nil {
		return common.Address{}, err
	}
	if r.Status != types.ReceiptStatusSuccessful {
		return common.Address{}, fmt.Errorf("deploy multicall failed")
	}
	return *address, nil
}

func (e *CCIPContractsDeployer) DeployLinkTokenContract() (*LinkToken, error) {
	address, _, instance, err := e.evmClient.DeployContract("Link Token", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return link_token_interface.DeployLinkToken(auth, backend)
	})

	if err != nil {
		return nil, err
	}
	return &LinkToken{
		client:     e.evmClient,
		instance:   instance.(*link_token_interface.LinkToken),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) DeployERC20TokenContract(deployerFn blockchain.ContractDeployer) (*ERC20Token, error) {
	address, _, _, err := e.evmClient.DeployContract("Custom ERC20 Token", deployerFn)
	if err != nil {
		return nil, err
	}
	err = e.evmClient.WaitForEvents()
	if err != nil {
		return nil, err
	}
	return e.NewERC20TokenContract(*address)
}

func (e *CCIPContractsDeployer) NewLinkTokenContract(addr common.Address) (*LinkToken, error) {
	token, err := link_token_interface.NewLinkToken(addr, e.evmClient.Backend())

	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "Link Token").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &LinkToken{
		client:     e.evmClient,
		instance:   token,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) NewERC20TokenContract(addr common.Address) (*ERC20Token, error) {
	token, err := erc20.NewERC20(addr, e.evmClient.Backend())

	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "ERC20 Token").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &ERC20Token{
		client:          e.evmClient,
		instance:        token,
		ContractAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) NewLockReleaseTokenPoolContract(addr common.Address) (
	*LockReleaseTokenPool,
	error,
) {
	pool, err := lock_release_token_pool.NewLockReleaseTokenPool(addr, e.evmClient.Backend())

	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "Native Token Pool").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &LockReleaseTokenPool{
		client:     e.evmClient,
		Instance:   pool,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) DeployLockReleaseTokenPoolContract(linkAddr string, armProxy common.Address) (
	*LockReleaseTokenPool,
	error,
) {
	log.Debug().Str("token", linkAddr).Msg("Deploying native token pool")
	token := common.HexToAddress(linkAddr)
	address, _, instance, err := e.evmClient.DeployContract("Native Token Pool", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return lock_release_token_pool.DeployLockReleaseTokenPool(
			auth,
			backend,
			token,
			[]common.Address{},
			armProxy,
			true)
	})

	if err != nil {
		return nil, err
	}
	return &LockReleaseTokenPool{
		client:     e.evmClient,
		Instance:   instance.(*lock_release_token_pool.LockReleaseTokenPool),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) DeployMockARMContract() (*common.Address, error) {
	address, _, _, err := e.evmClient.DeployContract("Mock ARM Contract", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return mock_arm_contract.DeployMockARMContract(auth, backend)
	})
	return address, err
}

func (e *CCIPContractsDeployer) NewARMContract(addr common.Address) (*ARM, error) {
	arm, err := arm_contract.NewARMContract(addr, e.evmClient.Backend())
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "Mock ARM Contract").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")

	return &ARM{
		client:     e.evmClient,
		Instance:   arm,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) NewCommitStore(addr common.Address) (
	*CommitStore,
	error,
) {
	ins, err := commit_store.NewCommitStore(addr, e.evmClient.Backend())
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "CommitStore").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &CommitStore{
		client:     e.evmClient,
		Instance:   ins,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) DeployCommitStore(sourceChainSelector, destChainSelector uint64, onRamp common.Address, armProxy common.Address) (*CommitStore, error) {
	address, _, instance, err := e.evmClient.DeployContract("CommitStore Contract", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return commit_store.DeployCommitStore(
			auth,
			backend,
			commit_store.CommitStoreStaticConfig{
				ChainSelector:       destChainSelector,
				SourceChainSelector: sourceChainSelector,
				OnRamp:              onRamp,
				ArmProxy:            armProxy,
			},
		)
	})
	if err != nil {
		return nil, err
	}
	return &CommitStore{
		client:     e.evmClient,
		Instance:   instance.(*commit_store.CommitStore),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) DeployReceiverDapp(revert bool) (
	*ReceiverDapp,
	error,
) {
	address, _, instance, err := e.evmClient.DeployContract("ReceiverDapp", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(auth, backend, revert)
	})
	if err != nil {
		return nil, err
	}
	return &ReceiverDapp{
		client:     e.evmClient,
		instance:   instance.(*maybe_revert_message_receiver.MaybeRevertMessageReceiver),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) NewReceiverDapp(addr common.Address) (
	*ReceiverDapp,
	error,
) {
	ins, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(addr, e.evmClient.Backend())
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "ReceiverDapp").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &ReceiverDapp{
		client:     e.evmClient,
		instance:   ins,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) DeployRouter(wrappedNative common.Address, armAddress common.Address) (
	*Router,
	error,
) {
	address, _, instance, err := e.evmClient.DeployContract("Router", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return router.DeployRouter(auth, backend, wrappedNative, armAddress)
	})
	if err != nil {
		return nil, err
	}
	return &Router{
		client:     e.evmClient,
		Instance:   instance.(*router.Router),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) NewRouter(addr common.Address) (
	*Router,
	error,
) {
	r, err := router.NewRouter(addr, e.evmClient.Backend())
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "Router").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	if err != nil {
		return nil, err
	}
	return &Router{
		client:     e.evmClient,
		Instance:   r,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) NewPriceRegistry(addr common.Address) (
	*PriceRegistry,
	error,
) {
	ins, err := price_registry.NewPriceRegistry(addr, e.evmClient.Backend())
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "PriceRegistry").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &PriceRegistry{
		client:     e.evmClient,
		Instance:   ins,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) DeployPriceRegistry(tokens []common.Address) (*PriceRegistry, error) {
	address, _, instance, err := e.evmClient.DeployContract("PriceRegistry", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return price_registry.DeployPriceRegistry(auth, backend, nil, tokens, 60*60*24*14)
	})
	if err != nil {
		return nil, err
	}
	return &PriceRegistry{
		client:     e.evmClient,
		Instance:   instance.(*price_registry.PriceRegistry),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) NewOnRamp(addr common.Address) (
	*OnRamp,
	error,
) {
	ins, err := evm_2_evm_onramp.NewEVM2EVMOnRamp(addr, e.evmClient.Backend())
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "OnRamp").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &OnRamp{
		client:     e.evmClient,
		Instance:   ins,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) DeployOnRamp(
	sourceChainSelector, destChainSelector uint64,
	tokensAndPools []evm_2_evm_onramp.InternalPoolUpdate,
	arm, router, priceRegistry common.Address,
	opts RateLimiterConfig,
	feeTokenConfig []evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs,
	tokenTransferFeeConfig []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs,
	linkTokenAddress common.Address,
) (
	*OnRamp,
	error,
) {
	address, _, instance, err := e.evmClient.DeployContract("OnRamp", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return evm_2_evm_onramp.DeployEVM2EVMOnRamp(
			auth,
			backend,
			evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
				LinkToken:         linkTokenAddress,
				ChainSelector:     sourceChainSelector, // source chain id
				DestChainSelector: destChainSelector,   // destinationChainSelector
				DefaultTxGasLimit: 200_000,
				MaxNopFeesJuels:   big.NewInt(0).Mul(big.NewInt(100_000_000), big.NewInt(1e18)),
				PrevOnRamp:        common.HexToAddress(""),
				ArmProxy:          arm,
			},
			evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
				Router:                            router,
				MaxNumberOfTokensPerMsg:           50,
				DestGasOverhead:                   350_000,
				DestGasPerPayloadByte:             16,
				DestDataAvailabilityOverheadGas:   33_596,
				DestGasPerDataAvailabilityByte:    16,
				DestDataAvailabilityMultiplierBps: 6840, // 0.684
				PriceRegistry:                     priceRegistry,
				MaxDataBytes:                      50000,
				MaxPerMsgGasLimit:                 4_000_000,
			},
			tokensAndPools,
			evm_2_evm_onramp.RateLimiterConfig{
				Capacity: opts.Capacity,
				Rate:     opts.Rate,
			},
			feeTokenConfig,
			tokenTransferFeeConfig,
			[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
		)
	})
	if err != nil {
		return nil, err
	}
	return &OnRamp{
		client:     e.evmClient,
		Instance:   instance.(*evm_2_evm_onramp.EVM2EVMOnRamp),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) NewOffRamp(addr common.Address) (
	*OffRamp,
	error,
) {
	ins, err := evm_2_evm_offramp.NewEVM2EVMOffRamp(addr, e.evmClient.Backend())
	log.Info().
		Str("Contract Address", addr.Hex()).
		Str("Contract Name", "OffRamp").
		Str("From", e.evmClient.GetDefaultWallet().Address()).
		Str("Network Name", e.evmClient.GetNetworkConfig().Name).
		Msg("New contract")
	return &OffRamp{
		client:     e.evmClient,
		Instance:   ins,
		EthAddress: addr,
	}, err
}

func (e *CCIPContractsDeployer) DeployOffRamp(sourceChainSelector, destChainSelector uint64, commitStore, onRamp common.Address, sourceToken, pools []common.Address, opts RateLimiterConfig, armProxy common.Address) (*OffRamp, error) {
	address, _, instance, err := e.evmClient.DeployContract("OffRamp Contract", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return evm_2_evm_offramp.DeployEVM2EVMOffRamp(
			auth,
			backend,
			evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
				CommitStore:         commitStore,
				ChainSelector:       destChainSelector,
				SourceChainSelector: sourceChainSelector,
				OnRamp:              onRamp,
				PrevOffRamp:         common.HexToAddress(""),
				ArmProxy:            armProxy,
			},
			sourceToken,
			pools,
			evm_2_evm_offramp.RateLimiterConfig{
				IsEnabled: true,
				Capacity:  opts.Capacity,
				Rate:      opts.Rate,
			},
		)
	})
	if err != nil {
		return nil, err
	}
	return &OffRamp{
		client:     e.evmClient,
		Instance:   instance.(*evm_2_evm_offramp.EVM2EVMOffRamp),
		EthAddress: *address,
	}, err
}

func (e *CCIPContractsDeployer) DeployWrappedNative() (*common.Address, error) {
	address, _, _, err := e.evmClient.DeployContract("WrappedNative", func(
		auth *bind.TransactOpts,
		backend bind.ContractBackend,
	) (common.Address, *types.Transaction, interface{}, error) {
		return weth9.DeployWETH9(auth, backend)
	})
	if err != nil {
		return nil, err
	}
	return address, err
}

var OCR2ParamsForCommit = contracts.OffChainAggregatorV2Config{
	DeltaProgress:                           2 * time.Minute,
	DeltaResend:                             5 * time.Second,
	DeltaRound:                              75 * time.Second,
	DeltaGrace:                              5 * time.Second,
	MaxDurationQuery:                        100 * time.Millisecond,
	MaxDurationObservation:                  35 * time.Second,
	MaxDurationReport:                       10 * time.Second,
	MaxDurationShouldAcceptFinalizedReport:  5 * time.Second,
	MaxDurationShouldTransmitAcceptedReport: 10 * time.Second,
}

var OCR2ParamsForExec = contracts.OffChainAggregatorV2Config{
	DeltaProgress:                           100 * time.Second,
	DeltaResend:                             5 * time.Second,
	DeltaRound:                              40 * time.Second,
	DeltaGrace:                              5 * time.Second,
	MaxDurationQuery:                        100 * time.Millisecond,
	MaxDurationObservation:                  20 * time.Second,
	MaxDurationReport:                       8 * time.Second,
	MaxDurationShouldAcceptFinalizedReport:  5 * time.Second,
	MaxDurationShouldTransmitAcceptedReport: 8 * time.Second,
}

func OffChainAggregatorV2ConfigWithNodes(numberNodes int, inflightExpiry time.Duration, cfg contracts.OffChainAggregatorV2Config) contracts.OffChainAggregatorV2Config {
	if numberNodes <= 4 {
		log.Err(fmt.Errorf("insufficient number of nodes (%d) supplied for OCR, need at least 5", numberNodes)).
			Int("Number Chainlink Nodes", numberNodes).
			Msg("You likely need more chainlink nodes to properly configure OCR, try 5 or more.")
	}
	s := make([]int, 0)
	for i := 0; i < numberNodes; i++ {
		s = append(s, 1)
	}
	faultyNodes := 0
	if numberNodes > 1 {
		faultyNodes = (numberNodes - 1) / 3
	}
	if faultyNodes == 0 {
		faultyNodes = 1
	}
	return contracts.OffChainAggregatorV2Config{
		DeltaProgress:                           cfg.DeltaProgress,
		DeltaResend:                             cfg.DeltaResend,
		DeltaRound:                              cfg.DeltaRound,
		DeltaGrace:                              cfg.DeltaGrace,
		DeltaStage:                              inflightExpiry,
		RMax:                                    3,
		S:                                       s,
		F:                                       faultyNodes,
		Oracles:                                 []ocrconfighelper2.OracleIdentityExtra{},
		MaxDurationQuery:                        cfg.MaxDurationQuery,
		MaxDurationObservation:                  cfg.MaxDurationObservation,
		MaxDurationReport:                       cfg.MaxDurationReport,
		MaxDurationShouldAcceptFinalizedReport:  cfg.MaxDurationShouldAcceptFinalizedReport,
		MaxDurationShouldTransmitAcceptedReport: cfg.MaxDurationShouldTransmitAcceptedReport,
		OnchainConfig:                           []byte{},
	}
}

func stripKeyPrefix(key string) string {
	chunks := strings.Split(key, "_")
	if len(chunks) == 3 {
		return chunks[2]
	}
	return key
}

func NewOffChainAggregatorV2ConfigForCCIPPlugin[T ccipconfig.OffchainConfig](
	nodes []*client.CLNodesWithKeys,
	offchainCfg T,
	onchainCfg abihelpers.AbiDefined,
	ocr2Params contracts.OffChainAggregatorV2Config,
	inflightExpiry time.Duration,
) (
	signers []common.Address,
	transmitters []common.Address,
	f_ uint8,
	onchainConfig_ []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
	err error,
) {
	oracleIdentities := make([]ocrconfighelper2.OracleIdentityExtra, 0)
	ocrConfig := OffChainAggregatorV2ConfigWithNodes(len(nodes), inflightExpiry, ocr2Params)
	var onChainKeys []ocrtypes2.OnchainPublicKey
	for i, nodeWithKeys := range nodes {
		ocr2Key := nodeWithKeys.KeysBundle.OCR2Key.Data
		offChainPubKeyTemp, err := hex.DecodeString(stripKeyPrefix(ocr2Key.Attributes.OffChainPublicKey))
		if err != nil {
			return nil, nil, 0, nil, 0, nil, err
		}
		formattedOnChainPubKey := stripKeyPrefix(ocr2Key.Attributes.OnChainPublicKey)
		cfgPubKeyTemp, err := hex.DecodeString(stripKeyPrefix(ocr2Key.Attributes.ConfigPublicKey))
		if err != nil {
			return nil, nil, 0, nil, 0, nil, err
		}
		cfgPubKeyBytes := [ed25519.PublicKeySize]byte{}
		copy(cfgPubKeyBytes[:], cfgPubKeyTemp)
		offChainPubKey := [curve25519.PointSize]byte{}
		copy(offChainPubKey[:], offChainPubKeyTemp)
		ethAddress := nodeWithKeys.KeysBundle.EthAddress
		p2pKeys := nodeWithKeys.KeysBundle.P2PKeys
		peerID := p2pKeys.Data[0].Attributes.PeerID
		oracleIdentities = append(oracleIdentities, ocrconfighelper2.OracleIdentityExtra{
			OracleIdentity: ocrconfighelper2.OracleIdentity{
				OffchainPublicKey: offChainPubKey,
				OnchainPublicKey:  common.HexToAddress(formattedOnChainPubKey).Bytes(),
				PeerID:            peerID,
				TransmitAccount:   ocrtypes2.Account(ethAddress),
			},
			ConfigEncryptionPublicKey: cfgPubKeyBytes,
		})
		onChainKeys = append(onChainKeys, oracleIdentities[i].OnchainPublicKey)
		transmitters = append(transmitters, common.HexToAddress(ethAddress))
	}
	signers, err = evm.OnchainPublicKeyToAddress(onChainKeys)
	if err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}
	ocrConfig.Oracles = oracleIdentities
	ocrConfig.ReportingPluginConfig, err = ccipconfig.EncodeOffchainConfig(offchainCfg)
	if err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}
	ocrConfig.OnchainConfig, err = abihelpers.EncodeAbiStruct(onchainCfg)
	if err != nil {
		return nil, nil, 0, nil, 0, nil, err
	}

	_, _, f_, onchainConfig_, offchainConfigVersion, offchainConfig, err = ocrconfighelper2.ContractSetConfigArgsForTests(
		ocrConfig.DeltaProgress,
		ocrConfig.DeltaResend,
		ocrConfig.DeltaRound,
		ocrConfig.DeltaGrace,
		ocrConfig.DeltaStage,
		ocrConfig.RMax,
		ocrConfig.S,
		ocrConfig.Oracles,
		ocrConfig.ReportingPluginConfig,
		ocrConfig.MaxDurationQuery,
		ocrConfig.MaxDurationObservation,
		ocrConfig.MaxDurationReport,
		ocrConfig.MaxDurationShouldAcceptFinalizedReport,
		ocrConfig.MaxDurationShouldTransmitAcceptedReport,
		ocrConfig.F,
		ocrConfig.OnchainConfig,
	)
	return
}
