package chains

import (
	"math/big"
	"time"

	"github.com/smartcontractkit/chainlink/core/assets"
)

type (

	// ChainSpecificConfig lists the config defaults specific to a particular chain ID
	ChainSpecificConfig struct {
		BalanceMonitorEnabled                      bool
		BalanceMonitorBlockDelay                   uint16
		BlockEmissionIdleWarningThreshold          time.Duration
		BlockHistoryEstimatorBatchSize             uint32
		BlockHistoryEstimatorBlockDelay            uint16
		BlockHistoryEstimatorBlockHistorySize      uint16
		BlockHistoryEstimatorTransactionPercentile uint16
		EthTxReaperInterval                        time.Duration
		EthTxReaperThreshold                       time.Duration
		EthTxResendAfterThreshold                  time.Duration
		FinalityDepth                              uint
		FlagsContractAddress                       string
		GasBumpPercent                             uint16
		GasBumpThreshold                           uint64
		GasBumpTxDepth                             uint16
		GasBumpWei                                 big.Int
		GasEstimatorMode                           string
		GasLimitDefault                            uint64
		GasLimitMultiplier                         float32
		GasLimitTransfer                           uint64
		GasPriceDefault                            big.Int
		HeadTrackerHistoryDepth                    uint
		HeadTrackerMaxBufferSize                   uint
		HeadTrackerSamplingInterval                time.Duration
		Layer2Type                                 string
		LinkContractAddress                        string
		LogBackfillBatchSize                       uint32
		MaxGasPriceWei                             big.Int
		MaxInFlightTransactions                    uint32
		MaxQueuedTransactions                      uint64
		MinGasPriceWei                             big.Int
		MinIncomingConfirmations                   uint32
		MinRequiredOutgoingConfirmations           uint64
		MinimumContractPayment                     *assets.Link
		NonceAutoSync                              bool
		OCRContractConfirmations                   uint16
		RPCDefaultBatchSize                        uint32
		set                                        bool
	}
)

// FallbackConfig represents the "base layer" of config defaults
// It can be overridden on a per-chain basis and may be used if the chain is unknown
var FallbackConfig ChainSpecificConfig

