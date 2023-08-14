package rhea

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/shared"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/ping_pong_demo"
)

const (
	RATE_LIMIT_CAPACITY_DOLLAR = 1e5
	// Rate = $100k per 10minute, Capacity = $100k
	RATE_LIMIT_RATE_DOLLAR = 167
	MAX_DATA_SIZE          = 50_000
	MAX_TOKEN_LENGTH       = 5
	MAX_TX_GAS_LIMIT       = 4_000_000
	// Network fee and min/max token transfer fee are multiples of $0.01 USD
	NETWORK_FEE            = 1_00
	MIN_TOKEN_TRANSFER_FEE = 1_00
	MAX_TOKEN_TRANSFER_FEE = 5000_00
	GAS_MULTIPLIER         = 1e18
	// Paying with LINK tokens carries a 10% discount on premium fees
	LINK_PREMIUM_MULTIPLIER    = 9e17
	DEFAULT_PREMIUM_MULTIPLIER = 1e18
	// For testnet we set this high to maximize self-healing
	// (e.g. for some reason a message doesn't get executed or is a strict seq failure and blocks a sequence
	// it may take a couple days to unblock)
	PERMISSIONLESS_EXEC_THRESHOLD_SEC = 7 * 24 * 3600
	DEST_GAS_OVERHEAD                 = 350_000
	DEST_GAS_PER_PAYLOAD_BYTE         = 16
	DEFAULT_GAS_LIMIT                 = 200_000
	MAX_NOP_FEES_LINK                 = 100_000_000
)

// DeployLanes will deploy all source and Destination chain contracts using the
// owner key. Only run this of the currently deployed contracts are outdated or
// when initializing a new chain.
func DeployLanes(t *testing.T, source *EvmDeploymentConfig, destination *EvmDeploymentConfig) {
	if !source.LaneConfig.DeploySettings.DeployLane && !source.LaneConfig.DeploySettings.DeployPingPongDapp &&
		!destination.LaneConfig.DeploySettings.DeployLane && !destination.LaneConfig.DeploySettings.DeployPingPongDapp {
		t.Log("No deployment needed")
		return
	}
	sourceChainSelector, destChainSelector := GetCCIPChainSelector(source.ChainConfig.EvmChainId), GetCCIPChainSelector(destination.ChainConfig.EvmChainId)

	// After running this code please update the configuration to reflect the newly
	// deployed contract addresses.
	// Deploy onRamps on both chains
	deploySourceContracts(t, source, destChainSelector, destination.ChainConfig.SupportedTokens)
	deploySourceContracts(t, destination, sourceChainSelector, source.ChainConfig.SupportedTokens)

	// Deploy commitStores and offRamps on both chains
	deployDestinationContracts(t, destination, sourceChainSelector, source.LaneConfig.OnRamp, source.ChainConfig.SupportedTokens)
	deployDestinationContracts(t, source, destChainSelector, destination.LaneConfig.OnRamp, destination.ChainConfig.SupportedTokens)

	setPriceRegistryPrices(t, source, destChainSelector)
	setPriceRegistryPrices(t, destination, sourceChainSelector)

	DeployPingPongDapps(t, source, destination)

	UpdateDeployedAt(t, source, destination)
}

func DeployUpgradeLanes(t *testing.T, source *EvmDeploymentConfig, destination *EvmDeploymentConfig) {
	if !source.UpgradeLaneConfig.DeploySettings.DeployLane && !destination.UpgradeLaneConfig.DeploySettings.DeployLane {
		t.Log("No deployment needed")
		return
	}
	sourceChainSelector, destChainSelector := GetCCIPChainSelector(source.ChainConfig.EvmChainId), GetCCIPChainSelector(destination.ChainConfig.EvmChainId)

	// After running this code please update the configuration to reflect the newly
	// deployed contract addresses.
	// Deploy onRamps on both chains
	deployUpgradeSourceContracts(t, source, destChainSelector, destination.ChainConfig.SupportedTokens)
	deployUpgradeSourceContracts(t, destination, sourceChainSelector, source.ChainConfig.SupportedTokens)

	// Deploy commitStores and offRamps on both chains
	deployUpgradeDestinationContracts(t, destination, sourceChainSelector, source.UpgradeLaneConfig.OnRamp, source.ChainConfig.SupportedTokens)
	deployUpgradeDestinationContracts(t, source, destChainSelector, destination.UpgradeLaneConfig.OnRamp, destination.ChainConfig.SupportedTokens)

	setPriceRegistryPrices(t, source, destChainSelector)
	setPriceRegistryPrices(t, destination, sourceChainSelector)

	DeployUpgradePingPongDapps(t, source, destination)

	UpdateDeployedAtUpgradeLane(t, source, destination)
}

