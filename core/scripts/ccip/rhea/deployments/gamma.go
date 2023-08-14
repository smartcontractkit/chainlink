package deployments

import (
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
)

var GammaChains = map[rhea.Chain]rhea.EvmDeploymentConfig{
	rhea.Ethereum: {ChainConfig: Prod_Ethereum},
	rhea.Optimism: {ChainConfig: Prod_Optimism},
	rhea.Arbitrum: {ChainConfig: Prod_Arbitrum},
	rhea.Avax:     {ChainConfig: Prod_Avax},
	rhea.Polygon:  {ChainConfig: Prod_Polygon},
}

var GammaChainMapping = map[rhea.Chain]map[rhea.Chain]rhea.EvmDeploymentConfig{
	rhea.Ethereum: {
		rhea.Avax:     Prod_EthereumToAvax,
		rhea.Optimism: Prod_EthereumToOptimism,
		rhea.Arbitrum: Prod_EthereumToArbitrum,
		rhea.Polygon:  Prod_EthereumToPolygon,
	},
	rhea.Avax: {
		rhea.Ethereum: Prod_AvaxToEthereum,
	},
	rhea.Optimism: {
		rhea.Ethereum: Prod_OptimismToEthereum,
	},
	rhea.Arbitrum: {
		rhea.Ethereum: Prod_ArbitrumToEthereum,
	},
	rhea.Polygon: {
		rhea.Ethereum: Prod_PolygonToEthereum,
	},
}

var Prod_Ethereum = rhea.EVMChainConfig{
	EvmChainId: 1,
	GasSettings: rhea.EVMGasSettings{
		EIP1559:   true,
		GasTipCap: rhea.DefaultGasTipFee,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x514910771AF9Ca656af840dff83E8264EcF986CA"),
			Pool:          gethcommon.HexToAddress("0x97AfF091eF4eb2AF981b9f50980aaEeb9cc80248"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
			Pool:          gethcommon.HexToAddress("0x456b59a6AC83213cD615C1867750CE0B28c1e87B"),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress("0x879FbE2c7943FdBDD4c14d71aED9FF65959ADDb9"),
	ARM:           gethcommon.HexToAddress("0x582d8f851fDc129901020E3F73aFAA8cb00423Eb"),
	PriceRegistry: gethcommon.HexToAddress("0xbE7C0Cf5C1464b33E24CE6244249D3AB76aCB0C3"),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.Ethereum),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.Ethereum),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.Ethereum),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
}

var Prod_Optimism = rhea.EVMChainConfig{
	EvmChainId: 10,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: true,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x350a791Bfc2C21F9Ed5d10980Dad2e2638ffa7f6"),
			Pool:          gethcommon.HexToAddress("0x7C8B10Fc5d45fFcF312773419A987Eb6B9a4a11e"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x4200000000000000000000000000000000000006"),
			Pool:          gethcommon.HexToAddress("0xC1e6A1271aa94697EF9eb0D33774cDE690a1e97B"),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress("0xDDfdfe69fD50aD26D3C7377623C04676CDB87555"),
	ARM:           gethcommon.HexToAddress("0x86034fb1Cc53dDF2CAD9622Ae829E9E2Fd47BDB8"),
	PriceRegistry: gethcommon.HexToAddress("0x2A2b7883bCBA00eB55702a939aE894709E89f589"),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.Optimism),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.Optimism),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.Optimism),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
}

var Prod_Avax = rhea.EVMChainConfig{
	EvmChainId: 43114,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: false,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x5947BB275c521040051D82396192181b413227A3"),
			Pool:          gethcommon.HexToAddress("0x2E53aC88A31A8568d423a012D4B71A36EfBaAB13"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WAVAX: {
			Token:         gethcommon.HexToAddress("0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7"),
			Pool:          gethcommon.HexToAddress("0xbc7Ea0B917cC1aD5683C6e65B67EFFb315545abf"),
			Price:         rhea.WAVAX.Price(),
			Decimals:      rhea.WAVAX.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WAVAX},
	WrappedNative: rhea.WAVAX,
	Router:        gethcommon.HexToAddress("0xb05844B1a7acE6b522F9733Fc67821f9E2b83601"),
	ARM:           gethcommon.HexToAddress("0x49911558C7D3496fA8bbf3Aa3D9A459236E09948"),
	PriceRegistry: gethcommon.HexToAddress("0xd76Bc4A633A2a56aE1F042b07F42EB9bdCbd2377"),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.Avax),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.Avax),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.Avax),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
}

var Prod_Arbitrum = rhea.EVMChainConfig{
	EvmChainId: 42161,
	GasSettings: rhea.EVMGasSettings{
		EIP1559:   true,
		GasTipCap: rhea.DefaultGasTipFee,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0xf97f4df75117a78c1A5a0DBb814Af92458539FB4"),
			Pool:          gethcommon.HexToAddress("0x0439AFf68a398cBA60Cd8dFD18a453803D1F57cE"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x82aF49447D8a07e3bd95BD0d56f35241523fBab1"),
			Pool:          gethcommon.HexToAddress("0x9aa615d204bfeaF7e1681FD811824e3c4739AC0d"),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress("0x9d7F8743c24aCFb88Bf862E1276820BEd8C9792b"),
	ARM:           gethcommon.HexToAddress("0x6Fe76e54A1AC4dd90acc44b09C43447F79d1b0bD"),
	PriceRegistry: gethcommon.HexToAddress("0x8A93719a3BAed1D1537AEB247194199B5704f2d8"),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.Arbitrum),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.Arbitrum),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.Arbitrum),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
}

