package mercury_v1

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustDecodeHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

var ChainViewBase = []Block{
	NewBlock(16634362, mustDecodeHex("6f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908172),
	NewBlock(16634361, mustDecodeHex("5f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908170),
}

// ChainView1 vs ChainView2 simulates a re-org based off of a common block 16634362
func MakeChainView1() []Block {
	return append([]Block{
		NewBlock(16634365, mustDecodeHex("9f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908181),
		NewBlock(16634364, mustDecodeHex("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908180),
		NewBlock(16634363, mustDecodeHex("7f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908176),
	}, ChainViewBase...)
}

var ChainView2 = append([]Block{
	NewBlock(16634365, mustDecodeHex("8e30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908180),
	NewBlock(16634364, mustDecodeHex("7e30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908177),
	NewBlock(16634363, mustDecodeHex("6e30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908173),
}, ChainViewBase...)

var ChainView3 = append([]Block{
	NewBlock(16634366, mustDecodeHex("9e30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908181),
	NewBlock(16634365, mustDecodeHex("8e30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908180),
	NewBlock(16634364, mustDecodeHex("7e30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908177),
	NewBlock(16634363, mustDecodeHex("6e30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), 1682908173),
}, ChainViewBase...)

func NewRawPAOS() []parsedAttributedObservation {
	return []parsedAttributedObservation{
		parsedAttributedObservation{
			Timestamp: 1676484822,
			Observer:  commontypes.OracleID(1),

			BenchmarkPrice: big.NewInt(345),
			Bid:            big.NewInt(343),
			Ask:            big.NewInt(347),
			PricesValid:    true,

			LatestBlocks: []Block{},

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		parsedAttributedObservation{
			Timestamp: 1676484826,
			Observer:  commontypes.OracleID(2),

			BenchmarkPrice: big.NewInt(335),
			Bid:            big.NewInt(332),
			Ask:            big.NewInt(336),
			PricesValid:    true,

			CurrentBlockNum:       16634364,
			CurrentBlockHash:      mustDecodeHex("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		parsedAttributedObservation{
			Timestamp: 1676484828,
			Observer:  commontypes.OracleID(3),

			BenchmarkPrice: big.NewInt(347),
			Bid:            big.NewInt(345),
			Ask:            big.NewInt(350),
			PricesValid:    true,

			CurrentBlockNum:       16634365,
			CurrentBlockHash:      mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		parsedAttributedObservation{
			Timestamp: 1676484830,
			Observer:  commontypes.OracleID(4),

			BenchmarkPrice: big.NewInt(346),
			Bid:            big.NewInt(347),
			Ask:            big.NewInt(350),
			PricesValid:    true,

			CurrentBlockNum:       16634365,
			CurrentBlockHash:      mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
	}
}

func NewValidLegacyParsedAttributedObservations() []PAO {
	return []PAO{
		parsedAttributedObservation{
			Timestamp: 1676484822,
			Observer:  commontypes.OracleID(1),

			BenchmarkPrice: big.NewInt(345),
			Bid:            big.NewInt(343),
			Ask:            big.NewInt(347),
			PricesValid:    true,

			CurrentBlockNum:       16634364,
			CurrentBlockHash:      mustDecodeHex("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		parsedAttributedObservation{
			Timestamp: 1676484826,
			Observer:  commontypes.OracleID(2),

			BenchmarkPrice: big.NewInt(335),
			Bid:            big.NewInt(332),
			Ask:            big.NewInt(336),
			PricesValid:    true,

			CurrentBlockNum:       16634364,
			CurrentBlockHash:      mustDecodeHex("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		parsedAttributedObservation{
			Timestamp: 1676484828,
			Observer:  commontypes.OracleID(3),

			BenchmarkPrice: big.NewInt(347),
			Bid:            big.NewInt(345),
			Ask:            big.NewInt(350),
			PricesValid:    true,

			CurrentBlockNum:       16634365,
			CurrentBlockHash:      mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
		parsedAttributedObservation{
			Timestamp: 1676484830,
			Observer:  commontypes.OracleID(4),

			BenchmarkPrice: big.NewInt(346),
			Bid:            big.NewInt(347),
			Ask:            big.NewInt(350),
			PricesValid:    true,

			CurrentBlockNum:       16634365,
			CurrentBlockHash:      mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			CurrentBlockValid:     true,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: true,
		},
	}
}

func NewInvalidParsedAttributedObservations() []PAO {
	return []PAO{
		parsedAttributedObservation{
			Timestamp: 1676484822,
			Observer:  commontypes.OracleID(1),

			BenchmarkPrice: big.NewInt(345),
			Bid:            big.NewInt(343),
			Ask:            big.NewInt(347),
			PricesValid:    false,

			CurrentBlockNum:       16634364,
			CurrentBlockHash:      mustDecodeHex("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			CurrentBlockValid:     false,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: false,
		},
		parsedAttributedObservation{
			Timestamp: 1676484826,
			Observer:  commontypes.OracleID(2),

			BenchmarkPrice: big.NewInt(335),
			Bid:            big.NewInt(332),
			Ask:            big.NewInt(336),
			PricesValid:    false,

			CurrentBlockNum:       16634364,
			CurrentBlockHash:      mustDecodeHex("8f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"),
			CurrentBlockTimestamp: 1682908180,
			CurrentBlockValid:     false,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: false,
		},
		parsedAttributedObservation{
			Timestamp: 1676484828,
			Observer:  commontypes.OracleID(3),

			BenchmarkPrice: big.NewInt(347),
			Bid:            big.NewInt(345),
			Ask:            big.NewInt(350),
			PricesValid:    false,

			CurrentBlockNum:       16634365,
			CurrentBlockHash:      mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			CurrentBlockValid:     false,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: false,
		},
		parsedAttributedObservation{
			Timestamp: 1676484830,
			Observer:  commontypes.OracleID(4),

			BenchmarkPrice: big.NewInt(346),
			Bid:            big.NewInt(347),
			Ask:            big.NewInt(350),
			PricesValid:    false,

			CurrentBlockNum:       16634365,
			CurrentBlockHash:      mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"),
			CurrentBlockTimestamp: 1682591344,
			CurrentBlockValid:     false,

			MaxFinalizedBlockNumber:      16634355,
			MaxFinalizedBlockNumberValid: false,
		},
	}
}

func Test_AggregateFunctions(t *testing.T) {
	f := 1
	invalidPaos := NewInvalidParsedAttributedObservations()
	validLegacyPaos := NewValidLegacyParsedAttributedObservations()

	t.Run("GetConsensusLatestBlock", func(t *testing.T) {
		makePAO := func(blocks []Block) PAO {
			return parsedAttributedObservation{LatestBlocks: blocks}
		}

		makeLegacyPAO := func(num int64, hash string, ts uint64) PAO {
			return parsedAttributedObservation{CurrentBlockNum: num, CurrentBlockHash: mustDecodeHex(hash), CurrentBlockTimestamp: ts, CurrentBlockValid: true}
		}

		t.Run("when all paos are using legacy 'current block'", func(t *testing.T) {
			t.Run("succeeds in the valid case", func(t *testing.T) {
				hash, num, ts, err := GetConsensusLatestBlock(validLegacyPaos, f)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("40044147503a81e9f2a225f4717bf5faf5dc574f69943bdcd305d5ed97504a7e"), hash)
				assert.Equal(t, 16634365, int(num))
				assert.Equal(t, uint64(1682591344), ts)
			})

			t.Run("if invalid, fails", func(t *testing.T) {
				_, _, _, err := GetConsensusLatestBlock(invalidPaos, f)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "cannot come to consensus on latest block number")
			})
			t.Run("if there are not at least f+1 in consensus about hash", func(t *testing.T) {
				_, _, _, err := GetConsensusLatestBlock(validLegacyPaos, 2)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "cannot come to consensus on latest block number")
				_, _, _, err = GetConsensusLatestBlock(validLegacyPaos, 3)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "cannot come to consensus on latest block number")
			})
			t.Run("if there are not at least f+1 in consensus about block number", func(t *testing.T) {
				badPaos := []PAO{
					parsedAttributedObservation{
						CurrentBlockNum:   100,
						CurrentBlockValid: true,
					},
					parsedAttributedObservation{
						CurrentBlockNum:   200,
						CurrentBlockValid: true,
					},
					parsedAttributedObservation{
						CurrentBlockNum:   300,
						CurrentBlockValid: true,
					},
					parsedAttributedObservation{
						CurrentBlockNum:   400,
						CurrentBlockValid: true,
					},
				}
				_, _, _, err := GetConsensusLatestBlock(badPaos, f)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "cannot come to consensus on latest block number")
			})
			t.Run("if there are not at least f+1 in consensus about timestamp", func(t *testing.T) {
				badPaos := []PAO{
					parsedAttributedObservation{
						CurrentBlockTimestamp: 100,
						CurrentBlockValid:     true,
					},
					parsedAttributedObservation{
						CurrentBlockTimestamp: 200,
						CurrentBlockValid:     true,
					},
					parsedAttributedObservation{
						CurrentBlockTimestamp: 300,
						CurrentBlockValid:     true,
					},
					parsedAttributedObservation{
						CurrentBlockTimestamp: 400,
						CurrentBlockValid:     true,
					},
				}
				_, _, _, err := GetConsensusLatestBlock(badPaos, f)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "cannot come to consensus on latest block number")
			})
			t.Run("in the event of an even split for block number/hash, take the higher block number", func(t *testing.T) {
				validFrom := int64(26014056)
				// values below are from a real observed case of this happening in the wild
				paos := []PAO{
					parsedAttributedObservation{
						Timestamp:                    1686759784,
						Observer:                     commontypes.OracleID(2),
						BenchmarkPrice:               big.NewInt(90700),
						Bid:                          big.NewInt(26200),
						Ask:                          big.NewInt(17500),
						PricesValid:                  true,
						CurrentBlockNum:              26014055,
						CurrentBlockHash:             mustDecodeHex("1a2b96ef9a29614c9fc4341a5ca6690ed8ee1a2cd6b232c90ba8bea65a4b93b5"),
						CurrentBlockTimestamp:        1686759784,
						CurrentBlockValid:            true,
						MaxFinalizedBlockNumber:      0,
						MaxFinalizedBlockNumberValid: false,
					},
					parsedAttributedObservation{
						Timestamp:                    1686759784,
						Observer:                     commontypes.OracleID(3),
						BenchmarkPrice:               big.NewInt(92000),
						Bid:                          big.NewInt(21300),
						Ask:                          big.NewInt(74700),
						PricesValid:                  true,
						CurrentBlockNum:              26014056,
						CurrentBlockHash:             mustDecodeHex("bdeb0181416f88812028c4e1ee9e049296c909c1ee15d57cf67d4ce869ed6518"),
						CurrentBlockTimestamp:        1686759784,
						CurrentBlockValid:            true,
						MaxFinalizedBlockNumber:      0,
						MaxFinalizedBlockNumberValid: false,
					},
					parsedAttributedObservation{
						Timestamp:                    1686759784,
						Observer:                     commontypes.OracleID(1),
						BenchmarkPrice:               big.NewInt(67300),
						Bid:                          big.NewInt(70100),
						Ask:                          big.NewInt(83200),
						PricesValid:                  true,
						CurrentBlockNum:              26014056,
						CurrentBlockHash:             mustDecodeHex("bdeb0181416f88812028c4e1ee9e049296c909c1ee15d57cf67d4ce869ed6518"),
						CurrentBlockTimestamp:        1686759784,
						CurrentBlockValid:            true,
						MaxFinalizedBlockNumber:      0,
						MaxFinalizedBlockNumberValid: false,
					},
					parsedAttributedObservation{
						Timestamp:                    1686759784,
						Observer:                     commontypes.OracleID(0),
						BenchmarkPrice:               big.NewInt(8600),
						Bid:                          big.NewInt(89100),
						Ask:                          big.NewInt(53300),
						PricesValid:                  true,
						CurrentBlockNum:              26014055,
						CurrentBlockHash:             mustDecodeHex("1a2b96ef9a29614c9fc4341a5ca6690ed8ee1a2cd6b232c90ba8bea65a4b93b5"),
						CurrentBlockTimestamp:        1686759784,
						CurrentBlockValid:            true,
						MaxFinalizedBlockNumber:      0,
						MaxFinalizedBlockNumberValid: false,
					},
				}
				hash, num, _, err := GetConsensusLatestBlock(paos, f)
				assert.NoError(t, err)
				assert.Equal(t, mustDecodeHex("bdeb0181416f88812028c4e1ee9e049296c909c1ee15d57cf67d4ce869ed6518"), hash)
				assert.Equal(t, int64(26014056), num)
				assert.GreaterOrEqual(t, num, validFrom)
			})
			t.Run("when there are multiple possible blocks meeting > f+1 hashes, takes the hash with the most block numbers in agreement", func(t *testing.T) {
				paos := []PAO{
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(41, "3333333333333333333333333333333333333333333333333333333333333333", 0),
					makeLegacyPAO(41, "3333333333333333333333333333333333333333333333333333333333333333", 0),
					makeLegacyPAO(41, "3333333333333333333333333333333333333333333333333333333333333333", 0),
					makeLegacyPAO(42, "1111111111111111111111111111111111111111111111111111111111111111", 1),
					makeLegacyPAO(42, "1111111111111111111111111111111111111111111111111111111111111111", 1),
					makeLegacyPAO(41, "1111111111111111111111111111111111111111111111111111111111111111", 1),
					makeLegacyPAO(43, "2222222222222222222222222222222222222222222222222222222222222222", 1),
					makeLegacyPAO(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
					makeLegacyPAO(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, f)
				assert.NoError(t, err)
				assert.Equal(t, mustDecodeHex("3333333333333333333333333333333333333333333333333333333333333333"), hash)
				assert.Equal(t, int64(42), num)
				assert.Equal(t, uint64(1), ts)
			})
			t.Run("in the event of an even split of numbers/hashes, takes the hash with the highest block number", func(t *testing.T) {
				paos := []PAO{
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(41, "2222222222222222222222222222222222222222222222222222222222222222", 1),
					makeLegacyPAO(41, "2222222222222222222222222222222222222222222222222222222222222222", 1),
					makeLegacyPAO(41, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, f)
				assert.NoError(t, err)
				assert.Equal(t, mustDecodeHex("3333333333333333333333333333333333333333333333333333333333333333"), hash)
				assert.Equal(t, int64(42), num)
				assert.Equal(t, uint64(1), ts)
			})
			t.Run("in the case where all block numbers are equal but timestamps differ, tie-breaks on latest timestamp", func(t *testing.T) {
				paos := []PAO{
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 2),
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 2),
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 2),
					makeLegacyPAO(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
					makeLegacyPAO(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
					makeLegacyPAO(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, f)
				assert.NoError(t, err)
				assert.Equal(t, mustDecodeHex("3333333333333333333333333333333333333333333333333333333333333333"), hash)
				assert.Equal(t, int64(42), num)
				assert.Equal(t, uint64(2), ts)
			})
			t.Run("in the case where all block numbers and timestamps are equal, tie-breaks by taking the 'lowest' hash", func(t *testing.T) {
				paos := []PAO{
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(42, "3333333333333333333333333333333333333333333333333333333333333333", 1),
					makeLegacyPAO(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
					makeLegacyPAO(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
					makeLegacyPAO(42, "2222222222222222222222222222222222222222222222222222222222222222", 1),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, f)
				assert.NoError(t, err)
				assert.Equal(t, mustDecodeHex("2222222222222222222222222222222222222222222222222222222222222222"), hash)
				assert.Equal(t, int64(42), num)
				assert.Equal(t, uint64(1), ts)
			})
		})

		t.Run("when there is a mix of PAOS, some with legacy 'current block' and some with LatestBlocks", func(t *testing.T) {
			t.Run("succeeds in the valid case where all agree", func(t *testing.T) {
				cv := MakeChainView1()
				paos := []PAO{
					makePAO(cv),
					makePAO(cv),
					makeLegacyPAO(cv[0].Num, hex.EncodeToString([]byte(cv[0].Hash)), cv[0].Ts),
					makeLegacyPAO(cv[0].Num, hex.EncodeToString([]byte(cv[0].Hash)), cv[0].Ts),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, f)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("9f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), hash)
				assert.Equal(t, 16634365, int(num))
				assert.Equal(t, 1682908181, int(ts))
			})

			t.Run("succeeds in the valid case with two different chain views, and returns the highest common block with f+1 observations", func(t *testing.T) {
				cv := MakeChainView1()
				cv2 := ChainView2
				paos := []PAO{
					makePAO(cv[1:]),
					makePAO(cv2),
					makeLegacyPAO(cv[3].Num, hex.EncodeToString([]byte(cv[3].Hash)), cv[3].Ts),
					makeLegacyPAO(cv2[0].Num, hex.EncodeToString([]byte(cv2[0].Hash)), cv2[0].Ts),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, 1)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("8e30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), hash)
				assert.Equal(t, 16634365, int(num))
				assert.Equal(t, 1682908180, int(ts))

				hash, num, ts, err = GetConsensusLatestBlock(paos, 2)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("6f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), hash)
				assert.Equal(t, 16634362, int(num))
				assert.Equal(t, 1682908172, int(ts))
			})
		})

		t.Run("when all PAOS are using LatestBlocks", func(t *testing.T) {
			t.Run("succeeds in the valid case where all agree", func(t *testing.T) {
				paos := []PAO{
					makePAO(MakeChainView1()),
					makePAO(MakeChainView1()),
					makePAO(MakeChainView1()),
					makePAO(MakeChainView1()),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, f)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("9f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), hash)
				assert.Equal(t, 16634365, int(num))
				assert.Equal(t, 1682908181, int(ts))
			})

			t.Run("succeeds in the valid case with two different chain views, and returns the highest common block with f+1 observations", func(t *testing.T) {
				paos := []PAO{
					makePAO(ChainView2),
					makePAO(ChainView2),
					makePAO(MakeChainView1()),
					makePAO(MakeChainView1()),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, 3)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("6f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), hash)
				assert.Equal(t, 16634362, int(num))
				assert.Equal(t, 1682908172, int(ts))
			})

			t.Run("succeeds in the case with many different chain views, and returns the highest common block with f+1 observations", func(t *testing.T) {
				paos := []PAO{
					makePAO(ChainView3),
					makePAO(ChainView2),
					makePAO(MakeChainView1()),
					makePAO(MakeChainView1()),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, 1)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("9f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), hash)
				assert.Equal(t, 16634365, int(num))
				assert.Equal(t, 1682908181, int(ts))
			})

			t.Run("takes highest with at least f+1 when some observations are behind", func(t *testing.T) {
				paos := []PAO{
					makePAO(MakeChainView1()[2:]),
					makePAO(MakeChainView1()[1:]),
					makePAO(MakeChainView1()),
					makePAO(MakeChainView1()),
				}
				hash, num, ts, err := GetConsensusLatestBlock(paos, 3)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("7f30cda279821c5bb6f72f7ab900aa5118215ce59fcf8835b12d0cdbadc9d7b0"), hash)
				assert.Equal(t, 16634363, int(num))
				assert.Equal(t, 1682908176, int(ts))
			})

			t.Run("tie-breaks using smaller hash", func(t *testing.T) {
				cv1 := MakeChainView1()[0:2]
				cv2 := MakeChainView1()[0:1]
				cv3 := MakeChainView1()[0:3]
				cv4 := MakeChainView1()[0:3]

				cv1[0].Hash = string(mustDecodeHex("0000000000000000000000000000000000000000000000000000000000000000"))
				cv4[0].Hash = string(mustDecodeHex("0000000000000000000000000000000000000000000000000000000000000000"))

				paos := []PAO{
					makePAO(cv1),
					makePAO(cv2),
					makePAO(cv3),
					makePAO(cv4),
				}

				hash, num, ts, err := GetConsensusLatestBlock(paos, 1)

				require.NoError(t, err)
				assert.Equal(t, mustDecodeHex("0000000000000000000000000000000000000000000000000000000000000000"), hash)
				assert.Equal(t, 16634365, int(num))
				assert.Equal(t, 1682908181, int(ts))
			})

			t.Run("fails in the case where there is no common block with at least f+1 observations", func(t *testing.T) {
				paos := []PAO{
					makePAO(ChainView2[0:3]),
					makePAO(ChainView2[0:3]),
					makePAO(MakeChainView1()[0:3]),
					makePAO(MakeChainView1()[0:3]),
				}
				_, _, _, err := GetConsensusLatestBlock(paos, 3)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "cannot come to consensus on latest block number")
			})

			t.Run("if invalid, fails", func(t *testing.T) {
				_, _, _, err := GetConsensusLatestBlock(invalidPaos, f)
				require.Error(t, err)
				assert.Contains(t, err.Error(), "cannot come to consensus on latest block number")
			})
		})
	})

	t.Run("GetConsensusMaxFinalizedBlockNum", func(t *testing.T) {
		t.Run("in the valid case", func(t *testing.T) {
			num, err := GetConsensusMaxFinalizedBlockNum(validLegacyPaos, f)

			require.NoError(t, err)
			assert.Equal(t, 16634355, int(num))
		})

		t.Run("errors if there are not at least f+1 valid", func(t *testing.T) {
			_, err := GetConsensusMaxFinalizedBlockNum(invalidPaos, f)
			assert.EqualError(t, err, "fewer than f+1 observations have a valid maxFinalizedBlockNumber (got: 0/4, f=1)")
		})

		t.Run("errors if there are not at least f+1 in consensus about number", func(t *testing.T) {
			badPaos := []PAO{
				parsedAttributedObservation{
					MaxFinalizedBlockNumber:      100,
					MaxFinalizedBlockNumberValid: true,
				},
				parsedAttributedObservation{
					MaxFinalizedBlockNumber:      200,
					MaxFinalizedBlockNumberValid: true,
				},
				parsedAttributedObservation{
					MaxFinalizedBlockNumber:      300,
					MaxFinalizedBlockNumberValid: true,
				},
				parsedAttributedObservation{
					MaxFinalizedBlockNumber:      400,
					MaxFinalizedBlockNumberValid: true,
				},
			}

			_, err := GetConsensusMaxFinalizedBlockNum(badPaos, f)
			assert.EqualError(t, err, "no valid maxFinalizedBlockNumber with at least f+1 votes (got counts: map[100:1 200:1 300:1 400:1], f=1)")
		})
	})
}
