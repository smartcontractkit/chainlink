package config

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var (
	// DefaultGasFeeCap is the default value to use for Fee Cap in EIP-1559 transactions
	DefaultGasFeeCap                     = assets.GWei(100)
	DefaultGasLimit               uint32 = 500000
	DefaultGasPrice                      = assets.GWei(20)
	DefaultGasTip                        = assets.NewWeiI(1)                           // go-ethereum requires the tip to be at least 1 wei
	DefaultMinimumContractPayment        = assets.NewLinkFromJuels(10_000_000_000_000) // 0.00001 LINK

	MaxLegalGasPrice = assets.NewWei(utils.MaxUint256)
)

type (
	// chainSpecificConfigDefaultSet lists the config defaults specific to a particular chain ID
	// https://app.shortcut.com/chainlinklabs/story/33622/remove-legacy-config
	chainSpecificConfigDefaultSet struct {
		balanceMonitorEnabled                         bool
		blockEmissionIdleWarningThreshold             time.Duration
		blockHistoryEstimatorBatchSize                uint32
		blockHistoryEstimatorBlockDelay               uint16
		blockHistoryEstimatorBlockHistorySize         uint16
		blockHistoryEstimatorCheckInclusionBlocks     uint16
		blockHistoryEstimatorCheckInclusionPercentile uint16
		blockHistoryEstimatorTransactionPercentile    uint16
		chainType                                     config.ChainType
		eip1559DynamicFees                            bool
		ethTxReaperInterval                           time.Duration
		ethTxReaperThreshold                          time.Duration
		ethTxResendAfterThreshold                     time.Duration
		finalityDepth                                 uint32
		flagsContractAddress                          string
		gasBumpPercent                                uint16
		gasBumpThreshold                              uint64
		gasBumpTxDepth                                uint16
		gasBumpWei                                    assets.Wei
		gasEstimatorMode                              string
		gasFeeCapDefault                              assets.Wei
		gasLimitDefault                               uint32
		gasLimitMax                                   uint32
		gasLimitMultiplier                            float32
		gasLimitTransfer                              uint32
		gasLimitOCRJobType                            *uint32
		gasLimitDRJobType                             *uint32
		gasLimitVRFJobType                            *uint32
		gasLimitFMJobType                             *uint32
		gasLimitKeeperJobType                         *uint32
		gasPriceDefault                               assets.Wei
		gasTipCapDefault                              assets.Wei
		gasTipCapMinimum                              assets.Wei
		headTrackerHistoryDepth                       uint32
		headTrackerMaxBufferSize                      uint32
		headTrackerSamplingInterval                   time.Duration
		linkContractAddress                           string
		operatorFactoryAddress                        string
		logBackfillBatchSize                          uint32
		logKeepBlocksDepth                            uint32
		logPollInterval                               time.Duration
		maxGasPriceWei                                assets.Wei
		maxInFlightTransactions                       uint32
		maxQueuedTransactions                         uint64
		minGasPriceWei                                assets.Wei
		minIncomingConfirmations                      uint32
		minimumContractPayment                        *assets.Link
		nodeDeadAfterNoNewHeadersThreshold            time.Duration
		nodePollFailureThreshold                      uint32
		nodePollInterval                              time.Duration
		nodeSelectionMode                             string
		nodeSyncThreshold                             uint32

		nonceAutoSync       bool
		useForwarders       bool
		rpcDefaultBatchSize uint32
		// set true if fully configured
		complete bool

		// Chain specific OCR1 config
		ocrContractConfirmations              uint16
		ocrContractTransmitterTransmitTimeout time.Duration
		ocrDatabaseTimeout                    time.Duration
		ocrObservationGracePeriod             time.Duration

		// Chain specific OCR2 config
		ocr2AutomationGasLimit uint32
	}
)

var chainSpecificConfigDefaultSets map[int64]chainSpecificConfigDefaultSet

// fallbackDefaultSet represents the "base layer" of config defaults
// It can be overridden on a per-chain basis and may be used if the chain is unknown
var fallbackDefaultSet chainSpecificConfigDefaultSet

