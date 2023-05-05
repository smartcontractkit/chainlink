package mercury

import (
	"math/big"

	"github.com/smartcontractkit/libocr/commontypes"
)

func NewParsedAttributedObservations() []ParsedAttributedObservation {
	return []ParsedAttributedObservation{
		ParsedAttributedObservation{
			Timestamp:             1676484822,
			BenchmarkPrice:        big.NewInt(345),
			Bid:                   big.NewInt(343),
			Ask:                   big.NewInt(347),
			CurrentBlockNum:       16634364,
			CurrentBlockHash:      mustDecodeHex("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			ValidFromBlockNum:     16634355,
			Observer:              commontypes.OracleID(1),
		},
		ParsedAttributedObservation{
			Timestamp:             1676484826,
			BenchmarkPrice:        big.NewInt(335),
			Bid:                   big.NewInt(332),
			Ask:                   big.NewInt(336),
			CurrentBlockNum:       16634364,
			CurrentBlockHash:      mustDecodeHex("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			ValidFromBlockNum:     16634355,
			Observer:              commontypes.OracleID(2),
		},
		ParsedAttributedObservation{
			Timestamp:             1676484828,
			BenchmarkPrice:        big.NewInt(347),
			Bid:                   big.NewInt(345),
			Ask:                   big.NewInt(350),
			CurrentBlockNum:       16634365,
			CurrentBlockHash:      mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			ValidFromBlockNum:     16634355,
			Observer:              commontypes.OracleID(3),
		},
		ParsedAttributedObservation{
			Timestamp:             1676484830,
			BenchmarkPrice:        big.NewInt(346),
			Bid:                   big.NewInt(347),
			Ask:                   big.NewInt(350),
			CurrentBlockNum:       16634365,
			CurrentBlockHash:      mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			ValidFromBlockNum:     16634355,
			Observer:              commontypes.OracleID(4),
		},
	}
}