func EnableUpgradeOffRamps(t *testing.T, source *EvmDeploymentConfig, destination *EvmDeploymentConfig, ocr2Update func(*EvmDeploymentConfig, *EvmDeploymentConfig)) {
	sourceChainSelector, destChainSelector := GetCCIPChainSelector(source.ChainConfig.EvmChainId), GetCCIPChainSelector(destination.ChainConfig.EvmChainId)

	// Attach Upgrade OffRamp to the current Router
	attachUpgradeOffRampsToRouter(t, source, destChainSelector)
	attachUpgradeOffRampsToRouter(t, destination, sourceChainSelector)

	// Allowing price registries to be used by new commit stores
	setPriceRegistryUpdater(t, source)
	setPriceRegistryUpdater(t, destination)

	// Update OCR2 configs with proper Router addresses
	ocr2Update(source, destination)
	ocr2Update(destination, source)
}

func EnableUpgradeOnRamps(t *testing.T, source *EvmDeploymentConfig, destination *EvmDeploymentConfig, dynamicConfig func(src *EvmDeploymentConfig, dst *EvmDeploymentConfig)) {
	// Update DynamicConfig to OnRamps (via setDynamicConfig) - sets router address to the new one
	dynamicConfig(source, destination)
	dynamicConfig(destination, source)

	sourceChainSelector, destChainSelector := GetCCIPChainSelector(source.ChainConfig.EvmChainId), GetCCIPChainSelector(destination.ChainConfig.EvmChainId)
	// Attach Upgrade OnRamp to the current Router
	attachUpgradeOnRampsToRouter(t, source, destChainSelector)
	attachUpgradeOnRampsToRouter(t, destination, sourceChainSelector)
}

func deploySourceContracts(t *testing.T, source *EvmDeploymentConfig, destChainSelector uint64, destSupportedTokens map[Token]EVMBridgedToken) {
	if source.LaneConfig.DeploySettings.DeployLane {
		// Updates source.OnRamp if any new contracts are deployed
		deployOnRamp(t, source.OnlyEvmConfig(), &source.LaneConfig, destChainSelector, destSupportedTokens, common.HexToAddress(""))
		setOnRampRouter(t, source.ChainConfig.Router, source.LaneConfig.OnRamp, source, destChainSelector)
	}

	// Skip if we reuse both the onRamp and the token pools
	if source.LaneConfig.DeploySettings.DeployLane || source.ChainConfig.DeploySettings.DeployTokenPools {
		setOnRampOnTokenPools(t, source, source.LaneConfig.OnRamp)
	}
	source.Logger.Infof("%s contracts deployed as source chain", ccip.ChainName(int64(source.ChainConfig.EvmChainId)))
}

func deployUpgradeSourceContracts(t *testing.T, source *EvmDeploymentConfig, destChainSelector uint64, destSupportedTokens map[Token]EVMBridgedToken) {
	sourceUpgradeLane := &source.UpgradeLaneConfig

	if sourceUpgradeLane.DeploySettings.DeployLane {
		// Updates source.OnRamp if any new contracts are deployed
		deployOnRamp(t, source.OnlyEvmConfig(), sourceUpgradeLane, destChainSelector, destSupportedTokens, source.LaneConfig.OnRamp)
		setOnRampRouter(t, source.ChainConfig.UpgradeRouter, sourceUpgradeLane.OnRamp, source, destChainSelector)
	}

	// Skip if we reuse both the onRamp and the token pools
	if sourceUpgradeLane.DeploySettings.DeployLane || source.ChainConfig.DeploySettings.DeployTokenPools {
		setOnRampOnTokenPools(t, source, sourceUpgradeLane.OnRamp)
	}
	source.Logger.Infof("%s contracts deployed as source chain", ccip.ChainName(int64(source.ChainConfig.EvmChainId)))
}