func setConfigs() {
	// --------------------------IMPORTANT---------------------------
	// All config sets should "inherit" from FallbackConfig and overwrite
	// fields as necessary. Do not create a new ChainSpecificConfig from
	// scratch for a particular chain, since it may accidentally contain zero
	// values.
	// Be sure to copy and --not modify-- FallbackConfig!
	// TODO: Warn if any of these are overridden by user-specified config vars
	// See: https://app.clubhouse.io/chainlinklabs/story/11090/warn-if-nop-has-overridden-any-default-config-var
	// TODO: We should probably move these to TOML or JSON files
	// See: https://app.clubhouse.io/chainlinklabs/story/11091/chain-configs-should-move-to-toml-json-files

	FallbackConfig = ChainSpecificConfig{
		BalanceMonitorEnabled:                      true,
		BalanceMonitorBlockDelay:                   1,
		BlockEmissionIdleWarningThreshold:          1 * time.Minute,
		BlockHistoryEstimatorBatchSize:             4, // FIXME: Workaround `websocket: read limit exceeded` until https://app.clubhouse.io/chainlinklabs/story/6717/geth-websockets-can-sometimes-go-bad-under-heavy-load-proposal-for-eth-node-balancer
		BlockHistoryEstimatorBlockDelay:            1,
		BlockHistoryEstimatorBlockHistorySize:      24,
		BlockHistoryEstimatorTransactionPercentile: 60,
		EthTxReaperInterval:                        1 * time.Hour,
		EthTxReaperThreshold:                       168 * time.Hour,
		EthTxResendAfterThreshold:                  1 * time.Minute,
		FinalityDepth:                              50,
		GasBumpPercent:                             20,
		GasBumpThreshold:                           3,
		GasBumpTxDepth:                             10,
		GasBumpWei:                                 *assets.GWei(5),
		GasEstimatorMode:                           "BlockHistory",
		GasLimitDefault:                            500000,
		GasLimitMultiplier:                         1.0,
		GasLimitTransfer:                           21000,
		GasPriceDefault:                            *assets.GWei(20),
		HeadTrackerHistoryDepth:                    100,
		HeadTrackerMaxBufferSize:                   3,
		HeadTrackerSamplingInterval:                1 * time.Second,
		LinkContractAddress:                        "",
		LogBackfillBatchSize:                       100,
		MaxGasPriceWei:                             *assets.GWei(5000),
		MaxInFlightTransactions:                    16,
		MaxQueuedTransactions:                      250,
		MinGasPriceWei:                             *assets.GWei(1),
		MinIncomingConfirmations:                   3,
		MinRequiredOutgoingConfirmations:           12,
		MinimumContractPayment:                     assets.NewLink(100000000000000), // 0.0001 LINK
		NonceAutoSync:                              true,
		OCRContractConfirmations:                   4,
		RPCDefaultBatchSize:                        100,
		set:                                        true,
	}

	mainnet := FallbackConfig
	mainnet.LinkContractAddress = "0x514910771AF9Ca656af840dff83E8264EcF986CA"
	mainnet.MinimumContractPayment = assets.NewLink(1000000000000000000) // 1 LINK
	// NOTE: There are probably other variables we can tweak for Kovan and other
	// test chains, but the defaults have been working fine and if it ain't
	// broke, don't fix it.
	kovan := mainnet
	kovan.LinkContractAddress = "0xa36085F69e2889c224210F603D836748e7dC0088"
	goerli := mainnet
	goerli.LinkContractAddress = "0x326c977e6efc84e512bb9c30f76e30c160ed06fb"
	rinkeby := mainnet
	rinkeby.LinkContractAddress = "0x01BE23585060835E02B77ef475b0Cc51aA1e0709"

	// xDai currently uses AuRa (like Parity) consensus so finality rules will be similar to parity
	// See: https://www.poa.network/for-users/whitepaper/poadao-v1/proof-of-authority
	// NOTE: xDai is planning to move to Honeybadger BFT which might have different finality guarantees
	// https://www.xdaichain.com/for-validators/consensus/honeybadger-bft-consensus
	// For worst case re-org depth on AuRa, assume 2n+2 (see: https://github.com/poanetwork/wiki/wiki/Aura-Consensus-Protocol-Audit)
	// With xDai's current maximum of 19 validators then 40 blocks is the maximum possible re-org)
	// The mainnet default of 50 blocks is ok here
	xDaiMainnet := FallbackConfig
	xDaiMainnet.GasBumpThreshold = 3 // 15s delay since feeds update every minute in volatile situations
	xDaiMainnet.GasPriceDefault = *assets.GWei(1)
	xDaiMainnet.MinGasPriceWei = *assets.GWei(1) // 1 Gwei is the minimum accepted by the validators (unless whitelisted)
	xDaiMainnet.MaxGasPriceWei = *assets.GWei(500)
	xDaiMainnet.LinkContractAddress = "0xE2e73A1c69ecF83F464EFCE6A5be353a37cA09b2"

	// BSC uses Clique consensus with ~3s block times
	// Clique offers finality within (N/2)+1 blocks where N is number of signers
	// There are 21 BSC validators so theoretically finality should occur after 21/2+1 = 11 blocks
	bscMainnet := FallbackConfig
	bscMainnet.BalanceMonitorBlockDelay = 2
	bscMainnet.FinalityDepth = 50   // Keeping this >> 11 because it's not expensive and gives us a safety margin
	bscMainnet.GasBumpThreshold = 5 // 15s delay since feeds update every minute in volatile situations
	bscMainnet.GasBumpWei = *assets.GWei(5)
	bscMainnet.GasPriceDefault = *assets.GWei(5)
	bscMainnet.HeadTrackerHistoryDepth = 100
	bscMainnet.HeadTrackerSamplingInterval = 1 * time.Second
	bscMainnet.BlockEmissionIdleWarningThreshold = 15 * time.Second
	bscMainnet.MinGasPriceWei = *assets.GWei(1)
	bscMainnet.EthTxResendAfterThreshold = 1 * time.Minute
	bscMainnet.BlockHistoryEstimatorBlockDelay = 2
	bscMainnet.BlockHistoryEstimatorBlockHistorySize = 24
	bscMainnet.LinkContractAddress = "0x404460c6a5ede2d891e8297795264fde62adbb75"
	bscMainnet.MinIncomingConfirmations = 3
	bscMainnet.MinRequiredOutgoingConfirmations = 12

	hecoMainnet := bscMainnet

	// Polygon has a 1s block time and looser finality guarantees than ereum.
	// Re-orgs have been observed at 64 blocks or even deeper
	polygonMainnet := FallbackConfig
	polygonMainnet.BalanceMonitorBlockDelay = 13 // equivalent of 1 eth block seems reasonable
	polygonMainnet.FinalityDepth = 200           // A sprint is 64 blocks long and doesn't guarantee finality. To be safe we take three sprints (192 blocks) plus a safety margin
	polygonMainnet.GasBumpThreshold = 5          // 10s delay since feeds update every minute in volatile situations
	polygonMainnet.GasBumpWei = *assets.GWei(20)
	polygonMainnet.GasPriceDefault = *assets.GWei(1)
	polygonMainnet.HeadTrackerHistoryDepth = 250 // FinalityDepth + safety margin
	polygonMainnet.HeadTrackerSamplingInterval = 1 * time.Second
	polygonMainnet.BlockEmissionIdleWarningThreshold = 15 * time.Second
	polygonMainnet.MaxQueuedTransactions = 2000 // Since re-orgs on Polygon can be so large, we need a large safety buffer to allow time for the queue to clear down before we start dropping transactions
	polygonMainnet.MinGasPriceWei = *assets.GWei(1)
	polygonMainnet.EthTxResendAfterThreshold = 5 * time.Minute // 5 minutes is roughly 300 blocks on Polygon. Since re-orgs occur often and can be deep we want to avoid overloading the node with a ton of re-sent unconfirmed transactions.
	polygonMainnet.BlockHistoryEstimatorBlockDelay = 10
	polygonMainnet.BlockHistoryEstimatorBlockHistorySize = 24
	polygonMainnet.LinkContractAddress = "0xb0897686c545045afc77cf20ec7a532e3120e0f1"
	polygonMainnet.MinIncomingConfirmations = 5
	polygonMainnet.MinRequiredOutgoingConfirmations = 12
	polygonMumbai := polygonMainnet
	polygonMumbai.LinkContractAddress = "0x326C977E6efc84E512bB9C30f76E30c160eD06FB"

	// Arbitrum is an L2 chain. Pending proper L2 support, for now we rely on their sequencer
	arbitrumMainnet := FallbackConfig
	arbitrumMainnet.GasBumpThreshold = 0 // Disable gas bumping on arbitrum
	arbitrumMainnet.GasLimitDefault = 7000000
	arbitrumMainnet.GasLimitTransfer = 800000            // estimating gas returns 695,344 so 800,000 should be safe with some buffer
	arbitrumMainnet.GasPriceDefault = *assets.GWei(1000) // Arbitrum uses something like a Vickrey auction model where gas price represents a "max bid". In practice we usually pay much less
	arbitrumMainnet.MaxGasPriceWei = *assets.GWei(1000)  // Fix the gas price
	arbitrumMainnet.MinGasPriceWei = *assets.GWei(1000)  // Fix the gas price
	arbitrumMainnet.GasEstimatorMode = "FixedPrice"
	arbitrumMainnet.Layer2Type = "Arbitrum"
	arbitrumMainnet.BlockHistoryEstimatorBlockHistorySize = 0 // Force an error if someone set GAS_UPDATER_ENABLED=true by accident; we never want to run the block history estimator on arbitrum
	arbitrumMainnet.LinkContractAddress = "0xf97f4df75117a78c1A5a0DBb814Af92458539FB4"
	arbitrumMainnet.OCRContractConfirmations = 1
	arbitrumRinkeby := arbitrumMainnet
	arbitrumRinkeby.LinkContractAddress = "0x615fBe6372676474d9e6933d310469c9b68e9726"

	// Optimism is an L2 chain. Pending proper L2 support, for now we rely on their sequencer
	optimismMainnet := FallbackConfig
	optimismMainnet.BalanceMonitorBlockDelay = 0
	optimismMainnet.BlockHistoryEstimatorBlockHistorySize = 0 // Force an error if someone set GAS_UPDATER_ENABLED=true by accident; we never want to run the block history estimator on optimism
	optimismMainnet.EthTxResendAfterThreshold = 15 * time.Second
	optimismMainnet.FinalityDepth = 1    // Sequencer offers absolute finality as long as no re-org longer than 20 blocks occurs on main chain this event would require special handling (new txm)
	optimismMainnet.GasBumpThreshold = 0 // Never bump gas on optimism
	optimismMainnet.GasEstimatorMode = "Optimism"
	optimismMainnet.HeadTrackerHistoryDepth = 10
	optimismMainnet.HeadTrackerSamplingInterval = 1 * time.Second
	optimismMainnet.Layer2Type = "Optimism"
	optimismMainnet.LinkContractAddress = "" // TBD
	optimismMainnet.LinkContractAddress = "0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6"
	optimismMainnet.MinIncomingConfirmations = 1
	optimismMainnet.MinRequiredOutgoingConfirmations = 0
	optimismMainnet.OCRContractConfirmations = 1
	optimismKovan := optimismMainnet
	optimismKovan.LinkContractAddress = "0x4911b761993b9c8c0d14Ba2d86902AF6B0074F5B"
	optimismKovan.BlockEmissionIdleWarningThreshold = 30 * time.Minute

	// Fantom
	fantomMainnet := FallbackConfig
	fantomMainnet.GasPriceDefault = *assets.GWei(15)
	fantomMainnet.LinkContractAddress = "0x6f43ff82cca38001b6699a8ac47a2d0e66939407"
	fantomMainnet.MinIncomingConfirmations = 3
	fantomMainnet.MinRequiredOutgoingConfirmations = 2
	fantomTestnet := fantomMainnet
	fantomTestnet.LinkContractAddress = "0xfafedb041c0dd4fa2dc0d87a6b0979ee6fa7af5f"

	// RSK
	// RSK prices its txes in sats not wei
	rskMainnet := FallbackConfig
	rskMainnet.GasPriceDefault = *big.NewInt(50000000) // It's about 100 times more expensive than Wei, very roughly speaking
	rskMainnet.MaxGasPriceWei = *big.NewInt(50000000000)
	rskMainnet.MinGasPriceWei = *big.NewInt(0)
	rskMainnet.MinimumContractPayment = assets.NewLink(1000000000000000)
	rskMainnet.LinkContractAddress = "0x14adae34bef7ca957ce2dde5add97ea050123827"
	rskTestnet := rskMainnet
	rskTestnet.LinkContractAddress = "0x8bbbd80981fe76d44854d8df305e8985c19f0e78"

	// Avalanche
	avalancheMainnet := FallbackConfig
	avalancheMainnet.LinkContractAddress = "0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6" // TBD
	avalancheMainnet.FinalityDepth = 1
	avalancheMainnet.GasEstimatorMode = "FixedPrice"
	avalancheMainnet.GasPriceDefault = *big.NewInt(225000000000) // 225 Gwei
	avalancheMainnet.MaxGasPriceWei = *big.NewInt(225000000000)
	avalancheMainnet.MinGasPriceWei = *big.NewInt(225000000000)
	avalancheMainnet.MinIncomingConfirmations = 1
	avalancheMainnet.MinRequiredOutgoingConfirmations = 1
	avalancheMainnet.OCRContractConfirmations = 1

	avalancheFuji := avalancheMainnet
	avalancheFuji.LinkContractAddress = "0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846"

	EthMainnet.config = mainnet
	EthRinkeby.config = rinkeby
	EthGoerli.config = goerli
	EthKovan.config = kovan
	OptimismMainnet.config = optimismMainnet
	OptimismKovan.config = optimismKovan
	ArbitrumMainnet.config = arbitrumMainnet
	ArbitrumRinkeby.config = arbitrumRinkeby
	BSCMainnet.config = bscMainnet
	HecoMainnet.config = hecoMainnet
	FantomMainnet.config = fantomMainnet
	FantomTestnet.config = fantomTestnet
	PolygonMainnet.config = polygonMainnet
	PolygonMumbai.config = polygonMumbai
	XDaiMainnet.config = xDaiMainnet
	RSKMainnet.config = rskMainnet
	RSKTestnet.config = rskTestnet
	AvalancheFuji.config = avalancheFuji
	AvalancheMainnet.config = avalancheMainnet
}
