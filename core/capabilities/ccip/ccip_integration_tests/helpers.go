package ccip_integration_tests

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-ccip/chainconfig"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccip_integration_tests/integrationhelpers"
	cctypes "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/types"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"

	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_config"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_multi_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/nonce_manager"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ocr3_config_encoder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/token_admin_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	kcr "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/keystone/generated/capabilities_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/link_token"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"

	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/stretchr/testify/require"
)

var (
	homeChainID                = chainsel.GETH_TESTNET.EvmChainID
	ccipSendRequestedTopic     = evm_2_evm_multi_onramp.EVM2EVMMultiOnRampCCIPSendRequested{}.Topic()
	commitReportAcceptedTopic  = evm_2_evm_multi_offramp.EVM2EVMMultiOffRampCommitReportAccepted{}.Topic()
	executionStateChangedTopic = evm_2_evm_multi_offramp.EVM2EVMMultiOffRampExecutionStateChanged{}.Topic()
)

const (
	CapabilityLabelledName = "ccip"
	CapabilityVersion      = "v1.0.0"
	NodeOperatorID         = 1

	// These constants drive what is set in the plugin offchain configs.
	FirstBlockAge                           = 8 * time.Hour
	RemoteGasPriceBatchWriteFrequency       = 30 * time.Minute
	BatchGasLimit                           = 6_500_000
	RelativeBoostPerWaitHour                = 1.5
	InflightCacheExpiry                     = 10 * time.Minute
	RootSnoozeTime                          = 30 * time.Minute
	BatchingStrategyID                      = 0
	DeltaProgress                           = 30 * time.Second
	DeltaResend                             = 10 * time.Second
	DeltaInitial                            = 20 * time.Second
	DeltaRound                              = 2 * time.Second
	DeltaGrace                              = 2 * time.Second
	DeltaCertifiedCommitRequest             = 10 * time.Second
	DeltaStage                              = 10 * time.Second
	Rmax                                    = 3
	MaxDurationQuery                        = 50 * time.Millisecond
	MaxDurationObservation                  = 5 * time.Second
	MaxDurationShouldAcceptAttestedReport   = 10 * time.Second
	MaxDurationShouldTransmitAcceptedReport = 10 * time.Second
)

func e18Mult(amount uint64) *big.Int {
	return new(big.Int).Mul(uBigInt(amount), uBigInt(1e18))
}

func uBigInt(i uint64) *big.Int {
	return new(big.Int).SetUint64(i)
}

type homeChain struct {
	backend            *backends.SimulatedBackend
	owner              *bind.TransactOpts
	chainID            uint64
	capabilityRegistry *kcr.CapabilitiesRegistry
	ccipConfig         *ccip_config.CCIPConfig
}

type onchainUniverse struct {
	backend            *backends.SimulatedBackend
	owner              *bind.TransactOpts
	chainID            uint64
	linkToken          *link_token.LinkToken
	weth               *weth9.WETH9
	router             *router.Router
	rmnProxy           *arm_proxy_contract.ARMProxyContract
	rmn                *mock_arm_contract.MockARMContract
	onramp             *evm_2_evm_multi_onramp.EVM2EVMMultiOnRamp
	offramp            *evm_2_evm_multi_offramp.EVM2EVMMultiOffRamp
	priceRegistry      *price_registry.PriceRegistry
	tokenAdminRegistry *token_admin_registry.TokenAdminRegistry
	nonceManager       *nonce_manager.NonceManager
	receiver           *maybe_revert_message_receiver.MaybeRevertMessageReceiver
}

type requestData struct {
	destChainSelector uint64
	receiverAddress   common.Address
	data              []byte
}

func (u *onchainUniverse) SendCCIPRequests(t *testing.T, requestDatas []requestData) {
	for _, reqData := range requestDatas {
		msg := router.ClientEVM2AnyMessage{
			Receiver:     common.LeftPadBytes(reqData.receiverAddress.Bytes(), 32),
			Data:         reqData.data,
			TokenAmounts: nil, // TODO: no tokens for now
			FeeToken:     u.weth.Address(),
			ExtraArgs:    nil, // TODO: no extra args for now, falls back to default
		}
		fee, err := u.router.GetFee(&bind.CallOpts{Context: testutils.Context(t)}, reqData.destChainSelector, msg)
		require.NoError(t, err)
		_, err = u.weth.Deposit(&bind.TransactOpts{
			From:   u.owner.From,
			Signer: u.owner.Signer,
			Value:  fee,
		})
		require.NoError(t, err)
		u.backend.Commit()
		_, err = u.weth.Approve(u.owner, u.router.Address(), fee)
		require.NoError(t, err)
		u.backend.Commit()

		t.Logf("Sending CCIP request from chain %d (selector %d) to chain selector %d",
			u.chainID, getSelector(u.chainID), reqData.destChainSelector)
		_, err = u.router.CcipSend(u.owner, reqData.destChainSelector, msg)
		require.NoError(t, err)
		u.backend.Commit()
	}
}

