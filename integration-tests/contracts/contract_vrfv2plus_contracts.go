package contracts

import "math/big"

type VRFV2PlusLoadTestMetrics struct {
	RequestCount                         *big.Int
	FulfilmentCount                      *big.Int
	AverageFulfillmentInMillions         *big.Int
	SlowestFulfillment                   *big.Int
	FastestFulfillment                   *big.Int
	ResponseTimesInBlocks                []uint32
	AverageResponseTimeInSecondsMillions *big.Int
	SlowestResponseTimeInSeconds         *big.Int
	FastestResponseTimeInSeconds         *big.Int
}
