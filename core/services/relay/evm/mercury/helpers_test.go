package mercury

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	v1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

var (
	sampleFeedID       = [32]uint8{28, 145, 107, 74, 167, 229, 124, 167, 182, 138, 225, 191, 69, 101, 63, 86, 182, 86, 253, 58, 163, 53, 239, 127, 174, 105, 107, 102, 63, 27, 132, 114}
	sampleClientPubKey = hexutil.MustDecode("0x724ff6eae9e900270edfff233e16322a70ec06e1a6e62a81ef13921f398f6c93")
)

func buildSampleV1Report(p int64) []byte {
	feedID := sampleFeedID
	timestamp := uint32(42)
	bp := big.NewInt(p)
	bid := big.NewInt(243)
	ask := big.NewInt(244)
	currentBlockNumber := uint64(143)
	currentBlockHash := utils.NewHash()
	currentBlockTimestamp := uint64(123)
	validFromBlockNum := uint64(142)

	b, err := v1.ReportTypes.Pack(feedID, timestamp, bp, bid, ask, currentBlockNumber, currentBlockHash, currentBlockTimestamp, validFromBlockNum)
	if err != nil {
		panic(err)
	}
	return b
}

var sampleReports [][]byte

func init() {
	sampleReports = make([][]byte, 4)
	for i := 0; i < len(sampleReports); i++ {
		sampleReports[i] = buildSampleV1Report(int64(i))
	}
}