var Prod_Polygon = rhea.EVMChainConfig{
	EvmChainId: 137,
	GasSettings: rhea.EVMGasSettings{
		EIP1559:   true,
		GasTipCap: rhea.DefaultGasTipFee,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0xb0897686c545045aFc77CF20eC7A532E3120E0F1"),
			Pool:          gethcommon.HexToAddress("0x555c837e73a1BF378105910B90b0f1eFD8687a87"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WMATIC: {
			Token:         gethcommon.HexToAddress("0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270"),
			Pool:          gethcommon.HexToAddress("0x5dAc2d2CBF3102Ca1848DA174265a2C329D188ed"),
			Price:         rhea.WMATIC.Price(),
			Decimals:      rhea.WMATIC.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WMATIC},
	WrappedNative: rhea.WMATIC,
	Router:        gethcommon.HexToAddress("0x75C682287ff98881dbD2ACB00F9E788658afbB0d"),
	ARM:           gethcommon.HexToAddress("0xC94EcD6A8b0190F2CaDb9097F68DE9372610701f"),
	PriceRegistry: gethcommon.HexToAddress("0xC6F837904da7CF1f46ADcB6ced39697d3B49ec90"),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.Polygon),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.Polygon),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.Polygon),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
}

var Prod_EthereumToAvax = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Ethereum,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x5B0c089ABf23b9f078a59C683aBf015773F25f66"),
		OffRamp:      gethcommon.HexToAddress("0x689A036eB17bEE3d4AD451E72855377E96c04175"),
		CommitStore:  gethcommon.HexToAddress("0x26f247A4Bd8dbF18675cf9E0C8fFBa654830336a"),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3491247,
		},
	},
}

var Prod_AvaxToEthereum = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Avax,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x82834E4D676a1d7A1e1969d0356515E973e6b460"),
		OffRamp:      gethcommon.HexToAddress("0x468462860d4d66C385F131c4ac9f2842D6FFc4AD"),
		CommitStore:  gethcommon.HexToAddress("0x9166cb4C167a9E63c0DD290c8fE1aBEAb24b3227"),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    21936491,
		},
	},
}

var Prod_EthereumToOptimism = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Ethereum,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x1c0C7858c7aD7a6B3f4aA813ee81E56D7405c712"),
		OffRamp:      gethcommon.HexToAddress("0xe4DB0cc096674eAC31332bc13Ca6084eF3D936BF"),
		CommitStore:  gethcommon.HexToAddress("0x7876C28E8a5Ca615046507DE36318C21562F4199"),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3491307,
		},
	},
}

var Prod_OptimismToEthereum = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Optimism,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0xf06A2e32477363bCAcAe5b86479e176Ca83D3f9d"),
		OffRamp:      gethcommon.HexToAddress("0x8Aa45E35Fa2142c04C10675d4D6D0ca9b3Fd5964"),
		CommitStore:  gethcommon.HexToAddress("0x7fa6DBA8Bab4aa0b0F6cb03632F4781d7F606427"),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    9366229,
		},
	},
}

var Prod_ArbitrumToEthereum = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Arbitrum,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x6264F5C5Bc1C0201159A5Bcd6486d9c6C2f75439"),
		OffRamp:      gethcommon.HexToAddress("0x16afbD47a0F851C78e43B7dB4932aA3efA60de41"),
		CommitStore:  gethcommon.HexToAddress("0x179592d39135D33a6Fc82a1678D2402a3CF7c151"),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    19905564,
		},
	},
}

var Prod_EthereumToArbitrum = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Ethereum,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0xe5C9121cc3f796A8446B9d35B0D53B67EB4c1Ab2"),
		OffRamp:      gethcommon.HexToAddress("0x282cE350aa31f068409d6BE82355e25267AA1cBF"),
		CommitStore:  gethcommon.HexToAddress("0x1B63f1372827feBb500F53e3DA7d9285F6bbEBC8"),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3491412,
		},
	},
}

var Prod_PolygonToEthereum = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Polygon,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x84347C236f4D4fb27673929899E554AB1151aa73"),
		OffRamp:      gethcommon.HexToAddress("0x3D0031BedE92258e55AbfD15bc74786ce71D2eae"),
		CommitStore:  gethcommon.HexToAddress("0xa108f62DfE10bb5cDe5ff39d7ebB86688fFbFa19"),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    35662240,
		},
	},
}

var Prod_EthereumToPolygon = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Ethereum,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x70349b74888B1364EE4862D9eF8cb1Af7ab47464"),
		OffRamp:      gethcommon.HexToAddress("0x7449e1074BBDe836fA8E74AaB51cCe8e66D1d902"),
		CommitStore:  gethcommon.HexToAddress("0xb3E3bcA2a8ea266d0436D2038EA636C1d83DA3c1"),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3497174,
		},
	},
}
