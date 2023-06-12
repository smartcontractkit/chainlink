package reportcodec

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/smartcontractkit/libocr/commontypes"

	relaymercury "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury"
)

func NewValidParsedAttributedObservations() []relaymercury.ParsedAttributedObservation {
	return []relaymercury.ParsedAttributedObservation{
		relaymercury.ParsedAttributedObservation{
			Timestamp: 1676484822,
			Observer:  commontypes.OracleID(1),

			BenchmarkPrice: big.NewInt(345),
			Bid:            big.NewInt(343),
			Ask:            big.NewInt(347),
			PricesValid:    true,

			CurrentBlockNum:       16634364,
			CurrentBlockHash:      hexutil.MustDecode("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		relaymercury.ParsedAttributedObservation{
			Timestamp: 1676484826,
			Observer:  commontypes.OracleID(2),

			BenchmarkPrice: big.NewInt(335),
			Bid:            big.NewInt(332),
			Ask:            big.NewInt(336),
			PricesValid:    true,

			CurrentBlockNum:       16634364,
			CurrentBlockHash:      hexutil.MustDecode("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		relaymercury.ParsedAttributedObservation{
			Timestamp: 1676484828,
			Observer:  commontypes.OracleID(3),

			BenchmarkPrice: big.NewInt(347),
			Bid:            big.NewInt(345),
			Ask:            big.NewInt(350),
			PricesValid:    true,

			CurrentBlockNum:       16634365,
			CurrentBlockHash:      hexutil.MustDecode("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		relaymercury.ParsedAttributedObservation{
			Timestamp: 1676484830,
			Observer:  commontypes.OracleID(4),

			BenchmarkPrice: big.NewInt(346),
			Bid:            big.NewInt(347),
			Ask:            big.NewInt(350),
			PricesValid:    true,

			CurrentBlockNum:       16634365,
			CurrentBlockHash:      hexutil.MustDecode("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
	}
}

func NewInvalidParsedAttributedObservations() []relaymercury.ParsedAttributedObservation {
	invalidPaos := NewValidParsedAttributedObservations()
	for i := range invalidPaos {
		invalidPaos[i].PricesValid = false
		invalidPaos[i].CurrentBlockValid = false
		invalidPaos[i].MaxFinalizedBlockNumberValid = false
	}
	return invalidPaos
}
