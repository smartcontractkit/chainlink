package config

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
)

var (
	DefaultMinimumContractPayment        = assets.NewLinkFromJuels(100000000000000) // 0.0001 LINK
	DefaultGasLimit               uint64 = 500000
	DefaultGasPrice                      = assets.GWei(20)
	DefaultGasTip                        = assets.GWei(0)
)

type (
	// chainSpecificConfigDefaultSet lists the config defaults specific to a particular chain ID
	chainSpecificConfigDefaultSet struct {
		balanceMonitorEnabled                      bool
		balanceMonitorBlockDelay                   uint16
		blockEmissionIdleWarningThreshold          time.Duration
		blockHistoryEstimatorBatchSize             uint32
		blockHistoryEstimatorBlockDelay            uint16
		blockHistoryEstimatorBlockHistorySize      uint16
		blockHistoryEstimatorTransactionPercentile uint16
		eip1559DynamicFees                         bool
		ethTxReaperInterval                        time.Duration
		ethTxReaperThreshold                       time.Duration
		ethTxResendAfterThreshold                  time.Duration
		finalityDepth                              uint32
		flagsContractAddress                       string
		gasBumpPercent                             uint16
		gasBumpThreshold                           uint64
		gasBumpTxDepth                             uint16
		gasBumpWei                                 big.Int
		gasEstimatorMode                           string
		gasLimitDefault                            uint64
		gasLimitMultiplier                         float32
		gasLimitTransfer                           uint64
		gasPriceDefault                            big.Int
		gasTipCapDefault                           big.Int
		gasTipCapMinimum                           big.Int
		headTrackerHistoryDepth                    uint32
		headTrackerMaxBufferSize                   uint32
		headTrackerSamplingInterval                time.Duration
		linkContractAddress                        string
		logBackfillBatchSize                       uint32
		maxGasPriceWei                             big.Int
		maxInFlightTransactions                    uint32
		maxQueuedTransactions                      uint64
		minGasPriceWei                             big.Int
		minIncomingConfirmations                   uint32
		minRequiredOutgoingConfirmations           uint64
		minimumContractPayment                     *assets.Link
		nonceAutoSync                              bool
		ocrContractConfirmations                   uint16
		rpcDefaultBatchSize                        uint32
		// set true indicates its not the empty struct
		set bool
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
		blockHistoryEstimatorBlockHistorySize:      16,
		blockHistoryEstimatorTransactionPercentile: 60,
		eip1559DynamicFees:                         false,
		ethTxReaperInterval:                        1 * time.Hour,
		ethTxReaperThreshold:                       168 * time.Hour,
		ethTxResendAfterThreshold:                  1 * time.Minute,
		finalityDepth:                              50,
		gasBumpPercent:                             20,
		gasBumpThreshold:                           3,
		gasBumpTxDepth:                             10,
		gasBumpWei:                                 *assets.GWei(5),
		gasEstimatorMode:                           "BlockHistory",
		gasLimitDefault:                            DefaultGasLimit,
		gasLimitMultiplier:                         1.0,
		gasLimitTransfer:                           21000,
		gasPriceDefault:                            *DefaultGasPrice,
		gasTipCapDefault:                           *DefaultGasTip,
		gasTipCapMinimum:                           *big.NewInt(0),
		headTrackerHistoryDepth:                    100,
		headTrackerMaxBufferSize:                   3,
		headTrackerSamplingInterval:                1 * time.Second,
		linkContractAddress:                        "",
		logBackfillBatchSize:                       100,
		maxGasPriceWei:                             *assets.GWei(5000),
		maxInFlightTransactions:                    16,
		maxQueuedTransactions:                      250,
		minGasPriceWei:                             *assets.GWei(1),
		minIncomingConfirmations:                   3,
		minRequiredOutgoingConfirmations:           12,
		minimumContractPayment:                     DefaultMinimumContractPayment,
		nonceAutoSync:                              true,
		ocrContractConfirmations:                   4,
		rpcDefaultBatchSize:                        100,
		set:                                        true,
	}

	mainnet := fallbackDefaultSet
	mainnet.linkContractAddress = "0x514910771AF9Ca656af840dff83E8264EcF986CA"
	mainnet.minimumContractPayment = assets.NewLinkFromJuels(100000000000000000) // 0.1 LINK
	mainnet.blockHistoryEstimatorBlockHistorySize = 12                           // mainnet has longer block times than everything else, so ideally this is kept small to keep it responsive
	// NOTE: There are probably other variables we can tweak for Kovan and other
	// test chains, but the defaults have been working fine and if it ain't
	// broke, don't fix it.
	ropsten := mainnet
	ropsten.linkContractAddress = "0x20fe562d797a42dcb3399062ae9546cd06f63280"
	kovan := mainnet
	kovan.linkContractAddress = "0xa36085F69e2889c224210F603D836748e7dC0088"
	goerli := mainnet
	goerli.linkContractAddress = "0x326c977e6efc84e512bb9c30f76e30c160ed06fb"
	rinkeby := mainnet
	rinkeby.linkContractAddress = "0x01BE23585060835E02B77ef475b0Cc51aA1e0709"

	// xDai currently uses AuRa (like Parity) consensus so finality rules will be similar to parity
	// See: https://www.poa.network/for-users/whitepaper/poadao-v1/proof-of-authority
	// NOTE: xDai is planning to move to Honeybadger BFT which might have different finality guarantees
	// https://www.xdaichain.com/for-validators/consensus/honeybadger-bft-consensus
	// For worst case re-org depth on AuRa, assume 2n+2 (see: https://github.com/poanetwork/wiki/wiki/Aura-Consensus-Protocol-Audit)
	// With xDai's current maximum of 19 validators then 40 blocks is the maximum possible re-org)
	// The mainnet default of 50 blocks is ok here
	xDaiMainnet := fallbackDefaultSet
	xDaiMainnet.gasBumpThreshold = 3 // 15s delay since feeds update every minute in volatile situations
	xDaiMainnet.gasPriceDefault = *assets.GWei(1)
	xDaiMainnet.minGasPriceWei = *assets.GWei(1) // 1 Gwei is the minimum accepted by the validators (unless whitelisted)
	xDaiMainnet.maxGasPriceWei = *assets.GWei(500)
	xDaiMainnet.linkContractAddress = "0xE2e73A1c69ecF83F464EFCE6A5be353a37cA09b2"

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

	hecoMainnet := bscMainnet

	// Polygon has a 1s block time and looser finality guarantees than ereum.
	// Re-orgs have been observed at 64 blocks or even deeper
	polygonMainnet := fallbackDefaultSet
	polygonMainnet.balanceMonitorBlockDelay = 13 // equivalent of 1 eth block seems reasonable
	polygonMainnet.finalityDepth = 200           // A sprint is 64 blocks long and doesn't guarantee finality. To be safe we take three sprints (192 blocks) plus a safety margin
	polygonMainnet.gasBumpThreshold = 5          // 10s delay since feeds update every minute in volatile situations
	polygonMainnet.gasBumpWei = *assets.GWei(20)
	polygonMainnet.gasPriceDefault = *assets.GWei(1)
	polygonMainnet.headTrackerHistoryDepth = 250 // FinalityDepth + safety margin
	polygonMainnet.headTrackerSamplingInterval = 1 * time.Second
	polygonMainnet.blockEmissionIdleWarningThreshold = 15 * time.Second
	polygonMainnet.maxQueuedTransactions = 2000 // Since re-orgs on Polygon can be so large, we need a large safety buffer to allow time for the queue to clear down before we start dropping transactions
	polygonMainnet.minGasPriceWei = *assets.GWei(1)
	polygonMainnet.ethTxResendAfterThreshold = 5 * time.Minute // 5 minutes is roughly 300 blocks on Polygon. Since re-orgs occur often and can be deep we want to avoid overloading the node with a ton of re-sent unconfirmed transactions.
	polygonMainnet.blockHistoryEstimatorBlockDelay = 10        // Must be set to something large here because Polygon has so many re-orgs that otherwise we are constantly refetching
	polygonMainnet.blockHistoryEstimatorBlockHistorySize = 24
	polygonMainnet.linkContractAddress = "0xb0897686c545045afc77cf20ec7a532e3120e0f1"
	polygonMainnet.minIncomingConfirmations = 5
	polygonMainnet.minRequiredOutgoingConfirmations = 12
	polygonMumbai := polygonMainnet
	polygonMumbai.linkContractAddress = "0x326C977E6efc84E512bB9C30f76E30c160eD06FB"

	// Arbitrum is an L2 chain. Pending proper L2 support, for now we rely on their sequencer
	arbitrumMainnet := fallbackDefaultSet
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
	optimismMainnet.ethTxResendAfterThreshold = 15 * time.Second
	optimismMainnet.finalityDepth = 1    // Sequencer offers absolute finality as long as no re-org longer than 20 blocks occurs on main chain this event would require special handling (new txm)
	optimismMainnet.gasBumpThreshold = 0 // Never bump gas on optimism
	optimismMainnet.gasEstimatorMode = "Optimism"
	optimismMainnet.headTrackerHistoryDepth = 10
	optimismMainnet.headTrackerSamplingInterval = 1 * time.Second
	optimismMainnet.linkContractAddress = "0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6"
	optimismMainnet.minIncomingConfirmations = 1
	optimismMainnet.minRequiredOutgoingConfirmations = 0
	optimismMainnet.ocrContractConfirmations = 1
	optimismKovan := optimismMainnet
	optimismKovan.blockEmissionIdleWarningThreshold = 30 * time.Minute
	optimismKovan.linkContractAddress = "0x4911b761993b9c8c0d14Ba2d86902AF6B0074F5B"

	// Fantom
	fantomMainnet := fallbackDefaultSet
	fantomMainnet.gasPriceDefault = *assets.GWei(15)
	fantomMainnet.linkContractAddress = "0x6f43ff82cca38001b6699a8ac47a2d0e66939407"
	fantomMainnet.minIncomingConfirmations = 3
	fantomMainnet.minRequiredOutgoingConfirmations = 2
	fantomTestnet := fantomMainnet
	fantomTestnet.linkContractAddress = "0xfafedb041c0dd4fa2dc0d87a6b0979ee6fa7af5f"

	// RSK
	// RSK prices its txes in sats not wei
	rskMainnet := fallbackDefaultSet
	rskMainnet.gasPriceDefault = *big.NewInt(50000000) // It's about 100 times more expensive than Wei, very roughly speaking
	rskMainnet.linkContractAddress = "0x14adae34bef7ca957ce2dde5add97ea050123827"
	rskMainnet.maxGasPriceWei = *big.NewInt(50000000000)
	rskMainnet.minGasPriceWei = *big.NewInt(0)
	rskMainnet.minimumContractPayment = assets.NewLinkFromJuels(1000000000000000)
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

	avalancheFuji := avalancheMainnet
	avalancheFuji.linkContractAddress = "0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846"

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
}