func deployDestinationContracts(t *testing.T, client *EvmDeploymentConfig, sourceChainSelector uint64, onRamp common.Address, supportedTokens map[Token]EVMBridgedToken) {
	// Updates destClient.LaneConfig.CommitStore if any new contracts are deployed
	deployCommitStore(t, client.OnlyEvmConfig(), &client.LaneConfig, sourceChainSelector, onRamp)

	if client.LaneConfig.DeploySettings.DeployLane || client.ChainConfig.DeploySettings.DeployPriceRegistry {
		setPriceRegistryUpdater(t, client)
	}

	if client.LaneConfig.DeploySettings.DeployLane {
		// Updates destClient.LaneConfig.OffRamp if any new contracts are deployed
		deployOffRamp(t, client.OnlyEvmConfig(), sourceChainSelector, supportedTokens, &client.LaneConfig, onRamp, common.HexToAddress(""))
	}

	if client.LaneConfig.DeploySettings.DeployLane || client.ChainConfig.DeploySettings.DeployTokenPools {
		setOffRampOnTokenPools(t, client.OnlyEvmConfig(), &client.LaneConfig)
	}

	if client.LaneConfig.DeploySettings.DeployLane || client.ChainConfig.DeploySettings.DeployRouter {
		setOffRampOnRouter(t, client.ChainConfig.Router, client.LaneConfig.OffRamp, client, sourceChainSelector)
	}

	client.Logger.Infof("%s contracts fully deployed as destination chain", ccip.ChainName(int64(client.ChainConfig.EvmChainId)))
}

func deployUpgradeDestinationContracts(t *testing.T, client *EvmDeploymentConfig, sourceChainSelector uint64, onRamp common.Address, supportedTokens map[Token]EVMBridgedToken) {
	laneConfig := &client.UpgradeLaneConfig
	deployCommitStore(t, client.OnlyEvmConfig(), laneConfig, sourceChainSelector, onRamp)

	if laneConfig.DeploySettings.DeployLane || client.ChainConfig.DeploySettings.DeployPriceRegistry {
		setPriceRegistryUpdater(t, client)
	}

	if laneConfig.DeploySettings.DeployLane {
		// Updates destClient.LaneConfig.OffRamp if any new contracts are deployed
		deployOffRamp(t, client.OnlyEvmConfig(), sourceChainSelector, supportedTokens, laneConfig, onRamp, client.LaneConfig.OffRamp)
	}

	if laneConfig.DeploySettings.DeployLane || client.ChainConfig.DeploySettings.DeployTokenPools {
		setOffRampOnTokenPools(t, client.OnlyEvmConfig(), laneConfig)
	}

	if laneConfig.DeploySettings.DeployLane || client.ChainConfig.DeploySettings.DeployRouter {
		setOffRampOnRouter(t, client.ChainConfig.UpgradeRouter, laneConfig.OffRamp, client, sourceChainSelector)
	}

	client.Logger.Infof("%s contracts fully deployed as destination chain", ccip.ChainName(int64(client.ChainConfig.EvmChainId)))
}