type chainBase struct {
	backend *backends.SimulatedBackend
	owner   *bind.TransactOpts
}

// createUniverses does the following:
// 1. Creates 1 home chain and `numChains`-1 non-home chains
// 2. Sets up home chain with the capability registry and the CCIP config contract
// 2. Deploys the CCIP contracts to all chains.
// 3. Sets up the initial configurations for the contracts on all chains.
// 4. Wires the chains together.
//
// Conceptually one universe is ONE chain with all the contracts deployed on it and all the dependencies initialized.
func createUniverses(
	t *testing.T,
	numChains int,
) (homeChainUni homeChain, universes map[uint64]onchainUniverse) {
	chains := createChains(t, numChains)

	homeChainBase, ok := chains[homeChainID]
	require.True(t, ok, "home chain backend not available")
	// Set up home chain first
	homeChainUniverse := setupHomeChain(t, homeChainBase.owner, homeChainBase.backend)

	// deploy the ccip contracts on all chains
	universes = make(map[uint64]onchainUniverse)
	for chainID, base := range chains {
		owner := base.owner
		backend := base.backend
		// deploy the CCIP contracts
		linkToken := deployLinkToken(t, owner, backend, chainID)
		rmn := deployMockARMContract(t, owner, backend, chainID)
		rmnProxy := deployARMProxyContract(t, owner, backend, rmn.Address(), chainID)
		weth := deployWETHContract(t, owner, backend, chainID)
		rout := deployRouter(t, owner, backend, weth.Address(), rmnProxy.Address(), chainID)
		priceRegistry := deployPriceRegistry(t, owner, backend, linkToken.Address(), weth.Address(), big.NewInt(1e18), chainID)
		tokenAdminRegistry := deployTokenAdminRegistry(t, owner, backend, chainID)
		nonceManager := deployNonceManager(t, owner, backend, chainID)

		// ======================================================================
		//							OnRamp
		// ======================================================================
		onRampAddr, _, _, err := evm_2_evm_multi_onramp.DeployEVM2EVMMultiOnRamp(
			owner,
			backend,
			evm_2_evm_multi_onramp.EVM2EVMMultiOnRampStaticConfig{
				ChainSelector:      getSelector(chainID),
				RmnProxy:           rmnProxy.Address(),
				NonceManager:       nonceManager.Address(),
				TokenAdminRegistry: tokenAdminRegistry.Address(),
			},
			evm_2_evm_multi_onramp.EVM2EVMMultiOnRampDynamicConfig{
				Router:        rout.Address(),
				PriceRegistry: priceRegistry.Address(),
				// `withdrawFeeTokens` onRamp function is not part of the message flow
				// so we can set this to any address
				FeeAggregator: testutils.NewAddress(),
			},
		)
		require.NoErrorf(t, err, "failed to deploy onramp on chain id %d", chainID)
		backend.Commit()
		onramp, err := evm_2_evm_multi_onramp.NewEVM2EVMMultiOnRamp(onRampAddr, backend)
		require.NoError(t, err)

		// ======================================================================
		//							OffRamp
		// ======================================================================
		offrampAddr, _, _, err := evm_2_evm_multi_offramp.DeployEVM2EVMMultiOffRamp(
			owner,
			backend,
			evm_2_evm_multi_offramp.EVM2EVMMultiOffRampStaticConfig{
				ChainSelector:      getSelector(chainID),
				RmnProxy:           rmnProxy.Address(),
				TokenAdminRegistry: tokenAdminRegistry.Address(),
				NonceManager:       nonceManager.Address(),
			},
			evm_2_evm_multi_offramp.EVM2EVMMultiOffRampDynamicConfig{
				Router:        rout.Address(),
				PriceRegistry: priceRegistry.Address(),
			},
			// Source chain configs will be set up later once we have all chains
			[]evm_2_evm_multi_offramp.EVM2EVMMultiOffRampSourceChainConfigArgs{},
		)
		require.NoErrorf(t, err, "failed to deploy offramp on chain id %d", chainID)
		backend.Commit()
		offramp, err := evm_2_evm_multi_offramp.NewEVM2EVMMultiOffRamp(offrampAddr, backend)
		require.NoError(t, err)

		receiverAddress, _, _, err := maybe_revert_message_receiver.DeployMaybeRevertMessageReceiver(
			owner,
			backend,
			false,
		)
		require.NoError(t, err, "failed to deploy MaybeRevertMessageReceiver on chain id %d", chainID)
		backend.Commit()
		receiver, err := maybe_revert_message_receiver.NewMaybeRevertMessageReceiver(receiverAddress, backend)
		require.NoError(t, err)

		universe := onchainUniverse{
			backend:            backend,
			owner:              owner,
			chainID:            chainID,
			linkToken:          linkToken,
			weth:               weth,
			router:             rout,
			rmnProxy:           rmnProxy,
			rmn:                rmn,
			onramp:             onramp,
			offramp:            offramp,
			priceRegistry:      priceRegistry,
			tokenAdminRegistry: tokenAdminRegistry,
			nonceManager:       nonceManager,
			receiver:           receiver,
		}
		// Set up the initial configurations for the contracts
		setupUniverseBasics(t, universe)

		universes[chainID] = universe
	}

	// Once we have all chains created and contracts deployed, we can set up the initial configurations and wire chains together
	connectUniverses(t, universes)

	// print out all contract addresses for debugging purposes
	for chainID, uni := range universes {
		t.Logf("Chain ID: %d\n Chain Selector: %d\n LinkToken: %s\n WETH: %s\n Router: %s\n RMNProxy: %s\n RMN: %s\n OnRamp: %s\n OffRamp: %s\n PriceRegistry: %s\n TokenAdminRegistry: %s\n NonceManager: %s\n",
			chainID,
			getSelector(chainID),
			uni.linkToken.Address().Hex(),
			uni.weth.Address().Hex(),
			uni.router.Address().Hex(),
			uni.rmnProxy.Address().Hex(),
			uni.rmn.Address().Hex(),
			uni.onramp.Address().Hex(),
			uni.offramp.Address().Hex(),
			uni.priceRegistry.Address().Hex(),
			uni.tokenAdminRegistry.Address().Hex(),
			uni.nonceManager.Address().Hex(),
		)
	}

	// print out topic hashes of relevant events for debugging purposes
	t.Logf("Topic hash of CommitReportAccepted: %s", commitReportAcceptedTopic.Hex())
	t.Logf("Topic hash of ExecutionStateChanged: %s", executionStateChangedTopic.Hex())
	t.Logf("Topic hash of CCIPSendRequested: %s", ccipSendRequestedTopic.Hex())

	return homeChainUniverse, universes
}

