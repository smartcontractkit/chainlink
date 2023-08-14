package deployments

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var BetaChains = map[rhea.Chain]rhea.EvmDeploymentConfig{
	rhea.AvaxFuji:       {ChainConfig: Beta_AvaxFuji},
	rhea.OptimismGoerli: {ChainConfig: Beta_OptimismGoerli},
	rhea.Sepolia:        {ChainConfig: Beta_Sepolia},
	rhea.ArbitrumGoerli: {ChainConfig: Beta_ArbitrumGoerli},
}

var BetaChainMapping = map[rhea.Chain]map[rhea.Chain]rhea.EvmDeploymentConfig{
	rhea.AvaxFuji: {
		rhea.OptimismGoerli: Beta_AvaxFujiToOptimismGoerli,
		rhea.ArbitrumGoerli: Beta_AvaxFujiToArbitrumGoerli,
		rhea.Sepolia:        Beta_AvaxFujiToSepolia,
	},
	rhea.OptimismGoerli: {
		rhea.AvaxFuji:       Beta_OptimismGoerliToAvaxFuji,
		rhea.ArbitrumGoerli: Beta_OptimismGoerliToArbitrumGoerli,
		rhea.Sepolia:        Beta_OptimismGoerliToSepolia,
	},
	rhea.ArbitrumGoerli: {
		rhea.AvaxFuji:       Beta_ArbitrumGoerliToAvaxFuji,
		rhea.OptimismGoerli: Beta_ArbitrumGoerliToOptimismGoerli,
		rhea.Sepolia:        Beta_ArbitrumGoerliToSepolia,
	},
	rhea.Sepolia: {
		rhea.AvaxFuji:       Beta_SepoliaToAvaxFuji,
		rhea.OptimismGoerli: Beta_SepoliaToOptimismGoerli,
		rhea.ArbitrumGoerli: Beta_SepoliaToArbitrumGoerli,
	},
}

var Beta_Sepolia = rhea.EVMChainConfig{
	EvmChainId: 11155111,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: false,
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x779877A7B0D9E8603169DdbD7836e478b4624789"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
			PriceFeed: rhea.PriceFeed{
				Aggregator: gethcommon.HexToAddress("0x5A2734CC0341ea6564dF3D00171cc99C63B1A7d3"),
				Multiplier: big.NewInt(1e10),
			},
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x097D90c9d3E0B50Ca60e1ae45F6A81010f9FB534"),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
			PriceFeed: rhea.PriceFeed{
				Aggregator: gethcommon.HexToAddress("0x719E22E3D4b690E5d96cCb40619180B5427F14AE"),
				Multiplier: big.NewInt(1e10),
			},
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress(""),
	UpgradeRouter: gethcommon.HexToAddress(""),
	ARM:           gethcommon.HexToAddress(""),
	PriceRegistry: gethcommon.HexToAddress(""),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.Sepolia),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.Sepolia),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB,
		MaxGasPrice:              getMaxGasPrice(rhea.Sepolia),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
	ARMConfig: &arm_contract.ARMConfig{
		Voters: []arm_contract.ARMVoter{
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xcee3d4ac88c3ec196a1b09b5142c0f83d7e5bd61"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xefb9e871f768a9535dca47946a44d73f074e0ca1"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x9bd9bfe1c863c50c732dfd1e4b2c7d6dec1dd293"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xaff9ad86d7c119a9707a95888ce1a0955b5af8f1"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x3b809ea936b365438d371e6794f250536f1dda0d"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xd911d677c758201e84e55a00f232f58b6587e978"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x61183156e620fe24fad874a936d2cdfebb7b9cc5"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x285d081480996bbe5dd11256f0d09321c7e93a58"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x3436679441013bc0bc4a2a9b091ec76a447e8e35"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x1b94b22f4f0623f3bc3ac91f1a4112c65f942f38"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x809de34394cda8a78185e7f8e7ad7a3f7a699813"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x7dab8add75b3c0900895432fe8816d232317b8b6"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000006"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
		},
		BlessWeightThreshold: 2,
		CurseWeightThreshold: 2,
	},
	DeploySettings: rhea.ChainDeploySettings{
		DeployARM:           false,
		DeployTokenPools:    false,
		DeployRouter:        false,
		DeployUpgradeRouter: false,
		DeployPriceRegistry: false,
	},
}

