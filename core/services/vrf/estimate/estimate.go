package estimate

import "math/big"

// GasProofVerification is an upper limit on the gas used for verifying the VRF proof on-chain.
// It can be used to estimate the amount of LINK needed to fulfill a request.
const GasProofVerification uint32 = 200_000

// JuelsNeeded estimates the amount of link needed to fulfill a request
// given the callback gas limit, the gas price, and the wei per unit link.
func JuelsNeeded(callbackGasLimit uint32, maxGasPriceWei, weiPerUnitLink *big.Int) *big.Int {
	maxGasUsed := big.NewInt(int64(callbackGasLimit + GasProofVerification))
	costWei := new(big.Float).SetInt(
		new(big.Int).Set(maxGasUsed).
			Mul(maxGasUsed, maxGasPriceWei),
	)
	costLink := new(big.Float).Set(costWei).Quo(
		costWei,
		new(big.Float).SetInt(weiPerUnitLink),
	)
	costJuelsFloat := new(big.Float).Set(costLink).Mul(
		costLink,
		big.NewFloat(1e18),
	)
	costJuels, _ := costJuelsFloat.Int(nil)
	return costJuels
}