// Creates 1 home chain and `numChains`-1 non-home chains
func createChains(t *testing.T, numChains int) map[uint64]chainBase {
	chains := make(map[uint64]chainBase)

	homeChainOwner := testutils.MustNewSimTransactor(t)
	homeChainBackend := backends.NewSimulatedBackend(core.GenesisAlloc{
		homeChainOwner.From: core.GenesisAccount{
			Balance: assets.Ether(10_000).ToInt(),
		},
	}, 30e6)
	tweakChainTimestamp(t, homeChainBackend, FirstBlockAge)

	chains[homeChainID] = chainBase{
		owner:   homeChainOwner,
		backend: homeChainBackend,
	}

	for chainID := chainsel.TEST_90000001.EvmChainID; len(chains) < numChains && chainID < chainsel.TEST_90000020.EvmChainID; chainID++ {
		owner := testutils.MustNewSimTransactor(t)
		backend := backends.NewSimulatedBackend(core.GenesisAlloc{
			owner.From: core.GenesisAccount{
				Balance: assets.Ether(10_000).ToInt(),
			},
		}, 30e6)

		tweakChainTimestamp(t, backend, FirstBlockAge)

		chains[chainID] = chainBase{
			owner:   owner,
			backend: backend,
		}
	}

	return chains
}

// CCIP relies on block timestamps, but SimulatedBackend uses by default clock starting from 1970-01-01
// This trick is used to move the clock closer to the current time. We set first block to be X hours ago.
// Tests create plenty of transactions so this number can't be too low, every new block mined will tick the clock,
// if you mine more than "X hours" transactions, SimulatedBackend will panic because generated timestamps will be in the future.
func tweakChainTimestamp(t *testing.T, backend *backends.SimulatedBackend, tweak time.Duration) {
	blockTime := time.Unix(int64(backend.Blockchain().CurrentHeader().Time), 0)
	sinceBlockTime := time.Since(blockTime)
	diff := sinceBlockTime - tweak
	err := backend.AdjustTime(diff)
	require.NoError(t, err, "unable to adjust time on simulated chain")
	backend.Commit()
	backend.Commit()
}

func setupHomeChain(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend) homeChain {
	// deploy the capability registry on the home chain
	crAddress, _, _, err := kcr.DeployCapabilitiesRegistry(owner, backend)
	require.NoError(t, err, "failed to deploy capability registry on home chain")
	backend.Commit()

	capabilityRegistry, err := kcr.NewCapabilitiesRegistry(crAddress, backend)
	require.NoError(t, err)

	ccAddress, _, _, err := ccip_config.DeployCCIPConfig(owner, backend, crAddress)
	require.NoError(t, err)
	backend.Commit()

	capabilityConfig, err := ccip_config.NewCCIPConfig(ccAddress, backend)
	require.NoError(t, err)

	_, err = capabilityRegistry.AddCapabilities(owner, []kcr.CapabilitiesRegistryCapability{
		{
			LabelledName:          CapabilityLabelledName,
			Version:               CapabilityVersion,
			CapabilityType:        2, // consensus. not used (?)
			ResponseType:          0, // report. not used (?)
			ConfigurationContract: ccAddress,
		},
	})
	require.NoError(t, err, "failed to add capabilities to the capability registry")
	backend.Commit()

	// Add NodeOperator, for simplicity we'll add one NodeOperator only
	// First NodeOperator will have NodeOperatorId = 1
	_, err = capabilityRegistry.AddNodeOperators(owner, []kcr.CapabilitiesRegistryNodeOperator{
		{
			Admin: owner.From,
			Name:  "NodeOperator",
		},
	})
	require.NoError(t, err, "failed to add node operator to the capability registry")
	backend.Commit()

	return homeChain{
		backend:            backend,
		owner:              owner,
		chainID:            homeChainID,
		capabilityRegistry: capabilityRegistry,
		ccipConfig:         capabilityConfig,
	}
}

