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
	AllowList: []gethcommon.Address{
		// ==============  INTERNAL ==============
		gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"), // deployer key
		gethcommon.HexToAddress("0x9BE566ad50021129F00Ee7219FcEE28490a85656"), // batch testing key
		gethcommon.HexToAddress("0xd54ba5d998479352f375940E5A2A18272714d434"), // batch testing key
		gethcommon.HexToAddress("0x28C70D03e471a2f1D1cad1DC35e7D90AAd2Ac512"), // batch testing key
		gethcommon.HexToAddress("0x5d39fF1Ae4Ab23E3640aa87a5C050483b53b9030"), // batch testing key
		gethcommon.HexToAddress("0x50C38847c059a7c829F7AEee969C652922bd139B"), // batch testing key
		gethcommon.HexToAddress("0x63fc8eE3Dc2326BC17A5E618872C1a4342Bcca09"), // batch testing key
		gethcommon.HexToAddress("0x68f740b79B9abe81628a654f8f733dd4ccE44DFB"), // batch testing key
		gethcommon.HexToAddress("0x0c55B0d8f41E6094a3d0F737c73E892ED0A52D8f"), // batch testing key
		gethcommon.HexToAddress("0x37ffDEe6Dc234E0D1d66571E2c2405aEfd661A6f"), // batch testing key
		gethcommon.HexToAddress("0x450F58153db2289B422e7629Eb4a70cFF77aA72f"), // batch testing key
		// Ping pong
		gethcommon.HexToAddress(""), // SepoliaToAvaxFuji.PingPongDapp,
		gethcommon.HexToAddress(""), // SepoliaToOptimismGoerli.PingPongDapp,
		gethcommon.HexToAddress(""), // SepoliaToArbitrum.PingPongDapp,
		gethcommon.HexToAddress(""), // SepoliaToPolygonMumbai.PingPongDapp,
		// Personal
		gethcommon.HexToAddress("0xEa94AA1318796b5C01a9A37faCBc65423fb2c520"), // Anindita Ghosh
		gethcommon.HexToAddress("0x25D7214ae75F169263921a1cAaf7E6F033210E24"), // Chris Cushman
		gethcommon.HexToAddress("0x498533848239DDc6Bb5Cf7aEF63c97f3f5513ed2"), // Pramod - DApp Sepolia->Fuji
		gethcommon.HexToAddress("0x8e5267453b0aa137Be1Fc976755E6A9bD2a2E029"), // Amine (DevRel) 1
		gethcommon.HexToAddress("0x9d087fC03ae39b088326b67fA3C788236645b717"), // Amine (DevRel) 2
		gethcommon.HexToAddress("0x8fDEA7A82D7861144D027e4eb2acCCf4eB37bb05"), // Andrej Rakic
		gethcommon.HexToAddress("0x208AA722Aca42399eaC5192EE778e4D42f4E5De3"), // Zubin Pratap
		gethcommon.HexToAddress("0x52eE5a881287486573cF5CB5e7E7D92F30b03014"), // Zubin Pratap
		gethcommon.HexToAddress("0x44794725885F23cf36deE43554Ad204fb634A057"), // Frank Kong
		gethcommon.HexToAddress("0x0F1aF0A5d727b53dA44bBDE16843B3BA7F98Af68"), // Frank Kong
		gethcommon.HexToAddress("0x5803a251D118899dF7B403769e72532dBE854712"), // Amine El Manaa
		gethcommon.HexToAddress("0xcD936a39336a2E2c5a011137E46c8120dcaE0d65"), // Internal devrel proxy

		// ==============  EXTERNAL ==============
		gethcommon.HexToAddress("0xd65113b9B1EeD81113EaF41DC0D2d34fCa31522C"), // BetaUser - Multimedia
		gethcommon.HexToAddress("0x217F4Eb693C54cA36Cfd80DA4DAAE6f7A5535e9C"), // BetaUser - Cozy Labs
		gethcommon.HexToAddress("0xB22107572f5A5352dDC1B4fc9630083FBfAE2022"), // BetaUser - Cozy Labs
		gethcommon.HexToAddress("0xB0AC8F6AF9712CF369934A811A79550DA046Fc51"), // BetaUser - InsurAce
		gethcommon.HexToAddress("0x244d07fe4DFa30b4EE376751FDC793aE844c5dE6"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0x8264AcEE321ac02549aff7fA05A4Ae7a2e92A6f1"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0x012a3fda37649945Cc72D725168FcB57A469bA6A"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0x552acA1343A6383aF32ce1B7c7B1b47959F7ad90"), // BetaUser - Sommelier Finance
		gethcommon.HexToAddress("0x8e0866aacCF880E45249e932a094c821Ef4dE5f7"), // BetaUser - OpenZeppelin
		gethcommon.HexToAddress("0x9bf889acd6dd651bd897b6ff7a6ecde84a4b29aa"), // BetaUser - ANZ
		gethcommon.HexToAddress("0x9E945BB44B7E264c579e7f0c1FC28FBb39a32386"), // BetaUser - ANZ
		gethcommon.HexToAddress("0x309bdb4F7608584653D1bE804E8420fA0302911b"), // BetaUser - ANZ
		gethcommon.HexToAddress("0x066AFe67f2762C4009637c5ac10C789738cc7488"), // BetaUser - Tristero
		gethcommon.HexToAddress("0x6d818effaE3B40a89AEEb0e0FbA1827EFf77e0E1"), // BetaUser - Tristero
		gethcommon.HexToAddress("0x1C4310602DEFc04117980080b1807eac15687649"), // BetaUser - Zaros (ZD Labs)
		gethcommon.HexToAddress("0x4d2F1C99BCE324B9Ba486d704A0235A754D188a2"), // BetaUser - Aave (BGD Labs)
		gethcommon.HexToAddress("0x289F4D1e83BE7bb8A493D55622cE09D72D2A16e6"), // BetaUser - Steadefi
		gethcommon.HexToAddress("0x651c84ACc85D7a4506FD5dd6EB94d050c7ED2fe7"), // BetaUser - Lendvest
		gethcommon.HexToAddress("0xf62FD6119EBAFEdAAa7a75C1713Bca98729f163D"), // BetaUser - Fidelity Digital Assets
		gethcommon.HexToAddress("0x0D7a3a17E2E160287D3e7e74c4A1B22422156642"), // BetaUser - RiseWorks
		gethcommon.HexToAddress("0xc5f502Ae5972c938940b33308f8845cbe80211B5"), // BetaUser - Robolabs
		gethcommon.HexToAddress("0x87F45de79da4c3356591d74619693E372D525F1b"), // BankToken 1 (BANK) - internal testing contract for SWIFT POC
		gethcommon.HexToAddress("0x784c400D6fF625051d2f587dC0276E3A1ffD9cda"), // BankToken 2 (BANK) - internal testing contract for SWIFT POC
		gethcommon.HexToAddress("0xF92E4b278380f39fADc24483C7baC61b73EE93F2"), // BetaUser - SWIFT (BondToken)
		gethcommon.HexToAddress("0xAa6f663a14b8dA1EB9CF021379f4Ba6BF536268A"), // BetaUser - Fidelity Digital Assets
		gethcommon.HexToAddress("0xB781A9EFC6bd4Cf0dbE547D20151A405673F4CDe"), // BetaUser - RiseWorks
		gethcommon.HexToAddress("0xe764C455e3Bd05Eb7Cf53Ec8491dca0e91486D24"), // BetaUser - Synthetix v3 core
		gethcommon.HexToAddress("0x8e52262f91ef7049adfD8d1E608172fAC57995c3"), // BetaUser - Synthetix v3 core
		gethcommon.HexToAddress("0x2A45BaE1E58AaD3261af187b7dAde90889c039Dc"), // BetaUser - Synthetix v3 core
		gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"), // BetaUser - Alongside
		gethcommon.HexToAddress("0x0819BBae96c2C0F15477D212e063303221Cf24b9"), // BetaUser - Oddz
		gethcommon.HexToAddress("0x38104E1bB27A06306B72162047F585B3e6D27484"), // BetaUser - Oddz
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x779877A7B0D9E8603169DdbD7836e478b4624789"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x097D90c9d3E0B50Ca60e1ae45F6A81010f9FB534"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.CACHEGOLD: {
			Token:         gethcommon.HexToAddress("0x997BCCAE553112CD023592691d41687a3f1EfA7C"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.CACHEGOLD.Price(),
			Decimals:      rhea.CACHEGOLD.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.ANZ: {
			Token:         gethcommon.HexToAddress("0x92eA346B7a2AaB84e6AaB03b80E2421eeFB04685"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.ANZ.Price(),
			Decimals:      rhea.ANZ.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.InsurAce: {
			Token:         gethcommon.HexToAddress("0xb7c8bCA891143221a34DB60A26639785C4839040"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.InsurAce.Price(),
			Decimals:      rhea.InsurAce.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		//rhea.ZUSD: {
		//	Token:         gethcommon.HexToAddress("0x09ae935D80E190403C61Cc5d854Fbf6a7b4a559a"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.ZUSD.Price(),
		//	Decimals:      rhea.ZUSD.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
		rhea.STEADY: {
			Token:         gethcommon.HexToAddress("0x82abB1864326A8A7e1A357FFA2270D09CCb867B9"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.STEADY.Price(),
			Decimals:      rhea.STEADY.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		//rhea.SUPER: {
		//	Token:         gethcommon.HexToAddress("0xCb4B3f72B5b6D0b7072aFDDf18FE61A0d569EC39"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.SUPER.Price(),
		//	Decimals:      rhea.SUPER.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
		rhea.BankToken: {
			Token:         gethcommon.HexToAddress("0x784c400D6fF625051d2f587dC0276E3A1ffD9cda"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.BankToken.Price(),
			Decimals:      rhea.BankToken.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.BondToken: {
			Token:         gethcommon.HexToAddress("0xF92E4b278380f39fADc24483C7baC61b73EE93F2"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.BondToken.Price(),
			Decimals:      rhea.BondToken.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.SNXUSD: {
			Token:         gethcommon.HexToAddress("0x585d8E269A250aCBf7D4884A1a31D3b596B46D8B"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.SNXUSD.Price(),
			Decimals:      rhea.SNXUSD.Decimals(),
			TokenPoolType: rhea.BurnMint,
			PoolAllowList: []gethcommon.Address{
				gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"),
				gethcommon.HexToAddress("0x2A45BaE1E58AaD3261af187b7dAde90889c039Dc"),
			},
		},
		//rhea.FUGAZIUSDC: {
		//	Token:         gethcommon.HexToAddress("0x832bA6abcAdC68812be372F4ef20aAC268bA20B7"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.FUGAZIUSDC.Price(),
		//	Decimals:      rhea.FUGAZIUSDC.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
		rhea.Alongside: {
			Token:         gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.Alongside.Price(),
			Decimals:      rhea.Alongside.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress(""),
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
	AllowList: []gethcommon.Address{
		// ==============  INTERNAL ==============
		gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"), // deployer key
		// Ping pong
		gethcommon.HexToAddress(""), // OptimismGoerliToAvaxFuji.PingPongDapp,
		gethcommon.HexToAddress(""), // OptimismGoerliToSepolia.PingPongDapp,
		gethcommon.HexToAddress(""), // OptimismGoerliToArbitrumGoerli.PingPongDapp,
		// Personal
		gethcommon.HexToAddress("0xEa94AA1318796b5C01a9A37faCBc65423fb2c520"), // Anindita Ghosh
		gethcommon.HexToAddress("0x8fDEA7A82D7861144D027e4eb2acCCf4eB37bb05"), // Andrej Rakic
		gethcommon.HexToAddress("0x208AA722Aca42399eaC5192EE778e4D42f4E5De3"), // Zubin Pratap
		gethcommon.HexToAddress("0x52eE5a881287486573cF5CB5e7E7D92F30b03014"), // Zubin Pratap
		gethcommon.HexToAddress("0xabFD23063251A6481D65e8244237996d3D4d7b59"), // Internal devrel proxy

		// ==============  EXTERNAL ==============
		gethcommon.HexToAddress("0x3FcFF7d9f88C64905e2cD9960c7452b5E6690E13"), // BetaUser - AAVE
		gethcommon.HexToAddress("0x1b5D803Be089e43110Faf54c6b4eC40409Cc7450"), // BetaUser - Multimedia
		gethcommon.HexToAddress("0xE8Cc2Bd6082387a7AC749176b1Fe19377f420740"), // BetaUser - Multimedia (AA wallet)
		gethcommon.HexToAddress("0x244d07fe4DFa30b4EE376751FDC793aE844c5dE6"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0x8264AcEE321ac02549aff7fA05A4Ae7a2e92A6f1"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0x012a3fda37649945Cc72D725168FcB57A469bA6A"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0xF7726C9F7D2a9433CF8E46640821bebAbbE020b3"), // BetaUser - Zaros (ZD Labs)
		gethcommon.HexToAddress("0xF640cEA278E94708c358D79e5872AFda56010117"), // BetaUser - Aave (BGD Labs)
		gethcommon.HexToAddress("0x69D235A7E01aBdf463D7d886492229b75A4F1BC6"), // BetaUser - Steadefi
		gethcommon.HexToAddress("0xDdcE30979147091F26513C495EEE1bfa6C0a6730"), // BetaUser - RiseWorks
		gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"), // BetaUser - Alongside
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0xdc2CC710e42857672E7907CF474a69B63B93089f"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x4200000000000000000000000000000000000006"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		//rhea.CACHEGOLD: {
		//	Token:         gethcommon.HexToAddress("0xa6446C6f492f31A33bC68249ae59F8871123a777"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.CACHEGOLD.Price(),
		//	Decimals:      rhea.CACHEGOLD.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
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
		rhea.Alongside: {
			Token:         gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.Alongside.Price(),
			Decimals:      rhea.Alongside.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress(""),
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
	AllowList: []gethcommon.Address{
		// ==============  INTERNAL ==============
		gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"), // deployer key
		gethcommon.HexToAddress("0xEa94AA1318796b5C01a9A37faCBc65423fb2c520"), // Test Script 0xEa94AA1318796b5C01a9A37faCBc65423fb2c520
		// Ping pong
		gethcommon.HexToAddress(""), // AvaxFujiToSepolia.PingPongDapp,
		gethcommon.HexToAddress(""), // AvaxFujiToOptimismGoerli.PingPongDapp,
		gethcommon.HexToAddress(""), // AvaxFujiToPolygonMumbai.PingPongDapp,
		// Personal
		gethcommon.HexToAddress("0xEa94AA1318796b5C01a9A37faCBc65423fb2c520"), // Anindita Ghosh
		gethcommon.HexToAddress("0x594D8E57D8801069C77AAB90222a9162E908AA63"), // Pramod - Dapp Fuji->OptimismGoerli
		gethcommon.HexToAddress("0xFE5394A63433A3975b1936dEc92DAa161FEE7463"), // Pramod - DApp Fuji->Sepolia
		gethcommon.HexToAddress("0x912519a7E5e2e2309b1e60F540683c6661757A0C"), // Amine (DevRel) 1
		gethcommon.HexToAddress("0x9d087fC03ae39b088326b67fA3C788236645b717"), // Amine (DevRel) 2
		gethcommon.HexToAddress("0x8fDEA7A82D7861144D027e4eb2acCCf4eB37bb05"), // Andrej Rakic
		gethcommon.HexToAddress("0x208AA722Aca42399eaC5192EE778e4D42f4E5De3"), // Zubin Pratap
		gethcommon.HexToAddress("0x52eE5a881287486573cF5CB5e7E7D92F30b03014"), // Zubin Pratap
		gethcommon.HexToAddress("0x00104e54E037453daE202d6a92A3f75B2fdC2737"), // Amine El Manaa
		gethcommon.HexToAddress("0x447Fd5eC2D383091C22B8549cb231a3bAD6d3fAf"), // Internal devrel proxy

		// ==============  EXTERNAL ==============
		gethcommon.HexToAddress("0x1b5D803Be089e43110Faf54c6b4eC40409Cc7450"), // BetaUser - Multimedia
		gethcommon.HexToAddress("0xE8Cc2Bd6082387a7AC749176b1Fe19377f420740"), // BetaUser - Multimedia (AA wallet)
		gethcommon.HexToAddress("0xa78ceF54da82D6279b20457F4D46294AfF59C871"), // BetaUser - Flash Liquidity
		gethcommon.HexToAddress("0x6613fd61bbfEF3291f2D7C7203Ceab212e880DbB"), // BetaUser - Flash Liquidity
		gethcommon.HexToAddress("0xa294275E5Bb4A786a3305f4276645290cCC7419B"), // BetaUser - Flash Liquidity
		gethcommon.HexToAddress("0xcA218DCFD26990223a2eDA70f3A568eaae22c051"), // BetaUser - Cozy Labs
		gethcommon.HexToAddress("0xD0fB066847d5DBc760E9575f79d9A044385e4079"), // BetaUser - Cozy Labs
		gethcommon.HexToAddress("0xD93C3Ae0949f905846FdfFc2b5b8A0a047dda59f"), // BetaUser - InsurAce
		gethcommon.HexToAddress("0x244d07fe4DFa30b4EE376751FDC793aE844c5dE6"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0x8264AcEE321ac02549aff7fA05A4Ae7a2e92A6f1"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0x012a3fda37649945Cc72D725168FcB57A469bA6A"), // BetaUser - CACHE.gold
		gethcommon.HexToAddress("0x1b38148B8DfdeA0B3D80C45F0d8569889504f0B5"), // BetaUser - Sommelier Finance
		gethcommon.HexToAddress("0xe0534662Ff1182a1C32E400d2b64723817344Ab4"), // BetaUser - Sommelier Finance
		gethcommon.HexToAddress("0x4986fD36b6b16f49b43282Ee2e24C5cF90ed166d"), // BetaUser - Sommelier Finance
		gethcommon.HexToAddress("0xc7a5d29248cf53b094106ca1d29634b34ad0fede"), // BetaUser - Tristero
		gethcommon.HexToAddress("0x4A5D71F7027684d473a1110a412B510354aF33e7"), // BetaUser - Aave (BGD Labs)
		gethcommon.HexToAddress("0x44eb6D97e98CE35eEFBD5764aa786f10121bC5e4"), // BetaUser - ANZ
		gethcommon.HexToAddress("0xa707480A11f12569b888306F2F118716d3BC29A1"), // BetaUser - Lendvest
		gethcommon.HexToAddress("0xbcFA8eAB1fCe576F1Ef71772E46519e0ADC06623"), // BetaUser - Lendvest
		gethcommon.HexToAddress("0xd35468ab2547a5ba9c9b809e67a35bcc5b89d2fe"), // BetaUser - Lendvest
		gethcommon.HexToAddress("0x9344AeA9b3270d51c9603d3054E421386dFaacB8"), // BetaUser - Fidelity Digital Assets
		gethcommon.HexToAddress("0x89Eccc61B2d35eACCe08284CF22c2D6487B80A3A"), // BetaUser - Robolabs
		gethcommon.HexToAddress("0xAa6f663a14b8dA1EB9CF021379f4Ba6BF536268A"), // BetaUser - Fidelity Digital Assets
		gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"), // BetaUser - Alongside
		gethcommon.HexToAddress("0xA9cb37191A089C8f8c24fC3e1F2f761De93FA827"), // MintDAO - Sender (EOA)
		gethcommon.HexToAddress("0xaFE336062eD69c108232c303fBa9b2b1c709fd9d"), // MintDAO - Sender (Proxy)
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x0b9d5D9136855f6FEc3c0993feE6E9CE8a297846"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WAVAX: {
			Token:         gethcommon.HexToAddress("0xd00ae08403B9bbb9124bB305C09058E32C39A48c"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.WAVAX.Price(),
			Decimals:      rhea.WAVAX.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		//rhea.CACHEGOLD: {
		//	Token:         gethcommon.HexToAddress("0xD16eD805F3eCe986d9541afaD3E59De2F3732517"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.CACHEGOLD.Price(),
		//	Decimals:      rhea.CACHEGOLD.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
		rhea.ANZ: {
			Token:         gethcommon.HexToAddress("0xe3d06cb8eac016749281f45e779ac2976baa02ed"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.ANZ.Price(),
			Decimals:      rhea.ANZ.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		rhea.InsurAce: {
			Token:         gethcommon.HexToAddress("0xda305ab72858939758d5a711494cd447d2d8842e"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.InsurAce.Price(),
			Decimals:      rhea.InsurAce.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		//rhea.SUPER: {
		//	Token:         gethcommon.HexToAddress("0xCb4B3f72B5b6D0b7072aFDDf18FE61A0d569EC39"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.SUPER.Price(),
		//	Decimals:      rhea.SUPER.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
		rhea.BankToken: {
			Token:         gethcommon.HexToAddress("0x7130aac4827a8b085ffe701a7d4749e2b452a837"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.BankToken.Price(),
			Decimals:      rhea.BankToken.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		rhea.BondToken: {
			Token:         gethcommon.HexToAddress("0x56e01ecb119c45ff14248f6ebc27c05d4a72d4f9"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.BondToken.Price(),
			Decimals:      rhea.BondToken.Decimals(),
			TokenPoolType: rhea.Wrapped,
		},
		//rhea.FUGAZIUSDC: {
		//	Token:         gethcommon.HexToAddress("0x150a0ee7393294442EE4d4F5C7d637af01dF93ee"),
		//	Pool:          gethcommon.HexToAddress(""),
		//	Price:         rhea.FUGAZIUSDC.Price(),
		//	Decimals:      rhea.FUGAZIUSDC.Decimals(),
		//	TokenPoolType: rhea.Legacy,
		//},
		rhea.Alongside: {
			Token:         gethcommon.HexToAddress("0xB3c3977B0aC329A9035889929482a4c635B50573"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.Alongside.Price(),
			Decimals:      rhea.Alongside.Decimals(),
			TokenPoolType: rhea.BurnMint,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WAVAX},
	WrappedNative: rhea.WAVAX,
	Router:        gethcommon.HexToAddress(""),
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
	AllowList: []gethcommon.Address{
		// ==============  INTERNAL ==============
		gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"), // deployer key
		// Ping pong
		gethcommon.HexToAddress(""), // ArbitrumGoerliToSepolia.PingPongDapp,
		gethcommon.HexToAddress(""), // ArbitrumGoerliToOptimismGoerli.PingPongDapp,
		// Personal
		gethcommon.HexToAddress("0x8fDEA7A82D7861144D027e4eb2acCCf4eB37bb05"), // Andrej Rakic
		gethcommon.HexToAddress("0x208AA722Aca42399eaC5192EE778e4D42f4E5De3"), // Zubin Pratap
		gethcommon.HexToAddress("0x52eE5a881287486573cF5CB5e7E7D92F30b03014"), // Zubin Pratap
		gethcommon.HexToAddress("0x82FAB72c5Baf6f15f89540EfBb7A62Cb410c300C"), // Internal devrel proxy

		// ==============  EXTERNAL ==============
		gethcommon.HexToAddress("0xF5022eDd1B827E6EA4bBdb961212ECD7F315ed88"), // BetaUser - RiseWorks
		gethcommon.HexToAddress("0x0D7a3a17E2E160287D3e7e74c4A1B22422156642"), // BetaUser - RiseWorks
		gethcommon.HexToAddress("0x63e430dBd88C1bBFBc97336b4357Aa5Aea83367e"), // BetaUser - RiseWorks
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0xd14838A68E8AFBAdE5efb411d5871ea0011AFd28"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WETH: {
			Token:         gethcommon.HexToAddress("0x32d5D5978905d9c6c2D4C417F0E06Fe768a4FB5a"),
			Price:         rhea.WETH.Price(),
			Decimals:      rhea.WETH.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WETH},
	WrappedNative: rhea.WETH,
	Router:        gethcommon.HexToAddress(""),
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
	AllowList: []gethcommon.Address{
		// ==============  INTERNAL ==============
		gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"), // deployer key
		gethcommon.HexToAddress("0xEa94AA1318796b5C01a9A37faCBc65423fb2c520"), // Test Script 0xEa94AA1318796b5C01a9A37faCBc65423fb2c520
		// Ping pong
		gethcommon.HexToAddress(""), // PolygonMumbaiToSepolia.PingPongDapp,
		gethcommon.HexToAddress(""), // PolygonMumbaiToAvax.PingPongDapp,
		// Personal
		gethcommon.HexToAddress("0xEa94AA1318796b5C01a9A37faCBc65423fb2c520"), // Anindita Ghosh
		gethcommon.HexToAddress("0x8fDEA7A82D7861144D027e4eb2acCCf4eB37bb05"), // Andrej Rakic
		gethcommon.HexToAddress("0x208AA722Aca42399eaC5192EE778e4D42f4E5De3"), // Zubin Pratap
		gethcommon.HexToAddress("0x52eE5a881287486573cF5CB5e7E7D92F30b03014"), // Zubin Pratap
		gethcommon.HexToAddress("0x44794725885F23cf36deE43554Ad204fb634A057"), // Frank Kong
		gethcommon.HexToAddress("0x0F1aF0A5d727b53dA44bBDE16843B3BA7F98Af68"), // Frank Kong
		gethcommon.HexToAddress("0xA4285EC042b198aeb0C68679c94a615c4d82DAd0"), // Amine El Manaa
		gethcommon.HexToAddress("0xBdcc3f1D0B4c78F1fe03C91D9498f3DAEeE6948B"), // Internal devrel proxy

		// ==============  EXTERNAL ==============
		gethcommon.HexToAddress("0xe764C455e3Bd05Eb7Cf53Ec8491dca0e91486D24"), // BetaUser - Synthetix v3 core
		gethcommon.HexToAddress("0x8e52262f91ef7049adfD8d1E608172fAC57995c3"), // BetaUser - Synthetix v3 core
		gethcommon.HexToAddress("0x6De1e981d2137f7839840e2140dBB3A05F05B770"), // BetaUser - Flash Liquidity
		gethcommon.HexToAddress("0x2A45BaE1E58AaD3261af187b7dAde90889c039Dc"), // BetaUser - Synthetix v3 core
		gethcommon.HexToAddress("0x0819BBae96c2C0F15477D212e063303221Cf24b9"), // BetaUser - Oddz
		gethcommon.HexToAddress("0x38104E1bB27A06306B72162047F585B3e6D27484"), // BetaUser - Oddz
		gethcommon.HexToAddress("0xA9cb37191A089C8f8c24fC3e1F2f761De93FA827"), // MintDAO - Sender (EOA)
		gethcommon.HexToAddress("0xaFE336062eD69c108232c303fBa9b2b1c709fd9d"), // MintDAO - Sender (Proxy)
	},
	SupportedTokens: map[rhea.Token]rhea.EVMBridgedToken{
		rhea.LINK: {
			Token:         gethcommon.HexToAddress("0x326C977E6efc84E512bB9C30f76E30c160eD06FB"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.LINK.Price(),
			Decimals:      rhea.LINK.Decimals(),
			TokenPoolType: rhea.LockRelease,
		},
		rhea.WMATIC: {
			Token:         gethcommon.HexToAddress("0x9c3C9283D3e44854697Cd22D3Faa240Cfb032889"),
			Price:         rhea.WMATIC.Price(),
			Decimals:      rhea.WMATIC.Decimals(),
			TokenPoolType: rhea.FeeTokenOnly,
		},
		rhea.SNXUSD: {
			Token:         gethcommon.HexToAddress("0x585d8E269A250aCBf7D4884A1a31D3b596B46D8B"),
			Pool:          gethcommon.HexToAddress(""),
			Price:         rhea.SNXUSD.Price(),
			Decimals:      rhea.SNXUSD.Decimals(),
			TokenPoolType: rhea.BurnMint,
			PoolAllowList: []gethcommon.Address{
				gethcommon.HexToAddress("0xda9e8e71bb750a996af33ebb8abb18cd9eb9dc75"),
				gethcommon.HexToAddress("0x2A45BaE1E58AaD3261af187b7dAde90889c039Dc"),
			},
		},
	},
	FeeTokens:     []rhea.Token{rhea.LINK, rhea.WMATIC},
	WrappedNative: rhea.WMATIC,
	Router:        gethcommon.HexToAddress(""),
	ARM:           gethcommon.HexToAddress(""),
	PriceRegistry: gethcommon.HexToAddress(""),
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
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3655211,
		},
	},
}

var Prod_AvaxFujiToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    22876561,
		},
	},
}

var Prod_SepoliaToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3655106,
		},
	},
}

var Prod_OptimismGoerliToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    10439785,
		},
	},
}

var Prod_OptimismGoerliToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    10439147,
		},
	},
}

var Prod_AvaxFujiToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    22875290,
		},
	},
}

var Prod_ArbitrumGoerliToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_ArbitrumGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    25022239,
		},
	},
}

var Prod_SepoliaToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3654960,
		},
	},
}

var Prod_PolygonMumbaiToSepolia = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_PolygonMumbai,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    36631143,
		},
	},
}

var Prod_SepoliaToPolygonMumbai = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_Sepolia,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    3654802,
		},
	},
}

var Prod_AvaxFujiToPolygonMumbai = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_AvaxFuji,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    22874574,
		},
	},
}

var Prod_PolygonMumbaiToAvaxFuji = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_PolygonMumbai,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    36631628,
		},
	},
}

var Prod_OptimismGoerliToArbitrumGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_OptimismGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    0,
		},
	},
}

var Prod_ArbitrumGoerliToOptimismGoerli = rhea.EvmDeploymentConfig{
	ChainConfig: Prod_ArbitrumGoerli,
	LaneConfig: rhea.EVMLaneConfig{
		OnRamp:       gethcommon.HexToAddress(""),
		OffRamp:      gethcommon.HexToAddress(""),
		CommitStore:  gethcommon.HexToAddress(""),
		PingPongDapp: gethcommon.HexToAddress(""),
		DeploySettings: rhea.LaneDeploySettings{
			DeployLane:         false,
			DeployPingPongDapp: false,
			DeployedAtBlock:    0,
		},
	},
}