var Beta_OptimismGoerli = rhea.EVMChainConfig{
	EvmChainId: 420,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: false,
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:          gethcommon.HexToAddress("0xdc2CC710e42857672E7907CF474a69B63B93089f"),
			Pool:           gethcommon.HexToAddress(""),
			TokenPoolType:  rhea.LockRelease,
			TokenPriceType: rhea.PriceFeeds,
			Price:          rhea.LINK.Price(),
			Decimals:       rhea.LINK.Decimals(),
			PriceFeed: rhea.PriceFeed{
				Aggregator: gethcommon.HexToAddress("0x53AFfFfA77006432146b667C67FA77b5D405793b"),
				Multiplier: big.NewInt(1e10),
			},
		},
		rhea.WETH: {
			Token:          gethcommon.HexToAddress("0x4200000000000000000000000000000000000006"),
			Price:          rhea.WETH.Price(),
			Decimals:       rhea.WETH.Decimals(),
			TokenPoolType:  rhea.FeeTokenOnly,
			TokenPriceType: rhea.TokenPrices,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress(""),
	UpgradeRouter: gethcommon.HexToAddress(""),
	ARM:           gethcommon.HexToAddress(""),
	PriceRegistry: gethcommon.HexToAddress(""),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.OptimismGoerli),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.OptimismGoerli),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.OptimismGoerli),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
	ARMConfig: &arm_contract.ARMConfig{
		Voters: []arm_contract.ARMVoter{
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x5318adc442a78983ad51f6a04de37f1ec5161ce5"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xa4bc432862724ac4f89f09ac49a7ba5f5d2f86c4"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xa2235d4e15fb6a7bc5544f640bc9e86459fe7fd9"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x7b3eff7928a252d0117d94994f11b2f5662f854c"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x119fb25b52e579fa3b73516f8957b9f01e363357"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xdf66bffcadb3485812d7ff3ca9cc345e2d94d4a5"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x3c560408bd37d480bc6ba0bc809d432250d3bd0a"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x35c376146bda951dd8988c180ef71d00e157f318"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x76ed187733982c10641fa6e58871d6f77cf0d4e1"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x992bbf33efc78d4e034d48a5aa374a319c6e78fe"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x67b22898dfa78efccde2dcdc301cdc5a6bbc77b8"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x2ce367fb5014e5f0e8b383eff7bf7bb74391296e"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000006"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
		},
		BlessWeightThreshold: 2,
		CurseWeightThreshold: 2,
	},
	DeploySettings: rhea.ChainDeploySettings{
		DeployARM:           false,
		DeployTokenPools:    false,
		DeployRouter:        false,
		DeployPriceRegistry: false,
		DeployUpgradeRouter: false,
	},
}

var Beta_AvaxFuji = rhea.EVMChainConfig{
	EvmChainId: 43113,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: false,
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:          gethcommon.HexToAddress("0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846"),
			Pool:           gethcommon.HexToAddress(""),
			Price:          rhea.LINK.Price(),
			Decimals:       rhea.LINK.Decimals(),
			TokenPoolType:  rhea.LockRelease,
			TokenPriceType: rhea.TokenPrices,
			PriceFeed: rhea.PriceFeed{
				Aggregator: gethcommon.HexToAddress("0x5F4a4f309Aefb6fb0Ab927A0421D0342fF92f194"),
				Multiplier: big.NewInt(1e10),
			},
		},
		rhea.WAVAX: {
			Token:          gethcommon.HexToAddress("0xd00ae08403B9bbb9124bB305C09058E32C39A48c"),
			Price:          rhea.WAVAX.Price(),
			Decimals:       rhea.WAVAX.Decimals(),
			TokenPoolType:  rhea.FeeTokenOnly,
			TokenPriceType: rhea.TokenPrices,
			PriceFeed: rhea.PriceFeed{
				Aggregator: gethcommon.HexToAddress("0x6C2441920404835155f33d88faf0545B895871b1"),
				Multiplier: big.NewInt(1e10),
			},
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WAVAX},
	WrappedNative: rhea.WAVAX,
	Router:        gethcommon.HexToAddress(""),
	UpgradeRouter: gethcommon.HexToAddress(""),
	ARM:           gethcommon.HexToAddress(""),
	PriceRegistry: gethcommon.HexToAddress(""),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.AvaxFuji),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.AvaxFuji),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.AvaxFuji),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
	ARMConfig: &arm_contract.ARMConfig{
		Voters: []arm_contract.ARMVoter{
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x22ad1e7ca671f07b502654a4a9605bd440a374ae"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x949754c662d24079f6ed024ae8cc152a03e02a83"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x7c10ef26cc03093ac39629da2c154b95b9093378"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x2a31f2831e9ad6fd7b24bb1d47637e59052b03c1"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xf3bb58c4595e39eec9518bfc5c3651e293e784e8"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x042f5b5f61697c155d2ef45e497250a01b0b8db9"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x8e7d8c061dcd9c86260a187d0a61d5c2e0ee09c7"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x306ca690533d30b418f90740265074186b532249"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xe99c75aad15e1f59d9e18fca9ac9e7d61c00c2a9"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xd02d00f3846860371bead4c8bf4a44b673aef319"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x9f82b506ac8e9c08f6b37afcceae93c34557474f"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x99271e6e4c88d4ac23510debb6ca4ff0412d5fba"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000006"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
		},
		BlessWeightThreshold: 2,
		CurseWeightThreshold: 2,
	},
	DeploySettings: rhea.ChainDeploySettings{
		DeployARM:           false,
		DeployTokenPools:    false,
		DeployRouter:        false,
		DeployUpgradeRouter: false,
		DeployPriceRegistry: false,
	},
}