func sortP2PIDS(p2pIDs [][32]byte) {
	sort.Slice(p2pIDs, func(i, j int) bool {
		return bytes.Compare(p2pIDs[i][:], p2pIDs[j][:]) < 0
	})
}

func (h *homeChain) AddNodes(
	t *testing.T,
	p2pIDs [][32]byte,
	capabilityIDs [][32]byte,
) {
	// Need to sort, otherwise _checkIsValidUniqueSubset onChain will fail
	sortP2PIDS(p2pIDs)
	var nodeParams []kcr.CapabilitiesRegistryNodeParams
	for _, p2pID := range p2pIDs {
		nodeParam := kcr.CapabilitiesRegistryNodeParams{
			NodeOperatorId:      NodeOperatorID,
			Signer:              p2pID, // Not used in tests
			P2pId:               p2pID,
			HashedCapabilityIds: capabilityIDs,
		}
		nodeParams = append(nodeParams, nodeParam)
	}
	_, err := h.capabilityRegistry.AddNodes(h.owner, nodeParams)
	require.NoError(t, err, "failed to add node operator oracles")
	h.backend.Commit()
}

func AddChainConfig(
	t *testing.T,
	h homeChain,
	chainSelector uint64,
	p2pIDs [][32]byte,
	f uint8,
) ccip_config.CCIPConfigTypesChainConfigInfo {
	// Need to sort, otherwise _checkIsValidUniqueSubset onChain will fail
	sortP2PIDS(p2pIDs)
	// First Add ChainConfig that includes all p2pIDs as readers
	encodedExtraChainConfig, err := chainconfig.EncodeChainConfig(chainconfig.ChainConfig{
		GasPriceDeviationPPB:    ccipocr3.NewBigIntFromInt64(1000),
		DAGasPriceDeviationPPB:  ccipocr3.NewBigIntFromInt64(0),
		FinalityDepth:           10,
		OptimisticConfirmations: 1,
	})
	require.NoError(t, err)
	chainConfig := integrationhelpers.SetupConfigInfo(chainSelector, p2pIDs, f, encodedExtraChainConfig)
	inputConfig := []ccip_config.CCIPConfigTypesChainConfigInfo{
		chainConfig,
	}
	_, err = h.ccipConfig.ApplyChainConfigUpdates(h.owner, nil, inputConfig)
	require.NoError(t, err)
	h.backend.Commit()
	return chainConfig
}