func deployOnRamp(t *testing.T, client EvmConfig, laneConfig *EVMLaneConfig, destChainSelector uint64, destSupportedTokens map[Token]EVMBridgedToken, prevOnRamp common.Address) {
	if !laneConfig.DeploySettings.DeployLane {
		client.Logger.Infof("Skipping OnRamp deployment, using onRamp on %s", laneConfig.OnRamp)
		return
	}

	var tokensAndPools []evm_2_evm_onramp.InternalPoolUpdate
	var tokenTransferFeeConfig []evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs
	for token, tokenConfig := range client.ChainConfig.SupportedTokens {
		// Don't add tokens that are only configured for paying fees
		if tokenConfig.TokenPoolType == FeeTokenOnly {
			continue
		}
		if destToken, ok := destSupportedTokens[token]; !ok || destToken.TokenPoolType == FeeTokenOnly {
			// If the token is not supported on the destination chain we
			// should not enable it for this ramp. If we enable the token,
			// txs could be sent but not executed, keeping the tokens in limbo.
			continue
		}

		tokensAndPools = append(tokensAndPools, evm_2_evm_onramp.InternalPoolUpdate{
			Token: tokenConfig.Token,
			Pool:  tokenConfig.Pool,
		})
		tokenTransferFeeConfig = append(tokenTransferFeeConfig, evm_2_evm_onramp.EVM2EVMOnRampTokenTransferFeeConfigArgs{
			Token: tokenConfig.Token,
			Ratio: 5_0, // 5 bps
		})
	}

	var feeTokenConfig []evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs

	for _, feeToken := range client.ChainConfig.FeeTokens {
		tokenConfig := client.ChainConfig.SupportedTokens[feeToken]
		gasMultiplier := uint64(GAS_MULTIPLIER)
		premiumMultiplier := uint64(DEFAULT_PREMIUM_MULTIPLIER)
		// Let link cost 10% of the non-link fee. This helps with our ping pong running out of funds.
		if feeToken == LINK {
			gasMultiplier = 1e17
			premiumMultiplier = LINK_PREMIUM_MULTIPLIER
		}

		feeTokenConfig = append(feeTokenConfig, evm_2_evm_onramp.EVM2EVMOnRampFeeTokenConfigArgs{
			Token:                  tokenConfig.Token,
			NetworkFeeUSD:          NETWORK_FEE,
			MinTokenTransferFeeUSD: MIN_TOKEN_TRANSFER_FEE,
			MaxTokenTransferFeeUSD: MAX_TOKEN_TRANSFER_FEE,
			GasMultiplier:          gasMultiplier,
			PremiumMultiplier:      premiumMultiplier,
			Enabled:                true,
		})
	}

	client.Logger.Infof("Deploying OnRamp: destinationChains %+v, tokensAndPools %+v", destChainSelector, tokensAndPools)
	onRampAddress, tx, _, err := evm_2_evm_onramp.DeployEVM2EVMOnRamp(
		client.Owner,  // user
		client.Client, // client
		evm_2_evm_onramp.EVM2EVMOnRampStaticConfig{
			LinkToken:         client.ChainConfig.SupportedTokens[LINK].Token,
			ChainSelector:     GetCCIPChainSelector(client.ChainConfig.EvmChainId),
			DestChainSelector: destChainSelector,
			DefaultTxGasLimit: DEFAULT_GAS_LIMIT,
			MaxNopFeesJuels:   big.NewInt(0).Mul(big.NewInt(MAX_NOP_FEES_LINK), big.NewInt(1e18)),
			PrevOnRamp:        prevOnRamp,
			ArmProxy:          client.ChainConfig.ARMProxy,
		},
		evm_2_evm_onramp.EVM2EVMOnRampDynamicConfig{
			Router:                client.ChainConfig.Router,
			MaxTokensLength:       MAX_TOKEN_LENGTH,
			DestGasOverhead:       DEST_GAS_OVERHEAD,
			DestGasPerPayloadByte: DEST_GAS_PER_PAYLOAD_BYTE,
			PriceRegistry:         client.ChainConfig.PriceRegistry,
			MaxDataSize:           MAX_DATA_SIZE,
			MaxGasLimit:           MAX_TX_GAS_LIMIT,
		},
		tokensAndPools,
		[]common.Address{}, // allow list
		evm_2_evm_onramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  UsdToRateLimitValue(RATE_LIMIT_CAPACITY_DOLLAR),
			Rate:      UsdToRateLimitValue(RATE_LIMIT_RATE_DOLLAR),
		},
		feeTokenConfig,
		tokenTransferFeeConfig,
		[]evm_2_evm_onramp.EVM2EVMOnRampNopAndWeight{},
	)
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)

	client.Logger.Infof(fmt.Sprintf("Onramp deployed on %s in tx %s", onRampAddress.String(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash())))
	laneConfig.OnRamp = onRampAddress
}