func init() {
	setChainSpecificConfigDefaultSets()
}

func setChainSpecificConfigDefaultSets() {
	// --------------------------IMPORTANT---------------------------
	// All config sets should "inherit" from fallbackDefaultSet and overwrite
	// fields as necessary. Do not create a new chainSpecificConfigDefaultSet from
	// scratch for a particular chain, since it may accidentally contain zero
	// values.
	// Be sure to copy and --not modify-- fallbackDefaultSet!
	// TODO: We should probably move these to TOML or JSON files
	// See: https://app.clubhouse.io/chainlinklabs/story/11091/chain-chainSpecificConfigDefaultSets-should-move-to-toml-json-files

	fallbackDefaultSet = chainSpecificConfigDefaultSet{
		balanceMonitorEnabled:                 true,
		blockEmissionIdleWarningThreshold:     1 * time.Minute,
		blockHistoryEstimatorBatchSize:        4, // FIXME: Workaround `websocket: read limit exceeded` until https://app.clubhouse.io/chainlinklabs/story/6717/geth-websockets-can-sometimes-go-bad-under-heavy-load-proposal-for-eth-node-balancer
		blockHistoryEstimatorBlockDelay:       1,
		blockHistoryEstimatorBlockHistorySize: 8,
		// Connectivity checker is conservative by default: if a transaction
		// has been above 90% mark for 12 blocks, it is relatively safe to
		// assume a connectivity issue.
		blockHistoryEstimatorCheckInclusionBlocks:     12,
		blockHistoryEstimatorCheckInclusionPercentile: 90,
		blockHistoryEstimatorTransactionPercentile:    60,
		chainType:                             "",
		eip1559DynamicFees:                    false,
		ethTxReaperInterval:                   1 * time.Hour,
		ethTxReaperThreshold:                  168 * time.Hour,
		ethTxResendAfterThreshold:             1 * time.Minute,
		finalityDepth:                         50,
		gasBumpPercent:                        20,
		gasBumpThreshold:                      3,
		gasBumpTxDepth:                        10,
		gasBumpWei:                            *assets.GWei(5),
		gasEstimatorMode:                      "BlockHistory",
		gasFeeCapDefault:                      *DefaultGasFeeCap,
		gasLimitDefault:                       DefaultGasLimit,
		gasLimitMax:                           DefaultGasLimit, // equal since no effect other than Arbitrum
		gasLimitMultiplier:                    1.0,
		gasLimitTransfer:                      21000,
		gasPriceDefault:                       *DefaultGasPrice,
		gasTipCapDefault:                      *DefaultGasTip,
		gasTipCapMinimum:                      *assets.NewWeiI(1),
		headTrackerHistoryDepth:               100,
		headTrackerMaxBufferSize:              3,
		headTrackerSamplingInterval:           1 * time.Second,
		linkContractAddress:                   "",
		logBackfillBatchSize:                  100,
		logKeepBlocksDepth:                    100_000,
		logPollInterval:                       15 * time.Second,
		maxGasPriceWei:                        *MaxLegalGasPrice,
		maxInFlightTransactions:               16,
		maxQueuedTransactions:                 250,
		minGasPriceWei:                        *assets.GWei(1),
		minIncomingConfirmations:              3,
		minimumContractPayment:                DefaultMinimumContractPayment,
		nodeDeadAfterNoNewHeadersThreshold:    3 * time.Minute,
		nodePollFailureThreshold:              5,
		nodePollInterval:                      10 * time.Second,
		nodeSelectionMode:                     client.NodeSelectionMode_HighestHead,
		nodeSyncThreshold:                     5,
		nonceAutoSync:                         true,
		ocrContractConfirmations:              4,
		ocrContractTransmitterTransmitTimeout: 10 * time.Second,
		ocrDatabaseTimeout:                    10 * time.Second,
		ocrObservationGracePeriod:             1 * time.Second,
		ocr2AutomationGasLimit:                5_300_000, // 5.3M: 5M upkeep gas limit + 300K overhead
		operatorFactoryAddress:                "",
		rpcDefaultBatchSize:                   100,
		useForwarders:                         false,
		complete:                              true,
	}

	mainnet := fallbackDefaultSet
	mainnet.blockHistoryEstimatorBlockHistorySize = 4 // EIP-1559 does well on a smaller block history size
	mainnet.blockHistoryEstimatorTransactionPercentile = 50
	mainnet.eip1559DynamicFees = true // enable EIP-1559 on Eth Mainnet and all testnets
	mainnet.linkContractAddress = "0x514910771AF9Ca656af840dff83E8264EcF986CA"
	mainnet.minimumContractPayment = assets.NewLinkFromJuels(100000000000000000) // 0.1 LINK
	mainnet.operatorFactoryAddress = "0x3e64cd889482443324f91bfa9c84fe72a511f48a"

	// NOTE: There are probably other variables we can tweak for Kovan and other
	// test chains, but the defaults have been working fine and if it ain't
	// broke, don't fix it.
	ropsten := mainnet
	ropsten.linkContractAddress = "0x20fe562d797a42dcb3399062ae9546cd06f63280"
	ropsten.operatorFactoryAddress = ""
	kovan := mainnet
	kovan.linkContractAddress = "0xa36085F69e2889c224210F603D836748e7dC0088"
	kovan.operatorFactoryAddress = "0x8007e24251b1D2Fc518Eb843A701d9cD21fe0aA3"
	// WONTFIX: Kovan has strange behaviour with EIP1559, see: https://app.shortcut.com/chainlinklabs/story/34098/kovan-can-emit-blocks-that-violate-assumptions-in-block-history-estimator
	// This is a WONTFIX because support for Kovan will soon be dropped
	kovan.eip1559DynamicFees = false
	goerli := mainnet
	goerli.linkContractAddress = "0x326c977e6efc84e512bb9c30f76e30c160ed06fb"
	goerli.eip1559DynamicFees = true
	goerli.operatorFactoryAddress = ""
	rinkeby := mainnet
	rinkeby.linkContractAddress = "0x01BE23585060835E02B77ef475b0Cc51aA1e0709"
	// WONTFIX: Rinkeby has not been tested with EIP1559
	// This is a WONTFIX because support for Rinkeby will soon be dropped
	rinkeby.eip1559DynamicFees = false
	rinkeby.operatorFactoryAddress = ""
	sepolia := mainnet
	sepolia.linkContractAddress = "0xb227f007804c16546Bd054dfED2E7A1fD5437678"
	sepolia.operatorFactoryAddress = "" // doesn't exist yet
	sepolia.eip1559DynamicFees = true

	// simulated chain is actually a local client that "pretends" to be a blockchain
	// see: https://goethereumbook.org/en/client-simulated/
	// generally speaking, this is only used in tests
	simulated := fallbackDefaultSet
	simulated.blockEmissionIdleWarningThreshold = 0
	simulated.nodeDeadAfterNoNewHeadersThreshold = 0 // Assume simulated chain can never die
	simulated.ethTxResendAfterThreshold = 0
	simulated.gasFeeCapDefault = *assets.GWei(100000)
	simulated.maxGasPriceWei = *assets.GWei(100000) // must be the same as gasFeeCapDefault in FixedPrice mode with gas bumping disabled
	simulated.finalityDepth = 1                     // Simulated does not have re-orgs
	simulated.gasBumpThreshold = 0                  // Never bump gas
	simulated.gasEstimatorMode = "FixedPrice"
	simulated.headTrackerHistoryDepth = 10
	simulated.headTrackerSamplingInterval = 1 * time.Second
	simulated.minIncomingConfirmations = 1
	simulated.minGasPriceWei = *assets.NewWeiI(0)
	simulated.ocrContractConfirmations = 1
	simulated.headTrackerMaxBufferSize = 100
	simulated.headTrackerSamplingInterval = 0
	simulated.ethTxReaperThreshold = 0
	simulated.minimumContractPayment = assets.NewLinkFromJuels(100)

	// xDai currently uses AuRa (like Parity) consensus so finality rules will be similar to parity
	// See: https://www.poa.network/for-users/whitepaper/poadao-v1/proof-of-authority
	// NOTE: xDai is planning to move to Honeybadger BFT which might have different finality guarantees
	// https://www.xdaichain.com/for-validators/consensus/honeybadger-bft-consensus
	// For worst case re-org depth on AuRa, assume 2n+2 (see: https://github.com/poanetwork/wiki/wiki/Aura-Consensus-Protocol-Audit)
	// With xDai's current maximum of 19 validators then 40 blocks is the maximum possible re-org)
	// The mainnet default of 50 blocks is ok here
	xDaiMainnet := fallbackDefaultSet
	xDaiMainnet.chainType = config.ChainXDai
	xDaiMainnet.gasBumpThreshold = 3 // 15s delay since feeds update every minute in volatile situations
	xDaiMainnet.gasPriceDefault = *assets.GWei(1)
	xDaiMainnet.minGasPriceWei = *assets.GWei(1) // 1 Gwei is the minimum accepted by the validators (unless whitelisted)
	xDaiMainnet.maxGasPriceWei = *assets.GWei(500)
	xDaiMainnet.linkContractAddress = "0xE2e73A1c69ecF83F464EFCE6A5be353a37cA09b2"
	xDaiMainnet.logPollInterval = 5 * time.Second

	// BSC uses Clique consensus with ~3s block times
	// Clique offers finality within (N/2)+1 blocks where N is number of signers
	// There are 21 BSC validators so theoretically finality should occur after 21/2+1 = 11 blocks
	bscMainnet := fallbackDefaultSet
	bscMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
	bscMainnet.nodeDeadAfterNoNewHeadersThreshold = 30 * time.Second
	bscMainnet.blockHistoryEstimatorBlockDelay = 2
	bscMainnet.blockHistoryEstimatorBlockHistorySize = 24
	bscMainnet.ethTxResendAfterThreshold = 1 * time.Minute
	bscMainnet.finalityDepth = 50   // Keeping this >> 11 because it's not expensive and gives us a safety margin
	bscMainnet.gasBumpThreshold = 5 // 15s delay since feeds update every minute in volatile situations
	bscMainnet.gasBumpWei = *assets.GWei(5)
	bscMainnet.gasPriceDefault = *assets.GWei(5)
	bscMainnet.headTrackerHistoryDepth = 100
	bscMainnet.headTrackerSamplingInterval = 1 * time.Second
	bscMainnet.linkContractAddress = "0x404460c6a5ede2d891e8297795264fde62adbb75"
	bscMainnet.minGasPriceWei = *assets.GWei(1)
	bscMainnet.minIncomingConfirmations = 3
	bscMainnet.ocrDatabaseTimeout = 2 * time.Second
	bscMainnet.ocrContractTransmitterTransmitTimeout = 2 * time.Second
	bscMainnet.ocrObservationGracePeriod = 500 * time.Millisecond
	bscMainnet.logPollInterval = 3 * time.Second
	bscMainnet.nodeSyncThreshold = 10

	hecoMainnet := bscMainnet

	// Polygon has a 1s block time and looser finality guarantees than ethereum.
	// Re-orgs have been observed at 64 blocks or even deeper
	polygonMainnet := fallbackDefaultSet
	polygonMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
	polygonMainnet.nodeDeadAfterNoNewHeadersThreshold = 30 * time.Second
	polygonMainnet.finalityDepth = 500  // It is quite common to see re-orgs on polygon go several hundred blocks deep. See: https://polygonscan.com/blocks_forked
	polygonMainnet.gasBumpThreshold = 5 // 10s delay since feeds update every minute in volatile situations
	polygonMainnet.gasBumpWei = *assets.GWei(20)
	polygonMainnet.headTrackerHistoryDepth = 2000 // Polygon suffers from a tremendous number of re-orgs, we need to set this to something very large to be conservative enough
	polygonMainnet.headTrackerSamplingInterval = 1 * time.Second
	polygonMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
	polygonMainnet.maxQueuedTransactions = 5000                // Since re-orgs on Polygon can be so large, we need a large safety buffer to allow time for the queue to clear down before we start dropping transactions
	polygonMainnet.gasPriceDefault = *assets.GWei(30)          // Many Polygon RPC providers set a minimum of 30 GWei on mainnet to prevent spam
	polygonMainnet.minGasPriceWei = *assets.GWei(30)           // Many Polygon RPC providers set a minimum of 30 GWei on mainnet to prevent spam
	polygonMainnet.ethTxResendAfterThreshold = 1 * time.Minute // Matic nodes under high mempool pressure are liable to drop txes, we need to ensure we keep sending them
	polygonMainnet.blockHistoryEstimatorBlockDelay = 10        // Must be set to something large here because Polygon has so many re-orgs that otherwise we are constantly refetching
	polygonMainnet.blockHistoryEstimatorBlockHistorySize = 24
	polygonMainnet.linkContractAddress = "0xb0897686c545045afc77cf20ec7a532e3120e0f1"
	polygonMainnet.minIncomingConfirmations = 5
	polygonMainnet.logPollInterval = 1 * time.Second
	polygonMainnet.nodeSyncThreshold = 10
	polygonMumbai := polygonMainnet
	polygonMumbai.gasPriceDefault = *assets.GWei(1)
	polygonMumbai.minGasPriceWei = *assets.GWei(1)
	polygonMumbai.linkContractAddress = "0x326C977E6efc84E512bB9C30f76E30c160eD06FB"

	// Arbitrum is an L2 chain. Pending proper L2 support, for now we rely on their sequencer
	arbitrumMainnet := fallbackDefaultSet
	arbitrumMainnet.blockEmissionIdleWarningThreshold = 0
	arbitrumMainnet.nodeDeadAfterNoNewHeadersThreshold = 0 // Arbitrum only emits blocks when a new tx is received, so this method of liveness detection is not useful
	arbitrumMainnet.chainType = config.ChainArbitrum
	arbitrumMainnet.gasBumpThreshold = 0 // Disable gas bumping on arbitrum
	arbitrumMainnet.gasEstimatorMode = "Arbitrum"
	arbitrumMainnet.gasLimitMax = 1_000_000_000
	arbitrumMainnet.minGasPriceWei = *assets.NewWeiI(0)          // Arbitrum uses the suggested gas price so we don't want to place any limits on the minimum
	arbitrumMainnet.gasPriceDefault = *assets.NewWeiI(100000000) // 0.1 gwei
	arbitrumMainnet.gasFeeCapDefault = *assets.GWei(1000)
	arbitrumMainnet.blockHistoryEstimatorBlockHistorySize = 0 // Force an error if someone set GAS_UPDATER_ENABLED=true by accident; we never want to run the block history estimator on arbitrum
	arbitrumMainnet.linkContractAddress = "0xf97f4df75117a78c1A5a0DBb814Af92458539FB4"
	arbitrumMainnet.logPollInterval = 1 * time.Second
	arbitrumMainnet.nodeSyncThreshold = 10
	arbitrumMainnet.ocrContractConfirmations = 1
	arbitrumRinkeby := arbitrumMainnet
	arbitrumRinkeby.linkContractAddress = "0x615fBe6372676474d9e6933d310469c9b68e9726"
	arbitrumGoerli := arbitrumRinkeby
	arbitrumGoerli.linkContractAddress = "0xd14838A68E8AFBAdE5efb411d5871ea0011AFd28"

	// Optimism is an L2 chain. Pending proper L2 support, for now we rely on their sequencer
	optimismMainnet := fallbackDefaultSet
	optimismMainnet.blockEmissionIdleWarningThreshold = 0
	optimismMainnet.nodeDeadAfterNoNewHeadersThreshold = 0    // Optimism only emits blocks when a new tx is received, so this method of liveness detection is not useful
	optimismMainnet.blockHistoryEstimatorBlockHistorySize = 0 // Force an error if someone set GAS_UPDATER_ENABLED=true by accident; we never want to run the block history estimator on optimism
	optimismMainnet.chainType = config.ChainOptimism
	optimismMainnet.ethTxResendAfterThreshold = 15 * time.Second
	optimismMainnet.finalityDepth = 1    // Sequencer offers absolute finality as long as no re-org longer than 20 blocks occurs on main chain this event would require special handling (new txm)
	optimismMainnet.gasBumpThreshold = 0 // Never bump gas on optimism
	optimismMainnet.gasEstimatorMode = "L2Suggested"
	optimismMainnet.headTrackerHistoryDepth = 10
	optimismMainnet.headTrackerSamplingInterval = 1 * time.Second
	optimismMainnet.linkContractAddress = "0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6"
	optimismMainnet.minIncomingConfirmations = 1
	optimismMainnet.minGasPriceWei = *assets.NewWeiI(0) // Optimism uses the L2Suggested estimator; we don't want to place any limits on the minimum gas price
	optimismMainnet.nodeSyncThreshold = 10
	optimismMainnet.ocrContractConfirmations = 1
	optimismMainnet.ocr2AutomationGasLimit = 6_500_000 // 5M (upkeep limit) + 1.5M. Optimism requires a larger overhead than normal chains
	optimismKovan := optimismMainnet
	optimismKovan.blockEmissionIdleWarningThreshold = 30 * time.Minute
	optimismKovan.linkContractAddress = "0x4911b761993b9c8c0d14Ba2d86902AF6B0074F5B"

	// Optimism's Bedrock upgrade uses EIP-1559.
	optimismBedrock := fallbackDefaultSet
	optimismBedrock.eip1559DynamicFees = true
	optimismBedrock.blockEmissionIdleWarningThreshold = 30 * time.Second
	optimismBedrock.nodeDeadAfterNoNewHeadersThreshold = 60 * time.Second // Bedrock produces blocks every two seconds, regardless of whether there are transactions to put in them or not.
	optimismBedrock.blockHistoryEstimatorBlockHistorySize = 24
	optimismBedrock.chainType = config.ChainOptimismBedrock
	optimismBedrock.ethTxResendAfterThreshold = 30 * time.Second
	optimismBedrock.logPollInterval = 2 * time.Second
	optimismBedrock.ocr2AutomationGasLimit = 6_500_000 // 5M (upkeep limit) + 1.5M. Optimism requires a larger overhead than normal chains
	// Bedrock supports a 10 block reorg resistance in L1. However, it considers a block final when it is included in a final block in L1.
	// L1 finality with PoS: Every 32 slots(epoch) each validator on the network has the opportunity to vote in favor of the epoch.
	// It takes two justified epochs for those epochs, and all the blocks inside of them, to be considered finalized.
	// (The proper way to consider finalization would be to mark an L2 block final when it gets included in a final L1 block, which requires special handling (new txm))
	optimismBedrock.finalityDepth = 200
	optimismBedrock.headTrackerHistoryDepth = 300
	optimismBedrock.nodeSyncThreshold = 10
	optimismGoerli := optimismBedrock
	optimismGoerli.minGasPriceWei = *assets.NewWeiI(1) // Gas prices were significantly reduced after the upgrade.
	optimismGoerli.linkContractAddress = "0xdc2CC710e42857672E7907CF474a69B63B93089f"

	// Fantom
	fantomMainnet := fallbackDefaultSet
	fantomMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
	fantomMainnet.blockHistoryEstimatorBlockDelay = 2
	fantomMainnet.gasPriceDefault = *assets.GWei(15)
	fantomMainnet.linkContractAddress = "0x6f43ff82cca38001b6699a8ac47a2d0e66939407"
	fantomMainnet.logPollInterval = 1 * time.Second
	fantomMainnet.minIncomingConfirmations = 3
	fantomMainnet.nodeDeadAfterNoNewHeadersThreshold = 30 * time.Second
	fantomMainnet.ocr2AutomationGasLimit = 3_800_000 // 3.5M (upkeep limit) + 300K. Fantom has a lower max gas limit than other chains
	fantomTestnet := fantomMainnet
	fantomTestnet.linkContractAddress = "0xfafedb041c0dd4fa2dc0d87a6b0979ee6fa7af5f"
	fantomTestnet.blockEmissionIdleWarningThreshold = 0
	fantomTestnet.nodeDeadAfterNoNewHeadersThreshold = 0 // Fantom testnet only emits blocks when a new tx is received, so this method of liveness detection is not useful

	// RSK
	// RSK prices its txes in sats not wei
	rskMainnet := fallbackDefaultSet
	rskMainnet.gasPriceDefault = *assets.NewWeiI(50000000) // It's about 100 times more expensive than Wei, very roughly speaking
	rskMainnet.linkContractAddress = "0x14adae34bef7ca957ce2dde5add97ea050123827"
	rskMainnet.maxGasPriceWei = *assets.NewWeiI(50000000000)
	rskMainnet.gasFeeCapDefault = *assets.NewWeiI(100000000) // rsk does not yet support EIP-1559 but this allows validation to pass
	rskMainnet.minGasPriceWei = *assets.NewWeiI(0)
	rskMainnet.minimumContractPayment = assets.NewLinkFromJuels(1000000000000000)
	rskMainnet.logPollInterval = 30 * time.Second
	rskTestnet := rskMainnet
	rskTestnet.linkContractAddress = "0x8bbbd80981fe76d44854d8df305e8985c19f0e78"

	// Avalanche
	avalancheMainnet := fallbackDefaultSet
	avalancheMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
	avalancheMainnet.nodeDeadAfterNoNewHeadersThreshold = 30 * time.Second
	avalancheMainnet.linkContractAddress = "0x5947BB275c521040051D82396192181b413227A3"
	avalancheMainnet.finalityDepth = 1
	avalancheMainnet.gasEstimatorMode = "BlockHistory"
	avalancheMainnet.gasPriceDefault = *assets.GWei(25)
	avalancheMainnet.minGasPriceWei = *assets.GWei(25)
	avalancheMainnet.blockHistoryEstimatorBlockHistorySize = 24 // Average block time of 2s
	avalancheMainnet.blockHistoryEstimatorBlockDelay = 2
	avalancheMainnet.minIncomingConfirmations = 1
	avalancheMainnet.ocrContractConfirmations = 1
	avalancheMainnet.logPollInterval = 3 * time.Second

	avalancheFuji := avalancheMainnet
	avalancheFuji.linkContractAddress = "0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846"

	// Harmony
	harmonyMainnet := fallbackDefaultSet
	harmonyMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
	harmonyMainnet.nodeDeadAfterNoNewHeadersThreshold = 30 * time.Second
	harmonyMainnet.linkContractAddress = "0x218532a12a389a4a92fC0C5Fb22901D1c19198aA"
	harmonyMainnet.gasPriceDefault = *assets.GWei(5)
	harmonyMainnet.minIncomingConfirmations = 1
	harmonyMainnet.logPollInterval = 2 * time.Second
	harmonyTestnet := harmonyMainnet
	harmonyTestnet.linkContractAddress = "0x8b12Ac23BFe11cAb03a634C1F117D64a7f2cFD3e"

	// OKExChain
	okxMainnet := fallbackDefaultSet
	okxTestnet := okxMainnet

	// Metis is an L2 chain based on Optimism.
	metisMainnet := fallbackDefaultSet
	metisMainnet.blockEmissionIdleWarningThreshold = 0
	metisMainnet.nodeDeadAfterNoNewHeadersThreshold = 0
	metisMainnet.blockHistoryEstimatorBlockHistorySize = 0 // Force an error if someone set GAS_UPDATER_ENABLED=true by accident; we never want to run the block history estimator on metis
	metisMainnet.chainType = config.ChainMetis
	metisMainnet.finalityDepth = 1    // Sequencer offers absolute finality
	metisMainnet.gasBumpThreshold = 0 // Never bump gas on metis
	metisMainnet.gasEstimatorMode = "L2Suggested"
	metisMainnet.linkContractAddress = ""
	metisMainnet.minIncomingConfirmations = 1
	metisMainnet.minGasPriceWei = *assets.NewWeiI(0) // Metis uses the L2Suggested estimator; we don't want to place any limits on the minimum gas price
	metisMainnet.ocrContractConfirmations = 1
	metisMainnet.nodeSyncThreshold = 10
	metisRinkeby := metisMainnet
	metisRinkeby.linkContractAddress = ""

	// Klaytn implements a special dynamic gas price model. It only charges the base fee and refunds the remaining.
	// Max gas price is 750ston(gwei), although it can change by Governance.
	// According to this: https://medium.com/klaytn/dynamic-gas-fee-pricing-mechanism-1dac83d2689 there are two ways to set proper gas fees:
	// Use the return value from eth_gasPrice method, or send transaction with max gas price and get refunded. First one is more future proof.
	klaytnMainnet := fallbackDefaultSet
	klaytnMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second // Klaytn has 1s block time
	klaytnMainnet.nodeDeadAfterNoNewHeadersThreshold = 30 * time.Second
	klaytnMainnet.gasPriceDefault = *assets.GWei(750)
	klaytnMainnet.finalityDepth = 1    // Klaytn offers instant finality
	klaytnMainnet.gasBumpThreshold = 0 // Never bump gas
	klaytnMainnet.gasEstimatorMode = "L2Suggested"
	klaytnMainnet.minIncomingConfirmations = 1
	klaytnMainnet.ocrContractConfirmations = 1
	klaytnTestnet := klaytnMainnet

	chainSpecificConfigDefaultSets = make(map[int64]chainSpecificConfigDefaultSet)
	chainSpecificConfigDefaultSets[1] = mainnet
	chainSpecificConfigDefaultSets[3] = ropsten
	chainSpecificConfigDefaultSets[4] = rinkeby
	chainSpecificConfigDefaultSets[5] = goerli
	chainSpecificConfigDefaultSets[11155111] = sepolia
	chainSpecificConfigDefaultSets[42] = kovan
	chainSpecificConfigDefaultSets[10] = optimismMainnet
	chainSpecificConfigDefaultSets[69] = optimismKovan
	chainSpecificConfigDefaultSets[420] = optimismGoerli
	chainSpecificConfigDefaultSets[42161] = arbitrumMainnet
	chainSpecificConfigDefaultSets[421611] = arbitrumRinkeby
	chainSpecificConfigDefaultSets[421613] = arbitrumGoerli
	chainSpecificConfigDefaultSets[56] = bscMainnet
	chainSpecificConfigDefaultSets[128] = hecoMainnet
	chainSpecificConfigDefaultSets[250] = fantomMainnet
	chainSpecificConfigDefaultSets[4002] = fantomTestnet
	chainSpecificConfigDefaultSets[137] = polygonMainnet
	chainSpecificConfigDefaultSets[80001] = polygonMumbai
	chainSpecificConfigDefaultSets[100] = xDaiMainnet
	chainSpecificConfigDefaultSets[30] = rskMainnet
	chainSpecificConfigDefaultSets[31] = rskTestnet
	chainSpecificConfigDefaultSets[43113] = avalancheFuji
	chainSpecificConfigDefaultSets[43114] = avalancheMainnet
	chainSpecificConfigDefaultSets[1666600000] = harmonyMainnet
	chainSpecificConfigDefaultSets[1666700000] = harmonyTestnet
	chainSpecificConfigDefaultSets[65] = okxTestnet
	chainSpecificConfigDefaultSets[66] = okxMainnet
	chainSpecificConfigDefaultSets[588] = metisRinkeby
	chainSpecificConfigDefaultSets[1088] = metisMainnet
	chainSpecificConfigDefaultSets[8217] = klaytnMainnet
	chainSpecificConfigDefaultSets[1001] = klaytnTestnet

	chainSpecificConfigDefaultSets[1337] = simulated

	// sanity check
	for id, c := range chainSpecificConfigDefaultSets {
		if !c.complete {
			panic(fmt.Sprintf("chain %d configuration incomplete - "+
				"start from fallbackDefaultSet instead of zero value", id))
		}
	}
}