func (h *homeChain) AddDON(
	t *testing.T,
	ccipCapabilityID [32]byte,
	chainSelector uint64,
	uni onchainUniverse,
	f uint8,
	bootstrapP2PID [32]byte,
	p2pIDs [][32]byte,
	oracles []confighelper2.OracleIdentityExtra,
) {
	// Get OCR3 Config from helper
	var schedule []int
	for range oracles {
		schedule = append(schedule, 1)
	}

	tabi, err := ocr3_config_encoder.IOCR3ConfigEncoderMetaData.GetAbi()
	require.NoError(t, err)

	// Add DON on capability registry contract
	var ocr3Configs []ocr3_config_encoder.CCIPConfigTypesOCR3Config
	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		var encodedOffchainConfig []byte
		var err2 error
		if pluginType == cctypes.PluginTypeCCIPCommit {
			encodedOffchainConfig, err2 = pluginconfig.EncodeCommitOffchainConfig(pluginconfig.CommitOffchainConfig{
				RemoteGasPriceBatchWriteFrequency: *commonconfig.MustNewDuration(RemoteGasPriceBatchWriteFrequency),
				// TODO: implement token price writes
				// TokenPriceBatchWriteFrequency:     *commonconfig.MustNewDuration(tokenPriceBatchWriteFrequency),
			})
			require.NoError(t, err2)
		} else {
			encodedOffchainConfig, err2 = pluginconfig.EncodeExecuteOffchainConfig(pluginconfig.ExecuteOffchainConfig{
				BatchGasLimit:             BatchGasLimit,
				RelativeBoostPerWaitHour:  RelativeBoostPerWaitHour,
				MessageVisibilityInterval: *commonconfig.MustNewDuration(FirstBlockAge),
				InflightCacheExpiry:       *commonconfig.MustNewDuration(InflightCacheExpiry),
				RootSnoozeTime:            *commonconfig.MustNewDuration(RootSnoozeTime),
				BatchingStrategyID:        BatchingStrategyID,
			})
			require.NoError(t, err2)
		}
		signers, transmitters, configF, _, offchainConfigVersion, offchainConfig, err2 := ocr3confighelper.ContractSetConfigArgsForTests(
			DeltaProgress,
			DeltaResend,
			DeltaInitial,
			DeltaRound,
			DeltaGrace,
			DeltaCertifiedCommitRequest,
			DeltaStage,
			Rmax,
			schedule,
			oracles,
			encodedOffchainConfig,
			MaxDurationQuery,
			MaxDurationObservation,
			MaxDurationShouldAcceptAttestedReport,
			MaxDurationShouldTransmitAcceptedReport,
			int(f),
			[]byte{}, // empty OnChainConfig
		)
		require.NoError(t, err2, "failed to create contract config")

		signersBytes := make([][]byte, len(signers))
		for i, signer := range signers {
			signersBytes[i] = signer
		}

		transmittersBytes := make([][]byte, len(transmitters))
		for i, transmitter := range transmitters {
			// anotherErr because linting doesn't want to shadow err
			parsed, anotherErr := common.ParseHexOrString(string(transmitter))
			require.NoError(t, anotherErr)
			transmittersBytes[i] = parsed
		}

		ocr3Configs = append(ocr3Configs, ocr3_config_encoder.CCIPConfigTypesOCR3Config{
			PluginType:            uint8(pluginType),
			ChainSelector:         chainSelector,
			F:                     configF,
			OffchainConfigVersion: offchainConfigVersion,
			OfframpAddress:        uni.offramp.Address().Bytes(),
			BootstrapP2PIds:       [][32]byte{bootstrapP2PID},
			P2pIds:                p2pIDs,
			Signers:               signersBytes,
			Transmitters:          transmittersBytes,
			OffchainConfig:        offchainConfig,
		})
	}

	encodedCall, err := tabi.Pack("exposeOCR3Config", ocr3Configs)
	require.NoError(t, err)

	// Trim first four bytes to remove function selector.
	encodedConfigs := encodedCall[4:]

	// commit so that we have an empty block to filter events from
	h.backend.Commit()

	_, err = h.capabilityRegistry.AddDON(h.owner, p2pIDs, []kcr.CapabilitiesRegistryCapabilityConfiguration{
		{
			CapabilityId: ccipCapabilityID,
			Config:       encodedConfigs,
		},
	}, false, false, f)
	require.NoError(t, err)
	h.backend.Commit()

	endBlock := h.backend.Blockchain().CurrentBlock().Number.Uint64()
	iter, err := h.capabilityRegistry.FilterConfigSet(&bind.FilterOpts{
		Start: h.backend.Blockchain().CurrentBlock().Number.Uint64() - 1,
		End:   &endBlock,
	})
	require.NoError(t, err, "failed to filter config set events")
	var donID uint32
	for iter.Next() {
		donID = iter.Event.DonId
		break
	}
	require.NotZero(t, donID, "failed to get donID from config set event")

	var signerAddresses []common.Address
	for _, oracle := range oracles {
		signerAddresses = append(signerAddresses, common.BytesToAddress(oracle.OnchainPublicKey))
	}

	var transmitterAddresses []common.Address
	for _, oracle := range oracles {
		transmitterAddresses = append(transmitterAddresses, common.HexToAddress(string(oracle.TransmitAccount)))
	}

	// get the config digest from the ccip config contract and set config on the offramp.
	var offrampOCR3Configs []evm_2_evm_multi_offramp.MultiOCR3BaseOCRConfigArgs
	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocrConfig, err1 := h.ccipConfig.GetOCRConfig(&bind.CallOpts{
			Context: testutils.Context(t),
		}, donID, uint8(pluginType))
		require.NoError(t, err1, "failed to get OCR3 config from ccip config contract")
		require.Len(t, ocrConfig, 1, "expected exactly one OCR3 config")
		offrampOCR3Configs = append(offrampOCR3Configs, evm_2_evm_multi_offramp.MultiOCR3BaseOCRConfigArgs{
			ConfigDigest:                   ocrConfig[0].ConfigDigest,
			OcrPluginType:                  uint8(pluginType),
			F:                              f,
			IsSignatureVerificationEnabled: pluginType == cctypes.PluginTypeCCIPCommit,
			Signers:                        signerAddresses,
			Transmitters:                   transmitterAddresses,
		})
	}

	uni.backend.Commit()

	_, err = uni.offramp.SetOCR3Configs(uni.owner, offrampOCR3Configs)
	require.NoError(t, err, "failed to set ocr3 configs on offramp")
	uni.backend.Commit()

	for _, pluginType := range []cctypes.PluginType{cctypes.PluginTypeCCIPCommit, cctypes.PluginTypeCCIPExec} {
		ocrConfig, err := uni.offramp.LatestConfigDetails(&bind.CallOpts{
			Context: testutils.Context(t),
		}, uint8(pluginType))
		require.NoError(t, err, "failed to get latest commit OCR3 config")
		require.Equalf(t, offrampOCR3Configs[pluginType].ConfigDigest, ocrConfig.ConfigInfo.ConfigDigest, "%s OCR3 config digest mismatch", pluginType.String())
		require.Equalf(t, offrampOCR3Configs[pluginType].F, ocrConfig.ConfigInfo.F, "%s OCR3 config F mismatch", pluginType.String())
		require.Equalf(t, offrampOCR3Configs[pluginType].IsSignatureVerificationEnabled, ocrConfig.ConfigInfo.IsSignatureVerificationEnabled, "%s OCR3 config signature verification mismatch", pluginType.String())
		if pluginType == cctypes.PluginTypeCCIPCommit {
			// only commit will set signers, exec doesn't need them.
			require.Equalf(t, offrampOCR3Configs[pluginType].Signers, ocrConfig.Signers, "%s OCR3 config signers mismatch", pluginType.String())
		}
		require.Equalf(t, offrampOCR3Configs[pluginType].Transmitters, ocrConfig.Transmitters, "%s OCR3 config transmitters mismatch", pluginType.String())
	}

	t.Logf("set ocr3 config on the offramp, signers: %+v, transmitters: %+v", signerAddresses, transmitterAddresses)
}