func deployOffRamp(t *testing.T, client EvmConfig, sourceChainSelector uint64, sourceTokens map[Token]EVMBridgedToken, laneConfig *EVMLaneConfig, onRamp common.Address, prevOffRamp common.Address) {
	if !laneConfig.DeploySettings.DeployLane {
		client.Logger.Infof("Skipping OffRamp deployment, using offRamp on %s", laneConfig.OnRamp)
		return
	}

	var syncedSourceTokens []common.Address
	var syncedDestPools []common.Address

	for tokenName, tokenConfig := range sourceTokens {
		if sourceToken, ok := client.ChainConfig.SupportedTokens[tokenName]; ok && sourceToken.TokenPoolType != FeeTokenOnly {
			syncedSourceTokens = append(syncedSourceTokens, tokenConfig.Token)
			syncedDestPools = append(syncedDestPools, client.ChainConfig.SupportedTokens[tokenName].Pool)
		} else {
			client.Logger.Warnf("Token %s not supported by destination chain", tokenName)
		}
	}

	client.Logger.Infof("Deploying OffRamp")
	offRampAddress, tx, _, err := evm_2_evm_offramp.DeployEVM2EVMOffRamp(
		client.Owner,
		client.Client,
		evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
			CommitStore:         laneConfig.CommitStore,
			ChainSelector:       GetCCIPChainSelector(client.ChainConfig.EvmChainId),
			SourceChainSelector: sourceChainSelector,
			OnRamp:              onRamp,
			PrevOffRamp:         prevOffRamp,
			ArmProxy:            client.ChainConfig.ARMProxy,
		},
		syncedSourceTokens,
		syncedDestPools,
		evm_2_evm_offramp.RateLimiterConfig{
			IsEnabled: true,
			Capacity:  UsdToRateLimitValue(RATE_LIMIT_CAPACITY_DOLLAR),
			Rate:      UsdToRateLimitValue(RATE_LIMIT_RATE_DOLLAR),
		},
	)
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)

	client.Logger.Infof("OffRamp contract deployed on %s in tx: %s", offRampAddress.Hex(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	laneConfig.OffRamp = offRampAddress

	client.Logger.Infof(fmt.Sprintf("Offramp configured for already deployed router in tx %s", helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash())))
}

func deployCommitStore(t *testing.T, client EvmConfig, lane *EVMLaneConfig, sourceChainSelector uint64, onRamp common.Address) {
	if !lane.DeploySettings.DeployLane {
		client.Logger.Infof("Skipping CommitStore deployment, using CommitStore on %s", lane.CommitStore)
		return
	}

	client.Logger.Infof("Deploying commitStore")
	commitStoreAddress, tx, _, err := commit_store.DeployCommitStore(
		client.Owner,  // user
		client.Client, // client
		commit_store.CommitStoreStaticConfig{
			ChainSelector:       GetCCIPChainSelector(client.ChainConfig.EvmChainId),
			SourceChainSelector: sourceChainSelector,
			OnRamp:              onRamp,
			ArmProxy:            client.ChainConfig.ARMProxy,
		},
	)
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
	client.Logger.Infof("CommitStore deployed on %s in tx: %s", commitStoreAddress.Hex(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))
	lane.CommitStore = commitStoreAddress
}

