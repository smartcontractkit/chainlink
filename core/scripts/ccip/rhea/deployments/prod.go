package deployments

import (
	gethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip/rhea"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

var ProdChains = map[rhea.Chain]rhea.EvmDeploymentConfig{
	rhea.AvaxFuji:       {ChainConfig: Prod_AvaxFuji},
	rhea.OptimismGoerli: {ChainConfig: Prod_OptimismGoerli},
	rhea.Sepolia:        {ChainConfig: Prod_Sepolia},
	rhea.ArbitrumGoerli: {ChainConfig: Prod_ArbitrumGoerli},
	rhea.PolygonMumbai:  {ChainConfig: Prod_PolygonMumbai},
}

var ProdChainMapping = map[rhea.Chain]map[rhea.Chain]rhea.EvmDeploymentConfig{
	rhea.Sepolia: {
		rhea.AvaxFuji:       Prod_SepoliaToAvaxFuji,
		rhea.OptimismGoerli: Prod_SepoliaToOptimismGoerli,
		rhea.ArbitrumGoerli: Prod_SepoliaToArbitrumGoerli,
		rhea.PolygonMumbai:  Prod_SepoliaToPolygonMumbai,
	},
	rhea.AvaxFuji: {
		rhea.Sepolia:        Prod_AvaxFujiToSepolia,
		rhea.OptimismGoerli: Prod_AvaxFujiToOptimismGoerli,
		rhea.PolygonMumbai:  Prod_AvaxFujiToPolygonMumbai,
	},
	rhea.OptimismGoerli: {
		rhea.Sepolia:        Prod_OptimismGoerliToSepolia,
		rhea.AvaxFuji:       Prod_OptimismGoerliToAvaxFuji,
		rhea.ArbitrumGoerli: Prod_OptimismGoerliToArbitrumGoerli,
	},
	rhea.ArbitrumGoerli: {
		rhea.Sepolia:        Prod_ArbitrumGoerliToSepolia,
		rhea.OptimismGoerli: Prod_ArbitrumGoerliToOptimismGoerli,
	},
	rhea.PolygonMumbai: {
		rhea.Sepolia:  Prod_PolygonMumbaiToSepolia,
		rhea.AvaxFuji: Prod_PolygonMumbaiToAvaxFuji,
	},
}

var Prod_Sepolia = rhea.EVMChainConfig{
	EvmChainId: 11155111,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: false,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.A_DC: {
			Token:         gethcommon.HexToAddress("0x5B0bCa59dB458c159e5CbbE977119797F290F355"),
			Pool:          gethcommon.HexToAddress("0x660b0f3feacd3a6de68c28b091f4548f6f75b457"),
			Price:         rhea.A_DC.Price(),
			Decimals:      rhea.A_DC.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.Alongside: {
			Token:         gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"),
			Pool:          gethcommon.HexToAddress("0xac8cfc3762a979628334a0e4c1026244498e821b"),
			Price:         rhea.Alongside.Price(),
			Decimals:      rhea.Alongside.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.BankToken: {
			Token:         gethcommon.HexToAddress("0x784c400D6fF625051d2f587dC0276E3A1ffD9cda"),
			Pool:          gethcommon.HexToAddress("0x5f217ce93e206d6f13b342aeef53a084fa957745"),
			Price:         rhea.BankToken.Price(),
			Decimals:      rhea.BankToken.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.BetSwirl: {
			Token:         gethcommon.HexToAddress("0x94025780a1ab58868d9b2dbbb775f44b32e8e6e5"),
			Pool:          gethcommon.HexToAddress("0x3A1f9cc20b1301F5bbB1C98374c0FBdf6583CEDd"),
			Price:         rhea.BetSwirl.Price(),
			Decimals:      rhea.BetSwirl.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.BondToken: {
			Token:         gethcommon.HexToAddress("0xF92E4b278380f39fADc24483C7baC61b73EE93F2"),
			Pool:          gethcommon.HexToAddress("0x919b1d308e4477c88350c336537ec5ac9ee76d9a"),
			Price:         rhea.BondToken.Price(),
			Decimals:      rhea.BondToken.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.CACHEGOLD: {
			Token:         gethcommon.HexToAddress("0x94095e6514411C65E7809761F21eF0febe69A977"),
			Pool:          gethcommon.HexToAddress("0x5a80462d5e15fcf75cacf8d0dbbf16b70476d029"),
			Price:         rhea.CACHEGOLD.Price(),
			Decimals:      rhea.CACHEGOLD.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.CCIP_BnM: {
			// NOTE this should be the custom burn_mint_erc677_helper contract, not the default burn_mint_erc677
			Token:    gethcommon.HexToAddress("0xFd57b4ddBf88a4e07fF4e34C487b99af2Fe82a05"),
			Pool:     gethcommon.HexToAddress("0x38d1ef9619cd40cf5482c045660ae7c82ada062c"),
			Price:    rhea.CCIP_BnM.Price(),
			Decimals: rhea.CCIP_BnM.Decimals(),
			// Wrapped is used to ensure new pool deployments will automatically grant burn/mint permissions
			TokenPoolType: rhea.Wrapped,
		},
		rhea.FUGAZIUSDC: {
			Token:         gethcommon.HexToAddress("0x832bA6abcAdC68812be372F4ef20aAC268bA20B7"),
			Pool:          gethcommon.HexToAddress("0x0ea0d7b2b78dd3a926fc76d6875a287f0aeb158f"),
			Price:         rhea.FUGAZIUSDC.Price(),
			Decimals:      rhea.FUGAZIUSDC.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.InsurAce: {
			Token:         gethcommon.HexToAddress("0xb7c8bCA891143221a34DB60A26639785C4839040"),
			Pool:          gethcommon.HexToAddress("0xa04c2cbbfa7bf7adcbde911216a0ba7e3f1e36b3"),
			Price:         rhea.InsurAce.Price(),
			Decimals:      rhea.InsurAce.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x779877A7B0D9E8603169DdbD7836e478b4624789"),
			Pool:          gethcommon.HexToAddress("0x5344b4bf5ae39038a591866d2853b2b1db622911"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.NZ_DC: {
			Token:         gethcommon.HexToAddress("0x72C6333A5A99BCB3394DcCd879d6D8FE8766A297"),
			Pool:          gethcommon.HexToAddress("0x21509dfda83a72e444cc18bc57e6961f1af93959"),
			Price:         rhea.NZ_DC.Price(),
			Decimals:      rhea.NZ_DC.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.SG_DC: {
			Token:         gethcommon.HexToAddress("0x1c1383CFdb6D64e884696C533c6D8A8c42033D79"),
			Pool:          gethcommon.HexToAddress("0x2a0e16a7d9a027f0aa77a64d362b05c35824db0d"),
			Price:         rhea.SG_DC.Price(),
			Decimals:      rhea.SG_DC.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.STEADY: {
			Token:         gethcommon.HexToAddress("0x82abB1864326A8A7e1A357FFA2270D09CCb867B9"),
			Pool:          gethcommon.HexToAddress("0x5c0b55dbd1335a7c96653788cf545a8c08148496"),
			Price:         rhea.STEADY.Price(),
			Decimals:      rhea.STEADY.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x097D90c9d3E0B50Ca60e1ae45F6A81010f9FB534"),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.CCIP_LnM: {
			// NOTE this should be the custom burn_mint_erc677_helper contract, not the default burn_mint_erc677
			Token:    gethcommon.HexToAddress("0x466D489b6d36E7E3b824ef491C225F5830E81cC1"),
			Pool:     gethcommon.HexToAddress("0x3637220fccd067927766a40475f2e8fade33f590"),
			Price:    rhea.CCIP_LnM.Price(),
			Decimals: rhea.CCIP_LnM.Decimals(),
			// Wrapped is used to ensure new pool deployments will automatically grant burn/mint permissions
			TokenPoolType: rhea.Wrapped,
		},
		rhea.SNXUSD: {
			Token:         gethcommon.HexToAddress("0x1b791d05E437C78039424749243F5A79E747525e"),
			Pool:          gethcommon.HexToAddress("0x9b65749b38278060c5787cce0391ac7f1094c8e8"),
			Price:         rhea.SNXUSD.Price(),
			Decimals:      rhea.SNXUSD.Decimals(),
			TokenPoolType: rhea.BurnMint,
			PoolAllowList: []gethcommon.Address{
				gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"),
				gethcommon.HexToAddress("0x2A45BaE1E58AaD3261af187b7dAde90889c039Dc"),
				gethcommon.HexToAddress("0x76490713314fCEC173f44e99346F54c6e92a8E42"), // BetaUser - Synthetix v3 core
			},
		},
		//rhea.ZUSD: {
		//	Token:         gethcommon.HexToAddress("0x09ae935D80E190403C61Cc5d854Fbf6a7b4a559a"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.ZUSD.Price(),
		//	Decimals:      rhea.ZUSD.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
		//rhea.SUPER: {
		//	Token:         gethcommon.HexToAddress("0xCb4B3f72B5b6D0b7072aFDDf18FE61A0d569EC39"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.SUPER.Price(),
		//	Decimals:      rhea.SUPER.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress("0xd0daae2231e9cb96b94c8512223533293c3693bf"),
	ARM:           gethcommon.HexToAddress("0xb4d360459f32dd641ef5a6985ffbac5c4e5521aa"),
	ARMProxy:      gethcommon.HexToAddress("0xba3f6251de62ded61ff98590cb2fdf6871fbb991"),
	PriceRegistry: gethcommon.HexToAddress("0x8737a1c3d55779d03b7a08188e97af87b4110946"),
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
			// Infra-testnet-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x8e327be2a8bb2e95b7e281ec8fbdb327ea7cbbb1"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x1936092090584fdf4542df8cc9b3ba695ef2cf88"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x671dea470d173e3a3fab9a463a9d85c9032da5e5"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xca1dc0e1ef2a413c4672f3dfa28922a097ba32a2"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x4679297e452b4b09ff2e351ddac3eff9c7999a17"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xaa4d66d0db8ac802ba5ecfa4291c2b7aabe23a3f"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x254090c0355c60aa7409c362196f65925393760e"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x916bfcdb65d7216e869fea39eb6bbc5b15e61768"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-3
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x94b7a50a85e1127cdeab48915d87fa930314e7d9"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xf109b0a93d2d352e117b5066605116fe87faf344"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-4
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x8d708fbb28b1f39ded877972a952c77346df5fcc"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xf0001cab768289b64151e544ea864966b3985c0f"),
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
	},
}

var Prod_OptimismGoerli = rhea.EVMChainConfig{
	EvmChainId: 420,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: true,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.Alongside: {
			Token:         gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"),
			Pool:          gethcommon.HexToAddress("0x8bcd622ac003160ea239c82e1b0e09364d77b1ac"),
			Price:         rhea.Alongside.Price(),
			Decimals:      rhea.Alongside.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.CACHEGOLD: {
			Token:         gethcommon.HexToAddress("0x94095e6514411C65E7809761F21eF0febe69A977"),
			Pool:          gethcommon.HexToAddress("0x60a434ae77d30c2e1d737fa50bc20753a621e0b6"),
			Price:         rhea.CACHEGOLD.Price(),
			Decimals:      rhea.CACHEGOLD.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.CCIP_BnM: {
			// NOTE this should be the custom burn_mint_erc677_helper contract, not the default burn_mint_erc677
			Token:    gethcommon.HexToAddress("0xaBfE9D11A2f1D61990D1d253EC98B5Da00304F16"),
			Pool:     gethcommon.HexToAddress("0x8668ab4eb1dffe11db7491ebce633b050bb29cda"),
			Price:    rhea.CCIP_BnM.Price(),
			Decimals: rhea.CCIP_BnM.Decimals(),
			// Wrapped is used to ensure new pool deployments will automatically grant burn/mint permissions
			TokenPoolType: rhea.Wrapped,
		},
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0xdc2CC710e42857672E7907CF474a69B63B93089f"),
			Pool:          gethcommon.HexToAddress("0xdecfaf632175915bdf38c00d9d9746e8a90a56c4"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x4200000000000000000000000000000000000006"),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.CCIP_LnM: {
			Token:         gethcommon.HexToAddress("0x835833d556299cdec623e7980e7369145b037591"),
			Pool:          gethcommon.HexToAddress("0xf66d20ac7b981e249fce8fb8ddae3974f5559735"),
			Price:         rhea.CCIP_LnM.Price(),
			Decimals:      rhea.CCIP_LnM.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		rhea.SNXUSD: {
			Token:         gethcommon.HexToAddress("0xe487Ad4291019b33e2230F8E2FB1fb6490325260"),
			Pool:          gethcommon.HexToAddress("0xd23c2ef3a533040b57cadaf33ccb111edbaca018"),
			Price:         rhea.SNXUSD.Price(),
			Decimals:      rhea.SNXUSD.Decimals(),
			TokenPoolType: rhea.BurnMint,
			PoolAllowList: []gethcommon.Address{
				gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"),
				gethcommon.HexToAddress("0x2A45BaE1E58AaD3261af187b7dAde90889c039Dc"),
				gethcommon.HexToAddress("0x76490713314fCEC173f44e99346F54c6e92a8E42"), // BetaUser - Synthetix v3 core
			},
		},

		//rhea.ZUSD: {
		//	Token:         gethcommon.HexToAddress("0x740ba2E7f25c036ED0b19b83c9Da2cB8D756f9D5"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.ZUSD.Price(),
		//	Decimals:      rhea.ZUSD.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
		//rhea.STEADY: {
		//	Token:         gethcommon.HexToAddress("0x615c83D5FEdafAEa641f1cC1a91ea09111EF0158"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.STEADY.Price(),
		//	Decimals:      rhea.STEADY.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress("0xeb52e9ae4a9fb37172978642d4c141ef53876f26"),
	ARM:           gethcommon.HexToAddress("0xeaf6968fab9c54ac31c3679f120705b5019d3546"),
	ARMProxy:      gethcommon.HexToAddress("0x4eb4dbdb3c3b56e5e209abf9c424a3834f2087d0"),
	PriceRegistry: gethcommon.HexToAddress("0x490f3b46fba6af0d7499867a73469a077251c2bb"),
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
			// Infra-testnet-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x4d9987720ec678aa1271621ffe617771288e436f"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x0e3594b19fb2b7ceb4e6872a9393b407579702b8"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xebe7b72da8ade2e1ed2077d51a933767029bf513"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xb83efe598c14c80004dc75d2728ecc52c0113315"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x7cf0592c4eda6b839b635ec0269df4d2e51ba1f4"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xf5c77f0fbf0be8559b6ddf752fca4342eb8e254b"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xc4e34bcc4b46b8e7fe02c8ac6fa8129f5027f3ef"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xb94388d24dde6e7155cecec5d6474a1ce14f0127"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-3
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x92b9eba2ab3e89b68a157cdb7e8146aa8e7a735c"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x4a26fec4350aab068a43f767468b83513a45ab52"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-4
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x3540e1f097a7d6bfce706224dcc965a26c5baa7c"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xc57589520ba1b98b0f09816b431f5866ebb58b90"),
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
	},
}

var Prod_AvaxFuji = rhea.EVMChainConfig{
	EvmChainId: 43113,
	GasSettings: rhea.EVMGasSettings{
		EIP1559: false,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.A_DC: {
			Token:         gethcommon.HexToAddress("0x2DFfDe4CEb3E17d27D19cE4add6b351044c3d290"),
			Pool:          gethcommon.HexToAddress("0xffef5d7868416491c3f7ebeee835a2872871e31f"),
			Price:         rhea.A_DC.Price(),
			Decimals:      rhea.A_DC.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.Alongside: {
			Token:         gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"),
			Pool:          gethcommon.HexToAddress("0xaef84b05e96e2aafac9a347e60ae7e9a414fb649"),
			Price:         rhea.Alongside.Price(),
			Decimals:      rhea.Alongside.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.BetSwirl: {
			Token:         gethcommon.HexToAddress("0x94025780a1ab58868d9b2dbbb775f44b32e8e6e5"),
			Pool:          gethcommon.HexToAddress("0xF4B0a2Ef2E77f981d5e2Ff45E3AFE2A199AD5EE2"),
			Price:         rhea.BetSwirl.Price(),
			Decimals:      rhea.BetSwirl.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.BankToken: {
			Token:         gethcommon.HexToAddress("0x0147cba76c478aa46a76b8e2d2fdbd789d63b773"),
			Pool:          gethcommon.HexToAddress("0xd12e98b53446048e5ec614df514bc6838c5a8010"),
			Price:         rhea.BankToken.Price(),
			Decimals:      rhea.BankToken.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		rhea.BondToken: {
			Token:         gethcommon.HexToAddress("0x8737a1c3d55779d03b7a08188e97af87b4110946"),
			Pool:          gethcommon.HexToAddress("0xfac166b229ca504c254bf89449da10d08d44cf69"),
			Price:         rhea.BondToken.Price(),
			Decimals:      rhea.BondToken.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		rhea.CACHEGOLD: {
			Token:         gethcommon.HexToAddress("0x94095e6514411C65E7809761F21eF0febe69A977"),
			Pool:          gethcommon.HexToAddress("0xa1fc992576a2cec26d8ca8c1240a82af0180da84"),
			Price:         rhea.CACHEGOLD.Price(),
			Decimals:      rhea.CACHEGOLD.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.CCIP_BnM: {
			// NOTE this should be the custom burn_mint_erc677_helper contract, not the default burn_mint_erc677
			Token:    gethcommon.HexToAddress("0xD21341536c5cF5EB1bcb58f6723cE26e8D8E90e4"),
			Pool:     gethcommon.HexToAddress("0xec1062cbdf4fbf31b3a6aac62b6f6f123bb70e12"),
			Price:    rhea.CCIP_BnM.Price(),
			Decimals: rhea.CCIP_BnM.Decimals(),
			// Wrapped is used to ensure new pool deployments will automatically grant burn/mint permissions
			TokenPoolType: rhea.Wrapped,
		},
		rhea.FUGAZIUSDC: {
			Token:         gethcommon.HexToAddress("0x150a0ee7393294442EE4d4F5C7d637af01dF93ee"),
			Pool:          gethcommon.HexToAddress("0x0040e5e502fe84de97b1b1cd7d33ca729d8b2a8b"),
			Price:         rhea.FUGAZIUSDC.Price(),
			Decimals:      rhea.FUGAZIUSDC.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.InsurAce: {
			Token:         gethcommon.HexToAddress("0x005d27b4ee87e1f7362916ffa54bc37a30729554"),
			Pool:          gethcommon.HexToAddress("0x96a0c308a3293f0a4425bab68f37bdf6661eedad"),
			Price:         rhea.InsurAce.Price(),
			Decimals:      rhea.InsurAce.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846"),
			Pool:          gethcommon.HexToAddress("0x658af0d8ecbb13c5fd5b545ac7316e50cc07cf6e"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.NZ_DC: {
			Token:         gethcommon.HexToAddress("0x3d1263bAa30c696f1d7eff8E63962674A43E5980"),
			Pool:          gethcommon.HexToAddress("0x9c03068935e61fc1070f8a1d7afe13f799422301"),
			Price:         rhea.NZ_DC.Price(),
			Decimals:      rhea.NZ_DC.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.SG_DC: {
			Token:         gethcommon.HexToAddress("0xA6C01A7E40e784934739BAC707DE1d6463e76cb7"),
			Pool:          gethcommon.HexToAddress("0xe359d83ed390c4143820230e18a5e535d93a9a6d"),
			Price:         rhea.SG_DC.Price(),
			Decimals:      rhea.SG_DC.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.WAVAX: {
			Token:         gethcommon.HexToAddress("0xd00ae08403B9bbb9124bB305C09058E32C39A48c"),
			Price:         rhea.WAVAX.Price(),
			Decimals:      rhea.WAVAX.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.CCIP_LnM: {
			Token:         gethcommon.HexToAddress("0x70f5c5c40b873ea597776da2c21929a8282a3b35"),
			Pool:          gethcommon.HexToAddress("0x583dbe5f15dea93f321826d856994e53e01cd498"),
			Price:         rhea.CCIP_LnM.Price(),
			Decimals:      rhea.CCIP_LnM.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		//rhea.SUPER: {
		//	Token:         gethcommon.HexToAddress("0xCb4B3f72B5b6D0b7072aFDDf18FE61A0d569EC39"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.SUPER.Price(),
		//	Decimals:      rhea.SUPER.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WAVAX},
	WrappedNative: rhea.WAVAX,
	Router:        gethcommon.HexToAddress("0x554472a2720e5e7d5d3c817529aba05eed5f82d8"),
	ARM:           gethcommon.HexToAddress("0x0ea0d7b2b78dd3a926fc76d6875a287f0aeb158f"),
	ARMProxy:      gethcommon.HexToAddress("0xac8cfc3762a979628334a0e4c1026244498e821b"),
	PriceRegistry: gethcommon.HexToAddress("0xe42ecce39ce5bd2bbf2443660ba6979eeafd48df"),
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
			// Infra-testnet-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xac591a80ff5a81c512a5bb52c77e2513eca77245"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xc5202b4be2f03ec895773c7ade8e79e9794ac214"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x722915bc6373d35bd051e1d61ff15edd2f0b0aae"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x829a26e035a0d5e217960db184f9742c076c9ddd"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x1874d82f4a25e2f2633106afd08baecbf3b52468"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xbf131623483b4f0ac00371cce9c4f3b59339390e"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x13e103d39e1970317e8d9dc05583bafb7b08f79e"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x112f239515839d8349a62fbfb48bae79ffbbad74"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-3
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xfd1d95eb730ac6895e585e10bfc4a67332887ae5"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x6ca0ef2b278839a505b548a9fa5e20026e9bbe82"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-4
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x1acbf34cd5e91784e918f1ea5c9cc2fd923c882b"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xc2150f9d6b38cefc0a95174026558eea6c05f4a3"),
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
	},
}

var Prod_ArbitrumGoerli = rhea.EVMChainConfig{
	EvmChainId: 421613,
	GasSettings: rhea.EVMGasSettings{
		EIP1559:   true,
		GasTipCap: rhea.DefaultGasTipFee,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.BetSwirl: {
			Token:         gethcommon.HexToAddress("0x94025780a1ab58868d9b2dbbb775f44b32e8e6e5"),
			Pool:          gethcommon.HexToAddress("0x3e6733c15199AB1058474642598dd52a3aE8237D"),
			Price:         rhea.BetSwirl.Price(),
			Decimals:      rhea.BetSwirl.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.CCIP_BnM: {
			// NOTE this should be the custom burn_mint_erc677_helper contract, not the default burn_mint_erc677
			Token:    gethcommon.HexToAddress("0x0579b4c1C8AcbfF13c6253f1B10d66896Bf399Ef"),
			Pool:     gethcommon.HexToAddress("0xf399f6a4ea83442f97f480118ebd56d1aed767b9"),
			Price:    rhea.CCIP_BnM.Price(),
			Decimals: rhea.CCIP_BnM.Decimals(),
			// Wrapped is used to ensure new pool deployments will automatically grant burn/mint permissions
			TokenPoolType: rhea.Wrapped,
		},
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0xd14838A68E8AFBAdE5efb411d5871ea0011AFd28"),
			Pool:          gethcommon.HexToAddress("0x044a6b4b561af69d2319a2f4be5ec327a6975d0a"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x32d5D5978905d9c6c2D4C417F0E06Fe768a4FB5a"),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.CCIP_LnM: {
			Token:         gethcommon.HexToAddress("0x0e14dbe2c8e1121902208be173a3fb91bb125cdb"),
			Pool:          gethcommon.HexToAddress("0xa77aefaba6161f907299dc2be79a60c9e80e9b91"),
			Price:         rhea.CCIP_LnM.Price(),
			Decimals:      rhea.CCIP_LnM.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		rhea.SNXUSD: {
			Token:         gethcommon.HexToAddress("0x1b791d05E437C78039424749243F5A79E747525e"),
			Pool:          gethcommon.HexToAddress("0xd7d47c0e62029a1a3eb8c08691c8c9863fe766c2"),
			Price:         rhea.SNXUSD.Price(),
			Decimals:      rhea.SNXUSD.Decimals(),
			TokenPoolType: rhea.BurnMint,
			PoolAllowList: []gethcommon.Address{
				gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"),
				gethcommon.HexToAddress("0x2A45BaE1E58AaD3261af187b7dAde90889c039Dc"),
				gethcommon.HexToAddress("0x76490713314fCEC173f44e99346F54c6e92a8E42"), // BetaUser - Synthetix v3 core
			},
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress("0x88e492127709447a5abefdab8788a15b4567589e"),
	ARM:           gethcommon.HexToAddress("0x8af4204e30565df93352fe8e1de78925f6664da7"),
	ARMProxy:      gethcommon.HexToAddress("0x3cc9364260d80f09ccac1ee6b07366db598900e6"),
	PriceRegistry: gethcommon.HexToAddress("0x114a20a10b43d4115e5aeef7345a1a71d2a60c57"),
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
			// Infra-testnet-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x9d2fefa5791ee383884df019e08d2fc307e776b1"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x7203cf710b0fa2128c13e12672286a287890ec22"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x3e8eee517dba13675fff1b4f2ea4210902cba81b"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x02cd3b5011567b70d41809907e672dcdd05285ee"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x0a6e792d5c6f813e399341740cf7a368b6d66e6f"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x263a09238c91cabe5e19642cd1c81fa567406c36"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x1db6a4237b54fdfc7bd3c50c055e468c85832fea"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x153a6331278bda690b4f7365d514a7a27017463a"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-3
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xec994ad324eb0c21c8184c059398094349d63693"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x989a3fdbff84130d919942cf65a0e3942bd6e5ed"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-4
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0xa2306e99f35fc1d150dcd78bf2af7f224e158fe0"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xd1ad57dd0401947b4e6eb48e6989c559e47ce7a2"),
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
	},
}

var Prod_PolygonMumbai = rhea.EVMChainConfig{
	EvmChainId: 80001,
	GasSettings: rhea.EVMGasSettings{
		EIP1559:   true,
		GasTipCap: rhea.DefaultGasTipFee,
	},
	AllowList: []gethcommon.Address{},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.CACHEGOLD: {
			Token:         gethcommon.HexToAddress("0x94095e6514411C65E7809761F21eF0febe69A977"),
			Pool:          gethcommon.HexToAddress("0xd12e98b53446048e5ec614df514bc6838c5a8010"),
			Price:         rhea.CACHEGOLD.Price(),
			Decimals:      rhea.CACHEGOLD.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
		rhea.CCIP_BnM: {
			// NOTE this should be the custom burn_mint_erc677_helper contract, not the default burn_mint_erc677
			Token:    gethcommon.HexToAddress("0xf1E3A5842EeEF51F2967b3F05D45DD4f4205FF40"),
			Pool:     gethcommon.HexToAddress("0xa6c88f12ae1aa9c333e86ccbdd2957cac2e5f58c"),
			Price:    rhea.CCIP_BnM.Price(),
			Decimals: rhea.CCIP_BnM.Decimals(),
			// Wrapped is used to ensure new pool deployments will automatically grant burn/mint permissions
			TokenPoolType: rhea.Wrapped,
		},
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x326C977E6efc84E512bB9C30f76E30c160eD06FB"),
			Pool:          gethcommon.HexToAddress("0x6fce09b2e74f649a4494a1844219cb0d86cfe8b7"),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.WMATIC: {
			Token:         gethcommon.HexToAddress("0x9c3C9283D3e44854697Cd22D3Faa240Cfb032889"),
			Price:         rhea.WMATIC.Price(),
			Decimals:      rhea.WMATIC.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.CCIP_LnM: {
			Token:         gethcommon.HexToAddress("0xc1c76a8c5bfde1be034bbcd930c668726e7c1987"),
			Pool:          gethcommon.HexToAddress("0x83369f8586ba000a87db278549b9a2370dc626b6"),
			Price:         rhea.CCIP_LnM.Price(),
			Decimals:      rhea.CCIP_LnM.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		rhea.SNXUSD: {
			Token:         gethcommon.HexToAddress("0x1b791d05E437C78039424749243F5A79E747525e"),
			Pool:          gethcommon.HexToAddress("0xb8b8592aaf82bd42190aa8b629c6afa35a433461"),
			Price:         rhea.SNXUSD.Price(),
			Decimals:      rhea.SNXUSD.Decimals(),
			TokenPoolType: rhea.BurnMint,
			PoolAllowList: []gethcommon.Address{
				gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"),
				gethcommon.HexToAddress("0x2A45BaE1E58AaD3261af187b7dAde90889c039Dc"),
				gethcommon.HexToAddress("0x76490713314fCEC173f44e99346F54c6e92a8E42"), // BetaUser - Synthetix v3 core
			},
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WMATIC},
	WrappedNative: rhea.WMATIC,
	Router:        gethcommon.HexToAddress("0x70499c328e1e2a3c41108bd3730f6670a44595d1"),
	ARM:           gethcommon.HexToAddress("0x917a6913f785094f8b06785aa8a884f922a650d8"),
	ARMProxy:      gethcommon.HexToAddress("0x235ce3408845a4767a2eaa2a3d8ef0848d283f1f"),
	PriceRegistry: gethcommon.HexToAddress("0x9bd312170aa145ef98453940dc9ab894235b063e"),
	TunableChainValues: rhea.TunableChainValues{
		FinalityDepth:            getFinalityDepth(rhea.PolygonMumbai),
		OptimisticConfirmations:  getOptimisticConfirmations(rhea.PolygonMumbai),
		BatchGasLimit:            BATCH_GAS_LIMIT,
		RelativeBoostPerWaitHour: RELATIVE_BOOST_PER_WAIT_HOUR,
		FeeUpdateHeartBeat:       models.MustMakeDuration(FEE_UPDATE_HEARTBEAT),
		FeeUpdateDeviationPPB:    FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN,
		MaxGasPrice:              getMaxGasPrice(rhea.PolygonMumbai),
		InflightCacheExpiry:      models.MustMakeDuration(INFLIGHT_CACHE_EXPIRY),
		RootSnoozeTime:           models.MustMakeDuration(ROOT_SNOOZE_TIME),
	},
	ARMConfig: &arm_contract.ARMConfig{
		Voters: []arm_contract.ARMVoter{
			// Infra-testnet-1
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x6b33111e07a15f51a82bf30708e95dc2169ec100"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x4e6239f0fdda2b81d4bc790c959caffcc47d8436"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000001"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-2
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x19efeb472fc308b777bef6282ad77688ff954181"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xf69d7ba7a60b148a08881744cd1a415582703ad5"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000002"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Kostis-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x8bccd4cd06f0b50e78da446ef0ac61f3b43aefc5"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xe5c553af74b9badb9e7d52ebf5a52a6c556fba10"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000003"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Xueyuan-0
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x76118c18bcaa561cea1bedf558cb9a11bfd7bf2c"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xa19eadc34225c8a521047a7027421e868e73a577"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000004"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-3
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x5495ce7608f1ab70784169f8619880c3c2ce054e"),
				CurseVoteAddr:   gethcommon.HexToAddress("0xc3d66c89008dd26c9e98a915245847a2a1ea09b2"),
				CurseUnvoteAddr: gethcommon.HexToAddress("0x0000000000000000000000000000000000000005"),
				BlessWeight:     1,
				CurseWeight:     1,
			},
			// Infra-testnet-4
			{
				BlessVoteAddr:   gethcommon.HexToAddress("0x70701d3ff4f6c56a36a44e6b1b2874d69cc4864f"),
				CurseVoteAddr:   gethcommon.HexToAddress("0x3ca8aac84e5543e0c6faa6312407d54001cf723d"),
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
	},
}

var Prod_SepoliaToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0xe42ecce39ce5bd2bbf2443660ba6979eeafd48df"),
		OffRamp:      gethcommon.HexToAddress("0x1f06781450e994b0005ce2922fca78e2c72d4353"),
		CommitStore:  gethcommon.HexToAddress("0xf3855a07bf75c4e0b4ddbeb7784badc9dd2ca274"),
		PingPongDapp: gethcommon.HexToAddress("0x7a7783e6073175f58db4d5f8bb40ea44065246db"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3816446,
		},
	},
}

var Prod_AvaxFujiToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0xa799c1855875e79b2e1752412058b485ee51aec4"),
		OffRamp:      gethcommon.HexToAddress("0x61c67e7b7c90ed1a44dabb26c33900270df7a144"),
		CommitStore:  gethcommon.HexToAddress("0xb4407405465a5dab21fe6e6b748b42a2dccc5e9d"),
		PingPongDapp: gethcommon.HexToAddress("0xc1f01a6d0e8382f2c5d394923a7e79693354934b"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    23753755,
		},
	},
}

var Prod_SepoliaToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x365408d655a6cdf8bc668fd6cebe1bb16a403ca6"),
		OffRamp:      gethcommon.HexToAddress("0x0d3299ee55d493b8d9aafc834a6fd5dcbc4a409a"),
		CommitStore:  gethcommon.HexToAddress("0x1576d23f986ecb572a4c839ba6758ca05c1eadc2"),
		PingPongDapp: gethcommon.HexToAddress("0x37b27863a14781acf41b787cf9ec3fb65d1c5885"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3816777,
		},
	},
}

var Prod_OptimismGoerliToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x64877f0b53e801adeb8d65f9706f7b134b82971c"),
		OffRamp:      gethcommon.HexToAddress("0xdc4606e96c37b877f2c9ddda82104c85a198a82d"),
		CommitStore:  gethcommon.HexToAddress("0xb7019c10bd604768c9cf5b3d086a2e661559a189"),
		PingPongDapp: gethcommon.HexToAddress("0x3b7a30028bf7ce52ad75b0afb142beef02deeecd"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    11482670,
		},
	},
}

var Prod_OptimismGoerliToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0xf33303166911ff86ad1a5ea94d459b4ba0ba8cc9"),
		OffRamp:      gethcommon.HexToAddress("0xee8ce182ea0c0edecf06c2a032a17b2058fc5a04"),
		CommitStore:  gethcommon.HexToAddress("0x2980de4ce178bc8bb6840abd2ef0e2a7c8e7272f"),
		PingPongDapp: gethcommon.HexToAddress("0x227c3699f9f0d6d55c38551a7d7feaea82efdd66"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    11482157,
		},
	},
}

var Prod_AvaxFujiToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x097076fbd8573418c77d2600606ad063c0e3cc7c"),
		OffRamp:      gethcommon.HexToAddress("0x0f287140d86335b37ae2ad0707992ecd4202d5b7"),
		CommitStore:  gethcommon.HexToAddress("0x5a7fa03e52628a0a6f0ab637f10ba45b68f9ad33"),
		PingPongDapp: gethcommon.HexToAddress("0x4ce7b0782966d58ebc5e1804ca6de3244dac9ad9"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    23755228,
		},
	},
}

var Prod_ArbitrumGoerliToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_ArbitrumGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0xc8b93b46bf682c39b3f65aa1c135bc8a95a5e43a"),
		OffRamp:      gethcommon.HexToAddress("0x7a0bb92bc8663abe6296d0162a9b41a2cb2e0358"),
		CommitStore:  gethcommon.HexToAddress("0x7eef73aca8657aaefd509a97ee75aa6740046e75"),
		PingPongDapp: gethcommon.HexToAddress("0x9b451300c94c7328bdb56a514f83205ea789136f"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    31280095,
		},
	},
}

var Prod_SepoliaToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x2c0b51f491ceaefe8c24c0199fce62d7b040470a"),
		OffRamp:      gethcommon.HexToAddress("0x1d649a11fa14024f9fa2058a6b5b473ea308b688"),
		CommitStore:  gethcommon.HexToAddress("0xc677d898f06cee7b5f6ecbd0f72df5125cebbfc9"),
		PingPongDapp: gethcommon.HexToAddress("0x65b51ba5c9233465f118285e5fb2110c52ad6b27"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3935771,
		},
	},
}

var Prod_PolygonMumbaiToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_PolygonMumbai,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0xa83f2cecb391779b59022eded6ebba0d7ec01f20"),
		OffRamp:      gethcommon.HexToAddress("0xbe582db704bd387222c70ca2e5a027e5e2c06fb7"),
		CommitStore:  gethcommon.HexToAddress("0xf06ff5d2084295909119ca541e93635e7d582ffc"),
		PingPongDapp: gethcommon.HexToAddress("0x044a6b4b561af69d2319a2f4be5ec327a6975d0a"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    37526794,
		},
	},
}

var Prod_SepoliaToPolygonMumbai = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x23438be3256369316e53fd2ef1cd2bdfaf22f6ad"),
		OffRamp:      gethcommon.HexToAddress("0x026fb7c16f1d0082809ff2335715f27e1e074ff6"),
		CommitStore:  gethcommon.HexToAddress("0x290789a55e2e26480f9c04c583d1d5c682aba49a"),
		PingPongDapp: gethcommon.HexToAddress("0xf66fcb898e838a997547ae58fd6882b9bbfdc399"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3819110,
		},
	},
}

var Prod_AvaxFujiToPolygonMumbai = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x762aabc808270fadfdd9e4186739920d68106673"),
		OffRamp:      gethcommon.HexToAddress("0x31cf2040d53f178d168997c658d1a7fc5fa7d215"),
		CommitStore:  gethcommon.HexToAddress("0xa60821b061116054672d102c0b59290910fb51e2"),
		PingPongDapp: gethcommon.HexToAddress("0xcffb4c676a996daefa2a8a6d404a55f59ecc7ce8"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    23754063,
		},
	},
}

var Prod_PolygonMumbaiToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_PolygonMumbai,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x48601a62aa4cb289d006ff4c14023e9d8a7e5a88"),
		OffRamp:      gethcommon.HexToAddress("0xf11e96f85e1038c429d32a877e2225d37cde10e2"),
		CommitStore:  gethcommon.HexToAddress("0x6b6b328cb1467d906389a1bbe54359c56000422c"),
		PingPongDapp: gethcommon.HexToAddress("0x3bd38d308045a39253b502f1789e95c703e27f77"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    37511834,
		},
	},
}

var Prod_OptimismGoerliToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x6bb8d729c35f29df532eb3998ddace336187c84b"),
		OffRamp:      gethcommon.HexToAddress("0xee55842b1d68224d9eef238d4736e851db613630"),
		CommitStore:  gethcommon.HexToAddress("0x4f57b2d4b3b42f09cd7ef48254d2c31b6b525763"),
		PingPongDapp: gethcommon.HexToAddress("0x2af63f50fa3f97f4aa94d28327a759ca86b33bf8"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    12260052,
		},
	},
}

var Prod_ArbitrumGoerliToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_ArbitrumGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress("0x782a7ba95215f2f7c3dd4c153cbb2ae3ec2d3215"),
		OffRamp:      gethcommon.HexToAddress("0xff4b0c64c50d2d7b444cb28699df03ed4bbaf44f"),
		CommitStore:  gethcommon.HexToAddress("0xb69923dfb790e622084b774b99bb45f68904d6a4"),
		PingPongDapp: gethcommon.HexToAddress("0x57c0059fc3f98aa0a5ce4fb5d2882d81d839e74f"),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    31280095,
		},
	},
}