func connectUniverses(
	t *testing.T,
	universes map[uint64]onchainUniverse,
) {
	for _, uni := range universes {
		wireRouter(t, uni, universes)
		wirePriceRegistry(t, uni, universes)
		wireOffRamp(t, uni, universes)
		initRemoteChainsGasPrices(t, uni, universes)
	}
}

// setupUniverseBasics sets up the initial configurations for the CCIP contracts on a single chain.
// 1. Mint 1000 LINK to the owner
// 2. Set the price registry with local token prices
// 3. Authorize the onRamp and offRamp on the nonce manager
func setupUniverseBasics(t *testing.T, uni onchainUniverse) {
	// =============================================================================
	//			Universe specific  updates/configs
	//		These updates are specific to each universe and are set up here
	//      These updates don't depend on other chains
	// =============================================================================
	owner := uni.owner
	// =============================================================================
	//							Mint 1000 LINK to owner
	// =============================================================================
	_, err := uni.linkToken.GrantMintRole(owner, owner.From)
	require.NoError(t, err)
	_, err = uni.linkToken.Mint(owner, owner.From, e18Mult(1000))
	require.NoError(t, err)
	uni.backend.Commit()

	// =============================================================================
	//						Price updates for tokens
	//			These are the prices of the fee tokens of local chain in USD
	// =============================================================================
	tokenPriceUpdates := []price_registry.InternalTokenPriceUpdate{
		{
			SourceToken: uni.linkToken.Address(),
			UsdPerToken: e18Mult(20),
		},
		{
			SourceToken: uni.weth.Address(),
			UsdPerToken: e18Mult(4000),
		},
	}
	_, err = uni.priceRegistry.UpdatePrices(owner, price_registry.InternalPriceUpdates{
		TokenPriceUpdates: tokenPriceUpdates,
	})
	require.NoErrorf(t, err, "failed to update prices in price registry on chain id %d", uni.chainID)
	uni.backend.Commit()

	_, err = uni.priceRegistry.ApplyAuthorizedCallerUpdates(owner, price_registry.AuthorizedCallersAuthorizedCallerArgs{
		AddedCallers: []common.Address{
			uni.offramp.Address(),
		},
	})
	require.NoError(t, err, "failed to authorize offramp on price registry")
	uni.backend.Commit()

	// =============================================================================
	//						Authorize OnRamp & OffRamp on NonceManager
	//	Otherwise the onramp will not be able to call the nonceManager to get next Nonce
	// =============================================================================
	authorizedCallersAuthorizedCallerArgs := nonce_manager.AuthorizedCallersAuthorizedCallerArgs{
		AddedCallers: []common.Address{
			uni.onramp.Address(),
			uni.offramp.Address(),
		},
	}
	_, err = uni.nonceManager.ApplyAuthorizedCallerUpdates(owner, authorizedCallersAuthorizedCallerArgs)
	require.NoError(t, err)
	uni.backend.Commit()
}

// As we can't change router contract. The contract was expecting onRamp and offRamp per lane and not per chain
// In the new architecture we have only one onRamp and one offRamp per chain.
// hence we add the mapping for all remote chains to the onRamp/offRamp contract of the local chain
func wireRouter(t *testing.T, uni onchainUniverse, universes map[uint64]onchainUniverse) {
	owner := uni.owner
	var (
		routerOnrampUpdates  []router.RouterOnRamp
		routerOfframpUpdates []router.RouterOffRamp
	)
	for remoteChainID := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		routerOnrampUpdates = append(routerOnrampUpdates, router.RouterOnRamp{
			DestChainSelector: getSelector(remoteChainID),
			OnRamp:            uni.onramp.Address(),
		})
		routerOfframpUpdates = append(routerOfframpUpdates, router.RouterOffRamp{
			SourceChainSelector: getSelector(remoteChainID),
			OffRamp:             uni.offramp.Address(),
		})
	}
	_, err := uni.router.ApplyRampUpdates(owner, routerOnrampUpdates, []router.RouterOffRamp{}, routerOfframpUpdates)
	require.NoErrorf(t, err, "failed to apply ramp updates on router on chain id %d", uni.chainID)
	uni.backend.Commit()
}