func deployPingPongDapp(t *testing.T, router common.Address, client EvmConfig, lane *EVMLaneConfig) {
	fundingAmount := big.NewInt(1e18)

	if !lane.DeploySettings.DeployPingPongDapp {
		return
	}

	feeToken := client.ChainConfig.SupportedTokens[LINK].Token
	client.Logger.Infof("Deploying source chain ping pong dapp")

	pingPongDappAddress, tx, _, err := ping_pong_demo.DeployPingPongDemo(
		client.Owner,
		client.Client,
		router,
		feeToken,
	)
	shared.RequireNoError(t, err)

	err = shared.WaitForMined(client.Logger, client.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
	client.Logger.Infof("Ping pong deployed on %s in tx: %s", pingPongDappAddress.Hex(), helpers.ExplorerLink(int64(client.ChainConfig.EvmChainId), tx.Hash()))

	lane.PingPongDapp = pingPongDappAddress
	err = FundPingPong(client, lane, fundingAmount, feeToken)
	shared.RequireNoError(t, err)
}

func finalizePingPongDaapDeployment(t *testing.T, sourceClient EvmConfig, sourceLane *EVMLaneConfig, destClient EvmConfig, destLane *EVMLaneConfig) {
	if !sourceLane.DeploySettings.DeployPingPongDapp && !destLane.DeploySettings.DeployPingPongDapp {
		sourceClient.Logger.Infof("Skipping ping pong deployment")
		return
	}

	pingDapp, err := ping_pong_demo.NewPingPongDemo(sourceLane.PingPongDapp, sourceClient.Client)
	shared.RequireNoError(t, err)

	tx, err := pingDapp.SetCounterpart(sourceClient.Owner, GetCCIPChainSelector(destClient.ChainConfig.EvmChainId), destLane.PingPongDapp)
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(sourceClient.Logger, sourceClient.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
	sourceClient.Logger.Infof("Ping pong dapp configured in tx: %s", helpers.ExplorerLink(int64(sourceClient.ChainConfig.EvmChainId), tx.Hash()))

	pongDapp, err := ping_pong_demo.NewPingPongDemo(destLane.PingPongDapp, destClient.Client)
	shared.RequireNoError(t, err)

	tx, err = pongDapp.SetCounterpart(destClient.Owner, GetCCIPChainSelector(sourceClient.ChainConfig.EvmChainId), sourceLane.PingPongDapp)
	shared.RequireNoError(t, err)
	err = shared.WaitForMined(destClient.Logger, destClient.Client, tx.Hash(), true)
	shared.RequireNoError(t, err)
	destClient.Logger.Infof("Ping pong dapp configured in tx: %s", helpers.ExplorerLink(int64(destClient.ChainConfig.EvmChainId), tx.Hash()))
}

func DeployPingPongDapps(t *testing.T, sourceClient *EvmDeploymentConfig, destClient *EvmDeploymentConfig) {
	deployPingPongDapp(t, sourceClient.ChainConfig.Router, sourceClient.OnlyEvmConfig(), &sourceClient.LaneConfig)
	deployPingPongDapp(t, destClient.ChainConfig.Router, destClient.OnlyEvmConfig(), &destClient.LaneConfig)
	finalizePingPongDaapDeployment(t, sourceClient.OnlyEvmConfig(), &sourceClient.LaneConfig, destClient.OnlyEvmConfig(), &destClient.LaneConfig)
}

func DeployUpgradePingPongDapps(t *testing.T, sourceClient *EvmDeploymentConfig, destClient *EvmDeploymentConfig) {
	deployPingPongDapp(t, sourceClient.ChainConfig.UpgradeRouter, sourceClient.OnlyEvmConfig(), &sourceClient.UpgradeLaneConfig)
	deployPingPongDapp(t, destClient.ChainConfig.UpgradeRouter, destClient.OnlyEvmConfig(), &destClient.UpgradeLaneConfig)
	finalizePingPongDaapDeployment(t, sourceClient.OnlyEvmConfig(), &sourceClient.UpgradeLaneConfig, destClient.OnlyEvmConfig(), &destClient.UpgradeLaneConfig)
}

func UsdToRateLimitValue(usd int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(1e18), big.NewInt(usd))
}
