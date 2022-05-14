package config

import (
	"fmt"
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/config"
)

var (
	// DefaultGasFeeCap is the default value to use for Fee Cap in EIP-1559 transactions
	DefaultGasFeeCap                     = assets.GWei(100)
	DefaultGasLimit               uint64 = 500000
	DefaultGasPrice                      = assets.GWei(20)
	DefaultGasTip                        = big.NewInt(1)                           // go-ethereum requires the tip to be at least 1 wei
	DefaultMinimumContractPayment        = assets.NewLinkFromJuels(10000000000000) // 0.00001 LINK
)

type (
	// chainSpecificConfigDefaultSet lists the config defaults specific to a particular chain ID
	chainSpecificConfigDefaultSet struct {
		balanceMonitorEnabled                          bool
		balanceMonitorBlockDelay                       uint16
		blockEmissionIdleWarningThreshold              time.Duration
		blockHistoryEstimatorBatchSize                 uint32
		blockHistoryEstimatorBlockDelay                uint16
		blockHistoryEstimatorBlockHistorySize          uint16
		blockHistoryEstimatorEIP1559FeeCapBufferBlocks *uint16
		blockHistoryEstimatorTransactionPercentile     uint16
		chainType                                      config.ChainType
		eip1559DynamicFees                             bool
		ethTxReaperInterval                            time.Duration
		ethTxReaperThreshold                           time.Duration
		ethTxResendAfterThreshold                      time.Duration
		finalityDepth                                  uint32
		flagsContractAddress                           string
		gasBumpPercent                                 uint16
		gasBumpThreshold                               uint64
		gasBumpTxDepth                                 uint16
		gasBumpWei                                     big.Int
		gasEstimatorMode                               string
		gasFeeCapDefault                               big.Int
		gasLimitDefault                                uint64
		gasLimitMultiplier                             float32
		gasLimitTransfer                               uint64
		gasPriceDefault                                big.Int
		gasTipCapDefault                               big.Int
		gasTipCapMinimum                               big.Int
		headTrackerHistoryDepth                        uint32
		headTrackerMaxBufferSize                       uint32
		headTrackerSamplingInterval                    time.Duration
		linkContractAddress                            string
		logBackfillBatchSize                           uint32
		logPollInterval                                time.Duration
		maxGasPriceWei                                 big.Int
		maxInFlightTransactions                        uint32
		maxQueuedTransactions                          uint64
		minGasPriceWei                                 big.Int
		minIncomingConfirmations                       uint32
		minRequiredOutgoingConfirmations               uint64
		minimumContractPayment                         *assets.Link
		nodeDeadAfterNoNewHeadersThreshold             time.Duration
		nodePollFailureThreshold                       uint32
		nodePollInterval                               time.Duration

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
		balanceMonitorEnabled:                      true,
		balanceMonitorBlockDelay:                   1,
		blockEmissionIdleWarningThreshold:          1 * time.Minute,
		blockHistoryEstimatorBatchSize:             4, // FIXME: Workaround `websocket: read limit exceeded` until https://app.clubhouse.io/chainlinklabs/story/6717/geth-websockets-can-sometimes-go-bad-under-heavy-load-proposal-for-eth-node-balancer
		blockHistoryEstimatorBlockDelay:            1,
		blockHistoryEstimatorBlockHistorySize:      8,
		blockHistoryEstimatorTransactionPercentile: 60,
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
		gasLimitMultiplier:                    1.0,
		gasLimitTransfer:                      21000,
		gasPriceDefault:                       *DefaultGasPrice,
		gasTipCapDefault:                      *DefaultGasTip,
		gasTipCapMinimum:                      *big.NewInt(1),
		headTrackerHistoryDepth:               100,
		headTrackerMaxBufferSize:              3,
		headTrackerSamplingInterval:           1 * time.Second,
		linkContractAddress:                   "",
		logBackfillBatchSize:                  100,
		logPollInterval:                       15 * time.Second,
		maxGasPriceWei:                        *assets.GWei(100000),
		maxInFlightTransactions:               16,
		maxQueuedTransactions:                 250,
		minGasPriceWei:                        *assets.GWei(1),
		minIncomingConfirmations:              3,
		minRequiredOutgoingConfirmations:      12,
		minimumContractPayment:                DefaultMinimumContractPayment,
		nodeDeadAfterNoNewHeadersThreshold:    3 * time.Minute,
		nodePollFailureThreshold:              5,
		nodePollInterval:                      10 * time.Second,
		nonceAutoSync:                         true,
		useForwarders:                         false,
		ocrContractConfirmations:              4,
		ocrContractTransmitterTransmitTimeout: 10 * time.Second,
		ocrDatabaseTimeout:                    10 * time.Second,
		ocrObservationGracePeriod:             1 * time.Second,
		rpcDefaultBatchSize:                   100,
		complete:                              true,
	}

	mainnet := fallbackDefaultSet
	mainnet.blockHistoryEstimatorBlockHistorySize = 4 // EIP-1559 does well on a smaller block history size
	mainnet.blockHistoryEstimatorTransactionPercentile = 50
	mainnet.eip1559DynamicFees = true // enable EIP-1559 on Eth Mainnet and all testnets
	mainnet.linkContractAddress = "0x514910771AF9Ca656af840dff83E8264EcF986CA"
	mainnet.minimumContractPayment = assets.NewLinkFromJuels(100000000000000000) // 0.1 LINK
	// NOTE: There are probably other variables we can tweak for Kovan and other
	// test chains, but the defaults have been working fine and if it ain't
	// broke, don't fix it.
	ropsten := mainnet
	ropsten.linkContractAddress = "0x20fe562d797a42dcb3399062ae9546cd06f63280"
	kovan := mainnet
	kovan.linkContractAddress = "0xa36085F69e2889c224210F603D836748e7dC0088"
	kovan.eip1559DynamicFees = false // FIXME: Kovan has strange behaviour with EIP1559, see: https://app.shortcut.com/chainlinklabs/story/34098/kovan-can-emit-blocks-that-violate-assumptions-in-block-history-estimator
	goerli := mainnet
	goerli.linkContractAddress = "0x326c977e6efc84e512bb9c30f76e30c160ed06fb"
	goerli.eip1559DynamicFees = false // TODO: EIP1559 on goerli has not been adequately tested, see: https://app.shortcut.com/chainlinklabs/story/34098/kovan-can-emit-blocks-that-violate-assumptions-in-block-history-estimator
	rinkeby := mainnet
	rinkeby.linkContractAddress = "0x01BE23585060835E02B77ef475b0Cc51aA1e0709"
	rinkeby.eip1559DynamicFees = false // TODO: EIP1559 on rinkeby has not been adequately tested, see: https://app.shortcut.com/chainlinklabs/story/34098/kovan-can-emit-blocks-that-violate-assumptions-in-block-history-estimator

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
	bscMainnet.balanceMonitorBlockDelay = 2
	bscMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
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
	bscMainnet.minRequiredOutgoingConfirmations = 12
	bscMainnet.ocrDatabaseTimeout = 2 * time.Second
	bscMainnet.ocrContractTransmitterTransmitTimeout = 2 * time.Second
	bscMainnet.ocrObservationGracePeriod = 500 * time.Millisecond
	bscMainnet.logPollInterval = 3 * time.Second

	hecoMainnet := bscMainnet

	// Polygon has a 1s block time and looser finality guarantees than ethereum.
	// Re-orgs have been observed at 64 blocks or even deeper
	polygonMainnet := fallbackDefaultSet
	polygonMainnet.balanceMonitorBlockDelay = 13 // equivalent of 1 eth block seems reasonable
	polygonMainnet.finalityDepth = 500           // It is quite common to see re-orgs on polygon go several hundred blocks deep. See: https://polygonscan.com/blocks_forked
	polygonMainnet.gasBumpThreshold = 5          // 10s delay since feeds update every minute in volatile situations
	polygonMainnet.gasBumpWei = *assets.GWei(20)
	polygonMainnet.gasPriceDefault = *assets.GWei(1)
	polygonMainnet.headTrackerHistoryDepth = 2000 // Polygon suffers from a tremendous number of re-orgs, we need to set this to something very large to be conservative enough
	polygonMainnet.headTrackerSamplingInterval = 1 * time.Second
	polygonMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
	polygonMainnet.maxQueuedTransactions = 5000                // Since re-orgs on Polygon can be so large, we need a large safety buffer to allow time for the queue to clear down before we start dropping transactions
	polygonMainnet.maxGasPriceWei = *assets.UEther(200)        // 200,000 GWei
	polygonMainnet.gasPriceDefault = *assets.GWei(30)          // Many Polygon RPC providers set a minimum of 30 GWei on mainnet to prevent spam
	polygonMainnet.minGasPriceWei = *assets.GWei(30)           // Many Polygon RPC providers set a minimum of 30 GWei on mainnet to prevent spam
	polygonMainnet.ethTxResendAfterThreshold = 1 * time.Minute // Matic nodes under high mempool pressure are liable to drop txes, we need to ensure we keep sending them
	polygonMainnet.blockHistoryEstimatorBlockDelay = 10        // Must be set to something large here because Polygon has so many re-orgs that otherwise we are constantly refetching
	polygonMainnet.blockHistoryEstimatorBlockHistorySize = 24
	polygonMainnet.linkContractAddress = "0xb0897686c545045afc77cf20ec7a532e3120e0f1"
	polygonMainnet.minIncomingConfirmations = 5
	polygonMainnet.minRequiredOutgoingConfirmations = 12
	polygonMainnet.logPollInterval = 1 * time.Second
	polygonMumbai := polygonMainnet
	polygonMumbai.gasPriceDefault = *assets.GWei(1)
	polygonMumbai.minGasPriceWei = *assets.GWei(1)
	polygonMumbai.linkContractAddress = "0x326C977E6efc84E512bB9C30f76E30c160eD06FB"

	// Arbitrum is an L2 chain. Pending proper L2 support, for now we rely on their sequencer
	arbitrumMainnet := fallbackDefaultSet
	arbitrumMainnet.chainType = config.ChainArbitrum
	arbitrumMainnet.gasBumpThreshold = 0 // Disable gas bumping on arbitrum
	arbitrumMainnet.gasLimitDefault = 7000000
	arbitrumMainnet.gasLimitTransfer = 800000            // estimating gas returns 695,344 so 800,000 should be safe with some buffer
	arbitrumMainnet.gasPriceDefault = *assets.GWei(1000) // Arbitrum uses something like a Vickrey auction model where gas price represents a "max bid". In practice we usually pay much less
	arbitrumMainnet.maxGasPriceWei = *assets.GWei(1000)  // Fix the gas price
	arbitrumMainnet.minGasPriceWei = *assets.GWei(1000)  // Fix the gas price
	arbitrumMainnet.gasEstimatorMode = "FixedPrice"
	arbitrumMainnet.blockHistoryEstimatorBlockHistorySize = 0 // Force an error if someone set GAS_UPDATER_ENABLED=true by accident; we never want to run the block history estimator on arbitrum
	arbitrumMainnet.linkContractAddress = "0xf97f4df75117a78c1A5a0DBb814Af92458539FB4"
	arbitrumMainnet.ocrContractConfirmations = 1
	arbitrumRinkeby := arbitrumMainnet
	arbitrumRinkeby.linkContractAddress = "0x615fBe6372676474d9e6933d310469c9b68e9726"

	// Optimism is an L2 chain. Pending proper L2 support, for now we rely on their sequencer
	optimismMainnet := fallbackDefaultSet
	optimismMainnet.balanceMonitorBlockDelay = 0
	optimismMainnet.blockHistoryEstimatorBlockHistorySize = 0 // Force an error if someone set GAS_UPDATER_ENABLED=true by accident; we never want to run the block history estimator on optimism
	optimismMainnet.chainType = config.ChainOptimism
	optimismMainnet.ethTxResendAfterThreshold = 15 * time.Second
	optimismMainnet.finalityDepth = 1    // Sequencer offers absolute finality as long as no re-org longer than 20 blocks occurs on main chain this event would require special handling (new txm)
	optimismMainnet.gasBumpThreshold = 0 // Never bump gas on optimism
	optimismMainnet.gasEstimatorMode = "Optimism2"
	optimismMainnet.headTrackerHistoryDepth = 10
	optimismMainnet.headTrackerSamplingInterval = 1 * time.Second
	optimismMainnet.linkContractAddress = "0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6"
	optimismMainnet.minIncomingConfirmations = 1
	optimismMainnet.minGasPriceWei = *big.NewInt(0) // Optimism uses the Optimism2 estimator; we don't want to place any limits on the minimum gas price
	optimismMainnet.minRequiredOutgoingConfirmations = 0
	optimismMainnet.ocrContractConfirmations = 1
	optimismKovan := optimismMainnet
	optimismKovan.blockEmissionIdleWarningThreshold = 30 * time.Minute
	optimismKovan.linkContractAddress = "0x4911b761993b9c8c0d14Ba2d86902AF6B0074F5B"

	// Fantom
	fantomMainnet := fallbackDefaultSet
	fantomMainnet.gasPriceDefault = *assets.GWei(15)
	fantomMainnet.maxGasPriceWei = *assets.GWei(200000)
	fantomMainnet.linkContractAddress = "0x6f43ff82cca38001b6699a8ac47a2d0e66939407"
	fantomMainnet.minIncomingConfirmations = 3
	fantomMainnet.minRequiredOutgoingConfirmations = 2
	fantomMainnet.logPollInterval = 1 * time.Second
	fantomTestnet := fantomMainnet
	fantomTestnet.linkContractAddress = "0xfafedb041c0dd4fa2dc0d87a6b0979ee6fa7af5f"

	// RSK
	// RSK prices its txes in sats not wei
	rskMainnet := fallbackDefaultSet
	rskMainnet.gasPriceDefault = *big.NewInt(50000000) // It's about 100 times more expensive than Wei, very roughly speaking
	rskMainnet.linkContractAddress = "0x14adae34bef7ca957ce2dde5add97ea050123827"
	rskMainnet.maxGasPriceWei = *big.NewInt(50000000000)
	rskMainnet.gasFeeCapDefault = *big.NewInt(100000000) // rsk does not yet support EIP-1559 but this allows validation to pass
	rskMainnet.minGasPriceWei = *big.NewInt(0)
	rskMainnet.minimumContractPayment = assets.NewLinkFromJuels(1000000000000000)
	rskMainnet.logPollInterval = 30 * time.Second
	rskTestnet := rskMainnet
	rskTestnet.linkContractAddress = "0x8bbbd80981fe76d44854d8df305e8985c19f0e78"

	// Avalanche
	avalancheMainnet := fallbackDefaultSet
	avalancheMainnet.linkContractAddress = "0x5947BB275c521040051D82396192181b413227A3"
	avalancheMainnet.finalityDepth = 1
	avalancheMainnet.gasEstimatorMode = "BlockHistory"
	avalancheMainnet.gasPriceDefault = *assets.GWei(25)
	avalancheMainnet.maxGasPriceWei = *assets.GWei(1000)
	avalancheMainnet.minGasPriceWei = *assets.GWei(25)
	avalancheMainnet.blockHistoryEstimatorBlockHistorySize = 24 // Average block time of 2s
	avalancheMainnet.blockHistoryEstimatorBlockDelay = 2
	avalancheMainnet.minIncomingConfirmations = 1
	avalancheMainnet.minRequiredOutgoingConfirmations = 1
	avalancheMainnet.ocrContractConfirmations = 1
	avalancheMainnet.logPollInterval = 3 * time.Second

	avalancheFuji := avalancheMainnet
	avalancheFuji.linkContractAddress = "0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846"

	// Harmony
	harmonyMainnet := fallbackDefaultSet
	harmonyMainnet.linkContractAddress = "0x218532a12a389a4a92fC0C5Fb22901D1c19198aA"
	harmonyMainnet.gasPriceDefault = *assets.GWei(5)
	harmonyMainnet.minIncomingConfirmations = 1
	harmonyMainnet.minRequiredOutgoingConfirmations = 2
	harmonyMainnet.logPollInterval = 2 * time.Second
	harmonyTestnet := harmonyMainnet
	harmonyTestnet.linkContractAddress = "0x8b12Ac23BFe11cAb03a634C1F117D64a7f2cFD3e"

	// OKExChain
	okxMainnet := fallbackDefaultSet
	okxTestnet := okxMainnet

	chainSpecificConfigDefaultSets = make(map[int64]chainSpecificConfigDefaultSet)
	chainSpecificConfigDefaultSets[1] = mainnet
	chainSpecificConfigDefaultSets[3] = ropsten
	chainSpecificConfigDefaultSets[4] = rinkeby
	chainSpecificConfigDefaultSets[5] = goerli
	chainSpecificConfigDefaultSets[42] = kovan
	chainSpecificConfigDefaultSets[10] = optimismMainnet
	chainSpecificConfigDefaultSets[69] = optimismKovan
	chainSpecificConfigDefaultSets[42161] = arbitrumMainnet
	chainSpecificConfigDefaultSets[421611] = arbitrumRinkeby
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

	// sanity check
	for id, c := range chainSpecificConfigDefaultSets {
		if !c.complete {
			panic(fmt.Sprintf("chain %d configuration incomplete - "+
				"start from fallbackDefaultSet instead of zero value", id))
		}
	}
}