// Setting OnRampDestChainConfigs
func wirePriceRegistry(t *testing.T, uni onchainUniverse, universes map[uint64]onchainUniverse) {
	owner := uni.owner
	var priceRegistryDestChainConfigArgs []price_registry.PriceRegistryDestChainConfigArgs
	for remoteChainID := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		priceRegistryDestChainConfigArgs = append(priceRegistryDestChainConfigArgs, price_registry.PriceRegistryDestChainConfigArgs{
			DestChainSelector: getSelector(remoteChainID),
			DestChainConfig:   defaultPriceRegistryDestChainConfig(t),
		})
	}
	_, err := uni.priceRegistry.ApplyDestChainConfigUpdates(owner, priceRegistryDestChainConfigArgs)
	require.NoErrorf(t, err, "failed to apply dest chain config updates on price registry on chain id %d", uni.chainID)
	uni.backend.Commit()
}

// Setting OffRampSourceChainConfigs
func wireOffRamp(t *testing.T, uni onchainUniverse, universes map[uint64]onchainUniverse) {
	owner := uni.owner
	var offrampSourceChainConfigArgs []evm_2_evm_multi_offramp.EVM2EVMMultiOffRampSourceChainConfigArgs
	for remoteChainID, remoteUniverse := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		offrampSourceChainConfigArgs = append(offrampSourceChainConfigArgs, evm_2_evm_multi_offramp.EVM2EVMMultiOffRampSourceChainConfigArgs{
			SourceChainSelector: getSelector(remoteChainID), // for each destination chain, add a source chain config
			IsEnabled:           true,
			OnRamp:              remoteUniverse.onramp.Address().Bytes(),
		})
	}
	_, err := uni.offramp.ApplySourceChainConfigUpdates(owner, offrampSourceChainConfigArgs)
	require.NoErrorf(t, err, "failed to apply source chain config updates on offramp on chain id %d", uni.chainID)
	uni.backend.Commit()
	for remoteChainID, remoteUniverse := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		sourceCfg, err2 := uni.offramp.GetSourceChainConfig(&bind.CallOpts{}, getSelector(remoteChainID))
		require.NoError(t, err2)
		require.True(t, sourceCfg.IsEnabled, "source chain config should be enabled")
		require.Equal(t, remoteUniverse.onramp.Address(), common.BytesToAddress(sourceCfg.OnRamp), "source chain config onRamp address mismatch")
	}
}

func getSelector(chainID uint64) uint64 {
	selector, err := chainsel.SelectorFromChainId(chainID)
	if err != nil {
		panic(err)
	}
	return selector
}

// initRemoteChainsGasPrices sets the gas prices for all chains except the local chain in the local price registry
func initRemoteChainsGasPrices(t *testing.T, uni onchainUniverse, universes map[uint64]onchainUniverse) {
	var gasPriceUpdates []price_registry.InternalGasPriceUpdate
	for remoteChainID := range universes {
		if remoteChainID == uni.chainID {
			continue
		}
		gasPriceUpdates = append(gasPriceUpdates,
			price_registry.InternalGasPriceUpdate{
				DestChainSelector: getSelector(remoteChainID),
				UsdPerUnitGas:     big.NewInt(2e12),
			},
		)
	}
	_, err := uni.priceRegistry.UpdatePrices(uni.owner, price_registry.InternalPriceUpdates{
		GasPriceUpdates: gasPriceUpdates,
	})
	require.NoError(t, err)
}

func defaultPriceRegistryDestChainConfig(t *testing.T) price_registry.PriceRegistryDestChainConfig {
	// https://github.com/smartcontractkit/ccip/blob/c4856b64bd766f1ddbaf5d13b42d3c4b12efde3a/contracts/src/v0.8/ccip/libraries/Internal.sol#L337-L337
	/*
		```Solidity
			// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
			bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
		```
	*/
	evmFamilySelector, err := hex.DecodeString("2812d52c")
	require.NoError(t, err)
	return price_registry.PriceRegistryDestChainConfig{
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      256,
		MaxPerMsgGasLimit:                 3_000_000,
		DestGasOverhead:                   50_000,
		DefaultTokenFeeUSDCents:           1,
		DestGasPerPayloadByte:             10,
		DestDataAvailabilityOverheadGas:   0,
		DestGasPerDataAvailabilityByte:    100,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       125_000,
		DefaultTokenDestBytesOverhead:     32,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            1,
		NetworkFeeUSDCents:                1,
		ChainFamilySelector:               [4]byte(evmFamilySelector),
	}
}