var Beta_ArbitrumGoerli = rhea.EVMChainConfig{
	EvmChainId: 421613,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: true,
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:          gethcommon.HexToAddress("0xd14838A68E8AFBAdE5efb411d5871ea0011AFd28"),
			Pool:           gethcommon.HexToAddress(""),
			Price:          rhea.LINK.Price(),
			Decimals:       rhea.LINK.Decimals(),
			TokenPoolType:  rhea.LockRelease,
			TokenPriceType: rhea.TokenPrices,
			PriceFeed: rhea.PriceFeed{
				Aggregator: gethcommon.HexToAddress("0xb1D4538B4571d411F07960EF2838Ce337FE1E80E"),
				Multiplier: big.NewInt(1e10),
			},
		},
		rhea.WETH: {
			Token:          gethcommon.HexToAddress("0x32d5D5978905d9c6c2D4C417F0E06Fe768a4FB5a"),
			Price:          rhea.WETH.Price(),
			Decimals:       rhea.WETH.Decimals(),
			TokenPoolType:  rhea.FeeTokenOnly,
			TokenPriceType: rhea.TokenPrices,
			PriceFeed: rhea.PriceFeed{
				Aggregator: gethcommon.HexToAddress("0xC975dEfb12C5e83F2C7E347831126cF136196447"),
				Multiplier: big.NewInt(1e10),
			},
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress(""),
	UpgradeRouter: gethcommon.HexToAddress(""),
	ARM:           gethcommon.HexToAddress(""),
	PriceRegistry: gethcommon.HexToAddress(""),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.ArbitrumGoerli),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.ArbitrumGoerli),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.ArbitrumGoerli),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
	ARMConfig: &arm_contract.ARMConfig{
		Voters: []arm_contract.ARMVoter{
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xc8e7532913d78f5b2874e48455e5df367116524e"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xbf16a12b390f40e2dcbb4fb533a699346102e9fb"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x482520152d83daa7005f547ed2019675f7581b01"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xec61b7ddc4fb36f06a8b066b16513ac7519c162f"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x4f6af3a1cf3d83c537164f5c4792e70e77dd4cf3"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x39cb1995898baca087335558d0aef66589172cf5"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x76c7d8712b6f8539fba326519e39844bab19bf32"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x8bfa8c4db766be6d8c5cded7000a292d06e9848e"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x4fd9c1f9d3b74ea2980d5d5f7b26f2dcb627cb9d"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x6280eec1ca26c07f986df51aa7b303a00fdb2e8a"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x2eb4b9c1f512f6ab04afbef99cdb44fcf9f3cb5e"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x7814f9a0a2687fec23134852b9cacae03e50acbf"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000006"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
		},
		BlessWeightThreshold: 2,
		CurseWeightThreshold: 2,
	},
	DeploySettings: rhea.ChainDeploySettings{
		DeployARM:           false,
		DeployTokenPools:    false,
		DeployRouter:        false,
		DeployUpgradeRouter: false,
		DeployPriceRegistry: false,
	},
}

var Beta_OptimismGoerliToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    10655209,
		},
	},
}

var Beta_AvaxFujiToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    23076286,
		},
	},
}

var Beta_ArbitrumGoerliToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_ArbitrumGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    25964429,
		},
	},
}

var Beta_AvaxFujiToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    23077287,
		},
	},
}

var Beta_ArbitrumGoerliToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_ArbitrumGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    25966637,
		},
	},
}

var Beta_OptimismGoerliToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    10656674,
		},
	},
}

var Beta_AvaxFujiToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    23078007,
		},
	},
}

var Beta_SepoliaToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3689886,
		},
	},
}

var Beta_OptimismGoerliToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    10657183,
		},
	},
}

var Beta_SepoliaToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3689917,
		},
	},
}

var Beta_ArbitrumGoerliToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_ArbitrumGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    25970391,
		},
	},
}

var Beta_SepoliaToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Beta_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3689948,
		},
	},
}
