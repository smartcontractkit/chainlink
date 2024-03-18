package test

import "math/big"

var commonGasPriceEstimator = commonStaticGasPriceEstimator{
	commonStaticGasPriceEstimatorConfig: commonStaticGasPriceEstimatorConfig{
		getGasPriceResponse: big.NewInt(7),

		denoteInUSDRequest: denoteInUSDRequest{
			p:                  big.NewInt(8),
			wrappedNativePrice: big.NewInt(9),
		},
		denoteInUSDResponse: denoteInUSDResponse{
			result: big.NewInt(10),
		},

		medianRequest: medianRequest{
			gasPrices: []*big.Int{big.NewInt(11), big.NewInt(13), big.NewInt(17)},
		},
		medianResponse: big.NewInt(13),
	},
}

type commonStaticGasPriceEstimator struct {
	commonStaticGasPriceEstimatorConfig
}

type commonStaticGasPriceEstimatorConfig struct {
	getGasPriceResponse *big.Int

	denoteInUSDRequest  denoteInUSDRequest
	denoteInUSDResponse denoteInUSDResponse

	medianRequest
	medianResponse *big.Int
}

type denoteInUSDRequest struct {
	p                  *big.Int
	wrappedNativePrice *big.Int
}

type denoteInUSDResponse struct {
	result *big.Int
}

type medianRequest struct {
	gasPrices []*big.Int
}

type deviatesRequest struct {
	p1 *big.Int
	p2 *big.Int
}