func deployLinkToken(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *link_token.LinkToken {
	linkAddr, _, _, err := link_token.DeployLinkToken(owner, backend)
	require.NoErrorf(t, err, "failed to deploy link token on chain id %d", chainID)
	backend.Commit()
	linkToken, err := link_token.NewLinkToken(linkAddr, backend)
	require.NoError(t, err)
	return linkToken
}

func deployMockARMContract(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *mock_arm_contract.MockARMContract {
	rmnAddr, _, _, err := mock_arm_contract.DeployMockARMContract(owner, backend)
	require.NoErrorf(t, err, "failed to deploy mock arm on chain id %d", chainID)
	backend.Commit()
	rmn, err := mock_arm_contract.NewMockARMContract(rmnAddr, backend)
	require.NoError(t, err)
	return rmn
}

func deployARMProxyContract(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, rmnAddr common.Address, chainID uint64) *arm_proxy_contract.ARMProxyContract {
	rmnProxyAddr, _, _, err := arm_proxy_contract.DeployARMProxyContract(owner, backend, rmnAddr)
	require.NoErrorf(t, err, "failed to deploy arm proxy on chain id %d", chainID)
	backend.Commit()
	rmnProxy, err := arm_proxy_contract.NewARMProxyContract(rmnProxyAddr, backend)
	require.NoError(t, err)
	return rmnProxy
}

func deployWETHContract(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *weth9.WETH9 {
	wethAddr, _, _, err := weth9.DeployWETH9(owner, backend)
	require.NoErrorf(t, err, "failed to deploy weth contract on chain id %d", chainID)
	backend.Commit()
	weth, err := weth9.NewWETH9(wethAddr, backend)
	require.NoError(t, err)
	return weth
}

func deployRouter(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, wethAddr, rmnProxyAddr common.Address, chainID uint64) *router.Router {
	routerAddr, _, _, err := router.DeployRouter(owner, backend, wethAddr, rmnProxyAddr)
	require.NoErrorf(t, err, "failed to deploy router on chain id %d", chainID)
	backend.Commit()
	rout, err := router.NewRouter(routerAddr, backend)
	require.NoError(t, err)
	return rout
}

func deployPriceRegistry(
	t *testing.T,
	owner *bind.TransactOpts,
	backend *backends.SimulatedBackend,
	linkAddr,
	wethAddr common.Address,
	maxFeeJuelsPerMsg *big.Int,
	chainID uint64,
) *price_registry.PriceRegistry {
	priceRegistryAddr, _, _, err := price_registry.DeployPriceRegistry(
		owner,
		backend,
		price_registry.PriceRegistryStaticConfig{
			MaxFeeJuelsPerMsg:  maxFeeJuelsPerMsg,
			LinkToken:          linkAddr,
			StalenessThreshold: 24 * 60 * 60, // 24 hours
		},
		[]common.Address{
			owner.From, // owner can update prices in this test
		}, // price updaters, will be set to offramp later
		[]common.Address{linkAddr, wethAddr}, // fee tokens
		// empty for now, need to fill in when testing token transfers
		[]price_registry.PriceRegistryTokenPriceFeedUpdate{},
		// empty for now, need to fill in when testing token transfers
		[]price_registry.PriceRegistryTokenTransferFeeConfigArgs{},
		[]price_registry.PriceRegistryPremiumMultiplierWeiPerEthArgs{
			{
				PremiumMultiplierWeiPerEth: 9e17, // 0.9 ETH
				Token:                      linkAddr,
			},
			{
				PremiumMultiplierWeiPerEth: 1e18,
				Token:                      wethAddr,
			},
		},
		// Destination chain configs will be set up later once we have all chains
		[]price_registry.PriceRegistryDestChainConfigArgs{},
	)
	require.NoErrorf(t, err, "failed to deploy price registry on chain id %d", chainID)
	backend.Commit()
	priceRegistry, err := price_registry.NewPriceRegistry(priceRegistryAddr, backend)
	require.NoError(t, err)
	return priceRegistry
}

func deployTokenAdminRegistry(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *token_admin_registry.TokenAdminRegistry {
	tarAddr, _, _, err := token_admin_registry.DeployTokenAdminRegistry(owner, backend)
	require.NoErrorf(t, err, "failed to deploy token admin registry on chain id %d", chainID)
	backend.Commit()
	tokenAdminRegistry, err := token_admin_registry.NewTokenAdminRegistry(tarAddr, backend)
	require.NoError(t, err)
	return tokenAdminRegistry
}

func deployNonceManager(t *testing.T, owner *bind.TransactOpts, backend *backends.SimulatedBackend, chainID uint64) *nonce_manager.NonceManager {
	nonceManagerAddr, _, _, err := nonce_manager.DeployNonceManager(owner, backend, []common.Address{owner.From})
	require.NoErrorf(t, err, "failed to deploy nonce_manager on chain id %d", chainID)
	backend.Commit()
	nonceManager, err := nonce_manager.NewNonceManager(nonceManagerAddr, backend)
	require.NoError(t, err)
	return nonceManager
}
